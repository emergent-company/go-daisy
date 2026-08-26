#!/usr/bin/env bash
# Universal Visual Comparison Tool
# Usage: ./compare.sh [--quick|--full] [--pages-only] [--output-dir DIR]
# Compares two web implementations page-by-page with structural + visual diffing.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUTPUT_DIR="${OUTPUT_DIR:-/tmp/visual-compare}"
REPORT_DIR="$OUTPUT_DIR/report"

# === Config ===
GODAISY_URL="${GODAISY_URL:-http://localhost:11001}"
NEXUS_REF_URL="${NEXUS_REF_URL:-http://localhost:11002}"
NEXUS_HTML_DIR="${NEXUS_HTML_DIR:-/root/nexus-html/html}"

# List of pages to compare: "label|gd_path|nx_path"
PAGES=(
  "landing|/landing|landing.html"
  "auth-login|/auth/login|auth-login.html"
  "auth-register|/auth/register|auth-register.html"
  "auth-forgot|/auth/forgot-password|auth-forgot-password.html"
  "auth-reset|/auth/reset-password|auth-reset-password.html"
  "dash-ecommerce|/dashboards/ecommerce|dashboards-ecommerce.html"
  "dash-crm|/dashboards/crm|dashboards-crm.html"
  "eco-products|/apps/ecommerce/products|apps-ecommerce-products.html"
  "eco-orders|/apps/ecommerce/orders|apps-ecommerce-orders.html"
  "eco-sellers|/apps/ecommerce/sellers|apps-ecommerce-sellers.html"
  "eco-customers|/apps/ecommerce/customers|apps-ecommerce-customers.html"
  "eco-shops|/apps/ecommerce/shops|apps-ecommerce-shops.html"
  "eco-products-create|/apps/ecommerce/products/create|apps-ecommerce-products-create.html"
  "eco-products-edit|/apps/ecommerce/products/1|apps-ecommerce-products-edit.html"
  "eco-order-detail|/apps/ecommerce/orders/1|apps-ecommerce-order-details.html"
  "eco-sellers-create|/apps/ecommerce/sellers/create|apps-ecommerce-sellers-create.html"
  "eco-sellers-edit|/apps/ecommerce/sellers/1|apps-ecommerce-sellers-edit.html"
  "eco-customers-create|/apps/ecommerce/customers/create|apps-ecommerce-customers-create.html"
  "eco-customers-edit|/apps/ecommerce/customers/1|apps-ecommerce-customers-edit.html"
  "eco-shops-create|/apps/ecommerce/shops/create|apps-ecommerce-shops-create.html"
  "eco-shops-edit|/apps/ecommerce/shops/1|apps-ecommerce-shops-edit.html"
  "genai-home|/apps/gen-ai/home|apps-gen-ai-home.html"
  "genai-content|/apps/gen-ai/content|apps-gen-ai-content.html"
  "genai-image|/apps/gen-ai/image|apps-gen-ai-image.html"
  "genai-library|/apps/gen-ai/library|apps-gen-ai-library.html"
  "app-chat|/apps/chat|apps-chat.html"
  "app-file-manager|/apps/file-manager|apps-file-manager.html"
  "page-settings|/pages/settings|pages-settings.html"
  "page-help|/pages/get-help|pages-get-help.html"
)

mkdir -p "$OUTPUT_DIR" "$REPORT_DIR"

# === Colors ===
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}=================================================${NC}"
echo -e "${CYAN}  Visual Comparison Tool v1.0${NC}"
echo -e "${CYAN}=================================================${NC}"
echo ""

# === Stage 1: Health Check ===
echo -e "${YELLOW}[1/5] Health check${NC}"
GD_OK=$(curl -s -o /dev/null -w "%{http_code}" "$GODAISY_URL/dashboards/ecommerce" 2>/dev/null || echo "000")
NX_OK=$(curl -s -o /dev/null -w "%{http_code}" "$NEXUS_REF_URL/dashboards-ecommerce.html" 2>/dev/null || echo "000")

if [ "$GD_OK" != "200" ]; then
  echo -e "  ${RED}✗${NC} go-daisy server not running at $GODAISY_URL"
  exit 1
fi
if [ "$NX_OK" != "200" ]; then
  echo -e "  ${RED}✗${NC} nexus-html reference not running at $NEXUS_REF_URL"
  echo "    Start with: cd $NEXUS_HTML_DIR && python3 -m http.server 11002 &"
  exit 1
