package metricsui

const styles = `
:root {
  color-scheme: light dark;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  background: Canvas;
  color: CanvasText;
  --line: color-mix(in srgb, CanvasText 18%, Canvas);
  --line-strong: color-mix(in srgb, CanvasText 34%, Canvas);
  --muted: color-mix(in srgb, CanvasText 64%, Canvas);
  --wash: color-mix(in srgb, LinkText 7%, Canvas);
  --waiting: #8b5e00;
  --running: #0066a6;
  --completed: #08783f;
  --skipped: #625f58;
  --danger: #a23d00;
  --daily: #4f55a5;
  --interval: #7545a0;
  --entity: #08783f;
  --event: #9a5700;
}
@media (prefers-color-scheme: dark) {
  :root {
    --waiting: #e3b341;
    --running: #65b5e8;
    --completed: #66d296;
    --skipped: #c0bdb5;
    --danger: #f09670;
    --daily: #b5baff;
    --interval: #d0a7f2;
    --entity: #66d296;
    --event: #e3b341;
  }
}
* { box-sizing: border-box; }
body { margin: 0; line-height: 1.35; background: Canvas; color: CanvasText; }
[hidden], #metrics-sse { display: none !important; }
a { color: LinkText; text-underline-offset: .15em; }
:focus-visible { outline: 2px solid LinkText; outline-offset: 2px; }
.topbar { position: sticky; top: 0; z-index: 2; display: flex; justify-content: space-between; gap: 1rem; align-items: center; padding: .45rem max(1rem, calc((100vw - 1300px) / 2)); border-top: 2px solid LinkText; border-bottom: 1px solid var(--line-strong); background: Canvas; }
.topbar nav { display: flex; gap: 0; border: 1px solid var(--line-strong); }
.topbar a { padding: .3rem .55rem; border-right: 1px solid var(--line-strong); border-radius: 0; text-decoration: none; }
.topbar a:last-child { border-right: 0; }
.topbar a.active { color: Canvas; background: LinkText; }
.product { display: inline-flex; align-items: center; color: var(--muted); font-size: .72rem; letter-spacing: .04em; text-transform: uppercase; }
.product::before { margin-right: .5rem; color: LinkText; content: '■ ■'; font-size: .55rem; letter-spacing: .12rem; }
.muted, small, .relative { color: var(--muted); }
main { max-width: 1300px; margin: 0 auto; padding: .8rem 1rem 1.25rem; }
h1 { margin: 0; font-size: 1.25rem; }
.summary { display: grid; grid-template-columns: minmax(12rem, .75fr) minmax(30rem, 2fr) minmax(13rem, auto); align-items: stretch; margin-bottom: .65rem; padding: .6rem 0; border-top: 1px solid var(--line-strong); border-bottom: 1px solid var(--line-strong); }
.summary > div:first-child { align-self: center; padding-right: .8rem; }
.summary p, .catalog-heading p { margin: .2rem 0 0; font-size: .82rem; }
.summary dl { display: grid; grid-template-columns: repeat(3, minmax(8rem, 1fr)); margin: 0; border-left: 1px solid var(--line); }
.summary dl > div { min-width: 0; padding: .15rem .7rem; border-right: 1px solid var(--line); }
.summary dt { color: var(--muted); font-size: .7rem; letter-spacing: .035em; text-transform: uppercase; }
.summary dd { margin: .08rem 0 0; overflow-wrap: anywhere; font-size: .82rem; font-weight: 700; }
.summary-actions { display: grid; gap: .35rem; align-content: center; padding-left: .7rem; }
.summary-controls { display: flex; flex-wrap: wrap; gap: .28rem; align-content: start; }
.summary-controls a, .scope a { padding: .22rem .38rem; border: 1px solid var(--line-strong); border-radius: 0; color: CanvasText; font-size: .75rem; text-decoration: none; white-space: nowrap; }
.summary-controls a strong { color: LinkText; }
.summary-controls a[aria-current="page"], .scope a[aria-current="page"] { color: CanvasText; background: color-mix(in srgb, LinkText 14%, Canvas); font-weight: 700; }
.warning, .error { margin: 0 0 .65rem; padding: .5rem .65rem; border: 1px solid var(--line-strong); border-left: .3rem solid #b56a00; background: color-mix(in srgb, #e9ae48 12%, Canvas); }
.error { border-left-color: #b02020; background: color-mix(in srgb, #b02020 8%, Canvas); }
.filters { display: grid; gap: .5rem; align-items: end; margin-bottom: .5rem; padding: .5rem; border: 1px solid var(--line-strong); border-radius: 0; background: var(--wash); }
.filters[data-filter-kind="activity"] { grid-template-columns: minmax(15rem, 2fr) repeat(3, minmax(8rem, 1fr)) auto; }
.filters[data-filter-kind="automations"] { grid-template-columns: minmax(15rem, 2fr) minmax(10rem, 1fr) auto; }
.filters label { display: grid; gap: .12rem; color: var(--muted); font-size: .7rem; letter-spacing: .035em; text-transform: uppercase; }
input, select, button { border-radius: 0; font: inherit; }
input, select { width: 100%; min-height: 2rem; max-width: 100%; padding: .25rem .35rem; border: 1px solid var(--line-strong); background: Canvas; color: CanvasText; font-size: .78rem; }
button, .load-older { min-height: 2rem; padding: .3rem .55rem; border: 1px solid var(--line-strong); border-radius: 0; background: Canvas; color: CanvasText; cursor: pointer; }
button:hover, .load-older:hover { color: Canvas; background: LinkText; border-color: LinkText; }
.filter-note { grid-column: 1 / -1; color: var(--muted); font-size: .78rem; }
.scope { display: flex; gap: 0; width: fit-content; margin-bottom: .5rem; border: 1px solid var(--line-strong); }
.scope a { border: 0; border-right: 1px solid var(--line-strong); }
.scope a:last-child { border-right: 0; }
.automation-counts { margin-bottom: .55rem; }
.new-updates { position: sticky; top: 3.2rem; z-index: 1; width: fit-content; margin: .4rem auto; }
.feed { --activity-columns: 1.1rem 11rem minmax(11rem, 1fr) minmax(13rem, 1.3fr) 10rem; border: 1px solid var(--line-strong); }
.feed-header, .evaluation > summary, .system-event { display: grid; grid-template-columns: var(--activity-columns); gap: .5rem; align-items: center; }
.feed-header { padding: .32rem .4rem; border-bottom: 1px solid var(--line-strong); color: var(--muted); background: var(--wash); font-size: .68rem; font-weight: 700; letter-spacing: .055em; text-transform: uppercase; }
.feed-header span:last-child { text-align: right; }
.evaluation, .system-event { border-bottom: 1px solid var(--line); }
.evaluation:last-child, .system-event:last-child { border-bottom: 0; }
.evaluation > summary { grid-template-areas: 'arrow time automation trigger result'; min-height: 3.15rem; padding: .3rem .4rem; border-left: 3px solid transparent; cursor: pointer; list-style: none; font-size: .84rem; }
.evaluation > summary:hover, .evaluation[open] > summary { background: var(--wash); }
.evaluation[open] > summary { border-left-color: LinkText; }
.evaluation > summary::-webkit-details-marker { display: none; }
.evaluation > summary::before { grid-area: arrow; display: grid; place-items: center; width: .9rem; height: .9rem; border: 1px solid var(--line-strong); content: '›'; font-size: .9rem; line-height: 1; transform: rotate(0deg); }
.evaluation[open] > summary::before { transform: rotate(90deg); }
.evaluation-main { grid-area: automation; display: grid; min-width: 0; overflow-wrap: anywhere; }
.evaluation-main small { margin-top: .05rem; font-size: .7rem; }
.evaluation .time { grid-area: time; display: grid; min-width: 0; font-size: .76rem; font-variant-numeric: tabular-nums; }
.evaluation .time time, .evaluation .time .relative { white-space: nowrap; }
.evaluation .time .relative { font-size: .7rem; }
.trigger { grid-area: trigger; min-width: 0; overflow-wrap: anywhere; font-size: .78rem; }
.evaluation-meta { grid-area: result; display: flex; flex-wrap: wrap; justify-content: flex-end; align-items: center; gap: .25rem .4rem; min-width: 0; }
.evaluation-meta > .badge:first-child { min-width: 5.5rem; justify-content: center; }
.badge, .type-chip { display: inline-flex; align-items: center; gap: .28rem; padding: .1rem .32rem; border: 1px solid currentColor; border-radius: 0; background: color-mix(in srgb, currentColor 8%, Canvas); font-size: .68rem; font-weight: 700; letter-spacing: .025em; text-transform: uppercase; white-space: nowrap; }
.badge::before { width: .3rem; height: .3rem; flex: 0 0 auto; background: currentColor; content: ''; }
.waiting { color: var(--waiting); }
.running { color: var(--running); }
.completed { color: var(--completed); }
.skipped { color: var(--skipped); }
.canceled, .incomplete, .warning-badge, .trace-error { color: var(--danger); }
.tooltip { border-bottom: 1px dotted currentColor; color: var(--danger); cursor: help; }
.duration { min-width: 3.5rem; color: var(--muted); font-size: .74rem; font-variant-numeric: tabular-nums; text-align: right; white-space: nowrap; }
.trace { margin: .15rem 0 .35rem 2.6rem; padding: .15rem .6rem .15rem 1rem; border-left: 2px solid var(--line-strong); list-style: none; }
.trace li { display: grid; grid-template-columns: minmax(9rem, auto) 1fr auto; gap: .45rem; align-items: start; padding: .25rem 0; font-size: .78rem; }
.trace li + li { border-top: 1px dotted var(--line); }
.trace-loading, .trace-load-error { color: var(--muted); }
.trace-load-error { color: var(--danger); }
.raw { grid-column: 1 / -1; margin-top: .15rem; }
.raw summary { width: fit-content; color: var(--muted); cursor: pointer; font-size: .76rem; }
pre { max-width: 100%; overflow: auto; margin: .35rem 0 0; padding: .5rem; border: 1px solid var(--line); border-radius: 0; background: var(--wash); font-size: .75rem; white-space: pre-wrap; overflow-wrap: anywhere; }
.system-event { grid-template-areas: 'marker time automation trigger result'; min-height: 2.8rem; padding: .3rem .4rem; border-left: 3px solid var(--line); background: var(--wash); font-size: .84rem; }
.system-event::before { grid-area: marker; width: .38rem; height: .38rem; justify-self: center; background: LinkText; content: ''; }
.system-event > time { grid-area: time; min-width: 0; font-size: .76rem; font-variant-numeric: tabular-nums; }
.system-event > strong { grid-area: automation; font-size: .75rem; text-transform: uppercase; }
.system-event > span { grid-area: trigger; min-width: 0; overflow-wrap: anywhere; font-size: .78rem; }
.system-event .raw { grid-area: result; justify-self: end; margin-top: 0; }
.system-event .raw[open] { grid-column: 2 / -1; grid-row: 2; justify-self: stretch; }
.load-older, .end, .empty { display: block; width: fit-content; margin: .75rem auto; }
.end { color: var(--muted); font-size: .75rem; }
.empty { padding: .6rem; border: 1px dashed var(--line-strong); }
.catalog-heading { margin-bottom: .55rem; padding-bottom: .5rem; border-bottom: 1px solid var(--line-strong); }
.catalog-table-wrap { overflow-x: auto; border: 1px solid var(--line-strong); }
.catalog-table-wrap:focus-visible { outline-offset: 2px; }
.sr-only, .scroll-hint { position: absolute; width: 1px; height: 1px; overflow: hidden; margin: -1px; padding: 0; border: 0; clip: rect(0 0 0 0); white-space: nowrap; }
.catalog { width: 100%; min-width: 850px; border-collapse: collapse; }
.catalog th, .catalog td { padding: .38rem .45rem; border-bottom: 1px solid var(--line); text-align: left; vertical-align: top; }
.catalog th { color: var(--muted); background: var(--wash); font-size: .7rem; letter-spacing: .035em; text-transform: uppercase; }
.catalog tbody tr:last-child td { border-bottom: 0; }
.catalog tbody tr:hover { background: var(--wash); }
.catalog th.automation-key, .catalog td.automation-key { width: 1%; }
.catalog th.automation-key, .catalog td.automation-key > details > summary { white-space: nowrap; }
.catalog .automation-time { width: 1%; white-space: nowrap; font-size: .74rem; font-variant-numeric: tabular-nums; text-align: right; }
.catalog td.automation-time .relative { display: block; font-size: .68rem; }
.catalog details > summary { width: fit-content; cursor: pointer; }
.type-chip[data-type="daily"] { color: var(--daily); }
.type-chip[data-type="interval"] { color: var(--interval); }
.type-chip[data-type="entity"] { color: var(--entity); }
.type-chip[data-type="event"] { color: var(--event); }
.automation-details { width: 100%; min-width: 0; margin-top: .4rem; padding: .2rem 0 .2rem .65rem; border-left: 2px solid var(--line-strong); white-space: normal; }
.automation-details p { margin: .35rem 0; }
@media (max-width: 1000px) {
  .summary { grid-template-columns: minmax(12rem, .7fr) minmax(28rem, 2fr); }
  .summary-actions { grid-column: 1 / -1; grid-template-columns: repeat(2, minmax(0, 1fr)); margin-top: .55rem; padding: .55rem 0 0; border-top: 1px solid var(--line); }
  .filters[data-filter-kind="activity"] { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .filters[data-filter-kind="activity"] .search-control, .filter-note { grid-column: 1 / -1; }
}
@media (max-width: 850px) {
  .topbar { padding: .4rem .75rem; }
  main { padding: .65rem .75rem 1rem; }
  .summary { grid-template-columns: 1fr; }
  .summary > div:first-child { padding: 0 0 .5rem; }
  .summary dl { grid-template-columns: repeat(2, minmax(0, 1fr)); border-top: 1px solid var(--line); border-left: 0; }
  .summary dl > div { padding: .4rem .55rem; border-bottom: 1px solid var(--line); }
  .summary-actions { grid-column: auto; }
  .feed-header { display: none; }
  .evaluation > summary { grid-template-columns: 1.1rem minmax(0, 1fr); grid-template-areas: 'arrow time' 'arrow automation' 'arrow trigger' 'arrow result'; gap: .12rem .5rem; padding: .4rem; }
  .evaluation > summary::before { align-self: start; margin-top: .1rem; }
  .evaluation-meta { justify-content: flex-start; }
  .trace { margin-left: 2.4rem; }
  .trace li { grid-template-columns: 1fr; gap: .16rem; }
  .system-event { grid-template-columns: 1.1rem minmax(0, 1fr); grid-template-areas: 'marker time' 'marker automation' 'marker trigger' 'marker result'; gap: .12rem .5rem; }
  .system-event::before { align-self: start; margin-top: .35rem; }
  .system-event .raw { justify-self: start; }
  .system-event .raw[open] { grid-column: 2; grid-row: 5; }
  .scroll-hint { position: static; width: auto; height: auto; overflow: visible; margin: 0; padding: .3rem .45rem; border-bottom: 1px solid var(--line); clip: auto; color: var(--muted); background: var(--wash); font-size: .68rem; white-space: normal; }
}
@media (max-width: 620px) {
  .product { font-size: 0; }
  .product::before { margin-right: 0; }
  .summary dl, .summary-actions { grid-template-columns: 1fr; }
  .summary dl > div { border-right: 0; }
  .filters[data-filter-kind="activity"], .filters[data-filter-kind="automations"] { grid-template-columns: 1fr; }
  .filters[data-filter-kind="activity"] .search-control, .filter-note { grid-column: auto; }
  .filters button { width: 100%; }
  .scope { width: 100%; }
  .scope a { flex: 1; text-align: center; }
}
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after { scroll-behavior: auto !important; transition: none !important; }
}
@media (forced-colors: active) {
  .badge, .type-chip, .topbar nav, .scope { border-color: CanvasText; }
  .badge::before, .system-event::before { background: CanvasText; }
}
`

