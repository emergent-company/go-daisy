const cheerio = require('cheerio');
const { chromium } = require('playwright');
const fs = require('fs');

function getSig($el, $) {
  const tag = ($el[0]?.name || '').toLowerCase();
  const cls = ($el.attr('class') || '').trim();
  const id = ($el.attr('id') || '').trim();
  const children = $el.children().length;
  const text = children === 0 ? ($el.text() || '').trim().substring(0, 40) : '';
  return tag + (id ? '#' + id : '') + (cls ? '.' + cls.split(/\s+/).slice(0, 3).join('.') : '') + (text ? '"><' + text + '>' : '');
}

function buildTree($el, $, depth) {
  const tag = ($el[0]?.name || '').toLowerCase();
  if (!tag) return null;
  const node = {
    t: tag,
    c: ($el.attr('class') || '').trim(),
    x: $el.children().length === 0 ? ($el.text() || '').trim().substring(0, 50) : '',
    sig: getSig($el, $),
    ch: [],
  };
  $el.children().each((i, c) => {
    if (i < 100) {
      const child = buildTree($(c), $, depth + 1);
      if (child) node.ch.push(child);
    }
  });
  return node;
}

function walk(tree, prefix) {
  if (!tree) return [];
  const p = prefix || '/';
  const entries = [{ p, s: tree.sig, t: tree.t, c: tree.c, x: tree.x }];
  tree.ch.forEach((c, i) => entries.push(...walk(c, p + '/' + tree.t + '[' + i + ']')));
  return entries;
}

(async () => {
  const b = await chromium.launch({ headless: true });
  const p1 = await b.newPage();
  const p2 = await b.newPage();
  
  // Both at once
  await p1.goto('http://localhost:11002/ui-components-button.html', { waitUntil: 'networkidle', timeout: 15000 });
  await p2.goto('http://localhost:11001/ui/components/button', { waitUntil: 'networkidle', timeout: 15000 });
  
  const nxHtml = await p1.content();
  const gdHtml = await p2.content();
  await b.close();
  
  const $nx = cheerio.load(nxHtml);
  const $gd = cheerio.load(gdHtml);
  
  const nxTree = buildTree($nx('#layout-content').first(), $nx, 0);
  const gdTree = buildTree($gd('#layout-content').first(), $gd, 0);
  
  const nxWalk = walk(nxTree, '');
  const gdWalk = walk(gdTree, '');
  
  console.log('NX nodes: ' + nxWalk.length + ' | GD nodes: ' + gdWalk.length + '\n');
  
  // Show first 15 nodes from each
  console.log('First 10 NX nodes:');
  for (const n of nxWalk.slice(0, 10))
    console.log('  ' + n.t.padEnd(10) + (n.c.substring(0, 60) || '(no class)').padEnd(65) + (n.x ? '"' + n.x + '"' : ''));
  
  console.log('\nFirst 10 GD nodes:');
  for (const n of gdWalk.slice(0, 10))
    console.log('  ' + n.t.padEnd(10) + (n.c.substring(0, 60) || '(no class)').padEnd(65) + (n.x ? '"' + n.x + '"' : ''));
  
})();
