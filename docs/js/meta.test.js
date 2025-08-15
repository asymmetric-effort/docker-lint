/*
 * Copyright (c) 2025 Asymmetric Effort, LLC. MIT License.
 * Verify meta tag placeholder exists in index.html.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import {JSDOM} from 'jsdom';

const html = fs.readFileSync(new URL('../index.html', import.meta.url), 'utf8');
const dom = new JSDOM(html);

// Test that meta tag placeholder exists
// and is ready for CI substitution.
test('index.html contains commit meta tag placeholder', () => {
  const meta = dom.window.document.head.querySelector("meta[name='docker-lint:commit']");
  assert.ok(meta, 'meta tag docker-lint:commit not found');
  assert.equal(meta.getAttribute('content'), '__GIT_COMMIT__');
});
