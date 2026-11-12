import type { Page } from '@playwright/test';
import { test, expect } from './fixtures';

// Regression: a dropdown whose trigger is a plain ui.Button opened exactly once
// per page load. ui.Dropdown renders [data-gd-menu-root] but only
// ui.DropdownTrigger renders [data-gd-menu-trigger], so go-daisy-a11y.js never
// wired the root. Its delegated document click listener still called close() on
// every root it did not contain, stamping data-gd-menu-closed; daisyUI turns
// that marker into visibility:hidden on :focus-within, and only the handlers
// wire() installs clear it. Every later click focused the trigger and left the
// panel hidden until a reload.
//
// The gallery's dropdown/Examples story renders exactly that shape (two plain
// buttons, no data-gd-menu-trigger), so it doubles as the fixture.

function rootOf(page: Page, name: string) {
  return page.locator('div.dropdown', { has: page.getByRole('button', { name, exact: true }) });
}

test.describe('dropdown with a plain button trigger (unwired root)', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/gallery/render/dropdown/examples');
  });

  test('reopens after an outside click and never keeps a closed marker', async ({ page }) => {
    const root = rootOf(page, 'Options');
    const trigger = root.getByRole('button', { name: 'Options', exact: true });
    const menu = root.locator('[data-gd-menu]');

    // This is the shape that used to break: no a11y wiring, pure daisyUI CSS.
    expect(await root.getAttribute('data-gd-menu-wired')).toBeNull();
    await expect(menu).toBeHidden();

    await trigger.click();
    await expect(menu).toBeVisible();

    // Any click elsewhere on the page is what used to pin the dropdown shut.
    await page.mouse.click(4, 4);
    await expect(menu).toBeHidden();
    expect(await root.getAttribute('data-gd-menu-closed')).toBeNull();

    // The regression: this click focused the trigger but the panel stayed hidden.
    await trigger.click();
    await expect(menu).toBeVisible();

    await page.mouse.click(4, 4);
    await expect(menu).toBeHidden();

    await trigger.click();
    await expect(menu).toBeVisible();
  });

  test('the second unwired dropdown on the page reopens too', async ({ page }) => {
    const root = rootOf(page, 'Actions');
    const trigger = root.getByRole('button', { name: 'Actions', exact: true });
    const menu = root.locator('[data-gd-menu]');

    await trigger.click();
    await expect(menu).toBeVisible();
    await page.mouse.click(4, 4);
    await expect(menu).toBeHidden();
    await trigger.click();
    await expect(menu).toBeVisible();
  });
});

test.describe('dropdown with a ui.DropdownTrigger (wired root)', () => {
  test('still closes on an outside click and reopens', async ({ page }) => {
    await page.goto('/gallery/render/dropdown/interactive');

    const root = page.locator('div.dropdown', { has: page.locator('[data-gd-menu-trigger]') });
    const trigger = root.locator('[data-gd-menu-trigger]');
    const menu = root.locator('[data-gd-menu]');

    await expect(trigger).toHaveAttribute('aria-expanded', 'false');

    await trigger.click();
    await expect(menu).toBeVisible();
    await expect(trigger).toHaveAttribute('aria-expanded', 'true');

    // A wired root is closed by the runtime, which is what keeps aria-expanded
    // truthful; the fix must not turn that off.
    await page.mouse.click(4, 4);
    await expect(menu).toBeHidden();
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');

    await trigger.click();
    await expect(menu).toBeVisible();
  });
});
