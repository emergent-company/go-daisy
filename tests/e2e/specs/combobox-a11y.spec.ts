import { test, expect } from './fixtures';

const TRIGGER = 'div[role="button"][aria-haspopup="listbox"]';
const ID = 'demo-combo-labelled';

test.describe('Combobox accessible trigger', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/gallery/render/combobox/with-label');
    await expect(page.locator(TRIGGER)).toBeVisible({ timeout: 10000 });
  });

  test('label is programmatically associated and controls the listbox', async ({ page }) => {
    const trigger = page.locator(TRIGGER);

    await expect(trigger).toHaveAttribute('id', ID);
    // A <div> is not labelable, so <label for> cannot associate with it; the
    // label carries an id and the trigger points at it with aria-labelledby.
    await expect(trigger).toHaveAttribute('aria-labelledby', `${ID}-label`);
    await expect(page.locator(`#${ID}-label`)).toHaveText('Industry');
    await expect(trigger).toHaveAttribute('tabindex', '0');
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');

    await expect(page.locator(`#${ID}-listbox`)).toHaveAttribute('role', 'listbox');
    await expect(trigger).toHaveAttribute('aria-controls', `${ID}-listbox`);
  });

  test('opens with ArrowDown while closed without freezing the tab', async ({ page }) => {
    // Regression: the trigger's ArrowDown handler opens the popover, and the
    // root keydown handler then runs before Alpine has painted it. The old
    // unbounded do/while scanned fully hidden options forever and hung the
    // renderer.
    const trigger = page.locator(TRIGGER);
    await trigger.focus();
    await page.keyboard.press('ArrowDown');

    await expect(trigger).toHaveAttribute('aria-expanded', 'true');
    await expect(page.locator('.combobox-popover')).toBeVisible();
  });

  test('Escape closes and Enter selects the highlighted option', async ({ page }) => {
    const trigger = page.locator(TRIGGER);

    await trigger.focus();
    await page.keyboard.press('Enter');
    await expect(trigger).toHaveAttribute('aria-expanded', 'true');

    await page.keyboard.press('Escape');
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');

    await page.keyboard.press('Enter');
    await expect(trigger).toHaveAttribute('aria-expanded', 'true');

    await page.keyboard.press('ArrowDown');
    await page.keyboard.press('Enter');
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');
    await expect(trigger).toContainText('Technology');
  });
});
