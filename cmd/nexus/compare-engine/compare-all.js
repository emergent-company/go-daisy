#!/usr/bin/env node
const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

const GD_URL = process.env.GD_URL || 'http://localhost:11001';
const NX_URL = process.env.NX_URL || 'http://localhost:11002';
const OUT = process.env.OUT || '/tmp/compare-report';

const PAGES = [
  { n: 'Login', g: '/auth/login', x: '/auth-login.html' },
  { n: 'Register', g: '/auth/register', x: '/auth-register.html' },
  { n: 'ForgotPwd', g: '/auth/forgot-password', x: '/auth-forgot-password.html' },
  { n: 'ResetPwd', g: '/auth/reset-password', x: '/auth-reset-password.html' },
  { n: 'Ecom Dash', g: '/dashboards/ecommerce', x: '/dashboards-ecommerce.html' },
  { n: 'CRM Dash', g: '/dashboards/crm', x: '/dashboards-crm.html' },
  { n: 'Products L', g: '/apps/ecommerce/products', x: '/apps-ecommerce-products.html' },
  { n: 'Products C', g: '/apps/ecommerce/products/create', x: '/apps-ecommerce-products-create.html' },
  { n: 'Orders L', g: '/apps/ecommerce/orders', x: '/apps-ecommerce-orders.html' },
  { n: 'Order Dtl', g: '/apps/ecommerce/orders/1', x: '/apps-ecommerce-order-details.html' },
  { n: 'Sellers L', g: '/apps/ecommerce/sellers', x: '/apps-ecommerce-sellers.html' },
  { n: 'Sellers C', g: '/apps/ecommerce/sellers/create', x: '/apps-ecommerce-sellers-create.html' },
  { n: 'Sellers E', g: '/apps/ecommerce/sellers/1', x: '/apps-ecommerce-sellers-edit.html' },
  { n: 'Cust L', g: '/apps/ecommerce/customers', x: '/apps-ecommerce-customers.html' },
  { n: 'Cust C', g: '/apps/ecommerce/customers/create', x: '/apps-ecommerce-customers-create.html' },
  { n: 'Cust E', g: '/apps/ecommerce/customers/1', x: '/apps-ecommerce-customers-edit.html' },
  { n: 'Shops L', g: '/apps/ecommerce/shops', x: '/apps-ecommerce-shops.html' },
  { n: 'Shops C', g: '/apps/ecommerce/shops/create', x: '/apps-ecommerce-shops-create.html' },
  { n: 'Shops E', g: '/apps/ecommerce/shops/1', x: '/apps-ecommerce-shops-edit.html' },
  { n: 'Chat', g: '/apps/chat', x: '/apps-chat.html' },
  { n: 'File Mgr', g: '/apps/file-manager', x: '/apps-file-manager.html' },
  { n: 'GenAI Home', g: '/apps/gen-ai/home', x: '/apps-gen-ai-home.html' },
  { n: 'GenAI Cont', g: '/apps/gen-ai/content', x: '/apps-gen-ai-content.html' },
  { n: 'GenAI Img', g: '/apps/gen-ai/image', x: '/apps-gen-ai-image.html' },
  { n: 'GenAI Lib', g: '/apps/gen-ai/library', x: '/apps-gen-ai-library.html' },
  { n: 'Settings', g: '/pages/settings', x: '/pages-settings.html' },
  { n: 'Get Help', g: '/pages/get-help', x: '/pages-get-help.html' },
  { n: 'Landing', g: '/landing', x: '/landing.html' },
  // UI Components (21)
  { n: 'C-Accordion', g: '/ui/components/accordion', x: '/ui-components-accordion.html' },
  { n: 'C-Alert', g: '/ui/components/alert', x: '/ui-components-alert.html' },
  { n: 'C-Avatar', g: '/ui/components/avatar', x: '/ui-components-avatar.html' },
  { n: 'C-Badge', g: '/ui/components/badge', x: '/ui-components-badge.html' },
  { n: 'C-Breadcrumb', g: '/ui/components/breadcrumb', x: '/ui-components-breadcrumb.html' },
  { n: 'C-Button', g: '/ui/components/button', x: '/ui-components-button.html' },
  { n: 'C-Countdown', g: '/ui/components/countdown', x: '/ui-components-countdown.html' },
  { n: 'C-Drawer', g: '/ui/components/drawer', x: '/ui-components-drawer.html' },
  { n: 'C-Dropdown', g: '/ui/components/dropdown', x: '/ui-components-dropdown.html' },
  { n: 'C-Indicator', g: '/ui/components/indicator', x: '/ui-components-indicator.html' },
  { n: 'C-Loading', g: '/ui/components/loading', x: '/ui-components-loading.html' },
  { n: 'C-Menu', g: '/ui/components/menu', x: '/ui-components-menu.html' },
  { n: 'C-Modal', g: '/ui/components/modal', x: '/ui-components-modal.html' },
  { n: 'C-Pagination', g: '/ui/components/pagination', x: '/ui-components-pagination.html' },
  { n: 'C-Progress', g: '/ui/components/progress', x: '/ui-components-progress.html' },
  { n: 'C-Step', g: '/ui/components/step', x: '/ui-components-step.html' },
  { n: 'C-Tab', g: '/ui/components/tab', x: '/ui-components-tab.html' },
  { n: 'C-Table', g: '/ui/components/table', x: '/ui-components-table.html' },
  { n: 'C-Timeline', g: '/ui/components/timeline', x: '/ui-components-timeline.html' },
  { n: 'C-Toast', g: '/ui/components/toast', x: '/ui-components-toast.html' },
  { n: 'C-Tooltip', g: '/ui/components/tooltip', x: '/ui-components-tooltip.html' },
  // UI Forms (12)
  { n: 'F-Checkbox', g: '/ui/forms/checkbox', x: '/ui-forms-checkbox.html' },
  { n: 'F-Fieldset', g: '/ui/forms/fieldset', x: '/ui-forms-fieldset.html' },
  { n: 'F-FileInput', g: '/ui/forms/file-input', x: '/ui-forms-file-input.html' },
  { n: 'F-Input', g: '/ui/forms/input', x: '/ui-forms-input.html' },
  { n: 'F-Label', g: '/ui/forms/label', x: '/ui-forms-label.html' },
  { n: 'F-Radio', g: '/ui/forms/radio', x: '/ui-forms-radio.html' },
  { n: 'F-Range', g: '/ui/forms/range', x: '/ui-forms-range.html' },
  { n: 'F-Rating', g: '/ui/forms/rating', x: '/ui-forms-rating.html' },
  { n: 'F-Select', g: '/ui/forms/select', x: '/ui-forms-select.html' },
  { n: 'F-Textarea', g: '/ui/forms/textarea', x: '/ui-forms-textarea.html' },
  { n: 'F-Toggle', g: '/ui/forms/toggle', x: '/ui-forms-toggle.html' },
  { n: 'F-Validator', g: '/ui/forms/validator', x: '/ui-forms-validator.html' },
  // UI Charts (5)
  { n: 'Ch-Area', g: '/ui/charts/area', x: '/ui-charts-apex-area.html' },
  { n: 'Ch-Bar', g: '/ui/charts/bar', x: '/ui-charts-apex-bar.html' },
  { n: 'Ch-Column', g: '/ui/charts/column', x: '/ui-charts-apex-column.html' },
  { n: 'Ch-Line', g: '/ui/charts/line', x: '/ui-charts-apex-line.html' },
  { n: 'Ch-Pie', g: '/ui/charts/pie', x: '/ui-charts-apex-pie.html' },
];

