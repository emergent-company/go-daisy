/**
 * DOM runtime tests for the bundled component scripts.
 *
 * The Go tests only prove markup is emitted; these tests prove the JavaScript
 * that consumes that markup actually works. They run under jsdom against the
 * same files staticfs ships, with fixtures that mirror the component output
 * (see components/ui/dropdown.templ, context-menu.templ, resizable.templ).
 *
 * Run: node --test tests/js/runtime.test.js   (or `task test:js`)
 */
'use strict';

const fs = require('node:fs');
const path = require('node:path');
const test = require('node:test');
const assert = require('node:assert/strict');
const { JSDOM } = require('jsdom');

const JS_DIR = path.resolve(__dirname, '..', '..', 'staticfs', 'static');

const DROPDOWN = `
<div data-gd-menu-root id="root">
  <button data-gd-menu-trigger id="trigger">Actions</button>
  <ul role="menu" data-gd-menu id="menu">
    <li role="none"><a role="menuitem" data-gd-menu-item id="rename" href="#">Rename</a></li>
    <li role="none"><a role="menuitem" data-gd-menu-item id="duplicate" href="#">Duplicate</a></li>
    <li role="none"><a role="menuitem" data-gd-menu-item id="delete" href="#">Delete</a></li>
  </ul>
</div>`;

const CONTEXT_MENU = `
<div data-gd-context-root id="ctx">
  <div id="region" tabindex="0">right click me</div>
  <ul role="menu" data-gd-context-menu id="ctxmenu" hidden>
    <li role="none"><a role="menuitem" data-gd-menu-item id="open" href="#">Open</a></li>
    <li role="none"><a role="menuitem" data-gd-menu-item id="remove" href="#">Remove</a></li>
  </ul>
</div>`;

const RESIZABLE = `
<div data-gd-resizable id="split" data-orientation="horizontal"
     data-initial="50" data-min="10" data-max="90" style="--gd-split:50%">
  <div data-gd-resizable-start>start</div>
  <div role="separator" tabindex="0" data-gd-resizable-handle aria-valuenow="50"></div>
  <div>end</div>
</div>`;

function setup(body) {
  const dom = new JSDOM(`<!DOCTYPE html><html><body>${body}</body></html>`, {
    pretendToBeVisual: true,
    url: 'http://localhost/',
  });
  // The scripts are written for a browser global scope. jsdom's `window.eval`
  // does not bind `window` inside the evaluated code, so run them with the
  // window/document in scope explicitly — same globals a page provides.
  for (const file of ['go-daisy-a11y.js', 'go-daisy-context-menu.js', 'go-daisy-resizable.js']) {
    const code = fs.readFileSync(path.join(JS_DIR, 'js', file), 'utf8');
    new Function('window', 'document', code)(dom.window, dom.window.document);
  }
  return dom;
}

function press(window, target, key, opts = {}) {
  const event = new window.KeyboardEvent('keydown', {
    key,
    bubbles: true,
    cancelable: true,
    ...opts,
  });
  target.dispatchEvent(event);
  return event;
}

test('dropdown: ArrowDown opens, roving focus and typeahead work', () => {
  const { window } = setup(DROPDOWN);
  const doc = window.document;
  const trigger = doc.getElementById('trigger');

  press(window, trigger, 'ArrowDown');
  assert.equal(doc.activeElement.id, 'rename', 'first item should be focused');
  assert.equal(trigger.getAttribute('aria-expanded'), 'true');

  press(window, doc.activeElement, 'ArrowDown');
  assert.equal(doc.activeElement.id, 'duplicate');
  press(window, doc.activeElement, 'End');
  assert.equal(doc.activeElement.id, 'delete');
  press(window, doc.activeElement, 'ArrowDown');
  assert.equal(doc.activeElement.id, 'rename', 'navigation wraps');
  press(window, doc.activeElement, 'ArrowUp');
  assert.equal(doc.activeElement.id, 'delete', 'navigation wraps backwards');

  press(window, doc.activeElement, 'd');
  assert.equal(doc.activeElement.id, 'duplicate', 'typeahead matches by prefix');
});

