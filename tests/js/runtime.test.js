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
const { JSDOM, VirtualConsole } = require('jsdom');

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

const TOGGLE_GROUP = `
<div role="group" data-gd-toggle-group data-orientation="" data-spacing="2" data-gd-toggle-values="[]">
  <button type="button" data-gd-toggle-value="left" aria-pressed="false">Left</button>
  <button type="button" data-gd-toggle-value="center" aria-pressed="false">Center</button>
  <button type="button" data-gd-toggle-value="right" aria-pressed="false">Right</button>
</div>`;

const JS_FILES = ['go-daisy-a11y.js', 'go-daisy-context-menu.js', 'go-daisy-resizable.js', 'go-daisy-toggle.js', 'go-daisy-messagescroller.js', 'go-daisy-alert-dialog.js', 'go-daisy-menubar.js', 'go-daisy-navigation-menu.js', 'go-daisy-questionnaire.js'];

const MESSAGE_SCROLLER = `
<div data-gd-messagescroller id="scroller">
  <div data-gd-messagescroller-viewport id="viewport"><div>one</div></div>
  <button type="button" data-gd-messagescroller-jump id="jump" hidden aria-label="Scroll to latest"></button>
</div>`;

function loadScripts(dom, extraScripts = []) {
  for (const file of [...JS_FILES, ...extraScripts]) {
    const code = fs.readFileSync(path.join(JS_DIR, 'js', file), 'utf8');
    new Function('window', 'document', code)(dom.window, dom.window.document);
  }
}

function setup(body, extraScripts = [], { onError } = {}) {
  // jsdom reports exceptions thrown inside event listeners to the virtual
  // console instead of rethrowing them, so a test that needs to prove "this
  // listener did not throw" must collect them here.
  const virtualConsole = onError ? new VirtualConsole() : undefined;
  if (virtualConsole) virtualConsole.on('jsdomError', onError);
  const dom = new JSDOM(`<!DOCTYPE html><html><body>${body}</body></html>`, {
    pretendToBeVisual: true,
    url: 'http://localhost/',
    virtualConsole,
  });
  loadScripts(dom, extraScripts);
  return dom;
}

