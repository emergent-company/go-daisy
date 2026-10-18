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
| `--control-height-sm` | 2rem (32px) | small buttons/inputs, sidebar rows |
| `--control-height-md` | 2.5rem (40px) | default controls |
| `--table-row-height` | 2.75rem (44px) | table/datatable rows |

A theme overrides any of them in its own `:root`/`[data-theme]` block:

```css
[data-theme="dense"] {
  --layout-sidebar-width: 14rem;
  --table-row-height: 2.5rem;
}
```

## Pairing with variants

Tokens carry **measurements**; variants carry **structure**. Nexus parity work
uses both: e.g. table density is `DataTable{Density: ...}` (structural) while
row height comes from `--table-row-height` (themeable).
