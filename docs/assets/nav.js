// Scroll-linked menu color. The bar stays cream over cream sections; as another color
// scrolls up toward it (a dark band, the orange finale, a deck's own surface) the bar
// blends into that color in step with the scroll. The bar samples the real background
// color at points under it and just below it, so every page and every themed chapter
// works without being told what color it is.
(function () {
  var nav = document.querySelector('.site-nav');
  if (!nav) return;

  var NAV_H = nav.getBoundingClientRect().height || 58;
  var LOOKAHEAD = 150;                       // px below the bar that start to tint it
  var SAMPLE_X = 4;                          // the page margin: section backgrounds, not cards
  var STEP = 14;
  var INK = [20, 20, 19], CREAM = [243, 239, 230];

  function parse(css) {
    var m = css.match(/rgba?\(([^)]+)\)/);
    if (!m) return null;
    var p = m[1].split(',').map(Number);
    return { rgb: p.slice(0, 3), a: p.length > 3 ? p[3] : 1 };
  }
  function backgroundAt(y) {
    var els = document.elementsFromPoint(SAMPLE_X, y);
    for (var i = 0; i < els.length; i++) {
      if (els[i].closest('.site-nav, .rail')) continue;
      var c = parse(getComputedStyle(els[i]).backgroundColor);
      if (c && c.a > 0.9) return c.rgb;
    }
    return CREAM;
  }
  function lum(rgb) {
    var f = rgb.map(function (v) { v /= 255; return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4); });
    return 0.2126 * f[0] + 0.7152 * f[1] + 0.0722 * f[2];
  }
  function contrast(a, b) { var x = lum(a), y = lum(b); return (Math.max(x, y) + 0.05) / (Math.min(x, y) + 0.05); }
  function mix(a, b, t) { return a.map(function (v, i) { return Math.round(v + (b[i] - v) * t); }); }
  function css(rgb, a) { return a === undefined ? 'rgb(' + rgb.join(',') + ')' : 'rgba(' + rgb.join(',') + ',' + a + ')'; }

  var queued = false;
  function paint() {
    queued = false;
    var sum = [0, 0, 0], wsum = 0;
    for (var y = 2; y < NAV_H + LOOKAHEAD; y += STEP) {
      var w = y <= NAV_H ? 1 : 1 - (y - NAV_H) / LOOKAHEAD;     // full weight under the bar, fading below it
      var c = backgroundAt(y);
      sum[0] += c[0] * w; sum[1] += c[1] * w; sum[2] += c[2] * w; wsum += w;
    }
    var bg = sum.map(function (v) { return Math.round(v / wsum); });
    // ink or cream, whichever reads better on the blended color; blend between them near the crossover
    var s = Math.min(1, Math.max(0, (contrast(bg, CREAM) - contrast(bg, INK) + 0.8) / 1.6));
    var ink = mix(INK, CREAM, s), dim = mix(ink, bg, 0.16);   // links stay near the ink color, so contrast holds mid-blend
    var st = nav.style;
    st.setProperty('--nav-bg', css(bg, 0.97));
    st.setProperty('--nav-ink', css(ink));
    st.setProperty('--nav-dim', css(dim));
    st.setProperty('--nav-line', css(ink, 0.16));
  }
  function queue() { if (!queued) { queued = true; requestAnimationFrame(paint); } }

  window.addEventListener('scroll', queue, { passive: true });
  window.addEventListener('resize', queue);
  window.addEventListener('load', queue);
  paint();
})();
