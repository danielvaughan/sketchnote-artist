import { test, expect } from '@playwright/test';

test('generates sketchnote from youtube url', async ({ page }) => {
  // Allow plenty of time for the AI agents to run
  test.setTimeout(180000);

  await page.goto('/');

  // Using "Me at the zoo" as it's short and reliable
  await page.fill('#youtubeUrl', 'https://www.youtube.com/watch?v=jNQXAC9IVRw');
  await page.click('#generateBtn');

  // Verify loader appears
  await expect(page.locator('.loader')).toBeVisible();

  // Verify result image appears eventually
  // The backend returns an image path, frontend sets src
  const resultImage = page.locator('#resultImage');
  await expect(resultImage).toBeVisible({ timeout: 150000 });

  // Verify src is populated correctly
  const src = await resultImage.getAttribute('src');
  expect(src).toContain('/images/');
});
