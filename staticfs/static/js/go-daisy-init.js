(function () {
  'use strict';
  if (window.gdOnce) return;              // idempotent under repeated script-tag loads after a swap
  var seen = Object.create(null);
  var behaviours = [];

  // Run fn exactly once per document lifetime, keyed by key.
  window.gdOnce = function (key, fn) {
    if (seen[key]) return false;
    seen[key] = true;
    try { fn(); } catch (e) { if (window.console) console.error('[go-daisy] init "' + key + '" failed', e); }
    return true;
  };

  // Run init once, then re-run rescan after every htmx swap so behaviours bound to
  // specific elements survive DOM replacement. rescan must be idempotent: track the
  // elements it has already handled (e.g. a data attribute or WeakSet), never
  // double-bind, and never throw on removed nodes.
  window.gdBehaviour = function (key, init, rescan) {
    var first = !seen[key];
    if (first) {
      seen[key] = true;
      try { init(); } catch (e) { if (window.console) console.error('[go-daisy] init "' + key + '" failed', e); }
    }
    // Register the rescan exactly once per key, even if several instances of a
    // component render their script on the same page.
    if (typeof rescan === 'function' && first) behaviours.push({ key: key, rescan: rescan });
    return first;
  };

  window.gdRescan = function () {
    for (var i = 0; i < behaviours.length; i++) {
      try { behaviours[i].rescan(); }
      catch (e) { if (window.console) console.error('[go-daisy] rescan "' + behaviours[i].key + '" failed', e); }
    }
  };

  // htmx 4 names its events with colon segments (htmx:after:settle); htmx 3 used
  // camelCase (htmx:afterSettle). Listen for both spellings so this runtime works
  // against either, whichever the host page bundles. A host fires only one of the
  // two spellings, so this cannot double-fire.
  //
  // Deliberately NOT listening for after:swap/afterSwap: htmx fires swap *and*
  // settle on every swap, and the first rescan wins for elements that mark
  // themselves as handled before binding. Binding on settle is the original
  // behaviour and the later of the two, so dimension-sensitive widgets (charts,
  // carousels) measure a fully settled DOM.
  document.addEventListener('htmx:after:settle', window.gdRescan);
  document.addEventListener('htmx:afterSettle', window.gdRescan);
})();
