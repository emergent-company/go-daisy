/*
 * go-daisy command palette runtime.
 *
 * The palette is a CSS-only modal: an `input.modal-toggle[data-command-palette]`
 * followed by a sibling `.modal`. CSS can show and hide that dialog, but it
 * cannot move focus into it, keep Tab inside it, or give it back on close.
 * This file adds that behaviour and owns the local/remote item filter.
 *
 * Contract:
 *   input.modal-toggle[data-command-palette]   toggle; may be any palette id
 *   .modal (sibling)                            the dialog surface
 *   .modal input[type="text"]                   receives initial focus
 *   [data-command-palette-keep-open]            item opt-out from click-to-close
 *
 * Behaviour:
 *   - Cmd/Ctrl+K toggles the first palette on the page.
 *   - Opening traps Tab/Shift+Tab inside the dialog (gdFocus.trap), focuses the
 *     input, and marks the dialog aria-modal via the server-rendered markup.
 *   - Escape, the Esc/backdrop labels, Cmd/Ctrl+K and activating a result all
 *     close the palette and restore focus to the element that opened it.
 *   - filterPaletteItems(input) filters local items and lazily fetches
 *     data-items-url results on first keystroke.
 *
 * Listeners are document-level and the initialiser is idempotent, so an HTMX
 * swap that re-inserts the component does not double-bind.
 */
