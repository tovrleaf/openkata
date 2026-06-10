(function() {
	var main = document.querySelector('main[data-artifact-name]');
	if (!main) return;

	var artifactPath = main.dataset.artifactPath;
	var artifactName = main.dataset.artifactName;

	// Version select navigation
	var versionSelect = document.querySelector('.version-select');
	if (versionSelect) {
		versionSelect.addEventListener('change', function() {
			window.location.href = '/' + artifactPath + '/' + artifactName + '/' + this.value;
		});
	}

	// Tab switching
	document.querySelectorAll('.tabs .tab').forEach(function(btn) {
		btn.addEventListener('click', function() {
			switchTab(btn.dataset.tab);
		});
	});

	function switchTab(name) {
		document.querySelectorAll('.tabs .tab').forEach(function(b) { b.classList.remove('active'); });
		document.querySelectorAll('.tab-panel').forEach(function(p) { p.classList.remove('active'); });
		var activeBtn = document.querySelector('.tabs .tab[data-tab="'+name+'"]');
		if (activeBtn) activeBtn.classList.add('active');
		var panel = document.getElementById('panel-' + name);
		if (panel) panel.classList.add('active');
		if (!switchTab._initializing) {
			history.pushState(null, '', '#' + name);
		}
	}

	var hash = location.hash.replace('#', '');
	if (hash) {
		var parts = hash.split(':');
		var tabName = parts[0];
		var fileId = parts[1] || '';
		if (document.querySelector('.tabs .tab[data-tab="'+tabName+'"]')) {
			switchTab._initializing = true;
			switchTab(tabName);
			switchTab._initializing = false;
		}
		if (fileId) {
			var target = document.getElementById('file-' + fileId);
			if (target) {
				target.open = true;
				target.scrollIntoView({ behavior: 'smooth', block: 'start' });
			}
		}
	}

	window.addEventListener('popstate', function() {
		var hash = location.hash.replace('#', '');
		if (hash) {
			var parts = hash.split(':');
			var tabName = parts[0];
			if (document.querySelector('.tabs .tab[data-tab="'+tabName+'"]')) {
				switchTab._initializing = true;
				switchTab(tabName);
				switchTab._initializing = false;
			}
		}
	});

	// Relative link interception — resolve to Files tab
	document.querySelectorAll('.tab-panel').forEach(function(panel) {
		panel.addEventListener('click', function(e) {
			var link = e.target.closest('a');
			if (!link) return;
			var href = link.getAttribute('href');
			if (!href || href.startsWith('http') || href.startsWith('#') || href.startsWith('/')) return;
			e.preventDefault();
			highlightFile(href);
		});
	});

	function highlightFile(path) {
		path = path.replace(/^\.\//, '');
		var id = 'file-' + path.replace(/\//g, '-').replace(/\./g, '-');
		var target = document.getElementById(id);
		if (!target) {
			id = 'file-references-' + path.replace(/\//g, '-').replace(/\./g, '-');
			target = document.getElementById(id);
		}
		if (target) {
			switchTab('files');
			target.open = true;
			target.scrollIntoView({ behavior: 'smooth', block: 'start' });
		}
		// Highlight in tree
		var treeLink = document.querySelector('.file-tree-link[href="#' + id + '"]');
		if (treeLink) {
			document.querySelectorAll('.file-tree-link').forEach(function(l) { l.classList.remove('file-highlight'); });
			treeLink.classList.add('file-highlight');
		}
	}

	// File tree navigation - scroll to file block and open it
	var fileTree = document.getElementById('file-tree');
	if (fileTree) {
		fileTree.addEventListener('click', function(e) {
			var link = e.target.closest('.file-tree-link');
			if (!link) return;
			e.preventDefault();
			var target = document.querySelector(link.getAttribute('href'));
			if (target) {
				target.open = true;
				target.scrollIntoView({ behavior: 'smooth', block: 'start' });
				history.replaceState(null, '', '#files:' + target.id.replace('file-', ''));
			}
		});
	}

	// Preview/Code toggle within file blocks
	document.querySelectorAll('.file-mode-btn[data-view]').forEach(function(btn) {
		btn.addEventListener('click', function(e) {
			e.preventDefault();
			e.stopPropagation();
			var block = btn.closest('.file-block');
			block.querySelectorAll('.file-mode-btn[data-view]').forEach(function(b) { b.classList.remove('active'); });
			btn.classList.add('active');
			var preview = block.querySelector('.file-preview');
			var code = block.querySelector('.file-code');
			if (btn.dataset.view === 'preview' && preview) {
				preview.style.display = 'block';
				if (code) code.style.display = 'none';
			} else if (btn.dataset.view === 'code') {
				if (preview) preview.style.display = 'none';
				if (code) code.style.display = 'block';
			}
			block.open = true;
		});
	});

	// Activate Preview button when file block opens, clear on close
	document.querySelectorAll('.file-block').forEach(function(block) {
		block.addEventListener('toggle', function() {
			if (block.open) {
				var previewBtn = block.querySelector('.file-mode-btn[data-view="preview"]');
				if (previewBtn && !block.querySelector('.file-mode-btn.active')) {
					previewBtn.classList.add('active');
				}
				history.replaceState(null, '', '#files:' + block.id.replace('file-', ''));
			} else {
				block.querySelectorAll('.file-mode-btn').forEach(function(b) { b.classList.remove('active'); });
				history.replaceState(null, '', '#files');
			}
		});
	});
})();
