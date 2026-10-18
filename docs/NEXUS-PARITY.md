# Nexus parity

`cmd/nexus` is the reference application. The rule is **bidirectional**: the
library absorbs Nexus's structure and measurements; Nexus then consumes the
library instead of hand-rolling chrome. Nexus-specific branding and image assets
stay in the app.

## Where a difference belongs

| Difference | Solution |
|---|---|
| Measurement (width, height, padding, row height) | token in `components/css/custom.css`, derived from `--size-field` (see `docs/TOKENS.md`) |
| Structure (slot, toolbar, density, icon position, shell variant) | prop / variant on the component |
| New visual object with its own behaviour | new component + boundary wrapper |

Brand-specific names are not exported. The first pass shipped `*Nexus` layout
names; the supported API uses semantic names (`AdminShell`,
`TopbarVariant("app", ...)`).

## Mapping

| Nexus markup | Library API |
|---|---|
| sidebar + content + right drawer shell, hidden toggle checkbox, backdrop | `layout.AdminShell(AdminShellProps{Sidebar, Topbar, Rightbar, Content})` |
| application topbar (hamburger, search pill, theme, settings, notifications, profile) | `layout.TopbarVariant("app", TopbarVariantOpts{...})` |
| search pill + mobile search icon | `nav.SearchTrigger(SearchTriggerProps{DialogID, Placeholder})` |
| orders / products tables inside cards | `table.TableWithProps` + `table.TableCard{Toolbar, Footer}` + `table.TableSelectAll/Row` |
| compact table rows | `Table{Density: table.DensityCompact}` (height from `--table-row-height`) |
| chat page (contact card, peer header, bubbles, composer) | `pages.ChatLayout(ChatLayoutProps{...})` + `pages.ChatContactItem` |
| marketing landing page | `pages.LandingTopbar/Hero/FeatureGrid/Testimonials/Pricing/FAQ/Footer` |
| dev-only component ids | `devmode.Attrs` + Nexus devmode middleware (`NEXUS_DEV=0` disables) |
| element-level feedback | `emergent-feedback.js` mounted from `layout.PageProps.BodyEnd` |

## Verification

1. `task build` (templ + Tailwind + codegen + Go build).
2. Gallery: the `table`, `page-chat`, `page-landing` and `admin-shell` stories
   render the same components Nexus uses.
3. Nexus: `cd cmd/nexus && go run .`, then check `/dashboards/ecommerce`,
   `/apps/chat`, `/landing`.
4. Visual comparison: `cmd/nexus/compare.sh` against `nexus-html`
   (see `cmd/nexus-compare/`).
5. `docs/TOKENS.md` must not claim a consumer that does not exist.

## Remaining gaps

- File manager and parts of the dashboards still contain raw DaisyUI markup;
  they convert as components are extracted, not before.
- The Nexus landing showcase strip has no dedicated section component.
- `nexus-compare` snapshots need re-baselining once Nexus renders through
  `AdminShell` (devmode adds `data-component` attributes to the DOM).
