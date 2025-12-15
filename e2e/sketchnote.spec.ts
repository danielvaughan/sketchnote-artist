import { test, expect } from '@playwright/test';

test('generate sketchnote', async ({ page }) => {
  // 1. Navigate to the Dev URL
  // We expect the URL to be passed via environment variable for flexibility
  const serviceUrl = process.env.SERVICE_URL;
  if (!serviceUrl) {
    throw new Error('SERVICE_URL environment variable is required');
  }

  console.log(`Navigating to: ${serviceUrl}`);
  await page.goto(serviceUrl);

  // 2. Initial State Verification
  await expect(page).toHaveTitle(/Sketchnote Artist/i);
  const urlInput = page.locator('#youtubeUrl');
  await expect(urlInput).toBeVisible();

  // 3. Enter YouTube URL
  // Testing with "Me at the zoo" - short, stable video
  const testVideoUrl = 'https://www.youtube.com/watch?v=jNQXAC9IVRw';
  await urlInput.fill(testVideoUrl);

  // 4. Submit via Enter Key (Testing new UI feature)
  await urlInput.press('Enter');

  // 5. Verify Processing State
  const statusMsg = page.locator('#statusMessage');
  await expect(statusMsg).toContainText(/Watching video/i, { timeout: 10000 });

  // 6. Verify Result
  // This might take regular processing time (~30-60s depending on backend)
  // We'll increase timeout for the final result
  await expect(statusMsg).toHaveText('Sketchnote created successfully!', { timeout: 120000 });

  const resultImage = page.locator('#resultImage');
  await expect(resultImage).toBeVisible();

  // Optional: Take a screenshot of the result
  await page.screenshot({ path: 'e2e-result.png' });
});
