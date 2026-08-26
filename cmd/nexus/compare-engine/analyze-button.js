const { chromium } = require('playwright');

async function getTree(p, url) {
  await p.setViewportSize({ width: 1440, height: 900 });
  await p.goto(url, { waitUntil: 'networkidle', timeout: 15000 });
  await p.waitForTimeout(500);
  return await p.evaluate(() => {
    function build(el, d) {
      if (!el || d > 25) return null;
      const n = { t: el.tagName?.toLowerCase() || '#t', c: el.className || '', x: '', ch: [] };
      if (el.children?.length === 1 && el.children[0]?.nodeType === 3)
        n.x = el.textContent.trim().substring(0, 80);
      if (el.children) for (let i = 0; i < el.children.length && n.ch.length < 100; i++) {
        const c = build(el.children[i], d + 1);
        if (c) n.ch.push(c);
      }
      return n;
    }
    return build(document.querySelector('#layout-content'), 0);
  });
}

function clsDiff(a, b) {
  a = (a || '').trim(); b = (b || '').trim();
  if (!a && !b) return '';
  const s1 = new Set(a.split(/\s+/).filter(Boolean));
  const s2 = new Set(b.split(/\s+/).filter(Boolean));
  const o1 = [...s1].filter(x => !s2.has(x));
  const o2 = [...s2].filter(x => !s1.has(x));
  const parts = [];
  if (o1.length) parts.push(`NX-only: ${o1.join(' ')}`);
  if (o2.length) parts.push(`GD-only: ${o2.join(' ')}`);
  return parts.join(' | ');
}

function walk(nx, gd, px, dp) {
  if (!nx || !gd || dp > 20) return [];
  const r = [];
  const p = px || '(root)';
  
  if (nx.t !== gd.t) {
    r.push(`🔴 TAG: ${p} — NX=<${nx.t}> GD=<${gd.t}>`);
    return r;
  }
  
  if (nx.c !== gd.c) {
    const d = clsDiff(nx.c, gd.c);
    if (d) r.push(`🔵 CLASS: ${p} <${nx.t}> — ${d}`);
  }
  
  if (nx.x !== gd.x && (nx.x || gd.x))
    r.push(`🟡 TEXT: ${p} — NX="${nx.x}" vs GD="${gd.x}"`);
    
  if (nx.ch.length !== gd.ch.length)
    r.push(`⚪ COUNT: ${p} <${nx.t}> — NX has ${nx.ch.length} children, GD has ${gd.ch.length}`);
  
  const max = Math.max(nx.ch.length, gd.ch.length);
  for (let i = 0; i < max; i++)
    r.push(...walk(nx.ch[i], gd.ch[i], p + '/' + (nx.t||'') + '[' + i + ']', dp + 1));
  
  return r;
}

async function main() {
  console.log('=== Button Page Element-Level Analysis ===\n');
  const b = await chromium.launch({ headless: true });
  
  const [nx, gd] = await Promise.all([
    getTree(await b.newPage(), 'http://localhost:11002/ui-components-button.html'),
    getTree(await b.newPage(), 'http://localhost:11001/ui/components/button'),
  ]);
  
  await b.close();
  
  const diffs = walk(nx, gd, '', 0);
  
  console.log(`${diffs.length} differences found:\n`);
  for (const d of diffs) console.log(`  ${d}`);
  
  // Summary
  const tags = diffs.filter(x => x.startsWith('🔴')).length;
  const classes = diffs.filter(x => x.startsWith('🔵')).length;
  const texts = diffs.filter(x => x.startsWith('🟡')).length;
  const counts = diffs.filter(x => x.startsWith('⚪')).length;
  
  console.log(`\n  Tags: ${tags} | Classes: ${classes} | Text: ${texts} | Children: ${counts}`);
}
main().catch(e => console.error(e.message));
