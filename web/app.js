const API_BASE = window.location.origin;

async function generateSketchnote() {
  const urlInput = document.getElementById('youtubeUrl');
  const generateBtn = document.getElementById('generateBtn');
  const statusMsg = document.getElementById('statusMessage');
  const resultSection = document.getElementById('resultSection');
  const resultImage = document.getElementById('resultImage');
  const btnText = generateBtn.querySelector('.btn-text');
  const loader = generateBtn.querySelector('.loader');

  const videoUrl = urlInput.value.trim();

  if (!videoUrl) {
    statusMsg.innerText = "Please enter a valid YouTube URL.";
    statusMsg.style.color = "#ec4899"; // Pink for error
    return;
  }

  // Reset UI
  statusMsg.innerText = "Curating content and sketching... (this may take a minute)";
  statusMsg.style.color = "#cbd5e1";
  resultSection.classList.add('hidden');
  generateBtn.disabled = true;
  btnText.classList.add('hidden');
  loader.classList.remove('hidden');

  try {
    // 1. Create Session
    const sessionRes = await fetch(`${API_BASE}/apps/sketchnote-artist/users/web-guest/sessions`, {
      method: 'POST'
    });

    if (!sessionRes.ok) {
      throw new Error(`Failed to create session: ${sessionRes.statusText}`);
    }

    const sessionData = await sessionRes.json();
    const sessionId = sessionData.id;

    console.log(`Session created: ${sessionId}`);

    // 2. Run Agent
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

    // 3. Parse result for filename
    // The last message usually contains the summary/tool output.
    // We look for patterns like "Visual_Brief_Title.png" or "Image successfully saved to..."
    // Or if the agent output text references the file.

    let filename = null;

    // Strategy: Search through the response text for the explicit confirmation message we added to the prompt.
    // Pattern: "I have successfully generated the sketchnote: [Filename]"
    const responseText = JSON.stringify(runData);

    // Logic: Look for the specific phrase. Capture everything until a quote, newline, or end of string.
    // This allows spaces and other chars in filenames.
    // We unescape double quotes since we are searching a stringified JSON.
    const match = responseText.match(/I have successfully generated the sketchnote: (.+?)(?:"|\\n|$|\\r)/);

    if (match && match[1]) {
      // The captured group might have trailing chars if the regex was loose, but our non-greedy match should be good.
      // We also trim just in case.
      filename = match[1].trim();
      // Remove any potential trailing period if the agent added one (common LLM behavior)
      if (filename.endsWith('.')) {
        filename = filename.slice(0, -1);
      }
    } else {
      // Try fallback to the old visual brief pattern just in case
      const fallback = responseText.match(/Visual_Brief_[a-zA-Z0-9_\-]+\.png/);
      if (fallback) {
        filename = fallback[0];
      } else {
        throw new Error("Could not find generated image filename in response.");
      }
    }

    // 4. Display Result
    resultImage.src = `${API_BASE}/images/${filename}`;
    resultSection.classList.remove('hidden');
    statusMsg.innerText = "Sketchnote generated successfully!";
    statusMsg.style.color = "#06b6d4"; // Cyan for success

  } catch (error) {
    console.error(error);
    statusMsg.innerText = `Error: ${error.message}`;
    statusMsg.style.color = "#ec4899";
  } finally {
    generateBtn.disabled = false;
    btnText.classList.remove('hidden');
    loader.classList.add('hidden');
  }
}

function downloadImage() {
  const img = document.getElementById('resultImage');
  const link = document.createElement('a');
  link.href = img.src;
  link.download = img.src.split('/').pop();
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}
