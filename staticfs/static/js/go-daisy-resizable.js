/*
 * go-daisy resizable split-view runtime.
 *
 * Contract:
 *   [data-gd-resizable]              root, carries --gd-split
 *   data-orientation="horizontal|vertical"
 *   data-initial / data-min / data-max   percentages
 *   data-persist                     optional localStorage key
 *   [data-gd-resizable-start]        first pane
 *   [data-gd-resizable-handle]       role="separator" divider
 *
 * Pointer drag, keyboard resize (arrows, Home/End, Shift = coarse), and
 * double-click reset. aria-valuenow tracks the split so screen readers report
 * the current position.
 */
(function () {
  'use strict';
  if (window.__gdResizableRuntime) return;
  window.__gdResizableRuntime = true;

  var FINE = 2;
  var COARSE = 10;

  function num(el, name, fallback) {
    var raw = el.getAttribute('data-' + name);
    var value = parseFloat(raw);
    return isNaN(value) ? fallback : value;
  }

  function splitOf(root) {
    var value = parseFloat(getComputedStyle(root).getPropertyValue('--gd-split'));
    return isNaN(value) ? 50 : value;
  }

  function setSplit(root, pct) {
    var min = num(root, 'min', 10);
    var max = num(root, 'max', 90);
    var value = Math.max(min, Math.min(max, pct));
    root.style.setProperty('--gd-split', value.toFixed(2) + '%');
    var handle = root.querySelector('[data-gd-resizable-handle]');
    if (handle) handle.setAttribute('aria-valuenow', String(Math.round(value)));
    if (root.getAttribute('data-persist')) {
      try { localStorage.setItem(root.getAttribute('data-persist'), value.toFixed(2)); } catch (e) {}
    }
    return value;
  }

  function isVertical(root) {
    return root.getAttribute('data-orientation') === 'vertical';
  }

  function init(root) {
    if (!root || root.getAttribute('data-gd-resizable-wired') === 'true') return;
    var handle = root.querySelector('[data-gd-resizable-handle]');
    if (!handle) return;
    root.setAttribute('data-gd-resizable-wired', 'true');

    if (root.getAttribute('data-persist')) {
      try {
        var stored = parseFloat(localStorage.getItem(root.getAttribute('data-persist')));
        if (!isNaN(stored)) setSplit(root, stored);
      } catch (e) {}
    }

    var dragging = false;

    function fromPointer(event) {
      var rect = root.getBoundingClientRect();
      var pct = isVertical(root)
        ? ((event.clientY - rect.top) / rect.height) * 100
        : ((event.clientX - rect.left) / rect.width) * 100;
      setSplit(root, pct);
    }

    handle.addEventListener('pointerdown', function (event) {
      if (event.button !== undefined && event.button !== 0) return;
      dragging = true;
      event.preventDefault();
      try { handle.setPointerCapture(event.pointerId); } catch (e) {}
      document.documentElement.setAttribute('data-gd-resizing', isVertical(root) ? 'row' : 'col');
      fromPointer(event);
    });

    handle.addEventListener('pointermove', function (event) {
      if (!dragging) return;
      fromPointer(event);
    });

    function endDrag(event) {
      if (!dragging) return;
      dragging = false;
      try { handle.releasePointerCapture(event.pointerId); } catch (e) {}
      document.documentElement.removeAttribute('data-gd-resizing');
    }

    handle.addEventListener('pointerup', endDrag);
    handle.addEventListener('pointercancel', endDrag);

    handle.addEventListener('keydown', function (event) {
      var step = event.shiftKey ? COARSE : FINE;
      var current = splitOf(root);
      var next = null;
      var vertical = isVertical(root);
      switch (event.key) {
        case 'ArrowLeft': if (!vertical) next = current - step; break;
        case 'ArrowRight': if (!vertical) next = current + step; break;
        case 'ArrowUp': if (vertical) next = current - step; break;
        case 'ArrowDown': if (vertical) next = current + step; break;
        case 'Home': next = num(root, 'min', 10); break;
        case 'End': next = num(root, 'max', 90); break;
        default: return;
      }
      if (next === null) return;
      event.preventDefault();
      setSplit(root, next);
    });

    handle.addEventListener('dblclick', function () {
      setSplit(root, num(root, 'initial', 50));
    });
  }

  function rescan() {
    Array.prototype.forEach.call(document.querySelectorAll('[data-gd-resizable]'), init);
  }

  window.gdResizable = { init: init, rescan: rescan };

  document.addEventListener('htmx:after:settle', rescan);
  document.addEventListener('htmx:afterSettle', rescan);
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', rescan);
  } else {
    rescan();
  }
})();
