/*
 * Copyright (c) 2025 Asymmetric Effort, LLC.
 * Licensed under the MIT License.
 */
import {JSDOM} from 'jsdom';
import assert from 'node:assert/strict';
import {setBusy, load} from '../website/js/app.js';

const dom = new JSDOM(`<section id="content"></section><article id="content-article"></article>`, {url: 'http://localhost'});
const {document} = dom.window;
const content = document.getElementById('content');
const article = document.getElementById('content-article');

setBusy(content, true);
assert.equal(content.getAttribute('aria-busy'), 'true');

setBusy(content, false);
assert.equal(content.getAttribute('aria-busy'), 'false');

global.fetch = async () => ({ ok: true, text: async () => '<h2>Test</h2>' });
await load('content/README.html', content, article);
assert.equal(content.getAttribute('aria-busy'), 'false');
assert.match(article.innerHTML, /Test/);

console.log('app.test.js passed');
