import { chromium } from '@playwright/test';
const b = await chromium.launch();
for (const [name, url, sel] of [
  ['ref-main', 'http://localhost:11002/apps-gen-ai-home.html', '#layout-main'],
  ['got-body', 'http://localhost:11001/apps/gen-ai/home', 'body'],
  ['got-content-pane', 'http://localhost:11001/apps/gen-ai/home', '#layout-content'],
  ['got-sidebar', 'http://localhost:11001/apps/gen-ai/home', '#_layout-sidebar'],
  ['got-topbar', 'http://localhost:11001/apps/gen-ai/home', '#_layout-topbar'],
]) {
  const p = await b.newPage({ viewport: { width: 1440, height: 900 } });
  await p.goto(url, { waitUntil: 'networkidle' });
  const v = await p.evaluate((s) => {
    const el = document.querySelector(s);
    return el ? getComputedStyle(el).backgroundColor : 'MISSING';
  }, sel);
  console.log(name + ' ' + sel + ' = ' + v);
  await p.close();
}
await b.close();
