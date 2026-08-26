const { chromium } = require('playwright');
const pixelmatch = require('pixelmatch').default;
const { PNG } = require('pngjs');
const fs = require('fs');
const path = require('path');

const BASELINES = path.join(__dirname, 'screenshots', 'baselines');
const GODAISY = path.join(__dirname, 'screenshots', 'godaisy');
const DIFFS = path.join(__dirname, 'screenshots', 'diffs');

const PAGES = [
  { name: 'ecommerce', gdRoute: '/dashboards/ecommerce', nexusFile: 'dashboards-ecommerce.html' },
  { name: 'crm', gdRoute: '/dashboards/crm', nexusFile: 'dashboards-crm.html' },
  { name: 'products', gdRoute: '/apps/ecommerce/products', nexusFile: 'apps-ecommerce-products.html' },
  { name: 'chat', gdRoute: '/apps/chat', nexusFile: 'apps-chat.html' },
  { name: 'file-manager', gdRoute: '/apps/file-manager', nexusFile: 'apps-file-manager.html' },
];

const NEXUS_HTML_ROOT = '/root/nexus-html/html';
const GODAISY_BASE = 'http://localhost:11001';
const VIEWPORT = { width: 1440, height: 900 };

async function takeScreenshot(page, url, filepath) {
  await page.goto(url, { waitUntil: 'networkidle', timeout: 15000 });
  await page.setViewportSize(VIEWPORT);
  await page.waitForTimeout(500);
  await page.screenshot({ path: filepath, fullPage: true });
}

async function compareImages(baseline, actual, diff) {
  const img1 = PNG.sync.read(fs.readFileSync(baseline));
  const img2 = PNG.sync.read(fs.readFileSync(actual));
  const { width, height } = img1;
  const diffImg = new PNG({ width, height });
  const mismatched = pixelmatch(img1.data, img2.data, diffImg.data, width, height, {
    threshold: 0.1,
    diffColor: [255, 0, 0],
  });
  fs.writeFileSync(diff, PNG.sync.write(diffImg));
  const pct = ((mismatched / (width * height)) * 100).toFixed(2);
  return { mismatched, pct };
}

async function main() {
  console.log('=== Nexus Pixel Comparison (Playwright) ===\n');
  const browser = await chromium.launch({ headless: true });
  
  try {
    console.log('[1/3] Capturing nexus-html baselines...');
    const page = await browser.newPage();
    for (const p of PAGES) {
      const url = `file://${NEXUS_HTML_ROOT}/${p.nexusFile}`;
      await takeScreenshot(page, url, path.join(BASELINES, `${p.name}.png`));
    }
    await page.close();

    console.log('\n[2/3] Capturing go-daisy screenshots...');
    const page2 = await browser.newPage();
    for (const p of PAGES) {
      await takeScreenshot(page2, `${GODAISY_BASE}${p.gdRoute}`, path.join(GODAISY, `${p.name}.png`));
    }
    await page2.close();

    console.log('\n[3/3] Comparing...');
    for (const p of PAGES) {
      const baseline = path.join(BASELINES, `${p.name}.png`);
      const actual = path.join(GODAISY, `${p.name}.png`);
      const diff = path.join(DIFFS, `${p.name}.png`);
      if (!fs.existsSync(baseline) || !fs.existsSync(actual)) {
        console.log(`  ${p.name}: SKIP`);
        continue;
      }
      const r = await compareImages(baseline, actual, diff);
      const status = r.pct < 1 ? 'GOOD' : r.pct < 5 ? 'WARN' : 'DIFF';
      console.log(`  ${p.name}: ${r.mismatched}px (${r.pct}%) ${status}`);
    }
  } finally {
    await browser.close();
  }
}

main().catch(e => { console.error(e); process.exit(1); });
