import { test, expect } from './fixtures';

// Regression guard for the behaviour-init runtime wiring.
//
// Component inline scripts call gdOnce/gdBehaviour at parse time, so the runtime
// (staticfs/static/js/go-daisy-init.js) must be loaded BEFORE the component's
// script in every document that renders one. The gallery's standalone preview
// document (renderSnippetPage) has no page shell, so it has to inject the
// runtime itself — it previously did not, and every migrated component threw
// "gdBehaviour is not defined" inside the preview iframe.
const RUNTIME_SLUGS = ['sortable-list', 'swiper-carousel', 'theme-toggle'];

test.describe('Preview runtime wiring', () => {
  for (const slug of RUNTIME_SLUGS) {
    test(`preview document for ${slug} boots the behaviour-init runtime`, async ({ page }) => {
      const errors: string[] = [];
      page.on('pageerror', (e) => errors.push(e.message));

      await test.step(`load /gallery/render/${slug}`, async () => {
        await page.goto(`/gallery/render/${slug}`);
        await page.waitForLoadState('load');
      });

      await test.step('runtime globals are defined', async () => {
        expect(await page.evaluate(() => typeof window.gdOnce)).toBe('function');
        expect(await page.evaluate(() => typeof window.gdBehaviour)).toBe('function');
      });

      await test.step('no missing-runtime error', async () => {
        expect(errors.filter((m) => /gdOnce|gdBehaviour|gdRescan/.test(m))).toEqual([]);
      });
    });
  }

  test('detail-page preview iframe for sortable-list boots the runtime', async ({ page }) => {
    const errors: string[] = [];
    page.on('pageerror', (e) => errors.push(e.message));

    await page.goto('/gallery/sortable-list');
    await expect(page.locator('#preview-frame')).toBeAttached({ timeout: 10000 });

    const frame = page.frames().find((f) => f.url().includes('/gallery/render/sortable-list'));
    expect(frame, 'preview iframe should be attached').toBeTruthy();
    await frame!.waitForLoadState('load');

    expect(await frame!.evaluate(() => typeof window.gdOnce)).toBe('function');
    expect(errors.filter((m) => /gdOnce|gdBehaviour|gdRescan/.test(m))).toEqual([]);
  });
});