async function getTree(p, url, sel) {
  try {
    await p.setViewportSize({ width: 1440, height: 900 });
    await p.goto(url, { waitUntil: 'networkidle', timeout: 15000 });
    await p.waitForTimeout(300);
  } catch(e) { return { err: e.message }; }

  return await p.evaluate(s => {
    function build(el, d) {
      if (!el || d > 25) return null;
      const n = { t: el.tagName?.toLowerCase() || '#t', c: el.className || '', x: '', ch: [] };
      if (el.children?.length === 1 && el.children[0]?.nodeType === 3)
        n.x = el.textContent.trim().substring(0, 60);
      if (el.children) for (let i = 0; i < el.children.length && n.ch.length < 80; i++) {
        const c = build(el.children[i], d + 1);
        if (c) n.ch.push(c);
      }
      return n;
    }
    return build(document.querySelector(s), 0);
  }, sel);
}

function clsDiff(a, b) {
  a = (a || '').trim(); b = (b || '').trim();
  if (!a && !b) return { s: 0, l: '', r: '' };
  const s1 = new Set(a.split(/\s+/).filter(Boolean));
  const s2 = new Set(b.split(/\s+/).filter(Boolean));
  return {
    s: [...s1].filter(x => s2.has(x)).length,
    l: [...s1].filter(x => !s2.has(x)).join(' '),
    r: [...s2].filter(x => !s1.has(x)).join(' '),
  };
}

