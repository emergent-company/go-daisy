const { chromium } = require('playwright');
const cheerio = require('cheerio');
const fs = require('fs');

function buildTree($el, $) {
  const tag = ($el[0]?.name || '').toLowerCase();
  if (!tag) return null;
  const node = { t: tag, c: ($el.attr('class') || '').trim(), x: '', sig: '', ch: [] };
  node.x = $el.children().length === 0 ? ($el.text() || '').trim().substring(0, 50) : '';
  node.sig = node.t + (node.c ? '.' + node.c.split(/\s+/).slice(0, 3).join('.') : '') + (node.x ? '">' + node.x + '<"' : '');
  $el.children().each(function (i, c) { if (i < 60) { const ch = buildTree($(c), $); if (ch) node.ch.push(ch); } });
  return node;
}
function flatSigs(tree) { const r = []; function w(n) { if (!n) return; r.push({ s: n.sig, t: n.t, c: n.c, x: n.x }); for (const c of n.ch) w(c); } w(tree); return r; }
function lcsMatch(seq1, seq2) {
  const m = seq1.length, n = seq2.length;
  const dp = Array(m+1).fill(null).map(() => Array(n+1).fill(0));
  for (let i = 1; i <= m; i++) for (let j = 1; j <= n; j++) dp[i][j] = seq1[i-1].s === seq2[j-1].s ? dp[i-1][j-1] + 1 : Math.max(dp[i-1][j], dp[i][j-1]);
  const matched1 = new Set(), matched2 = new Set();
  let i = m, j = n;
  while (i > 0 && j > 0) { if (seq1[i-1].s === seq2[j-1].s) { matched1.add(i-1); matched2.add(j-1); i--; j--; } else if (dp[i-1][j] >= dp[i][j-1]) i--; else j--; }
  return { matched: matched1.size, total: Math.max(m, n) };
}

(async () => {
  const b = await chromium.launch({ headless: true });
  const p1 = await b.newPage(); const p2 = await b.newPage();
  await p1.goto('http://localhost:11002/ui-components-badge.html', { waitUntil: 'networkidle' });
  await p2.goto('http://localhost:11001/ui/components/badge', { waitUntil: 'networkidle' });
  const $nx = cheerio.load(await p1.content()); const $gd = cheerio.load(await p2.content());
  await b.close();
  const nxS = flatSigs(buildTree($nx('#layout-content').first(), $nx));
  const gdS = flatSigs(buildTree($gd('#layout-content').first(), $gd));
  const r = lcsMatch(nxS, gdS);
  console.log('Badge: nx=' + nxS.length + ' gd=' + gdS.length + ' matched=' + r.matched + '/' + r.total + ' (' + Math.round(r.matched/r.total*100) + '%)');
})();
