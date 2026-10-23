// Docs sidebar — open on wide screens, collapsed on narrow ones.
// Without JavaScript the sidebar stays collapsed and every link works.
(function() {
  var sidebar = document.getElementById('docs-sidebar');
  if (!sidebar || !window.matchMedia) return;

  var wide = window.matchMedia('(min-width: 769px)');

  function sync() {
    sidebar.open = wide.matches;
  }

  sync();
  if (wide.addEventListener) {
    wide.addEventListener('change', sync);
  } else {
    wide.addListener(sync);
  }
})();