function walk(nx, gd, px, dp) {
  if (!nx || !gd || dp > 25) return [];
  const r = [];
  if (nx.t !== gd.t) { r.push({ tp: 'tag', p: px, l: `<${nx.t}>`, r: `<${gd.t}>` }); return r; }
  if (nx.c !== gd.c) {
    const d = clsDiff(nx.c, gd.c);
    if (d.l || d.r) r.push({ tp: 'cls', p: px, t: nx.t, s: d.s, l: d.l, r: d.r });
  }
  if (nx.x !== gd.x && (nx.x || gd.x)) r.push({ tp: 'txt', p: px, l: nx.x, r: gd.x });
  if (nx.ch.length !== gd.ch.length) r.push({ tp: 'cnt', p: px, t: nx.t, l: nx.ch.length, r: gd.ch.length });
  const max = Math.max(nx.ch.length, gd.ch.length);
  for (let i = 0; i < max; i++)
    r.push(...walk(nx.ch[i], gd.ch[i], (px||'') + '/' + (nx.t||'') + '[' + i + ']', dp + 1));
  return r;
}

function score(d) {
  const t = d.length;
  const w = d.filter(x => x.tp === 'tag').length * 10 +
            d.filter(x => x.tp === 'txt').length * 5 +
            d.filter(x => x.tp === 'cls').length * 2 +
            d.filter(x => x.tp === 'cnt').length * 1;
  return { t, w, s: Math.max(0, 100 - w) };
}