fi
echo -e "  ${GREEN}✓${NC} go-daisy: $GODAISY_URL"
echo -e "  ${GREEN}✓${NC} nexus-ref: $NEXUS_REF_URL"

# === Stage 2: Page Crawling ===
echo ""
echo -e "${YELLOW}[2/5] Crawling pages${NC}"

GD_PAGES=()
NX_PAGES=()
MISSING=()
TOTAL=${#PAGES[@]}
PASS=0

for page in "${PAGES[@]}"; do
  IFS='|' read -r label gd_path nx_path <<< "$page"
  
  gd_code=$(curl -s -o /dev/null -w "%{http_code}" "$GODAISY_URL$gd_path" 2>/dev/null || echo "000")
  nx_code=$(curl -s -o /dev/null -w "%{http_code}" "$NEXUS_REF_URL/$nx_path" 2>/dev/null || echo "000")
  
  if [ "$gd_code" = "200" ] && [ "$nx_code" = "200" ]; then
    echo -e "  ${GREEN}✓${NC} $label"
    PASS=$((PASS + 1))
  elif [ "$gd_code" != "200" ]; then
    echo -e "  ${RED}✗${NC} $label (go-daisy: $gd_code)"
    MISSING+=("$label:$gd_path")
  else
    echo -e "  ${RED}✗${NC} $label (nexus-ref: $nx_code)"
    MISSING+=("$label:$gd_path")
  fi
done

echo "  ${PASS}/${TOTAL} pages accessible"

# === Stage 3: Structural Comparison ===
echo ""
echo -e "${YELLOW}[3/5] Structural comparison${NC}"

RESULTS_JSON="$OUTPUT_DIR/results.json"
echo "[" > "$RESULTS_JSON"

FIRST=true
for page in "${PAGES[@]}"; do
  IFS='|' read -r label gd_path nx_path <<< "$page"
  
  gd_html=$(curl -s "$GODAISY_URL$gd_path" 2>/dev/null || echo "")
  nx_html=$(cat "$NEXUS_HTML_DIR/$nx_path" 2>/dev/null || echo "")
  
  gd_size=${#gd_html}
  nx_size=${#nx_html}
  
  # Quick structural diff: count matching CSS classes
  gd_classes=$(echo "$gd_html" | grep -oP 'class="[^"]*"' | sort -u | wc -l)
  nx_classes=$(echo "$nx_html" | grep -oP 'class="[^"]*"' | sort -u | wc -l)
  
  # Check key elements
  checks=()
  for term in "id=\"layout-content\"" "table" "card" "btn" "badge"; do
    gd_has=$(echo "$gd_html" | grep -c "$term" || echo 0)
    nx_has=$(echo "$nx_html" | grep -c "$term" || echo 0)
    if [ "$gd_has" != "0" ] && [ "$nx_has" != "0" ]; then
      checks+=("$term:✓")
    elif [ "$gd_has" != "0" ]; then
      checks+=("$term:gd-only")
    elif [ "$nx_has" != "0" ]; then
      checks+=("$term:nx-only")
    fi
  done
  
  # Build JSON
  if [ "$FIRST" = "false" ]; then echo "," >> "$RESULTS_JSON"; fi
  FIRST=false
  
  cat >> "$RESULTS_JSON" << JSONENTRY
  {
    "label": "$label",
    "path": "$gd_path",
    "gd_size": $gd_size,
    "nx_size": $nx_size,
    "gd_classes": $gd_classes,
    "nx_classes": $nx_classes,
    "checks": [$(printf '"%s",' "${checks[@]}" | sed 's/,$//')]
  }
JSONENTRY

  size_diff=$((gd_size - nx_size))
  diff_pct=$(( size_diff * 100 / (nx_size + 1) ))
  
  if [ $diff_pct -lt 10 ] && [ $diff_pct -gt -10 ]; then
    echo -e "  ${GREEN}✓${NC} $label (${gd_size}b vs ${nx_size}b, ${diff_pct}% diff)"
  else
    echo -e "  ${YELLOW}⚠${NC} $label (${gd_size}b vs ${nx_size}b, ${diff_pct}% diff)"
  fi
done

echo "]" >> "$RESULTS_JSON"
echo "  Results: $RESULTS_JSON"

# === Stage 4: Visual Comparison ===
echo ""
echo -e "${YELLOW}[4/5] Visual comparison${NC}"

VISUAL_DIR="$OUTPUT_DIR/visual"
mkdir -p "$VISUAL_DIR"

if command -v node &>/dev/null && [ -f "$SCRIPT_DIR/compare-full.cjs" ]; then
  node "$SCRIPT_DIR/compare-full.cjs" 2>&1 | head -20
  echo "  Visual results: $VISUAL_DIR"
else
  echo -e "  ${YELLOW}⚠${NC} Node.js visual comparison not available"
  echo "    Create compare-full.cjs with Playwright + pixelmatch for pixel-level comparison"
fi

# === Stage 5: Generate HTML Report ===
echo ""
echo -e "${YELLOW}[5/5] Generating report${NC}"

REPORT_HTML="$REPORT_DIR/index.html"
python3 << PYEOF > "$REPORT_HTML"
import json, sys, os

with open("$RESULTS_JSON") as f:
    results = json.load(f)

total = len(results)
pass_count = sum(1 for r in results if r["gd_size"] > 0 and r["nx_size"] > 0)
ok_count = sum(1 for r in results if abs(r["gd_size"] - r["nx_size"]) < r["nx_size"] * 0.3)

html = """<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Visual Comparison Report</title>
<style>
body{font-family:system-ui,sans-serif;background:#f5f5f5;padding:20px;color:#333}
h1{margin-bottom:5px}.sub{color:#666;font-size:14px;margin-bottom:20px}
.summary{display:flex;gap:15px;margin-bottom:25px}
.stat{padding:10px 20px;border-radius:8px;color:#fff;font-weight:bold;font-size:14px}
.stat-ok{background:#22c55e}.stat-warn{background:#f59e0b}.stat-err{background:#ef4444}
table{width:100%;border-collapse:collapse;background:#fff;border-radius:8px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,.1)}
th,td{padding:10px 15px;text-align:left;font-size:13px}
th{background:#e5e7eb;font-weight:600}
tr{border-bottom:1px solid #e5e7eb}
.pass{color:#22c55e}.warn{color:#f59e0b}.fail{color:#ef4444}
.bar{height:8px;border-radius:4px;min-width:80px;display:inline-block}
.bar-bg{background:#e5e7eb;border-radius:4px;width:120px;display:inline-block}
</style>
</head>
<body>
<h1>Visual Comparison Report</h1>
<div class="sub">go-daisy vs nexus-html | """ + str(total) + """ pages compared</div>
<div class="summary">
  <div class="stat stat-ok">""" + str(pass_count) + """ accessible</div>
  <div class="stat stat-warn">""" + str(ok_count) + """ similar size</div>
</div>
<table>
<tr><th>Page</th><th>go-daisy Size</th><th>nexus Size</th><th>Size Diff</th><th>Classes (GD/NX)</th><th>Elements</th></tr>
"""

for r in results:
    diff = r["gd_size"] - r["nx_size"]
    pct = abs(diff) * 100 // max(r["nx_size"], 1)
    cls = "pass" if pct < 30 else "warn" if pct < 60 else "fail"
    bar_w = min(100, pct) if pct < 101 else 100
    bar_color = "#22c55e" if pct < 30 else "#f59e0b" if pct < 60 else "#ef4444"
    checks = ", ".join(r["checks"]) if r["checks"] else "—"
    html += f"""<tr>
      <td><strong>{r['label']}</strong></td>
      <td>{r['gd_size']:,}b</td>
      <td>{r['nx_size']:,}b</td>
      <td class="{cls}">{diff:+,}b ({pct}%) <div class="bar-bg"><div class="bar" style="width:{bar_w}px;background:{bar_color}"></div></div></td>
      <td>{r['gd_classes']} / {r['nx_classes']}</td>
      <td style="font-size:11px">{checks}</td>
    </tr>"""

html += "</table></body></html>"
print(html)
PYEOF

echo -e "  ${GREEN}✓${NC} Report: $REPORT_HTML"

# === Summary ===
echo ""
echo -e "${CYAN}=================================================${NC}"
echo -e "${CYAN}  Summary${NC}"
echo -e "${CYAN}=================================================${NC}"
echo "  Pages tested: $PASS/$TOTAL"
echo "  Missing: ${#MISSING[@]}"
for m in "${MISSING[@]}"; do echo "    - $m"; done
echo ""
echo "  JSON report: $RESULTS_JSON"
echo "  HTML report: $REPORT_HTML"
echo "  Visual diffs: $VISUAL_DIR"
