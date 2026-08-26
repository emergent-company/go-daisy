const { chromium } = require('playwright');
const pixelmatch = require('pixelmatch').default;
const { PNG } = require('pngjs');
const fs = require('fs');
const path = require('path');

const OUT = '/tmp/button-compare';
if (!fs.existsSync(OUT)) fs.mkdirSync(OUT, { recursive: true });

async function run() {
  console.log('=== Button Page Screenshot Comparison ===\n');
  const browser = await chromium.launch({ headless: true });

  // Capture NEXUS button page
  console.log('[1] Capturing nexus-html button page...');
  const p1 = await browser.newPage();
  await p1.setViewportSize({ width: 1440, height: 900 });
  await p1.goto('http://localhost:11002/ui-components-button.html', { waitUntil: 'networkidle', timeout: 15000 });
  await p1.waitForTimeout(500);
  await p1.screenshot({ path: path.join(OUT, 'nexus-button.png'), fullPage: true });
  await p1.close();
  console.log('  ✓ nexus-button.png');

  // Capture GO-DAISY button page
  console.log('[2] Capturing go-daisy button page...');
  const p2 = await browser.newPage();
  await p2.setViewportSize({ width: 1440, height: 900 });
  await p2.goto('http://localhost:11001/ui/components/button', { waitUntil: 'networkidle', timeout: 15000 });
  await p2.waitForTimeout(500);
  await p2.screenshot({ path: path.join(OUT, 'godaisy-button.png'), fullPage: true });
  await p2.close();
  console.log('  ✓ godaisy-button.png');

  // Pixel diff
  console.log('[3] Computing diff...');
  const nx = PNG.sync.read(fs.readFileSync(path.join(OUT, 'nexus-button.png')));
  const gd = PNG.sync.read(fs.readFileSync(path.join(OUT, 'godaisy-button.png')));
  
  const w = Math.max(nx.width, gd.width);
  const h = Math.max(nx.height, gd.height);
  
  function pad(img, w, h) {
    if (img.width === w && img.height === h) return img;
    const p = new PNG({ width: w, height: h });
    PNG.bitblt(img, p, 0, 0, img.width, img.height, 0, 0);
    return p;
  }
  
  const nx_p = pad(nx, w, h);
  const gd_p = pad(gd, w, h);
  const diff = new PNG({ width: w, height: h });
  const mm = pixelmatch(nx_p.data, gd_p.data, diff.data, w, h, { 
    threshold: 0.1,
    diffColor: [255, 0, 0],
    diffColorAlt: [255, 255, 0],
  });
  fs.writeFileSync(path.join(OUT, 'diff.png'), PNG.sync.write(diff));
  
  const pct = ((mm / (w * h)) * 100).toFixed(2);
  console.log(`  Mismatched: ${mm}px (${pct}%)`);
  console.log(`  Dimensions: ${w}x${h}`);
  console.log(`\nScreenshots saved to: ${OUT}/`);
  console.log(`  nexus-button.png    - nexus reference`);
  console.log(`  godaisy-button.png  - go-daisy implementation`);
  console.log(`  diff.png            - pixel difference (red pixels = mismatch)`);
  
  await browser.close();
}
run().catch(e => { console.error(e); process.exit(1); });