const clientScript = `(function () {
  var pending = 0;
  function localize(root) {
    (root || document).querySelectorAll('time[data-local-time]').forEach(function (node) {
      var date = new Date(node.getAttribute('datetime'));
      if (!isNaN(date.getTime())) {
        var shortOptions = {month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit'};
        if (date.getFullYear() !== new Date().getFullYear()) shortOptions.year = 'numeric';
        if (node.hasAttribute('data-show-seconds')) shortOptions.second = '2-digit';
        var fullOptions = {weekday: 'long', year: 'numeric', month: 'long', day: 'numeric', hour: 'numeric', minute: '2-digit', second: '2-digit', timeZoneName: 'long'};
        var full = date.toLocaleString(undefined, fullOptions);
        node.textContent = date.toLocaleString(undefined, shortOptions);
        node.title = full + ' · ' + node.getAttribute('datetime');
        node.setAttribute('aria-label', full);
      }
    });
    (root || document).querySelectorAll('[data-relative]').forEach(function (node) {
      var then = Number(node.getAttribute('data-relative')), seconds = Math.max(0, Math.floor((Date.now() - then) / 1000));
      var text = seconds < 10 ? 'now' : seconds < 60 ? seconds + 's ago' : seconds < 3600 ? Math.floor(seconds / 60) + 'm ago' : seconds < 86400 ? Math.floor(seconds / 3600) + 'h ago' : Math.floor(seconds / 86400) + 'd ago';
      node.textContent = node.closest('.time, .automation-time') ? text : ' · ' + text;
    });
	(root || document).querySelectorAll('[data-uptime]').forEach(function (node) {
	  var seconds = Math.max(0, Math.floor((Date.now() - Number(node.getAttribute('data-uptime'))) / 1000));
	  var text = seconds < 60 ? seconds + 's uptime' : seconds < 3600 ? Math.floor(seconds / 60) + 'm uptime' : seconds < 86400 ? Math.floor(seconds / 3600) + 'h uptime' : Math.floor(seconds / 86400) + 'd uptime';
	  node.textContent = ' · ' + text;
	});
  }
  function syncExpanded() {
    var url = new URL(window.location.href);
    url.searchParams.delete('expand');
    document.querySelectorAll('details[data-evaluation][open], details[data-automation][open]').forEach(function (node) { url.searchParams.append('expand', node.getAttribute('data-evaluation') || node.getAttribute('data-automation')); });
    history.replaceState(null, '', url.pathname + (url.search ? url.search : ''));
  }
  function updateBanner() {
    document.querySelectorAll('[data-updates-banner]').forEach(function (banner) { banner.hidden = pending === 0; var count = banner.querySelector('[data-updates-count]'); if (count) count.textContent = pending; });
  }
  function refreshActivity(newest) {
    var target = document.getElementById('activity-content');
    if (!target || !window.htmx) return;
	var url = new URL(window.location.href);
	if (newest) { url.searchParams.delete('before'); history.replaceState(null, '', url.pathname + url.search); }
    pending = 0; updateBanner();
    htmx.ajax('GET', '/activity/fragment' + url.search, {target: '#activity-content', swap: 'outerHTML'});
  }
  function atNewestUninterrupted() { return !new URL(window.location.href).searchParams.has('before') && window.scrollY < 24 && document.querySelectorAll('details[data-evaluation][open]').length === 0; }
  document.addEventListener('DOMContentLoaded', function () { localize(document); updateBanner(); window.setInterval(function () { localize(document); }, 1000); });
  document.body.addEventListener('toggle', function (event) { if (event.target.matches('details[data-evaluation], details[data-automation]')) syncExpanded(); }, true);
  document.body.addEventListener('click', function (event) { var button = event.target.closest('[data-show-new]'); if (!button) return; event.preventDefault(); var reducedMotion = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches; window.scrollTo({top: 0, behavior: reducedMotion ? 'auto' : 'smooth'}); refreshActivity(true); });
  document.body.addEventListener('htmx:configRequest', function (event) {
    var path = event.detail.path.split('?')[0];
    if (path === '/activity/fragment' || path === '/activity/older') {
      event.detail.parameters.delete('expand');
      new URL(window.location.href).searchParams.getAll('expand').forEach(function (value) { event.detail.parameters.append('expand', value); });
    }
  });
  document.body.addEventListener('htmx:beforeSwap', function (event) {
    var path = event.detail.requestConfig && event.detail.requestConfig.path;
    if (path && path.indexOf('/activity/trace?') === 0 && event.detail.xhr.status >= 400) { event.detail.shouldSwap = true; event.detail.isError = false; }
  });
  document.body.addEventListener('htmx:afterSwap', function (event) { localize(event.target); updateBanner(); });
  document.body.addEventListener('htmx:sseMessage', function (event) {
    if (!event.detail || event.detail.type !== 'change' || !document.getElementById('activity-content')) return;
    if (atNewestUninterrupted()) refreshActivity(); else { pending++; updateBanner(); }
  });
}());`
