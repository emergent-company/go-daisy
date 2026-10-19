# Nexus Port Plan — rebuild the Nexus example from go-daisy components

**Status:** P0 implemented on `feat/nexus-port-p0` (gate live, commit `565c292`); P1–P8 not started. Audit re-based on `origin/master` @ `8fcec8e`.
**Scope:** `cmd/nexus` (36 routes / 69 reference pages) + `components/**` + `pages/**`
**Rule:** library first. If a widget is missing or visually off, add/fix it in the library (with `*WithBoundary`, `gallery:token`, gallery seed, `docs/PATTERNS.md` entry) **before** Nexus consumes it.

---

## 1. Definition of done

A page is "ported" when **all** of these hold:

1. Its handler renders a `.templ` page composed from library components (`@ui.`, `@form.`, `@nav.`, `@layout.`, `@table.`, `@modal.`, `@pages.`).
2. No `nxHTMLComponent(...)` call, no `nxhtml/*.html` embed on that route.
3. No hand-rolled DaisyUI composite widget in the page (`<table>`, `<progress>`, breadcrumb block, stat card helper, chat thread, file grid). Plain layout `<div>`/grid wrappers are allowed.
4. Every composite it uses exists in the library with boundary wrapper + gallery token annotation.
5. Visual: pixel diff ≤ 2% on static pages, ≤ 5% on chart pages (charts masked), structural score ≥ 95% against `nxhtml/` reference.
6. `task build`, `task lint`, gallery seed render clean.

Target: **69/69 reference pages** rendering this way; `nxHTMLFiles` embed and the 19 passthrough handlers deleted.

---

## 2. Current state (measured)

Measured by the P0 gate (`task nexus:compare`, commit `565c292`):

```
compared 67, ok 7, diff 14, passthrough 46, error 0, skipped 2
```

| Metric | Now |
|---|---|
| Reference pages (`nxhtml/*.html`) | 69 — 67 compared, 2 skipped (`index` = dashboard alias, `starter` = no route) |
| Routes registered | 36 |
| **Raw passthrough pages** | **46 / 67 (69%)** — 19 `nxHTMLComponent` call sites; the UI-detail fallbacks expand to 46 pages |
| Structurally diffing pages | 14 |
| Structurally matching pages | 7 (settings, get-help, 5 showcases) — note: matching structure ≠ library-composed |
| Library adoption today | `layout.AdminShell` + `TopbarVariant("app")` shell, `pages.ChatLayout`, `table.TableCard*` on 3 lists, `ui.StatCardFeatured` ×4, `ui.Chart` ×6, `ui.Button` ×14, `ui.IconSpan` ×86. Everything else is hand-written DaisyUI markup |
| Library component inventory | 258 `*WithBoundary` wrappers, 183 `.templ` files |
| Duplicated hand-rolled patterns in Nexus | page heading+breadcrumb **14**, breadcrumb markup **15**, stat card helpers **3 divergent**, raw table scaffolds **3 files**, status-badge logic **3**, auth field blocks **13** |
| Assets | `/images` and `/js` served from absolute `/root/nexus-html/html` path; ApexCharts 4.3 / SimpleBar 6.2 / FilePond 4.32 loaded **from CDN** in `shell.templ:275-286` |

### Structural gaps that block a faithful port

