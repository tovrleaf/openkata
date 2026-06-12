(function() {
  var overlay = document.getElementById('search-overlay');
  var openBtn = document.getElementById('search-open-btn');
  var closeBtn = document.getElementById('search-close-btn');
  var backdrop = overlay ? overlay.querySelector('.search-overlay__backdrop') : null;
  var searchInput = document.getElementById('search-input');
  var searchResults = document.getElementById('search-results');

  var catalogData = null;
  var fetching = false;

  if (!overlay || !openBtn) return;

  // --- Overlay open/close ---

  function openSearch() {
    overlay.hidden = false;
    searchInput.value = '';
    searchInput.focus();
    document.body.style.overflow = 'hidden';
    if (!catalogData && !fetching) fetchData();
  }

  function closeSearch() {
    overlay.hidden = true;
    document.body.style.overflow = '';
    if (searchResults) searchResults.innerHTML = '';
  }

  openBtn.addEventListener('click', openSearch);
  closeBtn.addEventListener('click', closeSearch);
  backdrop.addEventListener('click', closeSearch);

  document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape' && !overlay.hidden) {
      closeSearch();
    }
  });

  // --- Data fetching ---

  function fetchData(cb) {
    if (catalogData) { if (cb) cb(); return; }
    fetching = true;
    fetch('/static/versions.json')
      .then(function(r) { return r.json(); })
      .then(function(data) {
        catalogData = parseData(data);
        fetching = false;
        if (cb) cb();
      })
      .catch(function() { fetching = false; });
  }

  function parseData(data) {
    var items = [];
    ['skills', 'rules', 'profiles'].forEach(function(type) {
      if (!data[type]) return;
      Object.keys(data[type]).forEach(function(name) {
        var info = data[type][name];
        if (info.version === '0.0.0') return;
        items.push({
          name: name,
          type: type.slice(0, -1),
          path: '/' + type + '/' + name + '/',
          description: info.description || '',
          tags: info.tags || ''
        });
      });
    });
    items.sort(function(a, b) { return a.name.localeCompare(b.name); });
    return items;
  }

  // --- Search logic ---

  function tokenize(query) {
    return query.trim().split(/\s+/).filter(function(t) { return t.length > 0; });
  }

  function isTagToken(token) {
    return token.indexOf(':') > 0 && token.indexOf(' ') === -1;
  }

  function scoreItem(item, tokens) {
    var total = 0;
    for (var i = 0; i < tokens.length; i++) {
      var token = tokens[i];
      var score = 0;
      if (isTagToken(token)) {
        var tags = item.tags.split(',').map(function(t) { return t.trim().toLowerCase(); });
        var found = false;
        for (var j = 0; j < tags.length; j++) {
          if (tags[j] === token.toLowerCase()) { found = true; break; }
        }
        if (!found) return 0;
        score = 50;
      } else {
        var lower = token.toLowerCase();
        var nameMatch = item.name.toLowerCase().indexOf(lower) >= 0;
        var tagMatch = item.tags.toLowerCase().indexOf(lower) >= 0;
        var descMatch = item.description.toLowerCase().indexOf(lower) >= 0;
        if (!nameMatch && !tagMatch && !descMatch) return 0;
        if (nameMatch) score += 100;
        if (tagMatch) score += 50;
        if (descMatch) score += 25;
      }
      total += score;
    }
    return total;
  }

  function search(query) {
    if (!catalogData) return [];
    var tokens = tokenize(query);
    if (tokens.length === 0) return catalogData.slice();
    var results = [];
    for (var i = 0; i < catalogData.length; i++) {
      var score = scoreItem(catalogData[i], tokens);
      if (score > 0) results.push({ item: catalogData[i], score: score });
    }
    results.sort(function(a, b) { return b.score - a.score; });
    return results.map(function(r) { return r.item; });
  }

  // --- Highlighting ---

  function escapeHtml(str) {
    return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function highlightText(text, tokens) {
    if (!tokens.length) return escapeHtml(text);
    var freeTokens = tokens.filter(function(t) { return !isTagToken(t); });
    if (!freeTokens.length) return escapeHtml(text);
    var escaped = escapeHtml(text);
    freeTokens.forEach(function(token) {
      var regex = new RegExp('(' + escapeRegex(escapeHtml(token)) + ')', 'gi');
      escaped = escaped.replace(regex, '<mark class="search-highlight">$1</mark>');
    });
    return escaped;
  }

  function highlightTag(tag, tokens) {
    var tagTokens = tokens.filter(function(t) { return isTagToken(t); });
    var isMatched = false;
    for (var i = 0; i < tagTokens.length; i++) {
      if (tag.toLowerCase() === tagTokens[i].toLowerCase()) {
        isMatched = true;
        break;
      }
    }
    var freeTokens = tokens.filter(function(t) { return !isTagToken(t); });
    var html = escapeHtml(tag);
    freeTokens.forEach(function(token) {
      var regex = new RegExp('(' + escapeRegex(escapeHtml(token)) + ')', 'gi');
      html = html.replace(regex, '<mark class="search-highlight">$1</mark>');
    });
    if (isMatched) {
      return '<mark class="search-highlight">' + html + '</mark>';
    }
    return html;
  }

  function escapeRegex(str) {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  }

  // --- Overlay rendering ---

  function renderOverlayResults(query) {
    if (!searchResults) return;
    var results = search(query);
    var tokens = tokenize(query);
    var max = 7;
    var html = '';

    if (results.length === 0 && query.trim().length > 0) {
      html = '<div class="search-overlay__empty">No results found.</div>';
    } else {
      var shown = results.slice(0, max);
      for (var i = 0; i < shown.length; i++) {
        var item = shown[i];
        html += '<a href="' + item.path + '" class="search-overlay__item">';
        html += '<span class="search-overlay__item-name">' + highlightText(item.name, tokens) + '</span>';
        html += '<span class="search-overlay__item-type">' + item.type + '</span>';
        html += '</a>';
      }
      if (results.length > max) {
        var q = encodeURIComponent(query);
        html += '<a href="/catalog/?q=' + q + '" class="search-overlay__view-all">View all ' + results.length + ' results</a>';
      }
    }
    searchResults.innerHTML = html;
  }

  // --- Overlay input handling ---

  searchInput.addEventListener('input', function() {
    if (!catalogData) {
      fetchData(function() { renderOverlayResults(searchInput.value); });
      return;
    }
    renderOverlayResults(searchInput.value);
  });

  searchInput.addEventListener('keydown', function(e) {
    if (e.key === 'Enter') {
      var q = searchInput.value.trim();
      window.location.href = '/catalog/?q=' + encodeURIComponent(q);
    }
  });

  // --- Catalog page logic ---

  var catalogPage = document.getElementById('catalog-page');
  var catalogInput = document.getElementById('catalog-input');
  var catalogResults = document.getElementById('catalog-results');

  if (catalogPage && catalogInput && catalogResults) {
    fetchData(function() {
      var params = new URLSearchParams(window.location.search);
      var q = params.get('q') || '';
      catalogInput.value = q;
      renderCatalogResults(q);
    });

    catalogInput.addEventListener('input', function() {
      var q = catalogInput.value;
      var params = new URLSearchParams(window.location.search);
      if (q) {
        params.set('q', q);
      } else {
        params.delete('q');
      }
      history.replaceState(null, '', '?' + params.toString());
      renderCatalogResults(q);
    });
  }

  function renderCatalogResults(query) {
    if (!catalogResults || !catalogData) return;
    var results = search(query);
    var tokens = tokenize(query);
    var html = '';

    if (results.length === 0 && query.trim().length > 0) {
      html = '<div class="catalog-empty">No results found.</div>';
    } else {
      for (var i = 0; i < results.length; i++) {
        var item = results[i];
        html += '<div class="catalog-item">';
        html += '<div class="catalog-item__header">';
        html += '<a href="' + item.path + '" class="catalog-item__name">' + highlightText(item.name, tokens) + '</a>';
        html += '<span class="catalog-item__type">' + item.type + '</span>';
        html += '</div>';
        if (item.description) {
          html += '<p class="catalog-item__desc">' + highlightText(item.description, tokens) + '</p>';
        }
        if (item.tags) {
          var tagList = item.tags.split(',').map(function(t) { return t.trim(); }).filter(function(t) { return t; });
          html += '<div class="catalog-item__tags">';
          for (var j = 0; j < tagList.length; j++) {
            html += '<a href="/catalog/?q=' + encodeURIComponent(tagList[j]) + '" class="badge badge-green">' + highlightTag(tagList[j], tokens) + '</a>';
          }
          html += '</div>';
        }
        html += '</div>';
      }
    }
    catalogResults.innerHTML = html;
  }
  // --- Mobile nav menu ---
  var menuBtn = document.getElementById('nav-menu-btn');
  var menuPanel = document.getElementById('nav-menu-panel');

  if (menuBtn && menuPanel) {
    menuBtn.addEventListener('click', function() {
      var open = !menuPanel.hidden;
      menuPanel.hidden = open;
      menuBtn.setAttribute('aria-expanded', String(!open));
    });

    menuPanel.querySelectorAll('.nav-menu-link').forEach(function(link) {
      link.addEventListener('click', function() {
        menuPanel.hidden = true;
        menuBtn.setAttribute('aria-expanded', 'false');
      });
    });
  }
})();
