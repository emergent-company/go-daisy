#!/usr/bin/env python3
"""Content-Aware Comparison — compares #layout-content between implementations"""
import urllib.request, re, json, sys, os
from collections import Counter

GODAISY = "http://localhost:11001"
NEXUS_HTML = "/root/nexus-html/html"

def fetch_gd(path):
    try:
        return urllib.request.urlopen(f"{GODAISY}{path}", timeout=10).read().decode(errors='replace')
    except: return ""

def fetch_nx(file):
    try:
        with open(f"{NEXUS_HTML}/{file}") as f: return f.read()
    except: return ""

def extract_content(html):
    """Extract #layout-content area"""
    idx = html.find('id="layout-content"')
    if idx < 0: return "", ""
    end = html.find('<div class="drawer', idx)
    if end < 0:
        end = html.find('</div></div></div>', idx + 500)
    if end < 0: end = len(html)
    return html[idx:end], f"[{idx}:{end}]"

def elem_counts(html):
    counts = Counter()
    for m in re.finditer(r'<(\w+)[\s>]', html):
        counts[m.group(1)] += 1
    return counts

def class_set(html):
    classes = set()
    for m in re.finditer(r'class="([^"]+)"', html):
        for c in m.group(1).split():
            classes.add(c)
    return classes

def compare_page(label, gd_path, nx_file):
    gd_html = fetch_gd(gd_path)
    nx_html = fetch_nx(nx_file)
    
    gd_content, _ = extract_content(gd_html)
    nx_content, _ = extract_content(nx_html)
    
    gd_cls = class_set(gd_content)
    nx_cls = class_set(nx_content)
    gd_elems = elem_counts(gd_content)
    nx_elems = elem_counts(nx_content)
    
    shared = gd_cls & nx_cls
    match = round(len(shared) / max(len(gd_cls | nx_cls), 1) * 100, 1)
    
    # Element diffs
    elem_diff = {}
    all_tags = set(gd_elems.keys()) | set(nx_elems.keys())
    for tag in all_tags:
        g = gd_elems.get(tag, 0)
        n = nx_elems.get(tag, 0)
        if g != n:
            elem_diff[tag] = (g, n, g - n)
    
    return {
        'label': label,
        'gd_content_size': len(gd_content),
        'nx_content_size': len(nx_content),
        'gd_classes': len(gd_cls),
        'nx_classes': len(nx_cls),
        'shared_classes': len(shared),
        'match': match,
        'only_gd': sorted(gd_cls - nx_cls),
        'only_nx': sorted(nx_cls - gd_cls),
        'elem_diff': {k: v for k, v in sorted(elem_diff.items(), key=lambda x: -abs(x[1][2]))[:10]},
    }

# Test single page
if len(sys.argv) > 1:
    label, gd_path, nx_file = sys.argv[1], sys.argv[2], sys.argv[3]
    r = compare_page(label, gd_path, nx_file)
    print(f"=== {label} ===")
    print(f"Content: {r['gd_content_size']}b (GD) vs {r['nx_content_size']}b (NX)")
    print(f"Classes: {r['gd_classes']} (GD) vs {r['nx_classes']} (NX) — {r['match']}% match")
    print(f"Shared: {r['shared_classes']}")
    print(f"\nMissing from GD ({len(r['only_nx'])}):")
    for c in r['only_nx']:
        print(f"  - {c}")
    print(f"\nMissing from NX ({len(r['only_gd'])}):")
    for c in r['only_gd']:
        print(f"  + {c}")
    print(f"\nElement count diffs:")
    for tag, (g, n, d) in r['elem_diff'].items():
        print(f"  <{tag}>: {g} vs {n} ({d:+d})")
else:
    # Run all
    pages = [
        ("button", "/ui/components/button", "ui-components-button.html"),
        ("badge", "/ui/components/badge", "ui-components-badge.html"),
        ("alert", "/ui/components/alert", "ui-components-alert.html"),
        ("checkbox", "/ui/forms/checkbox", "ui-forms-checkbox.html"),
        ("input", "/ui/forms/input", "ui-forms-input.html"),
        ("select", "/ui/forms/select", "ui-forms-select.html"),
        ("dash-ecommerce", "/dashboards/ecommerce", "dashboards-ecommerce.html"),
        ("dash-crm", "/dashboards/crm", "dashboards-crm.html"),
    ]
    for label, gd, nx in pages:
        r = compare_page(label, gd, nx)
        bar = "█" * int(r['match']/2) + "░" * (50 - int(r['match']/2))
        print(f"{label:<20} {r['gd_content_size']:>6}b vs {r['nx_content_size']:>6}b | {r['match']:>5.1f}% {bar}")
