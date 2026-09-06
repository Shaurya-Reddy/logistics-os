import { defineConfig } from '@playwright/test';
export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  use: { baseURL: process.env.BASE_URL || 'http://127.0.0.1:8080', browserName: 'chromium' },
  reporter: [['list'], ['json', { outputFile: '../reports/browser.json' }]],
});
