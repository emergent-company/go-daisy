import { chromium } from '@playwright/test';
const b = await chromium.launch();
const p = await b.newPage({ viewport: { width: 1440, height: 900 } });
await p.goto('http://localhost:11001/apps/gen-ai/home', { waitUntil: 'networkidle' });
const info = await p.evaluate(() => {
  const el = [...document.querySelectorAll('p')].find(x => x.textContent.trim() === 'Blog Post Ideas');
  const chain = [];
  let n = el;
  while (n && n !== document.documentElement) {
    const cs = getComputedStyle(n);
    chain.push({ tag: n.tagName, cls: n.className, fs: cs.fontSize });
    n = n.parentElement;
  }
  const bodyFs = getComputedStyle(document.body).fontSize;
  const htmlFs = getComputedStyle(document.documentElement).fontSize;
  return { bodyFs, htmlFs, chain };
});
console.log(JSON.stringify(info, null, 2));
await b.close();
