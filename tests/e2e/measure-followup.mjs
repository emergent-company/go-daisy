import { chromium } from '@playwright/test';

async function measure(browser, url) {
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  await page.goto(url, { waitUntil: 'networkidle' });
  await page.waitForTimeout(500);
  const data = await page.evaluate(() => {
    const cs = (el) => el ? getComputedStyle(el) : null;
    const body = document.body;
    // content pane: for nexus it's #layout-content; for reference it may be #layout-main or .size-full
    const contentEl = document.querySelector('#layout-content') || document.querySelector('#layout-main') || document.querySelector('div.size-full');
    const cardBody = document.querySelector('.card .card-body, .card-body');
    const cardTitle = document.querySelector('.card-title');
    const suggestionTitle = [...document.querySelectorAll('.card-body p.font-medium')].find(p => /Blog Post|Watercolor|Email|Social/.test(p.textContent));
    const searchInput = [...document.querySelectorAll('input[type="search"]')].find(i => (i.getAttribute('placeholder')||'').startsWith('Search along'));
    const row = document.querySelector('tbody tr');
    const contentCell = [...document.querySelectorAll('td')].find(td => /Create stunning|Summarize lengthy/.test(td.textContent));
    return {
      bodyBg: cs(body)?.backgroundColor,
      contentBg: cs(contentEl)?.backgroundColor,
      contentEl: contentEl ? (contentEl.id || contentEl.className).toString().slice(0,40) : 'none',
      cardBodyPad: cardBody ? cs(cardBody).padding : null,
      cardTitle: cardTitle ? { fs: cs(cardTitle).fontSize, fw: cs(cardTitle).fontWeight, text: cardTitle.textContent.trim().slice(0,18) } : null,
      suggestionTitle: suggestionTitle ? { fs: cs(suggestionTitle).fontSize, fw: cs(suggestionTitle).fontWeight, text: suggestionTitle.textContent.trim().slice(0,18) } : null,
      search: searchInput ? { ph: searchInput.getAttribute('placeholder'), aria: searchInput.getAttribute('aria-label') } : null,
      rowWhiteSpace: row ? cs(row).whiteSpace : null,
      rowHeight: row ? row.getBoundingClientRect().height : null,
      contentCellWhiteSpace: contentCell ? cs(contentCell).whiteSpace : null,
      contentCellPWhiteSpace: contentCell ? cs(contentCell.querySelector('p')).whiteSpace : null,
    };
  });
  await page.close();
  return data;
}

const browser = await chromium.launch();
const targets = [
  ['home-ref',  'http://localhost:11002/apps-gen-ai-home.html'],
  ['home-got',  'http://localhost:11001/apps/gen-ai/home'],
  ['content-ref','http://localhost:11002/apps-gen-ai-content.html'],
  ['content-got','http://localhost:11001/apps/gen-ai/content'],
  ['image-ref', 'http://localhost:11002/apps-gen-ai-image.html'],
  ['image-got', 'http://localhost:11001/apps/gen-ai/image'],
  ['library-ref','http://localhost:11002/apps-gen-ai-library.html'],
  ['library-got','http://localhost:11001/apps/gen-ai/library'],
  ['products-got','http://localhost:11001/apps/ecommerce/products/create'],
  ['customers-got','http://localhost:11001/apps/ecommerce/customers'],
];
for (const [name, url] of targets) {
  const d = await measure(browser, url);
  console.log(name + ' :: ' + JSON.stringify(d));
}
await browser.close();
