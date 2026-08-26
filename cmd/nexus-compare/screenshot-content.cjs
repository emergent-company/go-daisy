const { chromium } = require('playwright');
const pixelmatch = require('pixelmatch').default;
const { PNG } = require('pngjs');
const fs = require('fs');
const path = require('path');

const OUT = '/tmp/button-compare';

async function run() {
  const browser = await chromium.launch({ headless: true });
  console.log('=== Content-Area Button Comparison ===\n');

  // NEXUS content area
  const p1 = await browser.newPage();
  await p1.setViewportSize({ width: 1440, height: 900 });
  await p1.goto('http://localhost:11002/ui-components-button.html', { waitUntil: 'networkidle', timeout: 15000 });
  await p1.waitForTimeout(500);
  const nxEl = await p1.$('#layout-content');
  await nxEl.screenshot({ path: path.join(OUT, 'content-nexus.png') });
  const nxBox = await nxEl.boundingBox();
  console.log(`  Nexus content: ${Math.round(nxBox.width)}x${Math.round(nxBox.height)}px`);
  await p1.close();

  // GO-DAISY content area
  const p2 = await browser.newPage();
  await p2.setViewportSize({ width: 1440, height: 900 });
  await p2.goto('http://localhost:11001/ui/components/button', { waitUntil: 'networkidle', timeout: 15000 });
  await p2.waitForTimeout(500);
  const gdEl = await p2.$('#layout-content');
  await gdEl.screenshot({ path: path.join(OUT, 'content-godaisy.png') });
  const gdBox = await gdEl.boundingBox();
  console.log(`  GoDaisy content: ${Math.round(gdBox.width)}x${Math.round(gdBox.height)}px`);
  await p2.close();

  // Diff content only
  const nx = PNG.sync.read(fs.readFileSync(path.join(OUT, 'content-nexus.png')));
  const gd = PNG.sync.read(fs.readFileSync(path.join(OUT, 'content-godaisy.png')));
  const w = Math.max(nx.width, gd.width);
  const h = Math.max(nx.height, gd.height);
  function pad(img, w, h) {
    if (img.width === w && img.height === h) return img;
    const p = new PNG({ width: w, height: h });
    PNG.bitblt(img, p, 0, 0, img.width, img.height, 0, 0);
    return p;
  }
  const diff = new PNG({ width: w, height: h });
  const mm = pixelmatch(pad(nx,w,h).data, pad(gd,w,h).data, diff.data, w, h, { threshold: 0.1, diffColor: [255,0,0] });
  fs.writeFileSync(path.join(OUT, 'content-diff.png'), PNG.sync.write(diff));
  console.log(`  Content diff: ${mm}px (${((mm/(w*h))*100).toFixed(2)}%)`);
  
  console.log(`\nScreenshots: ${OUT}/`);
  console.log(`  content-nexus.png   — nexus button demo (no shell)`);
  console.log(`  content-godaisy.png — go-daisy button demo (no shell)`);
  console.log(`  content-diff.png    — red = mismatch`);
  await browser.close();
}
run().catch(e => { console.error(e); process.exit(1); });