// setupScroller loads the scripts after stubbing the viewport's scroll
// geometry, so the scroller runtime can resolve distances in jsdom (which
// otherwise reports scrollHeight/clientHeight as 0).
function setupScroller(body, { scrollHeight = 100, clientHeight = 100 } = {}) {
  const dom = new JSDOM(`<!DOCTYPE html><html><body>${body}</body></html>`, {
    pretendToBeVisual: true,
    url: 'http://localhost/',
  });
  const viewport = dom.window.document.querySelector('[data-gd-messagescroller-viewport]');
  Object.defineProperty(viewport, 'scrollHeight', { configurable: true, get: () => scrollHeight });
  Object.defineProperty(viewport, 'clientHeight', { configurable: true, get: () => clientHeight });
  loadScripts(dom);
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

test('dropdown: Escape on a focused trigger with an open menu closes it', () => {
  const { window } = setup(DROPDOWN);
  const doc = window.document;
  const trigger = doc.getElementById('trigger');
  const root = doc.getElementById('root');

  // jsdom has no layout, so :focus-within never matches; drive the open state
  // with the explicit marker (daisyUI reaches it via :focus-within in a browser).
  root.setAttribute('data-gd-menu-open', '');
  trigger.setAttribute('aria-expanded', 'true');
  trigger.focus();

  press(window, trigger, 'Escape');

  assert.ok(!root.hasAttribute('data-gd-menu-open'), 'open marker removed');
  assert.equal(trigger.getAttribute('aria-expanded'), 'false');
  assert.ok(root.hasAttribute('data-gd-menu-closed'), 'closed marker set');
  assert.equal(doc.activeElement, trigger, 'focus stays on the trigger');
});

test('dropdown: Escape on a trigger whose menu is closed is a no-op', () => {
  const { window } = setup(DROPDOWN);
  const doc = window.document;
  const trigger = doc.getElementById('trigger');
  const root = doc.getElementById('root');

  // Fresh dropdown: nothing focused, no markers, aria-expanded starts false.
  press(window, trigger, 'Escape');

  assert.ok(!root.hasAttribute('data-gd-menu-closed'), 'no closed marker stamped');
  assert.equal(trigger.getAttribute('aria-expanded'), 'false', 'aria-expanded unchanged');

  press(window, trigger, 'ArrowDown');
  assert.equal(trigger.getAttribute('aria-expanded'), 'true', 'ArrowDown still opens');
  assert.ok(root.hasAttribute('data-gd-menu-open'), 'open marker set by ArrowDown');
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

test('toggle: single-select clears siblings', () => {
  const { window } = setup(TOGGLE_GROUP);
  const items = [...window.document.querySelectorAll('[data-gd-toggle-group] button')];

  items[1].click();
  assert.equal(items[1].getAttribute('aria-pressed'), 'true');
  assert.equal(items[0].getAttribute('aria-pressed'), 'false');
  assert.equal(items[2].getAttribute('aria-pressed'), 'false');

  items[2].click();
  assert.equal(items[2].getAttribute('aria-pressed'), 'true');
  assert.equal(items[1].getAttribute('aria-pressed'), 'false');
});

test('toggle: pressed state also sets btn-active for the visual state', () => {
  const { window } = setup(TOGGLE_GROUP);
  const items = [...window.document.querySelectorAll('[data-gd-toggle-group] button')];

  items[1].click();
  assert.equal(items[1].classList.contains('btn-active'), true);
  assert.equal(items[0].classList.contains('btn-active'), false);

  items[2].click();
  assert.equal(items[1].classList.contains('btn-active'), false, 'single-select clears the previous visual state');
  assert.equal(items[2].classList.contains('btn-active'), true);
});

test('toggle: initial values are applied to aria-pressed and btn-active', () => {
  const markup = TOGGLE_GROUP.replace('data-gd-toggle-values="[]"', 'data-gd-toggle-values="[&quot;center&quot;]"');
  const { window } = setup(markup);
  const items = [...window.document.querySelectorAll('[data-gd-toggle-group] button')];

  assert.equal(items[1].getAttribute('aria-pressed'), 'true');
  assert.equal(items[1].classList.contains('btn-active'), true);
  assert.equal(items[0].getAttribute('aria-pressed'), 'false');
});

test('toggle: multiple keeps both pressed independently', () => {
  const { window } = setup(TOGGLE_GROUP.replace('data-gd-toggle-group', 'data-gd-toggle-group data-multiple'));
  const items = [...window.document.querySelectorAll('[data-gd-toggle-group] button')];

  items[0].click();
  items[2].click();
  assert.equal(items[0].getAttribute('aria-pressed'), 'true');
  assert.equal(items[2].getAttribute('aria-pressed'), 'true');
  assert.equal(items[1].getAttribute('aria-pressed'), 'false');

  items[0].click();
  assert.equal(items[0].getAttribute('aria-pressed'), 'false');
});

test('toggle: controlled group does not mutate aria-pressed', () => {
  const { window } = setup(TOGGLE_GROUP.replace('data-gd-toggle-group', 'data-gd-toggle-group data-controlled'));
  const items = [...window.document.querySelectorAll('[data-gd-toggle-group] button')];

  items[0].click();
  assert.equal(items[0].getAttribute('aria-pressed'), 'false');
});

test('toggle: ArrowRight/ArrowLeft move focus in a horizontal group', () => {
  const { window } = setup(TOGGLE_GROUP);
  const doc = window.document;
  const items = [...doc.querySelectorAll('[data-gd-toggle-group] button')];

  items[0].focus();
  press(window, items[0], 'ArrowRight');
  assert.equal(doc.activeElement, items[1]);
  press(window, items[1], 'ArrowRight');
  assert.equal(doc.activeElement, items[2]);
  press(window, items[2], 'ArrowRight');
  assert.equal(doc.activeElement, items[0], 'navigation wraps');
  assert.equal(items[0].getAttribute('tabindex'), '0', 'roving tabindex follows arrow keys');
  assert.equal(items[2].getAttribute('tabindex'), '-1');
  press(window, items[0], 'ArrowLeft');
  assert.equal(doc.activeElement, items[2], 'navigation wraps backwards');
});

test('toggle: ArrowDown/ArrowUp move focus in a vertical group', () => {
  const { window } = setup(TOGGLE_GROUP.replace('data-orientation=""', 'data-orientation="vertical"'));
  const doc = window.document;
  const items = [...doc.querySelectorAll('[data-gd-toggle-group] button')];

  items[0].focus();
  press(window, items[0], 'ArrowDown');
  assert.equal(doc.activeElement, items[1]);
  press(window, items[1], 'ArrowUp');
  assert.equal(doc.activeElement, items[0]);
});

test('toggle: disabled items are skipped by roving tabindex and arrow navigation', () => {
  const markup = TOGGLE_GROUP.replace(
    '<button type="button" data-gd-toggle-value="left"',
    '<button type="button" disabled data-gd-toggle-value="left"',
  );
  const { window } = setup(markup);
  const doc = window.document;
  const items = [...doc.querySelectorAll('[data-gd-toggle-group] button')];

  assert.notEqual(items[0].getAttribute('tabindex'), '0', 'disabled first item must not own the tab stop');
  assert.equal(items[1].getAttribute('tabindex'), '0');

  items[1].focus();
  press(window, items[1], 'ArrowLeft');
  assert.equal(doc.activeElement, items[2], 'arrow navigation skips the disabled item');
  assert.equal(items[2].getAttribute('tabindex'), '0');
  assert.notEqual(items[0].getAttribute('tabindex'), '0');
});

test('toggle: roving tabindex exposes exactly one tabbable item', () => {
  const { window } = setup(TOGGLE_GROUP);
  const items = [...window.document.querySelectorAll('[data-gd-toggle-group] button')];
  assert.equal(items.filter((el) => el.getAttribute('tabindex') === '0').length, 1);
  assert.equal(items.filter((el) => el.getAttribute('tabindex') === '-1').length, 2);
});

test('toggle: gd-toggle-change fires with ordered values', () => {
  const { window } = setup(TOGGLE_GROUP.replace('data-gd-toggle-group', 'data-gd-toggle-group data-multiple'));
  const root = window.document.querySelector('[data-gd-toggle-group]');
  const items = [...root.querySelectorAll('button')];
  let last = null;
  root.addEventListener('gd-toggle-change', (e) => { last = e.detail.values; });

  items[2].click(); // right first
  items[0].click(); // then left
  assert.deepEqual(last, ['left', 'right'], 'values must be reported in DOM order, not click order');
});

test('message-scroller: jump button hidden at bottom, shown after scrolling up', () => {
  const { window } = setupScroller(MESSAGE_SCROLLER, { scrollHeight: 500, clientHeight: 100 });
  const doc = window.document;
  const viewport = doc.getElementById('viewport');
  const jump = doc.getElementById('jump');

  assert.equal(jump.hidden, true, 'jump button hidden while at the bottom');

  viewport.scrollTop = 0;
  viewport.dispatchEvent(new window.Event('scroll'));
  assert.equal(jump.hidden, false, 'jump button shown after scrolling up');
});

test('message-scroller: jump button click scrolls to bottom and hides', () => {
  const { window } = setupScroller(MESSAGE_SCROLLER, { scrollHeight: 500, clientHeight: 100 });
  const doc = window.document;
  const viewport = doc.getElementById('viewport');
  const jump = doc.getElementById('jump');

  viewport.scrollTop = 0;
  viewport.dispatchEvent(new window.Event('scroll'));
  assert.equal(jump.hidden, false);

  jump.click();
  assert.ok(viewport.scrollTop >= viewport.scrollHeight - viewport.clientHeight, 'click scrolls to the bottom');
  assert.equal(jump.hidden, true, 'button hides after jumping to the bottom');
});

test('message-scroller: static root does not auto-scroll on mutation', async () => {
  const markup = MESSAGE_SCROLLER.replace('data-gd-messagescroller id=', 'data-gd-messagescroller data-gd-messagescroller-static id=');
  const dom = new JSDOM(`<!DOCTYPE html><html><body>${markup}</body></html>`, {
    pretendToBeVisual: true,
    url: 'http://localhost/',
  });
  const { window } = dom;
  const viewport = dom.window.document.getElementById('viewport');

  let height = 100;
  Object.defineProperty(viewport, 'scrollHeight', { configurable: true, get: () => height });
  Object.defineProperty(viewport, 'clientHeight', { configurable: true, get: () => 100 });
  loadScripts(dom);

  assert.equal(viewport.scrollTop, 100, 'init scrolls to the bottom');

  height = 200;
  viewport.appendChild(window.document.createElement('div'));
  await new Promise((resolve) => setTimeout(resolve, 0));

  assert.equal(viewport.scrollTop, 100, 'static root must not auto-scroll on mutation');
});

test('message-scroller: non-static root auto-scrolls on mutation', async () => {
  const dom = new JSDOM(`<!DOCTYPE html><html><body>${MESSAGE_SCROLLER}</body></html>`, {
    pretendToBeVisual: true,
    url: 'http://localhost/',
  });
  const { window } = dom;
  const viewport = dom.window.document.getElementById('viewport');

  let height = 100;
  Object.defineProperty(viewport, 'scrollHeight', { configurable: true, get: () => height });
  Object.defineProperty(viewport, 'clientHeight', { configurable: true, get: () => 100 });
  loadScripts(dom);

  assert.equal(viewport.scrollTop, 100, 'init scrolls to the bottom');

  height = 200;
  viewport.appendChild(window.document.createElement('div'));
  await new Promise((resolve) => setTimeout(resolve, 0));

  assert.equal(viewport.scrollTop, 200, 'sticky root auto-scrolls to the bottom on mutation');
});

const COMMAND_PALETTE = `
<button id="opener" type="button">Search</button>
<input type="checkbox" id="spotlight-toggle" class="modal-toggle" data-command-palette="true">
<div class="modal" role="dialog" aria-modal="true" aria-label="Command palette">
  <div class="modal-box">
    <input type="text" id="spotlight-input" placeholder="Search">
    <div id="spotlight-results">
      <div data-palette-group="Navigate">
        <p>Navigate</p>
        <a id="item-one" href="#one" data-label="Overview">Overview</a>
        <a id="item-two" href="#two" data-label="Hosts">Hosts</a>
      </div>
    </div>
    <label for="spotlight-toggle">Esc</label>
  </div>
  <label class="modal-backdrop" for="spotlight-toggle"></label>
</div>`;

test('command palette: opening focuses the input and traps Tab', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const opener = doc.getElementById('opener');
  const toggle = doc.getElementById('spotlight-toggle');
  const input = doc.getElementById('spotlight-input');
  const last = doc.getElementById('item-two');

  opener.focus();
  window.gdCommandPalette.open('spotlight');

  assert.equal(toggle.checked, true, 'open checks the toggle');
  assert.equal(doc.activeElement, input, 'initial focus lands in the input');

  last.focus();
  const forward = press(window, last, 'Tab');
  assert.equal(forward.defaultPrevented, true, 'Tab from the last item is trapped');
  assert.equal(doc.activeElement, input, 'Tab wraps to the first focusable');

  const backward = press(window, input, 'Tab', { shiftKey: true });
  assert.equal(backward.defaultPrevented, true, 'Shift+Tab from the first item is trapped');
  assert.equal(doc.activeElement, last, 'Shift+Tab wraps to the last focusable');
});

test('command palette: Escape closes and restores focus to the opener', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const opener = doc.getElementById('opener');
  const toggle = doc.getElementById('spotlight-toggle');
  const input = doc.getElementById('spotlight-input');

  opener.focus();
  window.gdCommandPalette.open('spotlight');
  press(window, input, 'Escape');

  assert.equal(toggle.checked, false, 'Escape closes the palette');
  assert.equal(doc.activeElement, opener, 'focus returns to the element that opened it');
});

test('command palette: Cmd/Ctrl+K toggles and restores focus', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const opener = doc.getElementById('opener');
  const toggle = doc.getElementById('spotlight-toggle');
  const input = doc.getElementById('spotlight-input');

  opener.focus();
  const open = press(window, doc.body, 'k', { metaKey: true });
  assert.equal(open.defaultPrevented, true, 'browser find is suppressed');
  assert.equal(toggle.checked, true);
  assert.equal(doc.activeElement, input);

  press(window, doc.body, 'k', { metaKey: true });
  assert.equal(toggle.checked, false);
  assert.equal(doc.activeElement, opener);
});

