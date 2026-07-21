package galleryruntime

// devOverlayScript is injected into the preview iframe when dev mode is active.
// It reads [data-component] / [data-props] attributes emitted by devmode.ComponentBoundary
// and renders a floating badge on hover.  Hold Alt to cycle through ancestors.
// This script is self-contained vanilla JS with no external dependencies.
const devOverlayScript = `
<script>
(function() {
  // Map from data-component name → gallery slug for click-to-navigate.
  // Injected at render time from ComponentSlugsJSON() — single source of truth.
  var COMPONENT_SLUGS = __COMPONENT_SLUGS_JSON__;

  // Depth-indexed colour palette for nested component indicators.
  var DEPTH_COLORS = [
    {bg:'#3b82f6',text:'#fff'},  // blue-500
    {bg:'#22c55e',text:'#fff'},  // green-500
    {bg:'#f59e0b',text:'#fff'},  // amber-500
    {bg:'#a855f7',text:'#fff'},  // purple-500
    {bg:'#ec4899',text:'#fff'},  // pink-500
  ];

  // --- Badge element ---------------------------------------------------------
  var badge = document.createElement('div');
  badge.id = '__dev-overlay-badge__';
  badge.style.cssText = [
    'position:fixed','z-index:2147483647','pointer-events:none',
    'font-family:ui-monospace,monospace','font-size:11px','line-height:1.4',
    'padding:3px 8px','border-radius:4px','white-space:nowrap',
    'box-shadow:0 2px 8px rgba(0,0,0,.35)','transition:opacity .1s',
    'opacity:0','max-width:400px','overflow:hidden','text-overflow:ellipsis',
  ].join(';');
  document.body.appendChild(badge);

  // --- Highlight outline element ---------------------------------------------
  var highlight = document.createElement('div');
  highlight.id = '__dev-overlay-highlight__';
  highlight.style.cssText = [
    'position:fixed','z-index:2147483646','pointer-events:none',
    'border-radius:3px','opacity:0','transition:opacity .1s',
  ].join(';');
  document.body.appendChild(highlight);

  // --- State ----------------------------------------------------------------
  var ancestors = [];   // [data-component] chain from target up to root
  var altDepth  = 0;    // which ancestor to highlight when Alt is held
  var altHeld   = false;

  // --- Helpers --------------------------------------------------------------

  function collectAncestors(el) {
    var chain = [];
    var cur = el;
    while (cur && cur !== document.body) {
      if (cur.hasAttribute && cur.hasAttribute('data-component')) {
        chain.push(cur);
      }
      cur = cur.parentElement;
    }
    return chain; // innermost first
  }

  function depthColor(idx) {
    return DEPTH_COLORS[idx % DEPTH_COLORS.length];
  }

  function summariseProps(raw) {
    var obj;
    try { obj = JSON.parse(raw); } catch(e) { return ''; }
    if (!obj || typeof obj !== 'object') return '';
    var parts = [];
    var skip = {class:1, style:1, id:1, Class:1, Style:1, ID:1};
    var keys = Object.keys(obj);
    for (var i = 0; i < keys.length && parts.length < 3; i++) {
      var k = keys[i];
      var v = obj[k];
      if (skip[k]) continue;
      if (v === false || v === '' || v === null || v === undefined) continue;
      var vs = String(v);
      if (vs.length > 20) vs = vs.slice(0, 17) + '…';
      parts.push(k + '=' + vs);
    }
    return parts.length ? ' ' + parts.join(' ') : '';
  }

  function showBadge(target, x, y) {
    if (!target) { hideBadge(); return; }
    var name  = target.getAttribute('data-component') || '?';
    var props = target.getAttribute('data-props') || 'null';
    var idx   = ancestors.indexOf(target);
    var col   = depthColor(idx < 0 ? 0 : idx);
    var label = name + summariseProps(props);

    badge.textContent = label;
    badge.style.background = col.bg;
    badge.style.color = col.text;
    badge.style.opacity = '1';

    // Position: prefer top-left of element, clamp to viewport.
    var rect = target.getBoundingClientRect();
    var bw = badge.offsetWidth || 200;
    var bh = badge.offsetHeight || 22;
    var left = Math.max(4, Math.min(rect.left, window.innerWidth - bw - 4));
    var top  = rect.top - bh - 4;
    if (top < 4) top = rect.bottom + 4;
    badge.style.left = left + 'px';
    badge.style.top  = top  + 'px';

    // Show highlight box
    highlight.style.left   = rect.left   + 'px';
    highlight.style.top    = rect.top    + 'px';
    highlight.style.width  = rect.width  + 'px';
    highlight.style.height = rect.height + 'px';
    highlight.style.border = '2px dashed ' + col.bg;
    highlight.style.background = 'transparent';
    highlight.style.opacity = '1';
  }

  function hideBadge() {
    badge.style.opacity = '0';
    highlight.style.opacity = '0';
  }

  function currentTarget() {
    var depth = altHeld ? altDepth : 0;
    return ancestors[depth] || null;
  }

  // --- Mouse events ---------------------------------------------------------

  document.addEventListener('mousemove', function(e) {
    var el = document.elementFromPoint(e.clientX, e.clientY);
    if (!el || el === badge || el === highlight) return;

    var newAnc = collectAncestors(el);
    if (newAnc.length === 0) { hideBadge(); ancestors = []; document.body.style.cursor = ''; return; }

    // Only reset altDepth when we move to a genuinely different component tree.
    if (!ancestors.length || ancestors[0] !== newAnc[0]) {
      altDepth = 0;
    }
    ancestors = newAnc;
    var target = currentTarget();
    showBadge(target, e.clientX, e.clientY);

    // Show pointer cursor when the hovered component has a gallery page.
    var name = target ? (target.getAttribute('data-component') || '') : '';
    document.body.style.cursor = COMPONENT_SLUGS[name] ? 'pointer' : '';
  }, {passive: true});

  document.addEventListener('mouseleave', function() {
    hideBadge();
    document.body.style.cursor = '';
  });

  // --- Click: navigate parent to the component's gallery page ---------------
  document.addEventListener('click', function(e) {
    var target = currentTarget();
    if (!target) return;
    var name = target.getAttribute('data-component') || '';
    var slug = COMPONENT_SLUGS[name];
    if (!slug) return;
    e.preventDefault();
    e.stopPropagation();
    try {
      window.parent.location.href = '/gallery/' + slug;
    } catch(ex) {}
  });

  // --- Alt key: cycle ancestors ---------------------------------------------

  document.addEventListener('keydown', function(e) {
    if (e.key !== 'Alt') return;
    altHeld = true;
    if (ancestors.length > 1) {
      altDepth = (altDepth + 1) % ancestors.length;
      showBadge(currentTarget(), 0, 0);
    }
    e.preventDefault();
  });

  document.addEventListener('keyup', function(e) {
    if (e.key !== 'Alt') return;
    altHeld = false;
    altDepth = 0;
    showBadge(currentTarget(), 0, 0);
  });

  // --- Post message: tree data to parent gallery window ---------------------
  // After load, serialise the component tree and post it so the parent page
  // can render the Component Tree panel in the right column.
  function buildTree(root) {
    var nodes = [];
    function walk(el, depth) {
      if (!el) return;
      if (el.hasAttribute && el.hasAttribute('data-component')) {
        // Skip ComponentBoundary wrappers — they are <div style="display:contents">
        // that exist only to annotate the real component beneath them.
        var style = el.getAttribute('style') || '';
        var isBoundary = el.tagName === 'DIV' && style.indexOf('display:contents') !== -1;
        if (!isBoundary) {
          var name  = el.getAttribute('data-component');
          var props = el.getAttribute('data-props') || 'null';
          // Strip package prefix (e.g. "ui/Progress" → "Progress") so names
          // match the sidebar and COMPONENT_SLUGS entries.
          var slash = name.indexOf('/');
          if (slash !== -1) name = name.substring(slash + 1);
          var clone = el.cloneNode(true);
          clone.removeAttribute('data-component');
          clone.removeAttribute('data-props');
          var html = clone.outerHTML;
          nodes.push({name: name, props: props, depth: depth, html: html});
          depth++;
        }
      }
      var children = el.children;
      for (var i = 0; i < children.length; i++) {
        walk(children[i], depth);
      }
    }
    walk(root, 0);
    return nodes;
  }

  function postTree() {
    var tree = buildTree(document.body);
    try {
      window.parent.postMessage({type: '__dev_component_tree__', tree: tree}, '*');
    } catch(e) {}
  }

  // Expose postTree so the parent gallery page can re-request the tree
  // after an HTMX navigation reassigns the iframe src.
  window.postTree = postTree;

  if (document.readyState === 'complete') {
    postTree();
  } else {
    window.addEventListener('load', postTree);
  }
})();
</script>
`
