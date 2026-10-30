// Runs web/static/js/search.js against a small fake page and prints the
// results it renders as JSON. The scenario comes from stdin:
//   { "searchJS": path, "mode": "overlay" | "catalog", "query": string,
//     "versions": object, "docs": array | null, "docsNetworkError": bool }
// A null docs means the docs index answers 404.
const fs = require('fs');
const vm = require('vm');

const cfg = JSON.parse(fs.readFileSync(0, 'utf8'));

function element(id) {
  const handlers = {};
  return {
    id: id,
    hidden: true,
    value: '',
    innerHTML: '',
    style: {},
    focus() {},
    setAttribute() {},
    querySelector() { return element('child'); },
    querySelectorAll() { return []; },
    addEventListener(type, fn) { (handlers[type] = handlers[type] || []).push(fn); },
    fire(type) { (handlers[type] || []).forEach((fn) => fn({ key: '' })); },
  };
}

const ids = ['search-overlay', 'search-open-btn', 'search-close-btn', 'search-input',
  'search-results', 'nav-menu-btn', 'nav-menu-panel'];
if (cfg.mode === 'catalog') ids.push('catalog-page', 'catalog-input', 'catalog-results');
const elements = {};
ids.forEach((id) => { elements[id] = element(id); });

function fetch(url) {
  if (url === '/static/versions.json') {
    return Promise.resolve({ ok: true, json: () => Promise.resolve(cfg.versions) });
  }
  if (url === '/docs/search-index.json') {
    if (cfg.docsNetworkError) return Promise.reject(new Error('network down'));
    if (cfg.docs === null) {
      return Promise.resolve({ ok: false, status: 404, json: () => Promise.reject(new Error('not json')) });
    }
    return Promise.resolve({ ok: true, json: () => Promise.resolve(cfg.docs) });
  }
  return Promise.reject(new Error('unexpected fetch ' + url));
}

const context = {
  document: {
    getElementById: (id) => elements[id] || null,
    addEventListener() {},
    body: { style: {} },
  },
  window: { location: { search: cfg.mode === 'catalog' ? '?q=' + encodeURIComponent(cfg.query) : '', href: '' } },
  history: { replaceState() {} },
  fetch: fetch,
  URLSearchParams: URLSearchParams,
  encodeURIComponent: encodeURIComponent,
};

const wait = () => new Promise((resolve) => setTimeout(resolve, 20));

(async () => {
  vm.runInNewContext(fs.readFileSync(cfg.searchJS, 'utf8'), context);
  if (cfg.mode === 'overlay') {
    elements['search-open-btn'].fire('click');
    await wait();
    elements['search-input'].value = cfg.query;
    elements['search-input'].fire('input');
    await wait();
    console.log(JSON.stringify({ html: elements['search-results'].innerHTML }));
  } else {
    await wait();
    console.log(JSON.stringify({ html: elements['catalog-results'].innerHTML }));
  }
})().catch((err) => {
  console.error(err);
  process.exit(1);
});