test('command palette: activating a result closes it and restores focus', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const opener = doc.getElementById('opener');
  const toggle = doc.getElementById('spotlight-toggle');
  const item = doc.getElementById('item-one');

  opener.focus();
  window.gdCommandPalette.open('spotlight');
  item.dispatchEvent(new window.MouseEvent('click', { bubbles: true, cancelable: true }));

  assert.equal(toggle.checked, false, 'clicking a row dismisses the palette');
  assert.equal(doc.activeElement, opener, 'focus does not stay in the hidden dialog');
});

test('command palette: data-command-palette-keep-open leaves the palette open', () => {
  const fixture = COMMAND_PALETTE.replace(
    '<a id="item-two"',
    '<a id="item-keep" href="#keep" data-command-palette-keep-open data-label="Copy">Copy</a><a id="item-two"',
  );
  const { window } = setup(fixture, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const toggle = doc.getElementById('spotlight-toggle');

  window.gdCommandPalette.open('spotlight');
  doc.getElementById('item-keep').dispatchEvent(new window.MouseEvent('click', { bubbles: true, cancelable: true }));

  assert.equal(toggle.checked, true, 'opting out keeps the palette open');
});

test('command palette: filterPaletteItems still filters local rows', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const input = doc.getElementById('spotlight-input');

  input.value = 'over';
  window.filterPaletteItems(input);

  assert.equal(doc.getElementById('item-one').style.display, '', 'matching row stays visible');
  assert.equal(doc.getElementById('item-two').style.display, 'none', 'non-matching row is hidden');
});

test('command palette: opening an already-checked box installs the trap', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const opener = doc.getElementById('opener');
  const toggle = doc.getElementById('spotlight-toggle');
  const input = doc.getElementById('spotlight-input');

  opener.focus();
  toggle.checked = true; // initial Open prop, or app code setting it directly
  window.gdCommandPalette.open('spotlight');

  assert.equal(doc.activeElement, input, 'open() traps even without a change event');
});

test('command palette: close() releases the trap after a silent uncheck', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const opener = doc.getElementById('opener');
  const toggle = doc.getElementById('spotlight-toggle');

  opener.focus();
  window.gdCommandPalette.open('spotlight');

  // Apps such as the infra panel's dashboard.js dismiss with
  // `toggle.checked = false`, which emits no change event; they must then
  // call close() so the trap is released.
  toggle.checked = false;
  window.gdCommandPalette.close('spotlight');

  assert.equal(doc.activeElement, opener, 'close() releases the trap and restores focus');
});

test('command palette: an Open palette inserted by an HTMX settle is trapped', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;

  // Simulate a swap inserting a palette rendered with Open=true. Dispatch on
  // the inserted node so the event has to bubble to the document listener.
  const host = doc.createElement('div');
  host.innerHTML = `
<input type="checkbox" id="late-toggle" class="modal-toggle" data-command-palette="true" checked>
<div class="modal" role="dialog" aria-modal="true" aria-label="Late">
  <div class="modal-box"><input type="text" id="late-input"></div>
</div>`;
  doc.body.appendChild(host);
  host.dispatchEvent(new window.Event('htmx:after:settle', { bubbles: true }));

  assert.equal(doc.activeElement, doc.getElementById('late-input'), 'colon-separated settle traps the inserted palette');

  window.gdCommandPalette.close('late');
  doc.getElementById('late-toggle').checked = true;
  doc.dispatchEvent(new window.Event('htmx:afterSettle', { bubbles: true }));
  assert.equal(doc.activeElement, doc.getElementById('late-input'), 'camelCase settle traps too');
});

