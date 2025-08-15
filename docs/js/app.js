/*
 * Copyright (c) 2025 Asymmetric Effort, LLC. MIT License.
 * Client-side logic for Docker Lint site.
 */
/**
 * Toggle the aria-busy state on an element.
 * @param {HTMLElement} el - element to update.
 * @param {boolean} busy - busy state.
 */
export function setBusy(el, busy){
  if(!el) return;
  el.setAttribute('aria-busy', String(busy));
}
/**
 * Fetch HTML from a URL and inject into the article element.
 * @param {string} url - relative URL to fetch.
 * @param {{content: HTMLElement, article: HTMLElement}} elements - DOM nodes.
 */
export async function load(url, elements){
  const {content, article} = elements;
  setBusy(content, true);
  try{
    const res = await fetch(url, {credentials: 'omit'});
    if(!res.ok) throw new Error('HTTP '+res.status);
    const html = await res.text();
    article.innerHTML = html;
    content.scrollTop = 0;
  }finally{
    setBusy(content, false);
  }
}
/**
 * Initialize navigator click handling and deep-linking.
 * @param {Document} doc - document object.
 * @param {Window} win - window object.
 */
export function init(doc = document, win = window){
  const content = doc.getElementById('content');
  const article = doc.getElementById('content-article');
  const list = doc.getElementById('nav-list');
  list?.addEventListener('click', (e)=>{
    const a = e.target.closest('a.nav-link');
    if(!a) return;
    const url = a.dataset.src || a.getAttribute('href');
    if(!url || url.startsWith('http')) return;
    e.preventDefault();
    win.history.replaceState(null, '', '#' + encodeURIComponent(url));
    load(url, {content, article});
  });
  win.addEventListener('DOMContentLoaded', ()=>{
    const hash = decodeURIComponent(win.location.hash.replace(/^#/, ''));
    if(hash){ load(hash, {content, article}); }
  });
}
if(typeof window !== 'undefined' && typeof document !== 'undefined'){
  init();
}