1. **Runtime asset dependency on the source template.** `cmd/nexus/main.go:21-22` serves `/root/nexus-html/html/{images,js}` — the built binary is not self-contained, and vendor JS comes from cdnjs while `staticfs` already bundles htmx/alpine/morph.
2. **Shell is only half-library.** `layout.AdminShell` + `layout.TopbarVariant("app")` now host the frame, but the sidebar body (logo, promo, avatar dropdown), topbar controls, global search modal, and the right settings drawer are still hand-rolled in `cmd/nexus/internal/shell/shell.templ` (`layout.SidebarPromoCard`, `layout.SidebarAvatarDropdown`, `nav.SearchModal` exist but are unused).
3. **Elevation/style mismatch.** Library `ui.Card` emits `shadow-sm border border-base-200`; Nexus uses `card bg-base-100 shadow` (no border). `ui.Table` hardcodes `table-zebra`; Nexus tables are unstriped with `px-6` headers. `ui.StatCard*` are value-first with floating icon; Nexus KPI cards are label-first with an icon tile top-right and delta pill. None are drop-in today.
4. **No composite widgets for the app pages.** File manager, chat, gen-AI, file grids, activity feeds, pricing/feature cards, settings rows, help topics, prompt hero have no library component — these are why 19 routes fall back to raw HTML.
5. **Verification now sees the problem (P0).** The gate reports 46 passthroughs, 14 diffs, 7 structural matches and exits non-zero. Remaining limitation: it measures *structure*, not *component usage* — a hand-rolled page that copies the reference markup scores `ok`. `cmd/nexus/internal/parity/parity_test.go` is the composition guard and must be extended page-by-page as areas convert.

---

## 3. Workstream 0 — make the gate honest (DONE, commit `565c292`)

Delivered:

- `Pages()` covers all 69 reference files: 67 compared, 2 skipped with reasons; parameterised edit routes use `/…/1`.
- `nxHTMLComponent` emits `<!-- nx-raw:<page> -->`; any page containing it reports `passthrough` and can never score ok.
- Content extraction falls back to `<body>` when `#layout-content` is absent (auth pages, landing) — the previous false-OK path is gone.
- Lookahead sibling alignment replaces positional diffing (one inserted element no longer cascades).
- Gate entrypoint: flags `-base/-ref/-max-diff-pages/-json/-html/-filter`, summary counts, per-page lines, exit `1` on error/passthrough, `2` when diff pages exceed budget.
- Taskfile targets `nexus:build|serve|stop|compare`.
- Guard tests: `TestCompareCoversEveryReference`, `TestRawPassthroughIsDetectable`.

Baseline at P0:

```
compared 67, ok 7, diff 14, passthrough 46, error 0, skipped 2
```

Deferred to P2 (not gate-blocking): bundling images + vendor JS into `staticfs` and deleting the `/root/nexus-html` static mounts and CDN tags.

**Next:** P1. Pixel comparison stays a separate, optional target (`nexus:visual`) because it needs Playwright and pinned theme/masking; the structural gate is the required gate.

---

## 4. Component gap register

### 4.1 ADD (missing widgets) — 13

| # | Component | Package | Contents / API sketch | Consumers |
|---|---|---|---|---|
| 1 | `CardToolbar` | `ui` | Card header row: `Title`, `Subtitle`, `Actions templ.Component`, `Class`. Replaces `flex items-center justify-between px-5 pt-5` (~20 copies) | all cards, all lists |
| 2 | `ChartCard` | `ui` | `Title`, `Subtitle`, `Legend`, `Actions`, `Tabs`, `Chart templ.Component`, `Footer templ.Component`; body padding modes | ecom dash ×3, CRM ×3, 5× apex showcase |
| 3 | `ActivityFeed` + `ActivityItem` | `ui` | Item: icon/avatar, title, meta, timestamp, optional actions; density + connector line | order-details (Track Order), file-manager (File Activity), CRM social |
| 4 | `FileGrid` + `FileCard` + `StorageMeter` | `ui` | Folder/file tile (icon by type, name, meta, menu), grid/list switch, selection; storage bar + used/total | file-manager |
| 5 | `ConversationList` + `ConversationItem` + `ChatHeader` | `ui` | Avatar, name, last message, time, unread badge, online dot, active state, search + add slot | chat |
| 6 | `PromptHero` + `PromptSuggestionCard` | `ui` | Gradient greeting, subtitle, suggestion card (icon/title/cta) | gen-ai home |
| 7 | `PromptComposer` | `ui` | **Consolidates** `ChatComposer`+`form.PromptBar*`: textarea, model selector, attach/mic/image, token/usage hint, chips, send/stop, `Action`, `Method`, `MaxWidth` | gen-ai home/content/image, chat |
| 8 | `AIConfigCard` | `ui` | Brief textarea, keywords, tone/creativity(+dice)/language/voice selects, Generate/Stop, upgrade row | gen-ai content/image |
| 9 | `AIResultPanel` + `ImageResultGrid` | `ui` | History bubbles with hover actions (regenerate/copy/rate), "Older" divider, Clear; image grid with group headings, grid/list toggle, Load more | gen-ai content/image/library |
| 10 | `PricingCard`, `FeatureCard`, `AppPreviewShowcase` | `ui` | Plan card w/ feature checklist + CTA; icon feature tile; live demo control block with preview pane | landing |
| 11 | `SettingsSection` + `SettingsRow` + `ProfileHero` | `ui` | Section card, label/description/control row, hero with breadcrumb + title + action + watermark | settings, get-help |
| 12 | `HelpTopicCard` + `FaqGroup` | `ui` | Icon topic tile with arrow; accordion group with category header | get-help |
| 13 | `SettingsDrawer` | `layout` | Right offcanvas: theme cards, sidebar style, direction, reset — present on **all** reference pages | shell |