test('command palette: reopening after a silent uncheck focuses the input', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const opener = doc.getElementById('opener');
  const toggle = doc.getElementById('spotlight-toggle');
  const input = doc.getElementById('spotlight-input');

  opener.focus();
  window.gdCommandPalette.open('spotlight');
  toggle.checked = false; // silent close: no change event, no close()
  input.blur(); // a hidden dialog drops focus in a real browser
  opener.focus();
  window.gdCommandPalette.open('spotlight');

  assert.equal(doc.activeElement, input, 'a stale trap must not make the reopen skip the input');
});

test('command palette: Cmd+K closes the open palette, not the first', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;

  const host = doc.createElement('div');
  host.innerHTML = `
<input type="checkbox" id="second-toggle" class="modal-toggle" data-command-palette="true">
<div class="modal" role="dialog" aria-modal="true" aria-label="Second">
  <div class="modal-box"><input type="text" id="second-input"></div>
</div>`;
  doc.body.appendChild(host);

  window.gdCommandPalette.open('second');
  press(window, doc.body, 'k', { metaKey: true });

  assert.equal(doc.getElementById('second-toggle').checked, false, 'the open palette closes');
  assert.equal(doc.getElementById('spotlight-toggle').checked, false, 'the first palette is not opened');
});

test('command palette: Escape still closes a palette opened without a trap', () => {
  const { window } = setup(COMMAND_PALETTE, ['go-daisy-command-palette.js']);
  const doc = window.document;
  const toggle = doc.getElementById('spotlight-toggle');

  toggle.checked = true; // legacy path: set directly, no change event, no trap
  press(window, doc.body, 'Escape');

  assert.equal(toggle.checked, false, 'document-level Escape fallback still dismisses');
});

const ALERT_DIALOG = `
<button type="button" data-gd-alert-dialog-open="alert" id="trigger" class="btn">Delete account</button>
<dialog id="alert" data-gd-alert-dialog aria-labelledby="alert-title" aria-describedby="alert-desc" class="modal">
  <form method="dialog" data-slot="alert-dialog-content" class="modal-box">
    <div data-slot="alert-dialog-header" class="flex flex-col gap-1">
      <h2 id="alert-title" class="text-lg font-semibold">Delete account</h2>
      <p id="alert-desc" class="text-sm text-muted">This cannot be undone.</p>
    </div>
    <div data-slot="alert-dialog-footer" class="modal-action">
      <button type="submit" value="cancel" id="cancel" class="btn">Cancel</button>
      <button type="submit" value="confirm" id="confirm" class="btn btn-error">Delete</button>
    </div>
    <button type="button" data-gd-alert-dialog-close id="closebtn" class="btn">Close</button>
  </form>
</dialog>`;

function alertDialogIsOpen(dialog) {
  return dialog.open === true || dialog.getAttribute('open') !== null;
}

function alertDialogIsClosed(dialog) {
  return dialog.open !== true && dialog.getAttribute('open') === null;
}

test('alert-dialog: trigger opens the dialog (showModal or open-attribute fallback)', () => {
  const { window } = setup(ALERT_DIALOG);
  const doc = window.document;
  const dialog = doc.getElementById('alert');
  const trigger = doc.getElementById('trigger');

  const hasShowModal = typeof dialog.showModal === 'function';
  let showModalCalls = 0;
  if (hasShowModal) {
    const original = dialog.showModal.bind(dialog);
    dialog.showModal = function () { showModalCalls++; return original(); };
  }

  trigger.click();

  assert.ok(alertDialogIsOpen(dialog), 'dialog must be open after the trigger click');
  if (hasShowModal) {
    assert.equal(showModalCalls, 1, 'showModal must be called once when supported');
  } else {
    assert.ok(dialog.classList.contains('modal-open'), 'fallback must add the modal-open class');
  }
});

test('alert-dialog: data-gd-alert-dialog-close closes the dialog', () => {
  const { window } = setup(ALERT_DIALOG);
  const doc = window.document;
  const dialog = doc.getElementById('alert');
  const trigger = doc.getElementById('trigger');
  const closeBtn = doc.getElementById('closebtn');

  trigger.click();
  assert.ok(alertDialogIsOpen(dialog), 'dialog must be open before closing');

  closeBtn.click();

  assert.ok(alertDialogIsClosed(dialog), 'close button must close the dialog');
  assert.ok(!dialog.classList.contains('modal-open'), 'fallback close must remove modal-open');
});

test('alert-dialog: close event triggers gd-alert-dialog-close with returnValue', () => {
  const { window } = setup(ALERT_DIALOG);
  const doc = window.document;
  const dialog = doc.getElementById('alert');

  let closedDetail = null;
  dialog.addEventListener('gd-alert-dialog-close', (e) => { closedDetail = e.detail; });

  dialog.returnValue = 'confirm';
  dialog.dispatchEvent(new window.Event('close', { bubbles: false }));

  assert.ok(closedDetail, 'gd-alert-dialog-close must fire on the native close event');
  assert.equal(closedDetail.value, 'confirm', 'detail.value must carry the returnValue');
});

const MENUBAR = `
<div role="menubar" data-gd-menubar id="bar">
  <div data-gd-menubar-menu data-gd-menu-root id="menu-file">
    <button type="button" data-gd-menu-trigger data-gd-menubar-trigger tabindex="0" aria-haspopup="menu" id="file-trigger">File</button>
    <ul role="menu" data-gd-menu data-gd-menubar-content id="file-content">
      <li role="none"><a role="menuitem" data-gd-menu-item id="new" href="#">New</a></li>
      <li role="none"><a role="menuitem" data-gd-menu-item id="open" href="#">Open</a></li>
    </ul>
  </div>
  <div data-gd-menubar-menu data-gd-menu-root id="menu-edit">
    <button type="button" data-gd-menu-trigger data-gd-menubar-trigger tabindex="0" aria-haspopup="menu" id="edit-trigger">Edit</button>
    <ul role="menu" data-gd-menu data-gd-menubar-content id="edit-content">
      <li role="none"><a role="menuitem" data-gd-menu-item id="undo" href="#">Undo</a></li>
      <li role="none"><a role="menuitem" data-gd-menu-item id="redo" href="#">Redo</a></li>
    </ul>
  </div>
  <div data-gd-menubar-menu data-gd-menu-root id="menu-view">
    <button type="button" data-gd-menu-trigger data-gd-menubar-trigger tabindex="0" aria-haspopup="menu" id="view-trigger">View</button>
    <ul role="menu" data-gd-menu data-gd-menubar-content id="view-content">
      <li role="none"><a role="menuitem" data-gd-menu-item id="zoom" href="#">Zoom</a></li>
    </ul>
  </div>
</div>
<button type="button" id="outside">outside</button>`;

