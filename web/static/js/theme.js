// Theme switcher — single button cycles themes
(function() {
  var themes = [
    { id: 'japandi', color: '#fcfcfc' },
    { id: 'hokusai', color: '#d9d1ba' },
    { id: 'indigo', color: '#1a1a2e' }
  ];

  var saved = localStorage.getItem('theme');
  if (saved) document.documentElement.setAttribute('data-theme', saved);

  function getNextTheme() {
    var current = document.documentElement.getAttribute('data-theme') || 'light';
    var idx = themes.findIndex(function(t) { return t.id === current; });
    return themes[(idx + 1) % themes.length];
  }

  function updateButton() {
    var btn = document.querySelector('.theme-switcher__btn');
    if (btn) {
      var next = getNextTheme();
      btn.style.background = next.color;
      btn.setAttribute('aria-label', next.id + ' theme');
    }
  }

  window.cycleTheme = function() {
    var next = getNextTheme();
    document.documentElement.setAttribute('data-theme', next.id);
    localStorage.setItem('theme', next.id);
    updateButton();
  };

  document.addEventListener('DOMContentLoaded', updateButton);
})();
