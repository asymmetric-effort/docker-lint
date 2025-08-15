/*
 * Copyright (c) 2025 Asymmetric Effort, LLC. MIT License.
 * Tests for client-side logic.
 */
import test from 'node:test';
import assert from 'node:assert/strict';
import {JSDOM} from 'jsdom';
import {setBusy, load} from './app.js';

test('setBusy toggles aria-busy attribute', () => {
  const dom = new JSDOM('<section id="content"></section>');
  const el = dom.window.document.getElementById('content');
  setBusy(el, true);
  assert.equal(el.getAttribute('aria-busy'), 'true');
});

test('load fetches HTML into article and clears busy state', async () => {
  const dom = new JSDOM('<section id="content" aria-busy="false"><article id="content-article"></article></section>', {url:'http://localhost'});
  const content = dom.window.document.getElementById('content');
  const article = dom.window.document.getElementById('content-article');
  global.fetch = async () => new Response('<h2>Loaded</h2>', {status:200});
  await load('/fake', {content, article});
  assert.equal(article.innerHTML, '<h2>Loaded</h2>');
  assert.equal(content.getAttribute('aria-busy'), 'false');
});