const MENUBAR_DISABLED_TRIGGER = `
<div role="menubar" data-gd-menubar id="bar">
  <div data-gd-menubar-menu data-gd-menu-root id="menu-a">
    <button type="button" role="menuitem" data-gd-menu-trigger data-gd-menubar-trigger tabindex="0" aria-haspopup="menu" id="a-trigger">A</button>
    <ul role="menu" data-gd-menu id="a-content">
      <li role="none"><a role="menuitem" data-gd-menu-item id="a-item" href="#">A item</a></li>
    </ul>
  </div>
  <div data-gd-menubar-menu data-gd-menu-root id="menu-b">
    <button type="button" role="menuitem" data-gd-menu-trigger data-gd-menubar-trigger tabindex="0" aria-haspopup="menu" id="b-trigger" disabled>B</button>
    <ul role="menu" data-gd-menu id="b-content">
      <li role="none"><a role="menuitem" data-gd-menu-item id="b-item" href="#">B item</a></li>
    </ul>
  </div>
  <div data-gd-menubar-menu data-gd-menu-root id="menu-c">
    <button type="button" role="menuitem" data-gd-menu-trigger data-gd-menubar-trigger tabindex="0" aria-haspopup="menu" id="c-trigger">C</button>
    <ul role="menu" data-gd-menu id="c-content">
      <li role="none"><a role="menuitem" data-gd-menu-item id="c-item" href="#">C item</a></li>
    </ul>
  </div>
</div>`;

test('menubar: tabbing out and back in does not auto-open a menu', () => {
  // gdMenu's own focusout clears data-gd-menu-closed whenever focus leaves a
  // root; the menubar runtime must re-stamp every closed menu so a returning
  // Tab does not reopen whatever had focus last.
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');
  const outside = doc.getElementById('outside');

  file.focus();
  outside.focus();
  file.focus();

  assert.ok(!doc.getElementById('menu-file').hasAttribute('data-gd-menu-open'), 'bar stays closed');
  assert.ok(doc.getElementById('menu-file').hasAttribute('data-gd-menu-closed'), 'closed marker is re-stamped');
});

// clickTrigger mimics the browser order for a real pointer click: jsdom's
// element.click() neither moves focus nor fires pointerdown, and gdMenu opens
// on focusin, so a bare .click() would prove nothing.
function clickTrigger(window, trigger) {
  trigger.dispatchEvent(new window.MouseEvent('pointerdown', { bubbles: true }));
  trigger.focus();
  trigger.click();
}

test('menubar: clicking a different trigger leaves exactly one menu open', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');
  const edit = doc.getElementById('edit-trigger');

  clickTrigger(window, file);
  clickTrigger(window, edit);

  assert.ok(!doc.getElementById('menu-file').hasAttribute('data-gd-menu-open'), 'previous menu closes');
  assert.ok(doc.getElementById('menu-edit').hasAttribute('data-gd-menu-open'), 'clicked menu opens');
  assert.equal(doc.querySelectorAll('[data-gd-menu-open]').length, 1, 'never two open menus');
});

test('menubar: Home and End jump to the first and last top-level item', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const edit = doc.getElementById('edit-trigger');

  edit.focus();
  press(window, edit, 'End');
  assert.equal(doc.activeElement.id, 'view-trigger');

  press(window, doc.activeElement, 'Home');
  assert.equal(doc.activeElement.id, 'file-trigger');
});

test('menubar: ArrowUp opens the menu at its last item', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');

  file.focus();
  press(window, file, 'ArrowUp');

  assert.equal(doc.activeElement.id, 'open', 'last item is focused');
  assert.ok(doc.getElementById('menu-file').hasAttribute('data-gd-menu-open'));
});

test('menubar: arrow movement skips a disabled top-level item', () => {
  const { window } = setup(MENUBAR_DISABLED_TRIGGER);
  const doc = window.document;
  const a = doc.getElementById('a-trigger');

  a.focus();
  press(window, a, 'ArrowRight');
  assert.equal(doc.activeElement.id, 'c-trigger', 'disabled B is skipped');

  press(window, doc.activeElement, 'ArrowRight');
  assert.equal(doc.activeElement.id, 'a-trigger', 'wraps past the disabled item');
});

test('menubar: clicking a trigger does not throw in the deferred aria-expanded sync', () => {
  // gdMenu's click handler reflects aria-expanded in a setTimeout. Run deferred
  // callbacks synchronously so a wrong isOpen() call site fails here instead of
  // surfacing as an uncaught page error in the browser.
  const errors = [];
  const { window } = setup(MENUBAR, [], { onError: (err) => errors.push(err) });
  const doc = window.document;
  const realSetTimeout = window.setTimeout;
  window.setTimeout = (fn) => {
    fn();
    return 0;
  };
  try {
    doc.getElementById('file-trigger').click();
    doc.getElementById('edit-trigger').click();
  } finally {
    window.setTimeout = realSetTimeout;
  }
  assert.deepEqual(errors.map((e) => e.message), [], 'trigger clicks must not throw');
});

test('menubar: exactly one trigger is tabbable (roving tabindex)', () => {
  const { window } = setup(MENUBAR);
  const triggers = [...window.document.querySelectorAll('[data-gd-menubar-trigger]')];
  assert.equal(triggers.filter((el) => el.getAttribute('tabindex') === '0').length, 1);
  assert.equal(triggers.filter((el) => el.getAttribute('tabindex') === '-1').length, 2);
});

test('menubar: ArrowRight moves focus without opening a menu', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');
  const edit = doc.getElementById('edit-trigger');

  file.focus();
  press(window, file, 'ArrowRight');

  assert.equal(doc.activeElement, edit, 'focus moves to the next trigger');
  assert.equal(doc.querySelector('[data-gd-menu-open]'), null, 'no menu opens on arrow movement');
});

test('menubar: ArrowDown opens the menu and focuses the first item', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');

  file.focus();
  press(window, file, 'ArrowDown');

  assert.equal(doc.activeElement.id, 'new', 'first item is focused');
  assert.ok(doc.getElementById('menu-file').hasAttribute('data-gd-menu-open'), 'menu is open');
  assert.equal(file.getAttribute('aria-expanded'), 'true');
});

