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

  test('opens with ArrowDown while closed without freezing or pre-highlighting', async ({ page }) => {
    // Regression: the trigger's ArrowDown handler opens the popover, and the
    // root keydown handler then runs before Alpine has painted it. The old
    // unbounded do/while scanned fully hidden options forever and hung the
    // renderer; a bounded scan that kept the last index also pre-highlighted the
    // final option, so a follow-up Enter picked "Retail" instead of the first
    // enabled option.
    const trigger = page.locator(TRIGGER);
    await trigger.focus();
    await page.keyboard.press('ArrowDown');

    await expect(trigger).toHaveAttribute('aria-expanded', 'true');
    await expect(page.locator('.combobox-popover')).toBeVisible();
    await expect(page.locator('[role="option"].bg-base-200')).toHaveCount(0);

    // Disabled options cannot be selected by pointer either.
    await page.locator('[role="option"][aria-disabled="true"]').click({ force: true });
    await expect(trigger).toContainText('Select industry...');
    await trigger.focus();

    // Next ArrowDown starts at the first option and skips the disabled one.
    await page.keyboard.press('ArrowDown');
    await expect(page.locator('[role="option"].bg-base-200')).toHaveText('Technology');

    await page.keyboard.press('Enter');
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');
    await expect(trigger).toContainText('Technology');

    // Selected options expose aria-selected; the clear-all control lives outside
    // the presentational trigger subtree.
    await expect(page.locator('[role="option"][aria-selected="true"]')).toHaveCount(1);

    // Clear-all sits outside the trigger and resets the selection.
    const clear = page.locator('button[aria-label="Clear all"]');
    await expect(clear).toBeVisible();
    await clear.click();
    await expect(trigger).toContainText('Select industry...');
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');
  });

  test('Escape closes while Space selects the highlighted option', async ({ page }) => {
    const trigger = page.locator(TRIGGER);

    await trigger.focus();
    await page.keyboard.press('Enter');
    await expect(trigger).toHaveAttribute('aria-expanded', 'true');

    // ArrowUp from the "nothing highlighted" sentinel wraps to the last
    // enabled option, not the second-to-last.
    await page.keyboard.press('ArrowUp');
    await expect(page.locator('[role="option"].bg-base-200')).toHaveText('Retail');

    await page.keyboard.press('Escape');
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');

    await page.keyboard.press('Enter');
    await expect(trigger).toHaveAttribute('aria-expanded', 'true');
    await page.keyboard.press('ArrowDown');
    await page.keyboard.press('Space');
    await expect(trigger).toHaveAttribute('aria-expanded', 'false');
    await expect(trigger).toContainText('Technology');
  });
});
