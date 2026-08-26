#!/usr/bin/env python3
"""Enhanced Visual Comparison - Compares DaisyUI/Tailwind classes between two implementations"""
import urllib.request, re, json, sys, os
from html.parser import HTMLParser
from collections import Counter

GODAISY = "http://localhost:11001"
NEXUS_HTML_DIR = "/root/nexus-html/html"
OUTPUT_DIR = "/tmp/visual-compare"

PAGES = [
    ("dash-ecommerce", "/dashboards/ecommerce", "dashboards-ecommerce.html"),
    ("dash-crm", "/dashboards/crm", "dashboards-crm.html"),
    ("eco-products", "/apps/ecommerce/products", "apps-ecommerce-products.html"),
    ("app-chat", "/apps/chat", "apps-chat.html"),
    ("page-settings", "/pages/settings", "pages-settings.html"),
    ("auth-login", "/auth/login", "auth-login.html"),
    ("landing", "/landing", "landing.html"),
]

class ClassExtractor(HTMLParser):
    """Extract all class attributes and their element tags"""
    def __init__(self):
        super().__init__()
        self.classes = Counter()
        self.elements = Counter()
        self.class_elements = {}  # class -> [tags]
    
    def handle_starttag(self, tag, attrs):
        self.elements[tag] += 1
        for name, val in attrs:
            if name == 'class':
                for cls in val.split():
                    self.classes[cls] += 1
                    if cls not in self.class_elements:
                        self.class_elements[cls] = set()
                    self.class_elements[cls].add(tag)

def fetch(url):
    try:
        resp = urllib.request.urlopen(url, timeout=10)
        return resp.read().decode(errors='replace')
    except:
        return ""

def fetch_nexus(path):
    try:
        with open(f"{NEXUS_HTML_DIR}/{path}") as f:
            return f.read()
    except:
        return ""

def daisyui_classes(classes):
    """Filter to DaisyUI/Tailwind relevant classes"""
    daisy_patterns = [
        r'^(btn|badge|card|alert|toast|modal|drawer|dropdown|menu|table|input|select|textarea|checkbox|radio|toggle|range|rating|steps?|tabs?|collapse|tooltip|indicator|divider|join|mask|loading|progress|radial-progress|countdown|kbd|stat|chat|breadcrumb|pagination|timeline|swap|carousel|diff|hero|mockup|fieldset|label|file-input|avatar|link|stats)',
        r'^(bg-|text-|border-|rounded-|shadow-|p-|m-|gap-|flex|grid-|w-|h-|size-|items-|justify-|space-|overflow-|z-|sticky|relative|absolute|hidden|inline|block|grow|shrink|truncate|font-|leading-|tracking-|uppercase|capitalize|transition|duration|opacity|hover:|group-|peer-|dark:|sm:|md:|lg:|xl:|2xl:|max-)',
    ]
    result = set()
    for cls in classes:
        for pattern in daisy_patterns:
            if re.match(pattern, cls):
                result.add(cls)
                break
    return result

def compare_page(label, gd_path, nx_path):
    gd_html = fetch(f"{GODAISY}{gd_path}")
    nx_html = fetch_nexus(nx_path)
    
    gd = ClassExtractor()
    nx = ClassExtractor()
    gd.feed(gd_html)
    nx.feed(nx_html)
    
    gd_daisy = daisyui_classes(gd.classes.keys())
    nx_daisy = daisyui_classes(nx.classes.keys())
    
    shared = gd_daisy & nx_daisy
    only_gd = gd_daisy - nx_daisy
    only_nx = nx_daisy - gd_daisy
    
    # Compare element counts
    element_diffs = []
    all_tags = set(gd.elements.keys()) | set(nx.elements.keys())
    for tag in sorted(all_tags):
        gd_count = gd.elements.get(tag, 0)
        nx_count = nx.elements.get(tag, 0)
        if gd_count != nx_count:
            element_diffs.append((tag, gd_count, nx_count))
    
    return {
        'label': label,
        'gd_size': len(gd_html),
        'nx_size': len(nx_html),
        'gd_total_classes': len(gd.classes),
        'nx_total_classes': len(nx.classes),
        'gd_daisy_classes': len(gd_daisy),
        'nx_daisy_classes': len(nx_daisy),
        'shared_classes': len(shared),
        'only_gd': sorted(only_gd)[:30],
        'only_nx': sorted(only_nx)[:30],
        'match_pct': round(len(shared) / max(len(gd_daisy | nx_daisy), 1) * 100, 1),
        'element_diffs': element_diffs[:20],
        'gd_top_classes': gd.classes.most_common(20),
        'nx_top_classes': nx.classes.most_common(20),
    }