test('menubar: ArrowRight inside an open menu opens the adjacent menu', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');

  file.focus();
  press(window, file, 'ArrowDown');
  press(window, doc.activeElement, 'ArrowRight');

  assert.equal(doc.activeElement.id, 'undo', 'first item of the next menu is focused');
  assert.ok(doc.getElementById('menu-edit').hasAttribute('data-gd-menu-open'), 'next menu is open');
  assert.ok(!doc.getElementById('menu-file').hasAttribute('data-gd-menu-open'), 'previous menu is closed');
});

test('menubar: ArrowLeft wraps from the first menu to the last', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');

  file.focus();
  press(window, file, 'ArrowLeft');

  assert.equal(doc.activeElement.id, 'view-trigger', 'wraps to the last trigger');
  assert.equal(doc.querySelector('[data-gd-menu-open]'), null, 'still closed after wrap');
});

test('menubar: Escape closes and restores focus to the trigger', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');

  file.focus();
  press(window, file, 'ArrowDown');
  press(window, doc.activeElement, 'Escape');

  assert.equal(doc.activeElement, file, 'focus returns to the trigger');
  assert.ok(!doc.getElementById('menu-file').hasAttribute('data-gd-menu-open'), 'menu is closed');
});

test('menubar: Escape on a focused open trigger closes it and keeps focus', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');
  const menuFile = doc.getElementById('menu-file');

  // Open without moving focus into the surface (the hover-follow path):
  // focus the trigger, then gdMenu.open with toFirst=false leaves focus on it.
  file.focus();
  window.gdMenu.open(menuFile, false);

  assert.ok(menuFile.hasAttribute('data-gd-menu-open'), 'menu is open');
  assert.equal(doc.activeElement, file, 'focus stays on the trigger while open');

  press(window, file, 'Escape');

  assert.ok(!menuFile.hasAttribute('data-gd-menu-open'), 'menu closes');
  assert.ok(menuFile.hasAttribute('data-gd-menu-closed'), 'closed marker set');
  assert.equal(file.getAttribute('aria-expanded'), 'false', 'aria-expanded cleared');
  assert.equal(doc.activeElement, file, 'focus stays on the trigger');

  // The bar still behaves: arrows move between triggers, ArrowDown reopens.
  press(window, file, 'ArrowRight');
  assert.equal(doc.activeElement.id, 'edit-trigger', 'ArrowRight moves to the next trigger');
  press(window, doc.activeElement, 'ArrowLeft');
  assert.equal(doc.activeElement.id, 'file-trigger', 'ArrowLeft moves back');
  press(window, file, 'ArrowDown');
  assert.ok(menuFile.hasAttribute('data-gd-menu-open'), 'ArrowDown reopens the menu');
});

test('menubar: hovering another trigger while a menu is open follows', () => {
  const { window } = setup(MENUBAR);
  const doc = window.document;
  const file = doc.getElementById('file-trigger');
  const edit = doc.getElementById('edit-trigger');

  file.focus();
  press(window, file, 'ArrowDown');
  edit.dispatchEvent(new window.MouseEvent('pointerover', { bubbles: true }));

  assert.ok(doc.getElementById('menu-edit').hasAttribute('data-gd-menu-open'), 'hovered menu opens');
  assert.ok(!doc.getElementById('menu-file').hasAttribute('data-gd-menu-open'), 'previous menu closes');
});

const NAVIGATION_MENU = `
<nav data-gd-navigation-menu id="nav">
  <ul id="nav-list">
    <li data-gd-nav-item id="item-products">
      <button type="button" data-gd-nav-trigger aria-controls="products-content" id="products-trigger">Products</button>
      <div data-gd-nav-content id="products-content">
        <a data-gd-nav-link id="analytics" href="#">Analytics</a>
        <a data-gd-nav-link id="automation" href="#">Automation</a>
      </div>
    </li>
    <li data-gd-nav-item id="item-solutions">
      <button type="button" data-gd-nav-trigger aria-controls="solutions-content" id="solutions-trigger">Solutions</button>
      <div data-gd-nav-content id="solutions-content">
        <a data-gd-nav-link id="startups" href="#">Startups</a>
      </div>
    </li>
  </ul>
</nav>
<button type="button" id="nav-outside">outside</button>`;

test('navigation menu: runtime marks the root ready and stamps aria-expanded', () => {
  const { window } = setup(NAVIGATION_MENU);
  const doc = window.document;
  assert.equal(doc.getElementById('nav').getAttribute('data-gd-nav-ready'), 'true');
  assert.equal(doc.getElementById('products-trigger').getAttribute('aria-expanded'), 'false');
  assert.equal(doc.getElementById('solutions-trigger').getAttribute('aria-expanded'), 'false');
});

test('navigation menu: click toggles an item and closes the others', () => {
  const { window } = setup(NAVIGATION_MENU);
  const doc = window.document;
  const products = doc.getElementById('products-trigger');
  const solutions = doc.getElementById('solutions-trigger');

  products.click();
  assert.ok(doc.getElementById('item-products').hasAttribute('data-gd-nav-open'));
  assert.equal(products.getAttribute('aria-expanded'), 'true');

  solutions.click();
  assert.ok(doc.getElementById('item-solutions').hasAttribute('data-gd-nav-open'));
  assert.ok(!doc.getElementById('item-products').hasAttribute('data-gd-nav-open'), 'previous item closes');
  assert.equal(products.getAttribute('aria-expanded'), 'false');

  solutions.click();
  assert.ok(!doc.getElementById('item-solutions').hasAttribute('data-gd-nav-open'), 'second click closes');
  assert.equal(solutions.getAttribute('aria-expanded'), 'false');
});

test('navigation menu: Enter toggles once and the synthetic click is consumed', () => {
  const { window } = setup(NAVIGATION_MENU);
  const doc = window.document;
  const products = doc.getElementById('products-trigger');

  press(window, products, 'Enter');
  assert.ok(doc.getElementById('item-products').hasAttribute('data-gd-nav-open'), 'Enter opens');

  // A real button activation fires keydown then click; the click must not undo
  // the keyboard toggle.
  products.click();
  assert.ok(doc.getElementById('item-products').hasAttribute('data-gd-nav-open'), 'synthetic click is consumed');

  // A later genuine click still toggles.
  press(window, products, ' ');
  products.click();
  assert.ok(!doc.getElementById('item-products').hasAttribute('data-gd-nav-open'), 'next click closes');
});

test('navigation menu: ArrowDown opens and focuses the first link', () => {
  const { window } = setup(NAVIGATION_MENU);
  const doc = window.document;
  const products = doc.getElementById('products-trigger');

  press(window, products, 'ArrowDown');

  assert.equal(doc.activeElement.id, 'analytics');
  assert.ok(doc.getElementById('item-products').hasAttribute('data-gd-nav-open'));
});

