const { chromium } = require('playwright');
const pixelmatch = require('pixelmatch').default;
const { PNG } = require('pngjs');
const fs = require('fs');
const path = require('path');

const OUT = path.join(__dirname, 'screenshots');
const NEXUS = '/root/nexus-html/html';
const GODAISY = 'http://localhost:11001';

const COMPONENTS = [
  // Sidebar items
  { name: 'sidebar-header', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#_layout-sidebar > div:first-child', nxSelector: '#layout-sidebar > a:first-child', label: 'Sidebar Header' },
  { name: 'sidebar-menu', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#sidebar-menu', nxSelector: '#sidebar-menu', label: 'Sidebar Menu' },
  { name: 'sidebar-bottom', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#_layout-sidebar > div:last-child', nxSelector: '#layout-sidebar > div:last-child', label: 'Sidebar Bottom' },
  // Topbar items
  { name: 'topbar', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#_layout-topbar', nxSelector: '#layout-topbar', label: 'Topbar' },
  // Content - Ecommerce dashboard
  { name: 'stat-cards', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#main-content > div > div:first-of-type', nxSelector: '#layout-content > div:first-of-type > div:first-of-type', label: 'Stat Cards' },
  { name: 'revenue-chart', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#main-content .xl\\:col-span-7 .card', nxSelector: '#layout-content .xl\\:col-span-7 .card', label: 'Revenue Chart Card' },
  { name: 'recent-orders', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#main-content [aria-label=Card]', nxSelector: '#layout-content [aria-label=Card]', label: 'Recent Orders' },
  { name: 'quick-chat', gd: '/dashboards/ecommerce', nx: 'dashboards-ecommerce.html', selector: '#main-content .xl\\:col-span-2 .card', nxSelector: '#layout-content .xl\\:col-span-2 .card:first-of-type', label: 'Quick Chat' },
];

async function takeCompScreenshot(page, url, selector, filepath) {
  await page.goto(url, { waitUntil: 'networkidle', timeout: 15000 });
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.waitForTimeout(1000);
  try {
    const el = await page.$(selector);
    if (!el) { console.log(`    selector not found: ${selector}`); return false; }
    await el.screenshot({ path: filepath });
    return true;
  } catch(e) { console.log(`    error: ${e.message}`); return false; }
}

async function compareImages(baseline, actual, diff) {
  const img1 = PNG.sync.read(fs.readFileSync(baseline));
  const img2 = PNG.sync.read(fs.readFileSync(actual));
  const w = Math.max(img1.width, img2.width);
  const h = Math.max(img1.height, img2.height);
  // Pad smaller image
  function pad(img, w, h) {
    if (img.width === w && img.height === h) return img;
    const p = new PNG({ width: w, height: h });
    PNG.bitblt(img, p, 0, 0, img.width, img.height, 0, 0);
    return p;
  }
  const i1 = pad(img1, w, h);
  const i2 = pad(img2, w, h);
  const diffImg = new PNG({ width: w, height: h });
  const mismatched = pixelmatch(i1.data, i2.data, diffImg.data, w, h, { threshold: 0.1 });
  fs.writeFileSync(diff, PNG.sync.write(diffImg));
  return { mismatched, pct: ((mismatched / (w * h)) * 100).toFixed(2) };
}

async function main() {
  console.log('=== Component-level Pixel Comparison ===\n');
  const browser = await chromium.launch({ headless: true });
  
  const pageDirs = [
    path.join(OUT, 'comp-nx'),
    path.join(OUT, 'comp-gd'),
    path.join(OUT, 'comp-diff'),
  ];
  pageDirs.forEach(d => { if (!fs.existsSync(d)) fs.mkdirSync(d, { recursive: true }); });

  for (const comp of COMPONENTS) {
    console.log(`[${comp.label}]`);
    const nxFile = path.join(OUT, 'comp-nx', `${comp.name}.png`);
    const gdFile = path.join(OUT, 'comp-gd', `${comp.name}.png`);
    const diffFile = path.join(OUT, 'comp-diff', `${comp.name}.png`);

    const page1 = await browser.newPage();
    const nxOk = await takeCompScreenshot(page1, `file://${NEXUS}/${comp.nx}`, comp.nxSelector, nxFile);
    await page1.close();

    const page2 = await browser.newPage();
    const gdOk = await takeCompScreenshot(page2, `${GODAISY}${comp.gd}`, comp.selector, gdFile);
    await page2.close();

    if (nxOk && gdOk) {
      const r = await compareImages(nxFile, gdFile, diffFile);
      const status = r.pct < 1 ? '✓' : r.pct < 5 ? '⚠' : '✗';
      console.log(`  ${status} ${r.mismatched}px (${r.pct}%)`);
    } else {
      console.log(`  ✗ Failed to capture`);
    }
  }
  
  console.log('\nDone.');
  await browser.close();
}
main().catch(e => { console.error(e); process.exit(1); });
