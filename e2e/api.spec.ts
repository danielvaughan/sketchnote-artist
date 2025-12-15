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

    const runResp = await request.post(`${serviceUrl}/run`, {
      data: {
        appName: appName,
        userId: userId,
        sessionId: sessionId,
        newMessage: {
          role: 'user',
          parts: [
            { text: testVideoUrl }
          ]
        }
      },
      timeout: 300000 // 5 minutes for processing
    });

    if (!runResp.ok()) {
      throw new Error(`Run Agent failed: ${runResp.status()} ${runResp.statusText()}`);
    }

    const runData = await runResp.json();

    // Verify response structure
    // The artist agent returns a text response describing the image path
    expect(runData.responses).toBeDefined();
    expect(runData.responses.length).toBeGreaterThan(0);

    const lastResponse = runData.responses[runData.responses.length - 1];
    const textContent = lastResponse.parts[0].text;
    console.log(`Agent Response: ${textContent}`);

    expect(textContent).toContain('sketchnote here:');
    expect(textContent).toContain('.png');

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
