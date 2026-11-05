// Fechas del shell «kiban cloud flamboyant» (layout/flamboyant.go): todas como
// «28 oct 2026 12:02», o «28 oct 2026» si no llevan hora. Mes en tres letras
// minúsculas, día sin cero, año completo y hora de 24 h.
//
// Las vistas de cada herramienta todavía mandan sus fechas ya formateadas
// (DD/MM/AAAA HH:MM, AAAA-MM-DD, o el <time data-ds-localtime> que localiza
// Base, que se rearma de su atributo datetime). Mientras cada una no formatee
// con este patrón, aquí se reescriben en el texto que pinta el shell, también
// lo que llega después por htmx o por JS.
// No toca inputs, textareas, código (<pre>, <code>, el editor JSON) ni nada
// marcado con data-fs-sin-fechas.
(function () {
	var root = document.querySelector('[data-fs-root]');
	if (!root || !window.MutationObserver) return;
	// Los meses del idioma del request, de data-fs-months (flamboyant.templ).
	var MESES = (root.getAttribute('data-fs-months') || '').split(',');
	if (MESES.length !== 12) MESES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic'];
	// 28/10/2026, 28/10/2026 12:02, 28/10/2026, 12:02:09 (el de toLocaleString)
	var DMY = /\b(\d{1,2})\/(\d{1,2})\/(\d{4})(?:,?\s+(\d{1,2}):(\d{2})(?::\d{2})?)?\b/g;
	// 2026-10-28, 2026-10-28 12:02, 2026-10-28T12:02:09
	var YMD = /\b(\d{4})-(\d{2})-(\d{2})(?:[ T](\d{2}):(\d{2})(?::\d{2})?)?\b/g;
	var SALTA = { SCRIPT: 1, STYLE: 1, TEXTAREA: 1, CODE: 1, PRE: 1, NOSCRIPT: 1 };

	function dos(n) { return n < 10 ? '0' + n : '' + n; }
	function fecha(dia, mes, anio, hora, min) {
		dia = +dia; mes = +mes;
		if (mes < 1 || mes > 12 || dia < 1 || dia > 31) return null;
		var s = dia + ' ' + MESES[mes - 1] + ' ' + anio;
		return hora ? s + ' ' + dos(+hora) + ':' + min : s;
	}
	function convertir(t) {
		return t
			.replace(DMY, function (todo, d, m, a, h, mi) { return fecha(d, m, a, h, mi) || todo; })
			.replace(YMD, function (todo, a, m, d, h, mi) { return fecha(d, m, a, h, mi) || todo; });
	}
	function salta(el) {
		for (; el && el !== document.body; el = el.parentElement) {
			if (SALTA[el.tagName] || el.isContentEditable || el.hasAttribute('data-fs-sin-fechas') || (el.classList && el.classList.contains('ace_editor'))) return true;
		}
		return false;
	}
	// El <time data-ds-localtime> de Base se localiza con el idioma del
	// navegador (en inglés sale 04/14/2026): ese se arma del instante, no del texto.
	function local(el) {
		var d = new Date(el.getAttribute('datetime'));
		if (isNaN(d.getTime())) return;
		var nuevo = window.kibanFecha(d, true);
		if (el.textContent !== nuevo) el.textContent = nuevo;
	}
	function texto(n) {
		var v = n.nodeValue, p = n.parentElement;
		if (p && p.matches('time[data-ds-localtime][datetime]')) { local(p); return; }
		if (!v || !/\d[\/-]\d/.test(v) || salta(p)) return;
		var nuevo = convertir(v);
		if (nuevo !== v) n.nodeValue = nuevo;
	}
	function recorrer(nodo) {
		if (nodo.nodeType === 3) { texto(nodo); return; }
		if (nodo.nodeType !== 1 && nodo.nodeType !== 11) return;
		// Primero se juntan: rearmar un <time> reemplaza su nodo de texto y el
		// walker, parado en un nodo que ya no está en el árbol, cortaría ahí.
		var w = document.createTreeWalker(nodo, NodeFilter.SHOW_TEXT), nodos = [];
		for (var n = w.nextNode(); n; n = w.nextNode()) nodos.push(n);
		nodos.forEach(texto);
	}

	window.kibanFecha = function (d, conHora) {
		return fecha(d.getDate(), d.getMonth() + 1, d.getFullYear(), conHora ? d.getHours() + '' : '', dos(d.getMinutes()));
	};
	recorrer(root);
	new MutationObserver(function (cambios) {
		cambios.forEach(function (c) {
			if (c.type === 'characterData') texto(c.target);
			else c.addedNodes.forEach(recorrer);
		});
	}).observe(root, { subtree: true, childList: true, characterData: true });
})();
