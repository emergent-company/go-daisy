/*
 * go-daisy alert dialog runtime.
 *
 * Contract:
 *   [data-gd-alert-dialog-open]     button(s) that open the dialog by id
 *   [data-gd-alert-dialog-close]    button(s) that close their enclosing dialog
 *   dialog[data-gd-alert-dialog]    the <dialog> element itself
 *
 * Built on the native <dialog> element: showModal() provides the top layer,
 * focus trap, Escape-to-close and ::backdrop for free, so this runtime only
 * wires the trigger/close paths and rebroadcasts the native `close` event as a
 * bubbling `gd-alert-dialog-close` CustomEvent (detail.value = dialog.returnValue).
 *
 * An alert dialog must NOT close on backdrop click — the user has to pick an
 * action — so there is deliberately no outside-click handler here.
 *
 * Plain document-level delegation is swap-safe: the listeners survive htmx
 * swaps and bind by attribute at click time, so there is nothing to rescan and
 * no per-node wiring to re-apply after a swap. gdAlertDialog.init/rescan exist
 * only to mirror the shape of the other runtimes; both are no-ops.
 */
(function () {
  'use strict';
  if (window.__gdAlertDialogRuntime) return;
  window.__gdAlertDialogRuntime = true;

  function openDialog(dialog) {
    // showModal() throws InvalidStateError on an already-open dialog.
    if (dialog.open || dialog.hasAttribute('open')) return;
    if (typeof dialog.showModal === 'function') {
      dialog.showModal();
    } else {
      dialog.setAttribute('open', '');
      if (dialog.classList) dialog.classList.add('modal-open');
    }
    focusFirst(dialog);
  }

  function closeDialog(dialog) {
    if (typeof dialog.close === 'function') {
      dialog.close();
    } else {
      dialog.removeAttribute('open');
      if (dialog.classList) dialog.classList.remove('modal-open');
    }
  }

  // Focus the first focusable element inside the content so keyboard users land
  // somewhere useful. The native modal already traps focus; this covers the
  // no-showModal fallback and lets jsdom tests observe the initial focus.
  function focusFirst(dialog) {
    var root = dialog.querySelector('[data-slot="alert-dialog-content"]') || dialog;
    var candidates = root.querySelectorAll(
      'button:not([disabled]), [href], input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])',
    );
    for (var i = 0; i < candidates.length; i++) {
      var el = candidates[i];
      if (el.tabIndex === -1) continue;
      if (typeof el.focus === 'function') {
        el.focus();
        return;
      }
    }
  }

  function emitClose(dialog) {
    var detail = { value: dialog.returnValue || '' };
    var event;
    if (typeof window.CustomEvent === 'function') {
      event = new window.CustomEvent('gd-alert-dialog-close', { bubbles: true, detail: detail });
    } else {
      event = document.createEvent('CustomEvent');
      event.initCustomEvent('gd-alert-dialog-close', true, false, detail);
    }
    dialog.dispatchEvent(event);
  }

  document.addEventListener('click', function (event) {
    var target = event.target;
    if (!target || typeof target.closest !== 'function') return;

    var opener = target.closest('[data-gd-alert-dialog-open]');
    if (opener) {
      var dialog = document.getElementById(opener.getAttribute('data-gd-alert-dialog-open'));
      if (dialog) openDialog(dialog);
      return;
    }

    var closer = target.closest('[data-gd-alert-dialog-close]');
    if (closer) {
      var dlg = closer.closest('dialog');
      if (dlg) closeDialog(dlg);
    }
  });

  // The native `close` event (from Escape, a method="dialog" submit button, or
  // dialog.close()) does not bubble, so capture it at the document and
  // rebroadcast as a bubbling CustomEvent apps can listen for.
  document.addEventListener('close', function (event) {
    var dialog = event.target;
    if (!dialog || dialog.tagName !== 'DIALOG' || !dialog.hasAttribute('data-gd-alert-dialog')) {
      return;
    }
    emitClose(dialog);
  }, true);

  window.gdAlertDialog = {
    init: function () {},
    rescan: function () {},
  };
})();
