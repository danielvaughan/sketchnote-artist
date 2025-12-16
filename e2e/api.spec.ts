import { test, expect } from '@playwright/test';
import { execSync } from 'child_process';

test('api: generate sketchnote', async ({ request }) => {
  const serviceUrl = process.env.SERVICE_URL;
  expect(serviceUrl).toBeDefined();

  const appName = 'sketchnote-artist';
  const userId = 'api-test-user';
  const startTime = new Date().toISOString();

  try {
    // 1. Create Session
    const sessionResp = await request.post(`${serviceUrl}/apps/${appName}/users/${userId}/sessions`);
    if (!sessionResp.ok()) {
      throw new Error(`Failed to create session: ${sessionResp.status()} ${sessionResp.statusText()}`);
    }
    const sessionData = await sessionResp.json();
    const sessionId = sessionData.id;
    expect(sessionId).toBeDefined();
    console.log(`Created Session: ${sessionId}`);

    // 2. Run Agent
    const testVideoUrl = 'https://www.youtube.com/watch?v=jNQXAC9IVRw'; // Me at the zoo
    console.log(`Submitting video: ${testVideoUrl}`);

    // Call streaming endpoint
    const response = await request.post(`${serviceUrl}/run_sse`, {
      data: {
        appName: appName,
        userId: userId,
        sessionId: sessionId,
        newMessage: {
          role: 'user',
          parts: [
            { text: testVideoUrl }
          ]
        },
        streaming: true
      },
      headers: {
        'Accept': 'text/event-stream'
      },
      timeout: 300000 // 5 minutes
    });

    if (!response.ok()) {
      throw new Error(`Stream request failed: ${response.status()} ${response.statusText()}`);
    }

    // Read the full SSE stream text
    // Note: This waits for the server to close the connection, which happens when 'event: done' is sent and handler returns
    const streamText = await response.text();
    console.log(`Received stream length: ${streamText.length} bytes`);
    console.log("--- STREAM CONTENT START ---");
    console.log(streamText);
    console.log("--- STREAM CONTENT END ---");

    // Parse SSE events to find the final image
    const lines = streamText.split('\n');
    let foundImage = false;
    let finalContent = '';

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        try {
          const jsonStr = line.substring(6);
          if (jsonStr.trim() === '{}') continue; // Done event data

          const event = JSON.parse(jsonStr);

          // Check for Content with image
          if (event.content && event.content.parts) {
            for (const part of event.content.parts) {
              if (part.text) {
                finalContent += part.text;
                if (part.text.includes('.png')) {
                  foundImage = true;
                  console.log(`Found image in event: ${part.text}`);
                }
              }
            }
          }
        } catch (e) {
          // Ignore parse errors (e.g. partial lines if any, though .text() should be complete)
        }
      }
    }

    expect(foundImage).toBe(true);
    expect(finalContent).toContain('.png');

  } catch (error) {
    console.error("Test failed, fetching backend logs...");
    try {
      // Fetch logs from Cloud Run for the duration of the test
      // Filter for ERROR severity to find crashes or failures
      const filter = `resource.type="cloud_run_revision" AND resource.labels.service_name="sketchnote-artist-dev" AND timestamp>="${startTime}" AND severity>=ERROR`;
      const cmd = `gcloud logging read '${filter}' --project=sketchnote-artist-application --limit=10 --format="value(textPayload,jsonPayload.message)"`;

      console.log(`Running: ${cmd}`);
      const logs = execSync(cmd).toString();

      if (logs.trim()) {
        console.log("--- BACKEND ERROR LOGS ---");
        console.log(logs);
        console.log("--------------------------");
      } else {
        console.log("No backend error logs found.");
      }
    } catch (logError) {
      console.error("Failed to fetch logs:", logError);
    }
    throw error; // Re-throw original error to fail the test
  }
});
