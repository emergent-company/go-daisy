/*
 * go-daisy questionnaire runtime.
 *
 * A progressive-enhancement wizard over one native <form>. Every question is a
 * <fieldset> rendered server-side; with JavaScript off all are visible and the
 * last step's Submit posts the whole form. This runtime takes over the DOM to
 * drive step navigation, inline validation, and keyboard shortcuts.
 *
 * Contract:
 *   [data-gd-questionnaire]               root <form>
 *   [data-gd-questionnaire-item]          one <fieldset> step
 *   [data-gd-questionnaire-choice]        one labelled <input> row
 *   [data-gd-questionnaire-progress-text] "Question X of N" label
 *   [data-gd-questionnaire-progress-bar]  <progress> bar
 *   [data-gd-questionnaire-previous]      Previous button (hidden server-side)
 *   [data-gd-questionnaire-skip]          Skip button (hidden server-side)
 *   [data-gd-questionnaire-next]          Next button (hidden server-side)
 *   [data-gd-questionnaire-submit]        Submit button (visible server-side)
 *   [data-gd-questionnaire-error]         inline error <p> (hidden)
 *   [data-gd-questionnaire-shortcut]      kbd shortcut badge
 *
 * No-JS contract: only the nav buttons are hidden server-side. Every fieldset
 * is rendered visible, so a no-JS visitor reads every question and submits the
 * whole form with the Submit button.
 */