test('dropdown: exactly one menu item is tabbable', () => {
  const { window } = setup(DROPDOWN);
  const items = [...window.document.querySelectorAll('[role="menuitem"]')];
  assert.equal(items.filter((el) => el.getAttribute('tabindex') === '0').length, 1);
  assert.equal(items.filter((el) => el.getAttribute('tabindex') === '-1').length, 2);
});

test('dropdown: Escape closes, restores focus and stays closed', () => {
  const { window } = setup(DROPDOWN);
  const doc = window.document;
  const trigger = doc.getElementById('trigger');
  const root = doc.getElementById('root');

  press(window, trigger, 'ArrowDown');
  press(window, doc.activeElement, 'Escape');

  assert.equal(doc.activeElement, trigger, 'focus returns to the trigger');
  assert.equal(trigger.getAttribute('aria-expanded'), 'false');
  assert.ok(root.hasAttribute('data-gd-menu-closed'), 'closed marker blocks :focus-within reopening');
});

test('dropdown: clicking outside closes the open menu', () => {
  const { window } = setup(`${DROPDOWN}<button id="outside">elsewhere</button>`);
  const doc = window.document;
  const trigger = doc.getElementById('trigger');

  press(window, trigger, 'ArrowDown');
  assert.equal(trigger.getAttribute('aria-expanded'), 'true');

  doc.getElementById('outside').dispatchEvent(new window.MouseEvent('click', { bubbles: true }));
  assert.equal(trigger.getAttribute('aria-expanded'), 'false');
});

test('context menu: opens at the pointer with the first item focused', () => {
  const { window } = setup(CONTEXT_MENU);
  const doc = window.document;
  const menu = doc.getElementById('ctxmenu');

  const event = new window.MouseEvent('contextmenu', {
    bubbles: true,
    cancelable: true,
    clientX: 40,
    clientY: 30,
  });
  doc.getElementById('region').dispatchEvent(event);

  assert.ok(event.defaultPrevented, 'native menu is suppressed');
  assert.equal(menu.hidden, false, 'menu is revealed (hidden attribute is the toggle)');
  assert.notEqual(menu.style.left, '', 'menu is positioned');
  assert.equal(doc.activeElement.id, 'open', 'first item receives focus');
});

test('context menu: Escape closes and Shift+F10 opens from the keyboard', () => {
  const { window } = setup(CONTEXT_MENU);
  const doc = window.document;
  const menu = doc.getElementById('ctxmenu');
  const region = doc.getElementById('region');

  region.focus();
  const openEvent = press(window, region, 'F10', { shiftKey: true });
  assert.ok(openEvent.defaultPrevented);
  assert.equal(menu.hidden, false);

  press(window, doc.activeElement, 'Escape');
  assert.equal(menu.hidden, true, 'Escape closes the menu');
});

test('resizable: keyboard resizes, clamps and updates aria-valuenow', () => {
  const { window } = setup(RESIZABLE);
  const doc = window.document;
  const split = doc.getElementById('split');
  const handle = split.querySelector('[data-gd-resizable-handle]');

  const splitValue = () => parseFloat(split.style.getPropertyValue('--gd-split'));

  handle.focus();
  press(window, handle, 'ArrowRight');
  assert.equal(splitValue(), 52);

  press(window, handle, 'ArrowRight', { shiftKey: true });
  assert.equal(splitValue(), 62);

  press(window, handle, 'Home');
  assert.equal(splitValue(), 10, 'Home clamps to data-min');

  press(window, handle, 'End');
  assert.equal(splitValue(), 90, 'End clamps to data-max');

  press(window, handle, 'ArrowRight');
  assert.equal(splitValue(), 90, 'cannot exceed data-max');
  assert.equal(handle.getAttribute('aria-valuenow'), '90');
});

test('resizable: vertical orientation uses ArrowUp/ArrowDown', () => {
  const { window } = setup(RESIZABLE.replace('horizontal', 'vertical'));
  const doc = window.document;
  const split = doc.getElementById('split');
  const handle = split.querySelector('[data-gd-resizable-handle]');

  const splitValue = () => parseFloat(split.style.getPropertyValue('--gd-split'));

  handle.focus();
  press(window, handle, 'ArrowDown');
  assert.equal(splitValue(), 52);
  press(window, handle, 'ArrowUp');
  assert.equal(splitValue(), 50);
});
