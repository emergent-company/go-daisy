/*
 * go-daisy toggle group runtime.
 *
 * Contract:
 *   [data-gd-toggle-group]            root, role="group"
 *   button[data-gd-toggle-value]      items (single string value)
 *   data-orientation="vertical"       stacked (default horizontal)
 *   data-multiple                     multi-select (default single)
 *   data-controlled                   server owns aria-pressed (no client mutation)
 *   data-gd-toggle-values (JSON)      initial pressed values on the root
 *
 * Roving tabindex, Arrow/Home/End focus movement, click-toggled aria-pressed,
 * and a bubbling gd-toggle-change CustomEvent carrying the ordered pressed
 * values in DOM order.
 */
(function () {
  'use strict';
  if (window.__gdToggleRuntime) return;
  window.__gdToggleRuntime = true;

  var WIRED = 'data-gd-toggle-wired';

  function itemsOf(root) {
    return Array.prototype.slice.call(root.querySelectorAll('button[data-gd-toggle-value]'));
  }

  // Disabled buttons are not focusable, so they must never own the tab stop or
  // receive focus from arrow-key navigation.
  function navigableItemsOf(root) {
    return itemsOf(root).filter(function (item) { return !item.disabled; });
  }

  function isPressed(item) {
    return item.getAttribute('aria-pressed') === 'true';
  }

  // Keep the visual pressed state (btn-active) in sync with aria-pressed so an
  // uncontrolled group looks toggled without a server round-trip.
  function setPressed(item, on) {
    item.setAttribute('aria-pressed', on ? 'true' : 'false');
    if (item.classList) item.classList.toggle('btn-active', on);
  }

  function isVertical(root) {
    return root.getAttribute('data-orientation') === 'vertical';
  }

  function isMultiple(root) {
    return root.hasAttribute('data-multiple');
  }

  function isControlled(root) {
    return root.hasAttribute('data-controlled');
  }

  function valuesOf(root) {
    var values = [];
    var items = itemsOf(root);
    for (var i = 0; i < items.length; i++) {
      if (isPressed(items[i])) values.push(items[i].getAttribute('data-gd-toggle-value'));
    }
    return values;
  }

  function updateTabindex(root) {
    var items = navigableItemsOf(root);
    if (items.length === 0) return;
    var target = null;
    for (var i = 0; i < items.length; i++) {
      if (isPressed(items[i])) { target = items[i]; break; }
    }
    if (!target) target = items[0];
    for (var j = 0; j < items.length; j++) {
      items[j].setAttribute('tabindex', items[j] === target ? '0' : '-1');
    }
  }

  // Read the root's initial JSON and mark matching items pressed. Skipped for
  // controlled groups — the server owns aria-pressed there.
  function applyInitial(root) {
    var raw = root.getAttribute('data-gd-toggle-values');
    if (!raw) return;
    var initial;
    try { initial = JSON.parse(raw); } catch (e) { return; }
    if (!Array.isArray(initial)) return;
    var set = {};
    for (var i = 0; i < initial.length; i++) set[initial[i]] = true;
    var items = itemsOf(root);
    for (var j = 0; j < items.length; j++) {
      setPressed(items[j], Boolean(set[items[j].getAttribute('data-gd-toggle-value')]));
    }
  }

  function emitChange(root) {
    var detail = { values: valuesOf(root) };
    var event;
    if (typeof window.CustomEvent === 'function') {
      event = new window.CustomEvent('gd-toggle-change', { bubbles: true, detail: detail });
    } else {
      event = document.createEvent('CustomEvent');
      event.initCustomEvent('gd-toggle-change', true, false, detail);
    }
    root.dispatchEvent(event);
  }

  function onClick(root, item) {
    if (item.disabled) return;
    if (isControlled(root)) return; // server owns state
    if (isMultiple(root)) {
      setPressed(item, !isPressed(item));
    } else {
      var items = itemsOf(root);
      for (var i = 0; i < items.length; i++) {
        setPressed(items[i], items[i] === item);
      }
    }
    updateTabindex(root);
    emitChange(root);
  }

  function onKeydown(root, event) {
    var items = navigableItemsOf(root);
    if (items.length === 0) return;
    var vertical = isVertical(root);
    var current = Array.prototype.indexOf.call(items, event.target);
    if (current === -1) current = 0;
    var next = null;
    switch (event.key) {
      case 'ArrowRight':
        if (!vertical) next = (current + 1) % items.length;
        break;
      case 'ArrowLeft':
        if (!vertical) next = (current - 1 + items.length) % items.length;
        break;
      case 'ArrowDown':
        if (vertical) next = (current + 1) % items.length;
        break;
      case 'ArrowUp':
        if (vertical) next = (current - 1 + items.length) % items.length;
        break;
      case 'Home':
        next = 0;
        break;
      case 'End':
        next = items.length - 1;
        break;
      default:
        return;
    }
    if (next === null) return;
    event.preventDefault();
    for (var j = 0; j < items.length; j++) {
      items[j].setAttribute('tabindex', j === next ? '0' : '-1');
    }
    items[next].focus();
  }

  // A swap that replaces only the group's children leaves the root (and its
  // WIRED marker) in place, so new items are not wired. Re-rendering the group
  // root itself is the supported refresh path.
  function init(root) {
    if (!root || root.getAttribute(WIRED) === 'true') return;
    root.setAttribute(WIRED, 'true');

    if (!isControlled(root)) applyInitial(root);

    var items = itemsOf(root);
    for (var i = 0; i < items.length; i++) {
      (function (item) {
        item.addEventListener('click', function () { onClick(root, item); });
        item.addEventListener('keydown', function (event) { onKeydown(root, event); });
      })(items[i]);
    }
    updateTabindex(root);
  }

  function rescan() {
    Array.prototype.forEach.call(document.querySelectorAll('[data-gd-toggle-group]'), init);
  }

  window.gdToggle = { init: init, rescan: rescan };

  if (typeof window.gdBehaviour === 'function') {
    // go-daisy-init.js wires gdRescan to htmx settle itself, so init only needs
    // the first pass.
    window.gdBehaviour('toggleGroup', function () { rescan(); }, rescan);
  } else {
    document.addEventListener('DOMContentLoaded', rescan);
    document.addEventListener('htmx:after:settle', rescan);
    document.addEventListener('htmx:afterSettle', rescan);
    rescan();
  }
})();
