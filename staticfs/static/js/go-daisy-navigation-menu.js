/*
 * go-daisy navigation menu runtime.
 *
 * A site navigation panel is a list of links, not a menu, so this implements
 * the APG disclosure navigation pattern instead of menubar/menu roles:
 *
 *   - triggers are ordinary <button aria-expanded aria-controls> elements
 *     (no menubar/menu roles, no aria-haspopup="menu")
 *   - panels are plain <div>s of links
 *   - Enter/Space toggles, Escape closes and restores focus to the trigger,
 *     clicking outside closes, focus leaving the nav closes
 *   - no focus trap and no shared viewport: each panel is absolutely
 *     positioned under its own item
 *
 * Progressive enhancement: the server renders every panel VISIBLE in flow, and
 * this runtime stamps data-gd-nav-ready on the root. Only a rule scoped to
 * [data-gd-nav-ready] hides the panels, so a no-JS visitor reads the links.
 *
 * Contract:
 *   [data-gd-navigation-menu]  root <nav>
 *   [data-gd-nav-item]         one top-level item <li>
 *   [data-gd-nav-trigger]      its <button>
 *   [data-gd-nav-content]      its panel <div>
 *   [data-gd-nav-link]         a link inside a panel
 */
(function () {
  'use strict';
  if (window.__gdNavigationMenuRuntime) return;
  window.__gdNavigationMenuRuntime = true;

  var OPEN_MARKER = 'data-gd-nav-open';
  var KEYDOWN_MARKER = 'data-gd-nav-keydown-toggle';

  // Hover intent is opt-in for hover-capable pointers only. jsdom and old
  // browsers lack matchMedia entirely, so guard its absence rather than assume
  // it exists.
  var canHover = typeof window.matchMedia === 'function' &&
    window.matchMedia('(hover: hover)').matches;

  var hoverState = new WeakMap();

  function roots() {
    return Array.prototype.slice.call(document.querySelectorAll('[data-gd-navigation-menu]'));
  }

  function itemsOf(root) {
    return Array.prototype.slice.call(root.querySelectorAll('[data-gd-nav-item]'));
  }

  function triggersOf(root) {
    return Array.prototype.slice.call(root.querySelectorAll('[data-gd-nav-trigger]'));
  }

  function triggerOf(item) {
    return item.querySelector('[data-gd-nav-trigger]');
  }

  function isOpen(item) {
    return item.hasAttribute(OPEN_MARKER);
  }

  function itemOf(target) {
    if (!target || !target.closest) return null;
    return target.closest('[data-gd-nav-item]');
  }

  function rootOf(el) {
    if (!el || !el.closest) return null;
    return el.closest('[data-gd-navigation-menu]');
  }

  function open(item) {
    if (!item) return;
    closeOthers(rootOf(item), item);
    item.setAttribute(OPEN_MARKER, '');
    var trigger = triggerOf(item);
    if (trigger) trigger.setAttribute('aria-expanded', 'true');
  }

  function close(item, restoreFocus) {
    if (!item) return;
    item.removeAttribute(OPEN_MARKER);
    var trigger = triggerOf(item);
    if (trigger) {
      trigger.setAttribute('aria-expanded', 'false');
      if (restoreFocus) trigger.focus();
    }
  }

  function closeOthers(root, except) {
    if (!root) return;
    itemsOf(root).forEach(function (item) {
      if (item !== except && isOpen(item)) close(item, false);
    });
  }

  function closeAll(root) {
    closeOthers(root, null);
  }

  function toggle(item) {
    if (!item) return;
    if (isOpen(item)) close(item, false);
    else open(item);
  }

  function firstLink(item) {
    return item.querySelector('[data-gd-nav-link]');
  }

  function stateFor(root) {
    var s = hoverState.get(root);
    if (!s) {
      s = { openTimer: 0, closeTimer: 0 };
      hoverState.set(root, s);
    }
    return s;
  }

  function clearTimers(root) {
    var s = hoverState.get(root);
    if (!s) return;
    if (s.openTimer) { clearTimeout(s.openTimer); s.openTimer = 0; }
    if (s.closeTimer) { clearTimeout(s.closeTimer); s.closeTimer = 0; }
  }

  function indexOfTrigger(triggers, target) {
    var trigger = target.closest('[data-gd-nav-trigger]');
    if (trigger) return triggers.indexOf(trigger);
    var item = itemOf(target);
    if (!item) return -1;
    for (var i = 0; i < triggers.length; i++) {
      if (itemOf(triggers[i]) === item) return i;
    }
    return -1;
  }

  function anyOpen(root) {
    return itemsOf(root).some(isOpen);
  }

  function handleHorizontal(root, event) {
    var triggers = triggersOf(root);
    if (!triggers.length) return;
    var idx = indexOfTrigger(triggers, event.target);
    if (idx < 0) return;
    event.preventDefault();
    event.stopPropagation();

    var next;
    if (event.key === 'Home') {
      next = 0;
    } else if (event.key === 'End') {
      next = triggers.length - 1;
    } else {
      var dir = event.key === 'ArrowRight' ? 1 : -1;
      next = ((idx + dir) % triggers.length + triggers.length) % triggers.length;
    }
    var nextTrigger = triggers[next];
    var nextItem = itemOf(nextTrigger);

    if (anyOpen(root)) {
      // Once any item is open, moving horizontally opens the adjacent item and
      // focuses its first link.
      open(nextItem);
      var link = firstLink(nextItem);
      if (link) link.focus();
    } else {
      nextTrigger.focus();
    }
  }

  function handleTriggerKeydown(root, event) {
    var trigger = event.target.closest('[data-gd-nav-trigger]');
    if (!trigger) return;
    var item = itemOf(trigger);

    if (event.key === 'Enter' || event.key === ' ') {
      // Handle Enter/Space explicitly so the synthetic click the browser fires
      // for button activation does not toggle a second time. The marker is
      // consumed (or cleared) on the next click/pointerdown.
      event.preventDefault();
      trigger.setAttribute(KEYDOWN_MARKER, '');
      toggle(item);
    } else if (event.key === 'Escape') {
      close(item, true);
    } else if (event.key === 'ArrowDown') {
      event.preventDefault();
      open(item);
      var link = firstLink(item);
      if (link) link.focus();
    } else if (event.key === 'ArrowRight' || event.key === 'ArrowLeft' ||
               event.key === 'Home' || event.key === 'End') {
      handleHorizontal(root, event);
    }
  }

  function handleContentKeydown(event) {
    if (event.key !== 'Escape') return;
    var item = itemOf(event.target);
    if (!item) return;
    close(item, true);
  }

  function handleClick(event) {
    var trigger = event.target.closest && event.target.closest('[data-gd-nav-trigger]');
    if (!trigger) return;
    // Consume the one click that follows an Enter/Space keydown toggle.
    if (trigger.hasAttribute(KEYDOWN_MARKER)) {
      trigger.removeAttribute(KEYDOWN_MARKER);
      return;
    }
    toggle(itemOf(trigger));
  }

  function handleOutsidePointerdown(event) {
    var target = event.target;
    if (!target || !target.closest) return;
    if (target.closest('[data-gd-navigation-menu]')) return;
    roots().forEach(function (root) {
      closeAll(root);
      clearTimers(root);
    });
  }

  function handleFocusOut(event) {
    var root = event.target.closest && event.target.closest('[data-gd-navigation-menu]');
    if (!root) return;
    var next = event.relatedTarget;
    if (next && root.contains(next)) return;
    closeAll(root);
    clearTimers(root);
  }

  function scheduleOpen(root, item) {
    var s = stateFor(root);
    if (s.closeTimer) { clearTimeout(s.closeTimer); s.closeTimer = 0; }
    if (s.openTimer) clearTimeout(s.openTimer);
    s.openTimer = setTimeout(function () {
      s.openTimer = 0;
      if (!root.isConnected) return;
      open(item);
    }, 120);
  }

  function scheduleClose(root) {
    var s = stateFor(root);
    if (s.openTimer) { clearTimeout(s.openTimer); s.openTimer = 0; }
    if (s.closeTimer) clearTimeout(s.closeTimer);
    s.closeTimer = setTimeout(function () {
      s.closeTimer = 0;
      if (!root.isConnected) return;
      if (root.contains(document.activeElement)) return;
      closeAll(root);
    }, 200);
  }

  // pointerover/pointerout bubble (unlike pointerenter/pointerleave), which is
  // what a document-level delegation needs; the closest/relatedTarget guards
  // reproduce pointerenter/pointerleave semantics.
  function handlePointerOver(event) {
    if (!canHover) return;
    var root = event.target.closest && event.target.closest('[data-gd-navigation-menu]');
    if (!root) return;
    var item = itemOf(event.target);
    if (!item) return;
    scheduleOpen(root, item);
  }

  function handlePointerOut(event) {
    if (!canHover) return;
    var root = event.target.closest && event.target.closest('[data-gd-navigation-menu]');
    if (!root) return;
    var related = event.relatedTarget;
    if (related && root.contains(related)) return;
    scheduleClose(root);
  }

  // A real pointer interaction invalidates a pending keyboard-click suppression
  // so a genuine mouse click still toggles after an Enter/Space activation.
  function handlePointerDown(event) {
    var trigger = event.target.closest && event.target.closest('[data-gd-nav-trigger]');
    if (trigger) trigger.removeAttribute(KEYDOWN_MARKER);
  }

  function init() {
    roots().forEach(function (root) {
      root.setAttribute('data-gd-nav-ready', 'true');
    });
  }

  function bindDocument() {
    document.addEventListener('keydown', function (event) {
      var root = event.target.closest && event.target.closest('[data-gd-navigation-menu]');
      if (!root) return;
      if (event.target.closest('[data-gd-nav-content]')) {
        handleContentKeydown(event);
        return;
      }
      handleTriggerKeydown(root, event);
    });

    document.addEventListener('click', handleClick);
    document.addEventListener('pointerdown', handlePointerDown);
    document.addEventListener('pointerdown', handleOutsidePointerdown);
    document.addEventListener('focusout', handleFocusOut);
    document.addEventListener('pointerover', handlePointerOver);
    document.addEventListener('pointerout', handlePointerOut);
  }

  window.gdNavigationMenu = {
    open: open,
    close: close,
    closeAll: closeAll,
  };

  if (typeof window.gdBehaviour === 'function') {
    window.gdBehaviour('navigationMenu', init, init);
  } else {
    document.addEventListener('DOMContentLoaded', init);
    document.addEventListener('htmx:after:settle', init);
    document.addEventListener('htmx:afterSettle', init);
    init();
  }

  if (typeof window.gdOnce === 'function') {
    window.gdOnce('navigationMenuDoc', bindDocument);
  } else {
    bindDocument();
  }
})();
