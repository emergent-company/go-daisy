/*
 * go-daisy a11y primitives.
 *
 * One small, dependency-free layer for the keyboard and focus behaviour that
 * every menu-like widget needs. Components opt in with data attributes and a
 * tiny inline script that calls gdMenu / gdFocus, so the semantics stay in one
 * place instead of being re-implemented per component.
 *
 * Public surface:
 *   gdFocus.trap(container, { onEscape, initial })  -> release()
 *   gdFocus.roving(container, itemSelector)         -> { first, last, next, prev, at, index }
 *   gdFocus.typeahead()                             -> match(items, labelOf, ch)
 *   gdMenu.init(root)                               -> wires one menu root
 *   gdMenu.rescan()                                 -> wires every unwired [data-gd-menu-root]
 *
 * Menu contract (attribute driven):
 *   [data-gd-menu-root]                 container; may be any element
 *   [data-gd-menu-trigger]              button/role=button that toggles
 *   [data-gd-menu]                      the menu surface (role=menu)
 *   [data-gd-menu-item] or role=menuitem
 *
 * Behaviour: ArrowDown/Up/Home/End roving focus, Enter/Space activation,
 * Escape closes and restores focus, Tab closes and leaves, printable-key
 * typeahead, click-outside closes, aria-expanded kept in sync.
 */
