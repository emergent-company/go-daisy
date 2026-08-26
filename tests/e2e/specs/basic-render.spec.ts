import { test, expect } from './fixtures';

const BASIC_SLUGS = ['button', 'badge', 'card-real', 'alert'];

test.describe('Core component render', () => {
  for (const slug of BASIC_SLUGS) {
    test(`renders ${slug} detail page with iframe`, async ({ page }) => {
      await test.step(`navigate to /gallery/${slug}`, async () => {
        await page.goto(`/gallery/${slug}`);
        await page.waitForLoadState('domcontentloaded');
      });

      await test.step('page shell is visible', async () => {
        await expect(page.locator('#gallery-shell')).toBeVisible({ timeout: 10000 });
      });

      await test.step('component name heading is visible', async () => {
        await expect(page.locator('h1')).toBeVisible({ timeout: 5000 });
      });

      await test.step('preview iframe is present with non-empty src', async () => {
        const frame = page.locator('#preview-frame');
        await expect(frame).toBeAttached({ timeout: 10000 });
        const src = await frame.getAttribute('src');
        expect(src).toBeTruthy();
      });
    });
  }
});
