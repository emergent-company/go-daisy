const { chromium } = require('playwright');
const cheerio = require('cheerio');

function buildTree($el, $) {
  const tag = ($el[0]?.name || '').toLowerCase();
  if (!tag) return null;
  const node = { t: tag, c: ($el.attr('class') || '').trim(), x: '', sig: '', ch: [] };
  node.x = $el.children().length === 0 ? ($el.text() || '').trim().substring(0, 50) : '';
  node.sig = node.t + (node.c ? '.' + node.c.split(/\s+/).slice(0, 3).join('.') : '') + (node.x ? '">' + node.x + '<"' : '');
  $el.children().each((i, c) => { if (i < 80) { const ch = buildTree($(c), $); if (ch) node.ch.push(ch); } });
  return node;
}

function flatSigs(tree) {
  const r = [];
  function w(n) { if (!n) return; r.push({ s: n.sig, t: n.t, c: n.c, x: n.x }); for (const c of n.ch) w(c); } w(tree);
  return r;
}

function lcsMatch(seq1, seq2) {
  const m = seq1.length, n = seq2.length;
  const dp = Array(m+1).fill(null).map(() => Array(n+1).fill(0));
  for (let i = 1; i <= m; i++)
    for (let j = 1; j <= n; j++)
      dp[i][j] = seq1[i-1].s === seq2[j-1].s ? dp[i-1][j-1] + 1 : Math.max(dp[i-1][j], dp[i][j-1]);

  const matched1 = new Set(), matched2 = new Set();
  let i = m, j = n;
  while (i > 0 && j > 0) {
    if (seq1[i-1].s === seq2[j-1].s) { matched1.add(i-1); matched2.add(j-1); i--; j--; }
    else if (dp[i-1][j] >= dp[i][j-1]) i--;
    else j--;
  }

  const un1 = [], un2 = [];
  for (let i = 0; i < m; i++) if (!matched1.has(i)) un1.push(seq1[i]);
  for (let j = 0; j < n; j++) if (!matched2.has(j)) un2.push(seq2[j]);

  const classChanges = [];
  i = m; j = n;
  while (i > 0 && j > 0) {
    if (seq1[i-1].s === seq2[j-1].s) {
      if (seq1[i-1].c !== seq2[j-1].c) {
        const s1 = new Set((seq1[i-1].c||'').split(/\s+/).filter(Boolean));
        const s2 = new Set((seq2[i-1].c||'').split(/\s+/).filter(Boolean));
        const o1 = [...s1].filter(x => !s2.has(x));
        const o2 = [...s2].filter(x => !s1.has(x));
        if (o1.length || o2.length) classChanges.push({ t: seq1[i-1].t, x: seq1[i-1].x, o1: o1.sort().join(' '), o2: o2.sort().join(' ') });
      }
      i--; j--;
    } else if (dp[i-1][j] >= dp[i][j-1]) i--;
    else j--;
  }

  return { matched: matched1.size, total: Math.max(m, n), un1, un2, classChanges };
}

(async () => {
  console.log('=== LCS Semantic Tree Comparison ===\n');
  const b = await chromium.launch({ headless: true });
  const p1 = await b.newPage();
  const p2 = await b.newPage();
  
  await p1.goto('http://localhost:11002/ui-components-button.html', { waitUntil: 'networkidle', timeout: 15000 });
  await p2.goto('http://localhost:11001/ui/components/button', { waitUntil: 'networkidle', timeout: 15000 });
  
  const nxHtml = await p1.content();
  const gdHtml = await p2.content();
  await b.close();
  
  const $nx = cheerio.load(nxHtml);
  const $gd = cheerio.load(gdHtml);
  
  const nxTree = buildTree($nx('#layout-content').first(), $nx);
  const gdTree = buildTree($gd('#layout-content').first(), $gd);
  
  const nxS = flatSigs(nxTree);
  const gdS = flatSigs(gdTree);
  
  console.log('NX: ' + nxS.length + ' nodes | GD: ' + gdS.length + ' nodes\n');
  
  const r = lcsMatch(nxS, gdS);
  const pct = Math.round(r.matched / r.total * 100);
  const missingPenalty = r.un1.length * 3 + r.un2.length * 1;
  const classPenalty = r.classChanges.length * 2;
  const score = Math.max(0, Math.round(100 - missingPenalty - classPenalty));
  
  console.log('═══ Results ═══');
  console.log('LCS Match: ' + pct + '% (' + r.matched + '/' + r.total + ' nodes)');
  console.log('Score: ' + score + '%');
  console.log('Missing from GD: ' + r.un1.length + ' | Extra in GD: ' + r.un2.length + ' | Class diffs: ' + r.classChanges.length + '\n');

  if (r.un1.length) {
    console.log('Elements in NX but MISSING from GD (' + r.un1.length + '):');
    for (const n of r.un1.slice(0, 8))
      console.log('  <' + n.t + '> "' + n.x + '" [' + (n.c||'').substring(0,60) + ']');
  }
  
  if (r.un2.length) {
    console.log('\nElements in GD but NOT in NX (' + r.un2.length + '):');
    for (const n of r.un2.slice(0, 8))
      console.log('  <' + n.t + '> "' + n.x + '" [' + (n.c||'').substring(0,60) + ']');
  }
  
  if (r.classChanges.length) {
    console.log('\nSame position, different CSS classes (' + r.classChanges.length + '):');
    for (const d of r.classChanges.slice(0, 15)) {
      console.log('  <' + d.t + '> "' + d.x + '"');
      if (d.o1) console.log('    NX: +' + d.o1);
      if (d.o2) console.log('    GD: +' + d.o2);
    }
  }
})();
