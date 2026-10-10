// Shared behaviour for the landing, gallery and reference pages.
// Parallax layers, the chapter rail, and the install copy buttons. No dependencies.
(function () {
  var reduce = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  var narrow = window.matchMedia('(max-width: 900px)');

  // Copy-to-clipboard on install blocks.
  document.querySelectorAll('.copy-btn').forEach(function (btn) {
    btn.addEventListener('click', function () {
      var text = btn.parentElement.querySelector('pre').textContent;
      navigator.clipboard.writeText(text).then(function () {
        var label = btn.textContent;
        btn.textContent = 'Copied';
        setTimeout(function () { btn.textContent = label; }, 1500);
      });
    });
  });

  // Parallax: each [data-speed] layer shifts by (its host's distance from the viewport
  // centre) x speed. Faster layers move more, so they read as closer.
  var layers = Array.prototype.map.call(document.querySelectorAll('[data-speed]'), function (el) {
    return { el: el, speed: Number(el.dataset.speed), host: el.closest('.chapter, .hero-l') };
  });
  var queued = false;
  function paint() {
    queued = false;
    if (narrow.matches) { layers.forEach(function (l) { l.el.style.removeProperty('--py'); }); return; }
    var vh = window.innerHeight;
    var progress = new Map();
    layers.forEach(function (l) {
      if (!progress.has(l.host)) {
        var r = l.host.getBoundingClientRect();
        progress.set(l.host, (r.top + r.height / 2 - vh / 2) / vh);
      }
      l.el.style.setProperty('--py', (progress.get(l.host) * l.speed).toFixed(1) + 'px');
    });
  }
  function queue() { if (!queued) { queued = true; requestAnimationFrame(paint); } }
  if (layers.length && !reduce) {
    window.addEventListener('scroll', queue, { passive: true });
    window.addEventListener('resize', queue);
    paint();
  }

  // Chapter rail: mark the chapter nearest the viewport centre.
  var dots = new Map(Array.prototype.map.call(document.querySelectorAll('.rail a'), function (a) { return [a.dataset.for, a]; }));
  if (dots.size) {
    var seen = new IntersectionObserver(function (entries) {
      entries.filter(function (e) { return e.isIntersecting; }).forEach(function (e) {
        dots.forEach(function (a, id) { a.classList.toggle('is-active', id === e.target.id); });
      });
    }, { rootMargin: '-45% 0px -45% 0px' });
    document.querySelectorAll('.chapter').forEach(function (c) { seen.observe(c); });
  }

  // Mini accent-theme preview (a teaser of the gallery's full switcher).
  var wrap = document.getElementById('mini-switcher');
  if (wrap) {
    [
      { label: 'House', accent: '#3ea6ff' },
      { label: 'Wes Anderson', accent: '#f4a8b9' },
      { label: 'The Matrix', accent: '#22e06a' },
      { label: 'Breaking Bad', accent: '#7cb518' }
    ].forEach(function (t, i) {
      var btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'mini-swatch';
      btn.style.background = t.accent;
      btn.setAttribute('aria-pressed', i === 0 ? 'true' : 'false');
      btn.setAttribute('aria-label', t.label);
      btn.title = t.label;
      btn.addEventListener('click', function () {
        wrap.querySelectorAll('.mini-swatch').forEach(function (b) { b.setAttribute('aria-pressed', 'false'); });
        btn.setAttribute('aria-pressed', 'true');
      });
      wrap.appendChild(btn);
    });
  }
})();
