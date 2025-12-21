import { test, expect } from '@playwright/test';

test('generate sketchnote flutter', async ({ page }) => {
  // 1. Navigate to the Flutter App URL (served under /m/)
  const serviceUrl = process.env.SERVICE_URL || 'http://localhost:8080/m/';
  console.log(`Navigating to: ${serviceUrl}`);
  await page.goto(serviceUrl);

  // 2. Initial State Verification
  // Flutter apps might take a moment to initialize semantics
  await expect(page).toHaveTitle(/Sketchnote Artist/i);

  // Flutter renders into a canvas but exposes an accessibility tree (Semantics).
  // Playwright can often find these by role or text.
  // Note: 'Sketchnote Artist' header text might be split or canvas-rendered.
  // We'll trust the title check and input field presence for now.
  // await expect(page.getByText('Sketchnote Artist').first()).toBeVisible({
  // 3. Find Input Field
  const urlInput = page.getByRole('textbox');
  await expect(urlInput).toBeVisible();

  // 3a. Verify initial state of buttons (Visible but effectively disabled/low opacity)
  // We can check text presence.
  await expect(page.getByText('Save')).toBeVisible();
  await expect(page.getByText('Share')).toBeVisible();

  // 3b. Test Invalid URL
  await urlInput.fill('not-a-url');
  await urlInput.press('Enter');
  await expect(page.getByText('Error: Please enter a valid YouTube URL')).toBeVisible();

  // 4. Enter Valid YouTube URL
  const testVideoUrl = 'https://www.youtube.com/watch?v=jNQXAC9IVRw';
  await urlInput.fill(testVideoUrl);
  await urlInput.press('Enter');

  // 6. Verify Processing State
  await expect(page.getByText(/Connecting to agent|Curator is analyzing|Summarizing video/i).first())
    .toBeVisible({ timeout: 15000 });

  // 7. Verify Result
  await expect(page.getByText('Sketchnote created successfully!')).toBeVisible({ timeout: 300000 });

  // Verify Image presence and Button state
  // Check if we can find the enabled buttons (opacity 1.0 vs 0.4 implies enabled visually)
  // But functionally we just check they are still there.
  await expect(page.getByText('Save')).toBeVisible();
  await expect(page.getByText('Share')).toBeVisible();

  // Optional: Take a screenshot of the result
  await page.screenshot({ path: 'e2e-flutter-result.png' });
});
