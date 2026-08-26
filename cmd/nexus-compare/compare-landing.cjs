const { chromium } = require('playwright');
const pixelmatch = require('pixelmatch').default;
const { PNG } = require('pngjs');
const fs = require('fs');
const path = require('path');

const NEXUS = 'http://localhost:11002';
const GODAISY = 'http://localhost:11001';
const OUT = path.join(__dirname, 'screenshots', 'landing');

async function run() {
  console.log('=== Landing Page Comparison ===\n');
  if (!fs.existsSync(OUT)) fs.mkdirSync(OUT, { recursive: true });
  const browser = await chromium.launch({ headless: true });

  // Full page screenshots
  console.log('[1/2] Capturing screenshots...');
  
  // Nexus
  let p = await browser.newPage();
  await p.setViewportSize({ width: 1440, height: 900 });
  await p.goto(`${NEXUS}/landing.html`, { waitUntil: 'networkidle', timeout: 15000 });
  await p.waitForTimeout(1000);
  await p.screenshot({ path: path.join(OUT, 'nexus-full.png'), fullPage: true });
  const nxSize = (await p.screenshot({ fullPage: true })).length;
  await p.close();
  console.log(`  nexus: ${(nxSize/1024).toFixed(1)} KB`);

  // Go-daisy
  p = await browser.newPage();
  await p.setViewportSize({ width: 1440, height: 900 });
  await p.goto(`${GODAISY}/landing`, { waitUntil: 'networkidle', timeout: 15000 });
  await p.waitForTimeout(1000);
  await p.screenshot({ path: path.join(OUT, 'godaisy-full.png'), fullPage: true });
  const gdSize = (await p.screenshot({ fullPage: true })).length;
  await p.close();
  console.log(`  godaisy: ${(gdSize/1024).toFixed(1)} KB`);

  // Compare
  console.log('\n[2/2] Comparing...');
  const nx = PNG.sync.read(fs.readFileSync(path.join(OUT, 'nexus-full.png')));
  const gd = PNG.sync.read(fs.readFileSync(path.join(OUT, 'godaisy-full.png')));
  const w = Math.max(nx.width, gd.width);
  const h = Math.max(nx.height, gd.height);
  
  function pad(img, w, h) {
    if (img.width === w && img.height === h) return img;
    const p = new PNG({ width: w, height: h });
    PNG.bitblt(img, p, 0, 0, img.width, img.height, 0, 0);
    return p;
  }
  const nxp = pad(nx, w, h);
  const gdp = pad(gd, w, h);
  const diff = new PNG({ width: w, height: h });
  const mm = pixelmatch(nxp.data, gdp.data, diff.data, w, h, { threshold: 0.1 });
  fs.writeFileSync(path.join(OUT, 'diff.png'), PNG.sync.write(diff));

  const pct = ((mm / (w * h)) * 100).toFixed(2);
  console.log(`  Mismatched: ${mm}px (${pct}%)`);
  console.log(`  Dimensions: ${w}x${h}`);
  console.log(`  Status: ${pct < 5 ? 'GOOD' : pct < 15 ? 'WARN' : 'DIFF'}`);
  console.log(`\nScreenshots: ${OUT}/`);
  
  await browser.close();
}
run().catch(e => { console.error(e); process.exit(1); });
