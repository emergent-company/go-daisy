import { defineConfig } from '@playwright/test';

const PORT = process.env.GALLERY_E2E_PORT || '11001';

export default defineConfig({
  testDir: './specs',
  fullyParallel: true,
  retries: 1,
  use: {
    baseURL: `http://localhost:${PORT}/gallery`,
    headless: true,
    screenshot: 'only-on-failure',
    trace: 'on',
  },
  globalSetup: './helpers/global-setup.ts',
  globalTeardown: './helpers/global-teardown.ts',
});
