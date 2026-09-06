import { test, expect } from '@playwright/test';
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { resolve } from 'node:path';

test('production shell, keyboard retry, and manifest network accounting', async ({ page, baseURL }) => {
  const requests = new Set<string>();
  const errors: string[] = [];
  page.on('pageerror', error => errors.push(error.message));
  page.on('request', request => {
    const url = new URL(request.url());
    if (url.origin !== new URL(baseURL!).origin) errors.push(`remote request: ${url.href}`);
    if (/\.(js|css)$/.test(url.pathname)) requests.add(url.pathname);
  });
  await page.goto('/');
  await expect(page.getByRole('heading', { name: 'System ready' })).toBeVisible();
  await page.keyboard.press('Tab');
  await expect(page.getByRole('button', { name: 'Check again' })).toBeFocused();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('heading', { name: 'System ready' })).toBeVisible();
  const manifest = JSON.parse(readFileSync(resolve('dist/.vite/manifest.json'), 'utf8'));
  const expected = new Set<string>();
  const visited = new Set<string>();
  function visit(key: string) {
    if (visited.has(key)) return;
    visited.add(key);
    const chunk = manifest[key];
    expected.add('/' + chunk.file);
    for (const css of chunk.css || []) expected.add('/' + css);
    for (const child of chunk.imports || []) visit(child);
  }
  visit('index.html');
  expect([...requests].sort()).toEqual([...expected].sort());
  expect(errors).toEqual([]);
  mkdirSync('../reports', { recursive: true });
  writeFileSync('../reports/network.json', JSON.stringify({ expected: [...expected], requested: [...requests] }, null, 2));
  await page.screenshot({ path: '../reports/s0-shell.png', fullPage: true });
});

test('loading, failure, and recovery are explicit', async ({ page }) => {
  let release!: () => void;
  const waiting = new Promise<void>(resolve => { release = resolve; });
  await page.route('**/ready', async route => {
    await waiting;
    await route.fulfill({ status: 503, contentType: 'application/json', body: '{"status":"unavailable"}' });
  });
  await page.goto('/');
  await expect(page.getByRole('button', { name: 'Check again' })).toBeDisabled();
  await expect(page.getByRole('heading', { name: 'Checking connection…' })).toBeVisible();
  release();
  await expect(page.getByRole('heading', { name: 'Database unavailable' })).toBeVisible();
  await page.unroute('**/ready');
  await page.getByRole('button', { name: 'Check again' }).click();
  await expect(page.getByRole('heading', { name: 'System ready' })).toBeVisible();
});
