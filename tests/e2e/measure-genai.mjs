import { chromium } from '@playwright/test';

async function measure(browser, url) {
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  await page.goto(url, { waitUntil: 'networkidle' });
  await page.waitForTimeout(400);
  const data = await page.evaluate(() => {
    const cs = (el) => el ? getComputedStyle(el) : null;
    const bg = (sel) => {
      const el = document.querySelector(sel);
      return el ? cs(el).backgroundColor : 'MISSING:'+sel;
    };
    const cardTitle = document.querySelector('.card .card-title, .card-title');
    const suggestionTitle = document.querySelector('.card .card-body p.font-medium');
    const searchInput = document.querySelector('input[type="search"], label.input input');
    return {
      bodyBg: cs(document.body).backgroundColor,
      contentBg: bg('#layout-content'),
      cardBodyPad: (() => { const el = document.querySelector('.card .card-body'); return el ? cs(el).padding : 'none'; })(),
      cardTitle: cardTitle ? { fs: cs(cardTitle).fontSize, fw: cs(cardTitle).fontWeight, text: cardTitle.textContent.trim().slice(0,20) } : null,
      suggestionTitle: suggestionTitle ? { fs: cs(suggestionTitle).fontSize, fw: cs(suggestionTitle).fontWeight, text: suggestionTitle.textContent.trim().slice(0,20) } : null,
      searchPlaceholder: searchInput ? searchInput.getAttribute('placeholder') : null,
      searchAria: searchInput ? searchInput.getAttribute('aria-label') : null,
      tableCellWhiteSpace: (() => { const td = document.querySelector('td'); return td ? cs(td).whiteSpace : 'none'; })(),
    };
  });
  await page.close();
  return data;
}

const browser = await chromium.launch();
const targets = {
  'home-ref': 'http://localhost:11002/apps-gen-ai-home.html',
  'home-got': 'http://localhost:11001/apps/gen-ai/home',
  'content-ref': 'http://localhost:11002/apps-gen-ai-content.html',
  'content-got': 'http://localhost:11001/apps/gen-ai/content',
  'image-ref': 'http://localhost:11002/apps-gen-ai-image.html',
  'image-got': 'http://localhost:11001/apps/gen-ai/image',
  'library-ref': 'http://localhost:11002/apps-gen-ai-library.html',
  'library-got': 'http://localhost:11001/apps/gen-ai/library',
  'products-create-got': 'http://localhost:11001/apps/ecommerce/products/create',
  'customers-got': 'http://localhost:11001/apps/ecommerce/customers',
};
for (const [name, url] of Object.entries(targets)) {
  const d = await measure(browser, url);
  console.log(JSON.stringify({ name, ...d }));
}
await browser.close();
