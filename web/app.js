const API_BASE = window.location.origin;
const APP_NAME = "sketchnote-artist";
const USER_ID = "web-guest";

async function generateSketchnote() {
  const urlInput = document.getElementById('youtubeUrl');
  const generateBtn = document.getElementById('generateBtn');

  const idleState = document.getElementById('idleState');
  const progressSection = document.getElementById('progressSection');
  const resultSection = document.getElementById('resultSection');
  const resultImage = document.getElementById('resultImage');
  const statusMsg = document.getElementById('statusMessage');

  const videoUrl = urlInput.value.trim();

  if (!videoUrl) {
    alert("Please enter a valid YouTube URL");
    return;
  }

  // Reset UI State
  idleState.classList.add('hidden');
  resultSection.classList.add('hidden');
  progressSection.classList.remove('hidden');
  statusMsg.innerText = "Initializing...";
  statusMsg.style.color = "var(--text-secondary)";
  generateBtn.disabled = true;

  // Track duration for debugging timeouts
  const startTime = Date.now();

  try {
    // 1. Create Session
    const sessionResp = await fetch(`/apps/${APP_NAME}/users/${USER_ID}/sessions`, {
      method: 'POST'
    });

    if (!sessionResp.ok) throw new Error("Failed to create session");
    const sessionData = await sessionResp.json();
    const sessionId = sessionData.id;

    // 2. Start Streaming Request
    statusMsg.innerText = "Connecting to agent...";

    const response = await fetch('/stream-run', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        appName: APP_NAME,
        userId: USER_ID,
        sessionId: sessionId,
        newMessage: {
          role: "user",
          parts: [{ text: videoUrl }]
        }
      })
    });

    if (!response.ok) throw new Error(`Stream failed: ${response.statusText}`);

    // 3. Read the Stream
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n\n');
      buffer = lines.pop(); // Keep incomplete chunk

      for (const line of lines) {
        if (line.startsWith(':')) continue; // Ignore comments/heartbeats
        if (line.startsWith('event: done')) return; // Success!
        if (line.startsWith('event: error')) {
          // Locate data extraction more robustly
          const jsonStart = line.indexOf('{');
          if (jsonStart !== -1) {
            const data = JSON.parse(line.substring(jsonStart));
            throw new Error(data.error || "Unknown stream error");
          }
        }

        // Handle explicit status events from observability
        if (line.startsWith('event: status')) {
          const jsonStart = line.indexOf('{');
          if (jsonStart !== -1) {
            try {
              const data = JSON.parse(line.substring(jsonStart));
              if (data.message) {
                statusMsg.innerText = data.message;
              }
            } catch (e) {
              console.warn("Failed to parse status event", e);
            }
          }
          continue; // Skip trying to parse as 'data: '
        }

        if (line.startsWith('data: ')) {
          const jsonStr = line.substring(6);
          try {
            const event = JSON.parse(jsonStr);
            handleAgentEvent(event, statusMsg, resultImage, progressSection, resultSection);
          } catch (e) {
            console.warn("Failed to parse event", e);
          }
        }
      }
    }

  } catch (error) {
    const elapsed = ((Date.now() - startTime) / 1000).toFixed(1);
    console.error(error);

    // Check for "Failed to fetch" which usually means network error/timeout
    if (error.message === "Failed to fetch") {
      statusMsg.innerText = `Network Error: Failed to reach server (${elapsed}s). Check connection or VPN.`;
    } else {
      statusMsg.innerText = `Error: ${error.message} (${elapsed}s)`;
    }

    statusMsg.style.color = "#FF0000";
    progressSection.classList.add('hidden');
    idleState.classList.remove('hidden');
  } finally {
    generateBtn.disabled = false;
  }
}

function handleAgentEvent(event, statusMsg, resultImage, progressSection, resultSection) {
  // Handle different event types from ADK
  // This depends on the exact JSON structure of *session.Event

  // Example heuristics based on typical ADK event structure:
  // 1. Model Call (Thinking)
  if (event.Recall || (event.Models && event.Models.length > 0)) {
    statusMsg.innerText = "Curator is analyzing video...";
  }

  // 2. Tool Call (Summarizing)
  if (event.ToolCall) {
    const toolName = event.ToolCall.Name;
    if (toolName.includes('summarize')) {
      statusMsg.innerText = "Summarizing video content...";
    } else if (toolName.includes('generate_image')) {
      statusMsg.innerText = "Artist is sketching...";
    }
  }

  // 3. Model Response (Final Output)
  if (event.Content && event.Content.Parts) {
    for (const part of event.Content.Parts) {
      if (part.Text && part.Text.includes('.png')) {
        // Found image path!
        const filename = part.Text.trim();
        // Assuming filename matches locally served pattern
        // Safety check: is it a path or just name?
        // Extract basename if needed, but tool usually returns basics

        // Ensure we check for the actual filename
        // The text might be "Image saved to: foo.png"
        // We need to extract the .png part.
        const match = filename.match(/[\w-]+\.png/);
        if (match) {
          const cleanFilename = match[0];
          const imagePath = `/images/${cleanFilename}`;

          // Preload image to ensure it's ready
          const img = new Image();
          img.onload = () => {
            resultImage.src = imagePath;
            progressSection.classList.add('hidden');
            resultSection.classList.remove('hidden');
            statusMsg.innerText = "Sketchnote created successfully!";
          };
          img.onerror = () => {
            statusMsg.innerText = "Image generated but failed to load.";
          };
          img.src = imagePath;
        }
      }
    }
  }
}

function downloadImage() {
  const img = document.getElementById('resultImage');
  if (!img.src) return;
  const link = document.createElement('a');
  link.href = img.src;
  link.download = img.src.split('/').pop();
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

function shareImage() {
  const img = document.getElementById('resultImage');
  if (!img.src) return;

  // Use the absolute URL of the image
  const url = img.src;

  navigator.clipboard.writeText(url).then(() => {
    // Ideally use a proper toast notification, but for now we'll reuse the status message temporarily
    const statusMsg = document.getElementById('statusMessage');
    const originalText = statusMsg.innerText;
    statusMsg.innerText = "Link copied to clipboard!";
    statusMsg.style.color = "var(--accent-color)";

    setTimeout(() => {
      statusMsg.innerText = originalText;
      statusMsg.style.color = "var(--text-secondary)"; // Reset color if needed, though usually handled by CSS or original logic
    }, 2000);
  }).catch(err => {
    console.error('Failed to copy: ', err);
  });
}

function resetUI() {
  const urlInput = document.getElementById('youtubeUrl');
  const resultSection = document.getElementById('resultSection');
  const progressSection = document.getElementById('progressSection');
  const idleState = document.getElementById('idleState');
  const statusMsg = document.getElementById('statusMessage');

  urlInput.value = "";
  resultSection.classList.add('hidden');
  progressSection.classList.add('hidden');
  idleState.classList.remove('hidden');

  statusMsg.innerText = "Ready to create...";
  statusMsg.style.color = "var(--text-secondary)";
}

// Enable Enter key submission
document.addEventListener('DOMContentLoaded', () => {
  const urlInput = document.getElementById('youtubeUrl');
  if (urlInput) {
    urlInput.addEventListener('keydown', function (e) {
      if (e.key === 'Enter') {
        e.preventDefault();
        generateSketchnote();
      }
    });
  }
});
