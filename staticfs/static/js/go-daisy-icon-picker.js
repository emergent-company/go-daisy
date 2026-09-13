/* go-daisy icon picker runtime.
 *
 * Event-delegated on document and scoped to a [data-gd-icon-picker] root, so it
 * survives HTMX swaps and drives any number of picker instances on a page. The
 * Popover runtime owns open/close, Escape, click-away and positioning; this file
 * owns search filtering, roving-tabindex keyboard navigation, committing a
 * choice into the hidden input, and resetting to the default tile.
 */
(function () {
  if (window._gdIconPickerInit) return;
  window._gdIconPickerInit = true;

  function rootOf(node) {
    return node && node.closest ? node.closest("[data-gd-icon-picker]") : null;
  }

  function panelOf(node) {
    var root = rootOf(node);
    return root ? root.querySelector("[data-gd-icon-panel]") : null;
  }

  function options(panel) {
    return panel.querySelectorAll("[data-gd-icon-option]");
  }

  /* the tiles the current search leaves visible, in DOM order */
  function visibleOptions(panel) {
    var all = options(panel);
    var out = [];
    for (var i = 0; i < all.length; i++) {
      if (all[i].style.display !== "none") out.push(all[i]);
    }
    return out;
  }

  /* haystack for the search filter: value name, label and glyph class */
  function optionText(el) {
    return ((el.getAttribute("data-gd-icon-name") || "") + " " +
      (el.getAttribute("data-gd-icon-label") || "") + " " +
      (el.getAttribute("data-gd-icon-class") || "")).toLowerCase();
  }

  function filter(panel) {
    var input = panel.querySelector("[data-gd-icon-search]");
    var q = input ? (input.value || "").trim().toLowerCase() : "";
    var all = options(panel);
    var shown = 0;
    for (var i = 0; i < all.length; i++) {
      var hit = q === "" || optionText(all[i]).indexOf(q) !== -1;
      all[i].style.display = hit ? "" : "none";
      if (hit) shown++;
    }
    var empty = panel.querySelector("[data-gd-icon-empty]");
    if (empty) empty.classList.toggle("hidden", shown > 0);
  }

  function setSelected(panel, value) {
    var all = options(panel);
    for (var i = 0; i < all.length; i++) {
      var active = (all[i].getAttribute("data-gd-icon-value") || "") === value;
      all[i].setAttribute("aria-selected", active ? "true" : "false");
      all[i].setAttribute("tabindex", active ? "0" : "-1");
      all[i].classList.toggle("bg-primary/10", active);
      all[i].classList.toggle("text-primary", active);
      all[i].classList.toggle("ring-1", active);
      all[i].classList.toggle("ring-inset", active);
      all[i].classList.toggle("ring-primary/30", active);
    }
  }

  /* Commit a tile: write the hidden input, refresh the trigger, mark the tile
     selected, then close the popover. */
  function commit(option) {
    var root = rootOf(option);
    var panel = panelOf(option);
    if (!root || !panel) return;
    var value = option.getAttribute("data-gd-icon-value") || "";
    var glyphClass = option.getAttribute("data-gd-icon-class") ||
      root.getAttribute("data-gd-icon-picker-default-class") || "";
    var label = option.getAttribute("data-gd-icon-label") || "";
    var defaultVal = root.getAttribute("data-gd-icon-picker-default-value") || "";
    var isDefault = value === defaultVal;
    var hadFocus = document.activeElement === option;

    var hidden = root.querySelector("[data-gd-icon-input]");
    if (hidden) hidden.value = value;

    var labelEl = root.querySelector("[data-gd-icon-label]");
    if (labelEl) {
      labelEl.textContent = label;
      labelEl.classList.toggle("text-base-content/50", isDefault);
    }
    var glyph = root.querySelector("[data-gd-icon-glyph]");
    if (glyph) {
      while (glyph.firstChild) glyph.removeChild(glyph.firstChild);
      var span = document.createElement("span");
      span.className = "iconify size-5 " + glyphClass;
      span.setAttribute("aria-hidden", "true");
      glyph.appendChild(span);
    }
    setSelected(panel, value);
    var reset = panel.querySelector("[data-gd-icon-reset]");
    if (reset) reset.style.display = isDefault ? "none" : "";

    var trigger = root.querySelector("[data-gd-popover-trigger]");
    if (window.goDaisy && window.goDaisy.popover && trigger) {
      window.goDaisy.popover.close(trigger);
    }
    if (hadFocus && trigger) trigger.focus();
  }

  document.addEventListener("click", function (ev) {
    var option = ev.target.closest ? ev.target.closest("[data-gd-icon-option]") : null;
    if (option && rootOf(option)) {
      ev.preventDefault();
      commit(option);
      return;
    }
    var reset = ev.target.closest ? ev.target.closest("[data-gd-icon-reset]") : null;
    if (reset && rootOf(reset)) {
      ev.preventDefault();
      var panel = panelOf(reset);
      var def = panel ? panel.querySelector("[data-gd-icon-option][data-gd-icon-default]") : null;
      if (def) commit(def);
    }
  });

  document.addEventListener("input", function (ev) {
    var t = ev.target;
    if (t && t.hasAttribute && t.hasAttribute("data-gd-icon-search")) {
      var panel = panelOf(t);
      if (panel) filter(panel);
    }
  });

  /* On open, clear any stale search and focus the search box on fine pointers
     (never on touch: autofocusing would raise the on-screen keyboard). */
  document.addEventListener("click", function (ev) {
    var trigger = ev.target.closest ? ev.target.closest("[data-gd-popover-trigger]") : null;
    if (!trigger) return;
    var root = rootOf(trigger);
    if (!root) return;
    var panel = root.querySelector("[data-gd-icon-panel]");
    if (!panel) return;
    requestAnimationFrame(function () {
      if (panel.getAttribute("data-gd-popover-open") !== "true") return;
      var input = panel.querySelector("[data-gd-icon-search]");
      if (input) {
        input.value = "";
        filter(panel);
      }
      panel.scrollTop = 0;
      if (input && window.matchMedia && window.matchMedia("(pointer: fine)").matches) input.focus();
    });
  });

  /* Roving focus inside the listbox. Buttons commit on Enter/Space natively;
     arrow keys only move focus (they must never close the popover). Arrow
     left/right inside the search input stay caret movement. */
  document.addEventListener("keydown", function (ev) {
    var panel = panelOf(ev.target);
    if (!panel) return;
    var key = ev.key;
    var inSearch = ev.target && ev.target.hasAttribute && ev.target.hasAttribute("data-gd-icon-search");
    if (inSearch) {
      if (key === "Enter") {
        ev.preventDefault();
        var pick = visibleOptions(panel);
        if (pick.length) commit(pick[0]);
      } else if (key === "ArrowDown") {
        ev.preventDefault();
        var first = visibleOptions(panel);
        if (first.length) first[0].focus();
      }
      return;
    }
    var option = ev.target.closest ? ev.target.closest("[data-gd-icon-option]") : null;
    if (!option) return;
    var visible = visibleOptions(panel);
    var idx = visible.indexOf(option);
    if (idx < 0) return;
    var next = idx;
    switch (key) {
      case "ArrowDown":
      case "ArrowRight":
        next = Math.min(idx + 1, visible.length - 1);
        break;
      case "ArrowUp":
      case "ArrowLeft":
        next = Math.max(idx - 1, 0);
        break;
      case "Home":
        next = 0;
        break;
      case "End":
        next = visible.length - 1;
        break;
      default:
        return;
    }
    ev.preventDefault();
    if (visible[next]) visible[next].focus();
  });
})();
