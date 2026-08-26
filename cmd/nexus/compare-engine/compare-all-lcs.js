#!/usr/bin/env node
const { chromium } = require('playwright');
const cheerio = require('cheerio');
const fs = require('fs');

const GD_URL = 'http://localhost:11001';
const NX_URL = 'http://localhost:11002';
const OUT = '/tmp/compare-report';

const PAGES = [
  // Dashboards
  { n: 'Ecom Dash', g: '/dashboards/ecommerce', x: '/dashboards-ecommerce.html' },
  { n: 'CRM Dash', g: '/dashboards/crm', x: '/dashboards-crm.html' },
  // Ecommerce
  { n: 'Products', g: '/apps/ecommerce/products', x: '/apps-ecommerce-products.html' },
  { n: 'Prod Create', g: '/apps/ecommerce/products/create', x: '/apps-ecommerce-products-create.html' },
  { n: 'Orders', g: '/apps/ecommerce/orders', x: '/apps-ecommerce-orders.html' },
  { n: 'Sellers', g: '/apps/ecommerce/sellers', x: '/apps-ecommerce-sellers.html' },
  { n: 'Customers', g: '/apps/ecommerce/customers', x: '/apps-ecommerce-customers.html' },
  { n: 'Shops', g: '/apps/ecommerce/shops', x: '/apps-ecommerce-shops.html' },
  // Gen AI
  { n: 'GenAI Home', g: '/apps/gen-ai/home', x: '/apps-gen-ai-home.html' },
  { n: 'GenAI Content', g: '/apps/gen-ai/content', x: '/apps-gen-ai-content.html' },
  { n: 'GenAI Image', g: '/apps/gen-ai/image', x: '/apps-gen-ai-image.html' },
  { n: 'GenAI Library', g: '/apps/gen-ai/library', x: '/apps-gen-ai-library.html' },
  // Apps
  { n: 'Chat', g: '/apps/chat', x: '/apps-chat.html' },
  { n: 'File Mgr', g: '/apps/file-manager', x: '/apps-file-manager.html' },
  // Auth
  { n: 'Login', g: '/login', x: '/auth-login.html' },
  { n: 'Register', g: '/register', x: '/auth-register.html' },
  { n: 'Forgot Pwd', g: '/forgot-password', x: '/auth-forgot-password.html' },
  { n: 'Reset Pwd', g: '/reset-password', x: '/auth-reset-password.html' },
  // Landing
  { n: 'Landing', g: '/landing', x: '/landing.html' },
  // Pages
  { n: 'Settings', g: '/pages/settings', x: '/pages-settings.html' },
  { n: 'Get Help', g: '/pages/get-help', x: '/pages-get-help.html' },
  // UI Components
  { n: 'Button', g: '/ui/components/button', x: '/ui-components-button.html' },
  { n: 'Badge', g: '/ui/components/badge', x: '/ui-components-badge.html' },
  { n: 'Alert', g: '/ui/components/alert', x: '/ui-components-alert.html' },
  { n: 'Avatar', g: '/ui/components/avatar', x: '/ui-components-avatar.html' },
  { n: 'Accordion', g: '/ui/components/accordion', x: '/ui-components-accordion.html' },
  { n: 'Breadcrumb', g: '/ui/components/breadcrumb', x: '/ui-components-breadcrumb.html' },
  { n: 'Countdown', g: '/ui/components/countdown', x: '/ui-components-countdown.html' },
  { n: 'Drawer', g: '/ui/components/drawer', x: '/ui-components-drawer.html' },
  { n: 'Dropdown', g: '/ui/components/dropdown', x: '/ui-components-dropdown.html' },
  { n: 'Indicator', g: '/ui/components/indicator', x: '/ui-components-indicator.html' },
  { n: 'Loading', g: '/ui/components/loading', x: '/ui-components-loading.html' },
  { n: 'Menu', g: '/ui/components/menu', x: '/ui-components-menu.html' },
  { n: 'Modal', g: '/ui/components/modal', x: '/ui-components-modal.html' },
  { n: 'Pagination', g: '/ui/components/pagination', x: '/ui-components-pagination.html' },
  { n: 'Progress', g: '/ui/components/progress', x: '/ui-components-progress.html' },
  { n: 'Step', g: '/ui/components/step', x: '/ui-components-step.html' },
  { n: 'Tab', g: '/ui/components/tab', x: '/ui-components-tab.html' },
  { n: 'Table', g: '/ui/components/table', x: '/ui-components-table.html' },
  { n: 'Timeline', g: '/ui/components/timeline', x: '/ui-components-timeline.html' },
  { n: 'Toast', g: '/ui/components/toast', x: '/ui-components-toast.html' },
  { n: 'Tooltip', g: '/ui/components/tooltip', x: '/ui-components-tooltip.html' },
  // UI Forms
  { n: 'Checkbox', g: '/ui/forms/checkbox', x: '/ui-forms-checkbox.html' },
  { n: 'Input', g: '/ui/forms/input', x: '/ui-forms-input.html' },
  { n: 'Toggle', g: '/ui/forms/toggle', x: '/ui-forms-toggle.html' },
  { n: 'Select', g: '/ui/forms/select', x: '/ui-forms-select.html' },
  { n: 'File Input', g: '/ui/forms/file-input', x: '/ui-forms-file-input.html' },
  { n: 'Textarea', g: '/ui/forms/textarea', x: '/ui-forms-textarea.html' },
  { n: 'Radio', g: '/ui/forms/radio', x: '/ui-forms-radio.html' },
  { n: 'Range', g: '/ui/forms/range', x: '/ui-forms-range.html' },
  { n: 'Rating', g: '/ui/forms/rating', x: '/ui-forms-rating.html' },
  { n: 'Label', g: '/ui/forms/label', x: '/ui-forms-label.html' },
  { n: 'Fieldset', g: '/ui/forms/fieldset', x: '/ui-forms-fieldset.html' },
  { n: 'Validator', g: '/ui/forms/validator', x: '/ui-forms-validator.html' },
];