### 4.2 IMPROVE (exists, not drop-in) — 15

| Component | Problem today | Change |
|---|---|---|
| `nav.PageHeading` | `h1 text-2xl` + breadcrumbs above; Nexus is compact `h3 text-lg` with breadcrumb **right**, `hidden sm:inline`, actions baseline | add `Layout` enum (`Stacked`, `Split`), `Size` (`sm/md/lg`), `BreadcrumbPlacement`; keep default byte-identical |
| `nav.Breadcrumbs` | always full-size, `width:100%` | `Compact`, `HideOnMobile`, `Divider` options |
| `ui.Card` | `shadow-sm border` hardcoded | `Elevation` (`Flat/Raised/Outlined`), `BodyPadding`, `NoBody`; Nexus needs borderless `shadow` |
| `ui.StatCard` family | value-first, floating icon, no delta pill, sparkline is hand-rolled `[]int | add `Trend{Value,Direction}`, `Comparison`, `IconTile`, `Layout` (`ValueFirst`/`LabelFirst`); unify the 3 Nexus helper designs |
| `table.Table` / `TableWithProps` | zebra hardcoded, header padding fixed | `Zebra bool` (default off for app tables), `HeaderClass`, `Hover`, `Compact` |
| `table.DataGrid` | good base; no trailing actions, no kebab, batch actions POST-only, not in gallery/boundary | add `Actions templ.Component` (primary CTA + kebab), `Toolbar templ.Component`, `RowLink`, `EmptyRow`; register boundary + gallery + tokens |
| `table.TableCard*` | no header actions slot, footer fixed to 10/20/50/100 | `Actions`, `PageSizes []int` |
| `ui.Badge` | `badge-soft gap-1` vs Nexus `badge badge-soft px-1 gap-0.5` | size/density tokens; verify soft/outline parity |
| `ui.Accordion` | limited variants | arrow icon, `Multiple`, `Plus` styles (get-help 24 collapses, landing FAQ) |
| `ui.Tabs` | box/lift/bottom variants incomplete | add `Box`, `Lift`, `Bottom`, `Radio` (nexus uses `tabs-box`) |
| `form.FileUpload` | no preview thumbnail, progress, remove, avatar mode | `Preview`, `Avatar`, `Progress`, `MaxFiles`, drop target |
| `form.WizardStepper` | all panels rendered, no validation/HTMX | per-step validation hooks, `Action`, step persistence |
| `ui.Chat*` | `ChatBubble` actions hardcoded with no callbacks; `ChatWindow` no date dividers/auto-scroll; `ChatInput` single-line, no form action | richer item props, `DateDivider`, attachments/reactions, auto-scroll |
| `ui.Timeline` | only `Date/Label/Done` | icon/avatar, description, status, per-item action |
| `ui.Chart` | bare div + bespoke `data-chart`/`JSON.parse(dataset.options)` loop in shell; CDN ApexCharts | `Header`/`Footer` slots or pair with `ChartCard`; init from bundled script; `Height`, `Legend` props |
| `ui.EmptyState`, `ui.NotificationPanel`, `nav.SearchModal`, `layout.Sidebar*` | children/width/actions gaps; `SidebarMenuItemNexus` is Nexus-flavored naming leaking into the library; `SidebarPromoCard`/`SidebarAvatarDropdown` unused | wire real actions, add width/children props, promote Nexus sidebar helpers to a documented `App` variant |

### 4.3 UPDATE (cross-cutting, not per-component)

1. **Elevation + spacing tokens.** Pick one scale (`flat / raised / outlined`) and migrate library components and Nexus together. Root cause of most "visually different" reports.
2. **Theme pin.** Reference pages rely on JS-set `data-theme`; GD renders `data-theme=""`. Pin the theme used by the reference, port any custom theme vars into `components/css/custom.css` before measuring pixels.
3. **Assets.** Move `images/` + vendor JS (ApexCharts, SimpleBar, FilePond, Iconify) into `staticfs` (embedded, hashed); delete `e.Static("/images", "/root/nexus-html/...")` and CDN `<script>` tags. Nexus must build and run standalone.
4. **CSS consolidation.** Delete the inline `<style>` block in `shell.templ:12-20`; keep one rule set in `components/css/custom.css` (currently duplicated + conflicting z-index/transition).
5. **Dev-mode + gallery.** Every new component gets `*WithBoundary`, `gallery:token`/`gallery:hint`, seed entry, and a `docs/PATTERNS.md` line. No exceptions.

---

## 5. Page workstreams

Reference pages collapse into 9 templates, not 69 bespoke pages:

| WS | Pages | Covering template |
|---|---|---|
| W1 Shell | all | `layout.AppShell` + library sidebar/topbar/rightbar |
| W2 Dashboards | ecommerce, crm, index (dup) | `pages.DashboardPage` + `ChartCard`/`StatCard`/`ActivityFeed` |
| W3 Ecommerce CRUD | products ×3, orders ×2, customers ×3, sellers ×3, shops ×3, gen-ai library = **15** | `pages.CRUDListPage` (list/toolbar/table/pagination/delete dialog) + `pages.CRUDEditPage` (form grid/footer) + `OrderSummary` |
| W4 Apps | chat, file-manager, gen-ai ×4 = **6** | chat shell, file shell, AI two-column shell |
| W5 Site | auth ×4, settings, get-help, landing, starter = **8** | `pages.AuthSplitLayout`, `SettingsSection`, `ProfileHero`, pricing/feature/FAQ |
| W6 Showcase | 21 components + 12 forms + 5 charts = **38** | **generate** from gallery registry (see 5.1) |

### 5.1 Showcase pages must not be hand-ported

The 38 `ui-*` pages are demo galleries of the library itself. Hand-porting them duplicates the gallery and guarantees drift. Generate them: extend `galleryruntime` with a `RenderShowcasePage()` that consumes the same registry the gallery uses, and route `/ui/components/:name` etc. through it. One source of truth for demo markup; the gallery is then also a Nexus acceptance surface.

---

## 6. Phases

| Phase | Work | Gate |
|---|---|---|
| **P0** Gate ✅ | rewrite compare harness (69 pages, ported-vs-passthrough, exit codes), task targets | gate live: `compared 67, ok 7, diff 14, passthrough 46, skipped 2`, exit 1 |
| **P1** Foundation | elevation tokens, theme pin, `CardToolbar`, `PageHeading` split, `Breadcrumbs` compact, `StatCard` KPI, `ChartCard`, `table.Table` variants, `DataGrid` actions slot | 3 showcase pages rendered from library; component visual snapshots vs reference |
| **P2** Shell | `layout.AppShell` = sidebar + topbar + `SettingsDrawer` + `SearchModal` from library; delete inline CSS; bundle assets | shell pixel diff ≤ 2%; no hand-rolled topbar/rightbar |
| **P3** Dashboards | W2 (3 pages) | 3 pages ported, no raw helpers, pixel ≤ 5% (charts masked) |
| **P4** CRUD | W3 (15 pages) via 2 scaffolds | 15 pages ported; table/toolbar/pagination identical across pages |
| **P5** Apps | W4 (6 pages) | 6 pages ported; no `nxHTMLComponent` remains for apps |
| **P6** Site | W5 (8 pages) | 8 pages ported; auth pages have real `#layout-content` so they are actually compared |
| **P7** Showcase | W6 (38 pages) generated | all showcase routes served by registry renderer |
| **P8** Cleanup | delete `nxhtml/` embed + `nxHTMLComponent` + dead code (`CreateProductForm`, `EditProductForm`, `UIChartDetailPage`, `GenAI*` stubs), raw-HTML audit test, docs | 69/69 ported; audit test green; `go build ./...`, `task lint` |

