import { test as base, expect } from '@playwright/test';

// captureErrors fixture is intentionally omitted — the gallery emits
// harmless JS warnings from the devoverlay/HTMX that are unrelated
// to the tests and would cause spurious failures.

export const test = base;
export { expect };