function buildTree($el, $) {
  const tag = ($el[0]?.name || '').toLowerCase();
  if (!tag) return null;
  const node = { t: tag, c: ($el.attr('class') || '').trim(), x: '', sig: '', ch: [] };
  node.x = $el.children().length === 0 ? ($el.text() || '').replace(/\s+/g, ' ').trim().substring(0, 50) : '';
  node.sig = node.t + (node.c ? '.' + node.c.split(/\s+/).slice(0, 3).join('.') : '') + (node.x ? '">' + node.x + '<"' : '');
  $el.children().each(function (i, c) { if (i < 60) { const ch = buildTree($(c), $); if (ch) node.ch.push(ch); } });
  return node;
}

function flatSigs(tree) {
  const r = [];
  function w(n) { if (!n) return; r.push({ s: n.sig, t: n.t, c: n.c, x: n.x }); for (const c of n.ch) w(c); }
  w(tree);
  return r;
}

function lcsMatch(seq1, seq2) {
  const m = seq1.length, n = seq2.length;
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
  for (let i = 0; i < m; i++) if (!matched1.has(i)) un1.push(seq1[i]);
  for (let j = 0; j < n; j++) if (!matched2.has(j)) un2.push(seq2[j]);

  const classChanges = [];
  i = m; j = n;
  while (i > 0 && j > 0) {
    if (seq1[i-1].s === seq2[j-1].s) {
      if (seq1[i-1].c !== seq2[j-1].c) {
        const s1 = new Set((seq1[i-1].c || '').split(/\s+/).filter(Boolean));
        const s2 = new Set((seq2[j-1].c || '').split(/\s+/).filter(Boolean));
        const o1 = [].filter.call(s1, function(x) { return !s2.has(x); });
        const o2 = [].filter.call(s2, function(x) { return !s1.has(x); });
        if (o1.length || o2.length) classChanges.push({ t: seq1[i-1].t, x: seq1[i-1].x, o1: o1.sort().join(' '), o2: o2.sort().join(' ') });
      }
      i--; j--;
    } else if (dp[i-1][j] >= dp[i][j-1]) i--;
    else j--;
  }

  return { matched: matched1.size, total: Math.max(m, n), un1, un2, classChanges };
}

(async function() {
  console.log('=== LCS Comparison — ' + PAGES.length + ' pages ===\n');
  const b = await chromium.launch({ headless: true });
  const results = [];

  for (let idx = 0; idx < PAGES.length; idx++) {
    const pg = PAGES[idx];
    const label = String(idx+1).padStart(2) + '/' + PAGES.length + ' ' + pg.n;
    process.stdout.write('  ' + label.padEnd(28) + ' ');
    
    try {
      const p1 = await b.newPage();
      const p2 = await b.newPage();
      await p1.goto(NX_URL + pg.x, { waitUntil: 'networkidle', timeout: 15000 });
      await p2.goto(GD_URL + pg.g, { waitUntil: 'networkidle', timeout: 15000 });
      
      const $nx = cheerio.load(await p1.content());
      const $gd = cheerio.load(await p2.content());
      await p1.close();
      await p2.close();

      const nxRoot = $nx('#layout-content').first();
      const gdRoot = $gd('#layout-content').first();
      
      if (!nxRoot.length || !gdRoot.length) {
        console.log('SKIP (no #layout-content)');
        continue;
      }

      const nxS = flatSigs(buildTree(nxRoot, $nx));
      const gdS = flatSigs(buildTree(gdRoot, $gd));
      const r = lcsMatch(nxS, gdS);
      
      const pct = Math.round(r.matched / r.total * 100);
      const score = Math.max(0, pct - r.classChanges.length*2 - Math.abs(r.un1.length - r.un2.length)*0.5);
      const status = score > 90 ? '✓' : score > 70 ? '⚠' : '✗';
      
      console.log(status + ' ' + String(score).padStart(3) + '% (nx:' + nxS.length + ' gd:' + gdS.length + ' matched:' + r.matched + ' missing:' + r.un1.length + ' extra:' + r.un2.length + ' cls:' + r.classChanges.length + ')');
      
      results.push({ name: pg.n, nxNodes: nxS.length, gdNodes: gdS.length, matched: r.matched, missingFromGd: r.un1.length, extraInGd: r.un2.length, classChanges: r.classChanges.length, score });
    } catch(e) {
      console.log('ERR: ' + e.message);
    }
  }
  
  await b.close();

  // Summary
  const avg = results.reduce(function(a, r) { return a + r.score; }, 0) / results.length;
  console.log('\n' + '═'.repeat(60));
  console.log('  Average score: ' + avg.toFixed(1) + '% | ' + results.length + ' pages');
  
  const ordered = [].concat(results).sort(function(a, b) { return b.score - a.score; });
  for (const r of ordered) {
    const s = r.score > 90 ? '✓' : r.score > 70 ? '⚠' : '✗';
    console.log('  ' + s + ' ' + r.name.padEnd(20) + r.score + '%');
  }
  
  fs.mkdirSync(OUT, { recursive: true });
  fs.writeFileSync(OUT + '/lcs-results.json', JSON.stringify({ results, avg, time: new Date().toISOString() }, null, 2));
  console.log('\n  Report: ' + OUT + '/lcs-results.json');
})();