Dependency chain: P0 → P1 → {P2, P3..P7}. P2 is required before P3's visual gate. P7 depends on P1 only.

Rough effort: P0 1–2d, P1 3–4d, P2 2d, P3 2d, P4 4d, P5 5d, P6 3d, P7 2d, P8 2d → **~24–26 dev-days**, 8 PR-sized units.

---

## 7. Guardrails (prevent regression)

Add `cmd/nexus/audit_test.go` (no network, parses `.templ` sources):

- fail if `nxHTMLComponent` appears anywhere in `cmd/nexus/internal/handler`;
- fail if a page file contains raw `<table`, `<progress`, `<dialog`, `class="breadcrumbs"`, `class="card` composite scaffolds — allowlist by phase so the list shrinks monotonically;
- fail if a `@ui.`/`@form.` symbol used by Nexus has no `*WithBoundary` in the library;
- report a counter (`raw composites remaining`) in CI output so progress is visible.

Plus process rules:

- every library change ships with `*WithBoundary` + `gallery:token` + seed + `docs/PATTERNS.md`;
- nexus PRs that need a component change must land the library change first (separate PR), never an inline workaround;
- writer lanes run in worktrees (`/root/emergent.memory-wt/<slug>`); `cmd/nexus` currently has untracked parallel work (`components/ui/agent-*.templ`) — do not sweep it into port commits.

