const API_BASE = window.location.origin;

async function generateSketchnote() {
  const urlInput = document.getElementById('youtubeUrl');
  const generateBtn = document.getElementById('generateBtn');
  const btnIcon = generateBtn.querySelector('.material-icons');


  const idleState = document.getElementById('idleState');
  const progressSection = document.getElementById('progressSection');
  const resultSection = document.getElementById('resultSection');
  const resultImage = document.getElementById('resultImage');
  const statusMsg = document.getElementById('statusMessage');

  const videoUrl = urlInput.value.trim();

  if (!videoUrl) {
    statusMsg.innerText = "Please enter a valid YouTube URL.";
    statusMsg.style.color = "#FF0000";
    return;
  }

  // 1. UI Transition to Processing
  idleState.classList.add('hidden');
  resultSection.classList.add('hidden');
  progressSection.classList.remove('hidden');

  // Button State
  generateBtn.disabled = true;


  statusMsg.innerText = "Starting engine...";
  statusMsg.style.color = "var(--text-secondary)";

  try {
    // 2. Create Session
    statusMsg.innerText = "Connecting to agent...";
    const sessionRes = await fetch(`${API_BASE}/apps/sketchnote-artist/users/web-guest/sessions`, {
      method: 'POST'
    });

    if (!sessionRes.ok) throw new Error(`Failed to create session: ${sessionRes.statusText}`);

    const sessionData = await sessionRes.json();
    const sessionId = sessionData.id;
    console.log(`Session created: ${sessionId}`);

    // 3. Run Agent (Long Polling / Wait)
    statusMsg.innerText = "Watching video & sketching...";

    // Simulating "steps" by updating text occasionally if it takes too long? 
    // For now, keep it simple as requested.

    const runRes = await fetch(`${API_BASE}/run`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        appName: "sketchnote-artist",
        userId: "web-guest",
        sessionId: sessionId,
        newMessage: {
          role: "user",
          parts: [{ text: videoUrl }]
        }
      })
    });

    if (!runRes.ok) throw new Error("Agent execution failed");

    const runData = await runRes.json();
    console.log("Agent response:", runData);

    // 4. Extract Filename
    let filename = null;
    const responseText = JSON.stringify(runData); // Naive, but usually works for finding the filename in the JSON dump if structure varies

    // Regex from previous version
    const match = responseText.match(/(?:generated the sketchnote: |filename: |save it as |saved to )["']?([^"'\s,]+?\.(?:png|jpg|jpeg|webp))["']?/i);

    if (match && match[1]) {
      filename = match[1].trim();
      if (filename.endsWith('.')) filename = filename.slice(0, -1);
    } else {
      const fallback = responseText.match(/([a-zA-Z0-9_\-]+\.(?:png|jpg|webp))/);
      if (fallback) filename = fallback[0];
    }

    if (!filename) throw new Error("Could not find generated image filename.");

    // 5. Show Result
    // Add small delay to ensure FS write? (Usually immediate if response came back)
    resultImage.src = `${API_BASE}/images/${filename}`;

    progressSection.classList.add('hidden');
    resultSection.classList.remove('hidden');
    statusMsg.innerText = "Sketchnote created successfully!";

  } catch (error) {
    console.error(error);
    statusMsg.innerText = `Error: ${error.message}`;
    statusMsg.style.color = "#FF0000";
    progressSection.classList.add('hidden');
    idleState.classList.remove('hidden'); // Go back to idle on error
  } finally {
    generateBtn.disabled = false;

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
