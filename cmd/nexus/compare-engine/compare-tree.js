const { chromium } = require('playwright');

const GODAISY = 'http://localhost:11001/ui/components/button';
const NEXUS = 'http://localhost:11002/ui-components-button.html';

async function getDOMTree(page, url, sel) {
  await page.setViewportSize({ width: 1440, height: 900 });
  await page.goto(url, { waitUntil: 'networkidle', timeout: 15000 });
  await page.waitForTimeout(500);
  return await page.evaluate((selector) => {
    function build(el, depth) {
      if (!el || depth > 30) return null;
      const n = {
        t: el.tagName?.toLowerCase() || '#text',
        c: el.className || '',
        x: '',
        ch: []
      };
      if (el.children?.length === 1 && el.children[0]?.nodeType === 3) {
        n.x = el.textContent.trim().substring(0, 40);
      }
      if (el.children) {
        for (let i = 0; i < el.children.length && n.ch.length < 100; i++) {
          const child = build(el.children[i], depth + 1);
          if (child) n.ch.push(child);
        }
      }
      return n;
    }
    return build(document.querySelector(selector), 0);
  }, sel);
}

function classDiff(c1, c2) {
  const s1 = new Set((c1 || '').split(/\s+/).filter(Boolean));
  const s2 = new Set((c2 || '').split(/\s+/).filter(Boolean));
  const sh = [...s1].filter(x => s2.has(x));
  const o1 = [...s1].filter(x => !s2.has(x));
  const o2 = [...s2].filter(x => !s1.has(x));
  return { shared: sh.length, onlyLeft: o1, onlyRight: o2 };
}

function compare(nx, gd, path, depth) {
  if (!nx || !gd || depth > 25) return [];
  const r = [];
  
  if (nx.t !== gd.t) {
    r.push({ tp: 'tag', p: path, nx: nx.t, gd: gd.t });
    return r;
  }
  
  if (nx.c !== gd.c) {
    const d = classDiff(nx.c, gd.c);
    if (d.onlyLeft.length || d.onlyRight.length) {
      r.push({ tp: 'cls', p: path, tag: nx.t, sh: d.shared, n: d.onlyLeft, g: d.onlyRight });
    }
  }
  
  if (nx.x !== gd.x && (nx.x || gd.x)) {
    r.push({ tp: 'txt', p: path, nx: nx.x, gd: gd.x });
  }
  
  if (nx.ch.length !== gd.ch.length) {
    r.push({ tp: 'cnt', p: path, tag: nx.t, nx: nx.ch.length, gd: gd.ch.length });
  }
  
  const max = Math.max(nx.ch.length, gd.ch.length);
  for (let i = 0; i < max; i++) {
    const cp = (path || '') + '/' + (nx.t || 'x') + '[' + i + ']';
    r.push(...compare(nx.ch[i], gd.ch[i], cp, depth + 1));
  }
  return r;
}

async function main() {
  console.log('=== Element-Level Comparison POC ===\n');
  const b = await chromium.launch({ headless: true });

  const p1 = await b.newPage();
  const p2 = await b.newPage();
  console.log('[1] Loading pages...');
  const [nx, gd] = await Promise.all([
    getDOMTree(p1, NEXUS, '#layout-content'),
    getDOMTree(p2, GODAISY, '#layout-content'),
  ]);
  await p1.close();
  await p2.close();
  console.log(`  Nexus tree: ${JSON.stringify(nx).length}b, GoDaisy: ${JSON.stringify(gd).length}b\n`);

  console.log('[2] Diffing trees...\n');
  const diffs = compare(nx, gd, '', 0);
  
  const byType = {};
  for (const d of diffs) {
    if (!byType[d.tp]) byType[d.tp] = [];
    byType[d.tp].push(d);
  }

  // CLASS diffs — most actionable
  const cls = byType.cls || [];
  console.log(`═══ CLASS DIFFS (${cls.length}) ═══`);
  for (const d of cls.slice(0, 12)) {
    console.log(`  <${d.tag}> at ${d.p || '(root)'}`);
    if (d.n.length) console.log(`    NX-only: ${d.n.join(', ')}`);
    if (d.g.length) console.log(`    GD-only: ${d.g.join(', ')}`);
  }

  // TEXT diffs
  if (byType.txt) {
    console.log(`\n═══ TEXT DIFFS (${byType.txt.length}) ═══`);
    for (const d of byType.txt.slice(0, 8)) {
      console.log(`  ${d.p || '(root)'}`);
      console.log(`    NX: "${d.nx}" | GD: "${d.gd}"`);
    }
  }

  // CHILD COUNT diffs
  if (byType.cnt) {
    console.log(`\n═══ CHILD COUNT DIFFS (${byType.cnt.length}) ═══`);
    for (const d of byType.cnt.slice(0, 5)) {
      console.log(`  <${d.tag}> at ${d.p}: NX=${d.nx} GD=${d.gd} (Δ${d.nx-d.gd})`);
    }
  }

  // TAG mismatches
  if (byType.tag) {
    console.log(`\n═══ TAG MISMATCHES (${byType.tag.length}) ═══`);
    for (const d of byType.tag.slice(0, 5)) {
      console.log(`  ${d.p}: NX=<${d.nx}> GD=<${d.gd}>`);
    }
  }

  console.log(`\n═══════════════════════════════════`);
  console.log(`  Total: ${diffs.length} diffs`);
  console.log(`  Classes: ${cls.length} | Text: ${(byType.txt||[]).length} | Children: ${(byType.cnt||[]).length} | Tags: ${(byType.tag||[]).length}`);

  await b.close();
}
main().catch(e => { console.error(e); process.exit(1); });
