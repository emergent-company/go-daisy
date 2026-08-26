#!/usr/bin/env node
// Show exact node mismatches between NX and GD for a single page
const { chromium } = require('playwright');
const cheerio = require('cheerio');

const [,, nxUrl, gdUrl] = process.argv;
if (!nxUrl || !gdUrl) { console.log('Usage: node diff-page.js <NX-render-url> <GD-render-url>'); process.exit(1); }

function buildTree($el, $) {
  const tag = ($el[0]?.name || '').toLowerCase();
  if (!tag) return null;
  const node = { t: tag, c: ($el.attr('class') || '').trim(), x: '', sig: '', ch: [] };
  node.x = $el.children().length === 0 ? ($el.text() || '').replace(/\s+/g, ' ').trim().substring(0, 50) : '';
  node.sig = node.t + (node.c ? '.' + node.c.split(/\s+/).slice(0, 3).join('.') : '') + (node.x ? '>"' + node.x + '"<' : '');
  $el.children().each(function (i, c) { if (i < 60) { const ch = buildTree($(c), $); if (ch) node.ch.push(ch); } });
  return node;
}

function flatSigs(tree) {
  const r = []; function w(n) { if (!n) return; r.push({ s: n.sig, t: n.t, c: n.c, x: n.x }); for (const c of n.ch) w(c); } w(tree); return r;
}

function lcsMatch(seq1, seq2) {
  const m = seq1.length, n = seq2.length;
  if (m === 0 || n === 0) return { matched: new Set(), un1: seq1.map((_,i)=>i), un2: seq2.map((_,i)=>i) };
  const dp = Array(m+1).fill(null).map(function() { return Array(n+1).fill(0); });
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
  for (let i = 0; i < m; i++) if (!matched1.has(i)) un1.push(i);
  for (let j = 0; j < n; j++) if (!matched2.has(j)) un2.push(j);
  return { matched: matched1, un1, un2 };
}

(async () => {
  const b = await chromium.launch({ headless: true });
  const p1 = await b.newPage(), p2 = await b.newPage();
  await p1.goto('http://localhost:11002' + nxUrl, { waitUntil: 'networkidle', timeout: 10000 });
  await p2.goto('http://localhost:11001' + gdUrl, { waitUntil: 'networkidle', timeout: 10000 });
  const $nx = cheerio.load(await p1.content()), $gd = cheerio.load(await p2.content());
  await p1.close(); await p2.close(); await b.close();

  const nxT = buildTree($nx('#layout-content').first(), $nx), gdT = buildTree($gd('#layout-content').first(), $gd);
  const nxS = flatSigs(nxT), gdS = flatSigs(gdT);
  const r = lcsMatch(nxS, gdS);

  console.log(`NX: ${nxS.length} nodes, GD: ${gdS.length} nodes, Matched: ${r.matched.size}`);
  console.log(`Score: ${Math.round(r.matched.size / Math.max(nxS.length, gdS.length) * 100)}%`);

  if (r.un1.length > 0) {
    console.log(`\n--- Missing from GD (${r.un1.length} NX nodes not in GD) ---`);
    const grouped = {};
    for (const i of r.un1) {
      const n = nxS[i];
      const key = `<${n.t}> class="${n.c.substring(0,50)}" text="${(n.x||'').substring(0,30)}"`;
      grouped[key] = (grouped[key] || 0) + 1;
    }
    const sorted = Object.entries(grouped).sort((a,b) => b[1] - a[1]);
    for (const [k, v] of sorted.slice(0, 30)) console.log(`  ${v}x ${k}`);
    if (sorted.length > 30) console.log(`  ... and ${sorted.length - 30} more unique patterns`);
  }

  if (r.un2.length > 0) {
    console.log(`\n--- Extra in GD (${r.un2.length} GD nodes not in NX) ---`);
    const grouped = {};
    for (const j of r.un2) {
      const n = gdS[j];
      const key = `<${n.t}> class="${n.c.substring(0,50)}" text="${(n.x||'').substring(0,30)}"`;
      grouped[key] = (grouped[key] || 0) + 1;
    }
    const sorted = Object.entries(grouped).sort((a,b) => b[1] - a[1]);
    for (const [k, v] of sorted.slice(0, 20)) console.log(`  ${v}x ${k}`);
    if (sorted.length > 20) console.log(`  ... and ${sorted.length - 20} more unique patterns`);
  }

  // Class differences on matched nodes
  let classChanges = [];
  for (let i = 0; i < nxS.length; i++) {
    if (r.matched.has(i)) {
      const gdIdx = [...r.matched][[...r.matched].indexOf(i)];
      // Find the corresponding GD node
      for (let j = 0; j < gdS.length; j++) {
        if (gdS[j].s === nxS[i].s && !classChanges.some(c => c.j === j)) {
          if (gdS[j].c !== nxS[i].c) {
            classChanges.push({ i, j, nxC: nxS[i].c, gdC: gdS[j].c, tag: nxS[i].t });
          }
          break;
        }
      }
    }
  }
  if (classChanges.length > 0) {
    console.log(`\n--- Class differences on matched nodes (${classChanges.length}) ---`);
    for (const cc of classChanges.slice(0, 15)) {
      console.log(`  <${cc.tag}> NX="${cc.nxC}" vs GD="${cc.gdC}"`);
    }
  }
})();
