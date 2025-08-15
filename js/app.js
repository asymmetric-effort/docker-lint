/*
 * Copyright (c) 2025 Asymmetric Effort, LLC.
 * Licensed under the MIT License.
 */
/**
 * Set aria-busy attribute on an element.
 * @param {HTMLElement} el - Element to toggle.
 * @param {boolean} state - Busy state.
 */
export function setBusy(el, state){
  el.setAttribute('aria-busy', String(state));
}

/**
 * Fetch HTML from a URL and inject into article element.
 * @param {string} url - Source URL.
 * @param {HTMLElement} content - Container element.
 * @param {HTMLElement} article - Article element to update.
 * @returns {Promise<void>} Resolves when content loaded.
 */
export async function load(url, content, article){
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
 * Initialize navigator click handling and deep-link support.
 */
export function init(){
  const content = document.getElementById('content');
  const article = document.getElementById('content-article');
  const list = document.getElementById('nav-list');

  list?.addEventListener('click', (e)=>{
    const a = e.target.closest('a.nav-link');
    if(!a) return;
    const url = a.dataset.src || a.getAttribute('href');
    if(!url || url.startsWith('http')) return;
    e.preventDefault();
    history.replaceState(null, '', '#' + encodeURIComponent(url));
    load(url, content, article);
  });

  const hash = decodeURIComponent(location.hash.replace(/^#/, ''));
  if(hash){ load(hash, content, article); }
}

if(typeof window !== 'undefined'){
  window.addEventListener('DOMContentLoaded', init);
}
