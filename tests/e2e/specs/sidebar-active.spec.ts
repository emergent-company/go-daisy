import { test, expect } from './fixtures';

test.describe('Sidebar active state', () => {
  test('marks active page on initial load and updates on HTMX navigation', async ({ page }) => {
    await test.step('navigate to button component', async () => {
      await page.goto('/gallery/button');
      await expect(page.locator('#gallery-content')).toBeVisible({ timeout: 5000 });
    });
    await test.step('verify button sidebar item is active', async () => {
      const active = page.locator('#_layout-sidebar a.menu-item.active');
      await expect(active).toHaveCount(1);
      await expect(active).toHaveAttribute('href', '/gallery/button');
    });

    await test.step('click hero component in sidebar', async () => {
      await page.locator('a.menu-item[href="/gallery/hero"]').click();
      await expect(page.locator('#gallery-content h1')).toContainText('Hero', { timeout: 5000 });
    });
    await test.step('verify active class swapped from button to hero', async () => {
      await expect(
        page.locator('a.menu-item[href="/gallery/button"]')
      ).not.toHaveClass(/active/);
      await expect(
        page.locator('a.menu-item[href="/gallery/hero"]')
      ).toHaveClass(/active/);
    });

    await test.step('click SDK Reference link in sidebar', async () => {
      await page.locator('a.menu-item[href="/gallery/docs"]').click();
      await expect(page.locator('#gallery-content h1')).toContainText('SDK Reference', { timeout: 5000 });
    });
    await test.step('verify docs sidebar item is active', async () => {
      await expect(
        page.locator('a.menu-item[href="/gallery/docs"]')
      ).toHaveClass(/active/);
    });

    await test.step('navigate to index page', async () => {
      await page.goto('/gallery');
      await expect(page.locator('#gallery-shell')).toBeVisible();
    });
    await test.step('verify no sidebar items are active on index', async () => {
      await expect(page.locator('#_layout-sidebar .menu-item.active')).toHaveCount(0);
    });
  });
});
