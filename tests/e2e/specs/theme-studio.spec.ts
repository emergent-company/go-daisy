import { test, expect } from './fixtures';

// The theme studio is the only page that exercises the typography roles and
// the native colour pickers end to end: palette edits must land in the scoped
// preview, and Apply must push the result through the custom-theme seam.

async function bootStudio(page: import('@playwright/test').Page) {
  await page.goto('/gallery/theme-studio');
  await page.waitForFunction(() => (window as any).Alpine !== undefined);
  // Start from defaults rather than a draft left by a previous run.
  await page.evaluate(() => localStorage.removeItem('gallery-theme-studio'));
  await page.reload();
  await page.waitForFunction(() => (window as any).Alpine !== undefined);
  await expect(page.locator(`input[type="color"]`)).toHaveCount(20);
}

test.describe('Theme studio', () => {
  test('colour swatches open a picker and keep oklch tokens in oklch', async ({ page }) => {
    await bootStudio(page);

    const picker = page.locator('input[aria-label="Pick primary color"]');
    await expect(picker).toBeAttached();

    const text = page.locator('input[aria-label="primary value"]');
    const before = await text.inputValue();
    expect(before).toMatch(/^oklch\(/);

    // Native pickers only emit #rrggbb; the studio converts back to oklch.
    await picker.evaluate((el) => {
      const input = el as HTMLInputElement;
      input.value = '#ff0000';
      input.dispatchEvent(new Event('input', { bubbles: true }));
    });

    await expect(text).toHaveValue(/^oklch\(/);
    await expect(text).not.toHaveValue(before);

    // A hex-authored token stays hex on the next pick.
    await text.fill('#123456');
    await picker.evaluate((el) => {
      const input = el as HTMLInputElement;
      input.value = '#00ff00';
      input.dispatchEvent(new Event('input', { bubbles: true }));
    });
    await expect(text).toHaveValue('#00FF00');
  });

  test('one preview board showcases components and typography', async ({ page }) => {
    await bootStudio(page);

    const preview = page.locator('#studio-preview');
    await expect(preview).toHaveCount(1);
    await expect(preview.locator('table')).toBeVisible();
    await expect(preview.locator('input[type="range"]')).toBeVisible();
    await expect(preview.locator('.alert')).toHaveCount(4);
    await expect(preview.locator('[role="tablist"]')).toBeVisible();
    await expect(preview.locator('h1')).toHaveCount(1);
    await expect(preview.getByText('Data display')).toBeVisible();
  });

  test('preview sits full height beside the settings rail', async ({ page }) => {
    // Pin a width above the xl breakpoint so the two-column layout is the
    // contract under test, not the runner's default viewport.
    await page.setViewportSize({ width: 1440, height: 900 });
    await bootStudio(page);

    const preview = page.locator('#studio-preview');
    const settings = page.locator('#studio-settings');

    // The board grows with the page instead of scrolling inside the card.
    const scroll = await preview.evaluate((el) => ({
      overflowY: getComputedStyle(el).overflowY,
      scrollHeight: el.scrollHeight,
      clientHeight: el.clientHeight,
    }));
    expect(scroll.overflowY).toBe('visible');
    expect(scroll.scrollHeight).toBeLessThanOrEqual(scroll.clientHeight + 1);

    // Settings live in the right column on desktop.
    const previewBox = await preview.boundingBox();
    const settingsBox = await settings.boundingBox();
    expect(previewBox).not.toBeNull();
    expect(settingsBox).not.toBeNull();
    expect(settingsBox!.x).toBeGreaterThan(previewBox!.x + previewBox!.width - 1);
  });

  test('typography roles drive the preview and survive Apply to gallery', async ({ page }) => {
    await bootStudio(page);

    const headingFamily = page.locator('input[aria-label="Heading font-family"]');
    await expect(headingFamily).toBeVisible();
    await headingFamily.fill('Georgia, "Times New Roman", serif');

    const previewH1 = page.locator('[data-theme="gd-preview"] h1').first();
    await expect
      .poll(() => previewH1.evaluate((el) => getComputedStyle(el).fontFamily))
      .toContain('Georgia');

    await page.getByRole('button', { name: 'Apply to gallery' }).click();
    await page.waitForFunction(() => document.documentElement.getAttribute('data-theme') === 'custom');

    const applied = await page
      .locator('#gallery-content h1')
      .first()
      .evaluate((el) => getComputedStyle(el).fontFamily);
    expect(applied).toContain('Georgia');
  });

  test('export includes typography tokens in the theme block', async ({ page }) => {
    await bootStudio(page);

    await page.locator('input[aria-label="Body font-family"]').fill('ui-monospace, Menlo, monospace');
    const css = await page.locator('textarea[readonly]').inputValue();

    // Font stacks must NOT live inside the @plugin block: Tailwind's plugin
    // value parser splits top-level commas into repeated declarations ending
    // in the generic family (`sans-serif`). They belong in a scoped rule.
    const plugin = css.match(/@plugin[^{]*\{[\s\S]*?\n\}/)?.[0] ?? '';
    expect(plugin).toContain('--color-base-100:');
    expect(plugin).not.toContain('--font-');
    expect(plugin).not.toContain('--type-');

    const scoped = css.match(/\[data-theme="[^"]+"\]\s*\{[\s\S]*?\}/)?.[0] ?? '';
    expect(scoped).toContain('--font-body: ui-monospace, Menlo, monospace;');
    expect(scoped).toContain('--font-heading:');
    expect(scoped).toContain('--type-weight-heading:');

    // Round-trip: an imported block sets the role fields from its variables.
    await page.locator('details summary').filter({ hasText: 'Import CSS' }).click();
    await page
      .locator('textarea[placeholder*="daisyui/theme"]')
      .fill('@plugin "daisyui/theme" {\n  --font-heading: Georgia, serif;\n  --type-leading-body: 1.75;\n}');
    await page.getByRole('button', { name: 'Load', exact: true }).click();

    await expect(page.locator('input[aria-label="Heading font-family"]')).toHaveValue('Georgia, serif');
    await expect
      .poll(() => page.locator('[data-theme="gd-preview"] h1').first().evaluate((el) => getComputedStyle(el).fontFamily))
      .toContain('Georgia');
  });
});
