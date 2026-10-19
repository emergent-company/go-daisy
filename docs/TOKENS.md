# Design tokens

go-daisy exposes a small token layer on top of daisyUI so component sizes and
tones come from the theme instead of per-component constants.

## Rules

1. **Colour, radius, border** come from daisyUI: `--color-*`, `--radius-box`,
   `--radius-field`, `--border`, plus the line ladder in
   `components/css/custom.css` (`--color-line*`).
2. **Sizes** come from the tokens below, each derived from daisyUI's
   `--size-field` primitive so a theme can rescale the chrome in one place.
3. A difference in **structure** (toolbar slot, density mode, icon position,
   layout variant) is a component prop/variant — not a new token.
4. Tokens are declared in `:root`, not `@theme`: Tailwind prunes `@theme`
   variables it cannot see used, and component sizes are often consumed from Go
   strings or inline `style` attributes that a consumer's scanner never reads.

## Layout / control tokens

Declared in `components/css/custom.css`.

| Token | Default (`--size-field: 0.25rem`) | Consumers |
|---|---|---|
| `--layout-sidebar-width` | 16rem (256px) | `layout.Sidebar`, admin shell |
| `--layout-topbar-height` | 4rem (64px) | topbar, admin shell |
| `--layout-content-pad` | 1.5rem (24px) | `#layout-content` |
| `--control-height-sm` | 2rem (32px) | sidebar menu rows, sidebar items, collapse titles |
| `--control-height-md` | 2.5rem (40px) | reserved (no consumer yet) |
| `--table-row-height` | 2.75rem (44px) | `Table{Density: DensityCompact}` rows |

A theme overrides any of them in its own `:root`/`[data-theme]` block:

```css
[data-theme="dense"] {
  --layout-sidebar-width: 14rem;
  --table-row-height: 2.5rem;
}
```

Tokens are only declared when a real consumer exists or is landing in the same
change. `--control-height-md` is deliberately listed as reserved: default
button/input heights still come from daisyUI, and forcing a height would break
textareas and outlined fields.

## Pairing with variants

Tokens carry **measurements**; variants carry **structure**. Nexus parity work
uses both: e.g. table density is `DataTable{Density: ...}` (structural) while
row height comes from `--table-row-height` (themeable).

## Chart text tokens

Declared in `components/css/tokens.css` (plain `:root`, not `@theme`). They
bind ApexCharts' hard-coded grey default (`#373d3f`) to the theme's content ink.

| Token | Default | Consumer |
|---|---|---|
| `--chart-label-color` | `color-mix(in oklab, var(--color-base-content) 60%, transparent)` | axis tick / annotation labels (≈ `text-muted`) |
| `--chart-legend-color` | `color-mix(in oklab, var(--color-base-content) 80%, transparent)` | legend series names (≈ `text-body`) |
| `--chart-title-color` | `var(--color-base-content)` | chart title / subtitle |

`components/css/custom.css` maps ApexCharts' SVG/legend classes onto these
tokens (`[data-chart] .apexcharts-xaxis-label` → `fill`, etc.). Because the
tokens are plain variables referencing `--color-base-content`, a theme switch
applies immediately without a chart re-render. Chart labels are
theme-controlled: per-chart label overrides would need a future `ChartProps`
field (no such field exists yet).

## Chart tooltip tokens

ApexCharts hard-codes a fixed light tooltip surface in every theme
(`rgba(255,255,255,.96)` body, `#eceff1` title, `#90a4ae` axis-tooltip border),
which leaves light ink on a white surface in dark themes. The following tokens
(in the same plain `:root` chart block) bind the tooltip surface and ink to the
daisyUI base ramp, so dark themes stay legible:

| Token | Default | Consumer |
|---|---|---|
| `--chart-tooltip-bg` | `color-mix(in oklab, var(--color-base-100) 96%, transparent)` | tooltip body surface |
| `--chart-tooltip-title-bg` | `var(--color-base-200)` | tooltip title / axis tooltip surface |
| `--chart-tooltip-border` | `var(--color-base-300)` | tooltip + axis-tooltip border and arrows |
| `--chart-tooltip-text` | `var(--color-base-content)` | tooltip + axis-tooltip ink |

The tooltip and axis-tooltip backgrounds, borders, and ink now follow
`base-100`/`base-200`/`base-300`/`base-content` respectively, so a theme switch
restyles the tooltip surface (and its arrow outline) without a chart re-render.
