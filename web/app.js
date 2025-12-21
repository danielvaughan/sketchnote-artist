const API_BASE = window.location.origin;
const APP_NAME = "sketchnote-artist";
let USER_ID = "local-user"; // Default to local-user, updated from server

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
  // Display state will be shown only if video is valid or processing starts
  const displayStage = document.getElementById('displayStage');

  idleState.classList.add('hidden');
  resultSection.classList.add('hidden');
  progressSection.classList.remove('hidden');
  statusMsg.innerText = "Initializing...";
  statusMsg.style.color = "var(--text-secondary)";
  generateBtn.disabled = true;
  document.getElementById('clearBtn').classList.add('hidden'); // Hide clear during generation

  // Handle Video Player
  const videoId = extractVideoID(videoUrl);
  if (videoId) {
    const embedUrl = `https://www.youtube.com/embed/${videoId}?autoplay=1&origin=${window.location.origin}`;
    document.getElementById('youtubeEmbed').src = embedUrl;
    document.getElementById('videoPlayerContainer').classList.remove('hidden');
    displayStage.classList.remove('hidden'); // Show stage if video is valid
  }

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

    // Use shortened URL for agent submission if possible
    let submissionUrl = videoUrl;
    if (videoId) {
      submissionUrl = `https://youtu.be/${videoId}`;
    }

    const response = await fetch('/run_sse', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        appName: APP_NAME,
        userId: USER_ID,
        sessionId: sessionId,
        newMessage: {
          role: "user",
          parts: [{ text: submissionUrl }]
        },
        streaming: true
      })
    });

    if (!response.ok) throw new Error(`Stream failed: ${response.statusText}`);

    // 3. Read the Stream
    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    let accumulatedResponseText = ''; // Accumulate text across all chunks

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

                // Fast-path: Check status messages for image filenames
                if (data.message.includes('.png')) {
                  const match = data.message.match(/[\w\s.-]+\.png/);
                  if (match) {
                    loadImage(match[0].trim(), resultImage, progressSection, resultSection, statusMsg);
                  }
                }
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
            // Update the accumulator with any new text
            if (event.content && event.content.parts) {
              for (const part of event.content.parts) {
                if (part.text) {
                  accumulatedResponseText += part.text;
                }
              }
            }
            // Pass the accumulated text to the handler
            handleAgentEvent(event, accumulatedResponseText, statusMsg, resultImage, progressSection, resultSection);
          } catch (e) {
            console.warn("Failed to parse event", e);
          }
        }
      }
    }

    // After stream completes, if we are still 'watching', show a fallback message
    if (resultSection.classList.contains('hidden')) {
      progressSection.classList.add('hidden');
      if (statusMsg.innerText === "Watching video..." || statusMsg.innerText === "Initializing...") {
        statusMsg.innerText = "Process ended without creating a sketchnote.";
        statusMsg.style.color = "var(--accent-color)";
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
    toggleClearBtn(); // Restore clear button state
  }
}

function loadImage(filename, resultImage, progressSection, resultSection, statusMsg) {
  const imagePath = `/images/${filename}`;

  // Avoid re-triggering if already set to this path
  if (resultImage.src.endsWith(imagePath) && !resultSection.classList.contains('hidden')) return;

  console.log("Attempting to load image:", imagePath);

  const img = new Image();
  img.onload = () => {
    resultImage.src = imagePath;
    progressSection.classList.add('hidden');
    resultSection.classList.remove('hidden');
    statusMsg.innerText = "Sketchnote created successfully!";
    statusMsg.style.color = "var(--text-secondary)";
    console.log("Image loaded successfully:", imagePath);
  };
  img.onerror = () => {
    // Only show error if we haven't successfully loaded it yet
    if (resultSection.classList.contains('hidden')) {
      console.error("Failed to load image:", imagePath);
      statusMsg.innerText = "Image generated but failed to load.";
      statusMsg.style.color = "var(--accent-color)";
    }
  };
  img.src = imagePath;
}