def generate_report(results):
    html = """<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Class Comparison Report</title>
<style>
body{font-family:system-ui;background:#f5f5f5;padding:20px;color:#333}
h1{font-size:20px;margin-bottom:5px}.sub{color:#666;font-size:13px;margin-bottom:20px}
.summary{display:flex;gap:12px;margin-bottom:20px;flex-wrap:wrap}
.stat{padding:8px 16px;border-radius:6px;color:#fff;font-weight:600;font-size:13px}
.stat-ok{background:#22c55e}.stat-warn{background:#f59e0b}
.card{background:#fff;border-radius:8px;padding:16px;margin-bottom:14px;box-shadow:0 1px 3px rgba(0,0,0,.08)}
.card h3{font-size:15px;margin:0 0 8px 0}.card .meta{font-size:12px;color:#888;margin-bottom:10px}
.bar{display:inline-block;height:8px;border-radius:4px;vertical-align:middle}
.bar-bg{display:inline-block;width:150px;height:8px;background:#e5e7eb;border-radius:4px;vertical-align:middle}
.diff-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px;font-size:12px}
.diff-box{padding:8px;border-radius:4px;max-height:200px;overflow-y:auto}
.only-gd{background:#fef2f2;border:1px solid #fecaca}
.only-nx{background:#f0fdf4;border:1px solid #bbf7d0}
.class-tag{display:inline-block;background:#e5e7eb;padding:1px 6px;border-radius:3px;margin:2px;font-size:11px;font-family:monospace}
table{width:100%;border-collapse:collapse;font-size:12px}
th,td{padding:6px 10px;text-align:left;border-bottom:1px solid #e5e7eb}
th{background:#f9fafb;font-weight:600}
.match-high{color:#22c55e}.match-mid{color:#f59e0b}.match-low{color:#ef4444}
</style></head><body>
<h1>DaisyUI/Tailwind Class Comparison</h1>
<div class="sub">go-daisy vs nexus-html — """ + str(len(results)) + """ pages</div>
"""

    # Summary
    avg_match = sum(r['match_pct'] for r in results) / len(results)
    total_gd = sum(r['gd_daisy_classes'] for r in results)
    total_nx = sum(r['nx_daisy_classes'] for r in results)
    html += f"""<div class="summary">
<div class="stat stat-ok">Avg match: {avg_match:.1f}%</div>
<div class="stat stat-warn">GD classes: {total_gd}</div>
<div class="stat stat-warn">NX classes: {total_nx}</div>
</div>"""

    for r in results:
        match_color = "match-high" if r['match_pct'] > 80 else "match-mid" if r['match_pct'] > 50 else "match-low"
        pct = r['match_pct']
        bar_w = max(min(int(pct * 1.5), 150), 10)
        
        html += f"""<div class="card">
<h3>{r['label']} <span class="{match_color}">{pct}% class match</span></h3>
<div class="meta">
  Size: {r['gd_size']:,}b (GD) vs {r['nx_size']:,}b (NX) |
  Daisy classes: {r['gd_daisy_classes']} (GD) vs {r['nx_daisy_classes']} (NX) |
  Shared: {r['shared_classes']}
  <div class="bar-bg"><div class="bar" style="width:{bar_w}px;background:{"#22c55e" if pct > 80 else "#f59e0b" if pct > 50 else "#ef4444"}"></div></div>
</div>
<div class="diff-grid">
<div><strong>Only in go-daisy ({len(r['only_gd'])}):</strong>
<div class="diff-box only-gd">"""
        for c in r['only_gd']:
            html += f'<span class="class-tag">{c}</span>'
        html += """</div></div>
<div><strong>Only in nexus ({len(r['only_nx'])}):</strong>
<div class="diff-box only-nx">"""
        for c in r['only_nx']:
            html += f'<span class="class-tag">{c}</span>'
        html += """</div></div>
</div>"""
        
        if r['element_diffs']:
            html += """<div style="margin-top:10px"><strong>Element count differences:</strong>
<table><tr><th>Tag</th><th>GD</th><th>NX</th><th>Diff</th></tr>"""
            for tag, gd_c, nx_c in r['element_diffs'][:10]:
                diff = gd_c - nx_c
                html += f"<tr><td>&lt;{tag}&gt;</td><td>{gd_c}</td><td>{nx_c}</td><td>{diff:+d}</td></tr>"
            html += "</table></div>"
        
        html += "</div>"

    html += "</body></html>"
    return html

# Run
print("=== DaisyUI Class Comparison ===\n")
results = []
for label, gd_path, nx_path in PAGES:
    print(f"  {label}...", end=" ", flush=True)
    r = compare_page(label, gd_path, nx_path)
    results.append(r)
    print(f"{r['match_pct']}% match ({r['shared_classes']} shared, {len(r['only_gd'])} gd-only, {len(r['only_nx'])} nx-only)")

# Generate report
report = generate_report(results)
report_path = f"{OUTPUT_DIR}/class-report.html"
with open(report_path, 'w') as f:
    f.write(report)
json_path = f"{OUTPUT_DIR}/class-report.json"
with open(json_path, 'w') as f:
    json.dump(results, f, indent=2)

print(f"\n  JSON: {json_path}")
print(f"  HTML: {report_path}")
print(f"\n  Avg class match: {sum(r['match_pct'] for r in results)/len(results):.1f}%")
