import { test, expect } from './fixtures';

test.describe('Agentic chat page', () => {
  test('renders the run timeline, tool states, approval, and composer', async ({ page }) => {
    await test.step('navigate to the agentic chat sample', async () => {
      await page.goto('/gallery/page-agentic-chat');
      await page.waitForLoadState('domcontentloaded');
    });

    const frame = page.frameLocator('#preview-frame');

    await test.step('transcript and run marker render', async () => {
      await expect(frame.getByTestId('agentic-transcript')).toBeVisible({ timeout: 10000 });
      await expect(frame.getByText('Run started').first()).toBeVisible();
      await expect(frame.getByText('Thinking').first()).toBeVisible();
    });

    await test.step('plan, tool calls, and statuses render', async () => {
      await expect(frame.getByText('Plan', { exact: true })).toBeVisible();
      await expect(frame.getByText('Search knowledge')).toBeVisible();
      await expect(frame.getByText('Run cypher query (retry)')).toBeVisible();
      await expect(frame.getByText('Waiting for confirmation')).toHaveCount(0);
    });

    await test.step('error tool details expand', async () => {
      await frame.getByText('Run cypher query', { exact: true }).click();
      await expect(frame.getByText('Query timed out after 30s', { exact: false })).toBeVisible();
      await expect(frame.getByText('Error', { exact: true })).toBeVisible();
    });

    await test.step('approval and composer are present', async () => {
      await expect(frame.getByText('Approve tool call')).toBeVisible();
      await expect(frame.getByRole('button', { name: 'Approve' })).toBeVisible();
      await expect(frame.getByRole('button', { name: 'Reject' })).toBeVisible();
      await expect(frame.getByPlaceholder('Ask Memory anything…')).toBeVisible();
      await expect(frame.getByRole('button', { name: 'Stop generating' })).toBeVisible();
    });

    await test.step('session rail shows bucket badges', async () => {
      await expect(frame.getByText('Needs input').first()).toBeVisible();
      await expect(frame.getByText('Failed').first()).toBeVisible();
      await expect(frame.getByText('Done').first()).toBeVisible();
    });
  });
});
