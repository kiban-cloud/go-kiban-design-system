// Shell «kiban cloud flamboyant» (layout/flamboyant.go): botón de panel, menús
// de la topbar y el panel de Pendientes. Sin estado guardado: el menú arranca
// contraído en cada página y elegir una opción lo vuelve a contraer.
(function () {
	var root = document.querySelector('[data-fs-root]');
	if (!root) return;
	var mobile = window.matchMedia('(max-width: 767px)');
	var toggle = root.querySelector('[data-fs-toggle]');
	// Los textos vienen traducidos del servidor (flamboyant.templ, data-fs-i18n-*).
	function t(key) { return root.getAttribute('data-fs-i18n-' + key) || ''; }

	function railOpen() {
		return mobile.matches ? root.hasAttribute('data-fs-mobile-open') : root.hasAttribute('data-fs-open');
	}
	function syncToggle() {
		if (!toggle) return;
		var open = railOpen();
		var label = open ? t('collapse') : t('expand');
		toggle.setAttribute('aria-expanded', open ? 'true' : 'false');
		toggle.setAttribute('aria-label', label);
		toggle.setAttribute('data-tooltip', label);
	}
	function setRail(open) {
		root.toggleAttribute(mobile.matches ? 'data-fs-mobile-open' : 'data-fs-open', open);
		syncToggle();
		placeInbox();
	}

	// ---------- menús de la topbar (espacio y avatar) ----------
	function closeMenus(except) {
		root.querySelectorAll('[data-fs-menu]').forEach(function (m) {
			if (m === except) return;
			var p = m.querySelector('[data-fs-menu-panel]'), t = m.querySelector('[data-fs-menu-trigger]');
			if (p) p.classList.add('hidden');
			if (t) t.setAttribute('aria-expanded', 'false');
		});
	}
	function toggleMenu(menu) {
		var p = menu.querySelector('[data-fs-menu-panel]'), t = menu.querySelector('[data-fs-menu-trigger]');
		closeMenus(menu);
		var open = p.classList.toggle('hidden') === false;
		t.setAttribute('aria-expanded', open ? 'true' : 'false');
		var q = open && p.querySelector('input[type="search"]');
		if (q) setTimeout(function () { q.focus(); q.select(); }, 20);
	}

	// ---------- Pendientes ----------
	var inbox = root.querySelector('[data-fs-inbox]');
	var inboxBack = root.querySelector('[data-fs-inbox-back]');
	var inboxBtn = root.querySelector('[data-fs-inbox-toggle]');
	var filter = 'all';
	function inboxOpen() { return inbox && !inbox.hidden; }
	function placeInbox() {
		if (!inbox) return;
		// Pegado al menú: 64 px contraído, 14rem expandido; en celular, a todo lo ancho.
		var left = mobile.matches ? 0 : root.hasAttribute('data-fs-open') ? 224 : 64;
		inbox.style.left = left + 'px';
		inboxBack.style.left = left + 'px';
	}
	function setInbox(open) {
		if (!inbox) return;
		inbox.hidden = !open;
		inboxBack.hidden = !open;
		if (inboxBtn) inboxBtn.setAttribute('aria-expanded', open ? 'true' : 'false');
		if (open) {
			if (mobile.matches) root.removeAttribute('data-fs-mobile-open');
			syncToggle();
			placeInbox();
		}
	}
	function items() { return Array.prototype.slice.call(inbox.querySelectorAll('.ib-item')); }
	function chipsFor(list) {
		var mods = [];
		list.forEach(function (it) {
			var m = it.getAttribute('data-ib-mod');
			if (m && mods.indexOf(m) === -1) mods.push(m);
		});
		var unread = list.filter(function (it) { return it.hasAttribute('data-ib-unread'); }).length;
		var chips = [['all', t('all'), list.length], ['unread', t('unread'), unread]].concat(mods.map(function (m) {
			return [m, m, list.filter(function (it) { return it.getAttribute('data-ib-mod') === m; }).length];
		}));
		var box = inbox.querySelector('[data-ib-chips]');
		box.textContent = '';
		chips.forEach(function (c) {
			var b = document.createElement('button');
			b.type = 'button';
			b.className = 'ib-chip' + (filter === c[0] ? ' is-on' : '');
			b.setAttribute('data-ib-filter', c[0]);
			b.setAttribute('aria-pressed', filter === c[0] ? 'true' : 'false');
			b.appendChild(document.createTextNode(c[1]));
			var n = document.createElement('span');
			n.textContent = c[2];
			b.appendChild(n);
			box.appendChild(b);
		});
		var sub = inbox.querySelector('[data-ib-sub]');
		if (sub) sub.textContent = unread ? t('unread-count').replace('{n}', unread) : t('up-to-date');
	}
	function applyFilter() {
		var q = (inbox.querySelector('input[data-ib-q]').value || '').trim().toLowerCase();
		var list = items(), shown = 0;
		list.forEach(function (it) {
			var byChip = filter === 'all' || (filter === 'unread' ? it.hasAttribute('data-ib-unread') : it.getAttribute('data-ib-mod') === filter);
			var text = (it.getAttribute('data-ib-q') || it.textContent).toLowerCase();
			var hit = byChip && (!q || text.indexOf(q) !== -1);
			it.hidden = !hit;
			if (hit) shown++;
		});
		inbox.querySelectorAll('[data-ib-grp]').forEach(function (g) { g.hidden = !g.querySelector('.ib-item:not([hidden])'); });
		var noq = inbox.querySelector('[data-ib-noq]');
		if (noq) noq.hidden = !(list.length && shown === 0);
		inbox.querySelectorAll('[data-ib-filter]').forEach(function (b) {
			var on = b.getAttribute('data-ib-filter') === filter;
			b.classList.toggle('is-on', on);
			b.setAttribute('aria-pressed', on ? 'true' : 'false');
		});
	}
	document.addEventListener('htmx:afterSwap', function (e) {
		if (!inbox || !inbox.contains(e.target)) return;
		chipsFor(items());
		applyFilter();
	});

	// ---------- eventos ----------
	document.addEventListener('click', function (e) {
		if (!e.target.closest) return;
		var menu = e.target.closest('[data-fs-menu]');
		if (!menu) closeMenus();
		if (e.target.closest('[data-fs-toggle]')) {
			e.preventDefault();
			setRail(!railOpen());
			return;
		}
		var trigger = e.target.closest('[data-fs-menu-trigger]');
		if (trigger && menu) { toggleMenu(menu); return; }
		if (e.target.closest('[data-fs-inbox-toggle]')) { setInbox(!inboxOpen()); return; }
		if (e.target.closest('[data-fs-inbox-close]')) { setInbox(false); return; }
		var chip = e.target.closest('[data-ib-filter]');
		if (chip) { filter = chip.getAttribute('data-ib-filter'); applyFilter(); return; }
		if (e.target.closest('[data-fs-close-mobile]')) { root.removeAttribute('data-fs-mobile-open'); syncToggle(); return; }
		// Elegir una opción del menú lo vuelve a contraer.
		if (e.target.closest('.fs-nav a')) {
			root.removeAttribute('data-fs-open');
			root.removeAttribute('data-fs-mobile-open');
			syncToggle();
		}
	});
	document.addEventListener('input', function (e) {
		var t = e.target;
		if (!t.matches) return;
		if (t.matches('input[data-ib-q]')) { applyFilter(); return; }
		if (t.matches('input[data-fs-space-q]')) {
			var q = t.value.trim().toLowerCase(), any = false;
			root.querySelectorAll('[data-fs-space]').forEach(function (s) {
				var hit = !q || s.getAttribute('data-fs-space').indexOf(q) !== -1;
				s.hidden = !hit;
				if (hit) any = true;
			});
			var none = root.querySelector('[data-fs-space-none]');
			if (none) none.hidden = any;
		}
	});
	document.addEventListener('keydown', function (e) {
		if (e.key !== 'Escape') return;
		closeMenus();
		if (inboxOpen()) { setInbox(false); return; }
		if (root.hasAttribute('data-fs-mobile-open')) { root.removeAttribute('data-fs-mobile-open'); syncToggle(); }
	});
	var onChange = function () { if (!mobile.matches) root.removeAttribute('data-fs-mobile-open'); syncToggle(); placeInbox(); };
	if (mobile.addEventListener) mobile.addEventListener('change', onChange); else if (mobile.addListener) mobile.addListener(onChange);
	syncToggle();
})();
