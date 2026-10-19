import { chromium } from '@playwright/test';
const b = await chromium.launch();
const targets = [
  ['genai-home', 'http://localhost:11001/apps/gen-ai/home'],
  ['genai-content', 'http://localhost:11001/apps/gen-ai/content'],
  ['genai-image', 'http://localhost:11001/apps/gen-ai/image'],
  ['genai-library', 'http://localhost:11001/apps/gen-ai/library'],
  ['ecom-products-create', 'http://localhost:11001/apps/ecommerce/products/create'],
  ['ecom-customers', 'http://localhost:11001/apps/ecommerce/customers'],
];
for (const [name, url] of targets) {
  const p = await b.newPage({ viewport: { width: 1440, height: 900 } });
  await p.goto(url, { waitUntil: 'networkidle' });
  await p.waitForTimeout(400);
  await p.screenshot({ path: `/tmp/opencode/followup-${name}.png`, fullPage: true });
  await p.close();
}
await b.close();
console.log('done');