---

## 8. Decisions needed before P1

1. **Parity target.** Recommended: semantic parity + 2% pixel (static) / 5% (charts masked). Pixel-perfect everywhere is not achievable with bundled charts and relative timestamps.
2. **Showcase strategy.** Recommended: generate 38 pages from gallery registry (vs hand-port). Big effort saver.
3. **Library placement for app shells.** Recommended: `layout` app variants for sidebar/topbar/drawer, `pages` scaffolds for CRUD/auth. Alternative: keep them inside `cmd/nexus` and accept duplication.
4. **Strictness.** Recommended: raw layout `<div>` allowed; raw DaisyUI composite widgets forbidden (enforced by audit test).
5. **Nexus as dogfood vs shipped templates.** Some scaffolds (`CRUDListPage`, `AuthSplitLayout`, `DashboardPage`) are reusable enough to ship in `pages/`; confirm they should be public API.

---

## 9. Risks

- **Component sprawl.** 258 wrappers already; adding 13 more must be paired with consolidation (merge `ChatComposer`+`PromptBar`, retire duplicate stat/sparkline variants) or the library becomes unnavigable.
- **Pixel gate churn.** Chart animations, CDN assets, and theme defaults will produce false failures unless pinned/masked in P0.
- **Nexus-specific naming leaking into the library** (`SidebarMenuItemNexus`) — formalize as a documented `App` variant or the leak spreads with every new widget.
- **Parallel sessions.** `components/ui/misc.templ`, `seed_additions.go`, `pages/*` are currently dirty from another lane; port work must branch from clean `main` in a worktree.
