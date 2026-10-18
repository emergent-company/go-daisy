/*
 * go-daisy message scroller runtime.
 *
 * Contract:
 *   [data-gd-messagescroller]           root
 *   [data-gd-messagescroller-viewport]  scrollable content
 *   [data-gd-messagescroller-jump]      jump-to-latest button
 *   data-gd-messagescroller-static      disable stick-to-bottom auto-scroll
 *
 * Sticks the viewport to the bottom as content grows, reveals a jump button
 * when the user scrolls up, and scrolls smoothly to the bottom on click.
 */
(function () {
  'use strict';
  if (window.__gdMessagescrollerRuntime) return;
  window.__gdMessagescrollerRuntime = true;

  var WIRED = 'data-gd-messagescroller-wired';
  var STICKY_EDGE = 48; // px

  function viewportOf(root) {
    return root.querySelector('[data-gd-messagescroller-viewport]');
  }

  function jumpButtonOf(root) {
    return root.querySelector('[data-gd-messagescroller-jump]');
  }

  function distanceFromBottom(viewport) {
    return viewport.scrollHeight - viewport.clientHeight - viewport.scrollTop;
  }

  function isStatic(root) {
    return root.hasAttribute('data-gd-messagescroller-static');
  }

  function isSticky(viewport) {
    return distanceFromBottom(viewport) <= STICKY_EDGE;
  }

  function scrollToBottom(viewport, smooth) {
    if (typeof viewport.scrollTo === 'function') {
      viewport.scrollTo({ top: viewport.scrollHeight, behavior: smooth ? 'smooth' : 'auto' });
    } else {
      viewport.scrollTop = viewport.scrollHeight;
    }
  }

  function updateJumpButton(viewport, button) {
    if (!button) return;
    if (distanceFromBottom(viewport) > STICKY_EDGE) {
      button.hidden = false;
    } else {
      button.hidden = true;
    }
  }

  function init(root) {
    if (!root || root.getAttribute(WIRED) === 'true') return;

    var viewport = viewportOf(root);
    var button = jumpButtonOf(root);
    // Mark wired only once there is something to wire, so a malformed root can
    // be retried after the markup is fixed or swapped.
    if (!viewport) return;
    root.setAttribute(WIRED, 'true');

    var staticMode = isStatic(root);
    var sticky = true;

    scrollToBottom(viewport, false);
    updateJumpButton(viewport, button);

    viewport.addEventListener('scroll', function () {
      sticky = isSticky(viewport);
      updateJumpButton(viewport, button);
    }, { passive: true });

    if (button) {
      button.addEventListener('click', function () {
        sticky = true;
        scrollToBottom(viewport, true);
        updateJumpButton(viewport, button);
      });
    }

    // Stick to bottom as content grows, unless auto-scroll is disabled.
    if (!staticMode && typeof window.MutationObserver === 'function') {
      var observer = new window.MutationObserver(function () {
        if (sticky) {
          scrollToBottom(viewport, false);
        }
        updateJumpButton(viewport, button);
      });
      observer.observe(viewport, { childList: true, subtree: true });
    }
  }

  function rescan() {
    Array.prototype.forEach.call(document.querySelectorAll('[data-gd-messagescroller]'), init);
  }

  window.gdMessageScroller = { init: init, rescan: rescan };

  if (typeof window.gdBehaviour === 'function') {
    // go-daisy-init.js wires gdRescan to htmx settle itself, so init only needs
    // the first pass.
    window.gdBehaviour('messageScroller', function () { rescan(); }, rescan);
  } else {
    document.addEventListener('DOMContentLoaded', rescan);
    document.addEventListener('htmx:after:settle', rescan);
    document.addEventListener('htmx:afterSettle', rescan);
    rescan();
  }
})();