test('navigation menu: Escape from inside the panel closes and restores focus', () => {
  const { window } = setup(NAVIGATION_MENU);
  const doc = window.document;
  const products = doc.getElementById('products-trigger');

  press(window, products, 'ArrowDown');
  press(window, doc.activeElement, 'Escape');

  assert.equal(doc.activeElement, products);
  assert.ok(!doc.getElementById('item-products').hasAttribute('data-gd-nav-open'));
});

test('navigation menu: ArrowLeft/Right move between triggers, and switch panels when open', () => {
  const { window } = setup(NAVIGATION_MENU);
  const doc = window.document;
  const products = doc.getElementById('products-trigger');
  const solutions = doc.getElementById('solutions-trigger');

  products.focus();
  press(window, products, 'ArrowRight');
  assert.equal(doc.activeElement, solutions, 'moves focus when nothing is open');

  press(window, solutions, 'ArrowRight');
  assert.equal(doc.activeElement, products, 'wraps');

  products.click();
  press(window, products, 'ArrowRight');
  assert.ok(doc.getElementById('item-solutions').hasAttribute('data-gd-nav-open'), 'adjacent panel opens');
  assert.equal(doc.activeElement.id, 'startups', 'first link of the adjacent panel is focused');
});

test('navigation menu: outside pointerdown and focus leaving the nav close it', () => {
  const { window } = setup(NAVIGATION_MENU);
  const doc = window.document;
  const products = doc.getElementById('products-trigger');

  products.click();
  doc.getElementById('nav-outside').dispatchEvent(new window.MouseEvent('pointerdown', { bubbles: true }));
  assert.ok(!doc.getElementById('item-products').hasAttribute('data-gd-nav-open'), 'outside pointerdown closes');

  products.focus();
  products.click();
  assert.ok(doc.getElementById('item-products').hasAttribute('data-gd-nav-open'));
  doc.getElementById('nav-outside').focus();
  assert.ok(!doc.getElementById('item-products').hasAttribute('data-gd-nav-open'), 'focus leaving closes');
});

const QUESTIONNAIRE = `
<form id="q" method="post" data-gd-questionnaire data-gd-questionnaire-mode="letters" class="flex flex-col gap-4">
  <div data-gd-questionnaire-progress class="flex items-center gap-3">
    <span data-gd-questionnaire-progress-text>Question 1 of 3</span>
    <progress data-gd-questionnaire-progress-bar value="1" max="3"></progress>
  </div>
  <fieldset id="q-role" data-gd-questionnaire-item data-name="role" data-required="true">
    <div data-gd-questionnaire-choices>
      <label data-gd-questionnaire-choice><input type="radio" id="q-role-dev" name="role" value="dev"> <kbd data-gd-questionnaire-shortcut>a</kbd></label>
      <label data-gd-questionnaire-choice><input type="radio" id="q-role-design" name="role" value="design"> <kbd data-gd-questionnaire-shortcut>b</kbd></label>
    </div>
    <p id="q-role-error" data-gd-questionnaire-error role="alert" hidden></p>
  </fieldset>
  <fieldset id="q-interests" data-gd-questionnaire-item data-name="interests" data-multiple="true">
    <div data-gd-questionnaire-choices>
      <label data-gd-questionnaire-choice><input type="checkbox" id="q-int-htmx" name="interests" value="htmx"> <kbd data-gd-questionnaire-shortcut>a</kbd></label>
      <label data-gd-questionnaire-choice><input type="checkbox" id="q-int-alpine" name="interests" value="alpine" disabled> <kbd data-gd-questionnaire-shortcut>b</kbd></label>
    </div>
    <p id="q-interests-error" data-gd-questionnaire-error role="alert" hidden></p>
  </fieldset>
  <fieldset id="q-exp" data-gd-questionnaire-item data-name="experience" data-required="true">
    <div data-gd-questionnaire-choices>
      <label data-gd-questionnaire-choice><input type="radio" id="q-exp-beg" name="experience" value="beginner"> <kbd data-gd-questionnaire-shortcut>a</kbd></label>
    </div>
    <p id="q-exp-error" data-gd-questionnaire-error role="alert" hidden></p>
  </fieldset>
  <div>
    <button type="button" data-gd-questionnaire-previous disabled hidden>Previous</button>
    <button type="button" data-gd-questionnaire-skip hidden>Skip</button>
    <button type="button" data-gd-questionnaire-next hidden>Next</button>
    <button type="submit" data-gd-questionnaire-submit>Submit</button>
  </div>
</form>`;

function questionnaireEls(doc) {
  return {
    form: doc.getElementById('q'),
    role: doc.getElementById('q-role'),
    interests: doc.getElementById('q-interests'),
    exp: doc.getElementById('q-exp'),
    prev: doc.querySelector('[data-gd-questionnaire-previous]'),
    skip: doc.querySelector('[data-gd-questionnaire-skip]'),
    next: doc.querySelector('[data-gd-questionnaire-next]'),
    submit: doc.querySelector('[data-gd-questionnaire-submit]'),
    text: doc.querySelector('[data-gd-questionnaire-progress-text]'),
    bar: doc.querySelector('[data-gd-questionnaire-progress-bar]'),
  };
}

test('questionnaire: runtime reveals nav, shows step 1, hides the rest', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);

  assert.equal(els.role.hidden, false, 'step 1 visible');
  assert.equal(els.interests.hidden, true, 'step 2 hidden');
  assert.equal(els.exp.hidden, true, 'step 3 hidden');

  assert.equal(els.prev.hidden, false, 'Previous is revealed');
  assert.equal(els.prev.disabled, true, 'Previous is disabled on the first step');
  assert.equal(els.next.hidden, false, 'Next is revealed');
  assert.equal(els.submit.hidden, true, 'Submit hidden until the last step');

  assert.equal(els.text.textContent, 'Question 1 of 3');
  assert.equal(els.bar.value, 1);
});

test('questionnaire: Next advances and Previous goes back', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);

  doc.getElementById('q-role-dev').click();
  els.next.click();
  assert.equal(els.interests.hidden, false, 'step 2 visible after Next');
  assert.equal(els.role.hidden, true, 'step 1 hidden');
  assert.equal(els.prev.disabled, false, 'Previous enabled on step 2');

  els.prev.click();
  assert.equal(els.role.hidden, false, 'step 1 visible after Previous');
  assert.equal(els.interests.hidden, true);
  assert.equal(els.prev.disabled, true, 'Previous disabled again on step 1');
});

