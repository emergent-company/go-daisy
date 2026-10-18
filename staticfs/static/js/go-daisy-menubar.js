/*
 * go-daisy menubar runtime.
 *
 * Builds the WAI-ARIA menubar pattern on top of the shared gdMenu/gdFocus
 * primitives (go-daisy-a11y.js). Each top-level item is a [data-gd-menubar-menu]
 * that is also a [data-gd-menu-root], so gdMenu already owns open/close,
 * aria-expanded, roving focus within the menu, Escape-to-close and
 * click-outside-close. This runtime adds only the menubar-specific bits:
 *
 *   - roving tabindex across the top-level triggers
 *   - ArrowLeft/Right moves between top-level items without opening
 *   - the open-if-any-open rule: once any menu is open, moving to an adjacent
 *     trigger opens it; otherwise the bar stays closed
 *   - hover-follow: hovering a different trigger while a menu is open opens it
 *   - focus-leave closes every open menu
 *
 * Contract:
 *   [data-gd-menubar]                 root, role="menubar"
 *   [data-gd-menubar-menu]            one top-level item (also data-gd-menu-root)
 *   [data-gd-menubar-trigger]         its trigger (also data-gd-menu-trigger)
 *   [data-gd-menu]                    its menu surface (the content)
 *
 * The runtime suppresses gdMenu's open-on-focusin by stamping
 * data-gd-menu-closed on every menu at wire time and again whenever focus
 * leaves the bar (gdMenu clears the marker on its own focusout, which would
 * otherwise make a returning Tab auto-open a menu). Tabbing into the bar or
 * moving focus with the arrows while closed therefore never auto-opens, and
 * gdMenu clears the marker on trigger pointerdown so a real click still opens.
 */
