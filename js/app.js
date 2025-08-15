/*! (c) 2025 Asymmetric Effort, LLC. MIT License. */

(function () {
  const content = document.getElementById('content');
  const article = document.getElementById('content-article');
  const list = document.getElementById('nav-list');

  /**
   * Toggle the aria-busy state for the content pane.
   * @param {boolean} busy - True when content is loading.
   */
  function setBusy(busy) {
    content.setAttribute('aria-busy', String(busy));
  }

  /**
   * Load HTML from the given URL into the content pane.
   * @param {string} url - Relative URL of the document to fetch.
   */
  async function load(url) {
    setBusy(true);
    try {
      const res = await fetch(url, { credentials: 'omit' });
      if (!res.ok) throw new Error('HTTP ' + res.status);
      const html = await res.text();
      article.innerHTML = html;
      content.scrollTop = 0;
    } finally {
      setBusy(false);
    }
  }

  list?.addEventListener('click', (e) => {
    const a = e.target.closest('a.nav-link');
    if (!a) return;
    const url = a.dataset.src || a.getAttribute('href');
    if (!url || url.startsWith('http')) return;
    e.preventDefault();
    history.replaceState(null, '', '#' + encodeURIComponent(url));
    load(url);
  });

  window.addEventListener('DOMContentLoaded', () => {
    const hash = decodeURIComponent(location.hash.replace(/^#/, ''));
    if (hash) {
      load(hash);
    }
  });
})();
