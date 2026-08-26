const { chromium } = require('playwright');
const pixelmatch = require('pixelmatch').default;
const { PNG } = require('pngjs');
const fs = require('fs');
const path = require('path');

const OUT = path.join(__dirname, 'screenshots');
const NEXUS = 'http://localhost:11002';
const GODAISY = 'http://localhost:11001';

const TESTS = [
  // Full pages
  { name: 'ecommerce', gd: '/dashboards/ecommerce', nx: '/dashboards-ecommerce.html', type: 'full' },
  { name: 'crm', gd: '/dashboards/crm', nx: '/dashboards-crm.html', type: 'full' },
  // Components
  { name: 'comp-sidebar-header', gd: '/dashboards/ecommerce', nx: '/dashboards-ecommerce.html', type: 'comp', sel: '#_layout-sidebar > div:first-child', nxSel: '#layout-sidebar > a:first-child' },
  { name: 'comp-sidebar-menu', gd: '/dashboards/ecommerce', nx: '/dashboards-ecommerce.html', type: 'comp', sel: '#sidebar-menu', nxSel: '#sidebar-menu' },
  { name: 'comp-topbar', gd: '/dashboards/ecommerce', nx: '/dashboards-ecommerce.html', type: 'comp', sel: '#_layout-topbar', nxSel: '#layout-topbar' },
  { name: 'comp-statcards', gd: '/dashboards/ecommerce', nx: '/dashboards-ecommerce.html', type: 'comp', sel: '#main-content > div > div:nth-child(2) > div', nxSel: '#layout-content > div:nth-child(2) > div:nth-child(1) > div:first-child' },
  { name: 'comp-revenue', gd: '/dashboards/ecommerce', nx: '/dashboards-ecommerce.html', type: 'comp', sel: '#main-content .xl\\:col-span-7 .card', nxSel: '#layout-content .xl\\:col-span-7 .card' },
];

async function run() {
  console.log('=== Nexus Comparison ===\n');
  const browser = await chromium.launch({ headless: true });
  const results = [];

  for (const t of TESTS) {
    const name = t.name;
    console.log(`[${name}]`);
    const outDirs = ['full', 'comp'];
    outDirs.forEach(d => {
      if (!fs.existsSync(path.join(OUT, d, 'nx'))) fs.mkdirSync(path.join(OUT, d, 'nx'), { recursive: true });
      if (!fs.existsSync(path.join(OUT, d, 'gd'))) fs.mkdirSync(path.join(OUT, d, 'gd'), { recursive: true });
      if (!fs.existsSync(path.join(OUT, d, 'diff'))) fs.mkdirSync(path.join(OUT, d, 'diff'), { recursive: true });
    });

    const dir = t.type === 'full' ? 'full' : 'comp';
    const nxFile = path.join(OUT, dir, 'nx', `${name}.png`);
    const gdFile = path.join(OUT, dir, 'gd', `${name}.png`);
    const diffFile = path.join(OUT, dir, 'diff', `${name}.png`);

    // Nexus capture
    let p = await browser.newPage();
    await p.setViewportSize({ width: 1440, height: 900 });
    try {
      await p.goto(`${NEXUS}${t.nx}`, { waitUntil: 'networkidle', timeout: 15000 });
      await p.waitForTimeout(500);
      if (t.type === 'comp') {
        const el = await p.$(t.nxSel);
        if (el) await el.screenshot({ path: nxFile });
      } else {
        await p.screenshot({ path: nxFile, fullPage: true });
      }
    } catch(e) { console.log(`  nx error: ${e.message}`); }
    await p.close();

    // Go-daisy capture
    p = await browser.newPage();
    await p.setViewportSize({ width: 1440, height: 900 });
    try {
      await p.goto(`${GODAISY}${t.gd}`, { waitUntil: 'networkidle', timeout: 15000 });
      await p.waitForTimeout(500);
      if (t.type === 'comp') {
        const el = await p.$(t.sel);
        if (el) await el.screenshot({ path: gdFile });
      } else {
        await p.screenshot({ path: gdFile, fullPage: true });
      }
    } catch(e) { console.log(`  gd error: ${e.message}`); }
    await p.close();

    // Compare
    if (fs.existsSync(nxFile) && fs.existsSync(gdFile)) {
      const i1 = PNG.sync.read(fs.readFileSync(nxFile));
      const i2 = PNG.sync.read(fs.readFileSync(gdFile));
      const w = Math.max(i1.width, i2.width);
      const h = Math.max(i1.height, i2.height);
      function pad(img, w, h) {
        if (img.width === w && img.height === h) return img;
        const p = new PNG({ width: w, height: h });
        PNG.bitblt(img, p, 0, 0, img.width, img.height, 0, 0);
        return p;
      }
      const pi1 = pad(i1, w, h);
      const pi2 = pad(i2, w, h);
      const diffImg = new PNG({ width: w, height: h });
      const mm = pixelmatch(pi1.data, pi2.data, diffImg.data, w, h, { threshold: 0.1 });
      fs.writeFileSync(diffFile, PNG.sync.write(diffImg));
      const pct = ((mm / (w * h)) * 100).toFixed(2);
      const status = pct < 1 ? '✓' : pct < 5 ? '⚠' : '✗';
      console.log(`  ${status} ${mm}px (${pct}%)  ${w}x${h}`);
      results.push({ name, mismatched: mm, pct, width: w, height: h, status });
    }
  }

  console.log('\n=== Summary ===');
  for (const r of results) {
    console.log(`  ${r.status === '✓' ? '✓' : r.status === '⚠' ? '⚠' : '✗'} ${r.name.padEnd(25)} ${r.pct}%  (${r.mismatched}px)`);
  }
  
  fs.writeFileSync(path.join(__dirname, 'compare-results.json'), JSON.stringify(results, null, 2));
  console.log('\nResults: compare-results.json');
  await browser.close();
}
run().catch(e => { console.error(e); process.exit(1); });