function handleAgentEvent(event, accumulatedText, statusMsg, resultImage, progressSection, resultSection) {
  // Debug: Log every event to inspect structure
  console.log("Received Event:", event);

  // If we've already succeeded, don't let further text fragments mess up the UI
  if (!resultSection.classList.contains('hidden')) return;

  // 1. Model Call (Thinking) - Check for ADK 'modelCall' structure
  if (event.models && event.models.length > 0) {
    statusMsg.innerText = "Curator is analyzing video...";
    statusMsg.style.color = "var(--text-secondary)";
  }

  // Check content parts for Tool Calls and Text
  if (event.content && event.content.parts) {
    for (const part of event.content.parts) {

      // 2. Tool Call (Summarizing/Sketching)
      if (part.functionCall) {
        const toolName = part.functionCall.name;
        console.log("Tool Call Detected:", toolName);

        if (toolName.includes('summarize')) {
          statusMsg.innerText = "Summarizing video content...";
          document.getElementById('displayStage').classList.remove('hidden'); // Ensure visible when tool starts
        } else if (toolName.includes('generate_image')) {
          statusMsg.innerText = "Artist is sketching...";
        }
        statusMsg.style.color = "var(--text-secondary)";
      }

      // 3. General Text Parts (Feedback from Agent)
      if (part.text && !part.text.includes('.png')) {
        // Detect error phrases to highlight in red and handle UI states,
        // but do not display the text in the status message as per user request.
        if (part.text.toLowerCase().includes("can't process") ||
          part.text.toLowerCase().includes("invalid") ||
          part.text.toLowerCase().includes("error")) {
          statusMsg.innerText = part.text; // Still show explicit errors
          statusMsg.style.color = "var(--accent-color)";
          progressSection.classList.add('hidden'); // Stop animation if it's an error
          document.getElementById('idleState').classList.remove('hidden');
          document.getElementById('displayStage').classList.add('hidden'); // Hide stage on validation error
        }
      }

      // 4. Check Accumulated Text for Image Filename
      if (accumulatedText && accumulatedText.includes('.png')) {
        // Found image path!
        const match = accumulatedText.match(/[\w\s.-]+\.png/);
        if (match) {
          loadImage(match[0].trim(), resultImage, progressSection, resultSection, statusMsg);
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
  const videoPlayerContainer = document.getElementById('videoPlayerContainer');
  const displayStage = document.getElementById('displayStage');
  const clearBtn = document.getElementById('clearBtn');

  urlInput.value = "";
  resultSection.classList.add('hidden');
  progressSection.classList.add('hidden');
  videoPlayerContainer.classList.add('hidden');
  displayStage.classList.add('hidden'); // Hide the entire stage
  document.getElementById('youtubeEmbed').src = ""; // Stop video
  idleState.classList.remove('hidden');
  clearBtn.classList.add('hidden'); // Hide clear button

  statusMsg.innerText = "Ready to create...";
  statusMsg.style.color = "var(--text-secondary)";
}

function toggleClearBtn() {
  const urlInput = document.getElementById('youtubeUrl');
  const clearBtn = document.getElementById('clearBtn');
  if (urlInput.value.length > 0) {
    clearBtn.classList.remove('hidden');
  } else {
    clearBtn.classList.add('hidden');
  }
}

function extractVideoID(url) {
  const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*/;
  const match = url.match(regExp);
  return (match && match[2].length === 11) ? match[2] : null;
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
    urlInput.addEventListener('input', toggleClearBtn);
  }

  // Fetch user identity
  fetch('/me')
    .then(r => r.json())
    .then(data => {
      if (data.email) {
        USER_ID = data.email;
        console.log("Authenticated as:", USER_ID);
      }
    })
    .catch(e => console.warn("Failed to load user identity", e));

  fetch('/version')
    .then(r => r.json())
    .then(data => {
      document.getElementById('appVersion').innerText = `v${data.version}`;
    })
    .catch(e => console.warn("Failed to load version", e));
});
