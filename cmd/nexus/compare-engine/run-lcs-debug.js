const { chromium } = require('playwright');
const cheerio = require('cheerio');

function buildTree($el, $) {
  const tag = ($el[0]?.name || '').toLowerCase();
  console.log('  buildTree: tag=' + tag + ' class="' + ($el.attr('class')||'') + '" children=' + $el.children().length);
  if (!tag) return null;
  const node = { t: tag, c: ($el.attr('class') || '').trim(), x: '', sig: '', ch: [] };
  node.x = $el.children().length === 0 ? ($el.text() || '').trim().substring(0, 50) : '';
  node.sig = node.t + (node.c ? '.' + node.c.split(/\s+/).slice(0, 3).join('.') : '') + (node.x ? '">' + node.x + '<"' : '');
  $el.children().each((i, c) => { if (i < 80) { const ch = buildTree($(c), $); if (ch) node.ch.push(ch); } });
  return node;
}

function flatSigs(tree) {
  const r = [];
  (function w(n) { if (!n) return; r.push({ s: n.sig, t: n.t, c: n.c, x: n.x }); for (const c of n.ch) w(c); })();
  return r;
}

(async () => {
  const b = await chromium.launch({ headless: true });
  const p = await b.newPage();
  await p.goto('http://localhost:11001/ui/components/button', { waitUntil: 'networkidle', timeout: 15000 });
  const html = await p.content();
  await b.close();
  
  const $ = cheerio.load(html);
  console.log('HTML size: ' + html.length);
  console.log('layout-content found: ' + $('#layout-content').length);
  
  const root = $('#layout-content').first();
  console.log('root tag: ' + root[0]?.name);
  
  const tree = buildTree(root, $);
  console.log('tree produced: ' + (tree ? tree.t : 'NULL'));
  
  const sigs = flatSigs(tree);
  console.log('flatSigs: ' + sigs.length + ' entries');
})();