test('questionnaire: required step blocks Next and shows inline error', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);

  els.next.click();

  assert.equal(els.role.hidden, false, 'still on step 1');
  const error = els.role.querySelector('[data-gd-questionnaire-error]');
  assert.equal(error.hidden, false, 'error is revealed');
  assert.match(error.textContent, /Select an option to continue/);
  // aria-invalid belongs on the inputs; a fieldset is not a widget.
  assert.equal(doc.getElementById('q-role-dev').getAttribute('aria-invalid'), 'true');
  assert.equal(doc.getElementById('q-role-design').getAttribute('aria-invalid'), 'true');
  assert.equal(els.role.getAttribute('aria-invalid'), null);
  assert.equal(els.role.getAttribute('aria-describedby'), 'q-role-error', 'error is programmatically associated');
});

test('questionnaire: selecting a choice clears the error and allows Next', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);

  els.next.click();
  assert.equal(els.role.hidden, false, 'blocked before selection');

  doc.getElementById('q-role-dev').click();
  els.next.click();

  assert.equal(els.interests.hidden, false, 'advanced after selection');
  assert.equal(doc.getElementById('q-role-dev').getAttribute('aria-invalid'), null, 'invalid marker cleared');
  assert.equal(els.role.getAttribute('aria-describedby'), null, 'error association cleared');
});

test('questionnaire: a required step offers no Skip', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);

  assert.equal(els.skip.hidden, true, 'Skip hidden on the required step');

  // Even a programmatic skip must not bypass a required answer.
  window.gdQuestionnaire.next(els.form);
  assert.equal(els.role.hidden, false, 'stayed on the required step');

  // On the optional step 2 Skip appears and advances without validation.
  doc.getElementById('q-role-dev').click();
  els.next.click();
  assert.equal(els.interests.hidden, false);
  assert.equal(els.skip.hidden, false, 'Skip revealed on an optional step');
  els.skip.click();
  assert.equal(els.exp.hidden, false, 'Skip advances without validating');
});

test('questionnaire: submit visible only on the last step', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);

  doc.getElementById('q-role-dev').click();
  els.next.click(); // interests (multiple, not required) -> allowed
  els.next.click(); // experience (last step)
  assert.equal(els.exp.hidden, false, 'on the last step');
  assert.equal(els.next.hidden, true, 'Next hidden on the last step');
  assert.equal(els.submit.hidden, false, 'Submit revealed on the last step');
  assert.equal(els.text.textContent, 'Question 3 of 3');
  assert.equal(els.bar.value, 3);
});

test('questionnaire: letter shortcut selects the matching choice', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const dev = doc.getElementById('q-role-dev');
  const design = doc.getElementById('q-role-design');

  dev.focus();
  press(window, dev, 'b');
  assert.equal(dev.checked, false, 'badge "b" selects the second choice');
  assert.equal(design.checked, true, 'letter shortcut checks the second radio');
});

test('questionnaire: gd-questionnaire-step fires with the new index', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);
  const steps = [];
  els.form.addEventListener('gd-questionnaire-step', (e) => steps.push(e.detail.index));

  doc.getElementById('q-role-dev').click();
  els.next.click();
  els.next.click();

  assert.deepEqual(steps, [1, 2], 'step events carry the zero-based index');
});

test('questionnaire: htmx settle rescan restores the persisted index', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);

  doc.getElementById('q-role-dev').click();
  els.next.click(); // step 2
  assert.equal(els.interests.hidden, false);

  // Simulate an htmx settle re-rendering the fieldsets back to their
  // server-side all-visible state; the rescan must restore step 2.
  els.role.hidden = false;
  els.interests.hidden = false;
  els.exp.hidden = false;
  doc.dispatchEvent(new window.Event('htmx:after:settle', { bubbles: true }));

  assert.equal(els.role.hidden, true, 'rescan restores step 2 (step 1 hidden)');
  assert.equal(els.interests.hidden, false, 'step 2 stays the current step');
});

test('questionnaire: implicit submit on a non-final step advances instead of posting', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const form = doc.getElementById('q');

  // Choose a valid answer so Next can advance, then fire the submit event a
  // browser produces for Enter inside the radio group.
  doc.getElementById('q-role-dev').checked = true;
  const evt = new window.Event('submit', { bubbles: true, cancelable: true });
  form.dispatchEvent(evt);

  assert.equal(evt.defaultPrevented, true, 'half-filled wizard must not submit');
  assert.equal(doc.getElementById('q-role').hidden, true);
  assert.equal(doc.getElementById('q-interests').hidden, false, 'advanced to step 2');
  assert.equal(doc.querySelector('[data-gd-questionnaire-progress-text]').textContent, 'Question 2 of 3');
});

test('questionnaire: two wizards on one page keep independent state', () => {
  const second = QUESTIONNAIRE.replace(/\bq-/g, 'r-').replace(/id="q"/, 'id="r"');
  const { window } = setup(QUESTIONNAIRE + second);
  const doc = window.document;

  const q = doc.getElementById('q');
  const r = doc.getElementById('r');
  assert.ok(q && r, 'both forms rendered');

  // Advance the first one only.
  doc.getElementById('q-role-dev').click();
  q.querySelector('[data-gd-questionnaire-next]').click();

  assert.equal(doc.getElementById('q-role').hidden, true, 'first wizard advanced');
  assert.equal(doc.getElementById('r-role').hidden, false, 'second wizard unchanged');
  assert.equal(doc.getElementById('q-interests').hidden, false);
  assert.equal(doc.getElementById('r-interests').hidden, true);
});

test('questionnaire: an invalid submit dispatches no submit event', () => {
  const { window } = setup(QUESTIONNAIRE);
  const doc = window.document;
  const els = questionnaireEls(doc);
  let fired = 0;
  let reachedForm = 0;
  els.form.addEventListener('gd-questionnaire-submit', () => { fired++; });
  // htmx listens on the form itself; a blocked submit must not reach it.
  els.form.addEventListener('submit', () => { reachedForm++; });

  // Walk to the last step with choices filled, then clear the last answer so
  // the final submit is invalid.
  doc.getElementById('q-role-dev').click();
  els.next.click();
  els.next.click();
  assert.equal(els.exp.hidden, false);

  const evt = new window.Event('submit', { bubbles: true, cancelable: true });
  els.form.dispatchEvent(evt);

  assert.equal(evt.defaultPrevented, true, 'invalid submit blocked');
  assert.equal(reachedForm, 0, 'blocked submit never reaches a form-level listener');
  assert.equal(fired, 0, 'no submit intent event for a blocked form');

  doc.getElementById('q-exp-beg').click();
  const ok = new window.Event('submit', { bubbles: true, cancelable: true });
  els.form.dispatchEvent(ok);
  assert.equal(ok.defaultPrevented, false, 'valid submit proceeds');
  assert.equal(fired, 1, 'submit event fires once after validation');
});
