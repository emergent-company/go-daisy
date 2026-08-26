# go-daisy Roadmap — P0 to P3

Goal: most modern SSR UI library. Fast, non-repetitive, rich micro-interactions.

---

## P0 — Foundation (unlocks everything else)

### P0.1 Alpine.js Integration

**What:** Ship 15KB Alpine.js as bundled JS. Use for all client-side micro-interactions — toggles, transitions, local state. Keep it behind CDN or embed so zero Go code impact for users who don't need it.

**Why P0:** Every UI micro-interaction (modal open/close, dropdown toggle, tab switch, accordion expand, theme toggle) currently relies on CSS-only DaisyUI hacks. No transitions, no focus traps, no keyboard nav. Alpine fixes all of this with `x-show`, `x-transition`, `x-trap`, `x-data` — declaratively, no custom JS.

**Benefits:**
- Smooth enter/leave transitions for modals, dropdowns, toasts (DaisyUI CSS can't animate `display:none`)
- Focus trapping for modals (accessibility — currently impossible with CSS-only)
- Keyboard navigation for dropdowns/menus (Escape to close, arrows to navigate)
- Persistent theming via Alpine `$persist` plugin (localStorage, survives page reloads)
- Toast queue with auto-dismiss (no custom JS needed)
- Instant tab switching (no server roundtrip for simple tab panels)
- Zero server-side state duplication — server renders initial state, Alpine manages subsequent mutations

**Examples:**

```html
<!-- Current: CSS-only modal, no transition, no focus trap -->
<dialog class="modal" open>
  <div class="modal-box">...</div>
</dialog>

<!-- With Alpine: smooth transition + focus trap + Escape to close -->
<div x-data="{ open: false }">
  <button @click="open = true">Open</button>
  <div x-show="open" x-trap="open" x-transition @keydown.escape="open = false"
       class="modal" :class="open && 'modal-open'">
    <div class="modal-box" @click.outside="open = false">...</div>
  </div>
</div>
```

```html
<!-- Current: Tabs require full server roundtrip -->
<div role="tablist" class="tabs tabs-boxed">
  <a href="/tab1" class="tab">Tab 1</a>
  <a href="/tab2" class="tab">Tab 2</a>
</div>

<!-- With Alpine: instant client-side switching -->
<div x-data="{ tab: 'one' }">
  <div role="tablist" class="tabs tabs-boxed">
    <button @click="tab = 'one'" :class="tab === 'one' ? 'tab tab-active' : 'tab'">Tab 1</button>
    <button @click="tab = 'two'" :class="tab === 'two' ? 'tab tab-active' : 'tab'">Tab 2</button>
  </div>
  <div x-show="tab === 'one'">Content 1</div>
  <div x-show="tab === 'two'">Content 2</div>
</div>
```

```html
<!-- Accordion with smooth collapse -->
<div x-data="{ open: false }">
  <button @click="open = !open">Section</button>
  <div x-show="open" x-collapse>Content with animated height</div>
</div>
```

**Implementation plan:**
1. Add `alpinejs` to `package.json`, ship minified in `staticfs/static/js/`
2. Add Alpine `$persist`, `$focus`, `$collapse` plugins
3. Create `components/alpine/` package with Go helpers:
   ```go
   alpine.Attrs()       // → {"x-data": "..."}
   alpine.Toggle()      // → {"x-show": "open", "x-transition": ""}
   alpine.Trap()        // → {"x-trap": "open"}
   alpine.Persist()     // → {"x-persist": ""}
   ```
4. Update `Modal`, `Dropdown`, `Tabs`, `Accordion`, `ThemeToggle` to accept optional Alpine attrs
5. Update `layout.Page` to include Alpine `<script>` when `WithAlpine(true)` option set

**Size cost:** 15KB gzipped (~45KB uncompressed). Optional — only loaded when page opts in.

**Breaking change:** None. Alpine is additive — existing CSS-only behavior preserved. Components accept `WithAlpine bool` prop, default `false`.

---

### P0.2 DOM Morphing (idiomorph)

**What:** Replace raw `innerHTML` swaps with DOM morphing. When HTMX swaps content, morph the existing DOM in-place instead of destroying and recreating it. Preserves input focus, scroll position, video/audio playback state, CSS transition state.

**Why P0:** Current `innerHTML` swaps destroy all transient browser state. User typing in an input → loses cursor position on any nearby HTMX swap. Playing a video → resets on page refresh. Scrolled half-way down → jumps to top. Morphing eliminates all of this. Feels like a SPA without being one.

**Benefits:**
- Input focus/cursor preserved across HTMX navigations
- `<video>` / `<audio>` continue playing through swaps
- Scroll position stable (no "jump to top" on partial updates)
- CSS animations aren't interrupted mid-transition
- `hx-preserve` elements (marked by id) survive any swap
- Works for both full-page navigation (`RenderAuto`) and partial swaps

**How it works under the hood:**
1. HTMX fetches new HTML from server
2. Before swap, idiomorph diffs current DOM against new HTML
3. Only changed nodes are added/removed/updated
4. Identical nodes (same id + tag) are left untouched

**Examples:**

Without morphing (current):
```
User types "Hello" in search box → clicks "Load more" → HTMX swaps list area
→ search box innerHTML replaced → "Hello" text lost, cursor gone
```

With morphing:
```
User types "Hello" in search box → clicks "Load more" → HTMX swaps list area
→ list area morphs (only new items appended) → search box untouched → "Hello" preserved, cursor stays
```

```html
<!-- Mark elements that must never be destroyed -->
<video hx-preserve="true" src="..." controls></video>
<audio hx-preserve="true" src="..." controls></audio>
<div hx-preserve="true" id="chat-input">...</div>
```

```go
// New render helper
render.MorphSwap(c, targetSelector)
// → sets HX-Reswap header + triggers morphdom/idiomorph on client
```

**Implementation plan:**
1. Add `idiomorph` (or `morphdom`) to bundled JS (~5KB gzipped)
2. HTMX already supports morph extensions — wire `hx-swap-ext` via JS config
3. Add `render.MorphSwap()` helper that sets appropriate response headers
4. Add `hx-preserve` to all interactive components (modal body, form inputs, video/audio)
5. Default all `RenderAuto` swaps to morph mode (configurable)

**Size cost:** 5KB gzipped for idiomorph.

**Breaking change:** Potentially. Morphing can cause issues if server HTML IDs don't match client IDs exactly. Default to opt-in per-page with `meta` tag, similar to Turbo's `turbo-refresh-method: morph`.

---

## P1 — Polish Layer

### P1.1 View Transitions API

**What:** Browser-native page transition animations (crossfade, slide) without any JS library. Declare via CSS `@view-transition` rules. Transitions work between any two page states — full navigation or partial swaps.

**Why P1:** Zero JS, zero library weight. Browsers handle the entire animation pipeline. One `<meta>` tag + a few CSS rules = SPA-grade page transitions. Already shipping in Chrome, Edge, Safari — Firefox behind flag.

**Benefits:**
- Page-to-page crossfade (default, no CSS needed beyond meta tag)
- Directional transitions (slide left on forward, right on back)
- Shared-element transitions (animate hero image from list → detail)
- Works with HTMX `hx-swap` — not just full page loads
- `prefers-reduced-motion` respected automatically

**Examples:**

```html
<!-- Enable view transitions for all same-origin navigations -->
<meta name="view-transition" content="same-origin" />
```

```css
/* Directional transitions */
html[data-navigation-direction="forward"]::view-transition-old(root) {
  animation: slide-to-left 0.3s ease-out;
}
html[data-navigation-direction="back"]::view-transition-old(root) {
  animation: slide-to-right 0.3s ease-out;
}
html[data-navigation-direction="replace"]::view-transition-old(root) {
  animation: fade-out 0.15s ease-out;
}

/* Shared element: card → detail */
.card-hero {
  view-transition-name: hero-image;
}
.detail-hero {
  view-transition-name: hero-image;
}
```

```go
// Go helper in render package
render.EnableViewTransitions(c) // sets meta tag in head
// Or as layout.Page option:
layout.Page(layout.PageProps{
    ViewTransitions: true,
})
```

```templ
// Component: adds meta tag to head
templ ViewTransitionMeta() {
  <meta name="view-transition" content="same-origin"/>
}
```

**Implementation plan:**
1. Add `render.EnableViewTransitions()` helper
2. Add CSS utilities in `assets/app.css` for common transition patterns (fade, slide-left, slide-right, scale-up)
3. Set `data-navigation-direction` attribute on `<html>` during HTMX swaps
4. Add `ViewTransitionMeta` component (or integrate into `layout.Page`)

**Size cost:** 0 bytes JS. CSS only. No library.

**Breaking change:** None. Feature-detected. Unsupported browsers silently ignore.

---

### P1.2 Real-Time Updates (SSE + WebSocket via `stream` Package)

**What:** Server-pushed partial page updates. No polling, no client-side polling loop. Server emits HTML fragments that morph into the live page. Two transports: Server-Sent Events (simple, one-way) and WebSockets (bidirectional).

**Why P1:** Live dashboards, notification panels, chat windows, log tailers — all need server push. Current go-daisy requires polling via HTMX `hx-trigger="every 2s"` which is wasteful and slow. A `stream` Go package makes real-time trivial.

**Benefits:**
- Dust live dashboards — stock tickers, server metrics, CI status
- Live notification dropdown — new alerts appear without page refresh
- Chat window — new messages stream in from other users
- Log tailer — real-time log lines appended
- Progress bars — streaming upload/download progress
- Model-agnostic — same Go API for SSE or WebSocket
- Reuses existing Templ components — no duplicate rendering logic

**Examples:**

```go
// Server: broadcast to all connected clients
stream.Broadcast(
    stream.Append("#messages", MessageComponent(msg)),
    stream.Replace("#counter", CounterComponent(n)),
    stream.Remove("#alert-42"),
)

// Server: send to specific client
stream.Send(c, stream.Update("#status", StatusComponent("done")))

// Server: refresh entire page via morph
stream.Send(c, stream.Refresh(method: "morph", scroll: "preserve"))
```

```html
<!-- Client: connect to SSE stream -->
<div hx-ext="sse" sse-connect="/events" sse-swap="message" hx-swap="beforeend">
  <!-- Messages stream in here -->
</div>
```

```html
<!-- Client: WebSocket for bidirectional chat -->
<div hx-ext="ws" ws-connect="/chat/room/42">
  <div id="messages">...</div>
  <form ws-send>
    <input name="message">
    <button>Send</button>
  </form>
</div>
```

```go
// Go helper: generate stream message from Templ component
streamMsg := stream.Append("#items", table.TableRow(item))
// → <turbo-stream action="append" target="items">
//     <template><tr id="item-3">...</tr></template>
//   </turbo-stream>
```

**Implementation plan:**
1. Add `htmx-sse` and `htmx-ws` extensions to bundled JS
2. Create `stream/` package:
   ```go
   stream.Append(target, comp)    // AppendBeforeEnd
   stream.Prepend(target, comp)   // AppendBeforeBegin
   stream.Replace(target, comp)   // Replace (outerHTML)
   stream.Update(target, comp)    // Update (innerHTML)
   stream.Remove(target)          // Remove
   stream.Before(target, comp)    // BeforeBegin
   stream.After(target, comp)     // AfterEnd
   stream.Refresh()               // Full page refresh w/ morph
   stream.Morph(target, comp)     // Replace with morph
   stream.Broadcast(msgs...)      // Fan-out to all connections
   ```
3. SSE transport: Echo handler that keeps connections open, writes stream messages
4. WebSocket transport: Echo websocket upgrade, bidirectional
5. Auto-reconnect on connection loss (HTMX handles this)
6. Gallery demo: live counter, live notification panel, live chat

**Size cost:** ~8KB gzipped for htmx-sse + htmx-ws extensions.

**Breaking change:** None. New package, additive.

---

### P1.3 Alpine Data Component Helpers

**What:** Go functions that generate Alpine `x-data` JSON strings for common interaction patterns. Type-safe, no hand-writing JSON in Templ files.

**Why P1:** Writing `x-data="{ open: false, tab: 'one' }"` as raw strings in Templ is error-prone (missing quotes, syntax errors). Go structs → JSON → safe injection. Also provides type safety for Alpine state shape.

**Benefits:**
- Type-safe Alpine state generation from Go
- Shared patterns across the component library
- Gallery can showcase all Alpine interaction patterns
- No raw JSON strings in `.templ` files

**Examples:**

```go
// In Go handler:
state := alpine.Tabs("overview", "settings", "billing")  // active: "overview"
state := alpine.Toggle(false)                              // open: false
state := alpine.Dropdown(false)                            // open: false
state := alpine.Theme(true)                                // dark: true
state := alpine.Accordion(nil)                             // open: null (all closed)
state := alpine.Counter(0)                                 // count: 0
state := alpine.Search("", []string{"foo", "bar"})        // search + items
```

```templ
// In Templ:
<div { alpine.XData(alpine.Toggle(false)) }>
  <button @click="open = !open">Toggle</button>
  <div x-show="open" x-transition>Content</div>
</div>
```

```go
// Package API:
package alpine

func XData(state State) templ.Attributes {
    json, _ := json.Marshal(state)
    return templ.Attributes{"x-data": string(json)}
}

type State struct {
    Open  bool   `json:"open,omitempty"`
    Tab   string `json:"tab,omitempty"`
    Count int    `json:"count,omitempty"`
    // ...
}
```

**Implementation plan:**
1. Create `components/alpine/data.go` with state generators
2. Create `components/alpine/alpine.go` with `XData()` helper
3. Add `x-cloak` support in `assets/app.css`
4. Document common patterns in gallery

**Size cost:** ~1KB Go source. No JS cost.

**Breaking change:** None.

---

## P2 — Quality of Life

### P2.1 Prefetch + Progress Bar (Turbo Drive Patterns)

**What:** Preload links on hover before the user clicks. Show a CSS progress bar during slow (500ms+) page loads. Both patterns from Turbo Drive — zero-cost polish that makes apps feel faster.

**Why P2:** Low effort, high perceived performance gain. Prefetch cuts perceived latency by 500–800ms. Progress bar gives visual feedback during slow loads (no more "is this loading?" confusion).

**Benefits:**
- Prefetch: link pages load instantly (already in cache when user clicks)
- Progress bar: thin colored bar at top of page during navigation
- `aria-busy` on `<html>` during navigation (screen reader feedback)
- `data-turbo-prefetch="false"` to exclude expensive pages
- Respects `navigator.connection.saveData` (no prefetch on slow connections)

**Examples:**

```html
<!-- Progress bar (auto-shown by JS during slow HTMX requests) -->
<div id="progress-bar" class="fixed top-0 left-0 h-1 bg-primary z-50 hidden"
     style="width: 0%"></div>
```

```html
<!-- Prefetch: link loads in background on hover -->
<a href="/expensive-page" data-prefetch="false">No prefetch</a>
<a href="/fast-page">Prefetched on hover (default)</a>
```

```go
// render helper: inject progress bar into layout
layout.Page(layout.PageProps{
    ProgressBar: true,      // adds progress bar + aria-busy
    Prefetch:    true,      // enables link prefetch on hover
})
```

```js
// Bundled JS (~2KB): progress bar logic
document.addEventListener("htmx:beforeRequest", () => {
  document.documentElement.setAttribute("aria-busy", "true")
  showProgressBar()
})
document.addEventListener("htmx:afterSettle", () => {
  document.documentElement.removeAttribute("aria-busy")
  hideProgressBar()
})
```

**Implementation plan:**
1. Write ~2KB vanilla JS for progress bar + prefetch (no library needed)
2. Add `ProgressBar` component in `components/ui/`
3. Integrate into `layout.Page` via `ProgressBar: true` prop
4. Prefetch: `mouseenter` listener on `[data-prefetch]` links (HTMX `hx-get` in background)

**Size cost:** ~2KB gzipped JS. Optional.

**Breaking change:** None.

---

### P2.2 Web Components for Heavy Client Widgets

**What:** Custom elements (`<go-daisy-calendar>`, `<go-daisy-richtext>`, `<go-daisy-gantt>`) for widgets that need rich client-side JS. Self-contained, framework-agnostic. Server renders initial HTML + data attributes, custom element bootstraps the JS.

**Why P2:** Components like calendar, rich text editor, gantt chart, sortable list need significant JS. Embedding all that JS in Templ files is messy and non-reusable. Web Components encapsulate behavior — usable outside Go/Templ projects.

**Benefits:**
- Calendar: fully keyboard-navigable, range selection, locale-aware
- Rich text: toolbar, formatting, paste handling, image upload
- Gantt: drag-to-reschedule, dependency lines, zoom
- Sortable: drag-drop reorder, handle-based vs area-based
- Survives morphing (custom element state persists)
- Exportable standalone — use in any HTML page, not just Go projects
- Shadow DOM encapsulation prevents CSS conflicts

**Examples:**

```html
<!-- Server renders container + config; custom element boots the JS -->
<go-daisy-calendar
  data-value="2025-01-15"
  data-min="2025-01-01"
  data-locale="en-US"
></go-daisy-calendar>

<go-daisy-richtext
  data-placeholder="Write something..."
  data-toolbar="bold italic link image"
></go-daisy-richtext>

<go-daisy-gantt
  data-tasks='[{"id":"1","name":"Design","start":"2025-01-01","end":"2025-01-05"}]'
></go-daisy-gantt>
```

```go
// Go helper: create calendar element from props
ui.CalendarWidget(ui.CalendarProps{
    Value:  time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
    Min:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
    Locale: "en-US",
})
// → <go-daisy-calendar data-value="2025-01-15" ...></go-daisy-calendar>
```

**Implementation plan:**
1. Create `webcomponents/` directory for custom element source
2. Build pipeline: write in vanilla JS (no framework), ship as ES module
3. First candidates: `go-daisy-calendar` (replace `form/calendar.templ`), `go-daisy-gantt` (wrap frappe-gantt)
4. Go wrappers in `components/form/calendar.go` that generate the custom element HTML
5. Bundle system: optionally include in `staticfs` or load via CDN

**Size cost:** Per widget: 10–40KB gzipped (calendar, richtext heavy). Optional — only loaded when used.

**Breaking change:** Components like `CalendarWrapper` become thin wrappers around `<go-daisy-calendar>`. Back compat via deprecation period.

---

### P2.3 Focus Trap + Keyboard Navigation

**What:** Modal focus trapping (Tab/Shift+Tab cycle within modal, Escape to close). Dropdown keyboard navigation (arrow keys, Enter to select, Escape to close). Both accessible (WCAG 2.1 AA) and expected by modern users.

**Why P2:** Current DaisyUI modal and dropdown are CSS-only. No keyboard handling. Power users expect keyboard navigation. Screen reader users need focus trapping.

**Benefits:**
- Tab/Shift+Tab trapped inside open modal
- Escape closes modal, dropdown, tooltip
- Arrow keys navigate dropdown menu items
- Enter/Space activates focused dropdown item
- Focus returns to trigger element on close
- Works with Alpine `x-trap` plugin or vanilla JS
- WCAG 2.1 AA compliant

**Examples:**

```html
<!-- Alpine x-trap: focus locked inside modal -->
<div x-data="{ open: false }" @keydown.escape="open = false">
  <button @click="open = true">Open</button>
  <div x-show="open" x-trap.noreturn="open" x-transition>
    <button @click="open = false">Close</button>
    <input placeholder="First focusable">
    <button>Last focusable</button>
  </div>
</div>
```

```html
<!-- Dropdown with keyboard nav -->
<div x-data="{ open: false, activeIndex: 0 }"
     @keydown.escape="open = false"
     @keydown.arrow-down.prevent="activeIndex = (activeIndex + 1) % items.length"
     @keydown.arrow-up.prevent="activeIndex = (activeIndex - 1 + items.length) % items.length"
     @keydown.enter.prevent="selectItem(activeIndex); open = false">
  <button @click="open = !open" :aria-expanded="open">Menu</button>
  <div x-show="open">
    <template x-for="(item, idx) in items">
      <button :class="{ 'bg-base-200': idx === activeIndex }" x-text="item"></button>
    </template>
  </div>
</div>
```

**Implementation plan:**
1. Add Alpine `$focus` plugin for focus trapping (already part of P0.1)
2. Create `components/alpine/trap.go` — `alpine.Trap(showExpr)` helper
3. Update `Modal` component: add `WithKeyboard bool` prop
4. Update `DropdownMenu` component: add `WithKeyboard bool` prop
5. Add `aria-expanded`, `aria-haspopup`, `role="menu"`, `role="menuitem"` to menus

**Size cost:** Alpine `$focus` plugin: ~1KB. No new library.

**Breaking change:** None.

---

## P3 — Compounding Value

### P3.1 Toast Notification Queue

**What:** System for queuing, displaying, and auto-dismissing multiple toast notifications. Swipe-to-dismiss on mobile. Position variants (top-right, bottom-center, etc.). Action buttons in toasts.

**Why P3:** Current `Toast` component is single-fire — can't stack multiple toasts. No queue management. No auto-dismiss. Every app needs this. Building it as a reusable system saves every user from writing their own.

**Benefits:**
- Stack multiple toasts (newest on top/bottom)
- Auto-dismiss after configurable timeout
- Swipe-to-dismiss on touch devices
- Action buttons (undo, retry, view)
- Position variants: top-right, top-center, bottom-right, bottom-center
- Types: info, success, warning, error
- Works with HTMX OOB swaps or Alpine queue
- Works with real-time streams (server pushes toast)

**Examples:**

```go
// Server-side: append toast to response
render.Toast(c, render.ToastProps{
    Type:    "success",
    Message: "File saved.",
    Duration: 3000,
    Action:  render.ToastAction{Label: "Undo", URL: "/undo/42"},
})
```

```html
<!-- Client: Alpine-managed toast container -->
<div x-data="toastQueue()" id="toast-container" class="fixed top-4 right-4 z-50 space-y-2">
  <template x-for="toast in toasts" :key="toast.id">
    <div x-show="toast.visible" x-transition:enter="..." x-transition:leave="..."
         class="alert" :class="'alert-' + toast.type" @click="dismiss(toast.id)">
      <span x-text="toast.message"></span>
    </div>
  </template>
</div>
```

**Implementation plan:**
1. Create `components/alpine/toast-queue.js` — Alpine data component for queue
2. Update `render.AppendToast` to inject HTMX OOB swap into queue
3. Add `ui.ToastContainer` component to instantiate the queue
4. Add swipe-to-dismiss (touch events)
5. Add position variants via props

**Size cost:** ~2KB JS. Uses existing `Toast` Templ component.

**Breaking change:** None. `render.AppendToast` behavior preserved.

---

### P3.2 Asset Versioning + Cache Busting

**What:** Auto-detect CSS/JS file changes and force browser reload. Append content hash to filenames. Track assets with `data-asset-track="reload"` (Turbo pattern). When deployed version changes, `<head>` merge detects mismatch and triggers full reload.

**Why P3:** Current `staticfs` has a hash file but it's not wired to auto-reload. Users on stale cache see broken UI. Fixing this prevents a class of support issues.

**Benefits:**
- Auto-reload when CSS/JS change (no "clear your cache" support tickets)
- Content hashes prevent stale cache after deploy
- Per-asset tracking (only reload what changed)
- Works with CDN edge caching
- `[data-asset-track="dynamic"]` removes assets when they disappear

**Examples:**

```html
<!-- Assets with content hash + track -->
<link rel="stylesheet" href="/static/css/app.abc123.css" data-asset-track="reload">
<script src="/static/js/htmx.def456.js" data-asset-track="reload"></script>

<!-- Dynamic asset: removed when absent from next page -->
<link rel="stylesheet" href="/page-specific.ghi789.css" data-asset-track="dynamic">
```

```go
// Go helper: generate asset tag with hash
staticfs.Stylesheet("/css/app.css")  // → <link href="/css/app.abc123.css" data-asset-track="reload">
staticfs.Script("/js/htmx.js")      // → <script src="/js/htmx.def456.js" data-asset-track="reload">
```

**Implementation plan:**
1. Generate content hash on `task build:ui` → write to `staticfs_hash.go`
2. Add `staticfs.Stylesheet()` / `staticfs.Script()` helpers
3. Wire `data-asset-track` in `layout.Page`
4. Add JS logic: on HTMX page swap, compare tracked assets → mismatch → full reload

**Size cost:** ~1KB JS for asset tracking. ~200 bytes per asset tag (hash in filename).

**Breaking change:** None. Hash filenames use existing `staticfs` path prefix.

---

### P3.3 Server-Side Streaming (HTML over the wire)

**What:** Stream HTML chunks to client as they're generated. First paint in 50ms, rest streams in. No waiting for full page render. Works with both initial page load (HTTP chunked transfer) and HTMX partial swaps.

**Why P3:** Complex pages with slow database queries or external API calls block rendering until everything completes. Streaming lets server send layout + fast queries immediately, then slow queries stream in as they finish. User sees content progressively.

**Benefits:**
- First paint 200–500ms faster for slow pages
- Progressive rendering — skeleton → partial → complete
- No client-side JS needed (browser native chunked transfer)
- Works with HTMX (swap chunks as they arrive)
- Server can cancel slow queries if client disconnects

**Examples:**

```go
// Handler: stream page in chunks
func DashboardHandler(c echo.Context) error {
    c.Response().Header().Set("Content-Type", "text/html")
    c.Response().WriteHeader(200)

    // Chunk 1: layout skeleton (immediate)
    layout.Shell().Render(c.Request().Context(), c.Response())
    c.Response().Flush()

    // Chunk 2: fast queries (50ms)
    stats := fetchQuickStats()
    ui.StatCard(stats).Render(c.Request().Context(), c.Response())
    c.Response().Flush()

    // Chunk 3: slow query (300ms)
    chart := fetchChartData()
    ui.Chart(chart).Render(c.Request().Context(), c.Response())
    c.Response().Flush()

    return nil
}
```

```html
<!-- Client: swap as HTML arrives (no hx-swap delay) -->
<div hx-get="/dashboard" hx-swap="outerHTML swap:0ms">
  <!-- Chunks swap in progressively -->
</div>
```

**Implementation plan:**
1. Create `stream/html` package — writer that wraps `http.ResponseWriter` with Flush()
2. Helper: `stream.Render(c, components...)` — renders each component to writer, flushes after each
3. HTMX swap options: set `swap:0ms` for instant swap on partial receipt
4. Gallery demo: dashboard with 3 tiers of query speed

**Size cost:** ~500 bytes Go code. No JS.

**Breaking change:** None. New package.

---

### P3.4 `hx-on` Integration + Inline Event Handlers

**What:** Leverage HTMX 2.x `hx-on` attribute for inline event handlers. Respond to HTMX lifecycle events (`htmx:afterRequest`, `htmx:responseError`, etc.) directly in HTML without writing global JS listeners.

**Why P3:** Currently, JS behaviors that react to HTMX events require `document.addEventListener(...)` in external JS files. `hx-on` moves these inline — co-located with the HTML they affect. More maintainable, easier to delete when component is removed.

**Benefits:**
- Event handlers co-located with HTML (no hidden listeners)
- Auto-cleaned when element is removed from DOM
- Supports any HTMX event: `htmx:afterRequest`, `htmx:responseError`, `htmx:validation:failed`
- Can call Alpine methods or vanilla JS

**Examples:**

```html
<!-- Reset form after successful HTMX POST -->
<form hx-post="/items" hx-on::after-request="if (event.detail.successful) this.reset()">
  <input name="name">
  <button>Submit</button>
</form>

<!-- Show toast on error -->
<form hx-post="/items"
      hx-on::response-error="Alpine.store('toasts').add('Failed!', 'error')">
  ...
</form>

<!-- Focus input after swap -->
<input hx-get="/search" hx-trigger="keyup changed delay:500ms"
       hx-on::after-settle="this.focus()">
```

```go
// Go helper: generate hx-on attrs
render.OnAfterRequest("if (event.detail.successful) this.reset()")
render.OnResponseError("showError(event.detail.xhr.response)")

// In Templ:
<div { render.OnAfterRequest("...") }>
```

**Implementation plan:**
1. Add `render.On*` helper functions for common HTMX events:
   ```go
   render.OnAfterRequest(js string) templ.Attributes
   render.OnAfterSettle(js string) templ.Attributes
   render.OnResponseError(js string) templ.Attributes
   render.OnBeforeRequest(js string) templ.Attributes
   ```
2. Add `render.OnErrorToast(msg string)` — shows toast on any HTMX error
3. Document patterns in gallery

**Size cost:** ~500 bytes Go code. Part of existing HTMX.

**Breaking change:** None.

---

### P3.5 Reduced Motion + Accessibility Audit

**What:** Full accessibility pass — add `aria-*` attributes, `role` attributes, `prefers-reduced-motion` media query support. Ensure all components work with screen readers, keyboard-only navigation, and high-contrast mode.

**Why P3:** Go-daisy components have minimal ARIA. The visual polish is good but a11y makes it production-grade. Companies evaluating UI libraries increasingly require VPAT/accessibility statements.

**Benefits:**
- Screen reader compatibility for all interactive components
- Keyboard-only navigation (Tab, Enter, Escape, arrows)
- `prefers-reduced-motion` disables animations
- High contrast mode support via DaisyUI themes
- Focus indicators visible and logical
- `aria-live` regions for dynamic content (toasts, live regions)
- `aria-busy` for loading states

**Examples:**

```html
<!-- Before: minimal ARIA -->
<button class="btn" onclick="...">Delete</button>

<!-- After: full ARIA -->
<button class="btn" aria-label="Delete item"
        @click="confirmDelete" :aria-busy="deleting"
        :disabled="deleting">
  <span aria-hidden="true">🗑</span>
  <span>Delete</span>
</button>
```

```css
/* Reduced motion: disable all animations */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

**Implementation plan:**
1. Audit all components for missing ARIA:
   - `role` on interactive elements (button, menu, menuitem, tab, tabpanel, dialog, alert, status)
   - `aria-label` on icon-only buttons
   - `aria-expanded` on toggles
   - `aria-haspopup` on menus
   - `aria-live` on dynamic content regions
2. Add `@media (prefers-reduced-motion)` to `assets/app.css`
3. Add `prefers-contrast` / `forced-colors` media query support
4. Test with screen reader (NVDA/VoiceOver) + keyboard-only

**Size cost:** CSS rules (~500 bytes). ARIA attrs (inline, negligible).

**Breaking change:** None.

---

### P3.6 Go-Native SSE/WebSocket Hub

**What:** Go-native pub/sub hub for real-time messages. Register connections, broadcast to channels, handle disconnect/reconnect. Built on Echo/Slog. Similar to Rails ActionCable but Go-native and simpler.

**Why P3:** `stream` package (P1.2) handles message format. This handles transport lifecycle. Without it, each user implements their own connection management. A shared hub means one import and done.

**Benefits:**
- One `hub := streamhub.New()` — all connections managed
- Channels: `hub.Broadcast("orders", msg)` — only subscribers receive
- Auto-heartbeat, auto-reconnect, auto-cleanup on disconnect
- Slog integration for observability
- Works with any Go framework (Echo, net/http, Chi)
- In-memory for single-server, pluggable for multi-server (Redis pub/sub)

**Examples:**

```go
// Bootstrap once
hub := streamhub.New(streamhub.Options{
    Logger:          slog.Default(),
    Heartbeat:       30 * time.Second,
    MaxMessageSize:  64 * 1024,
})

e := echo.New()
e.GET("/events", func(c echo.Context) error {
    return hub.ServeSSE(c.Response(), c.Request())
})

// Anywhere: broadcast
hub.Broadcast("orders", stream.Append("#orders", orderRow))
hub.Broadcast("notifications", stream.Update("#badge", badge))
```

```go
// Client connects with optional auth + channels
e.GET("/events", func(c echo.Context) error {
    userID := c.Get("user_id").(string)
    return hub.ServeSSE(c.Response(), c.Request(),
        streamhub.WithChannels("orders:"+userID, "notifications"),
    )
})
```

**Implementation plan:**
1. Create `streamhub/` package with SSE hub
2. WebSocket hub (or use `gorilla/websocket`)
3. Channel-based pub/sub with fan-out
4. Heartbeat / ping-pong for dead connection detection
5. Optional Redis pub/sub backend for multi-server
6. Gallery demo: live dashboard with multiple channels

**Size cost:** ~500 lines Go. No JS. Optional dependency.

**Breaking change:** None.

---

## Hotwire Deep Dive — What Turbo/Stimulus Teach That We Haven't Mapped Yet

### 1. Stimulus API Design — `data-controller` + `data-action` + `data-target`

**What:** Stimulus is Alpine's predecessor in the Hotwire stack. Alpine won mindshare (simpler, no class definitions), but Stimulus has design patterns worth copying:

| Stimulus concept | Alpine equivalent | What Stimulus does better |
|---|---|---|
| `data-controller="clipboard"` | `x-data="clipboard()"` | Explicit named controllers — easier to grep, debug, test |
| `data-action="click->clipboard#copy"` | `@click="copy()"` | `event->controller#method` is self-documenting; shows event type AND handler |
| `data-clipboard-target="source"` | `x-ref="source"` / `$refs.source` | Targets are typed, pluralizable (`targets`), with connect/disconnect lifecycle |
| `data-clipboard-text-param="hello"` | `x-bind:data-text="'hello'"` | Typed params (String/Number/Boolean) auto-parsed from data attrs |
| `data-clipboard-outlet=".parent"` | `$dispatch` + parent listener | Cross-controller typed communication |
| `connect()` / `disconnect()` | `x-init` / `$watch` | Standard lifecycle hooks; disconnect for cleanup is explicit |
| `static values = { delay: Number }` | Manual `JSON.parse(el.dataset.delay)` | Typed, defaulted values auto-bound to controller instance |

**Why it matters for go-daisy:** Alpine is the right default JS layer. But large apps with many interactive components benefit from Stimulus's explicit naming. go-daisy should provide **both** integration paths — helpers that generate `data-controller`/`data-action` attrs for Stimulus users AND `x-data`/`x-on` attrs for Alpine users.

**Examples — Stimulus vs Alpine in go-daisy:**

```html
<!-- Stimulus path: explicit controller + actions -->
<div data-controller="modal"
     data-action="keydown.escape->modal#close click@window->modal#outside">
  <button data-action="click->modal#open">Open</button>
  <dialog data-modal-target="dialog" class="modal">...</dialog>
</div>
```

```html
<!-- Alpine path: inline expressions (same component, different mode) -->
<div x-data="{ open: false }" @keydown.escape="open = false">
  <button @click="open = true">Open</button>
  <dialog x-show="open" x-transition class="modal" @click.outside="open = false">...</dialog>
</div>
```

```templ
// Go: component supports both
templ Modal(script ScriptMode, props ModalProps) {
  if script == ScriptStimulus {
    <div data-controller="modal" data-action="keydown.escape->modal#close">
      ...
    </div>
  } else {
    <div { alpine.XData(alpine.Toggle(false)) }>
      ...
    </div>
  }
}
```

**Implementation:**
1. Ship Stimulus as optional alternative to Alpine (~5KB gzipped vs Alpine's 15KB)
2. Create `components/stimulus/` package with Go helpers:
   ```go
   stimulus.Controller("modal")           // → {"data-controller": "modal"}
   stimulus.Action("click->modal#open")   // → {"data-action": "click->modal#open"}
   stimulus.Target("dialog")              // → {"data-modal-target": "dialog"}
   stimulus.Targets("items")              // → {"data-modal-targets": "items"}
   stimulus.Param("text", "hello")        // → {"data-modal-text-param": "hello"}
   ```
3. Pre-built Stimulus controllers for: modal, dropdown, tabs, accordion, clipboard, theme
4. Go `ScriptMode` type: `ScriptNone | ScriptAlpine | ScriptStimulus`
5. Every interactive component accepts `ScriptMode`, defaults to `ScriptAlpine`

**Size cost:** Stimulus core: 5KB gzipped. Pre-built controllers: ~3KB. Optional.

**Priority:** P2. Alpine first (simpler, bigger ecosystem). Stimulus as option for teams that prefer explicit controller pattern.

---

### 2. Permanent Elements — Surviving Any Navigation

**What:** Turbo's `data-turbo-permanent` marks elements that survive ALL page navigations. HTMX equivalent: `hx-preserve="true"`. Unlike morphing (which preserves identical elements by id), permanent elements are NEVER replaced even if the incoming HTML has a matching element.

**Use cases:**
- Audio player while user browses podcast episodes
- Active video call modal during site navigation
- Live chat sidebar persisting through page changes
- Form with unsaved input that must not be destroyed by background list refresh

**Why morphing isn't enough:** Morphing preserves elements by matching IDs. But if a page navigation replaces the entire `<body>` and the new HTML has a different structure, the element gets replaced. `hx-preserve` is stronger — the element is cloned out before swap and re-inserted after.

**go-daisy action:**

```templ
// Component: wrap element that must survive any swap
<div id="audio-player" hx-preserve="true">
  <audio src={ src } controls></audio>
</div>
```

```go
// Go helper: mark any element as permanent
render.Permanent(id string) templ.Attributes
// → {"id": "audio-player", "hx-preserve": "true"}

// layout.Page integration:
layout.PersistentAreas = []string{"#audio-player", "#live-chat"}
```

**Pattern for go-daisy layout:**
```
[Page Shell]
├── [Navbar]             ← normal element, replaced on navigation
├── [Sidebar]            ← normal element, hx-target swaps content area
├── [Main Content]       ← hx-target region, swapped on navigation  
├── [Audio Player]       ← hx-preserve="true", NEVER touched
└── [Live Chat FAB]      ← hx-preserve="true", NEVER touched
```

**Priority:** Dovetail into P0.2 (Morphing). Same implementation cycle.

---

### 3. Turbo Frame Advanced Patterns

Turbo Frames have nuances HTMX's `hx-target` doesn't replicate exactly:

#### 3a. Frame Scoped Navigation

In Turbo, any link/form INSIDE a `<turbo-frame>` automatically targets that frame. No `hx-target` needed on every element — it's inherited.

HTMX equivalent: `hx-target="this"` on each element, or attribute inheritance via parent.

**go-daisy enhancement:**

```templ
// Frame component: scoped navigation context
templ Frame(id string, src string, loading string, content templ.Component) {
  <div id={ id } 
       hx-target={ "#" + id }
       hx-select={ "#" + id }
       hx-swap="outerHTML"
       { templ KV("hx-trigger", "load", src != "") }
       { templ KV("hx-get", src, src != "") }>
    @content
  </div>
}
```

```go
// Usage: all links/forms inside auto-target this frame
@ui.Frame("messages", "/messages/latest", "eager", messagesList())
```

#### 3b. Frame Loading States

Turbo sets `[aria-busy="true"]` on the frame during fetch, `[complete]` after. Enables CSS loading states.

```css
turbo-frame[aria-busy="true"] { opacity: 0.6; pointer-events: none; }
turbo-frame[complete] { /* loaded */ }
```

HTMX equivalent: `htmx-request` class on triggering element.

**go-daisy enhancement:**

```templ
templ FrameShell(id string, content templ.Component) {
  <div id={ id }
       class="frame-container"
       hx-ext="class-tools"
       classes="{ add htmx-request:opacity-50 }">
    <div class="htmx-indicator absolute inset-0 flex items-center justify-center bg-base-100/50">
      @ui.Loader(ui.LoaderVariantSpinner)
    </div>
    @content
  </div>
}
```

#### 3c. Lazy Loading Below the Fold

Turbo: `loading="lazy"` — frame only loads when scrolled into view.  
HTMX: `hx-trigger="revealed"` or `hx-trigger="intersect once"`.

```templ
templ LazyFrame(id string, src string, placeholder templ.Component) {
  <div id={ id }
       hx-get={ src }
       hx-trigger="intersect once threshold:0.1"
       hx-swap="outerHTML">
    @placeholder  // skeleton initially, replaced on scroll-into-view
  </div>
}
```

#### 3d. Frame Cache Segmentation

Turbo's biggest architectural insight: **each frame is independently cacheable**. A composite page with 4 frames = 4 independent cache entries. Change one badge — only that frame's cache invalidates.

In go-daisy's context: HTMX partial responses are individually cacheable. But go-daisy doesn't currently provide cache-control helpers. This is a server-side concern.

**go-daisy enhancement:**

```go
// render helper: set cache headers for frame response
render.CacheFrame(c, 30*time.Second)
// → Cache-Control: private, max-age=30
// → Vary: HX-Request

render.CacheSharedFrame(c, 5*time.Minute)
// → Cache-Control: public, max-age=300
// → Vary: HX-Request
```

**Priority:** P2. Frame patterns are already possible with HTMX attrs. Go helpers make them discoverable and consistent.

---

### 4. Form Submission Convention — 303 + 422

Turbo enforces a strict convention for form submissions:
- **POST/PUT/PATCH/DELETE → 303 redirect.** Never render directly.
- **422 Unprocessable Entity → render validation errors.** No redirect.
- **GET form submission → render directly** (with `data-turbo-frame` target).

This avoids the browser's "Are you sure you want to resubmit?" dialog on refresh.

go-daisy's current `RedirectAfterMutation` already does this, but it's not woven into the form components.

**go-daisy enhancement:**

```templ
// FormModal: auto-handle 422 rendering
templ FormModal(props FormModalProps) {
  <form hx-post={ props.Action } hx-target={ "#" + props.ID }
        hx-swap="outerHTML"
        hx-on::after-request="if (event.detail.successful && event.detail.xhr.status === 200) {
            // 303 redirect handled by HTMX. 200 means server forgot to redirect.
            console.warn('Form submitted with 200, expected 303 redirect');
        }">
    @children
  </form>
}
```

```go
// Go handler pattern (convention helper)
func CreateItem(c echo.Context) error {
    item, err := service.Create(form)
    if err != nil {
        // Render form with errors on 422
        return render.RenderPartial(c, form.FormWithErrors(form, err))
    }
    // Redirect on success (303)
    return render.RedirectAfterMutation(c, "/items/"+item.ID)
}
```

**Priority:** P2. Already mostly followed. Formalizing in form components prevents drift.

---

### 5. `broadcasts_refreshes` — The Simplest Real-Time Pattern

Turbo's `broadcasts_refreshes` is a breakthrough: instead of crafting surgical `<turbo-stream>` messages (append this, replace that), just broadcast "refresh the page with morphing." The browser diffs the current page against the new HTML and only updates what changed.

```ruby
class Calendar < ApplicationRecord
  broadcasts_refreshes  # Any change → PageRefreshStream → all clients morph
end
```

This is dramatically simpler than maintaining targeted `append`/`replace`/`remove` for every mutation.

**go-daisy equivalent:**

```go
// Instead of:
stream.Broadcast(
    stream.Append("#items", ItemRow(newItem)),
    stream.Replace("#count", CounterComponent(n+1)),
    stream.Update("#badge", BadgeComponent(count)),
)

// Just:
streamhub.BroadcastRefresh("dashboard")
// → all clients on "dashboard" channel re-fetch page with morph
// → morphing diffs old vs new DOM, only updates what changed
```

```go
// Hub API
hub := streamhub.New()
hub.BroadcastRefresh("orders")
hub.BroadcastRefreshWithMethod("orders", streamhub.MorphMethod, streamhub.PreserveScroll)
```

**Benefits over targeted streams:**
- Server doesn't need to know what changed — sends full HTML, morphing handles it
- Impossible to get out of sync (no missing selectors, no stale target IDs)
- Same template used for initial render AND refresh — zero duplication
- Works with any page structure change (list reorder, new column, layout shift)
- One Go line for any mutation anywhere

**When NOT to use:** High-frequency updates (stock ticker, cursor position). Use targeted streams for sub-second updates. Refresh+morph for everything else.

**Priority:** Amend P1.2 (Real-Time Updates). `broadcast_refresh` is the default; targeted streams for performance-critical cases.

---

### 6. Stimulus Lifecycle as Alpine Conventions

Stimulus's `connect()`/`disconnect()`/`[name]TargetConnected()` lifecycle is more structured than Alpine's `x-init`. We can adopt these as Alpine conventions in go-daisy.

**Pattern:**

```js
// go-daisy Alpine convention: init/destroy hooks
Alpine.data("goDaisyModal", () => ({
  open: false,
  init() {
    this.$watch("open", (value) => {
      if (value) {
        document.body.style.overflow = "hidden"        // lock scroll
        this.$nextTick(() => this.$refs.close?.focus()) // focus close btn
      } else {
        document.body.style.overflow = ""               // unlock scroll
      }
    })
  },
  destroy() {
    document.body.style.overflow = ""  // cleanup on element removal
  },
}))
```

```go
// Go: pre-built Alpine data components
alpine.ModalData(false)     // init/destroy with scroll lock + focus trap
alpine.DropdownData()       // init/destroy with keyboard nav
alpine.TabsData("tab1")     // init/destroy with arrow key navigation
alpine.ClipboardData()      // init/destroy with clipboard API
```

```templ
// Usage
<div { alpine.XData(alpine.ModalData(false)) }>
  <div x-show="open" x-trap="open" x-transition>...</div>
</div>
```

**Priority:** Amend P1.3 (Alpine Data Helpers). Pre-built data components with proper lifecycle.

---

### 7. Turbo's Asset Tracking Revisited

Beyond what P3.2 covers, Turbo has two asset tracking modes:

| Attribute | Behavior |
|---|---|
| `data-turbo-track="reload"` | If the asset's URL (hash) differs between old and new page → full page reload |
| `data-turbo-track="dynamic"` | Remove the asset element if it's absent from the new page's `<head>` |

`dynamic` mode solves a specific problem: page-specific CSS/JS that should NOT persist to other pages. Without it, navigating from `/dashboard` (with `dashboard.css`) to `/settings` leaves `dashboard.css` loaded → style conflicts.

**go-daisy equivalent:**

```templ
// layout.Page: track assets
templ Page(props PageProps) {
  <html>
    <head>
      <link rel="stylesheet" href="/static/css/app.abc123.css" data-asset-track="reload">
      if props.PageAssets != nil {
        for _, asset := range props.PageAssets {
          <link rel="stylesheet" href={ asset } data-asset-track="dynamic">
        }
      }
    </head>
    ...
  </html>
}
```

```go
// Usage: page-specific assets auto-cleaned on navigation
layout.Page(layout.PageProps{
    PageAssets: []string{"/static/css/dashboard.def456.css"},
})
```

**Priority:** Amend P3.2 (Asset Versioning).

---

### 8. Turbo's `turbo-visit-control` Meta Tag

Turbo lets any page opt out of Turbo Drive navigation with:
```html
<meta name="turbo-visit-control" content="reload">
```

HTMX equivalent: `hx-boost="false"` on parent element, or `HX-Refresh: true` response header.

**go-daisy helper:**

```go
// Force full page load for specific pages (e.g. login, legacy pages)
render.ForceReload() // → sets <meta name="turbo-visit-control" content="reload">
```

Use case: third-party JS that doesn't survive HTMX swaps (Google Maps, Stripe Elements, legacy jQuery plugins).

**Priority:** P3. Low effort, useful escape hatch.

---

### 9. CSRF Token Injection

Turbo reads CSRF token from `<meta name="csrf-token">` and injects it into all form requests as `X-CSRF-TOKEN` header. HTMX does NOT do this automatically.

**go-daisy helper:**

```templ
templ CSRFMeta(token string) {
  <meta name="csrf-token" content={ token }/>
}
```

```go
// Middleware: set CSRF token in context
func CSRFToken(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Set("csrf_token", generateCSRFToken(c))
        return next(c)
    }
}
```

```js
// Bundled JS (~500 bytes): auto-inject CSRF token into HTMX requests
document.addEventListener("htmx:configRequest", (event) => {
  const token = document.querySelector('meta[name="csrf-token"]')?.content
  if (token) event.detail.headers["X-CSRF-TOKEN"] = token
})
```

**Priority:** P3. Required for production apps. Trivial to implement.

---

### 10. Hotwire Summary — What We're Adopting vs. What We're Not

| Turbo/Stimulus concept | go-daisy action | Priority |
|---|---|---|
| DOM Morphing + idiomorph | P0.2 — core feature | P0 |
| Stimulus controllers | Optional alternative to Alpine, `data-controller`/`data-action` helpers | P2 |
| `data-turbo-permanent` | `hx-preserve` + `render.Permanent()` helper | P0.2 |
| Turbo Frames (scoped nav) | `Frame` component + inherited `hx-target` | P2 |
| Frame lazy loading | `LazyFrame` component + `intersect` trigger | P2 |
| Frame cache segmentation | `render.CacheFrame()` / `render.CacheSharedFrame()` | P2 |
| Form 303/422 convention | Enforce in `FormModal`, document pattern | P2 |
| `broadcasts_refreshes` | `streamhub.BroadcastRefresh()` — simplest real-time primitive | P1.2 |
| Asset tracking (reload + dynamic) | `data-asset-track` + `PageAssets` prop | P3.2 |
| `turbo-visit-control` | `render.ForceReload()` helper | P3 |
| CSRF token injection | `CSRFMeta` component + auto-inject JS | P3 |
| Stimulus lifecycle as Alpine convention | Pre-built Alpine data components with `init()`/`destroy()` | P1.3 |
| `data-action` syntax | None — Alpine `@event="handler()"` covers it | — |
| Turbo Native (iOS/Android) | Not applicable to go-daisy | — |

---

## P2.4 Chart Component Enrichment

**Source:** nexus-html dashboard template (15+ ApexCharts configurations).

**Current state:** `chart.templ` has 10 props — basic line/bar/pie/donut/radial with no stacking, no gradients, no annotations, no synced groups, no forecasts, no data labels. ~208 lines.

**Goal:** match nexus-html's production chart patterns with type-safe Go props. Every ApexCharts feature nexus-html uses should have a corresponding struct field.

### New Props (12 additions to ChartProps)

| Prop | Type | Default | nexus-html source |
|---|---|---|---|
| `Stacked` | `bool` | false | apex-bar.js stacked 100% |
| `StackType` | `string` | "normal" | "100%" for percentage stack |
| `FillType` | `string` | "solid" | "gradient", "pattern" from apex-pie.js |
| `FillOpacity` | `float64` | 1.0 | apex-pie.js gradient donut |
| `StrokeCurve` | `string` | "smooth" | apex-line.js stepline |
| `StrokeWidth` | `float64` | 2.0 | apex-line.js thin lines |
| `ShowDataLabels` | `bool` | false | apex-line.js label line chart |
| `ForecastCount` | `int` | 0 | ecommerce.js forecastDataPoints |
| `GroupID` | `string` | "" | apex-line.js synced group |
| `GoalValues` | `[]float64` | nil | apex-bar.js goal markers |
| `LegendPosition` | `string` | "bottom" | apex-pie.js right legend |
| `Annotations` | `[]ChartAnnotation` | nil | apex-line.js annotations |
| `Monochrome` | `bool` | false | apex-pie.js monochrome pie |
| `Horizontal` | `bool` | false | currently auto-set for ChartBar only |
| `Dumbbell` | `bool` | false | apex-column.js rangeBar dunbbell |

### New Types

```go
type ChartAnnotation struct {
    Type  string // "yaxis", "xaxis", "point"
    Value float64
    Label string
    Color string
}
```

### New ChartType Constants

```go
ChartRangeBar ChartType = "rangeBar"
ChartHeatmap  ChartType = "heatmap"
```

### Examples

Stacked percentage bar:
```go
ui.Chart(ui.ChartProps{
    ID: "stacked-bar", Type: ui.ChartBar, Stacked: true, StackType: "100%",
    Series: []ui.ChartSeries{
        {Name: "Labor", Data: []float64{30, 40, 25, 50}},
        {Name: "Materials", Data: []float64{20, 30, 45, 30}},
    },
    Categories: []string{"Q1", "Q2", "Q3", "Q4"},
})
```

Gradient donut:
```go
ui.Chart(ui.ChartProps{
    ID: "gradient-donut", Type: ui.ChartDonut, FillType: "gradient",
    FillOpacity: 0.9, LegendPosition: "right",
    Series: []ui.ChartSeries{{Name: "Sales", Data: []float64{44, 55, 41, 17, 15}}},
})
```

Forecast line chart:
```go
ui.Chart(ui.ChartProps{
    ID: "forecast", Type: ui.ChartLine, ForecastCount: 3,
    StrokeCurve: "stepline", ShowDataLabels: true,
    Series: []ui.ChartSeries{{Name: "Revenue", Data: []float64{20, 35, 50, 45, 60, 55, 70}}},
})
```

Synced chart group:
```go
ui.Chart(ui.ChartProps{ID: "orders", Type: ui.ChartLine, GroupID: "dashboard-sync", ...})
ui.Chart(ui.ChartProps{ID: "revenue", Type: ui.ChartLine, GroupID: "dashboard-sync", ...})
```

### Size cost

Zero JS changes. `chart.templ` grows ~100 lines (from 208 to ~300). `chartOptionsJSON` gains ~50 lines of conditional serialization. Backward compatible — all new props default to zero/nil values.

### Priority: P2 — ships alongside existing chart while expanding capability.

---

## P2.5 Fieldset + Layout Customizer

**Source:** nexus-html `ui-forms-fieldset.html` + `app.js` LayoutCustomizer class.

### Fieldset Component

Simple standalone `<fieldset>` + `<legend>` wrapper. Currently go-daisy embeds fieldset patterns inline in form examples but has no reusable component.

```templ
templ Fieldset(legend string) {
    <fieldset class="fieldset">
        if legend != "" {
            <legend class="fieldset-legend">{ legend }</legend>
        }
        { children... }
    </fieldset>
}
```

~15 lines of Templ. No dependencies. P3 priority.

### LayoutCustomizer

Unified theme/direction settings drawer combining:
- Theme picker (light/dark/contrast/material/dim/system)
- Sidebar theme (independent from main theme)  
- RTL/LTR direction toggle
- Fullscreen toggle
- Reset to defaults

Current go-daisy has individual `ThemeToggle`, `ThemeSwitcher`, `ThemeController` but no unified settings panel. This is a composite component combining existing primitives.

~80 lines of Templ + ~60 lines JS. P3 priority.

---

## Appendix: What We're NOT Doing

| Anti-pattern | Reason |
|---|---|
| Rewrite Templ components as Web Components | Server-rendered Templ is the right layer. Web Components only for JS-heavy widgets. |
| Adopt Turbo.js wholesale | HTMX already serves the same role. Borrow patterns, not the library. |
| Build a React/Vue/Svelte competitor | go-daisy is HTML-over-the-wire. That's the architecture. |
| Add a JS build pipeline (esbuild/webpack) | HTMX + Alpine + idiomorph are vanilla JS. No build step needed. |
| Abstract HTMX attributes behind Go wrappers | `hx-get`, `hx-target`, `hx-swap` are HTML attributes — write them in Templ. Wrap only complex patterns (streams, morph, SSE). |
| Recreate shadcn/ui's `cn()` class merger | DaisyUI + Tailwind v4 handle class conflicts. No utility needed. |
| Server-side JS rendering (React SSR, etc.) | go-daisy is Go-native. Templ IS the SSR. |
| npm package for go-daisy components | go-daisy is a Go module. Web Components can be published separately to npm if needed. |

---

## Appendix: Bundle Size Budget

| Component | Size (gzipped) | Required |
|---|---|---|
| HTMX | 14KB | Always |
| idiomorph | 5KB | P0 (opt-in per page) |
| Alpine.js core | 15KB | P0 (opt-in per page) |
| Alpine persist plugin | 1KB | P0 |
| Alpine focus plugin | 1KB | P2 |
| Alpine collapse plugin | 1KB | P0 |
| Progress bar + prefetch JS | 2KB | P2 (opt-in) |
| Toast queue JS | 2KB | P3 (opt-in) |
| Asset tracking JS | 1KB | P3 (opt-in) |
| **Total required** | **14KB** | HTMX only |
| **Total with all features** | **~42KB** | HTMX + Alpine + morph + extras |
| **Target ceiling** | **50KB** | Never exceed this |

---

## Appendix: Delivery Cadence

| Phase | Items | Target |
|---|---|---|
| P0 | Alpine + Morphing | Single release cycle |
| P1 | View Transitions + SSE/WS stream + Alpine helpers | After P0 stabilizes |
| P2 | Prefetch/Progress + Web Components + Focus trap | Incremental |
| P3 | Toast queue + Asset versioning + Streaming + hx-on + A11y + Hub | Ongoing |
