document.addEventListener('DOMContentLoaded', () => {
  // Main tabs (MCP / Download)
  const tabs = document.querySelectorAll('.gs-tab');
  const panels = document.querySelectorAll('.gs-panel');
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      tabs.forEach(t => t.classList.remove('gs-tab--active'));
      panels.forEach(p => p.classList.remove('gs-panel--active'));
      tab.classList.add('gs-tab--active');
      document.getElementById('panel-' + tab.dataset.tab).classList.add('gs-panel--active');
    });
  });

  // Client tabs
  const ctabs = document.querySelectorAll('.gs-ctab');
  const cpanels = document.querySelectorAll('.gs-cpanel');
  ctabs.forEach(ctab => {
    ctab.addEventListener('click', () => {
      ctabs.forEach(t => t.classList.remove('gs-ctab--active'));
      cpanels.forEach(p => p.classList.remove('gs-cpanel--active'));
      ctab.classList.add('gs-ctab--active');
      document.getElementById('cpanel-' + ctab.dataset.ctab).classList.add('gs-cpanel--active');
    });
  });
});