(function () {
  'use strict';
  if (window.__gdQuestionnaireRuntime) return;
  window.__gdQuestionnaireRuntime = true;

  var READY_MARKER = 'data-gd-questionnaire-ready';
  var INDEX_MARKER = 'data-gd-questionnaire-index';

  function roots() {
    return Array.prototype.slice.call(document.querySelectorAll('[data-gd-questionnaire]'));
  }

  function itemsOf(root) {
    return Array.prototype.slice.call(root.querySelectorAll('[data-gd-questionnaire-item]'));
  }

  function rootOf(el) {
    if (!el || !el.closest) return null;
    return el.closest('[data-gd-questionnaire]');
  }

  function choiceInputs(item) {
    return Array.prototype.slice.call(item.querySelectorAll('input[type=radio], input[type=checkbox]'));
  }

  function enabledInputs(item) {
    return choiceInputs(item).filter(function (el) { return !el.disabled; });
  }

  function isRequired(item) {
    return item.getAttribute('data-required') === 'true';
  }

  function isMultiple(item) {
    return item.getAttribute('data-multiple') === 'true';
  }

  function errorOf(item) {
    return item.querySelector('[data-gd-questionnaire-error]');
  }

  function prevBtn(root) { return root.querySelector('[data-gd-questionnaire-previous]'); }
  function skipBtn(root) { return root.querySelector('[data-gd-questionnaire-skip]'); }
  function nextBtn(root) { return root.querySelector('[data-gd-questionnaire-next]'); }
  function submitBtn(root) { return root.querySelector('[data-gd-questionnaire-submit]'); }
  function progressText(root) { return root.querySelector('[data-gd-questionnaire-progress-text]'); }
  function progressBar(root) { return root.querySelector('[data-gd-questionnaire-progress-bar]'); }

  function currentIndex(root) {
    var items = itemsOf(root);
    for (var i = 0; i < items.length; i++) {
      if (!items[i].hidden) return i;
    }
    return 0;
  }

  function clearError(item) {
    var err = errorOf(item);
    if (err) {
      err.textContent = '';
      err.hidden = true;
    }
    item.removeAttribute('aria-invalid');
  }

  function showError(item, message) {
    var err = errorOf(item);
    if (err) {
      err.textContent = message;
      err.hidden = false;
    }
    item.setAttribute('aria-invalid', 'true');
    var first = enabledInputs(item)[0];
    if (first) first.focus();
  }

  function validate(item) {
    var inputs = enabledInputs(item);
    if (isRequired(item) && !inputs.some(function (el) { return el.checked; })) {
      showError(item, isMultiple(item)
        ? 'Select at least one option to continue.'
        : 'Select an option to continue.');
      return false;
    }
    clearError(item);
    return true;
  }

  function dispatch(root, name, detail) {
    var evt;
    try {
      evt = new window.CustomEvent(name, { bubbles: true, detail: detail });
    } catch (e) {
      evt = document.createEvent('Event');
      evt.initEvent(name, true, false);
      evt.detail = detail;
    }
    root.dispatchEvent(evt);
  }

  function update(root) {
    var items = itemsOf(root);
    var idx = currentIndex(root);
    var total = items.length;

    var text = progressText(root);
    if (text) text.textContent = 'Question ' + (idx + 1) + ' of ' + total;
    var bar = progressBar(root);
    if (bar) {
      bar.max = total < 1 ? 1 : total;
      bar.value = idx + 1;
    }

    var prev = prevBtn(root);
    if (prev) {
      prev.removeAttribute('hidden');
      if (idx === 0) prev.setAttribute('disabled', '');
      else prev.removeAttribute('disabled');
    }

    var skip = skipBtn(root);
    if (skip) {
      if (idx < total - 1) skip.removeAttribute('hidden');
      else skip.setAttribute('hidden', '');
    }

    var next = nextBtn(root);
    if (next) {
      if (idx < total - 1) next.removeAttribute('hidden');
      else next.setAttribute('hidden', '');
    }

    var submit = submitBtn(root);
    if (submit) {
      if (idx < total - 1) submit.setAttribute('hidden', '');
      else submit.removeAttribute('hidden');
    }
  }

  // apply shows exactly one step, persists the index on the root (so an htmx
  // settle that re-renders the fieldsets does not reset to step 1), syncs the
  // nav/progress, and dispatches the step event unless silent.
  function apply(root, index, silent) {
    var items = itemsOf(root);
    if (!items.length) return;
    if (index < 0) index = 0;
    if (index >= items.length) index = items.length - 1;
    items.forEach(function (item, i) { item.hidden = i !== index; });
    root.setAttribute(INDEX_MARKER, String(index));
    update(root);
    if (!silent) dispatch(root, 'gd-questionnaire-step', { index: index });
  }

  function goTo(root, index) { apply(root, index, false); }

  function next(root) {
    var items = itemsOf(root);
    var idx = currentIndex(root);
    if (idx >= items.length - 1) return;
    if (!validate(items[idx])) return;
    goTo(root, idx + 1);
  }

  function previous(root) {
    var items = itemsOf(root);
    var idx = currentIndex(root);
    if (idx <= 0) return;
    clearError(items[idx]);
    goTo(root, idx - 1);
  }

  function skip(root) {
    var items = itemsOf(root);
    var idx = currentIndex(root);
    if (idx >= items.length - 1) return;
    clearError(items[idx]);
    goTo(root, idx + 1);
  }

  function handleChoiceChange(event) {
    var target = event.target;
    if (!target || !target.closest) return;
    var root = rootOf(target);
    if (!root) return;
    var item = target.closest('[data-gd-questionnaire-item]');
    if (!item) return;
    clearError(item);
    var items = itemsOf(root);
    dispatch(root, 'gd-questionnaire-change', {
      index: items.indexOf(item),
      name: item.getAttribute('data-name') || '',
      value: target.value,
    });
  }

  function shortcutMatch(item, mode, key) {
    var wanted = mode === 'letters' ? key.toLowerCase() : key;
    var labels = Array.prototype.slice.call(item.querySelectorAll('[data-gd-questionnaire-choice]'));
    for (var i = 0; i < labels.length; i++) {
      var input = labels[i].querySelector('input[type=radio], input[type=checkbox]');
      if (!input || input.disabled) continue;
      var kbd = labels[i].querySelector('[data-gd-questionnaire-shortcut]');
      if (!kbd) continue;
      var badge = (kbd.textContent || '').trim();
      if (mode === 'letters' && badge.toLowerCase() === wanted) return input;
      if (mode === 'numbers' && badge === wanted) return input;
    }
    return null;
  }

  function selectChoice(input) {
    if (input.type === 'radio') {
      input.checked = true;
    } else {
      input.checked = !input.checked;
    }
    input.dispatchEvent(new window.Event('change', { bubbles: true }));
    input.dispatchEvent(new window.Event('input', { bubbles: true }));
    input.focus();
  }

  function handleKeydown(event) {
    if (event.altKey || event.ctrlKey || event.metaKey) return;
    var target = event.target;
    var tag = (target && target.tagName || '').toLowerCase();
    if (tag === 'textarea' || tag === 'select') return;
    if (tag === 'input' && target.type !== 'radio' && target.type !== 'checkbox') return;
    if (target && target.isContentEditable) return;

    var root = rootOf(target);
    if (!root) return;

    var mode = root.getAttribute('data-gd-questionnaire-mode');
    if (mode !== 'letters' && mode !== 'numbers') return;

    var key = event.key;
    if (mode === 'letters') {
      if (!/^[a-zA-Z]$/.test(key)) return;
    } else if (!/^[1-9]$/.test(key)) {
      return;
    }

    var items = itemsOf(root);
    var item = items[currentIndex(root)];
    if (!item) return;

    var input = shortcutMatch(item, mode, key);
    if (!input) return;
    event.preventDefault();
    selectChoice(input);
  }

  function handleClick(event) {
    var btn = event.target.closest && event.target.closest('[data-gd-questionnaire-previous], [data-gd-questionnaire-skip], [data-gd-questionnaire-next]');
    if (!btn) return;
    var root = rootOf(btn);
    if (!root) return;
    if (btn.hasAttribute('data-gd-questionnaire-previous')) previous(root);
    else if (btn.hasAttribute('data-gd-questionnaire-skip')) skip(root);
    else next(root);
  }

  function handleSubmit(event) {
    var root = event.target;
    if (!root || !root.hasAttribute || !root.hasAttribute('data-gd-questionnaire')) return;
    var items = itemsOf(root);
    var idx = currentIndex(root);
    var item = items[idx];
    if (!item) return;

    // Implicit submission: the Submit button is hidden on every step but the
    // last, yet a hidden submit button is still the form's default button, so
    // Enter in a radio group fires submit. On a non-final step that must behave
    // like Next (validate and advance) instead of posting a half-filled wizard.
    if (idx < items.length - 1) {
      event.preventDefault();
      next(root);
      return;
    }

    dispatch(root, 'gd-questionnaire-submit', { index: idx });
    if (!validate(item)) {
      event.preventDefault();
    }
  }

  function init() {
    roots().forEach(function (root) {
      var wasReady = root.getAttribute(READY_MARKER) === 'true';
      root.setAttribute(READY_MARKER, 'true');

      if (!wasReady) {
        apply(root, 0, true);
        return;
      }

      // Rescan after an htmx settle: restore the persisted index if the step
      // fieldsets were re-rendered, rather than resetting to step 1.
      var saved = parseInt(root.getAttribute(INDEX_MARKER), 10);
      if (!isNaN(saved)) {
        apply(root, saved, true);
      } else {
        update(root);
      }
    });
  }

  function bindDocument() {
    document.addEventListener('change', handleChoiceChange);
    document.addEventListener('click', handleClick);
    document.addEventListener('keydown', handleKeydown);
    document.addEventListener('submit', handleSubmit);
    document.addEventListener('DOMContentLoaded', init);
  }

  // Public API for HTMX response scripts; root is the <form> argument.
  window.gdQuestionnaire = {
    next: function (root) { next(root); },
    previous: function (root) { previous(root); },
    goTo: function (root, index) { goTo(root, index); },
  };

  if (typeof window.gdBehaviour === 'function') {
    window.gdBehaviour('questionnaire', init, init);
  } else {
    document.addEventListener('DOMContentLoaded', init);
    document.addEventListener('htmx:after:settle', init);
    document.addEventListener('htmx:afterSettle', init);
    init();
  }

  if (typeof window.gdOnce === 'function') {
    window.gdOnce('questionnaireDoc', bindDocument);
  } else {
    bindDocument();
  }
})();