(function () {
  'use strict';
  if (window.__gdMenubarRuntime) return;
  window.__gdMenubarRuntime = true;

  var WIRED = 'data-gd-menubar-wired';

  function menubarRoots() {
    return Array.prototype.slice.call(document.querySelectorAll('[data-gd-menubar]'));
  }

  function menusOf(root) {
    return Array.prototype.slice.call(root.querySelectorAll('[data-gd-menubar-menu]'));
  }

  function triggersOf(root) {
    return Array.prototype.filter.call(
      root.querySelectorAll('[data-gd-menubar-trigger]'),
      function (el) {
        return !el.disabled && el.getAttribute('aria-disabled') !== 'true';
      }
    );
  }

  function isOpen(menu) {
    return menu.hasAttribute('data-gd-menu-open');
  }

  function anyOpen(root) {
    return menusOf(root).some(isOpen);
  }

  function openMenu(menu, toFirst) {
    if (window.gdMenu) window.gdMenu.open(menu, toFirst === true);
  }

  function closeMenu(menu, restoreFocus) {
    if (window.gdMenu) window.gdMenu.close(menu, restoreFocus === true);
  }

  function closeAll(root, except) {
    menusOf(root).forEach(function (menu) {
      if (menu !== except && isOpen(menu)) closeMenu(menu, false);
    });
  }

  // Re-stamp the closed marker on every menu in the bar. gdMenu's focusout
  // handler clears the marker whenever focus leaves a root, so after a Tab-out
  // a returning Tab would auto-open the menu that happened to have focus. The
  // bar owns "stay closed until asked" while focus is outside it.
  function stampClosed(root) {
    menusOf(root).forEach(function (menu) {
      if (!isOpen(menu)) menu.setAttribute('data-gd-menu-closed', '');
    });
  }

  function triggerMenu(trigger) {
    var menu = trigger.closest('[data-gd-menubar-menu]');
    return menu || null;
  }

  // A click on a trigger keeps that item open and closes the rest. gdMenu opens
  // on focusin, so the clicked menu may already be open by the time the click
  // event lands; the bar must never end up with two open menus. Clicking an
  // already-open trigger deliberately does NOT toggle it closed: hover-follow
  // opens the item under the pointer before the click, so a toggle would make a
  // plain "switch to this menu" click close it again.
  function handleTriggerClick(root, trigger) {
    var menu = triggerMenu(trigger);
    if (!menu) return;
    closeAll(root, menu);
  }

  function rovingFor(root) {
    return window.gdFocus.roving(root, '[data-gd-menubar-trigger]');
  }

  // The triggers ship with a static tabindex="0", so gdFocus.roving().prime()
  // (which only fills in missing tabindex) would leave them all tabbable. Force
  // the roving invariant directly: exactly the first trigger owns the tab stop.
  function primeTriggers(root) {
    var triggers = triggersOf(root);
    triggers.forEach(function (el, i) {
      el.setAttribute('tabindex', i === 0 ? '0' : '-1');
    });
  }

  function wire(root) {
    if (root.getAttribute(WIRED) === 'true') return;
    root.setAttribute(WIRED, 'true');
    menusOf(root).forEach(function (menu) {
      // Start closed so a Tab into the bar never auto-opens the first menu, and
      // keep it that way across focus-leave: gdMenu's own focusout would
      // otherwise clear the marker and a returning Tab would reopen the menu.
      menu.setAttribute('data-gd-menu-keep-closed', '');
      menu.setAttribute('data-gd-menu-closed', '');
    });
    primeTriggers(root);
  }

  function init() {
    menubarRoots().forEach(wire);
  }

  function indexOfTarget(triggers, target) {
    var trigger = target.closest('[data-gd-menubar-trigger]');
    if (trigger) return triggers.indexOf(trigger);
    var menu = target.closest('[data-gd-menubar-menu]');
    if (!menu) return -1;
    for (var i = 0; i < triggers.length; i++) {
      if (triggers[i].closest('[data-gd-menubar-menu]') === menu) return i;
    }
    return -1;
  }

  function handleHorizontal(root, event) {
    var triggers = triggersOf(root);
    if (!triggers.length) return;

    var idx = indexOfTarget(triggers, event.target);
    if (idx < 0) return;

    event.preventDefault();
    event.stopPropagation();

    var wasOpen = anyOpen(root);
    var next;
    if (event.key === 'Home') {
      next = 0;
    } else if (event.key === 'End') {
      next = triggers.length - 1;
    } else {
      next = ((event.key === 'ArrowRight' ? idx + 1 : idx - 1) % triggers.length + triggers.length) % triggers.length;
    }
    var nextTrigger = triggers[next];
    var nextMenu = nextTrigger.closest('[data-gd-menubar-menu]');

    closeAll(root, null);

    if (wasOpen) {
      rovingFor(root).at(next);
      openMenu(nextMenu, true);
    } else {
      // Stamp BEFORE focusing: gdMenu's focusin listener opens any menu without
      // the marker, and the arrow keys must not open while the bar is closed.
      nextMenu.setAttribute('data-gd-menu-closed', '');
      rovingFor(root).at(next);
    }
  }

  function bindDocument() {
    document.addEventListener('keydown', function (event) {
      var root = event.target.closest && event.target.closest('[data-gd-menubar]');
      if (!root) return;
      if (event.key !== 'ArrowRight' && event.key !== 'ArrowLeft' &&
          event.key !== 'Home' && event.key !== 'End') return;
      handleHorizontal(root, event);
    });

    document.addEventListener('click', function (event) {
      var root = event.target.closest && event.target.closest('[data-gd-menubar]');
      if (!root) return;
      var trigger = event.target.closest('[data-gd-menubar-trigger]');
      if (!trigger) return;
      handleTriggerClick(root, trigger);
    });

    document.addEventListener('pointerover', function (event) {
      var root = event.target.closest && event.target.closest('[data-gd-menubar]');
      if (!root) return;
      var menu = event.target.closest('[data-gd-menubar-menu]');
      if (!menu) return;
      if (!anyOpen(root) || isOpen(menu)) return;
      closeAll(root, menu);
      openMenu(menu, false);
    });

    document.addEventListener('focusout', function (event) {
      var root = event.target.closest && event.target.closest('[data-gd-menubar]');
      if (!root) return;
      var next = event.relatedTarget;
      if (next && root.contains(next)) return;
      closeAll(root, null);
      stampClosed(root);
    });
  }

  window.gdMenubar = { init: init, rescan: init, closeAll: closeAll };

  if (typeof window.gdBehaviour === 'function') {
    window.gdBehaviour('menubar', init, init);
  } else {
    document.addEventListener('DOMContentLoaded', init);
    document.addEventListener('htmx:after:settle', init);
    document.addEventListener('htmx:afterSettle', init);
    init();
  }

  if (typeof window.gdOnce === 'function') {
    window.gdOnce('menubarDoc', bindDocument);
  } else {
    bindDocument();
  }
})();
