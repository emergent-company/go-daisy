/*
 * go-daisy context menu runtime.
 *
 * Turns [data-gd-context-root] regions into right-click targets and positions
 * the [data-gd-context-menu] list at the pointer. Keyboard navigation, roving
 * tabindex and typeahead come from go-daisy-a11y.js (gdFocus.roving), so this
 * file only owns opening, positioning and closing.
 *
 * Contract:
 *   [data-gd-context-root]   region that owns a menu
 *   [data-gd-context-menu]   <ul role="menu"> inside that region, hidden by default
 *   [data-gd-context-id]     optional id echoed onto the menu (tests/hooks)
 */
(function () {
  'use strict';
  if (window.__gdContextMenuRuntime) return;
  window.__gdContextMenuRuntime = true;

  var EDGE = 8;
  var openMenu = null;
  var restoreTarget = null;

  function menuOf(root) {
    return root ? root.querySelector('[data-gd-context-menu]') : null;
  }

  function clamp(value, min, max) {
    return Math.max(min, Math.min(value, max));
  }

  function position(menu, x, y) {
    // Measure after unhiding so the size is real rather than 0x0.
    var rect = menu.getBoundingClientRect();
    var left = clamp(x, EDGE, Math.max(EDGE, window.innerWidth - rect.width - EDGE));
    var top = clamp(y, EDGE, Math.max(EDGE, window.innerHeight - rect.height - EDGE));
    menu.style.left = left + 'px';
    menu.style.top = top + 'px';
  }

  function close(restoreFocus) {
    if (!openMenu) return;
    var menu = openMenu;
    openMenu = null;
    menu.hidden = true;
    menu.removeAttribute('data-open');
    var root = menu.closest('[data-gd-context-root]');
    if (root) root.removeAttribute('data-gd-context-active');
    if (restoreFocus && restoreTarget && restoreTarget.focus) {
      try { restoreTarget.focus(); } catch (e) {}
    }
    restoreTarget = null;
  }

  function open(root, menu, x, y) {
    close(false);
    restoreTarget = document.activeElement;
    openMenu = menu;
    menu.hidden = false;
    menu.setAttribute('data-open', '');
    root.setAttribute('data-gd-context-active', '');
    position(menu, x, y);
    if (window.gdFocus) {
      window.gdFocus.roving(menu).first();
    } else {
      var first = menu.querySelector('[role="menuitem"]');
      if (first && first.focus) first.focus();
    }
  }

  document.addEventListener('contextmenu', function (event) {
    var root = event.target.closest && event.target.closest('[data-gd-context-root]');
    if (!root) return;
    var menu = menuOf(root);
    if (!menu) return;
    event.preventDefault();
    open(root, menu, event.clientX, event.clientY);
  });

  // Keyboard invocation: Shift+F10 and the dedicated ContextMenu key open the
  // menu anchored to the focused element, matching native platform behaviour.
  document.addEventListener('keydown', function (event) {
    var isMenuKey = event.key === 'ContextMenu' || (event.shiftKey && event.key === 'F10');
    if (!isMenuKey) return;
    var root = event.target.closest && event.target.closest('[data-gd-context-root]');
    if (!root) return;
    var menu = menuOf(root);
    if (!menu) return;
    event.preventDefault();
    var rect = event.target.getBoundingClientRect();
    open(root, menu, rect.left + rect.width / 2, rect.top + rect.height / 2);
  });

  document.addEventListener('click', function (event) {
    if (!openMenu) return;
    if (openMenu.contains(event.target)) {
      close(true);
      return;
    }
    close(false);
  });

  document.addEventListener('keydown', function (event) {
    if (!openMenu) return;
    if (event.key === 'Escape') {
      event.preventDefault();
      close(true);
    } else if (event.key === 'Tab') {
      close(false);
    }
  });

  window.addEventListener('resize', function () { close(false); });
  document.addEventListener('scroll', function () { close(false); }, true);
  document.addEventListener('htmx:before:request', function () { close(false); });
  document.addEventListener('htmx:beforeRequest', function () { close(false); });

  window.gdContextMenu = { close: function () { close(false); } };
})();