(function () {
  'use strict';

  if (window.__gdCommandPaletteInit) return;
  window.__gdCommandPaletteInit = true;

  // toggle checkbox -> release() from gdFocus.trap
  var releases = new WeakMap();

  function toggleFor(modal) {
    var cb = modal ? modal.previousElementSibling : null;
    if (!cb || !cb.matches || !cb.matches('input.modal-toggle[data-command-palette]')) return null;
    return cb;
  }

  function modalFor(cb) {
    var modal = cb.nextElementSibling;
    return modal && modal.classList && modal.classList.contains('modal') ? modal : null;
  }

  function inputFor(modal) {
    return modal ? modal.querySelector('input[type="text"]') : null;
  }

  function activate(cb) {
    if (releases.has(cb)) return;
    var modal = modalFor(cb);
    if (!modal) return;
    var input = inputFor(modal);
    if (window.gdFocus && window.gdFocus.trap) {
      releases.set(cb, window.gdFocus.trap(modal, {
        initial: input,
        onEscape: function () { close(cb); },
      }));
    } else if (input) {
      // gdFocus ships with the component, but never leave the palette
      // un-focused if the a11y runtime failed to load.
      try { input.focus(); } catch (e) {}
    }
  }

  function deactivate(cb) {
    var release = releases.get(cb);
    if (!release) return;
    releases.delete(cb);
    release();
  }

  // A checkbox toggled from script does not emit `change`. Dispatch one so the
  // change listener below handles every open/close route (button, Cmd/Ctrl+K,
  // Escape, backdrop, result click). Code that dismisses by assigning
  // `toggle.checked = false` directly must call gdCommandPalette.close(id) so
  // the trap is released and focus comes back.
  function sync(cb) {
    cb.dispatchEvent(new window.Event('change', { bubbles: true }));
  }

  function open(cb) {
    if (!cb) return;
    if (!cb.checked) {
      cb.checked = true;
      sync(cb);
      return;
    }
    // Already checked (initial Open, or app code set .checked directly) but
    // possibly never trapped: make sure focus and the trap are installed.
    activate(cb);
  }

  function close(cb) {
    if (!cb) return;
    if (cb.checked) {
      cb.checked = false;
      sync(cb);
      return;
    }
    // Apps can dismiss with `toggle.checked = false` and no change event, then
    // call gdCommandPalette.close(id): release the trap and restore focus even
    // though the box is already off.
    deactivate(cb);
  }

  function initOpen() {
    var toggles = document.querySelectorAll('input.modal-toggle[data-command-palette]');
    for (var i = 0; i < toggles.length; i++) {
      if (toggles[i].checked) activate(toggles[i]);
    }
  }

  document.addEventListener('change', function (e) {
    var cb = e.target;
    if (!cb || !cb.matches || !cb.matches('input.modal-toggle[data-command-palette]')) return;
    if (cb.checked) activate(cb);
    else deactivate(cb);
  });

  document.addEventListener('keydown', function (e) {
    if (!e.key) return;
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
      var cb = document.querySelector('input.modal-toggle[data-command-palette]');
      if (!cb) return;
      e.preventDefault();
      if (cb.checked) close(cb);
      else open(cb);
    } else if (e.key === 'Escape') {
      // The trap handles Escape while it is installed (and stops propagation).
      // This keeps palettes opened by setting .checked directly, before the
      // trap existed, dismissable like the old inline runtime was.
      var toggles = document.querySelectorAll('input.modal-toggle[data-command-palette]');
      for (var i = 0; i < toggles.length; i++) {
        if (toggles[i].checked) close(toggles[i]);
      }
    }
  });

  // Activating a result should dismiss the palette. Anything that stays open
  // after its action (e.g. a "copy" row) opts out with
  // [data-command-palette-keep-open].
  document.addEventListener('click', function (e) {
    var t = e.target;
    if (!t || !t.closest) return;
    var item = t.closest('a[href]');
    if (!item || !item.closest('.modal')) return;
    if (item.closest('[data-command-palette-keep-open]')) return;
    var cb = toggleFor(item.closest('.modal'));
    if (cb) close(cb);
  });

  // htmx 4 emits colon-separated names (htmx:after:settle); htmx 3 used
  // camelCase (htmx:afterSettle). Listen for both, as the other runtimes do.
  document.addEventListener('htmx:after:settle', initOpen);
  document.addEventListener('htmx:afterSettle', initOpen);

  window.gdCommandPalette = {
    open: function (id) {
      open(document.getElementById(id + '-toggle'));
    },
    close: function (id) {
      close(document.getElementById(id + '-toggle'));
    },
  };

  window.filterPaletteItems = function (input) {
    var q = input.value.toLowerCase().trim();
    var modal = input.closest('.modal');
    var results = modal ? modal.querySelector('[id$="-results"]') : document.getElementById(input.id.replace('-input', '-results'));
    if (!results) return;

    if (results.dataset.itemsUrl && !window['_paletteFetched' + results.id]) {
      window['_paletteFetched' + results.id] = true;
      var xhr = new XMLHttpRequest();
      xhr.open('GET', results.dataset.itemsUrl);
      xhr.onload = function () {
        if (xhr.status === 200) {
          try {
            var data = JSON.parse(xhr.responseText);
            window['_paletteRemote' + results.id] = data;
          } catch (e) {}
        }
        filterResults(q, results);
      };
      xhr.send();
      return;
    }
    filterResults(q, results);
  };

  function filterResults(q, results) {
    var remote = window['_paletteRemote' + results.id] || [];
    if (remote.length) {
      var existing = results.querySelector('[data-remote-results]');
      if (existing) existing.remove();
      var div = document.createElement('div');
      div.setAttribute('data-remote-results', '');
      var found = 0;
      remote.forEach(function (item) {
        var label = (item.label || item.title || '').toLowerCase();
        var desc = (item.description || '').toLowerCase();
        var section = (item.section || item.group || '').toLowerCase();
        if (q && !label.includes(q) && !desc.includes(q) && !section.includes(q)) return;
        found++;
        var a = document.createElement('a');
        a.className = 'flex items-center gap-3 px-3 py-2.5 rounded-box hover:bg-base-200 cursor-pointer transition-colors';
        a.href = item.href || '#';
        if (item.icon) {
          var icon = document.createElement('span');
          icon.className = 'iconify size-4 shrink-0 text-subtle ' + item.icon;
          a.appendChild(icon);
        }
        var body = document.createElement('div');
        body.className = 'flex-1 min-w-0';
        var title = document.createElement('p');
        title.className = 'text-sm font-medium text-base-content truncate';
        title.textContent = item.label || item.title;
        body.appendChild(title);
        if (item.description) {
          var d = document.createElement('p');
          d.className = 'text-xs text-faint truncate';
          d.textContent = item.description;
          body.appendChild(d);
        }
        a.appendChild(body);
        div.appendChild(a);
      });
      if (found === 0 && q) {
        var empty = document.createElement('p');
        empty.className = 'text-sm text-faint text-center py-4';
        empty.textContent = 'No results found.';
        div.appendChild(empty);
      }
      results.appendChild(div);
    }
    results.querySelectorAll('[data-palette-group]').forEach(function (group) {
      var items = group.querySelectorAll('[data-label]');
      var hasVisible = false;
      items.forEach(function (item) {
        var text = (item.dataset.label || '').toLowerCase();
        var match = !q || text.includes(q);
        item.style.display = match ? '' : 'none';
        if (match) hasVisible = true;
      });
      var header = group.querySelector('p');
      if (header) header.style.display = hasVisible ? '' : 'none';
    });
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initOpen);
  } else {
    initOpen();
  }
})();
