import { test, expect } from './fixtures';

test.describe('Agentic components', () => {
  test('approval card expands its arguments', async ({ page }) => {
    await page.goto('/gallery/approval-card');

    const frame = page.frameLocator('#preview-frame');
    await expect(frame.getByText('Approve tool call')).toBeVisible({ timeout: 10000 });
    await expect(frame.getByRole('button', { name: 'Approve' })).toBeVisible();
    await expect(frame.getByRole('button', { name: 'Reject' })).toBeVisible();

    await frame.getByText('Arguments', { exact: true }).click();
    await expect(frame.getByText('"incident": "inc_9f2"', { exact: false })).toBeVisible();
  });

  test('agent question supports single selection', async ({ page }) => {
    await page.goto('/gallery/agent-question');

    const frame = page.frameLocator('#preview-frame');
    await expect(frame.getByText('The agent has a question')).toBeVisible({ timeout: 10000 });

    await frame.getByText('Q2 2026', { exact: true }).click();
    await expect(frame.locator('input[value="q2"]')).toBeChecked();
    await expect(frame.locator('input[value="both"]')).not.toBeChecked();
  });
});