async function main() {
  console.log('Comparison Engine\n');
  const b = await chromium.launch({ headless: true });
  const results = [];

  for (let i = 0; i < PAGES.length; i++) {
    const pg = PAGES[i];
    process.stdout.write(`  ${String(i+1).padStart(2)}/${PAGES.length} ${pg.n.padEnd(15)} `);
    try {
      const [nx, gd] = await Promise.all([
        getTree(await b.newPage(), `${NX_URL}${pg.x}`, '#layout-content'),
        getTree(await b.newPage(), `${GD_URL}${pg.g}`, '#layout-content'),
      ]);
      if (nx.err) { console.log(`SKIP (NX: ${nx.err})`); continue; }
      if (gd.err) { console.log(`SKIP (GD: ${gd.err})`); continue; }
      const diffs = walk(nx, gd, '', 0);
      const s = score(diffs);
      const st = s.s > 95 ? '✓' : s.s > 80 ? '⚠' : '✗';
      console.log(`${st} ${s.s}% (${s.t}d: ${diffs.filter(x=>x.tp==='tag').length}t ${diffs.filter(x=>x.tp==='cls').length}c ${diffs.filter(x=>x.tp==='txt').length}x)`);
      results.push({ name: pg.n, score: s, diffs });
    } catch(e) { console.log(`ERR: ${e.message}`); }
  }
  await b.close();

  // Report
  const total = results.reduce((a,r) => a + r.score.s, 0) / results.length;
  const tbl = results.map(r => {
    const c = r.score.s > 95 ? '#22c55e' : r.score.s > 80 ? '#84cc16' : '#f59e0b';
    return `<tr>
<td><strong>${r.name}</strong></td>
<td style="color:${c};font-weight:bold;font-size:16px">${r.score.s > 95 ? '✓' : r.score.s > 80 ? '⚠' : '✗'}</td>
<td style="color:${c};font-weight:bold">${r.score.s}%</td>
<td>${r.score.t}d</td>
<td>${r.diffs.filter(x=>x.tp==='tag').length}</td>
<td>${r.diffs.filter(x=>x.tp==='cls').length}</td>
<td>${r.diffs.filter(x=>x.tp==='txt').length}</td>
<td>${r.diffs.filter(x=>x.tp==='cnt').length}</td>
</tr>`;
  }).join('');

  const detail = results.map((r, i) => {
    const dhtml = r.diffs.length ? r.diffs.map(d => {
      const icon = { tag: '🔴', cls: '🔵', txt: '🟡', cnt: '⚪' }[d.tp] || '•';
      const bg = { tag: '#fef2f2', cls: '#f0f9ff', txt: '#fefce8', cnt: '#f0fdf4' }[d.tp] || '#fff';
      let desc = `<span style="font-weight:600">${d.p || 'root'}</span>`;
      if (d.tp === 'tag') desc += ` — <b>${d.l}</b> vs <b>${d.r}</b>`;
      if (d.tp === 'cls') desc += ` &lt;${d.t}&gt; NX:<code>${d.l||'none'}</code> GD:<code>${d.r||'none'}</code>`;
      if (d.tp === 'txt') desc += ` — "${d.l}" vs "${d.r}"`;
      if (d.tp === 'cnt') desc += ` &lt;${d.t}&gt; children: NX=${d.l} GD=${d.r}`;
      return `<div style="background:${bg};margin:2px 0;padding:3px 8px;border-radius:4px;font-size:13px">${icon} ${desc}</div>`;
    }).join('') : '<p style="color:#22c55e">No differences</p>';
    return `<details style="background:#fff;padding:12px 16px;margin:8px 0;border-radius:8px;box-shadow:0 1px 2px rgba(0,0,0,.06)"><summary style="cursor:pointer;font-weight:600;font-size:14px">${r.name} — ${r.score.s}%</summary><div style="margin-top:8px">${dhtml}</div></details>`;
  }).join('');

  const html = `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Comparison Report</title>
<style>body{font-family:system-ui;background:#f5f5f5;padding:20px;color:#222}h1{margin-bottom:3px;}
.sub{color:#666;font-size:13px;margin-bottom:15px}.stat{display:inline-block;padding:6px 14px;border-radius:6px;color:#fff;font-weight:600;font-size:13px;margin-right:8px}
table{width:100%;border-collapse:collapse;background:#fff;border-radius:8px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,.08);font-size:12px;margin-bottom:20px}
th,td{padding:8px 12px;text-align:left;border-bottom:1px solid #e5e7eb}th{background:#f9fafb;font-weight:600}
.a{background:#22c55e}.b{background:#84cc16}.c{background:#f59e0b}
code{font-size:11px;padding:1px 4px;background:#e5e7eb;border-radius:3px;word-break:break-all}
</style></head><body>
<h1>Comparison Report</h1><p class="sub">go-daisy vs nexus-html • ${results.length} pages • Avg: ${total.toFixed(1)}%</p>
<div style="margin-bottom:15px">
<span class="stat a">${results.filter(r=>r.score.s>95).length} Excellent</span>
<span class="stat b">${results.filter(r=>r.score.s>80&&r.score.s<=95).length} Good</span>
<span class="stat c">${results.filter(r=>r.score.s<=80).length} Needs Work</span>
</div>
<table><tr><th>Page</th><th></th><th>Score</th><th>Total</th><th>Tags</th><th>Classes</th><th>Text</th><th>Children</th></tr>${tbl}</table>
<h2>Details</h2>${detail}</body></html>`;

  fs.mkdirSync(OUT, { recursive: true });
  fs.writeFileSync(path.join(OUT, 'index.html'), html);
  fs.writeFileSync(path.join(OUT, 'results.json'), JSON.stringify(results, null, 2));
  console.log(`\n  Avg: ${total.toFixed(1)}% | Report: ${OUT}/index.html\n`);
}
main().catch(e => { console.error(e); process.exit(1); });