(function () {
  'use strict';

  if (window.gdFocus && window.gdMenu) return;

  var SEARCH_RESET_MS = 600;

  function isVisible(el) {
    if (!el || el.nodeType !== 1) return false;
    if (el.hasAttribute('disabled') || el.getAttribute('aria-disabled') === 'true') return false;
    if (el.hidden) return false;
    var style = window.getComputedStyle(el);
    return style.display !== 'none' && style.visibility !== 'hidden';
  }

  function focusableIn(container) {
    var sel = 'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]),' +
      ' textarea:not([disabled]), [tabindex]:not([tabindex="-1"]), [contenteditable="true"]';
    return Array.prototype.filter.call(container.querySelectorAll(sel), isVisible);
  }

  /* ---------------------------------------------------------------- gdFocus */

  var gdFocus = window.gdFocus || {};

  /**
   * trap(container, { onEscape, initial }) keeps Tab focus inside container,
   * focuses `initial` (or the first focusable child), and returns release()
   * which restores the previously focused element.
   */
  gdFocus.trap = function (container, opts) {
    if (!container) return function () {};
    opts = opts || {};
    var previous = document.activeElement;
    var target = opts.initial || focusableIn(container)[0] || container;

    if (!container.hasAttribute('tabindex') && target === container) {
      container.setAttribute('tabindex', '-1');
    }
    try { target.focus(); } catch (e) {}

    function onKeydown(e) {
      if (e.key === 'Escape' && opts.onEscape) {
        e.stopPropagation();
        opts.onEscape(e);
        return;
      }
      if (e.key !== 'Tab') return;
      var items = focusableIn(container);
      if (!items.length) { e.preventDefault(); return; }
      var first = items[0];
      var last = items[items.length - 1];
      if (e.shiftKey && (document.activeElement === first || document.activeElement === container)) {
        e.preventDefault();
        last.focus();
      } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault();
        first.focus();
      }
    }

    container.addEventListener('keydown', onKeydown);

    function release() {
      container.removeEventListener('keydown', onKeydown);
      if (opts.restore !== false && previous && previous.focus) {
        try { previous.focus(); } catch (e) {}
      }
    }
    return release;
  };

  /**
   * roving(container, itemSelector) implements the WAI-ARIA roving tabindex
   * pattern: exactly one item is tabbable, arrows move focus with wraparound.
   */
  gdFocus.roving = function (container, itemSelector) {
    var selector = itemSelector || '[data-gd-menu-item], [role="menuitem"], [role="option"]';

    function items() {
      return Array.prototype.filter.call(container.querySelectorAll(selector), isVisible);
    }

    function at(index) {
      var list = items();
      if (!list.length) return null;
      var i = ((index % list.length) + list.length) % list.length;
      for (var n = 0; n < list.length; n++) {
        list[n].setAttribute('tabindex', n === i ? '0' : '-1');
      }
      try { list[i].focus(); } catch (e) {}
      return list[i];
    }

    function current() {
      var list = items();
      var index = list.indexOf(document.activeElement);
      return index;
    }

    function prime() {
      var list = items();
      list.forEach(function (el, i) {
        if (!el.hasAttribute('tabindex')) el.setAttribute('tabindex', i === 0 ? '0' : '-1');
      });
      return list;
    }

    return {
      items: items,
      prime: prime,
      index: current,
      at: at,
      first: function () { return at(0); },
      last: function () { return at(-1); },
      next: function () { return at(current() + 1); },
      prev: function () { return at(current() - 1); },
    };
  };

  /**
   * typeahead() returns a stateful matcher: feed it each printable character
   * and it resolves the next item whose text starts with the accumulated
   * buffer. Resets after SEARCH_RESET_MS of inactivity.
   */
  gdFocus.typeahead = function () {
    var buffer = '';
    var last = 0;
    return function (items, labelOf, ch) {
      var now = Date.now();
      if (now - last > SEARCH_RESET_MS) buffer = '';
      last = now;
      buffer += ch.toLowerCase();
      var needle = buffer;
      // Repeated identical characters cycle: "ddd" moves to the third d-item.
      var repeated = buffer.split('').every(function (c) { return c === buffer[0]; });
      var start = 0;
      if (repeated && buffer.length > 1) {
        needle = buffer[0];
        var active = items.indexOf(document.activeElement);
        start = active >= 0 ? active + 1 : 0;
      }
      for (var n = 0; n < items.length; n++) {
        var item = items[(start + n) % items.length];
        var label = (labelOf ? labelOf(item) : item.textContent) || '';
        if (label.trim().toLowerCase().indexOf(needle) === 0) return item;
      }
      return null;
    };
  };

  window.gdFocus = gdFocus;

  /* ----------------------------------------------------------------- gdMenu */

  function triggerOf(root) { return root.querySelector('[data-gd-menu-trigger]'); }
  function menuOf(root) { return root.querySelector('[data-gd-menu]'); }

  function setExpanded(root, open) {
    var trigger = triggerOf(root);
    if (trigger) trigger.setAttribute('aria-expanded', open ? 'true' : 'false');
    if (open) root.setAttribute('data-gd-menu-open', ''); else root.removeAttribute('data-gd-menu-open');
  }

  function close(root, restoreFocus) {
    var trigger = triggerOf(root);
    setExpanded(root, false);
    var active = document.activeElement;
    if (active && root.contains(active)) active.blur();
    if (restoreFocus && trigger && trigger.focus) {
      try { trigger.focus(); } catch (e) {}
    }
  }

  function open(root, toFirst) {
    var trigger = triggerOf(root);
    if (trigger) {
      try { trigger.focus(); } catch (e) {}
    }
    setExpanded(root, true);
    if (toFirst) {
      var menu = menuOf(root);
      if (menu) gdFocus.roving(menu).first();
    }
  }

  function wire(root) {
    if (root.getAttribute('data-gd-menu-wired') === 'true') return;
    var trigger = triggerOf(root);
    var menu = menuOf(root);
    if (!trigger || !menu) return;
    root.setAttribute('data-gd-menu-wired', 'true');

    var roving = gdFocus.roving(menu);
    roving.prime();
    var typeahead = gdFocus.typeahead();
    trigger.setAttribute('aria-expanded', 'false');
    if (!trigger.hasAttribute('aria-haspopup')) trigger.setAttribute('aria-haspopup', 'menu');

    function items() { return roving.items(); }
    function labelOf(el) { return el.getAttribute('data-gd-menu-label') || el.textContent; }

    function isOpen() {
      return root.hasAttribute('data-gd-menu-open') ||
        root.matches(':focus-within') ||
        root.contains(document.activeElement);
    }

    trigger.addEventListener('keydown', function (e) {
      if (e.key === 'ArrowDown' || e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        open(root, true);
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        open(root, false);
        roving.last();
      }
    });

    trigger.addEventListener('click', function () {
      // The trigger owns focus on click; reflect it in aria-expanded.
      window.setTimeout(function () { setExpanded(root, isOpen()); }, 0);
    });

    menu.addEventListener('keydown', function (e) {
      var list = items();
      if (!list.length) return;
      var handled = true;
      switch (e.key) {
        case 'ArrowDown': roving.next(); break;
        case 'ArrowUp': roving.prev(); break;
        case 'Home': roving.first(); break;
        case 'End': roving.last(); break;
        case 'Escape': close(root, true); break;
        case 'Tab': close(root, false); handled = false; break;
        default:
          if (e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey) {
            var hit = typeahead(list, labelOf, e.key);
            if (hit) { hit.focus(); } else { handled = false; }
          } else {
            handled = false;
          }
      }
      if (handled) e.preventDefault();
    });

    menu.addEventListener('focusout', function (e) {
      var next = e.relatedTarget;
      if (!next || !root.contains(next)) {
        window.setTimeout(function () {
          if (!root.contains(document.activeElement)) setExpanded(root, false);
        }, 0);
      }
    });

    document.addEventListener('click', function (e) {
      if (!root.contains(e.target)) close(root, false);
    });

    // Reflect the CSS-driven open state (DaisyUI opens on :focus-within).
    root.addEventListener('focusin', function () { setExpanded(root, true); });
  }

  var gdMenu = window.gdMenu || {
    init: wire,
    rescan: function () {
      Array.prototype.forEach.call(
        document.querySelectorAll('[data-gd-menu-root]'),
        wire
      );
    },
    close: function (root) { close(root, false); },
  };

  window.gdMenu = gdMenu;

  if (!window.__gdMenuBound) {
    window.__gdMenuBound = true;
    document.addEventListener('htmx:after:settle', gdMenu.rescan);
    document.addEventListener('htmx:afterSettle', gdMenu.rescan);
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', gdMenu.rescan);
    } else {
      gdMenu.rescan();
    }
  }
})();
