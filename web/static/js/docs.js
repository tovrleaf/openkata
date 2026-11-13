// Heading anchors — a click copies the section's URL and puts it in the
// address bar. Without the clipboard the link is followed as usual.
(function() {
  if (!navigator.clipboard || !window.history) return;

  document.addEventListener('click', function(e) {
    var link = e.target.closest && e.target.closest('.heading-anchor');
    if (!link) return;
    e.preventDefault();
    history.replaceState(null, '', link.getAttribute('href'));
    navigator.clipboard.writeText(link.href).then(function() {
      link.classList.add('is-copied');
      setTimeout(function() { link.classList.remove('is-copied'); }, 1500);
    }, function() {
      location.hash = link.getAttribute('href');
    });
  });
})();

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
