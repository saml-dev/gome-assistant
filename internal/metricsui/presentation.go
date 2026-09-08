package metricsui

const styles = `
:root { color-scheme: light dark; font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; background: Canvas; color: CanvasText; }
* { box-sizing: border-box; }
body { margin: 0; line-height: 1.4; }
a { color: LinkText; }
.topbar { position: sticky; top: 0; z-index: 2; display: flex; justify-content: space-between; gap: 1rem; align-items: center; padding: .6rem max(1rem, calc((100vw - 1300px) / 2)); border-bottom: 1px solid color-mix(in srgb, CanvasText 24%, Canvas); background: Canvas; }
.topbar nav { display: flex; gap: .35rem; }
.topbar a { padding: .35rem .55rem; border-radius: .25rem; text-decoration: none; }
.topbar a.active { color: Canvas; background: LinkText; }
.product, .muted, small, .relative { color: color-mix(in srgb, CanvasText 58%, Canvas); }
main { max-width: 1300px; margin: 0 auto; padding: 1rem; }
h1 { margin: 0; font-size: 1.3rem; }
.summary, .catalog-heading { display: grid; gap: 1rem; margin-bottom: .9rem; }
.summary { grid-template-columns: minmax(12rem, .8fr) 2fr auto; align-items: start; }
.summary p { margin: .25rem 0 0; }
.summary dl { display: grid; grid-template-columns: repeat(3, minmax(8rem, 1fr)); gap: .45rem 1rem; margin: 0; }
.summary dt { font-size: .75rem; color: color-mix(in srgb, CanvasText 58%, Canvas); }
.summary dd { margin: 0; overflow-wrap: anywhere; }
.summary-controls { display: flex; flex-wrap: wrap; gap: .4rem; align-content: start; }
.summary-controls a, .scope a { border: 1px solid color-mix(in srgb, CanvasText 25%, Canvas); border-radius: .25rem; padding: .25rem .4rem; text-decoration: none; white-space: nowrap; }
.warning, .error { padding: .65rem .8rem; border-left: .3rem solid #b56a00; margin-bottom: .9rem; background: color-mix(in srgb, #e9ae48 18%, Canvas); }
.error { border-color: #b02020; }
.filters { display: flex; flex-wrap: wrap; gap: .65rem; align-items: end; padding: .65rem; margin-bottom: .55rem; border: 1px solid color-mix(in srgb, CanvasText 18%, Canvas); }
.filters label { display: grid; gap: .15rem; font-size: .85rem; }
input, select, button { font: inherit; }
input, select { max-width: 100%; padding: .28rem .35rem; }
button, .load-older { padding: .34rem .55rem; border: 1px solid color-mix(in srgb, CanvasText 35%, Canvas); border-radius: .25rem; background: Canvas; color: CanvasText; cursor: pointer; }
.filter-note { font-size: .85rem; }
.scope { display: flex; gap: .3rem; margin-bottom: .6rem; }
.summary-controls a[aria-current="page"], .scope a[aria-current="page"] { background: color-mix(in srgb, LinkText 18%, Canvas); font-weight: 700; }
.automation-counts { margin-bottom: .7rem; }
.new-updates { margin: .5rem 0; }
.feed { --activity-columns: 1.1rem 11rem minmax(11rem, 1fr) minmax(13rem, 1.3fr) 10rem; border-top: 1px solid color-mix(in srgb, CanvasText 24%, Canvas); }
.feed-header, .evaluation > summary, .system-event { display: grid; grid-template-columns: var(--activity-columns); gap: .5rem; align-items: center; }
.feed-header { padding: .3rem .25rem; border-bottom: 1px solid color-mix(in srgb, CanvasText 18%, Canvas); color: color-mix(in srgb, CanvasText 58%, Canvas); font-size: .72rem; font-weight: 700; letter-spacing: .04em; text-transform: uppercase; }
.feed-header span:last-child { text-align: right; }
.evaluation { display: block; }
.evaluation, .system-event { border-bottom: 1px solid color-mix(in srgb, CanvasText 18%, Canvas); }
.evaluation > summary { grid-template-areas: 'arrow time automation trigger result'; padding: .35rem .25rem; cursor: pointer; list-style: none; font-size: .9rem; }
.evaluation > summary:hover { background: color-mix(in srgb, CanvasText 4%, Canvas); }
.evaluation > summary::-webkit-details-marker { display: none; }
.evaluation > summary::before { grid-area: arrow; content: '›'; font-size: 1.25rem; line-height: 1; transform: rotate(0deg); }
.evaluation[open] > summary::before { transform: rotate(90deg); }
.evaluation-main { grid-area: automation; display: grid; min-width: 0; overflow-wrap: anywhere; }
.evaluation .time { grid-area: time; display: grid; min-width: 0; font-variant-numeric: tabular-nums; }
.evaluation .time time, .evaluation .time .relative { white-space: nowrap; }
.trigger { grid-area: trigger; min-width: 0; overflow-wrap: anywhere; }
.evaluation-meta { grid-area: result; display: flex; flex-wrap: wrap; justify-content: flex-end; align-items: center; gap: .3rem .45rem; min-width: 0; }
.evaluation-meta > .badge:first-child { min-width: 5.5rem; text-align: center; }
.badge { display: inline-block; padding: .1rem .35rem; border-radius: .2rem; border: 1px solid currentColor; font-size: .78rem; white-space: nowrap; }
.waiting { color: #8b5e00; }.running { color: #0066a6; }.completed { color: #08783f; }.skipped { color: #625f58; }.canceled, .incomplete { color: #a23d00; }
.warning-badge, .trace-error { color: #a23d00; }
.tooltip { border-bottom: 1px dotted currentColor; cursor: help; }
.duration { min-width: 3.5rem; text-align: right; white-space: nowrap; color: color-mix(in srgb, CanvasText 65%, Canvas); font-variant-numeric: tabular-nums; }
.trace { margin: .65rem 0 .15rem 1.5rem; padding-left: 1rem; border-left: 2px solid color-mix(in srgb, CanvasText 18%, Canvas); }
.trace li { display: grid; grid-template-columns: minmax(9rem, auto) 1fr auto; gap: .45rem; padding: .28rem 0; align-items: start; }
.raw { grid-column: 1 / -1; margin-top: .2rem; }
.raw summary { cursor: pointer; font-size: .82rem; color: color-mix(in srgb, CanvasText 64%, Canvas); }
pre { max-width: 100%; overflow: auto; padding: .55rem; background: color-mix(in srgb, CanvasText 6%, Canvas); border: 1px solid color-mix(in srgb, CanvasText 14%, Canvas); white-space: pre-wrap; overflow-wrap: anywhere; }
.system-event { grid-template-areas: 'marker time automation trigger result'; padding: .35rem .25rem; background: color-mix(in srgb, CanvasText 4%, Canvas); font-size: .9rem; }
.system-event::before { grid-area: marker; content: ''; }
.system-event > time { grid-area: time; min-width: 0; }
.system-event > strong { grid-area: automation; }
.system-event > span { grid-area: trigger; min-width: 0; overflow-wrap: anywhere; }
.system-event .raw { grid-area: result; margin-top: 0; }
.load-older, .end, .empty { display: block; margin: .9rem auto; width: fit-content; }
.catalog-heading { margin-bottom: .7rem; }
.catalog-table-wrap { overflow-x: auto; }
.catalog { border-collapse: collapse; width: 100%; min-width: 850px; }
.catalog th, .catalog td { border-bottom: 1px solid color-mix(in srgb, CanvasText 20%, Canvas); padding: .45rem; text-align: left; vertical-align: top; }
.catalog th { font-size: .8rem; color: color-mix(in srgb, CanvasText 65%, Canvas); }
.catalog th.automation-key, .catalog td.automation-key > details > summary { white-space: nowrap; }
.catalog .automation-time { width: 1%; white-space: nowrap; font-variant-numeric: tabular-nums; }
.catalog td.automation-time .relative { display: block; }
.catalog details summary { cursor: pointer; }
.automation-details { min-width: 27rem; padding: .6rem 0; }
@media (max-width: 850px) {
  .topbar { padding: .55rem 1rem; }
  main { padding: .75rem; }
  .summary { grid-template-columns: 1fr; }
  .summary dl { grid-template-columns: repeat(2, minmax(8rem, 1fr)); }
  .feed-header { display: none; }
  .evaluation > summary { grid-template-columns: 1.1rem minmax(0, 1fr); grid-template-areas: 'arrow time' 'arrow automation' 'arrow trigger' 'arrow result'; gap: .15rem .5rem; }
  .evaluation > summary::before { align-self: start; margin-top: .15rem; }
  .evaluation-meta { justify-content: flex-start; }
  .trace li { grid-template-columns: 1fr; gap: .2rem; }
  .system-event { grid-template-columns: 1.1rem minmax(0, 1fr); grid-template-areas: 'marker time' 'marker automation' 'marker trigger' 'marker result'; gap: .15rem .5rem; }
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
  document.body.addEventListener('click', function (event) { var button = event.target.closest('[data-show-new]'); if (!button) return; event.preventDefault(); window.scrollTo({top: 0, behavior: 'smooth'}); refreshActivity(true); });
  document.body.addEventListener('htmx:configRequest', function (event) {
    if (event.detail.path === '/activity/fragment') new URL(window.location.href).searchParams.getAll('expand').forEach(function (value) { event.detail.parameters.append('expand', value); });
  });
  document.body.addEventListener('htmx:afterSwap', function (event) { localize(event.target); updateBanner(); });
  document.body.addEventListener('htmx:sseMessage', function (event) {
    if (!event.detail || event.detail.type !== 'change' || !document.getElementById('activity-content')) return;
    if (atNewestUninterrupted()) refreshActivity(); else { pending++; updateBanner(); }
  });
}());`
