import * as ye from "https://d3js.org/d3.v7.min.js"
/**
* @vue/shared v3.5.25
* (c) 2018-present Yuxi (Evan) You and Vue contributors
* @license MIT
**/
// @__NO_SIDE_EFFECTS__
function Js(e) {
  const t = /* @__PURE__ */ Object.create(null);
  for (const n of e.split(",")) t[n] = 1;
  return (n) => n in t;
}
const pe = {}, Dt = [], ot = () => {
}, yi = () => !1, Xn = (e) => e.charCodeAt(0) === 111 && e.charCodeAt(1) === 110 && // uppercase letter
  (e.charCodeAt(2) > 122 || e.charCodeAt(2) < 97), Ys = (e) => e.startsWith("onUpdate:"), xe = Object.assign, Xs = (e, t) => {
    const n = e.indexOf(t);
    n > -1 && e.splice(n, 1);
  }, il = Object.prototype.hasOwnProperty, ce = (e, t) => il.call(e, t), K = Array.isArray, $t = (e) => _n(e) === "[object Map]", Wt = (e) => _n(e) === "[object Set]", Sr = (e) => _n(e) === "[object Date]", Z = (e) => typeof e == "function", we = (e) => typeof e == "string", lt = (e) => typeof e == "symbol", ge = (e) => e !== null && typeof e == "object", _i = (e) => (ge(e) || Z(e)) && Z(e.then) && Z(e.catch), wi = Object.prototype.toString, _n = (e) => wi.call(e), ol = (e) => _n(e).slice(8, -1), Qn = (e) => _n(e) === "[object Object]", Qs = (e) => we(e) && e !== "NaN" && e[0] !== "-" && "" + parseInt(e, 10) === e, rn = /* @__PURE__ */ Js(
    // the leading comma is intentional so empty string "" is also included
    ",key,ref,ref_for,ref_key,onVnodeBeforeMount,onVnodeMounted,onVnodeBeforeUpdate,onVnodeUpdated,onVnodeBeforeUnmount,onVnodeUnmounted"
  ), Zn = (e) => {
    const t = /* @__PURE__ */ Object.create(null);
    return (n) => t[n] || (t[n] = e(n));
  }, ll = /-\w/g, Ne = Zn(
    (e) => e.replace(ll, (t) => t.slice(1).toUpperCase())
  ), al = /\B([A-Z])/g, ze = Zn(
    (e) => e.replace(al, "-$1").toLowerCase()
  ), es = Zn((e) => e.charAt(0).toUpperCase() + e.slice(1)), ys = Zn(
    (e) => e ? `on${es(e)}` : ""
  ), vt = (e, t) => !Object.is(e, t), Ln = (e, ...t) => {
    for (let n = 0; n < e.length; n++)
      e[n](...t);
  }, vi = (e, t, n, s = !1) => {
    Object.defineProperty(e, t, {
      configurable: !0,
      enumerable: !1,
      writable: s,
      value: n
    });
  }, xi = (e) => {
    const t = parseFloat(e);
    return isNaN(t) ? e : t;
  }, Cr = (e) => {
    const t = we(e) ? Number(e) : NaN;
    return isNaN(t) ? e : t;
  };
let Er;
const ts = () => Er || (Er = typeof globalThis < "u" ? globalThis : typeof self < "u" ? self : typeof window < "u" ? window : typeof global < "u" ? global : {});
function Zs(e) {
  if (K(e)) {
    const t = {};
    for (let n = 0; n < e.length; n++) {
      const s = e[n], r = we(s) ? dl(s) : Zs(s);
      if (r)
        for (const i in r)
          t[i] = r[i];
    }
    return t;
  } else if (we(e) || ge(e))
    return e;
}
const cl = /;(?![^(]*\))/g, ul = /:([^]+)/, fl = /\/\*[^]*?\*\//g;
function dl(e) {
  const t = {};
  return e.replace(fl, "").split(cl).forEach((n) => {
    if (n) {
      const s = n.split(ul);
      s.length > 1 && (t[s[0].trim()] = s[1].trim());
    }
  }), t;
}
function Lt(e) {
  let t = "";
  if (we(e))
    t = e;
  else if (K(e))
    for (let n = 0; n < e.length; n++) {
      const s = Lt(e[n]);
      s && (t += s + " ");
    }
  else if (ge(e))
    for (const n in e)
      e[n] && (t += n + " ");
  return t.trim();
}
const hl = "itemscope,allowfullscreen,formnovalidate,ismap,nomodule,novalidate,readonly", pl = /* @__PURE__ */ Js(hl);
function Si(e) {
  return !!e || e === "";
}
function gl(e, t) {
  if (e.length !== t.length) return !1;
  let n = !0;
  for (let s = 0; n && s < e.length; s++)
    n = wn(e[s], t[s]);
  return n;
}
function wn(e, t) {
  if (e === t) return !0;
  let n = Sr(e), s = Sr(t);
  if (n || s)
    return n && s ? e.getTime() === t.getTime() : !1;
  if (n = lt(e), s = lt(t), n || s)
    return e === t;
  if (n = K(e), s = K(t), n || s)
    return n && s ? gl(e, t) : !1;
  if (n = ge(e), s = ge(t), n || s) {
    if (!n || !s)
      return !1;
    const r = Object.keys(e).length, i = Object.keys(t).length;
    if (r !== i)
      return !1;
    for (const o in e) {
      const l = e.hasOwnProperty(o), a = t.hasOwnProperty(o);
      if (l && !a || !l && a || !wn(e[o], t[o]))
        return !1;
    }
  }
  return String(e) === String(t);
}
function er(e, t) {
  return e.findIndex((n) => wn(n, t));
}
const Ci = (e) => !!(e && e.__v_isRef === !0), se = (e) => we(e) ? e : e == null ? "" : K(e) || ge(e) && (e.toString === wi || !Z(e.toString)) ? Ci(e) ? se(e.value) : JSON.stringify(e, Ei, 2) : String(e), Ei = (e, t) => Ci(t) ? Ei(e, t.value) : $t(t) ? {
  [`Map(${t.size})`]: [...t.entries()].reduce(
    (n, [s, r], i) => (n[_s(s, i) + " =>"] = r, n),
    {}
  )
} : Wt(t) ? {
  [`Set(${t.size})`]: [...t.values()].map((n) => _s(n))
} : lt(t) ? _s(t) : ge(t) && !K(t) && !Qn(t) ? String(t) : t, _s = (e, t = "") => {
  var n;
  return (
    // Symbol.description in es2019+ so we need to cast here to pass
    // the lib: es2016 check
    lt(e) ? `Symbol(${(n = e.description) != null ? n : t})` : e
  );
};
/**
* @vue/reactivity v3.5.25
* (c) 2018-present Yuxi (Evan) You and Vue contributors
* @license MIT
**/
let De;
class ml {
  constructor(t = !1) {
    this.detached = t, this._active = !0, this._on = 0, this.effects = [], this.cleanups = [], this._isPaused = !1, this.parent = De, !t && De && (this.index = (De.scopes || (De.scopes = [])).push(
      this
    ) - 1);
  }
  get active() {
    return this._active;
  }
  pause() {
    if (this._active) {
      this._isPaused = !0;
      let t, n;
      if (this.scopes)
        for (t = 0, n = this.scopes.length; t < n; t++)
          this.scopes[t].pause();
      for (t = 0, n = this.effects.length; t < n; t++)
        this.effects[t].pause();
    }
  }
  /**
   * Resumes the effect scope, including all child scopes and effects.
   */
  resume() {
    if (this._active && this._isPaused) {
      this._isPaused = !1;
      let t, n;
      if (this.scopes)
        for (t = 0, n = this.scopes.length; t < n; t++)
          this.scopes[t].resume();
      for (t = 0, n = this.effects.length; t < n; t++)
        this.effects[t].resume();
    }
  }
  run(t) {
    if (this._active) {
      const n = De;
      try {
        return De = this, t();
      } finally {
        De = n;
      }
    }
  }
  /**
   * This should only be called on non-detached scopes
   * @internal
   */
  on() {
    ++this._on === 1 && (this.prevScope = De, De = this);
  }
  /**
   * This should only be called on non-detached scopes
   * @internal
   */
  off() {
    this._on > 0 && --this._on === 0 && (De = this.prevScope, this.prevScope = void 0);
  }
  stop(t) {
    if (this._active) {
      this._active = !1;
      let n, s;
      for (n = 0, s = this.effects.length; n < s; n++)
        this.effects[n].stop();
      for (this.effects.length = 0, n = 0, s = this.cleanups.length; n < s; n++)
        this.cleanups[n]();
      if (this.cleanups.length = 0, this.scopes) {
        for (n = 0, s = this.scopes.length; n < s; n++)
          this.scopes[n].stop(!0);
        this.scopes.length = 0;
      }
      if (!this.detached && this.parent && !t) {
        const r = this.parent.scopes.pop();
        r && r !== this && (this.parent.scopes[this.index] = r, r.index = this.index);
      }
      this.parent = void 0;
    }
  }
}
function bl() {
  return De;
}
let me;
const ws = /* @__PURE__ */ new WeakSet();
class Ti {
  constructor(t) {
    this.fn = t, this.deps = void 0, this.depsTail = void 0, this.flags = 5, this.next = void 0, this.cleanup = void 0, this.scheduler = void 0, De && De.active && De.effects.push(this);
  }
  pause() {
    this.flags |= 64;
  }
  resume() {
    this.flags & 64 && (this.flags &= -65, ws.has(this) && (ws.delete(this), this.trigger()));
  }
  /**
   * @internal
   */
  notify() {
    this.flags & 2 && !(this.flags & 32) || this.flags & 8 || ki(this);
  }
  run() {
    if (!(this.flags & 1))
      return this.fn();
    this.flags |= 2, Tr(this), Ri(this);
    const t = me, n = Ye;
    me = this, Ye = !0;
    try {
      return this.fn();
    } finally {
      Oi(this), me = t, Ye = n, this.flags &= -3;
    }
  }
  stop() {
    if (this.flags & 1) {
      for (let t = this.deps; t; t = t.nextDep)
        sr(t);
      this.deps = this.depsTail = void 0, Tr(this), this.onStop && this.onStop(), this.flags &= -2;
    }
  }
  trigger() {
    this.flags & 64 ? ws.add(this) : this.scheduler ? this.scheduler() : this.runIfDirty();
  }
  /**
   * @internal
   */
  runIfDirty() {
    Fs(this) && this.run();
  }
  get dirty() {
    return Fs(this);
  }
}
let Ai = 0, on, ln;
function ki(e, t = !1) {
  if (e.flags |= 8, t) {
    e.next = ln, ln = e;
    return;
  }
  e.next = on, on = e;
}
function tr() {
  Ai++;
}
function nr() {
  if (--Ai > 0)
    return;
  if (ln) {
    let t = ln;
    for (ln = void 0; t;) {
      const n = t.next;
      t.next = void 0, t.flags &= -9, t = n;
    }
  }
  let e;
  for (; on;) {
    let t = on;
    for (on = void 0; t;) {
      const n = t.next;
      if (t.next = void 0, t.flags &= -9, t.flags & 1)
        try {
          t.trigger();
        } catch (s) {
          e || (e = s);
        }
      t = n;
    }
  }
  if (e) throw e;
}
function Ri(e) {
  for (let t = e.deps; t; t = t.nextDep)
    t.version = -1, t.prevActiveLink = t.dep.activeLink, t.dep.activeLink = t;
}
function Oi(e) {
  let t, n = e.depsTail, s = n;
  for (; s;) {
    const r = s.prevDep;
    s.version === -1 ? (s === n && (n = r), sr(s), yl(s)) : t = s, s.dep.activeLink = s.prevActiveLink, s.prevActiveLink = void 0, s = r;
  }
  e.deps = t, e.depsTail = n;
}
function Fs(e) {
  for (let t = e.deps; t; t = t.nextDep)
    if (t.dep.version !== t.version || t.dep.computed && (Pi(t.dep.computed) || t.dep.version !== t.version))
      return !0;
  return !!e._dirty;
}
function Pi(e) {
  if (e.flags & 4 && !(e.flags & 16) || (e.flags &= -17, e.globalVersion === dn) || (e.globalVersion = dn, !e.isSSR && e.flags & 128 && (!e.deps && !e._dirty || !Fs(e))))
    return;
  e.flags |= 2;
  const t = e.dep, n = me, s = Ye;
  me = e, Ye = !0;
  try {
    Ri(e);
    const r = e.fn(e._value);
    (t.version === 0 || vt(r, e._value)) && (e.flags |= 128, e._value = r, t.version++);
  } catch (r) {
    throw t.version++, r;
  } finally {
    me = n, Ye = s, Oi(e), e.flags &= -3;
  }
}
function sr(e, t = !1) {
  const { dep: n, prevSub: s, nextSub: r } = e;
  if (s && (s.nextSub = r, e.prevSub = void 0), r && (r.prevSub = s, e.nextSub = void 0), n.subs === e && (n.subs = s, !s && n.computed)) {
    n.computed.flags &= -5;
    for (let i = n.computed.deps; i; i = i.nextDep)
      sr(i, !0);
  }
  !t && !--n.sc && n.map && n.map.delete(n.key);
}
function yl(e) {
  const { prevDep: t, nextDep: n } = e;
  t && (t.nextDep = n, e.prevDep = void 0), n && (n.prevDep = t, e.nextDep = void 0);
}
let Ye = !0;
const Li = [];
function ht() {
  Li.push(Ye), Ye = !1;
}
function pt() {
  const e = Li.pop();
  Ye = e === void 0 ? !0 : e;
}
function Tr(e) {
  const { cleanup: t } = e;
  if (e.cleanup = void 0, t) {
    const n = me;
    me = void 0;
    try {
      t();
    } finally {
      me = n;
    }
  }
}
let dn = 0;
class _l {
  constructor(t, n) {
    this.sub = t, this.dep = n, this.version = n.version, this.nextDep = this.prevDep = this.nextSub = this.prevSub = this.prevActiveLink = void 0;
  }
}
class rr {
  // TODO isolatedDeclarations "__v_skip"
  constructor(t) {
    this.computed = t, this.version = 0, this.activeLink = void 0, this.subs = void 0, this.map = void 0, this.key = void 0, this.sc = 0, this.__v_skip = !0;
  }
  track(t) {
    if (!me || !Ye || me === this.computed)
      return;
    let n = this.activeLink;
    if (n === void 0 || n.sub !== me)
      n = this.activeLink = new _l(me, this), me.deps ? (n.prevDep = me.depsTail, me.depsTail.nextDep = n, me.depsTail = n) : me.deps = me.depsTail = n, Fi(n);
    else if (n.version === -1 && (n.version = this.version, n.nextDep)) {
      const s = n.nextDep;
      s.prevDep = n.prevDep, n.prevDep && (n.prevDep.nextDep = s), n.prevDep = me.depsTail, n.nextDep = void 0, me.depsTail.nextDep = n, me.depsTail = n, me.deps === n && (me.deps = s);
    }
    return n;
  }
  trigger(t) {
    this.version++, dn++, this.notify(t);
  }
  notify(t) {
    tr();
    try {
      for (let n = this.subs; n; n = n.prevSub)
        n.sub.notify() && n.sub.dep.notify();
    } finally {
      nr();
    }
  }
}
function Fi(e) {
  if (e.dep.sc++, e.sub.flags & 4) {
    const t = e.dep.computed;
    if (t && !e.dep.subs) {
      t.flags |= 20;
      for (let s = t.deps; s; s = s.nextDep)
        Fi(s);
    }
    const n = e.dep.subs;
    n !== e && (e.prevSub = n, n && (n.nextSub = e)), e.dep.subs = e;
  }
}
const Ms = /* @__PURE__ */ new WeakMap(), Rt = Symbol(
  ""
), Is = Symbol(
  ""
), hn = Symbol(
  ""
);
function Re(e, t, n) {
  if (Ye && me) {
    let s = Ms.get(e);
    s || Ms.set(e, s = /* @__PURE__ */ new Map());
    let r = s.get(n);
    r || (s.set(n, r = new rr()), r.map = s, r.key = n), r.track();
  }
}
function ft(e, t, n, s, r, i) {
  const o = Ms.get(e);
  if (!o) {
    dn++;
    return;
  }
  const l = (a) => {
    a && a.trigger();
  };
  if (tr(), t === "clear")
    o.forEach(l);
  else {
    const a = K(e), u = a && Qs(n);
    if (a && n === "length") {
      const c = Number(s);
      o.forEach((f, p) => {
        (p === "length" || p === hn || !lt(p) && p >= c) && l(f);
      });
    } else
      switch ((n !== void 0 || o.has(void 0)) && l(o.get(n)), u && l(o.get(hn)), t) {
        case "add":
          a ? u && l(o.get("length")) : (l(o.get(Rt)), $t(e) && l(o.get(Is)));
          break;
        case "delete":
          a || (l(o.get(Rt)), $t(e) && l(o.get(Is)));
          break;
        case "set":
          $t(e) && l(o.get(Rt));
          break;
      }
  }
  nr();
}
function It(e) {
  const t = ue(e);
  return t === e ? t : (Re(t, "iterate", hn), Ge(e) ? t : t.map(Qe));
}
function ns(e) {
  return Re(e = ue(e), "iterate", hn), e;
}
function yt(e, t) {
  return gt(e) ? Ot(e) ? Ut(Qe(t)) : Ut(t) : Qe(t);
}
const wl = {
  __proto__: null,
  [Symbol.iterator]() {
    return vs(this, Symbol.iterator, (e) => yt(this, e));
  },
  concat(...e) {
    return It(this).concat(
      ...e.map((t) => K(t) ? It(t) : t)
    );
  },
  entries() {
    return vs(this, "entries", (e) => (e[1] = yt(this, e[1]), e));
  },
  every(e, t) {
    return ct(this, "every", e, t, void 0, arguments);
  },
  filter(e, t) {
    return ct(
      this,
      "filter",
      e,
      t,
      (n) => n.map((s) => yt(this, s)),
      arguments
    );
  },
  find(e, t) {
    return ct(
      this,
      "find",
      e,
      t,
      (n) => yt(this, n),
      arguments
    );
  },
  findIndex(e, t) {
    return ct(this, "findIndex", e, t, void 0, arguments);
  },
  findLast(e, t) {
    return ct(
      this,
      "findLast",
      e,
      t,
      (n) => yt(this, n),
      arguments
    );
  },
  findLastIndex(e, t) {
    return ct(this, "findLastIndex", e, t, void 0, arguments);
  },
  // flat, flatMap could benefit from ARRAY_ITERATE but are not straight-forward to implement
  forEach(e, t) {
    return ct(this, "forEach", e, t, void 0, arguments);
  },
  includes(...e) {
    return xs(this, "includes", e);
  },
  indexOf(...e) {
    return xs(this, "indexOf", e);
  },
  join(e) {
    return It(this).join(e);
  },
  // keys() iterator only reads `length`, no optimization required
  lastIndexOf(...e) {
    return xs(this, "lastIndexOf", e);
  },
  map(e, t) {
    return ct(this, "map", e, t, void 0, arguments);
  },
  pop() {
    return Qt(this, "pop");
  },
  push(...e) {
    return Qt(this, "push", e);
  },
  reduce(e, ...t) {
    return Ar(this, "reduce", e, t);
  },
  reduceRight(e, ...t) {
    return Ar(this, "reduceRight", e, t);
  },
  shift() {
    return Qt(this, "shift");
  },
  // slice could use ARRAY_ITERATE but also seems to beg for range tracking
  some(e, t) {
    return ct(this, "some", e, t, void 0, arguments);
  },
  splice(...e) {
    return Qt(this, "splice", e);
  },
  toReversed() {
    return It(this).toReversed();
  },
  toSorted(e) {
    return It(this).toSorted(e);
  },
  toSpliced(...e) {
    return It(this).toSpliced(...e);
  },
  unshift(...e) {
    return Qt(this, "unshift", e);
  },
  values() {
    return vs(this, "values", (e) => yt(this, e));
  }
};
function vs(e, t, n) {
  const s = ns(e), r = s[t]();
  return s !== e && !Ge(e) && (r._next = r.next, r.next = () => {
    const i = r._next();
    return i.done || (i.value = n(i.value)), i;
  }), r;
}
const vl = Array.prototype;
function ct(e, t, n, s, r, i) {
  const o = ns(e), l = o !== e && !Ge(e), a = o[t];
  if (a !== vl[t]) {
    const f = a.apply(e, i);
    return l ? Qe(f) : f;
  }
  let u = n;
  o !== e && (l ? u = function (f, p) {
    return n.call(this, yt(e, f), p, e);
  } : n.length > 2 && (u = function (f, p) {
    return n.call(this, f, p, e);
  }));
  const c = a.call(o, u, s);
  return l && r ? r(c) : c;
}
function Ar(e, t, n, s) {
  const r = ns(e);
  let i = n;
  return r !== e && (Ge(e) ? n.length > 3 && (i = function (o, l, a) {
    return n.call(this, o, l, a, e);
  }) : i = function (o, l, a) {
    return n.call(this, o, yt(e, l), a, e);
  }), r[t](i, ...s);
}
function xs(e, t, n) {
  const s = ue(e);
  Re(s, "iterate", hn);
  const r = s[t](...n);
  return (r === -1 || r === !1) && ar(n[0]) ? (n[0] = ue(n[0]), s[t](...n)) : r;
}
function Qt(e, t, n = []) {
  ht(), tr();
  const s = ue(e)[t].apply(e, n);
  return nr(), pt(), s;
}
const xl = /* @__PURE__ */ Js("__proto__,__v_isRef,__isVue"), Mi = new Set(
  /* @__PURE__ */ Object.getOwnPropertyNames(Symbol).filter((e) => e !== "arguments" && e !== "caller").map((e) => Symbol[e]).filter(lt)
);
function Sl(e) {
  lt(e) || (e = String(e));
  const t = ue(this);
  return Re(t, "has", e), t.hasOwnProperty(e);
}
class Ii {
  constructor(t = !1, n = !1) {
    this._isReadonly = t, this._isShallow = n;
  }
  get(t, n, s) {
    if (n === "__v_skip") return t.__v_skip;
    const r = this._isReadonly, i = this._isShallow;
    if (n === "__v_isReactive")
      return !r;
    if (n === "__v_isReadonly")
      return r;
    if (n === "__v_isShallow")
      return i;
    if (n === "__v_raw")
      return s === (r ? i ? Fl : ji : i ? $i : Di).get(t) || // receiver is not the reactive proxy, but has the same prototype
        // this means the receiver is a user proxy of the reactive proxy
        Object.getPrototypeOf(t) === Object.getPrototypeOf(s) ? t : void 0;
    const o = K(t);
    if (!r) {
      let a;
      if (o && (a = wl[n]))
        return a;
      if (n === "hasOwnProperty")
        return Sl;
    }
    const l = Reflect.get(
      t,
      n,
      // if this is a proxy wrapping a ref, return methods using the raw ref
      // as receiver so that we don't have to call `toRaw` on the ref in all
      // its class methods
      Le(t) ? t : s
    );
    if ((lt(n) ? Mi.has(n) : xl(n)) || (r || Re(t, "get", n), i))
      return l;
    if (Le(l)) {
      const a = o && Qs(n) ? l : l.value;
      return r && ge(a) ? Ds(a) : a;
    }
    return ge(l) ? r ? Ds(l) : or(l) : l;
  }
}
class Ni extends Ii {
  constructor(t = !1) {
    super(!1, t);
  }
  set(t, n, s, r) {
    let i = t[n];
    const o = K(t) && Qs(n);
    if (!this._isShallow) {
      const u = gt(i);
      if (!Ge(s) && !gt(s) && (i = ue(i), s = ue(s)), !o && Le(i) && !Le(s))
        return u || (i.value = s), !0;
    }
    const l = o ? Number(n) < t.length : ce(t, n), a = Reflect.set(
      t,
      n,
      s,
      Le(t) ? t : r
    );
    return t === ue(r) && (l ? vt(s, i) && ft(t, "set", n, s) : ft(t, "add", n, s)), a;
  }
  deleteProperty(t, n) {
    const s = ce(t, n);
    t[n];
    const r = Reflect.deleteProperty(t, n);
    return r && s && ft(t, "delete", n, void 0), r;
  }
  has(t, n) {
    const s = Reflect.has(t, n);
    return (!lt(n) || !Mi.has(n)) && Re(t, "has", n), s;
  }
  ownKeys(t) {
    return Re(
      t,
      "iterate",
      K(t) ? "length" : Rt
    ), Reflect.ownKeys(t);
  }
}
class Cl extends Ii {
  constructor(t = !1) {
    super(!0, t);
  }
  set(t, n) {
    return !0;
  }
  deleteProperty(t, n) {
    return !0;
  }
}
const El = /* @__PURE__ */ new Ni(), Tl = /* @__PURE__ */ new Cl(), Al = /* @__PURE__ */ new Ni(!0);
const Ns = (e) => e, kn = (e) => Reflect.getPrototypeOf(e);
function kl(e, t, n) {
  return function (...s) {
    const r = this.__v_raw, i = ue(r), o = $t(i), l = e === "entries" || e === Symbol.iterator && o, a = e === "keys" && o, u = r[e](...s), c = n ? Ns : t ? Ut : Qe;
    return !t && Re(
      i,
      "iterate",
      a ? Is : Rt
    ), {
      // iterator protocol
      next() {
        const { value: f, done: p } = u.next();
        return p ? { value: f, done: p } : {
          value: l ? [c(f[0]), c(f[1])] : c(f),
          done: p
        };
      },
      // iterable protocol
      [Symbol.iterator]() {
        return this;
      }
    };
  };
}
function Rn(e) {
  return function (...t) {
    return e === "delete" ? !1 : e === "clear" ? void 0 : this;
  };
}
function Rl(e, t) {
  const n = {
    get(r) {
      const i = this.__v_raw, o = ue(i), l = ue(r);
      e || (vt(r, l) && Re(o, "get", r), Re(o, "get", l));
      const { has: a } = kn(o), u = t ? Ns : e ? Ut : Qe;
      if (a.call(o, r))
        return u(i.get(r));
      if (a.call(o, l))
        return u(i.get(l));
      i !== o && i.get(r);
    },
    get size() {
      const r = this.__v_raw;
      return !e && Re(ue(r), "iterate", Rt), r.size;
    },
    has(r) {
      const i = this.__v_raw, o = ue(i), l = ue(r);
      return e || (vt(r, l) && Re(o, "has", r), Re(o, "has", l)), r === l ? i.has(r) : i.has(r) || i.has(l);
    },
    forEach(r, i) {
      const o = this, l = o.__v_raw, a = ue(l), u = t ? Ns : e ? Ut : Qe;
      return !e && Re(a, "iterate", Rt), l.forEach((c, f) => r.call(i, u(c), u(f), o));
    }
  };
  return xe(
    n,
    e ? {
      add: Rn("add"),
      set: Rn("set"),
      delete: Rn("delete"),
      clear: Rn("clear")
    } : {
      add(r) {
        !t && !Ge(r) && !gt(r) && (r = ue(r));
        const i = ue(this);
        return kn(i).has.call(i, r) || (i.add(r), ft(i, "add", r, r)), this;
      },
      set(r, i) {
        !t && !Ge(i) && !gt(i) && (i = ue(i));
        const o = ue(this), { has: l, get: a } = kn(o);
        let u = l.call(o, r);
        u || (r = ue(r), u = l.call(o, r));
        const c = a.call(o, r);
        return o.set(r, i), u ? vt(i, c) && ft(o, "set", r, i) : ft(o, "add", r, i), this;
      },
      delete(r) {
        const i = ue(this), { has: o, get: l } = kn(i);
        let a = o.call(i, r);
        a || (r = ue(r), a = o.call(i, r)), l && l.call(i, r);
        const u = i.delete(r);
        return a && ft(i, "delete", r, void 0), u;
      },
      clear() {
        const r = ue(this), i = r.size !== 0, o = r.clear();
        return i && ft(
          r,
          "clear",
          void 0,
          void 0
        ), o;
      }
    }
  ), [
    "keys",
    "values",
    "entries",
    Symbol.iterator
  ].forEach((r) => {
    n[r] = kl(r, e, t);
  }), n;
}
function ir(e, t) {
  const n = Rl(e, t);
  return (s, r, i) => r === "__v_isReactive" ? !e : r === "__v_isReadonly" ? e : r === "__v_raw" ? s : Reflect.get(
    ce(n, r) && r in s ? n : s,
    r,
    i
  );
}
const Ol = {
  get: /* @__PURE__ */ ir(!1, !1)
}, Pl = {
  get: /* @__PURE__ */ ir(!1, !0)
}, Ll = {
  get: /* @__PURE__ */ ir(!0, !1)
};
const Di = /* @__PURE__ */ new WeakMap(), $i = /* @__PURE__ */ new WeakMap(), ji = /* @__PURE__ */ new WeakMap(), Fl = /* @__PURE__ */ new WeakMap();
function Ml(e) {
  switch (e) {
    case "Object":
    case "Array":
      return 1;
    case "Map":
    case "Set":
    case "WeakMap":
    case "WeakSet":
      return 2;
    default:
      return 0;
  }
}
function Il(e) {
  return e.__v_skip || !Object.isExtensible(e) ? 0 : Ml(ol(e));
}
function or(e) {
  return gt(e) ? e : lr(
    e,
    !1,
    El,
    Ol,
    Di
  );
}
function Nl(e) {
  return lr(
    e,
    !1,
    Al,
    Pl,
    $i
  );
}
function Ds(e) {
  return lr(
    e,
    !0,
    Tl,
    Ll,
    ji
  );
}
function lr(e, t, n, s, r) {
  if (!ge(e) || e.__v_raw && !(t && e.__v_isReactive))
    return e;
  const i = Il(e);
  if (i === 0)
    return e;
  const o = r.get(e);
  if (o)
    return o;
  const l = new Proxy(
    e,
    i === 2 ? s : n
  );
  return r.set(e, l), l;
}
function Ot(e) {
  return gt(e) ? Ot(e.__v_raw) : !!(e && e.__v_isReactive);
}
function gt(e) {
  return !!(e && e.__v_isReadonly);
}
function Ge(e) {
  return !!(e && e.__v_isShallow);
}
function ar(e) {
  return e ? !!e.__v_raw : !1;
}
function ue(e) {
  const t = e && e.__v_raw;
  return t ? ue(t) : e;
}
function Dl(e) {
  return !ce(e, "__v_skip") && Object.isExtensible(e) && vi(e, "__v_skip", !0), e;
}
const Qe = (e) => ge(e) ? or(e) : e, Ut = (e) => ge(e) ? Ds(e) : e;
function Le(e) {
  return e ? e.__v_isRef === !0 : !1;
}
function he(e) {
  return $l(e, !1);
}
function $l(e, t) {
  return Le(e) ? e : new jl(e, t);
}
class jl {
  constructor(t, n) {
    this.dep = new rr(), this.__v_isRef = !0, this.__v_isShallow = !1, this._rawValue = n ? t : ue(t), this._value = n ? t : Qe(t), this.__v_isShallow = n;
  }
  get value() {
    return this.dep.track(), this._value;
  }
  set value(t) {
    const n = this._rawValue, s = this.__v_isShallow || Ge(t) || gt(t);
    t = s ? t : ue(t), vt(t, n) && (this._rawValue = t, this._value = s ? t : Qe(t), this.dep.trigger());
  }
}
function Bi(e) {
  return Le(e) ? e.value : e;
}
const Bl = {
  get: (e, t, n) => t === "__v_raw" ? e : Bi(Reflect.get(e, t, n)),
  set: (e, t, n, s) => {
    const r = e[t];
    return Le(r) && !Le(n) ? (r.value = n, !0) : Reflect.set(e, t, n, s);
  }
};
function Hi(e) {
  return Ot(e) ? e : new Proxy(e, Bl);
}
class Hl {
  constructor(t, n, s) {
    this.fn = t, this.setter = n, this._value = void 0, this.dep = new rr(this), this.__v_isRef = !0, this.deps = void 0, this.depsTail = void 0, this.flags = 16, this.globalVersion = dn - 1, this.next = void 0, this.effect = this, this.__v_isReadonly = !n, this.isSSR = s;
  }
  /**
   * @internal
   */
  notify() {
    if (this.flags |= 16, !(this.flags & 8) && // avoid infinite self recursion
      me !== this)
      return ki(this, !0), !0;
  }
  get value() {
    const t = this.dep.track();
    return Pi(this), t && (t.version = this.dep.version), this._value;
  }
  set value(t) {
    this.setter && this.setter(t);
  }
}
function Ul(e, t, n = !1) {
  let s, r;
  return Z(e) ? s = e : (s = e.get, r = e.set), new Hl(s, r, n);
}
const On = {}, jn = /* @__PURE__ */ new WeakMap();
let Tt;
function ql(e, t = !1, n = Tt) {
  if (n) {
    let s = jn.get(n);
    s || jn.set(n, s = []), s.push(e);
  }
}
function Vl(e, t, n = pe) {
  const { immediate: s, deep: r, once: i, scheduler: o, augmentJob: l, call: a } = n, u = (I) => r ? I : Ge(I) || r === !1 || r === 0 ? dt(I, 1) : dt(I);
  let c, f, p, g, m = !1, x = !1;
  if (Le(e) ? (f = () => e.value, m = Ge(e)) : Ot(e) ? (f = () => u(e), m = !0) : K(e) ? (x = !0, m = e.some((I) => Ot(I) || Ge(I)), f = () => e.map((I) => {
    if (Le(I))
      return I.value;
    if (Ot(I))
      return u(I);
    if (Z(I))
      return a ? a(I, 2) : I();
  })) : Z(e) ? t ? f = a ? () => a(e, 2) : e : f = () => {
    if (p) {
      ht();
      try {
        p();
      } finally {
        pt();
      }
    }
    const I = Tt;
    Tt = c;
    try {
      return a ? a(e, 3, [g]) : e(g);
    } finally {
      Tt = I;
    }
  } : f = ot, t && r) {
    const I = f, J = r === !0 ? 1 / 0 : r;
    f = () => dt(I(), J);
  }
  const A = bl(), H = () => {
    c.stop(), A && A.active && Xs(A.effects, c);
  };
  if (i && t) {
    const I = t;
    t = (...J) => {
      I(...J), H();
    };
  }
  let q = x ? new Array(e.length).fill(On) : On;
  const W = (I) => {
    if (!(!(c.flags & 1) || !c.dirty && !I))
      if (t) {
        const J = c.run();
        if (r || m || (x ? J.some((de, ie) => vt(de, q[ie])) : vt(J, q))) {
          p && p();
          const de = Tt;
          Tt = c;
          try {
            const ie = [
              J,
              // pass undefined as the old value when it's changed for the first time
              q === On ? void 0 : x && q[0] === On ? [] : q,
              g
            ];
            q = J, a ? a(t, 3, ie) : (
              // @ts-expect-error
              t(...ie)
            );
          } finally {
            Tt = de;
          }
        }
      } else
        c.run();
  };
  return l && l(W), c = new Ti(f), c.scheduler = o ? () => o(W, !1) : W, g = (I) => ql(I, !1, c), p = c.onStop = () => {
    const I = jn.get(c);
    if (I) {
      if (a)
        a(I, 4);
      else
        for (const J of I) J();
      jn.delete(c);
    }
  }, t ? s ? W(!0) : q = c.run() : o ? o(W.bind(null, !0), !0) : c.run(), H.pause = c.pause.bind(c), H.resume = c.resume.bind(c), H.stop = H, H;
}
function dt(e, t = 1 / 0, n) {
  if (t <= 0 || !ge(e) || e.__v_skip || (n = n || /* @__PURE__ */ new Map(), (n.get(e) || 0) >= t))
    return e;
  if (n.set(e, t), t--, Le(e))
    dt(e.value, t, n);
  else if (K(e))
    for (let s = 0; s < e.length; s++)
      dt(e[s], t, n);
  else if (Wt(e) || $t(e))
    e.forEach((s) => {
      dt(s, t, n);
    });
  else if (Qn(e)) {
    for (const s in e)
      dt(e[s], t, n);
    for (const s of Object.getOwnPropertySymbols(e))
      Object.prototype.propertyIsEnumerable.call(e, s) && dt(e[s], t, n);
  }
  return e;
}
/**
* @vue/runtime-core v3.5.25
* (c) 2018-present Yuxi (Evan) You and Vue contributors
* @license MIT
**/
function vn(e, t, n, s) {
  try {
    return s ? e(...s) : e();
  } catch (r) {
    ss(r, t, n);
  }
}
function at(e, t, n, s) {
  if (Z(e)) {
    const r = vn(e, t, n, s);
    return r && _i(r) && r.catch((i) => {
      ss(i, t, n);
    }), r;
  }
  if (K(e)) {
    const r = [];
    for (let i = 0; i < e.length; i++)
      r.push(at(e[i], t, n, s));
    return r;
  }
}
function ss(e, t, n, s = !0) {
  const r = t ? t.vnode : null, { errorHandler: i, throwUnhandledErrorInProduction: o } = t && t.appContext.config || pe;
  if (t) {
    let l = t.parent;
    const a = t.proxy, u = `https://vuejs.org/error-reference/#runtime-${n}`;
    for (; l;) {
      const c = l.ec;
      if (c) {
        for (let f = 0; f < c.length; f++)
          if (c[f](e, a, u) === !1)
            return;
      }
      l = l.parent;
    }
    if (i) {
      ht(), vn(i, null, 10, [
        e,
        a,
        u
      ]), pt();
      return;
    }
  }
  Wl(e, n, r, s, o);
}
function Wl(e, t, n, s = !0, r = !1) {
  if (r)
    throw e;
  console.error(e);
}
const Ie = [];
let rt = -1;
const jt = [];
let _t = null, Nt = 0;
const Ui = /* @__PURE__ */ Promise.resolve();
let Bn = null;
function rs(e) {
  const t = Bn || Ui;
  return e ? t.then(this ? e.bind(this) : e) : t;
}
function zl(e) {
  let t = rt + 1, n = Ie.length;
  for (; t < n;) {
    const s = t + n >>> 1, r = Ie[s], i = pn(r);
    i < e || i === e && r.flags & 2 ? t = s + 1 : n = s;
  }
  return t;
}
function cr(e) {
  if (!(e.flags & 1)) {
    const t = pn(e), n = Ie[Ie.length - 1];
    !n || // fast path when the job id is larger than the tail
      !(e.flags & 2) && t >= pn(n) ? Ie.push(e) : Ie.splice(zl(t), 0, e), e.flags |= 1, qi();
  }
}
function qi() {
  Bn || (Bn = Ui.then(Wi));
}
function Kl(e) {
  K(e) ? jt.push(...e) : _t && e.id === -1 ? _t.splice(Nt + 1, 0, e) : e.flags & 1 || (jt.push(e), e.flags |= 1), qi();
}
function kr(e, t, n = rt + 1) {
  for (; n < Ie.length; n++) {
    const s = Ie[n];
    if (s && s.flags & 2) {
      if (e && s.id !== e.uid)
        continue;
      Ie.splice(n, 1), n--, s.flags & 4 && (s.flags &= -2), s(), s.flags & 4 || (s.flags &= -2);
    }
  }
}
function Vi(e) {
  if (jt.length) {
    const t = [...new Set(jt)].sort(
      (n, s) => pn(n) - pn(s)
    );
    if (jt.length = 0, _t) {
      _t.push(...t);
      return;
    }
    for (_t = t, Nt = 0; Nt < _t.length; Nt++) {
      const n = _t[Nt];
      n.flags & 4 && (n.flags &= -2), n.flags & 8 || n(), n.flags &= -2;
    }
    _t = null, Nt = 0;
  }
}
const pn = (e) => e.id == null ? e.flags & 2 ? -1 : 1 / 0 : e.id;
function Wi(e) {
  try {
    for (rt = 0; rt < Ie.length; rt++) {
      const t = Ie[rt];
      t && !(t.flags & 8) && (t.flags & 4 && (t.flags &= -2), vn(
        t,
        t.i,
        t.i ? 15 : 14
      ), t.flags & 4 || (t.flags &= -2));
    }
  } finally {
    for (; rt < Ie.length; rt++) {
      const t = Ie[rt];
      t && (t.flags &= -2);
    }
    rt = -1, Ie.length = 0, Vi(), Bn = null, (Ie.length || jt.length) && Wi();
  }
}
let qe = null, zi = null;
function Hn(e) {
  const t = qe;
  return qe = e, zi = e && e.type.__scopeId || null, t;
}
function Gl(e, t = qe, n) {
  if (!t || e._n)
    return e;
  const s = (...r) => {
    s._d && Wn(-1);
    const i = Hn(t);
    let o;
    try {
      o = e(...r);
    } finally {
      Hn(i), s._d && Wn(1);
    }
    return o;
  };
  return s._n = !0, s._c = !0, s._d = !0, s;
}
function Ss(e, t) {
  if (qe === null)
    return e;
  const n = as(qe), s = e.dirs || (e.dirs = []);
  for (let r = 0; r < t.length; r++) {
    let [i, o, l, a = pe] = t[r];
    i && (Z(i) && (i = {
      mounted: i,
      updated: i
    }), i.deep && dt(o), s.push({
      dir: i,
      instance: n,
      value: o,
      oldValue: void 0,
      arg: l,
      modifiers: a
    }));
  }
  return e;
}
function Ct(e, t, n, s) {
  const r = e.dirs, i = t && t.dirs;
  for (let o = 0; o < r.length; o++) {
    const l = r[o];
    i && (l.oldValue = i[o].value);
    let a = l.dir[s];
    a && (ht(), at(a, n, 8, [
      e.el,
      l,
      e,
      t
    ]), pt());
  }
}
const Jl = Symbol("_vte"), Yl = (e) => e.__isTeleport, Xl = Symbol("_leaveCb");
function ur(e, t) {
  e.shapeFlag & 6 && e.component ? (e.transition = t, ur(e.component.subTree, t)) : e.shapeFlag & 128 ? (e.ssContent.transition = t.clone(e.ssContent), e.ssFallback.transition = t.clone(e.ssFallback)) : e.transition = t;
}
// @__NO_SIDE_EFFECTS__
function Ql(e, t) {
  return Z(e) ? (
    // #8236: extend call and options.name access are considered side-effects
    // by Rollup, so we have to wrap it in a pure-annotated IIFE.
    xe({ name: e.name }, t, { setup: e })
  ) : e;
}
function Ki(e) {
  e.ids = [e.ids[0] + e.ids[2]++ + "-", 0, 0];
}
const Un = /* @__PURE__ */ new WeakMap();
function an(e, t, n, s, r = !1) {
  if (K(e)) {
    e.forEach(
      (m, x) => an(
        m,
        t && (K(t) ? t[x] : t),
        n,
        s,
        r
      )
    );
    return;
  }
  if (cn(s) && !r) {
    s.shapeFlag & 512 && s.type.__asyncResolved && s.component.subTree.component && an(e, t, n, s.component.subTree);
    return;
  }
  const i = s.shapeFlag & 4 ? as(s.component) : s.el, o = r ? null : i, { i: l, r: a } = e, u = t && t.r, c = l.refs === pe ? l.refs = {} : l.refs, f = l.setupState, p = ue(f), g = f === pe ? yi : (m) => ce(p, m);
  if (u != null && u !== a) {
    if (Rr(t), we(u))
      c[u] = null, g(u) && (f[u] = null);
    else if (Le(u)) {
      u.value = null;
      const m = t;
      m.k && (c[m.k] = null);
    }
  }
  if (Z(a))
    vn(a, l, 12, [o, c]);
  else {
    const m = we(a), x = Le(a);
    if (m || x) {
      const A = () => {
        if (e.f) {
          const H = m ? g(a) ? f[a] : c[a] : a.value;
          if (r)
            K(H) && Xs(H, i);
          else if (K(H))
            H.includes(i) || H.push(i);
          else if (m)
            c[a] = [i], g(a) && (f[a] = c[a]);
          else {
            const q = [i];
            a.value = q, e.k && (c[e.k] = q);
          }
        } else m ? (c[a] = o, g(a) && (f[a] = o)) : x && (a.value = o, e.k && (c[e.k] = o));
      };
      if (o) {
        const H = () => {
          A(), Un.delete(e);
        };
        H.id = -1, Un.set(e, H), Ue(H, n);
      } else
        Rr(e), A();
    }
  }
}
function Rr(e) {
  const t = Un.get(e);
  t && (t.flags |= 8, Un.delete(e));
}
ts().requestIdleCallback;
ts().cancelIdleCallback;
const cn = (e) => !!e.type.__asyncLoader, Gi = (e) => e.type.__isKeepAlive;
function Zl(e, t) {
  Ji(e, "a", t);
}
function ea(e, t) {
  Ji(e, "da", t);
}
function Ji(e, t, n = Oe) {
  const s = e.__wdc || (e.__wdc = () => {
    let r = n;
    for (; r;) {
      if (r.isDeactivated)
        return;
      r = r.parent;
    }
    return e();
  });
  if (is(t, s, n), n) {
    let r = n.parent;
    for (; r && r.parent;)
      Gi(r.parent.vnode) && ta(s, t, n, r), r = r.parent;
  }
}
function ta(e, t, n, s) {
  const r = is(
    t,
    e,
    s,
    !0
    /* prepend */
  );
  dr(() => {
    Xs(s[t], r);
  }, n);
}
function is(e, t, n = Oe, s = !1) {
  if (n) {
    const r = n[e] || (n[e] = []), i = t.__weh || (t.__weh = (...o) => {
      ht();
      const l = xn(n), a = at(t, n, e, o);
      return l(), pt(), a;
    });
    return s ? r.unshift(i) : r.push(i), i;
  }
}
const mt = (e) => (t, n = Oe) => {
  (!bn || e === "sp") && is(e, (...s) => t(...s), n);
}, na = mt("bm"), fr = mt("m"), sa = mt(
  "bu"
), ra = mt("u"), ia = mt(
  "bum"
), dr = mt("um"), oa = mt(
  "sp"
), la = mt("rtg"), aa = mt("rtc");
function ca(e, t = Oe) {
  is("ec", e, t);
}
const Yi = "components";
function qn(e, t) {
  return Zi(Yi, e, !0, t) || e;
}
const Xi = Symbol.for("v-ndc");
function Qi(e) {
  return we(e) ? Zi(Yi, e, !1) || e : e || Xi;
}
function Zi(e, t, n = !0, s = !1) {
  const r = qe || Oe;
  if (r) {
    const i = r.type;
    {
      const l = Qa(
        i,
        !1
      );
      if (l && (l === t || l === Ne(t) || l === es(Ne(t))))
        return i;
    }
    const o = (
      // local registration
      // check instance[type] first which is resolved for options API
      Or(r[e] || i[e], t) || // global registration
      Or(r.appContext[e], t)
    );
    return !o && s ? i : o;
  }
}
function Or(e, t) {
  return e && (e[t] || e[Ne(t)] || e[es(Ne(t))]);
}
function Xe(e, t, n, s) {
  let r;
  const i = n, o = K(e);
  if (o || we(e)) {
    const l = o && Ot(e);
    let a = !1, u = !1;
    l && (a = !Ge(e), u = gt(e), e = ns(e)), r = new Array(e.length);
    for (let c = 0, f = e.length; c < f; c++)
      r[c] = t(
        a ? u ? Ut(Qe(e[c])) : Qe(e[c]) : e[c],
        c,
        void 0,
        i
      );
  } else if (typeof e == "number") {
    r = new Array(e);
    for (let l = 0; l < e; l++)
      r[l] = t(l + 1, l, void 0, i);
  } else if (ge(e))
    if (e[Symbol.iterator])
      r = Array.from(
        e,
        (l, a) => t(l, a, void 0, i)
      );
    else {
      const l = Object.keys(e);
      r = new Array(l.length);
      for (let a = 0, u = l.length; a < u; a++) {
        const c = l[a];
        r[a] = t(e[c], c, a, i);
      }
    }
  else
    r = [];
  return r;
}
const $s = (e) => e ? wo(e) ? as(e) : $s(e.parent) : null, un = (
  // Move PURE marker to new line to workaround compiler discarding it
  // due to type annotation
  /* @__PURE__ */ xe(/* @__PURE__ */ Object.create(null), {
  $: (e) => e,
  $el: (e) => e.vnode.el,
  $data: (e) => e.data,
  $props: (e) => e.props,
  $attrs: (e) => e.attrs,
  $slots: (e) => e.slots,
  $refs: (e) => e.refs,
  $parent: (e) => $s(e.parent),
  $root: (e) => $s(e.root),
  $host: (e) => e.ce,
  $emit: (e) => e.emit,
  $options: (e) => to(e),
  $forceUpdate: (e) => e.f || (e.f = () => {
    cr(e.update);
  }),
  $nextTick: (e) => e.n || (e.n = rs.bind(e.proxy)),
  $watch: (e) => va.bind(e)
})
), Cs = (e, t) => e !== pe && !e.__isScriptSetup && ce(e, t), ua = {
  get({ _: e }, t) {
    if (t === "__v_skip")
      return !0;
    const { ctx: n, setupState: s, data: r, props: i, accessCache: o, type: l, appContext: a } = e;
    if (t[0] !== "$") {
      const p = o[t];
      if (p !== void 0)
        switch (p) {
          case 1:
            return s[t];
          case 2:
            return r[t];
          case 4:
            return n[t];
          case 3:
            return i[t];
        }
      else {
        if (Cs(s, t))
          return o[t] = 1, s[t];
        if (r !== pe && ce(r, t))
          return o[t] = 2, r[t];
        if (ce(i, t))
          return o[t] = 3, i[t];
        if (n !== pe && ce(n, t))
          return o[t] = 4, n[t];
        js && (o[t] = 0);
      }
    }
    const u = un[t];
    let c, f;
    if (u)
      return t === "$attrs" && Re(e.attrs, "get", ""), u(e);
    if (
      // css module (injected by vue-loader)
      (c = l.__cssModules) && (c = c[t])
    )
      return c;
    if (n !== pe && ce(n, t))
      return o[t] = 4, n[t];
    if (
      // global properties
      f = a.config.globalProperties, ce(f, t)
    )
      return f[t];
  },
  set({ _: e }, t, n) {
    const { data: s, setupState: r, ctx: i } = e;
    return Cs(r, t) ? (r[t] = n, !0) : s !== pe && ce(s, t) ? (s[t] = n, !0) : ce(e.props, t) || t[0] === "$" && t.slice(1) in e ? !1 : (i[t] = n, !0);
  },
  has({
    _: { data: e, setupState: t, accessCache: n, ctx: s, appContext: r, props: i, type: o }
  }, l) {
    let a;
    return !!(n[l] || e !== pe && l[0] !== "$" && ce(e, l) || Cs(t, l) || ce(i, l) || ce(s, l) || ce(un, l) || ce(r.config.globalProperties, l) || (a = o.__cssModules) && a[l]);
  },
  defineProperty(e, t, n) {
    return n.get != null ? e._.accessCache[t] = 0 : ce(n, "value") && this.set(e, t, n.value, null), Reflect.defineProperty(e, t, n);
  }
};
function Pr(e) {
  return K(e) ? e.reduce(
    (t, n) => (t[n] = null, t),
    {}
  ) : e;
}
let js = !0;
function fa(e) {
  const t = to(e), n = e.proxy, s = e.ctx;
  js = !1, t.beforeCreate && Lr(t.beforeCreate, e, "bc");
  const {
    // state
    data: r,
    computed: i,
    methods: o,
    watch: l,
    provide: a,
    inject: u,
    // lifecycle
    created: c,
    beforeMount: f,
    mounted: p,
    beforeUpdate: g,
    updated: m,
    activated: x,
    deactivated: A,
    beforeDestroy: H,
    beforeUnmount: q,
    destroyed: W,
    unmounted: I,
    render: J,
    renderTracked: de,
    renderTriggered: ie,
    errorCaptured: Se,
    serverPrefetch: v,
    // public API
    expose: P,
    inheritAttrs: D,
    // assets
    components: Q,
    directives: te,
    filters: oe
  } = t;
  if (u && da(u, s, null), o)
    for (const X in o) {
      const G = o[X];
      Z(G) && (s[X] = G.bind(n));
    }
  if (r) {
    const X = r.call(n, n);
    ge(X) && (e.data = or(X));
  }
  if (js = !0, i)
    for (const X in i) {
      const G = i[X], b = Z(G) ? G.bind(n, n) : Z(G.get) ? G.get.bind(n, n) : ot, _ = !Z(G) && Z(G.set) ? G.set.bind(n) : ot, R = We({
        get: b,
        set: _
      });
      Object.defineProperty(s, X, {
        enumerable: !0,
        configurable: !0,
        get: () => R.value,
        set: (O) => R.value = O
      });
    }
  if (l)
    for (const X in l)
      eo(l[X], s, n, X);
  if (a) {
    const X = Z(a) ? a.call(n) : a;
    Reflect.ownKeys(X).forEach((G) => {
      ya(G, X[G]);
    });
  }
  c && Lr(c, e, "c");
  function z(X, G) {
    K(G) ? G.forEach((b) => X(b.bind(n))) : G && X(G.bind(n));
  }
  if (z(na, f), z(fr, p), z(sa, g), z(ra, m), z(Zl, x), z(ea, A), z(ca, Se), z(aa, de), z(la, ie), z(ia, q), z(dr, I), z(oa, v), K(P))
    if (P.length) {
      const X = e.exposed || (e.exposed = {});
      P.forEach((G) => {
        Object.defineProperty(X, G, {
          get: () => n[G],
          set: (b) => n[G] = b,
          enumerable: !0
        });
      });
    } else e.exposed || (e.exposed = {});
  J && e.render === ot && (e.render = J), D != null && (e.inheritAttrs = D), Q && (e.components = Q), te && (e.directives = te), v && Ki(e);
}
function da(e, t, n = ot) {
  K(e) && (e = Bs(e));
  for (const s in e) {
    const r = e[s];
    let i;
    ge(r) ? "default" in r ? i = Fn(
      r.from || s,
      r.default,
      !0
    ) : i = Fn(r.from || s) : i = Fn(r), Le(i) ? Object.defineProperty(t, s, {
      enumerable: !0,
      configurable: !0,
      get: () => i.value,
      set: (o) => i.value = o
    }) : t[s] = i;
  }
}
function Lr(e, t, n) {
  at(
    K(e) ? e.map((s) => s.bind(t.proxy)) : e.bind(t.proxy),
    t,
    n
  );
}
function eo(e, t, n, s) {
  let r = s.includes(".") ? ro(n, s) : () => n[s];
  if (we(e)) {
    const i = t[e];
    Z(i) && Ke(r, i);
  } else if (Z(e))
    Ke(r, e.bind(n));
  else if (ge(e))
    if (K(e))
      e.forEach((i) => eo(i, t, n, s));
    else {
      const i = Z(e.handler) ? e.handler.bind(n) : t[e.handler];
      Z(i) && Ke(r, i, e);
    }
}
function to(e) {
  const t = e.type, { mixins: n, extends: s } = t, {
    mixins: r,
    optionsCache: i,
    config: { optionMergeStrategies: o }
  } = e.appContext, l = i.get(t);
  let a;
  return l ? a = l : !r.length && !n && !s ? a = t : (a = {}, r.length && r.forEach(
    (u) => Vn(a, u, o, !0)
  ), Vn(a, t, o)), ge(t) && i.set(t, a), a;
}
function Vn(e, t, n, s = !1) {
  const { mixins: r, extends: i } = t;
  i && Vn(e, i, n, !0), r && r.forEach(
    (o) => Vn(e, o, n, !0)
  );
  for (const o in t)
    if (!(s && o === "expose")) {
      const l = ha[o] || n && n[o];
      e[o] = l ? l(e[o], t[o]) : t[o];
    }
  return e;
}
const ha = {
  data: Fr,
  props: Mr,
  emits: Mr,
  // objects
  methods: sn,
  computed: sn,
  // lifecycle
  beforeCreate: Me,
  created: Me,
  beforeMount: Me,
  mounted: Me,
  beforeUpdate: Me,
  updated: Me,
  beforeDestroy: Me,
  beforeUnmount: Me,
  destroyed: Me,
  unmounted: Me,
  activated: Me,
  deactivated: Me,
  errorCaptured: Me,
  serverPrefetch: Me,
  // assets
  components: sn,
  directives: sn,
  // watch
  watch: ga,
  // provide / inject
  provide: Fr,
  inject: pa
};
function Fr(e, t) {
  return t ? e ? function () {
    return xe(
      Z(e) ? e.call(this, this) : e,
      Z(t) ? t.call(this, this) : t
    );
  } : t : e;
}
function pa(e, t) {
  return sn(Bs(e), Bs(t));
}
function Bs(e) {
  if (K(e)) {
    const t = {};
    for (let n = 0; n < e.length; n++)
      t[e[n]] = e[n];
    return t;
  }
  return e;
}
function Me(e, t) {
  return e ? [...new Set([].concat(e, t))] : t;
}
function sn(e, t) {
  return e ? xe(/* @__PURE__ */ Object.create(null), e, t) : t;
}
function Mr(e, t) {
  return e ? K(e) && K(t) ? [.../* @__PURE__ */ new Set([...e, ...t])] : xe(
    /* @__PURE__ */ Object.create(null),
    Pr(e),
    Pr(t ?? {})
  ) : t;
}
function ga(e, t) {
  if (!e) return t;
  if (!t) return e;
  const n = xe(/* @__PURE__ */ Object.create(null), e);
  for (const s in t)
    n[s] = Me(e[s], t[s]);
  return n;
}
function no() {
  return {
    app: null,
    config: {
      isNativeTag: yi,
      performance: !1,
      globalProperties: {},
      optionMergeStrategies: {},
      errorHandler: void 0,
      warnHandler: void 0,
      compilerOptions: {}
    },
    mixins: [],
    components: {},
    directives: {},
    provides: /* @__PURE__ */ Object.create(null),
    optionsCache: /* @__PURE__ */ new WeakMap(),
    propsCache: /* @__PURE__ */ new WeakMap(),
    emitsCache: /* @__PURE__ */ new WeakMap()
  };
}
let ma = 0;
function ba(e, t) {
  return function (s, r = null) {
    Z(s) || (s = xe({}, s)), r != null && !ge(r) && (r = null);
    const i = no(), o = /* @__PURE__ */ new WeakSet(), l = [];
    let a = !1;
    const u = i.app = {
      _uid: ma++,
      _component: s,
      _props: r,
      _container: null,
      _context: i,
      _instance: null,
      version: ec,
      get config() {
        return i.config;
      },
      set config(c) {
      },
      use(c, ...f) {
        return o.has(c) || (c && Z(c.install) ? (o.add(c), c.install(u, ...f)) : Z(c) && (o.add(c), c(u, ...f))), u;
      },
      mixin(c) {
        return i.mixins.includes(c) || i.mixins.push(c), u;
      },
      component(c, f) {
        return f ? (i.components[c] = f, u) : i.components[c];
      },
      directive(c, f) {
        return f ? (i.directives[c] = f, u) : i.directives[c];
      },
      mount(c, f, p) {
        if (!a) {
          const g = u._ceVNode || Te(s, r);
          return g.appContext = i, p === !0 ? p = "svg" : p === !1 && (p = void 0), e(g, c, p), a = !0, u._container = c, c.__vue_app__ = u, as(g.component);
        }
      },
      onUnmount(c) {
        l.push(c);
      },
      unmount() {
        a && (at(
          l,
          u._instance,
          16
        ), e(null, u._container), delete u._container.__vue_app__);
      },
      provide(c, f) {
        return i.provides[c] = f, u;
      },
      runWithContext(c) {
        const f = Bt;
        Bt = u;
        try {
          return c();
        } finally {
          Bt = f;
        }
      }
    };
    return u;
  };
}
let Bt = null;
function ya(e, t) {
  if (Oe) {
    let n = Oe.provides;
    const s = Oe.parent && Oe.parent.provides;
    s === n && (n = Oe.provides = Object.create(s)), n[e] = t;
  }
}
function Fn(e, t, n = !1) {
  const s = Ka();
  if (s || Bt) {
    let r = Bt ? Bt._context.provides : s ? s.parent == null || s.ce ? s.vnode.appContext && s.vnode.appContext.provides : s.parent.provides : void 0;
    if (r && e in r)
      return r[e];
    if (arguments.length > 1)
      return n && Z(t) ? t.call(s && s.proxy) : t;
  }
}
const _a = Symbol.for("v-scx"), wa = () => Fn(_a);
function Ke(e, t, n) {
  return so(e, t, n);
}
function so(e, t, n = pe) {
  const { immediate: s, deep: r, flush: i, once: o } = n, l = xe({}, n), a = t && s || !t && i !== "post";
  let u;
  if (bn) {
    if (i === "sync") {
      const g = wa();
      u = g.__watcherHandles || (g.__watcherHandles = []);
    } else if (!a) {
      const g = () => {
      };
      return g.stop = ot, g.resume = ot, g.pause = ot, g;
    }
  }
  const c = Oe;
  l.call = (g, m, x) => at(g, c, m, x);
  let f = !1;
  i === "post" ? l.scheduler = (g) => {
    Ue(g, c && c.suspense);
  } : i !== "sync" && (f = !0, l.scheduler = (g, m) => {
    m ? g() : cr(g);
  }), l.augmentJob = (g) => {
    t && (g.flags |= 4), f && (g.flags |= 2, c && (g.id = c.uid, g.i = c));
  };
  const p = Vl(e, t, l);
  return bn && (u ? u.push(p) : a && p()), p;
}
function va(e, t, n) {
  const s = this.proxy, r = we(e) ? e.includes(".") ? ro(s, e) : () => s[e] : e.bind(s, s);
  let i;
  Z(t) ? i = t : (i = t.handler, n = t);
  const o = xn(this), l = so(r, i.bind(s), n);
  return o(), l;
}
function ro(e, t) {
  const n = t.split(".");
  return () => {
    let s = e;
    for (let r = 0; r < n.length && s; r++)
      s = s[n[r]];
    return s;
  };
}
const xa = (e, t) => t === "modelValue" || t === "model-value" ? e.modelModifiers : e[`${t}Modifiers`] || e[`${Ne(t)}Modifiers`] || e[`${ze(t)}Modifiers`];
function Sa(e, t, ...n) {
  if (e.isUnmounted) return;
  const s = e.vnode.props || pe;
  let r = n;
  const i = t.startsWith("update:"), o = i && xa(s, t.slice(7));
  o && (o.trim && (r = n.map((c) => we(c) ? c.trim() : c)), o.number && (r = n.map(xi)));
  let l, a = s[l = ys(t)] || // also try camelCase event handler (#2249)
    s[l = ys(Ne(t))];
  !a && i && (a = s[l = ys(ze(t))]), a && at(
    a,
    e,
    6,
    r
  );
  const u = s[l + "Once"];
  if (u) {
    if (!e.emitted)
      e.emitted = {};
    else if (e.emitted[l])
      return;
    e.emitted[l] = !0, at(
      u,
      e,
      6,
      r
    );
  }
}
const Ca = /* @__PURE__ */ new WeakMap();
function io(e, t, n = !1) {
  const s = n ? Ca : t.emitsCache, r = s.get(e);
  if (r !== void 0)
    return r;
  const i = e.emits;
  let o = {}, l = !1;
  if (!Z(e)) {
    const a = (u) => {
      const c = io(u, t, !0);
      c && (l = !0, xe(o, c));
    };
    !n && t.mixins.length && t.mixins.forEach(a), e.extends && a(e.extends), e.mixins && e.mixins.forEach(a);
  }
  return !i && !l ? (ge(e) && s.set(e, null), null) : (K(i) ? i.forEach((a) => o[a] = null) : xe(o, i), ge(e) && s.set(e, o), o);
}
function os(e, t) {
  return !e || !Xn(t) ? !1 : (t = t.slice(2).replace(/Once$/, ""), ce(e, t[0].toLowerCase() + t.slice(1)) || ce(e, ze(t)) || ce(e, t));
}
function Ir(e) {
  const {
    type: t,
    vnode: n,
    proxy: s,
    withProxy: r,
    propsOptions: [i],
    slots: o,
    attrs: l,
    emit: a,
    render: u,
    renderCache: c,
    props: f,
    data: p,
    setupState: g,
    ctx: m,
    inheritAttrs: x
  } = e, A = Hn(e);
  let H, q;
  try {
    if (n.shapeFlag & 4) {
      const I = r || s, J = I;
      H = it(
        u.call(
          J,
          I,
          c,
          f,
          g,
          p,
          m
        )
      ), q = l;
    } else {
      const I = t;
      H = it(
        I.length > 1 ? I(
          f,
          { attrs: l, slots: o, emit: a }
        ) : I(
          f,
          null
        )
      ), q = t.props ? l : Ea(l);
    }
  } catch (I) {
    fn.length = 0, ss(I, e, 1), H = Te(xt);
  }
  let W = H;
  if (q && x !== !1) {
    const I = Object.keys(q), { shapeFlag: J } = W;
    I.length && J & 7 && (i && I.some(Ys) && (q = Ta(
      q,
      i
    )), W = qt(W, q, !1, !0));
  }
  return n.dirs && (W = qt(W, null, !1, !0), W.dirs = W.dirs ? W.dirs.concat(n.dirs) : n.dirs), n.transition && ur(W, n.transition), H = W, Hn(A), H;
}
const Ea = (e) => {
  let t;
  for (const n in e)
    (n === "class" || n === "style" || Xn(n)) && ((t || (t = {}))[n] = e[n]);
  return t;
}, Ta = (e, t) => {
  const n = {};
  for (const s in e)
    (!Ys(s) || !(s.slice(9) in t)) && (n[s] = e[s]);
  return n;
};
function Aa(e, t, n) {
  const { props: s, children: r, component: i } = e, { props: o, children: l, patchFlag: a } = t, u = i.emitsOptions;
  if (t.dirs || t.transition)
    return !0;
  if (n && a >= 0) {
    if (a & 1024)
      return !0;
    if (a & 16)
      return s ? Nr(s, o, u) : !!o;
    if (a & 8) {
      const c = t.dynamicProps;
      for (let f = 0; f < c.length; f++) {
        const p = c[f];
        if (o[p] !== s[p] && !os(u, p))
          return !0;
      }
    }
  } else
    return (r || l) && (!l || !l.$stable) ? !0 : s === o ? !1 : s ? o ? Nr(s, o, u) : !0 : !!o;
  return !1;
}
function Nr(e, t, n) {
  const s = Object.keys(t);
  if (s.length !== Object.keys(e).length)
    return !0;
  for (let r = 0; r < s.length; r++) {
    const i = s[r];
    if (t[i] !== e[i] && !os(n, i))
      return !0;
  }
  return !1;
}
function ka({ vnode: e, parent: t }, n) {
  for (; t;) {
    const s = t.subTree;
    if (s.suspense && s.suspense.activeBranch === e && (s.el = e.el), s === e)
      (e = t.vnode).el = n, t = t.parent;
    else
      break;
  }
}
const oo = {}, lo = () => Object.create(oo), ao = (e) => Object.getPrototypeOf(e) === oo;
function Ra(e, t, n, s = !1) {
  const r = {}, i = lo();
  e.propsDefaults = /* @__PURE__ */ Object.create(null), co(e, t, r, i);
  for (const o in e.propsOptions[0])
    o in r || (r[o] = void 0);
  n ? e.props = s ? r : Nl(r) : e.type.props ? e.props = r : e.props = i, e.attrs = i;
}
function Oa(e, t, n, s) {
  const {
    props: r,
    attrs: i,
    vnode: { patchFlag: o }
  } = e, l = ue(r), [a] = e.propsOptions;
  let u = !1;
  if (
    // always force full diff in dev
    // - #1942 if hmr is enabled with sfc component
    // - vite#872 non-sfc component used by sfc component
    (s || o > 0) && !(o & 16)
  ) {
    if (o & 8) {
      const c = e.vnode.dynamicProps;
      for (let f = 0; f < c.length; f++) {
        let p = c[f];
        if (os(e.emitsOptions, p))
          continue;
        const g = t[p];
        if (a)
          if (ce(i, p))
            g !== i[p] && (i[p] = g, u = !0);
          else {
            const m = Ne(p);
            r[m] = Hs(
              a,
              l,
              m,
              g,
              e,
              !1
            );
          }
        else
          g !== i[p] && (i[p] = g, u = !0);
      }
    }
  } else {
    co(e, t, r, i) && (u = !0);
    let c;
    for (const f in l)
      (!t || // for camelCase
        !ce(t, f) && // it's possible the original props was passed in as kebab-case
        // and converted to camelCase (#955)
        ((c = ze(f)) === f || !ce(t, c))) && (a ? n && // for camelCase
          (n[f] !== void 0 || // for kebab-case
            n[c] !== void 0) && (r[f] = Hs(
              a,
              l,
              f,
              void 0,
              e,
              !0
            )) : delete r[f]);
    if (i !== l)
      for (const f in i)
        (!t || !ce(t, f)) && (delete i[f], u = !0);
  }
  u && ft(e.attrs, "set", "");
}
function co(e, t, n, s) {
  const [r, i] = e.propsOptions;
  let o = !1, l;
  if (t)
    for (let a in t) {
      if (rn(a))
        continue;
      const u = t[a];
      let c;
      r && ce(r, c = Ne(a)) ? !i || !i.includes(c) ? n[c] = u : (l || (l = {}))[c] = u : os(e.emitsOptions, a) || (!(a in s) || u !== s[a]) && (s[a] = u, o = !0);
    }
  if (i) {
    const a = ue(n), u = l || pe;
    for (let c = 0; c < i.length; c++) {
      const f = i[c];
      n[f] = Hs(
        r,
        a,
        f,
        u[f],
        e,
        !ce(u, f)
      );
    }
  }
  return o;
}
function Hs(e, t, n, s, r, i) {
  const o = e[n];
  if (o != null) {
    const l = ce(o, "default");
    if (l && s === void 0) {
      const a = o.default;
      if (o.type !== Function && !o.skipFactory && Z(a)) {
        const { propsDefaults: u } = r;
        if (n in u)
          s = u[n];
        else {
          const c = xn(r);
          s = u[n] = a.call(
            null,
            t
          ), c();
        }
      } else
        s = a;
      r.ce && r.ce._setProp(n, s);
    }
    o[
      0
      /* shouldCast */
    ] && (i && !l ? s = !1 : o[
      1
      /* shouldCastTrue */
    ] && (s === "" || s === ze(n)) && (s = !0));
  }
  return s;
}
const Pa = /* @__PURE__ */ new WeakMap();
function uo(e, t, n = !1) {
  const s = n ? Pa : t.propsCache, r = s.get(e);
  if (r)
    return r;
  const i = e.props, o = {}, l = [];
  let a = !1;
  if (!Z(e)) {
    const c = (f) => {
      a = !0;
      const [p, g] = uo(f, t, !0);
      xe(o, p), g && l.push(...g);
    };
    !n && t.mixins.length && t.mixins.forEach(c), e.extends && c(e.extends), e.mixins && e.mixins.forEach(c);
  }
  if (!i && !a)
    return ge(e) && s.set(e, Dt), Dt;
  if (K(i))
    for (let c = 0; c < i.length; c++) {
      const f = Ne(i[c]);
      Dr(f) && (o[f] = pe);
    }
  else if (i)
    for (const c in i) {
      const f = Ne(c);
      if (Dr(f)) {
        const p = i[c], g = o[f] = K(p) || Z(p) ? { type: p } : xe({}, p), m = g.type;
        let x = !1, A = !0;
        if (K(m))
          for (let H = 0; H < m.length; ++H) {
            const q = m[H], W = Z(q) && q.name;
            if (W === "Boolean") {
              x = !0;
              break;
            } else W === "String" && (A = !1);
          }
        else
          x = Z(m) && m.name === "Boolean";
        g[
          0
          /* shouldCast */
        ] = x, g[
        1
        /* shouldCastTrue */
        ] = A, (x || ce(g, "default")) && l.push(f);
      }
    }
  const u = [o, l];
  return ge(e) && s.set(e, u), u;
}
function Dr(e) {
  return e[0] !== "$" && !rn(e);
}
const hr = (e) => e === "_" || e === "_ctx" || e === "$stable", pr = (e) => K(e) ? e.map(it) : [it(e)], La = (e, t, n) => {
  if (t._n)
    return t;
  const s = Gl((...r) => pr(t(...r)), n);
  return s._c = !1, s;
}, fo = (e, t, n) => {
  const s = e._ctx;
  for (const r in e) {
    if (hr(r)) continue;
    const i = e[r];
    if (Z(i))
      t[r] = La(r, i, s);
    else if (i != null) {
      const o = pr(i);
      t[r] = () => o;
    }
  }
}, ho = (e, t) => {
  const n = pr(t);
  e.slots.default = () => n;
}, po = (e, t, n) => {
  for (const s in t)
    (n || !hr(s)) && (e[s] = t[s]);
}, Fa = (e, t, n) => {
  const s = e.slots = lo();
  if (e.vnode.shapeFlag & 32) {
    const r = t._;
    r ? (po(s, t, n), n && vi(s, "_", r, !0)) : fo(t, s);
  } else t && ho(e, t);
}, Ma = (e, t, n) => {
  const { vnode: s, slots: r } = e;
  let i = !0, o = pe;
  if (s.shapeFlag & 32) {
    const l = t._;
    l ? n && l === 1 ? i = !1 : po(r, t, n) : (i = !t.$stable, fo(t, r)), o = t;
  } else t && (ho(e, t), o = { default: 1 });
  if (i)
    for (const l in r)
      !hr(l) && o[l] == null && delete r[l];
}, Ue = ja;
function Ia(e) {
  return Na(e);
}
function Na(e, t) {
  const n = ts();
  n.__VUE__ = !0;
  const {
    insert: s,
    remove: r,
    patchProp: i,
    createElement: o,
    createText: l,
    createComment: a,
    setText: u,
    setElementText: c,
    parentNode: f,
    nextSibling: p,
    setScopeId: g = ot,
    insertStaticContent: m
  } = e, x = (d, h, w, T = null, S = null, C = null, F = void 0, L = null, k = !!h.dynamicChildren) => {
    if (d === h)
      return;
    d && !Zt(d, h) && (T = et(d), O(d, S, C, !0), d = null), h.patchFlag === -2 && (k = !1, h.dynamicChildren = null);
    const { type: E, ref: V, shapeFlag: N } = h;
    switch (E) {
      case ls:
        A(d, h, w, T);
        break;
      case xt:
        H(d, h, w, T);
        break;
      case Ts:
        d == null && q(h, w, T, F);
        break;
      case _e:
        Q(
          d,
          h,
          w,
          T,
          S,
          C,
          F,
          L,
          k
        );
        break;
      default:
        N & 1 ? J(
          d,
          h,
          w,
          T,
          S,
          C,
          F,
          L,
          k
        ) : N & 6 ? te(
          d,
          h,
          w,
          T,
          S,
          C,
          F,
          L,
          k
        ) : (N & 64 || N & 128) && E.process(
          d,
          h,
          w,
          T,
          S,
          C,
          F,
          L,
          k,
          bt
        );
    }
    V != null && S ? an(V, d && d.ref, C, h || d, !h) : V == null && d && d.ref != null && an(d.ref, null, C, d, !0);
  }, A = (d, h, w, T) => {
    if (d == null)
      s(
        h.el = l(h.children),
        w,
        T
      );
    else {
      const S = h.el = d.el;
      h.children !== d.children && u(S, h.children);
    }
  }, H = (d, h, w, T) => {
    d == null ? s(
      h.el = a(h.children || ""),
      w,
      T
    ) : h.el = d.el;
  }, q = (d, h, w, T) => {
    [d.el, d.anchor] = m(
      d.children,
      h,
      w,
      T,
      d.el,
      d.anchor
    );
  }, W = ({ el: d, anchor: h }, w, T) => {
    let S;
    for (; d && d !== h;)
      S = p(d), s(d, w, T), d = S;
    s(h, w, T);
  }, I = ({ el: d, anchor: h }) => {
    let w;
    for (; d && d !== h;)
      w = p(d), r(d), d = w;
    r(h);
  }, J = (d, h, w, T, S, C, F, L, k) => {
    if (h.type === "svg" ? F = "svg" : h.type === "math" && (F = "mathml"), d == null)
      de(
        h,
        w,
        T,
        S,
        C,
        F,
        L,
        k
      );
    else {
      const E = d.el && d.el._isVueCE ? d.el : null;
      try {
        E && E._beginPatch(), v(
          d,
          h,
          S,
          C,
          F,
          L,
          k
        );
      } finally {
        E && E._endPatch();
      }
    }
  }, de = (d, h, w, T, S, C, F, L) => {
    let k, E;
    const { props: V, shapeFlag: N, transition: U, dirs: Y } = d;
    if (k = d.el = o(
      d.type,
      C,
      V && V.is,
      V
    ), N & 8 ? c(k, d.children) : N & 16 && Se(
      d.children,
      k,
      null,
      T,
      S,
      Es(d, C),
      F,
      L
    ), Y && Ct(d, null, T, "created"), ie(k, d, d.scopeId, F, T), V) {
      for (const fe in V)
        fe !== "value" && !rn(fe) && i(k, fe, null, V[fe], C, T);
      "value" in V && i(k, "value", null, V.value, C), (E = V.onVnodeBeforeMount) && nt(E, T, d);
    }
    Y && Ct(d, null, T, "beforeMount");
    const ne = Da(S, U);
    ne && U.beforeEnter(k), s(k, h, w), ((E = V && V.onVnodeMounted) || ne || Y) && Ue(() => {
      E && nt(E, T, d), ne && U.enter(k), Y && Ct(d, null, T, "mounted");
    }, S);
  }, ie = (d, h, w, T, S) => {
    if (w && g(d, w), T)
      for (let C = 0; C < T.length; C++)
        g(d, T[C]);
    if (S) {
      let C = S.subTree;
      if (h === C || bo(C.type) && (C.ssContent === h || C.ssFallback === h)) {
        const F = S.vnode;
        ie(
          d,
          F,
          F.scopeId,
          F.slotScopeIds,
          S.parent
        );
      }
    }
  }, Se = (d, h, w, T, S, C, F, L, k = 0) => {
    for (let E = k; E < d.length; E++) {
      const V = d[E] = L ? wt(d[E]) : it(d[E]);
      x(
        null,
        V,
        h,
        w,
        T,
        S,
        C,
        F,
        L
      );
    }
  }, v = (d, h, w, T, S, C, F) => {
    const L = h.el = d.el;
    let { patchFlag: k, dynamicChildren: E, dirs: V } = h;
    k |= d.patchFlag & 16;
    const N = d.props || pe, U = h.props || pe;
    let Y;
    if (w && Et(w, !1), (Y = U.onVnodeBeforeUpdate) && nt(Y, w, h, d), V && Ct(h, d, w, "beforeUpdate"), w && Et(w, !0), (N.innerHTML && U.innerHTML == null || N.textContent && U.textContent == null) && c(L, ""), E ? P(
      d.dynamicChildren,
      E,
      L,
      w,
      T,
      Es(h, S),
      C
    ) : F || G(
      d,
      h,
      L,
      null,
      w,
      T,
      Es(h, S),
      C,
      !1
    ), k > 0) {
      if (k & 16)
        D(L, N, U, w, S);
      else if (k & 2 && N.class !== U.class && i(L, "class", null, U.class, S), k & 4 && i(L, "style", N.style, U.style, S), k & 8) {
        const ne = h.dynamicProps;
        for (let fe = 0; fe < ne.length; fe++) {
          const ae = ne[fe], Ce = N[ae], Ee = U[ae];
          (Ee !== Ce || ae === "value") && i(L, ae, Ce, Ee, S, w);
        }
      }
      k & 1 && d.children !== h.children && c(L, h.children);
    } else !F && E == null && D(L, N, U, w, S);
    ((Y = U.onVnodeUpdated) || V) && Ue(() => {
      Y && nt(Y, w, h, d), V && Ct(h, d, w, "updated");
    }, T);
  }, P = (d, h, w, T, S, C, F) => {
    for (let L = 0; L < h.length; L++) {
      const k = d[L], E = h[L], V = (
        // oldVNode may be an errored async setup() component inside Suspense
        // which will not have a mounted element
        k.el && // - In the case of a Fragment, we need to provide the actual parent
          // of the Fragment itself so it can move its children.
          (k.type === _e || // - In the case of different nodes, there is going to be a replacement
            // which also requires the correct parent container
            !Zt(k, E) || // - In the case of a component, it could contain anything.
            k.shapeFlag & 198) ? f(k.el) : (
          // In other cases, the parent container is not actually used so we
          // just pass the block element here to avoid a DOM parentNode call.
          w
        )
      );
      x(
        k,
        E,
        V,
        null,
        T,
        S,
        C,
        F,
        !0
      );
    }
  }, D = (d, h, w, T, S) => {
    if (h !== w) {
      if (h !== pe)
        for (const C in h)
          !rn(C) && !(C in w) && i(
            d,
            C,
            h[C],
            null,
            S,
            T
          );
      for (const C in w) {
        if (rn(C)) continue;
        const F = w[C], L = h[C];
        F !== L && C !== "value" && i(d, C, L, F, S, T);
      }
      "value" in w && i(d, "value", h.value, w.value, S);
    }
  }, Q = (d, h, w, T, S, C, F, L, k) => {
    const E = h.el = d ? d.el : l(""), V = h.anchor = d ? d.anchor : l("");
    let { patchFlag: N, dynamicChildren: U, slotScopeIds: Y } = h;
    Y && (L = L ? L.concat(Y) : Y), d == null ? (s(E, w, T), s(V, w, T), Se(
      // #10007
      // such fragment like `<></>` will be compiled into
      // a fragment which doesn't have a children.
      // In this case fallback to an empty array
      h.children || [],
      w,
      V,
      S,
      C,
      F,
      L,
      k
    )) : N > 0 && N & 64 && U && // #2715 the previous fragment could've been a BAILed one as a result
      // of renderSlot() with no valid children
      d.dynamicChildren ? (P(
        d.dynamicChildren,
        U,
        w,
        S,
        C,
        F,
        L
      ), // #2080 if the stable fragment has a key, it's a <template v-for> that may
        //  get moved around. Make sure all root level vnodes inherit el.
        // #2134 or if it's a component root, it may also get moved around
        // as the component is being moved.
        (h.key != null || S && h === S.subTree) && go(
          d,
          h,
          !0
          /* shallow */
        )) : G(
          d,
          h,
          w,
          V,
          S,
          C,
          F,
          L,
          k
        );
  }, te = (d, h, w, T, S, C, F, L, k) => {
    h.slotScopeIds = L, d == null ? h.shapeFlag & 512 ? S.ctx.activate(
      h,
      w,
      T,
      F,
      k
    ) : oe(
      h,
      w,
      T,
      S,
      C,
      F,
      k
    ) : le(d, h, k);
  }, oe = (d, h, w, T, S, C, F) => {
    const L = d.component = za(
      d,
      T,
      S
    );
    if (Gi(d) && (L.ctx.renderer = bt), Ga(L, !1, F), L.asyncDep) {
      if (S && S.registerDep(L, z, F), !d.el) {
        const k = L.subTree = Te(xt);
        H(null, k, h, w), d.placeholder = k.el;
      }
    } else
      z(
        L,
        d,
        h,
        w,
        S,
        C,
        F
      );
  }, le = (d, h, w) => {
    const T = h.component = d.component;
    if (Aa(d, h, w))
      if (T.asyncDep && !T.asyncResolved) {
        X(T, h, w);
        return;
      } else
        T.next = h, T.update();
    else
      h.el = d.el, T.vnode = h;
  }, z = (d, h, w, T, S, C, F) => {
    const L = () => {
      if (d.isMounted) {
        let { next: N, bu: U, u: Y, parent: ne, vnode: fe } = d;
        {
          const He = mo(d);
          if (He) {
            N && (N.el = fe.el, X(d, N, F)), He.asyncDep.then(() => {
              d.isUnmounted || L();
            });
            return;
          }
        }
        let ae = N, Ce;
        Et(d, !1), N ? (N.el = fe.el, X(d, N, F)) : N = fe, U && Ln(U), (Ce = N.props && N.props.onVnodeBeforeUpdate) && nt(Ce, ne, N, fe), Et(d, !0);
        const Ee = Ir(d), Be = d.subTree;
        d.subTree = Ee, x(
          Be,
          Ee,
          // parent may have changed if it's in a teleport
          f(Be.el),
          // anchor may have changed if it's in a fragment
          et(Be),
          d,
          S,
          C
        ), N.el = Ee.el, ae === null && ka(d, Ee.el), Y && Ue(Y, S), (Ce = N.props && N.props.onVnodeUpdated) && Ue(
          () => nt(Ce, ne, N, fe),
          S
        );
      } else {
        let N;
        const { el: U, props: Y } = h, { bm: ne, m: fe, parent: ae, root: Ce, type: Ee } = d, Be = cn(h);
        Et(d, !1), ne && Ln(ne), !Be && (N = Y && Y.onVnodeBeforeMount) && nt(N, ae, h), Et(d, !0);
        {
          Ce.ce && // @ts-expect-error _def is private
            Ce.ce._def.shadowRoot !== !1 && Ce.ce._injectChildStyle(Ee);
          const He = d.subTree = Ir(d);
          x(
            null,
            He,
            w,
            T,
            d,
            S,
            C
          ), h.el = He.el;
        }
        if (fe && Ue(fe, S), !Be && (N = Y && Y.onVnodeMounted)) {
          const He = h;
          Ue(
            () => nt(N, ae, He),
            S
          );
        }
        (h.shapeFlag & 256 || ae && cn(ae.vnode) && ae.vnode.shapeFlag & 256) && d.a && Ue(d.a, S), d.isMounted = !0, h = w = T = null;
      }
    };
    d.scope.on();
    const k = d.effect = new Ti(L);
    d.scope.off();
    const E = d.update = k.run.bind(k), V = d.job = k.runIfDirty.bind(k);
    V.i = d, V.id = d.uid, k.scheduler = () => cr(V), Et(d, !0), E();
  }, X = (d, h, w) => {
    h.component = d;
    const T = d.vnode.props;
    d.vnode = h, d.next = null, Oa(d, h.props, T, w), Ma(d, h.children, w), ht(), kr(d), pt();
  }, G = (d, h, w, T, S, C, F, L, k = !1) => {
    const E = d && d.children, V = d ? d.shapeFlag : 0, N = h.children, { patchFlag: U, shapeFlag: Y } = h;
    if (U > 0) {
      if (U & 128) {
        _(
          E,
          N,
          w,
          T,
          S,
          C,
          F,
          L,
          k
        );
        return;
      } else if (U & 256) {
        b(
          E,
          N,
          w,
          T,
          S,
          C,
          F,
          L,
          k
        );
        return;
      }
    }
    Y & 8 ? (V & 16 && Je(E, S, C), N !== E && c(w, N)) : V & 16 ? Y & 16 ? _(
      E,
      N,
      w,
      T,
      S,
      C,
      F,
      L,
      k
    ) : Je(E, S, C, !0) : (V & 8 && c(w, ""), Y & 16 && Se(
      N,
      w,
      T,
      S,
      C,
      F,
      L,
      k
    ));
  }, b = (d, h, w, T, S, C, F, L, k) => {
    d = d || Dt, h = h || Dt;
    const E = d.length, V = h.length, N = Math.min(E, V);
    let U;
    for (U = 0; U < N; U++) {
      const Y = h[U] = k ? wt(h[U]) : it(h[U]);
      x(
        d[U],
        Y,
        w,
        null,
        S,
        C,
        F,
        L,
        k
      );
    }
    E > V ? Je(
      d,
      S,
      C,
      !0,
      !1,
      N
    ) : Se(
      h,
      w,
      T,
      S,
      C,
      F,
      L,
      k,
      N
    );
  }, _ = (d, h, w, T, S, C, F, L, k) => {
    let E = 0;
    const V = h.length;
    let N = d.length - 1, U = V - 1;
    for (; E <= N && E <= U;) {
      const Y = d[E], ne = h[E] = k ? wt(h[E]) : it(h[E]);
      if (Zt(Y, ne))
        x(
          Y,
          ne,
          w,
          null,
          S,
          C,
          F,
          L,
          k
        );
      else
        break;
      E++;
    }
    for (; E <= N && E <= U;) {
      const Y = d[N], ne = h[U] = k ? wt(h[U]) : it(h[U]);
      if (Zt(Y, ne))
        x(
          Y,
          ne,
          w,
          null,
          S,
          C,
          F,
          L,
          k
        );
      else
        break;
      N--, U--;
    }
    if (E > N) {
      if (E <= U) {
        const Y = U + 1, ne = Y < V ? h[Y].el : T;
        for (; E <= U;)
          x(
            null,
            h[E] = k ? wt(h[E]) : it(h[E]),
            w,
            ne,
            S,
            C,
            F,
            L,
            k
          ), E++;
      }
    } else if (E > U)
      for (; E <= N;)
        O(d[E], S, C, !0), E++;
    else {
      const Y = E, ne = E, fe = /* @__PURE__ */ new Map();
      for (E = ne; E <= U; E++) {
        const ke = h[E] = k ? wt(h[E]) : it(h[E]);
        ke.key != null && fe.set(ke.key, E);
      }
      let ae, Ce = 0;
      const Ee = U - ne + 1;
      let Be = !1, He = 0;
      const tt = new Array(Ee);
      for (E = 0; E < Ee; E++) tt[E] = 0;
      for (E = Y; E <= N; E++) {
        const ke = d[E];
        if (Ce >= Ee) {
          O(ke, S, C, !0);
          continue;
        }
        let Fe;
        if (ke.key != null)
          Fe = fe.get(ke.key);
        else
          for (ae = ne; ae <= U; ae++)
            if (tt[ae - ne] === 0 && Zt(ke, h[ae])) {
              Fe = ae;
              break;
            }
        Fe === void 0 ? O(ke, S, C, !0) : (tt[Fe - ne] = E + 1, Fe >= He ? He = Fe : Be = !0, x(
          ke,
          h[Fe],
          w,
          null,
          S,
          C,
          F,
          L,
          k
        ), Ce++);
      }
      const Yt = Be ? $a(tt) : Dt;
      for (ae = Yt.length - 1, E = Ee - 1; E >= 0; E--) {
        const ke = ne + E, Fe = h[ke], gs = h[ke + 1], Xt = ke + 1 < V ? (
          // #13559, fallback to el placeholder for unresolved async component
          gs.el || gs.placeholder
        ) : T;
        tt[E] === 0 ? x(
          null,
          Fe,
          w,
          Xt,
          S,
          C,
          F,
          L,
          k
        ) : Be && (ae < 0 || E !== Yt[ae] ? R(Fe, w, Xt, 2) : ae--);
      }
    }
  }, R = (d, h, w, T, S = null) => {
    const { el: C, type: F, transition: L, children: k, shapeFlag: E } = d;
    if (E & 6) {
      R(d.component.subTree, h, w, T);
      return;
    }
    if (E & 128) {
      d.suspense.move(h, w, T);
      return;
    }
    if (E & 64) {
      F.move(d, h, w, bt);
      return;
    }
    if (F === _e) {
      s(C, h, w);
      for (let N = 0; N < k.length; N++)
        R(k[N], h, w, T);
      s(d.anchor, h, w);
      return;
    }
    if (F === Ts) {
      W(d, h, w);
      return;
    }
    if (T !== 2 && E & 1 && L)
      if (T === 0)
        L.beforeEnter(C), s(C, h, w), Ue(() => L.enter(C), S);
      else {
        const { leave: N, delayLeave: U, afterLeave: Y } = L, ne = () => {
          d.ctx.isUnmounted ? r(C) : s(C, h, w);
        }, fe = () => {
          C._isLeaving && C[Xl](
            !0
            /* cancelled */
          ), N(C, () => {
            ne(), Y && Y();
          });
        };
        U ? U(C, ne, fe) : fe();
      }
    else
      s(C, h, w);
  }, O = (d, h, w, T = !1, S = !1) => {
    const {
      type: C,
      props: F,
      ref: L,
      children: k,
      dynamicChildren: E,
      shapeFlag: V,
      patchFlag: N,
      dirs: U,
      cacheIndex: Y
    } = d;
    if (N === -2 && (S = !1), L != null && (ht(), an(L, null, w, d, !0), pt()), Y != null && (h.renderCache[Y] = void 0), V & 256) {
      h.ctx.deactivate(d);
      return;
    }
    const ne = V & 1 && U, fe = !cn(d);
    let ae;
    if (fe && (ae = F && F.onVnodeBeforeUnmount) && nt(ae, h, d), V & 6)
      Ae(d.component, w, T);
    else {
      if (V & 128) {
        d.suspense.unmount(w, T);
        return;
      }
      ne && Ct(d, null, h, "beforeUnmount"), V & 64 ? d.type.remove(
        d,
        h,
        w,
        bt,
        T
      ) : E && // #5154
        // when v-once is used inside a block, setBlockTracking(-1) marks the
        // parent block with hasOnce: true
        // so that it doesn't take the fast path during unmount - otherwise
        // components nested in v-once are never unmounted.
        !E.hasOnce && // #1153: fast path should not be taken for non-stable (v-for) fragments
        (C !== _e || N > 0 && N & 64) ? Je(
          E,
          h,
          w,
          !1,
          !0
        ) : (C === _e && N & 384 || !S && V & 16) && Je(k, h, w), T && j(d);
    }
    (fe && (ae = F && F.onVnodeUnmounted) || ne) && Ue(() => {
      ae && nt(ae, h, d), ne && Ct(d, null, h, "unmounted");
    }, w);
  }, j = (d) => {
    const { type: h, el: w, anchor: T, transition: S } = d;
    if (h === _e) {
      re(w, T);
      return;
    }
    if (h === Ts) {
      I(d);
      return;
    }
    const C = () => {
      r(w), S && !S.persisted && S.afterLeave && S.afterLeave();
    };
    if (d.shapeFlag & 1 && S && !S.persisted) {
      const { leave: F, delayLeave: L } = S, k = () => F(w, C);
      L ? L(d.el, C, k) : k();
    } else
      C();
  }, re = (d, h) => {
    let w;
    for (; d !== h;)
      w = p(d), r(d), d = w;
    r(h);
  }, Ae = (d, h, w) => {
    const { bum: T, scope: S, job: C, subTree: F, um: L, m: k, a: E } = d;
    $r(k), $r(E), T && Ln(T), S.stop(), C && (C.flags |= 8, O(F, d, h, w)), L && Ue(L, h), Ue(() => {
      d.isUnmounted = !0;
    }, h);
  }, Je = (d, h, w, T = !1, S = !1, C = 0) => {
    for (let F = C; F < d.length; F++)
      O(d[F], h, w, T, S);
  }, et = (d) => {
    if (d.shapeFlag & 6)
      return et(d.component.subTree);
    if (d.shapeFlag & 128)
      return d.suspense.next();
    const h = p(d.anchor || d.el), w = h && h[Jl];
    return w ? p(w) : h;
  };
  let Mt = !1;
  const Jt = (d, h, w) => {
    d == null ? h._vnode && O(h._vnode, null, null, !0) : x(
      h._vnode || null,
      d,
      h,
      null,
      null,
      null,
      w
    ), h._vnode = d, Mt || (Mt = !0, kr(), Vi(), Mt = !1);
  }, bt = {
    p: x,
    um: O,
    m: R,
    r: j,
    mt: oe,
    mc: Se,
    pc: G,
    pbc: P,
    n: et,
    o: e
  };
  return {
    render: Jt,
    hydrate: void 0,
    createApp: ba(Jt)
  };
}
function Es({ type: e, props: t }, n) {
  return n === "svg" && e === "foreignObject" || n === "mathml" && e === "annotation-xml" && t && t.encoding && t.encoding.includes("html") ? void 0 : n;
}
function Et({ effect: e, job: t }, n) {
  n ? (e.flags |= 32, t.flags |= 4) : (e.flags &= -33, t.flags &= -5);
}
function Da(e, t) {
  return (!e || e && !e.pendingBranch) && t && !t.persisted;
}
function go(e, t, n = !1) {
  const s = e.children, r = t.children;
  if (K(s) && K(r))
    for (let i = 0; i < s.length; i++) {
      const o = s[i];
      let l = r[i];
      l.shapeFlag & 1 && !l.dynamicChildren && ((l.patchFlag <= 0 || l.patchFlag === 32) && (l = r[i] = wt(r[i]), l.el = o.el), !n && l.patchFlag !== -2 && go(o, l)), l.type === ls && // avoid cached text nodes retaining detached dom nodes
        l.patchFlag !== -1 && (l.el = o.el), l.type === xt && !l.el && (l.el = o.el);
    }
}
function $a(e) {
  const t = e.slice(), n = [0];
  let s, r, i, o, l;
  const a = e.length;
  for (s = 0; s < a; s++) {
    const u = e[s];
    if (u !== 0) {
      if (r = n[n.length - 1], e[r] < u) {
        t[s] = r, n.push(s);
        continue;
      }
      for (i = 0, o = n.length - 1; i < o;)
        l = i + o >> 1, e[n[l]] < u ? i = l + 1 : o = l;
      u < e[n[i]] && (i > 0 && (t[s] = n[i - 1]), n[i] = s);
    }
  }
  for (i = n.length, o = n[i - 1]; i-- > 0;)
    n[i] = o, o = t[o];
  return n;
}
function mo(e) {
  const t = e.subTree.component;
  if (t)
    return t.asyncDep && !t.asyncResolved ? t : mo(t);
}
function $r(e) {
  if (e)
    for (let t = 0; t < e.length; t++)
      e[t].flags |= 8;
}
const bo = (e) => e.__isSuspense;
function ja(e, t) {
  t && t.pendingBranch ? K(e) ? t.effects.push(...e) : t.effects.push(e) : Kl(e);
}
const _e = Symbol.for("v-fgt"), ls = Symbol.for("v-txt"), xt = Symbol.for("v-cmt"), Ts = Symbol.for("v-stc"), fn = [];
let Ve = null;
function $(e = !1) {
  fn.push(Ve = e ? null : []);
}
function Ba() {
  fn.pop(), Ve = fn[fn.length - 1] || null;
}
let gn = 1;
function Wn(e, t = !1) {
  gn += e, e < 0 && Ve && t && (Ve.hasOnce = !0);
}
function yo(e) {
  return e.dynamicChildren = gn > 0 ? Ve || Dt : null, Ba(), gn > 0 && Ve && Ve.push(e), e;
}
function B(e, t, n, s, r, i) {
  return yo(
    M(
      e,
      t,
      n,
      s,
      r,
      i,
      !0
    )
  );
}
function mn(e, t, n, s, r) {
  return yo(
    Te(
      e,
      t,
      n,
      s,
      r,
      !0
    )
  );
}
function zn(e) {
  return e ? e.__v_isVNode === !0 : !1;
}
function Zt(e, t) {
  return e.type === t.type && e.key === t.key;
}
const _o = ({ key: e }) => e ?? null, Mn = ({
  ref: e,
  ref_key: t,
  ref_for: n
}) => (typeof e == "number" && (e = "" + e), e != null ? we(e) || Le(e) || Z(e) ? { i: qe, r: e, k: t, f: !!n } : e : null);
function M(e, t = null, n = null, s = 0, r = null, i = e === _e ? 0 : 1, o = !1, l = !1) {
  const a = {
    __v_isVNode: !0,
    __v_skip: !0,
    type: e,
    props: t,
    key: t && _o(t),
    ref: t && Mn(t),
    scopeId: zi,
    slotScopeIds: null,
    children: n,
    component: null,
    suspense: null,
    ssContent: null,
    ssFallback: null,
    dirs: null,
    transition: null,
    el: null,
    anchor: null,
    target: null,
    targetStart: null,
    targetAnchor: null,
    staticCount: 0,
    shapeFlag: i,
    patchFlag: s,
    dynamicProps: r,
    dynamicChildren: null,
    appContext: null,
    ctx: qe
  };
  return l ? (gr(a, n), i & 128 && e.normalize(a)) : n && (a.shapeFlag |= we(n) ? 8 : 16), gn > 0 && // avoid a block node from tracking itself
    !o && // has current parent block
    Ve && // presence of a patch flag indicates this node needs patching on updates.
    // component nodes also should always be patched, because even if the
    // component doesn't need to update, it needs to persist the instance on to
    // the next vnode so that it can be properly unmounted later.
    (a.patchFlag > 0 || i & 6) && // the EVENTS flag is only for hydration and if it is the only flag, the
    // vnode should not be considered dynamic due to handler caching.
    a.patchFlag !== 32 && Ve.push(a), a;
}
const Te = Ha;
function Ha(e, t = null, n = null, s = 0, r = null, i = !1) {
  if ((!e || e === Xi) && (e = xt), zn(e)) {
    const l = qt(
      e,
      t,
      !0
      /* mergeRef: true */
    );
    return n && gr(l, n), gn > 0 && !i && Ve && (l.shapeFlag & 6 ? Ve[Ve.indexOf(e)] = l : Ve.push(l)), l.patchFlag = -2, l;
  }
  if (Za(e) && (e = e.__vccOpts), t) {
    t = Ua(t);
    let { class: l, style: a } = t;
    l && !we(l) && (t.class = Lt(l)), ge(a) && (ar(a) && !K(a) && (a = xe({}, a)), t.style = Zs(a));
  }
  const o = we(e) ? 1 : bo(e) ? 128 : Yl(e) ? 64 : ge(e) ? 4 : Z(e) ? 2 : 0;
  return M(
    e,
    t,
    n,
    s,
    r,
    o,
    i,
    !0
  );
}
function Ua(e) {
  return e ? ar(e) || ao(e) ? xe({}, e) : e : null;
}
function qt(e, t, n = !1, s = !1) {
  const { props: r, ref: i, patchFlag: o, children: l, transition: a } = e, u = t ? qa(r || {}, t) : r, c = {
    __v_isVNode: !0,
    __v_skip: !0,
    type: e.type,
    props: u,
    key: u && _o(u),
    ref: t && t.ref ? (
      // #2078 in the case of <component :is="vnode" ref="extra"/>
      // if the vnode itself already has a ref, cloneVNode will need to merge
      // the refs so the single vnode can be set on multiple refs
      n && i ? K(i) ? i.concat(Mn(t)) : [i, Mn(t)] : Mn(t)
    ) : i,
    scopeId: e.scopeId,
    slotScopeIds: e.slotScopeIds,
    children: l,
    target: e.target,
    targetStart: e.targetStart,
    targetAnchor: e.targetAnchor,
    staticCount: e.staticCount,
    shapeFlag: e.shapeFlag,
    // if the vnode is cloned with extra props, we can no longer assume its
    // existing patch flag to be reliable and need to add the FULL_PROPS flag.
    // note: preserve flag for fragments since they use the flag for children
    // fast paths only.
    patchFlag: t && e.type !== _e ? o === -1 ? 16 : o | 16 : o,
    dynamicProps: e.dynamicProps,
    dynamicChildren: e.dynamicChildren,
    appContext: e.appContext,
    dirs: e.dirs,
    transition: a,
    // These should technically only be non-null on mounted VNodes. However,
    // they *should* be copied for kept-alive vnodes. So we just always copy
    // them since them being non-null during a mount doesn't affect the logic as
    // they will simply be overwritten.
    component: e.component,
    suspense: e.suspense,
    ssContent: e.ssContent && qt(e.ssContent),
    ssFallback: e.ssFallback && qt(e.ssFallback),
    placeholder: e.placeholder,
    el: e.el,
    anchor: e.anchor,
    ctx: e.ctx,
    ce: e.ce
  };
  return a && s && ur(
    c,
    a.clone(c)
  ), c;
}
function St(e = " ", t = 0) {
  return Te(ls, null, e, t);
}
function be(e = "", t = !1) {
  return t ? ($(), mn(xt, null, e)) : Te(xt, null, e);
}
function it(e) {
  return e == null || typeof e == "boolean" ? Te(xt) : K(e) ? Te(
    _e,
    null,
    // #3666, avoid reference pollution when reusing vnode
    e.slice()
  ) : zn(e) ? wt(e) : Te(ls, null, String(e));
}
function wt(e) {
  return e.el === null && e.patchFlag !== -1 || e.memo ? e : qt(e);
}
function gr(e, t) {
  let n = 0;
  const { shapeFlag: s } = e;
  if (t == null)
    t = null;
  else if (K(t))
    n = 16;
  else if (typeof t == "object")
    if (s & 65) {
      const r = t.default;
      r && (r._c && (r._d = !1), gr(e, r()), r._c && (r._d = !0));
      return;
    } else {
      n = 32;
      const r = t._;
      !r && !ao(t) ? t._ctx = qe : r === 3 && qe && (qe.slots._ === 1 ? t._ = 1 : (t._ = 2, e.patchFlag |= 1024));
    }
  else Z(t) ? (t = { default: t, _ctx: qe }, n = 32) : (t = String(t), s & 64 ? (n = 16, t = [St(t)]) : n = 8);
  e.children = t, e.shapeFlag |= n;
}
function qa(...e) {
  const t = {};
  for (let n = 0; n < e.length; n++) {
    const s = e[n];
    for (const r in s)
      if (r === "class")
        t.class !== s.class && (t.class = Lt([t.class, s.class]));
      else if (r === "style")
        t.style = Zs([t.style, s.style]);
      else if (Xn(r)) {
        const i = t[r], o = s[r];
        o && i !== o && !(K(i) && i.includes(o)) && (t[r] = i ? [].concat(i, o) : o);
      } else r !== "" && (t[r] = s[r]);
  }
  return t;
}
function nt(e, t, n, s = null) {
  at(e, t, 7, [
    n,
    s
  ]);
}
const Va = no();
let Wa = 0;
function za(e, t, n) {
  const s = e.type, r = (t ? t.appContext : e.appContext) || Va, i = {
    uid: Wa++,
    vnode: e,
    type: s,
    parent: t,
    appContext: r,
    root: null,
    // to be immediately set
    next: null,
    subTree: null,
    // will be set synchronously right after creation
    effect: null,
    update: null,
    // will be set synchronously right after creation
    job: null,
    scope: new ml(
      !0
      /* detached */
    ),
    render: null,
    proxy: null,
    exposed: null,
    exposeProxy: null,
    withProxy: null,
    provides: t ? t.provides : Object.create(r.provides),
    ids: t ? t.ids : ["", 0, 0],
    accessCache: null,
    renderCache: [],
    // local resolved assets
    components: null,
    directives: null,
    // resolved props and emits options
    propsOptions: uo(s, r),
    emitsOptions: io(s, r),
    // emit
    emit: null,
    // to be set immediately
    emitted: null,
    // props default value
    propsDefaults: pe,
    // inheritAttrs
    inheritAttrs: s.inheritAttrs,
    // state
    ctx: pe,
    data: pe,
    props: pe,
    attrs: pe,
    slots: pe,
    refs: pe,
    setupState: pe,
    setupContext: null,
    // suspense related
    suspense: n,
    suspenseId: n ? n.pendingId : 0,
    asyncDep: null,
    asyncResolved: !1,
    // lifecycle hooks
    // not using enums here because it results in computed properties
    isMounted: !1,
    isUnmounted: !1,
    isDeactivated: !1,
    bc: null,
    c: null,
    bm: null,
    m: null,
    bu: null,
    u: null,
    um: null,
    bum: null,
    da: null,
    a: null,
    rtg: null,
    rtc: null,
    ec: null,
    sp: null
  };
  return i.ctx = { _: i }, i.root = t ? t.root : i, i.emit = Sa.bind(null, i), e.ce && e.ce(i), i;
}
let Oe = null;
const Ka = () => Oe || qe;
let Kn, Us;
{
  const e = ts(), t = (n, s) => {
    let r;
    return (r = e[n]) || (r = e[n] = []), r.push(s), (i) => {
      r.length > 1 ? r.forEach((o) => o(i)) : r[0](i);
    };
  };
  Kn = t(
    "__VUE_INSTANCE_SETTERS__",
    (n) => Oe = n
  ), Us = t(
    "__VUE_SSR_SETTERS__",
    (n) => bn = n
  );
}
const xn = (e) => {
  const t = Oe;
  return Kn(e), e.scope.on(), () => {
    e.scope.off(), Kn(t);
  };
}, jr = () => {
  Oe && Oe.scope.off(), Kn(null);
};
function wo(e) {
  return e.vnode.shapeFlag & 4;
}
let bn = !1;
function Ga(e, t = !1, n = !1) {
  t && Us(t);
  const { props: s, children: r } = e.vnode, i = wo(e);
  Ra(e, s, i, t), Fa(e, r, n || t);
  const o = i ? Ja(e, t) : void 0;
  return t && Us(!1), o;
}
function Ja(e, t) {
  const n = e.type;
  e.accessCache = /* @__PURE__ */ Object.create(null), e.proxy = new Proxy(e.ctx, ua);
  const { setup: s } = n;
  if (s) {
    ht();
    const r = e.setupContext = s.length > 1 ? Xa(e) : null, i = xn(e), o = vn(
      s,
      e,
      0,
      [
        e.props,
        r
      ]
    ), l = _i(o);
    if (pt(), i(), (l || e.sp) && !cn(e) && Ki(e), l) {
      if (o.then(jr, jr), t)
        return o.then((a) => {
          Br(e, a);
        }).catch((a) => {
          ss(a, e, 0);
        });
      e.asyncDep = o;
    } else
      Br(e, o);
  } else
    vo(e);
}
function Br(e, t, n) {
  Z(t) ? e.type.__ssrInlineRender ? e.ssrRender = t : e.render = t : ge(t) && (e.setupState = Hi(t)), vo(e);
}
function vo(e, t, n) {
  const s = e.type;
  e.render || (e.render = s.render || ot);
  {
    const r = xn(e);
    ht();
    try {
      fa(e);
    } finally {
      pt(), r();
    }
  }
}
const Ya = {
  get(e, t) {
    return Re(e, "get", ""), e[t];
  }
};
function Xa(e) {
  const t = (n) => {
    e.exposed = n || {};
  };
  return {
    attrs: new Proxy(e.attrs, Ya),
    slots: e.slots,
    emit: e.emit,
    expose: t
  };
}
function as(e) {
  return e.exposed ? e.exposeProxy || (e.exposeProxy = new Proxy(Hi(Dl(e.exposed)), {
    get(t, n) {
      if (n in t)
        return t[n];
      if (n in un)
        return un[n](e);
    },
    has(t, n) {
      return n in t || n in un;
    }
  })) : e.proxy;
}
function Qa(e, t = !0) {
  return Z(e) ? e.displayName || e.name : e.name || t && e.__name;
}
function Za(e) {
  return Z(e) && "__vccOpts" in e;
}
const We = (e, t) => Ul(e, t, bn);
function en(e, t, n) {
  try {
    Wn(-1);
    const s = arguments.length;
    return s === 2 ? ge(t) && !K(t) ? zn(t) ? Te(e, null, [t]) : Te(e, t) : Te(e, null, t) : (s > 3 ? n = Array.prototype.slice.call(arguments, 2) : s === 3 && zn(n) && (n = [n]), Te(e, t, n));
  } finally {
    Wn(1);
  }
}
const ec = "3.5.25";
/**
* @vue/runtime-dom v3.5.25
* (c) 2018-present Yuxi (Evan) You and Vue contributors
* @license MIT
**/
let qs;
const Hr = typeof window < "u" && window.trustedTypes;
if (Hr)
  try {
    qs = /* @__PURE__ */ Hr.createPolicy("vue", {
      createHTML: (e) => e
    });
  } catch {
  }
const xo = qs ? (e) => qs.createHTML(e) : (e) => e, tc = "http://www.w3.org/2000/svg", nc = "http://www.w3.org/1998/Math/MathML", ut = typeof document < "u" ? document : null, Ur = ut && /* @__PURE__ */ ut.createElement("template"), sc = {
  insert: (e, t, n) => {
    t.insertBefore(e, n || null);
  },
  remove: (e) => {
    const t = e.parentNode;
    t && t.removeChild(e);
  },
  createElement: (e, t, n, s) => {
    const r = t === "svg" ? ut.createElementNS(tc, e) : t === "mathml" ? ut.createElementNS(nc, e) : n ? ut.createElement(e, { is: n }) : ut.createElement(e);
    return e === "select" && s && s.multiple != null && r.setAttribute("multiple", s.multiple), r;
  },
  createText: (e) => ut.createTextNode(e),
  createComment: (e) => ut.createComment(e),
  setText: (e, t) => {
    e.nodeValue = t;
  },
  setElementText: (e, t) => {
    e.textContent = t;
  },
  parentNode: (e) => e.parentNode,
  nextSibling: (e) => e.nextSibling,
  querySelector: (e) => ut.querySelector(e),
  setScopeId(e, t) {
    e.setAttribute(t, "");
  },
  // __UNSAFE__
  // Reason: innerHTML.
  // Static content here can only come from compiled templates.
  // As long as the user only uses trusted templates, this is safe.
  insertStaticContent(e, t, n, s, r, i) {
    const o = n ? n.previousSibling : t.lastChild;
    if (r && (r === i || r.nextSibling))
      for (; t.insertBefore(r.cloneNode(!0), n), !(r === i || !(r = r.nextSibling));)
        ;
    else {
      Ur.innerHTML = xo(
        s === "svg" ? `<svg>${e}</svg>` : s === "mathml" ? `<math>${e}</math>` : e
      );
      const l = Ur.content;
      if (s === "svg" || s === "mathml") {
        const a = l.firstChild;
        for (; a.firstChild;)
          l.appendChild(a.firstChild);
        l.removeChild(a);
      }
      t.insertBefore(l, n);
    }
    return [
      // first
      o ? o.nextSibling : t.firstChild,
      // last
      n ? n.previousSibling : t.lastChild
    ];
  }
}, rc = Symbol("_vtc");
function ic(e, t, n) {
  const s = e[rc];
  s && (t = (t ? [t, ...s] : [...s]).join(" ")), t == null ? e.removeAttribute("class") : n ? e.setAttribute("class", t) : e.className = t;
}
const Gn = Symbol("_vod"), So = Symbol("_vsh"), oc = {
  // used for prop mismatch check during hydration
  name: "show",
  beforeMount(e, { value: t }, { transition: n }) {
    e[Gn] = e.style.display === "none" ? "" : e.style.display, n && t ? n.beforeEnter(e) : tn(e, t);
  },
  mounted(e, { value: t }, { transition: n }) {
    n && t && n.enter(e);
  },
  updated(e, { value: t, oldValue: n }, { transition: s }) {
    !t != !n && (s ? t ? (s.beforeEnter(e), tn(e, !0), s.enter(e)) : s.leave(e, () => {
      tn(e, !1);
    }) : tn(e, t));
  },
  beforeUnmount(e, { value: t }) {
    tn(e, t);
  }
};
function tn(e, t) {
  e.style.display = t ? e[Gn] : "none", e[So] = !t;
}
const lc = Symbol(""), ac = /(?:^|;)\s*display\s*:/;
function cc(e, t, n) {
  const s = e.style, r = we(n);
  let i = !1;
  if (n && !r) {
    if (t)
      if (we(t))
        for (const o of t.split(";")) {
          const l = o.slice(0, o.indexOf(":")).trim();
          n[l] == null && In(s, l, "");
        }
      else
        for (const o in t)
          n[o] == null && In(s, o, "");
    for (const o in n)
      o === "display" && (i = !0), In(s, o, n[o]);
  } else if (r) {
    if (t !== n) {
      const o = s[lc];
      o && (n += ";" + o), s.cssText = n, i = ac.test(n);
    }
  } else t && e.removeAttribute("style");
  Gn in e && (e[Gn] = i ? s.display : "", e[So] && (s.display = "none"));
}
const qr = /\s*!important$/;
function In(e, t, n) {
  if (K(n))
    n.forEach((s) => In(e, t, s));
  else if (n == null && (n = ""), t.startsWith("--"))
    e.setProperty(t, n);
  else {
    const s = uc(e, t);
    qr.test(n) ? e.setProperty(
      ze(s),
      n.replace(qr, ""),
      "important"
    ) : e[s] = n;
  }
}
const Vr = ["Webkit", "Moz", "ms"], As = {};
function uc(e, t) {
  const n = As[t];
  if (n)
    return n;
  let s = Ne(t);
  if (s !== "filter" && s in e)
    return As[t] = s;
  s = es(s);
  for (let r = 0; r < Vr.length; r++) {
    const i = Vr[r] + s;
    if (i in e)
      return As[t] = i;
  }
  return t;
}
const Wr = "http://www.w3.org/1999/xlink";
function zr(e, t, n, s, r, i = pl(t)) {
  s && t.startsWith("xlink:") ? n == null ? e.removeAttributeNS(Wr, t.slice(6, t.length)) : e.setAttributeNS(Wr, t, n) : n == null || i && !Si(n) ? e.removeAttribute(t) : e.setAttribute(
    t,
    i ? "" : lt(n) ? String(n) : n
  );
}
function Kr(e, t, n, s, r) {
  if (t === "innerHTML" || t === "textContent") {
    n != null && (e[t] = t === "innerHTML" ? xo(n) : n);
    return;
  }
  const i = e.tagName;
  if (t === "value" && i !== "PROGRESS" && // custom elements may use _value internally
    !i.includes("-")) {
    const l = i === "OPTION" ? e.getAttribute("value") || "" : e.value, a = n == null ? (
      // #11647: value should be set as empty string for null and undefined,
      // but <input type="checkbox"> should be set as 'on'.
      e.type === "checkbox" ? "on" : ""
    ) : String(n);
    (l !== a || !("_value" in e)) && (e.value = a), n == null && e.removeAttribute(t), e._value = n;
    return;
  }
  let o = !1;
  if (n === "" || n == null) {
    const l = typeof e[t];
    l === "boolean" ? n = Si(n) : n == null && l === "string" ? (n = "", o = !0) : l === "number" && (n = 0, o = !0);
  }
  try {
    e[t] = n;
  } catch {
  }
  o && e.removeAttribute(r || t);
}
function mr(e, t, n, s) {
  e.addEventListener(t, n, s);
}
function fc(e, t, n, s) {
  e.removeEventListener(t, n, s);
}
const Gr = Symbol("_vei");
function dc(e, t, n, s, r = null) {
  const i = e[Gr] || (e[Gr] = {}), o = i[t];
  if (s && o)
    o.value = s;
  else {
    const [l, a] = hc(t);
    if (s) {
      const u = i[t] = mc(
        s,
        r
      );
      mr(e, l, u, a);
    } else o && (fc(e, l, o, a), i[t] = void 0);
  }
}
const Jr = /(?:Once|Passive|Capture)$/;
function hc(e) {
  let t;
  if (Jr.test(e)) {
    t = {};
    let s;
    for (; s = e.match(Jr);)
      e = e.slice(0, e.length - s[0].length), t[s[0].toLowerCase()] = !0;
  }
  return [e[2] === ":" ? e.slice(3) : ze(e.slice(2)), t];
}
let ks = 0;
const pc = /* @__PURE__ */ Promise.resolve(), gc = () => ks || (pc.then(() => ks = 0), ks = Date.now());
function mc(e, t) {
  const n = (s) => {
    if (!s._vts)
      s._vts = Date.now();
    else if (s._vts <= n.attached)
      return;
    at(
      bc(s, n.value),
      t,
      5,
      [s]
    );
  };
  return n.value = e, n.attached = gc(), n;
}
function bc(e, t) {
  if (K(t)) {
    const n = e.stopImmediatePropagation;
    return e.stopImmediatePropagation = () => {
      n.call(e), e._stopped = !0;
    }, t.map(
      (s) => (r) => !r._stopped && s && s(r)
    );
  } else
    return t;
}
const Yr = (e) => e.charCodeAt(0) === 111 && e.charCodeAt(1) === 110 && // lowercase letter
  e.charCodeAt(2) > 96 && e.charCodeAt(2) < 123, yc = (e, t, n, s, r, i) => {
    const o = r === "svg";
    t === "class" ? ic(e, s, o) : t === "style" ? cc(e, n, s) : Xn(t) ? Ys(t) || dc(e, t, n, s, i) : (t[0] === "." ? (t = t.slice(1), !0) : t[0] === "^" ? (t = t.slice(1), !1) : _c(e, t, s, o)) ? (Kr(e, t, s), !e.tagName.includes("-") && (t === "value" || t === "checked" || t === "selected") && zr(e, t, s, o, i, t !== "value")) : /* #11081 force set props for possible async custom element */ e._isVueCE && (/[A-Z]/.test(t) || !we(s)) ? Kr(e, Ne(t), s, i, t) : (t === "true-value" ? e._trueValue = s : t === "false-value" && (e._falseValue = s), zr(e, t, s, o));
  };
function _c(e, t, n, s) {
  if (s)
    return !!(t === "innerHTML" || t === "textContent" || t in e && Yr(t) && Z(n));
  if (t === "spellcheck" || t === "draggable" || t === "translate" || t === "autocorrect" || t === "sandbox" && e.tagName === "IFRAME" || t === "form" || t === "list" && e.tagName === "INPUT" || t === "type" && e.tagName === "TEXTAREA")
    return !1;
  if (t === "width" || t === "height") {
    const r = e.tagName;
    if (r === "IMG" || r === "VIDEO" || r === "CANVAS" || r === "SOURCE")
      return !1;
  }
  return Yr(t) && we(n) ? !1 : t in e;
}
const Xr = {};
// @__NO_SIDE_EFFECTS__
function Co(e, t, n) {
  let s = /* @__PURE__ */ Ql(e, t);
  Qn(s) && (s = xe({}, s, t));
  class r extends br {
    constructor(o) {
      super(s, o, n);
    }
  }
  return r.def = s, r;
}
const wc = typeof HTMLElement < "u" ? HTMLElement : class {
};
class br extends wc {
  constructor(t, n = {}, s = ti) {
    super(), this._def = t, this._props = n, this._createApp = s, this._isVueCE = !0, this._instance = null, this._app = null, this._nonce = this._def.nonce, this._connected = !1, this._resolved = !1, this._patching = !1, this._dirty = !1, this._numberProps = null, this._styleChildren = /* @__PURE__ */ new WeakSet(), this._ob = null, this.shadowRoot && s !== ti ? this._root = this.shadowRoot : t.shadowRoot !== !1 ? (this.attachShadow(
      xe({}, t.shadowRootOptions, {
        mode: "open"
      })
    ), this._root = this.shadowRoot) : this._root = this;
  }
  connectedCallback() {
    if (!this.isConnected) return;
    !this.shadowRoot && !this._resolved && this._parseSlots(), this._connected = !0;
    let t = this;
    for (; t = t && (t.parentNode || t.host);)
      if (t instanceof br) {
        this._parent = t;
        break;
      }
    this._instance || (this._resolved ? this._mount(this._def) : t && t._pendingResolve ? this._pendingResolve = t._pendingResolve.then(() => {
      this._pendingResolve = void 0, this._resolveDef();
    }) : this._resolveDef());
  }
  _setParent(t = this._parent) {
    t && (this._instance.parent = t._instance, this._inheritParentContext(t));
  }
  _inheritParentContext(t = this._parent) {
    t && this._app && Object.setPrototypeOf(
      this._app._context.provides,
      t._instance.provides
    );
  }
  disconnectedCallback() {
    this._connected = !1, rs(() => {
      this._connected || (this._ob && (this._ob.disconnect(), this._ob = null), this._app && this._app.unmount(), this._instance && (this._instance.ce = void 0), this._app = this._instance = null, this._teleportTargets && (this._teleportTargets.clear(), this._teleportTargets = void 0));
    });
  }
  _processMutations(t) {
    for (const n of t)
      this._setAttr(n.attributeName);
  }
  /**
   * resolve inner component definition (handle possible async component)
   */
  _resolveDef() {
    if (this._pendingResolve)
      return;
    for (let s = 0; s < this.attributes.length; s++)
      this._setAttr(this.attributes[s].name);
    this._ob = new MutationObserver(this._processMutations.bind(this)), this._ob.observe(this, { attributes: !0 });
    const t = (s, r = !1) => {
      this._resolved = !0, this._pendingResolve = void 0;
      const { props: i, styles: o } = s;
      let l;
      if (i && !K(i))
        for (const a in i) {
          const u = i[a];
          (u === Number || u && u.type === Number) && (a in this._props && (this._props[a] = Cr(this._props[a])), (l || (l = /* @__PURE__ */ Object.create(null)))[Ne(a)] = !0);
        }
      this._numberProps = l, this._resolveProps(s), this.shadowRoot && this._applyStyles(o), this._mount(s);
    }, n = this._def.__asyncLoader;
    n ? this._pendingResolve = n().then((s) => {
      s.configureApp = this._def.configureApp, t(this._def = s, !0);
    }) : t(this._def);
  }
  _mount(t) {
    this._app = this._createApp(t), this._inheritParentContext(), t.configureApp && t.configureApp(this._app), this._app._ceVNode = this._createVNode(), this._app.mount(this._root);
    const n = this._instance && this._instance.exposed;
    if (n)
      for (const s in n)
        ce(this, s) || Object.defineProperty(this, s, {
          // unwrap ref to be consistent with public instance behavior
          get: () => Bi(n[s])
        });
  }
  _resolveProps(t) {
    const { props: n } = t, s = K(n) ? n : Object.keys(n || {});
    for (const r of Object.keys(this))
      r[0] !== "_" && s.includes(r) && this._setProp(r, this[r]);
    for (const r of s.map(Ne))
      Object.defineProperty(this, r, {
        get() {
          return this._getProp(r);
        },
        set(i) {
          this._setProp(r, i, !0, !this._patching);
        }
      });
  }
  _setAttr(t) {
    if (t.startsWith("data-v-")) return;
    const n = this.hasAttribute(t);
    let s = n ? this.getAttribute(t) : Xr;
    const r = Ne(t);
    n && this._numberProps && this._numberProps[r] && (s = Cr(s)), this._setProp(r, s, !1, !0);
  }
  /**
   * @internal
   */
  _getProp(t) {
    return this._props[t];
  }
  /**
   * @internal
   */
  _setProp(t, n, s = !0, r = !1) {
    if (n !== this._props[t] && (this._dirty = !0, n === Xr ? delete this._props[t] : (this._props[t] = n, t === "key" && this._app && (this._app._ceVNode.key = n)), r && this._instance && this._update(), s)) {
      const i = this._ob;
      i && (this._processMutations(i.takeRecords()), i.disconnect()), n === !0 ? this.setAttribute(ze(t), "") : typeof n == "string" || typeof n == "number" ? this.setAttribute(ze(t), n + "") : n || this.removeAttribute(ze(t)), i && i.observe(this, { attributes: !0 });
    }
  }
  _update() {
    const t = this._createVNode();
    this._app && (t.appContext = this._app._context), Tc(t, this._root);
  }
  _createVNode() {
    const t = {};
    this.shadowRoot || (t.onVnodeMounted = t.onVnodeUpdated = this._renderSlots.bind(this));
    const n = Te(this._def, xe(t, this._props));
    return this._instance || (n.ce = (s) => {
      this._instance = s, s.ce = this, s.isCE = !0;
      const r = (i, o) => {
        this.dispatchEvent(
          new CustomEvent(
            i,
            Qn(o[0]) ? xe({ detail: o }, o[0]) : { detail: o }
          )
        );
      };
      s.emit = (i, ...o) => {
        r(i, o), ze(i) !== i && r(ze(i), o);
      }, this._setParent();
    }), n;
  }
  _applyStyles(t, n) {
    if (!t) return;
    if (n) {
      if (n === this._def || this._styleChildren.has(n))
        return;
      this._styleChildren.add(n);
    }
    const s = this._nonce;
    for (let r = t.length - 1; r >= 0; r--) {
      const i = document.createElement("style");
      s && i.setAttribute("nonce", s), i.textContent = t[r], this.shadowRoot.prepend(i);
    }
  }
  /**
   * Only called when shadowRoot is false
   */
  _parseSlots() {
    const t = this._slots = {};
    let n;
    for (; n = this.firstChild;) {
      const s = n.nodeType === 1 && n.getAttribute("slot") || "default";
      (t[s] || (t[s] = [])).push(n), this.removeChild(n);
    }
  }
  /**
   * Only called when shadowRoot is false
   */
  _renderSlots() {
    const t = this._getSlots(), n = this._instance.type.__scopeId;
    for (let s = 0; s < t.length; s++) {
      const r = t[s], i = r.getAttribute("name") || "default", o = this._slots[i], l = r.parentNode;
      if (o)
        for (const a of o) {
          if (n && a.nodeType === 1) {
            const u = n + "-s", c = document.createTreeWalker(a, 1);
            a.setAttribute(u, "");
            let f;
            for (; f = c.nextNode();)
              f.setAttribute(u, "");
          }
          l.insertBefore(a, r);
        }
      else
        for (; r.firstChild;) l.insertBefore(r.firstChild, r);
      l.removeChild(r);
    }
  }
  /**
   * @internal
   */
  _getSlots() {
    const t = [this];
    this._teleportTargets && t.push(...this._teleportTargets);
    const n = /* @__PURE__ */ new Set();
    for (const s of t) {
      const r = s.querySelectorAll("slot");
      for (let i = 0; i < r.length; i++)
        n.add(r[i]);
    }
    return Array.from(n);
  }
  /**
   * @internal
   */
  _injectChildStyle(t) {
    this._applyStyles(t.styles, t);
  }
  /**
   * @internal
   */
  _beginPatch() {
    this._patching = !0, this._dirty = !1;
  }
  /**
   * @internal
   */
  _endPatch() {
    this._patching = !1, this._dirty && this._instance && this._update();
  }
  /**
   * @internal
   */
  _removeChildStyle(t) {
  }
}
const Jn = (e) => {
  const t = e.props["onUpdate:modelValue"] || !1;
  return K(t) ? (n) => Ln(t, n) : t;
}, Ht = Symbol("_assign"), vc = {
  // #4096 array checkboxes need to be deep traversed
  deep: !0,
  created(e, t, n) {
    e[Ht] = Jn(n), mr(e, "change", () => {
      const s = e._modelValue, r = yn(e), i = e.checked, o = e[Ht];
      if (K(s)) {
        const l = er(s, r), a = l !== -1;
        if (i && !a)
          o(s.concat(r));
        else if (!i && a) {
          const u = [...s];
          u.splice(l, 1), o(u);
        }
      } else if (Wt(s)) {
        const l = new Set(s);
        i ? l.add(r) : l.delete(r), o(l);
      } else
        o(Eo(e, i));
    });
  },
  // set initial checked on mount to wait for true-value/false-value
  mounted: Qr,
  beforeUpdate(e, t, n) {
    e[Ht] = Jn(n), Qr(e, t, n);
  }
};
function Qr(e, { value: t, oldValue: n }, s) {
  e._modelValue = t;
  let r;
  if (K(t))
    r = er(t, s.props.value) > -1;
  else if (Wt(t))
    r = t.has(s.props.value);
  else {
    if (t === n) return;
    r = wn(t, Eo(e, !0));
  }
  e.checked !== r && (e.checked = r);
}
const xc = {
  // <select multiple> value need to be deep traversed
  deep: !0,
  created(e, { value: t, modifiers: { number: n } }, s) {
    const r = Wt(t);
    mr(e, "change", () => {
      const i = Array.prototype.filter.call(e.options, (o) => o.selected).map(
        (o) => n ? xi(yn(o)) : yn(o)
      );
      e[Ht](
        e.multiple ? r ? new Set(i) : i : i[0]
      ), e._assigning = !0, rs(() => {
        e._assigning = !1;
      });
    }), e[Ht] = Jn(s);
  },
  // set value in mounted & updated because <select> relies on its children
  // <option>s.
  mounted(e, { value: t }) {
    Zr(e, t);
  },
  beforeUpdate(e, t, n) {
    e[Ht] = Jn(n);
  },
  updated(e, { value: t }) {
    e._assigning || Zr(e, t);
  }
};
function Zr(e, t) {
  const n = e.multiple, s = K(t);
  if (!(n && !s && !Wt(t))) {
    for (let r = 0, i = e.options.length; r < i; r++) {
      const o = e.options[r], l = yn(o);
      if (n)
        if (s) {
          const a = typeof l;
          a === "string" || a === "number" ? o.selected = t.some((u) => String(u) === String(l)) : o.selected = er(t, l) > -1;
        } else
          o.selected = t.has(l);
      else if (wn(yn(o), t)) {
        e.selectedIndex !== r && (e.selectedIndex = r);
        return;
      }
    }
    !n && e.selectedIndex !== -1 && (e.selectedIndex = -1);
  }
}
function yn(e) {
  return "_value" in e ? e._value : e.value;
}
function Eo(e, t) {
  const n = t ? "_trueValue" : "_falseValue";
  return n in e ? e[n] : t;
}
const Sc = ["ctrl", "shift", "alt", "meta"], Cc = {
  stop: (e) => e.stopPropagation(),
  prevent: (e) => e.preventDefault(),
  self: (e) => e.target !== e.currentTarget,
  ctrl: (e) => !e.ctrlKey,
  shift: (e) => !e.shiftKey,
  alt: (e) => !e.altKey,
  meta: (e) => !e.metaKey,
  left: (e) => "button" in e && e.button !== 0,
  middle: (e) => "button" in e && e.button !== 1,
  right: (e) => "button" in e && e.button !== 2,
  exact: (e, t) => Sc.some((n) => e[`${n}Key`] && !t.includes(n))
}, To = (e, t) => {
  const n = e._withMods || (e._withMods = {}), s = t.join(".");
  return n[s] || (n[s] = (r, ...i) => {
    for (let o = 0; o < t.length; o++) {
      const l = Cc[t[o]];
      if (l && l(r, t)) return;
    }
    return e(r, ...i);
  });
}, Ec = /* @__PURE__ */ xe({ patchProp: yc }, sc);
let ei;
function Ao() {
  return ei || (ei = Ia(Ec));
}
const Tc = (...e) => {
  Ao().render(...e);
}, ti = (...e) => {
  const t = Ao().createApp(...e), { mount: n } = t;
  return t.mount = (s) => {
    const r = kc(s);
    if (!r) return;
    const i = t._component;
    !Z(i) && !i.render && !i.template && (i.template = r.innerHTML), r.nodeType === 1 && (r.textContent = "");
    const o = n(r, !1, Ac(r));
    return r instanceof Element && (r.removeAttribute("v-cloak"), r.setAttribute("data-v-app", "")), o;
  }, t;
};
function Ac(e) {
  if (e instanceof SVGElement)
    return "svg";
  if (typeof MathMLElement == "function" && e instanceof MathMLElement)
    return "mathml";
}
function kc(e) {
  return we(e) ? document.querySelector(e) : e;
}
function Rc(e) {
  const t = document.body || document.documentElement;
  let n = null;
  const s = (c) => {
    n = c;
  };
  document.addEventListener("mousemove", s), window.getComputedStyle(e.containerElement).position === "static" && ye.select(e.containerElement).style("position", "relative");
  const i = ye.select(e.containerElement).append("div").attr("class", "track-viewer-menu-button").attr("aria-label", "Track viewer options menu").attr("role", "button").attr("tabindex", "0").style("position", "absolute").style("top", "8px").style("right", "8px").style("width", "32px").style("height", "32px").style("background", "white").style("border", "1px solid #ccc").style("border-radius", "4px").style("cursor", "pointer").style("display", "flex").style("align-items", "center").style("justify-content", "center").style("z-index", "100").style("box-shadow", "0 2px 4px rgba(0,0,0,0.1)").style("opacity", "0").style("transition", "opacity 0.2s").style("pointer-events", "none").html("&#8942;").style("font-size", "20px").style("line-height", "1").style("user-select", "none"), o = ye.select(t).append("div").attr("class", "track-viewer-context-menu").style("position", "absolute").style("background", "white").style("border", "1px solid #ccc").style("border-radius", "4px").style("box-shadow", "0 2px 8px rgba(0,0,0,0.15)").style("padding", "4px 0").style("min-width", "200px").style("display", "none").style("z-index", "1001"), l = [
    { label: "Save as SVG", action: e.onSaveSVG },
    { label: "Save as PNG", action: e.onSavePNG },
    { type: "separator" },
    { label: "Show track labels", action: e.onToggleTrackLabels, checkbox: !0, checked: e.getShowTrackLabels },
    { label: "Show all annotation labels", action: e.onToggleAllAnnotationLabels, checkbox: !0, checked: e.getShowAllAnnotationLabels }
  ];
  l.forEach((c) => {
    if (c.type === "separator")
      o.append("div").style("height", "1px").style("background", "#e0e0e0").style("margin", "4px 0");
    else {
      const f = o.append("div").attr("class", "track-viewer-menu-item").style("padding", "8px 16px").style("cursor", "pointer").style("user-select", "none").style("display", "flex").style("align-items", "center").style("gap", "8px").on("mouseenter", function () {
        ye.select(this).style("background", "#f5f5f5");
      }).on("mouseleave", function () {
        ye.select(this).style("background", "transparent");
      }).on("click", () => {
        typeof c.action == "function" && c.action(), u();
      });
      if (c.checkbox) {
        const p = f.append("span").attr("class", "menu-checkbox").style("width", "16px").style("height", "16px").style("border", "1px solid #999").style("border-radius", "3px").style("display", "inline-flex").style("align-items", "center").style("justify-content", "center").style("flex-shrink", "0");
        c.checkboxElement = p, c.checked && c.checked() && p.html("&#10003;").style("background", "#2196F3").style("color", "white").style("border-color", "#2196F3");
      }
      f.append("span").text(c.label);
    }
  }), ye.select(e.containerElement).on("mouseenter", () => {
    i.style("opacity", "1").style("pointer-events", "all");
  }).on("mouseleave", () => {
    o.style("display") === "none" && i.style("opacity", "0").style("pointer-events", "none");
  });
  function a() {
    const c = i.node().getBoundingClientRect(), f = window.pageYOffset || document.documentElement.scrollTop, p = window.pageXOffset || document.documentElement.scrollLeft;
    o.style("display", "block").style("top", `${c.bottom + f + 4}px`).style("left", `${c.right + p - 200}px`);
  }
  function u() {
    o.style("display", "none");
    const c = e.containerElement.getBoundingClientRect();
    if (n) {
      const f = n.clientX, p = n.clientY;
      f >= c.left && f <= c.right && p >= c.top && p <= c.bottom || i.style("opacity", "0").style("pointer-events", "none");
    }
  }
  return i.on("click", (c) => {
    c.stopPropagation(), o.style("display") === "none" ? a() : u();
  }), ye.select(t).on("click.trackviewer-context", () => {
    u();
  }), {
    destroy() {
      document.removeEventListener("mousemove", s), ye.select(e.containerElement).on("mouseenter", null).on("mouseleave", null), ye.select(t).on("click.trackviewer-context", null), o.remove(), i.remove();
    },
    updateCheckbox(c, f) {
      const p = l.find((g) => g.label === c);
      p && p.checkboxElement && (f ? p.checkboxElement.html("&#10003;").style("background", "#2196F3").style("color", "white").style("border-color", "#2196F3") : p.checkboxElement.html("").style("background", "white").style("border-color", "#999"));
    },
    toggle() {
      o.style("display") === "none" ? a() : u();
    }
  };
}
const At = class At {
  constructor(t) {
    this.data = { tracks: [], annotations: [], primitives: [] }, this.showAllAnnotationLabels = !1, this.isDragging = !1, this.isResponsiveWidth = t.width === void 0;
    const n = typeof t.container == "string" ? document.querySelector(t.container) : t.container;
    if (!n)
      throw new Error("Container element not found");
    this.containerElement = n;
    const s = this.isResponsiveWidth ? this.calculateResponsiveWidth() : t.width || 800, r = t.showTrackLabels !== !1 ? 60 : 10, i = t.margin || { top: 20, right: 10, bottom: 20, left: r };
    this.config = {
      container: t.container,
      width: s,
      height: t.height || 300,
      margin: i,
      trackHeight: t.trackHeight || 30,
      domain: t.domain || [0, 100],
      zoomExtent: t.zoomExtent || [0.5, 20],
      showTrackLabels: t.showTrackLabels !== void 0 ? t.showTrackLabels : !0,
      onAnnotationClick: t.onAnnotationClick || (() => {
      }),
      onAnnotationHover: t.onAnnotationHover || (() => {
      }),
      onBackgroundClick: t.onBackgroundClick || (() => {
      })
    }, this.originalLeftMargin = this.config.margin.left, this.showTrackLabels = this.config.showTrackLabels, this.showTrackLabels || (this.config.margin.left = this.originalLeftMargin), this.currentTransform = ye.zoomIdentity, this.initialize(), this.isResponsiveWidth && this.setupAutoResize();
  }
  calculateResponsiveWidth() {
    const t = this.containerElement.clientWidth || 800;
    return Math.max(600, t - 40);
  }
  setupAutoResize() {
    let t;
    this.resizeObserver = new ResizeObserver((n) => {
      clearTimeout(t), t = setTimeout(() => {
        this.resize();
      }, 100);
    }), this.resizeObserver.observe(this.containerElement);
  }
  initialize() {
    const t = typeof this.config.container == "string" ? document.querySelector(this.config.container) : this.config.container;
    if (!t)
      throw new Error("Container element not found");
    const n = document.body || document.documentElement;
    this.tooltip = ye.select(n).append("div").attr("class", "track-viewer-tooltip").style("position", "absolute").style("background", "white").style("border", "1px solid #ccc").style("padding", "4px 8px").style("pointer-events", "none").style("display", "none").style("z-index", "1000"), this.svg = ye.select(t).append("svg").attr("width", this.config.width).attr("height", this.config.height), this.chart = this.svg.append("g").attr("transform", `translate(${this.config.margin.left},${this.config.margin.top})`);
    const s = this.config.width - this.config.margin.left - this.config.margin.right;
    this.clipId = `clip-${Math.random().toString(36).substring(2, 11)}`, this.svg.append("defs").append("clipPath").attr("id", this.clipId).append("rect").attr("x", 0).attr("y", -At.LABEL_PADDING).attr("width", s).attr("height", "100%"), this.xAxisGroup = this.chart.append("g").attr("class", "x-axis"), this.clippedChart = this.chart.append("g").attr("clip-path", `url(#${this.clipId})`), this.x = ye.scaleLinear().domain(this.config.domain).range([0, this.config.width - this.config.margin.left - this.config.margin.right]), this.initializeZoom(), this.createContextMenu();
  }
  initializeZoom() {
    const t = this.config.width - this.config.margin.left - this.config.margin.right, n = this.config.height - this.config.margin.top - this.config.margin.bottom;
    this.zoom = ye.zoom().scaleExtent(this.config.zoomExtent).translateExtent([[0, 0], [t, n]]).extent([[0, 0], [t, n]]).filter((s) => (s.type === "mousedown" && s.button === 0 && (this.isDragging = !0), s.type === "mouseup" && (this.isDragging = !1), !s.ctrlKey && !s.button)).on("zoom", (s) => {
      this.currentTransform = s.transform, this.drawTracks();
    }).on("start", (s) => {
      s.sourceEvent && s.sourceEvent.type === "mousedown" && (this.clippedChart.select(".chart-background").style("cursor", "grabbing"), this.svg.style("cursor", "grabbing"));
    }).on("end", () => {
      this.isDragging = !1, this.clippedChart.select(".chart-background").style("cursor", "grab"), this.svg.style("cursor", "");
    }), this.svg.call(this.zoom), this.clippedChart.insert("rect", ":first-child").attr("class", "chart-background").attr("width", t).attr("height", n).attr("y", -At.LABEL_PADDING).style("fill", "transparent").style("pointer-events", "all").style("cursor", "grab").on("click", () => {
      this.config.onBackgroundClick();
    });
  }
  createContextMenu() {
    this.contextMenuController = Rc({
      containerElement: this.containerElement,
      onSaveSVG: () => this.saveAsSVG(),
      onSavePNG: () => this.saveAsPNG(),
      getShowTrackLabels: () => this.showTrackLabels,
      getShowAllAnnotationLabels: () => this.showAllAnnotationLabels,
      onToggleTrackLabels: () => this.toggleTrackLabels(),
      onToggleAllAnnotationLabels: () => this.toggleAllAnnotationLabels()
    });
  }
  saveAsSVG() {
    const t = this.svg.node();
    if (!t) return;
    const n = t.cloneNode(!0), s = this.getInlineStyles(n), r = document.createElementNS("http://www.w3.org/2000/svg", "style");
    r.textContent = s, n.insertBefore(r, n.firstChild), ye.select(n).selectAll("[data-styled]").attr("data-styled", null);
    const o = new XMLSerializer().serializeToString(n), l = new Blob([o], { type: "image/svg+xml" }), a = URL.createObjectURL(l), u = document.createElement("a");
    u.href = a, u.download = "track-viewer.svg", u.click(), URL.revokeObjectURL(a);
  }
  saveAsPNG() {
    const t = this.svg.node();
    if (!t) return;
    const n = t.cloneNode(!0), s = this.getInlineStyles(n), r = document.createElementNS("http://www.w3.org/2000/svg", "style");
    r.textContent = s, n.insertBefore(r, n.firstChild), ye.select(n).selectAll("[data-styled]").attr("data-styled", null);
    const o = new XMLSerializer().serializeToString(n), l = 20, a = document.createElement("canvas"), u = a.getContext("2d");
    if (!u) return;
    const c = 2;
    a.width = (this.config.width + l * 2) * c, a.height = (this.config.height + l * 2) * c, u.scale(c, c);
    const f = new Image(), p = new Blob([o], { type: "image/svg+xml" }), g = URL.createObjectURL(p);
    f.onload = () => {
      u.fillStyle = "white", u.fillRect(0, 0, this.config.width + l * 2, this.config.height + l * 2), u.drawImage(f, l, l), URL.revokeObjectURL(g), a.toBlob((m) => {
        if (m) {
          const x = URL.createObjectURL(m), A = document.createElement("a");
          A.href = x, A.download = "track-viewer.png", A.click(), URL.revokeObjectURL(x);
        }
      });
    }, f.src = g;
  }
  getInlineStyles(t = this.svg.node()) {
    const n = Array.from(document.styleSheets), s = [];
    n.forEach((i) => {
      try {
        const o = i.cssRules || i.rules;
        o && Array.from(o).forEach((l) => {
          if (l.selectorText) {
            const a = l.selectorText;
            (a.includes(".track") || a.includes(".annotation") || a.includes(".axis") || a.includes(".primitive") || a.includes("svg") || a.includes("text") || a.includes("rect") || a.includes("path") || a.includes("circle") || a.includes("line") || a.includes("polygon")) && s.push(l.cssText);
          }
        });
      } catch {
      }
    });
    const r = ye.select(t);
    return r.selectAll(".track-label").each(function () {
      const i = this, o = window.getComputedStyle(i);
      i.hasAttribute("data-styled") || (i.setAttribute("data-styled", "true"), i.setAttribute(
        "style",
        `font-family: ${o.fontFamily}; font-size: ${o.fontSize}; fill: ${o.fill || "currentColor"};`
      ));
    }), r.selectAll(".annotation, .primitive").each(function () {
      const i = this, o = window.getComputedStyle(i), l = i.getAttribute("data-label");
      if (l && !i.querySelector("title")) {
        const u = document.createElementNS("http://www.w3.org/2000/svg", "title");
        u.textContent = l, i.insertBefore(u, i.firstChild), i.removeAttribute("data-label");
      }
      i.querySelectorAll("path, rect, circle, line, polygon, ellipse").forEach((u) => {
        const c = u, f = window.getComputedStyle(c);
        if (!c.hasAttribute("data-styled")) {
          c.setAttribute("data-styled", "true");
          const p = [], g = f.fill, m = o.fill;
          if (g && g !== "none" && g !== "rgb(0, 0, 0)" ? p.push(`fill: ${g}`) : m && m !== "none" && m !== "rgb(0, 0, 0)" && p.push(`fill: ${m}`), f.stroke && f.stroke !== "none" && p.push(`stroke: ${f.stroke}`), f.strokeWidth && p.push(`stroke-width: ${f.strokeWidth}`), f.opacity && f.opacity !== "1" && p.push(`opacity: ${f.opacity}`), p.length > 0) {
            const x = c.getAttribute("style") || "";
            c.setAttribute("style", x + " " + p.join("; ") + ";");
          }
        }
      });
    }), s.join(`
`);
  }
  toggleTrackLabels() {
    var n;
    if (this.showTrackLabels = !this.showTrackLabels, this.chart.select(".track-labels-container").style("display", this.showTrackLabels ? "" : "none"), this.showTrackLabels)
      this.updateMarginAndLayout();
    else {
      this.config.margin.left = this.originalLeftMargin, this.chart.attr("transform", `translate(${this.config.margin.left},${this.config.margin.top})`), this.x.range([0, this.config.width - this.config.margin.left - this.config.margin.right]);
      const s = this.config.width - this.config.margin.left - this.config.margin.right;
      this.svg.select("clipPath rect").attr("width", s), this.clippedChart.select(".chart-background").attr("width", s), this.drawTracks();
    }
    (n = this.contextMenuController) == null || n.updateCheckbox("Show track labels", this.showTrackLabels);
  }
  toggleAllAnnotationLabels() {
    var n;
    this.showAllAnnotationLabels = !this.showAllAnnotationLabels;
    const t = this.showAllAnnotationLabels;
    this.clippedChart.selectAll(".annotation-labels-layer .annotation-label-group").each(function () {
      const s = ye.select(this), i = s.select(".annotation-label").attr("data-show-label");
      i === "never" ? s.style("display", "none") : i === "always" ? s.style("display", "") : s.style("display", t ? "" : "none");
    }), (n = this.contextMenuController) == null || n.updateCheckbox("Show all annotation labels", this.showAllAnnotationLabels);
  }
  // Helper methods for track heights
  getTrackHeight(t) {
    return t.height || this.config.trackHeight;
  }
  getTotalTracksHeight() {
    return this.data.tracks.reduce((t, n) => t + this.getTrackHeight(n), 0);
  }
  getTrackYPosition(t) {
    let n = 0;
    for (let s = 0; s < t; s++)
      n += this.getTrackHeight(this.data.tracks[s]);
    return n;
  }
  updateHeight() {
    const t = this.getTotalTracksHeight() + this.config.margin.top + this.config.margin.bottom;
    this.config.height = t, this.svg.attr("height", t);
    const n = this.getTotalTracksHeight();
    this.xAxisGroup.attr("transform", `translate(0, ${n})`), this.svg.select("clipPath rect").attr("height", n + At.LABEL_PADDING), this.clippedChart.select(".chart-background").attr("height", n + At.LABEL_PADDING);
  }
  createTrackGroups() {
    this.trackGroups = this.clippedChart.selectAll(".track").data(this.data.tracks, (n) => n.id).join("g").attr("id", (n) => n.id).attr("class", "track").attr("transform", (n, s) => `translate(0, ${this.getTrackYPosition(s)})`);
    let t = this.chart.select(".track-labels-container");
    t.empty() && (t = this.chart.append("g").attr("class", "track-labels-container")), t.style("display", this.showTrackLabels ? "" : "none"), this.labelGroups = t.selectAll(".track-label-group").data(this.data.tracks, (n) => n.id).join("g").attr("id", (n) => `${n.id}-label`).attr("class", "track-label-group").attr("transform", (n, s) => `translate(0, ${this.getTrackYPosition(s)})`), this.labelGroups.selectAll(".track-label").data((n) => [n]).join("text").attr("class", "track-label").attr("x", -10).attr("y", (n) => this.getTrackHeight(n) / 2).attr("dy", "0.35em").attr("text-anchor", "end").text((n) => n.label);
  }
  drawTracks() {
    const t = this.currentTransform.rescaleX(this.x);
    if (this.xAxisGroup.call(ye.axisBottom(t)), this.xAxisGroup.selectAll("text").attr("class", "axis-label"), !this.trackGroups)
      return;
    const n = [];
    this.trackGroups.each((s, r, i) => {
      const o = ye.select(i[r]), l = this.getTrackYPosition(r), a = this.data.annotations.filter((g) => g.trackId === s.id), c = [
        ...(this.data.primitives || []).filter((g) => g.trackId === s.id).map((g) => ({ type: "primitive", data: g })),
        ...a.map((g) => ({ type: "annotation", data: g }))
      ], f = o.selectAll("g.element").data(c, (g) => g.data.id).join(
        // Enter: create new groups
        (g) => g.append("g"),
        // Update: existing groups (no action needed)
        (g) => g,
        // Exit: remove groups that are no longer in data
        (g) => g.remove()
      ).attr("id", (g) => g.data.id).attr("class", (g) => g.type === "primitive" ? `element ${g.type} ${g.data.type} ${g.data.class}` : `element ${g.type} ${g.data.type} ${g.data.classes.join(" ")}`), p = this;
      f.each(function (g) {
        const m = ye.select(this);
        m.selectAll("*").remove(), g.type === "primitive" ? p.renderPrimitive(m, g.data, t, s) : p.renderAnnotation(m, g.data, t, s, n, l);
      });
    }), this.renderAllLabels(n, t);
  }
  renderAnnotation(t, n, s, r, i, o = 0) {
    const l = s(n.start), a = Math.max(1, s(n.end) - s(n.start)), u = this.getTrackHeight(r), c = n.heightFraction !== void 0 ? u * n.heightFraction : u / 2, f = u - (n.fy !== void 0 ? n.fy * u : u / 2);
    let p;
    switch (n.type) {
      case "arrow":
        p = this.renderArrow(t, l, f, a, c, n.direction);
        break;
      case "circle":
        p = this.renderCircle(t, l, f, c);
        break;
      case "triangle":
        p = this.renderTriangle(t, l, f, c);
        break;
      case "pin":
        p = this.renderPin(t, l, f, c);
        break;
      case "box":
      default:
        p = this.renderBox(t, l, f, a, c, n.corner_radius || 0);
        break;
    }
    t.attr("data-label", n.label), p.attr("data-annotation-id", n.id).style("cursor", "pointer").style("pointer-events", "all"), n.fill && p.style("fill", n.fill), n.stroke && p.style("stroke", n.stroke), n.opacity !== void 0 && p.style("opacity", n.opacity), p.on("mouseover", (g) => {
      this.isDragging || (p.classed("hovered", !0), this.clippedChart.selectAll(`.annotation-label-group[data-annotation-id="${n.id}"]`).style("display", "block").select(".annotation-label-background").style("display", ""), n.tooltip !== void 0 && this.showTooltip(g, n, r), this.config.onAnnotationHover(n, r, g));
    }).on("mouseout", () => {
      p.classed("hovered", !1);
      const g = this.showAllAnnotationLabels;
      this.clippedChart.selectAll(`.annotation-label-group[data-annotation-id="${n.id}"]`).each(function () {
        const m = ye.select(this), A = m.select(".annotation-label").attr("data-show-label");
        A === "hover" && g ? m.style("display", "") : A === "hover" ? m.style("display", "none") : m.style("display", ""), m.select(".annotation-label-background").style("display", "none");
      }), this.hideTooltip();
    }).on("click", () => {
      this.config.onAnnotationClick(n, r);
    }), i && n.label && n.showLabel !== "never" && i.push({ annotation: n, trackData: r, x: l, y: f, width: a, height: c, trackY: o });
  }
  renderPrimitive(t, n, s, r) {
    const i = this.getTrackHeight(r);
    let o;
    switch (n.type) {
      case "horizontal-line":
        o = this.renderHorizontalLine(t, n, s, i);
        break;
      case "background":
        o = this.renderBackground(t, n, s, i);
        break;
      default:
        return;
    }
    o.attr("class", `primitive ${n.class}`), n.stroke ? o.style("stroke", n.stroke) : n.type === "horizontal-line" && o.style("stroke", "currentColor"), n.fill ? o.style("fill", n.fill) : n.type === "background" && o.style("fill", "none"), n.opacity !== void 0 && o.style("opacity", n.opacity);
  }
  renderHorizontalLine(t, n, s, r) {
    const i = s.range();
    let o, l;
    n.start !== void 0 && n.end !== void 0 ? (o = s(n.start), l = s(n.end)) : n.start !== void 0 && n.end === void 0 ? (o = s(n.start), l = i[1]) : n.start === void 0 && n.end !== void 0 ? (o = i[0], l = s(n.end)) : (o = i[0], l = i[1]);
    const a = r - (n.fy !== void 0 ? n.fy * r : r / 2), u = Math.round(a) + 0.5;
    return t.append("line").attr("x1", o).attr("x2", l).attr("y1", u).attr("y2", u).style("stroke-width", 1).style("shape-rendering", "crispEdges");
  }
  renderBackground(t, n, s, r) {
    const i = s.range();
    let o, l;
    return n.start !== void 0 && n.end !== void 0 ? (o = s(n.start), l = s(n.end)) : n.start !== void 0 && n.end === void 0 ? (o = s(n.start), l = i[1]) : n.start === void 0 && n.end !== void 0 ? (o = i[0], l = s(n.end)) : (o = i[0], l = i[1]), t.append("rect").attr("x", o).attr("y", 0).attr("width", l - o).attr("height", r);
  }
  renderBox(t, n, s, r, i, o = 0) {
    return t.append("rect").attr("x", n).attr("y", s - i / 2).attr("width", r).attr("height", i).attr("rx", o).attr("ry", o);
  }
  renderArrow(t, n, s, r, i, o) {
    const l = Math.min(r * 0.2, i * 0.5, 8), a = r - l;
    let u;
    const c = s - i / 2, f = s + i / 2;
    if (o === "right")
      u = `
        M ${n} ${c}
        L ${n + a} ${c}
        L ${n + r} ${s}
        L ${n + a} ${f}
        L ${n} ${f}
        Z
      `;
    else if (o === "left")
      u = `
        M ${n} ${s}
        L ${n + l} ${c}
        L ${n + r} ${c}
        L ${n + r} ${f}
        L ${n + l} ${f}
        Z
      `;
    else {
      const p = Math.min(r * 0.1, i * 0.3, 4);
      u = `
        M ${n + p} ${c}
        L ${n + r - p} ${c}
        L ${n + r} ${s}
        L ${n + r - p} ${f}
        L ${n + p} ${f}
        L ${n} ${s}
        Z
      `;
    }
    return t.append("path").attr("d", u);
  }
  renderCircle(t, n, s, r) {
    return t.append("circle").attr("cx", n).attr("cy", s).attr("r", r / 2).attr("class", "annotation-marker");
  }
  renderTriangle(t, n, s, r = 0.5) {
    const i = r * 0.8, o = [
      [n - i / 2, s + r],
      [n, s],
      [n + i / 2, s + r]
    ].map((l) => l.join(",")).join(" ");
    return t.append("polygon").attr("points", o).attr("class", "annotation-triangle");
  }
  renderPin(t, n, s, r) {
    const i = t.append("g").attr("class", "annotation-pin");
    return i.append("line").attr("x1", n).attr("y1", s).attr("x2", n).attr("y2", s - r).attr("stroke", "black").attr("class", "annotation-pin-line"), i.append("circle").attr("cx", n).attr("cy", s - r).attr("r", 3).attr("class", "annotation-pin-head"), i;
  }
  renderAnnotationLabel(t, n, s, r, i, o) {
    const l = n.showLabel || "hover";
    if (!n.label)
      return;
    const a = n.labelPosition || "above";
    let u, c, f = "middle";
    n.type === "circle" || n.type === "triangle" || n.type === "pin" ? (u = s, a === "above" ? n.type === "pin" ? c = r - o - 8 : c = r - o / 2 - 2 : c = r) : (u = s + i / 2, a === "above" ? c = r - o / 2 - 8 : c = r);
    const p = t.append("text").attr("x", u).attr("y", c).attr("dy", "0.35em").attr("text-anchor", f).attr("class", "annotation-label").attr("data-annotation-id", n.id).attr("data-show-label", l).style("fill", "currentColor").text(n.label);
    if (l === "hover") {
      const g = p.node();
      if (g && typeof g.getBBox == "function")
        try {
          const m = g.getBBox(), x = 0;
          t.insert("rect", "text").attr("x", m.x - x).attr("y", m.y - x).attr("width", m.width + x * 2).attr("height", m.height + x * 2).attr("rx", 2).attr("ry", 2).attr("class", "annotation-label-background").style("fill", "white").style("opacity", "0.8").style("display", "none");
        } catch {
        }
    }
    l === "always" ? t.style("display", "") : l === "never" ? t.style("display", "none") : t.style("display", this.showAllAnnotationLabels ? "" : "none");
  }
  renderAllLabels(t, n) {
    let s = this.clippedChart.select(".annotation-labels-layer");
    s.empty() && (s = this.clippedChart.append("g").attr("class", "annotation-labels-layer"));
    const r = s.selectAll(".annotation-label-group").data(t, (o) => o.annotation.id).join(
      (o) => o.append("g"),
      (o) => o,
      (o) => o.remove()
    ).attr("class", (o) => `annotation-label-group ${o.annotation.classes.join(" ")}`).attr("data-annotation-id", (o) => o.annotation.id).attr("transform", (o) => `translate(0, ${o.trackY})`).style("pointer-events", "none"), i = this;
    r.each(function (o) {
      const l = ye.select(this);
      l.selectAll("*").remove(), i.renderAnnotationLabel(l, o.annotation, o.x, o.y, o.width, o.height);
    });
  }
  showTooltip(t, n, s) {
    const r = n.tooltip !== void 0 ? n.tooltip : `<strong>${n.label}</strong><br/>Start: ${n.start}<br/>End: ${n.end}<br/>Track: ${s.label}`;
    this.tooltip.style("display", "block").style("left", `${t.pageX + 10}px`).style("top", `${t.pageY - 10}px`).html(r);
  }
  hideTooltip() {
    this.tooltip.style("display", "none");
  }
  // Calculate required left margin based on track label lengths
  calculateRequiredLeftMargin() {
    if (!this.showTrackLabels)
      return this.originalLeftMargin;
    if (this.data.tracks.length === 0)
      return this.originalLeftMargin;
    const t = this.svg.append("text").attr("class", "track-label").style("visibility", "hidden");
    let n = 0;
    return this.data.tracks.forEach((r) => {
      t.text(r.label);
      const i = t.node();
      if (i && typeof i.getBBox == "function")
        try {
          const o = i.getBBox();
          n = Math.max(n, o.width);
        } catch {
          n = Math.max(n, r.label.length * 8);
        }
      else
        n = Math.max(n, r.label.length * 8);
    }), t.remove(), Math.max(this.originalLeftMargin, n + 20);
  }
  updateMarginAndLayout() {
    const t = this.calculateRequiredLeftMargin();
    if (Math.abs(this.config.margin.left - t) > 5) {
      this.config.margin.left = t, this.chart.attr("transform", `translate(${this.config.margin.left},${this.config.margin.top})`), this.x.range([0, this.config.width - this.config.margin.left - this.config.margin.right]);
      const n = this.config.width - this.config.margin.left - this.config.margin.right;
      this.svg.select("clipPath rect").attr("width", n), this.clippedChart.select(".chart-background").attr("width", n);
    }
  }
  // Public API methods
  setData(t) {
    this.data = {
      tracks: t.tracks,
      annotations: t.annotations,
      primitives: t.primitives || []
    }, this.updateMarginAndLayout(), this.updateHeight(), this.createTrackGroups(), this.drawTracks();
  }
  addTrack(t, n, s) {
    this.data.tracks.push(t), n && this.data.annotations.push(...n), s && (this.data.primitives || (this.data.primitives = []), this.data.primitives.push(...s)), this.updateMarginAndLayout(), this.updateHeight(), this.createTrackGroups(), this.drawTracks();
  }
  removeTrack(t) {
    this.data.tracks = this.data.tracks.filter((n) => n.id !== t), this.data.annotations = this.data.annotations.filter((n) => n.trackId !== t), this.data.primitives = (this.data.primitives || []).filter((n) => n.trackId !== t), this.updateMarginAndLayout(), this.updateHeight(), this.createTrackGroups(), this.drawTracks();
  }
  addAnnotation(t) {
    this.data.annotations.push(t), this.drawTracks();
  }
  removeAnnotation(t) {
    this.data.annotations = this.data.annotations.filter((n) => n.id !== t), this.drawTracks();
  }
  addPrimitive(t) {
    this.data.primitives || (this.data.primitives = []), this.data.primitives.push(t), this.drawTracks();
  }
  removePrimitive(t) {
    this.data.primitives && (this.data.primitives = this.data.primitives.filter((n) => n.id !== t), this.drawTracks());
  }
  updateDomain(t) {
    this.config.domain = t, this.x.domain(t), this.drawTracks();
  }
  zoomTo(t, n) {
    const r = (this.config.width - this.config.margin.left - this.config.margin.right) / (this.x(n) - this.x(t)), i = -this.x(t) * r, o = ye.zoomIdentity.translate(i, 0).scale(r);
    this.svg.transition().duration(750).call(this.zoom.transform, o);
  }
  resetZoom() {
    this.svg.transition().duration(750).call(this.zoom.transform, ye.zoomIdentity);
  }
  destroy() {
    var t;
    this.resizeObserver && this.resizeObserver.disconnect(), (t = this.contextMenuController) == null || t.destroy(), this.tooltip.remove(), this.svg.remove();
  }
  getConfig() {
    return { ...this.config };
  }
  getData() {
    return {
      tracks: [...this.data.tracks],
      annotations: [...this.data.annotations],
      primitives: this.data.primitives ? [...this.data.primitives] : []
    };
  }
  resize() {
    if (!this.isResponsiveWidth)
      return;
    const t = this.calculateResponsiveWidth();
    if (Math.abs(this.config.width - t) > 10) {
      this.config.width = t, this.svg.attr("width", t), this.x.range([0, t - this.config.margin.left - this.config.margin.right]);
      const n = t - this.config.margin.left - this.config.margin.right;
      this.svg.select("clipPath rect").attr("width", n), this.clippedChart.select(".chart-background").attr("width", n);
      const s = this.config.height - this.config.margin.top - this.config.margin.bottom;
      this.zoom.translateExtent([[0, 0], [n, s]]).extent([[0, 0], [n, s]]), this.drawTracks();
    }
  }
};
At.LABEL_PADDING = 24;
let Vs = At;
const zt = (e, t) => {
  const n = e.__vccOpts || e;
  for (const [s, r] of t)
    n[s] = r;
  return n;
}, Oc = {
  name: "SimpleTable",
  props: {
    rows: {
      type: Array,
      required: !0,
      validator: (e) => e.every((t) => Array.isArray(t))
    }
  },
  setup(e) {
    const t = he(null), n = he("asc"), s = We(() => e.rows[0] || []), r = We(() => e.rows.slice(1)), i = (l) => {
      t.value === l ? n.value === "asc" ? n.value = "desc" : (t.value = null, n.value = "asc") : (t.value = l, n.value = "asc");
    }, o = We(() => t.value === null ? r.value : [...r.value].sort((l, a) => {
      const u = l[t.value], c = a[t.value];
      if (u === "" || u === null || u === void 0) return 1;
      if (c === "" || c === null || c === void 0) return -1;
      const f = parseFloat(String(u).replace(/[%,]/g, "")), p = parseFloat(String(c).replace(/[%,]/g, ""));
      let g = 0;
      return !isNaN(f) && !isNaN(p) ? g = f - p : g = String(u).localeCompare(String(c)), n.value === "asc" ? g : -g;
    }));
    return {
      headers: s,
      sortColumn: t,
      sortDirection: n,
      sortedDataRows: o,
      toggleSort: i
    };
  }
}, Pc = { class: "simple-table" }, Lc = ["onClick"], Fc = {
  key: 0,
  class: "sort-indicator"
};
function Mc(e, t, n, s, r, i) {
  return $(), B("table", Pc, [
    M("thead", null, [
      M("tr", null, [
        ($(!0), B(_e, null, Xe(s.headers, (o, l) => ($(), B("th", {
          key: l,
          onClick: (a) => s.toggleSort(l),
          class: Lt({
            sortable: !0,
            sorted: s.sortColumn === l,
            "sort-asc": s.sortColumn === l && s.sortDirection === "asc",
            "sort-desc": s.sortColumn === l && s.sortDirection === "desc"
          })
        }, [
          St(se(o) + " ", 1),
          s.sortColumn === l ? ($(), B("span", Fc, se(s.sortDirection === "asc" ? "▲" : "▼"), 1)) : be("", !0)
        ], 10, Lc))), 128))
      ])
    ]),
    M("tbody", null, [
      ($(!0), B(_e, null, Xe(s.sortedDataRows, (o, l) => ($(), B("tr", { key: l }, [
        ($(!0), B(_e, null, Xe(o, (a, u) => ($(), B("td", { key: u }, se(a), 1))), 128))
      ]))), 128))
    ])
  ]);
}
const Rs = /* @__PURE__ */ zt(Oc, [["render", Mc], ["__scopeId", "data-v-3dcf836d"]]), Ic = {
  name: "SortableTable",
  props: {
    // Array of header objects: [{ label: 'Column Name', cellClass: 'custom-class' }]
    headers: {
      type: Array,
      required: !0
    },
    // Array of row data (arrays matching header order)
    rows: {
      type: Array,
      required: !0
    },
    // Initial sort column index
    initialSortColumn: {
      type: Number,
      default: null
    },
    // Initial sort direction
    initialSortDirection: {
      type: String,
      default: "asc",
      validator: (e) => ["asc", "desc"].includes(e)
    }
  },
  setup(e) {
    const t = he(e.initialSortColumn), n = he(e.initialSortDirection), s = (i) => {
      t.value === i ? n.value === "asc" ? n.value = "desc" : (t.value = null, n.value = "asc") : (t.value = i, n.value = "asc");
    }, r = We(() => t.value === null ? e.rows : [...e.rows].sort((i, o) => {
      var m, x;
      const l = i[t.value], a = o[t.value], u = typeof l == "object" && l !== null ? l.children || ((m = l.props) == null ? void 0 : m.children) || "" : l, c = typeof a == "object" && a !== null ? a.children || ((x = a.props) == null ? void 0 : x.children) || "" : a;
      if (u === "" || u === null || u === void 0) return 1;
      if (c === "" || c === null || c === void 0) return -1;
      const f = parseFloat(String(u).replace(/[%,]/g, "")), p = parseFloat(String(c).replace(/[%,]/g, ""));
      let g = 0;
      return !isNaN(f) && !isNaN(p) ? g = f - p : g = String(u).localeCompare(String(c)), n.value === "asc" ? g : -g;
    }));
    return {
      sortColumn: t,
      sortDirection: n,
      sortedRows: r,
      toggleSort: s
    };
  }
}, Nc = { class: "sortable-table" }, Dc = ["onClick"], $c = {
  key: 0,
  class: "sort-indicator"
};
function jc(e, t, n, s, r, i) {
  return $(), B("table", Nc, [
    M("thead", null, [
      M("tr", null, [
        ($(!0), B(_e, null, Xe(n.headers, (o, l) => ($(), B("th", {
          key: l,
          onClick: (a) => s.toggleSort(l),
          class: Lt({
            sortable: !0,
            sorted: s.sortColumn === l,
            "sort-asc": s.sortColumn === l && s.sortDirection === "asc",
            "sort-desc": s.sortColumn === l && s.sortDirection === "desc"
          })
        }, [
          St(se(o.label) + " ", 1),
          s.sortColumn === l ? ($(), B("span", $c, se(s.sortDirection === "asc" ? "▲" : "▼"), 1)) : be("", !0)
        ], 10, Dc))), 128))
      ])
    ]),
    M("tbody", null, [
      ($(!0), B(_e, null, Xe(s.sortedRows, (o, l) => ($(), B("tr", { key: l }, [
        ($(!0), B(_e, null, Xe(o, (a, u) => {
          var c;
          return $(), B("td", {
            key: u,
            class: Lt((c = n.headers[u]) == null ? void 0 : c.cellClass)
          }, [
            typeof a == "object" && a !== null ? ($(), mn(Qi(a), { key: 0 })) : ($(), B(_e, { key: 1 }, [
              St(se(a), 1)
            ], 64))
          ], 2);
        }), 128))
      ]))), 128))
    ])
  ]);
}
const Bc = /* @__PURE__ */ zt(Ic, [["render", jc], ["__scopeId", "data-v-b8895c87"]]), Hc = {
  name: "FeatureDetails",
  components: {
    SimpleTable: Rs,
    SortableTable: Bc
  },
  props: {
    // The selected feature object
    feature: {
      type: Object,
      default: null
    },
    // All features (to find related PFAM domains)
    allFeatures: {
      type: Array,
      default: () => []
    },
    // Data provider for fetching MiBIG entries
    dataProvider: {
      type: Object,
      default: null
    },
    // Record information containing recordId
    recordInfo: {
      type: Object,
      default: null
    },
    // Region number for MiBIG entry retrieval
    regionNumber: {
      type: String,
      default: "1"
    }
  },
  emits: ["close"],
  setup(e) {
    const t = he([]), n = he(!1), s = he(null), r = async () => {
      var P, D, Q;
      if (t.value = [], n.value = !1, s.value = null, !e.feature || e.feature.type !== "CDS" || !e.dataProvider || !e.recordInfo) return;
      const v = (D = (P = e.feature.qualifiers) == null ? void 0 : P.locus_tag) == null ? void 0 : D[0];
      if (v) {
        n.value = !0;
        try {
          const te = await e.dataProvider.getMiBIGEntries(e.recordInfo.recordId, v, e.regionNumber);
          t.value = te.entries || [];
        } catch (te) {
          ((Q = te.response) == null ? void 0 : Q.status) === 404 ? t.value = [] : (console.error("Error fetching MiBIG entries:", te), s.value = "Failed to load MiBIG entries");
        } finally {
          n.value = !1;
        }
      }
    };
    Ke(() => e.feature, () => {
      r();
    }, { immediate: !0 });
    const i = We(() => {
      var P, D, Q, te;
      if (!e.feature || e.feature.type !== "CDS") return [];
      if (!e.allFeatures || e.allFeatures.length === 0) return [];
      const v = ((D = (P = e.feature.qualifiers) == null ? void 0 : P.locus_tag) == null ? void 0 : D[0]) || ((te = (Q = e.feature.qualifiers) == null ? void 0 : Q.gene) == null ? void 0 : te[0]);
      return v ? e.allFeatures.filter((oe) => {
        var le, z, X, G;
        return oe.type === "PFAM_domain" && (((z = (le = oe.qualifiers) == null ? void 0 : le.locus_tag) == null ? void 0 : z[0]) === v || ((G = (X = oe.qualifiers) == null ? void 0 : X.gene) == null ? void 0 : G[0]) === v);
      }).map((oe) => {
        var le, z, X, G, b, _, R, O, j, re, Ae, Je, et;
        return {
          id: ((X = (z = (le = oe.qualifiers) == null ? void 0 : le.db_xref) == null ? void 0 : z[0]) == null ? void 0 : X.replace("PFAM:", "")) || ((R = (_ = (b = (G = oe.qualifiers) == null ? void 0 : G.inference) == null ? void 0 : b[0]) == null ? void 0 : _.match(/PFAM:([^,\s]+)/)) == null ? void 0 : R[1]) || "Unknown",
          description: ((j = (O = oe.qualifiers) == null ? void 0 : O.description) == null ? void 0 : j[0]) || "",
          location: l(oe.qualifiers),
          score: ((Ae = (re = oe.qualifiers) == null ? void 0 : re.score) == null ? void 0 : Ae[0]) || "",
          evalue: ((et = (Je = oe.qualifiers) == null ? void 0 : Je.evalue) == null ? void 0 : et[0]) || ""
        };
      }) : [];
    }), o = We(() => {
      if (!e.feature) return [];
      const v = [], P = e.feature.qualifiers || {};
      return P.go_function && v.push(...P.go_function), P.go_process && v.push(...P.go_process), P.go_component && v.push(...P.go_component), v;
    }), l = (v) => {
      var Q, te;
      if (!v) return "";
      const P = (Q = v.protein_start) == null ? void 0 : Q[0], D = (te = v.protein_end) == null ? void 0 : te[0];
      return P && D ? `[${P}:${D}]` : "";
    }, a = (v) => {
      if (!v) return "";
      const P = v.match(/\[<?(\d+):>?(\d+)\](?:\(([+-])\))?/);
      if (!P) return v;
      const D = parseInt(P[1]), Q = parseInt(P[2]), te = P[3] || "", oe = Q - D;
      return `${D.toLocaleString()} - ${Q.toLocaleString()}${te ? ", " + te : ""} (total: ${oe.toLocaleString()} nt)`;
    }, u = () => {
      var oe, le, z, X, G, b;
      if (!e.feature) return "";
      const v = e.feature.qualifiers || {}, P = ((oe = v.locus_tag) == null ? void 0 : oe[0]) || "", D = ((le = v.product) == null ? void 0 : le[0]) || "", Q = ((z = v.gene) == null ? void 0 : z[0]) || "", te = ((X = v.description) == null ? void 0 : X[0]) || "";
      if (e.feature.type === "CDS") {
        if (P && D)
          return `${P} - ${D}`;
        if (P) return P;
        if (Q) return Q;
        if (D) return D;
      }
      if (e.feature.type === "PFAM_domain") {
        const _ = ((b = (G = v.db_xref) == null ? void 0 : G[0]) == null ? void 0 : b.replace("PFAM:", "")) || "PFAM domain";
        return te ? `${_} - ${te}` : _;
      }
      return te || D || e.feature.type || "Feature Details";
    }, c = (v) => v.split("_").map(
      (D, Q) => Q === 0 ? D.charAt(0).toUpperCase() + D.slice(1) : D
    ).join(" "), f = (v) => Array.isArray(v) ? v.join(", ") : v, p = () => {
      if (!e.feature || !e.feature.qualifiers) return !1;
      const v = [
        "locus_tag",
        "protein_id",
        "gene",
        "product",
        "gene_kind",
        "translation",
        "nucleotide",
        "go_function",
        "go_process",
        "go_component"
      ];
      return Object.keys(e.feature.qualifiers).filter((D) => !v.includes(D)).length > 0;
    }, g = () => {
      if (!e.feature || !e.feature.qualifiers) return {};
      const v = [
        "locus_tag",
        "protein_id",
        "gene",
        "product",
        "gene_kind",
        "translation",
        "nucleotide",
        "go_function",
        "go_process",
        "go_component"
      ], P = {};
      return Object.keys(e.feature.qualifiers).forEach((D) => {
        v.includes(D) || (P[D] = e.feature.qualifiers[D]);
      }), P;
    }, m = (v) => {
      if (!v) return "";
      const P = [];
      for (let D = 0; D < v.length; D += 60)
        P.push(v.slice(D, D + 60));
      return P.join(`
`);
    }, x = async (v, P) => {
      try {
        await navigator.clipboard.writeText(v), alert(`${P} sequence copied to clipboard!`);
      } catch (D) {
        console.error("Failed to copy to clipboard:", D), alert("Failed to copy to clipboard");
      }
    }, A = (v, P) => ({
      NRPS_PKS: H,
      sec_met_domain: q
    }[v] || W)(P), H = (v) => {
      const P = v.map((le) => {
        const z = le.split(/\.\s+/);
        if (z.length === 0)
          return { domain: le };
        const G = { domain: z[0].replace(/^Domain:\s*/, "").trim() };
        return z.slice(1).forEach((b) => {
          const _ = b.indexOf(":");
          if (_ > -1) {
            const R = b.substring(0, _).trim(), O = b.substring(_ + 1).trim();
            G[R] = O;
          }
        }), G;
      }), D = P[0] || {}, Q = Object.keys(D).filter((le) => le !== "domain"), oe = [
        ["", ...Q],
        ...P.map((le) => {
          const z = [le.domain || ""];
          return Q.forEach((X) => {
            z.push(le[X] || "");
          }), z;
        })
      ];
      return en(Rs, { rows: oe });
    }, q = (v) => {
      const P = v.map((le) => {
        const z = le.match(/^([^(]+)\s*\(([^)]+)\)/);
        if (!z)
          return { name: le };
        const X = z[1].trim(), G = z[2], b = { name: X };
        return G.split(/,\s*/).forEach((R) => {
          const O = R.indexOf(":");
          if (O > -1) {
            const j = R.substring(0, O).trim(), re = R.substring(O + 1).trim();
            b[j] = re;
          }
        }), b;
      }), D = P[0] || {}, Q = Object.keys(D).filter((le) => le !== "name"), oe = [
        ["", ...Q],
        ...P.map((le) => {
          const z = [le.name || ""];
          return Q.forEach((X) => {
            z.push(le[X] || "");
          }), z;
        })
      ];
      return en(Rs, { rows: oe });
    }, W = (v) => {
      const P = v.map(
        (D, Q) => en("li", { key: Q }, D)
      );
      return en("ul", { class: "qualifier-list" }, P);
    }, I = (v) => typeof v == "number" ? v.toExponential(2) : v, J = [
      { label: "Domain", cellClass: "pfam-id" },
      { label: "Description", cellClass: "pfam-description" },
      { label: "Location", cellClass: "pfam-location" },
      { label: "Score", cellClass: "pfam-score" },
      { label: "E-value", cellClass: "pfam-evalue" }
    ], de = We(() => i.value.map((v) => [
      v.id,
      v.description,
      v.location,
      v.score,
      v.evalue
    ])), ie = [
      { label: "MIBiG Protein", cellClass: "mibig-protein" },
      { label: "Description", cellClass: "mibig-description" },
      { label: "MIBiG Cluster", cellClass: "mibig-cluster" },
      { label: "Product", cellClass: "mibig-product" },
      { label: "% ID", cellClass: "mibig-percent" },
      { label: "BLAST Score", cellClass: "mibig-score" },
      { label: "% Coverage", cellClass: "mibig-percent" },
      { label: "E-value", cellClass: "mibig-evalue" }
    ], Se = We(() => t.value.map((v) => [
      v.mibig_protein,
      v.description,
      en("a", {
        href: `https://mibig.secondarymetabolites.org/go/${v.mibig_cluster}`,
        target: "_blank",
        rel: "noopener noreferrer"
      }, v.mibig_cluster),
      v.mibig_product,
      `${v.percent_identity.toFixed(1)}%`,
      v.blast_score.toFixed(1),
      `${v.percent_coverage.toFixed(1)}%`,
      I(v.evalue)
    ]));
    return {
      pfamDomains: i,
      goTerms: o,
      mibigEntries: t,
      mibigLoading: n,
      mibigError: s,
      formatLocation: a,
      getFeatureTitle: u,
      formatQualifierKey: c,
      formatQualifierValue: f,
      hasAdditionalQualifiers: p,
      getAdditionalQualifiers: g,
      formatSequence: m,
      copyToClipboard: x,
      renderQualifierValue: A,
      formatEvalue: I,
      pfamTableHeaders: J,
      pfamTableRows: de,
      mibigTableHeaders: ie,
      mibigTableRows: Se
    };
  }
}, Uc = {
  key: 0,
  class: "feature-details"
}, qc = { class: "feature-content" }, Vc = { class: "feature-header" }, Wc = { class: "info-section" }, zc = {
  key: 0,
  class: "info-row"
}, Kc = { class: "info-value" }, Gc = {
  key: 1,
  class: "info-row"
}, Jc = { class: "info-value" }, Yc = {
  key: 2,
  class: "info-row"
}, Xc = { class: "info-value" }, Qc = {
  key: 3,
  class: "info-row"
}, Zc = { class: "info-value" }, eu = {
  key: 4,
  class: "info-row"
}, tu = { class: "info-value" }, nu = {
  key: 5,
  class: "info-row"
}, su = { class: "info-value" }, ru = { class: "info-label" }, iu = { class: "info-value" }, ou = {
  key: 0,
  class: "expandable-list"
}, lu = { class: "expanded-content" }, au = {
  key: 6,
  class: "info-row"
}, cu = { class: "info-value" }, uu = { class: "expandable-list" }, fu = { class: "expanded-content" }, du = {
  key: 7,
  class: "info-row"
}, hu = { class: "info-value" }, pu = { class: "expandable-list" }, gu = { class: "expanded-content" }, mu = {
  key: 8,
  class: "info-row"
}, bu = { class: "info-value" }, yu = { class: "expandable-list" }, _u = { class: "expanded-content" }, wu = { class: "go-term-list" }, vu = {
  key: 0,
  class: "info-section sequences"
}, xu = {
  key: 0,
  class: "info-row"
}, Su = { class: "info-value" }, Cu = { class: "expandable-list" }, Eu = { class: "expanded-content" }, Tu = { class: "sequence-text" }, Au = {
  key: 1,
  class: "sequence-block"
}, ku = { class: "sequence-header" }, Ru = { class: "sequence-text" }, Ou = {
  key: 1,
  class: "no-selection"
};
function Pu(e, t, n, s, r, i) {
  var l, a, u, c, f, p, g, m, x, A, H, q, W, I;
  const o = qn("SortableTable");
  return n.feature ? ($(), B("div", Uc, [
    M("div", qc, [
      M("div", Vc, [
        M("h3", null, se(n.feature.type), 1)
      ]),
      M("div", Wc, [
        (a = (l = n.feature.qualifiers) == null ? void 0 : l.locus_tag) != null && a[0] ? ($(), B("div", zc, [
          t[2] || (t[2] = M("span", { class: "info-label" }, "Locus tag", -1)),
          M("span", Kc, se(n.feature.qualifiers.locus_tag[0]), 1)
        ])) : be("", !0),
        (c = (u = n.feature.qualifiers) == null ? void 0 : u.protein_id) != null && c[0] ? ($(), B("div", Gc, [
          t[3] || (t[3] = M("span", { class: "info-label" }, "Protein ID", -1)),
          M("span", Jc, se(n.feature.qualifiers.protein_id[0]), 1)
        ])) : be("", !0),
        (p = (f = n.feature.qualifiers) == null ? void 0 : f.gene) != null && p[0] ? ($(), B("div", Yc, [
          t[4] || (t[4] = M("span", { class: "info-label" }, "Gene", -1)),
          M("span", Xc, se(n.feature.qualifiers.gene[0]), 1)
        ])) : be("", !0),
        (m = (g = n.feature.qualifiers) == null ? void 0 : g.product) != null && m[0] ? ($(), B("div", Qc, [
          t[5] || (t[5] = M("span", { class: "info-label" }, "Product", -1)),
          M("span", Zc, se(n.feature.qualifiers.product[0]), 1)
        ])) : be("", !0),
        n.feature.location ? ($(), B("div", eu, [
          t[6] || (t[6] = M("span", { class: "info-label" }, "Location", -1)),
          M("span", tu, se(s.formatLocation(n.feature.location)), 1)
        ])) : be("", !0),
        (A = (x = n.feature.qualifiers) == null ? void 0 : x.gene_kind) != null && A[0] ? ($(), B("div", nu, [
          t[7] || (t[7] = M("span", { class: "info-label" }, "Gene kind", -1)),
          M("span", su, se(n.feature.qualifiers.gene_kind[0]), 1)
        ])) : be("", !0),
        ($(!0), B(_e, null, Xe(s.getAdditionalQualifiers(), (J, de) => ($(), B("div", {
          key: de,
          class: "info-row"
        }, [
          M("span", ru, se(s.formatQualifierKey(de)), 1),
          M("span", iu, [
            Array.isArray(J) && J.length > 1 ? ($(), B("details", ou, [
              M("summary", null, se(J.length) + " items", 1),
              M("div", lu, [
                ($(), mn(Qi(s.renderQualifierValue(de, J))))
              ])
            ])) : ($(), B(_e, { key: 1 }, [
              St(se(s.formatQualifierValue(J)), 1)
            ], 64))
          ])
        ]))), 128)),
        s.pfamDomains.length > 0 ? ($(), B("div", au, [
          t[8] || (t[8] = M("span", { class: "info-label" }, "Pfam hits", -1)),
          M("span", cu, [
            M("details", uu, [
              M("summary", null, se(s.pfamDomains.length) + " items", 1),
              M("div", fu, [
                Te(o, {
                  headers: s.pfamTableHeaders,
                  rows: s.pfamTableRows
                }, null, 8, ["headers", "rows"])
              ])
            ])
          ])
        ])) : be("", !0),
        n.feature.type === "CDS" && s.mibigEntries.length > 0 ? ($(), B("div", du, [
          t[9] || (t[9] = M("span", { class: "info-label" }, "MiBIG hits", -1)),
          M("span", hu, [
            M("details", pu, [
              M("summary", null, se(s.mibigEntries.length) + " items", 1),
              M("div", gu, [
                Te(o, {
                  headers: s.mibigTableHeaders,
                  rows: s.mibigTableRows
                }, null, 8, ["headers", "rows"])
              ])
            ])
          ])
        ])) : be("", !0),
        s.goTerms.length > 0 ? ($(), B("div", mu, [
          t[10] || (t[10] = M("span", { class: "info-label" }, "GO terms", -1)),
          M("span", bu, [
            M("details", yu, [
              M("summary", null, se(s.goTerms.length) + " items", 1),
              M("div", _u, [
                M("ul", wu, [
                  ($(!0), B(_e, null, Xe(s.goTerms, (J, de) => ($(), B("li", {
                    key: de,
                    class: "go-term"
                  }, se(J), 1))), 128))
                ])
              ])
            ])
          ])
        ])) : be("", !0)
      ]),
      n.feature.type === "CDS" ? ($(), B("div", vu, [
        (q = (H = n.feature.qualifiers) == null ? void 0 : H.translation) != null && q[0] ? ($(), B("div", xu, [
          t[12] || (t[12] = M("span", { class: "info-label" }, "AA sequence", -1)),
          M("span", Su, [
            M("details", Cu, [
              M("summary", null, [
                t[11] || (t[11] = St(" Show sequence ", -1)),
                M("button", {
                  onClick: t[0] || (t[0] = To((J) => s.copyToClipboard(n.feature.qualifiers.translation[0], "amino acid"), ["stop"])),
                  class: "copy-button copy-button-inline"
                }, " Copy to clipboard ")
              ]),
              M("div", Eu, [
                M("pre", Tu, se(s.formatSequence(n.feature.qualifiers.translation[0])), 1)
              ])
            ])
          ])
        ])) : be("", !0),
        (I = (W = n.feature.qualifiers) == null ? void 0 : W.nucleotide) != null && I[0] ? ($(), B("div", Au, [
          M("div", ku, [
            t[13] || (t[13] = M("span", { class: "sequence-label" }, "Nucleotide sequence", -1)),
            M("button", {
              onClick: t[1] || (t[1] = (J) => s.copyToClipboard(n.feature.qualifiers.nucleotide[0], "nucleotide")),
              class: "copy-button"
            }, " Copy to clipboard ")
          ]),
          M("pre", Ru, se(s.formatSequence(n.feature.qualifiers.nucleotide[0])), 1)
        ])) : be("", !0)
      ])) : be("", !0)
    ])
  ])) : ($(), B("div", Ou, [...t[14] || (t[14] = [
    M("p", null, "Click on a feature to view details", -1)
  ])]));
}
const Lu = /* @__PURE__ */ zt(Hc, [["render", Pu], ["__scopeId", "data-v-4af88c98"]]), Fu = {
  name: "SimpleDetails",
  props: {
    // The data object to display
    data: {
      type: Object,
      default: null
    },
    // The type of element for the header
    elementType: {
      type: String,
      default: ""
    }
  },
  setup(e) {
    const t = We(() => {
      if (!e.data) return {};
      const i = {}, o = ["_elementType"];
      return Object.keys(e.data).forEach((l) => {
        o.includes(l) || (i[l] = e.data[l]);
      }), i;
    }), n = We(() => {
      if (!e.data) return [];
      const i = ["_elementType"];
      return Object.keys(e.data).filter((o) => !i.includes(o));
    });
    return {
      displayData: t,
      displayKeys: n,
      formatLabel: (i) => {
        const o = i.replace(/_/g, " ").replace(/([A-Z])/g, " $1").split(" ").filter((l) => l.length > 0);
        return o.length === 0 ? "" : o[0].charAt(0).toUpperCase() + o[0].slice(1).toLowerCase() + (o.length > 1 ? " " + o.slice(1).map((l) => l.toLowerCase()).join(" ") : "");
      },
      formatValue: (i, o) => o == null ? "N/A" : typeof o == "number" ? i.toLowerCase().includes("evalue") && o < 0.01 ? o.toExponential(2) : i.toLowerCase().includes("score") || i.toLowerCase().includes("bitscore") ? o.toFixed(1) : i.toLowerCase().includes("percent") || i.toLowerCase().includes("identity") || i.toLowerCase().includes("coverage") ? `${o.toFixed(1)}%` : o.toLocaleString() : Array.isArray(o) ? o.length === 0 ? "None" : o.length === 1 ? String(o[0]) : o.join(", ") : typeof o == "boolean" ? o ? "Yes" : "No" : typeof o == "object" ? JSON.stringify(o, null, 2) : String(o)
    };
  }
}, Mu = {
  key: 0,
  class: "simple-details"
}, Iu = {
  key: 0,
  class: "details-header"
}, Nu = { class: "details-content" }, Du = { class: "info-section" }, $u = { class: "info-label" }, ju = { class: "info-value" }, Bu = {
  key: 1,
  class: "no-selection"
};
function Hu(e, t, n, s, r, i) {
  return n.data ? ($(), B("div", Mu, [
    n.elementType ? ($(), B("div", Iu, [
      M("h3", null, se(n.elementType), 1)
    ])) : be("", !0),
    M("div", Nu, [
      M("div", Du, [
        ($(!0), B(_e, null, Xe(s.displayKeys, (o) => ($(), B("div", {
          key: o,
          class: "info-row"
        }, [
          M("span", $u, se(s.formatLabel(o)), 1),
          M("span", ju, se(s.formatValue(o, s.displayData[o])), 1)
        ]))), 128))
      ])
    ])
  ])) : ($(), B("div", Bu, [...t[0] || (t[0] = [
    M("p", null, "Click on an element to view details", -1)
  ])]));
}
const Uu = /* @__PURE__ */ zt(Fu, [["render", Hu], ["__scopeId", "data-v-6d8f9976"]]), qu = {
  name: "RegionViewerComponent",
  components: {
    FeatureDetails: Lu,
    SimpleDetails: Uu
  },
  props: {
    // Current record information
    recordInfo: {
      type: Object,
      default: null
      // Expected shape: { recordId, filename, recordInfo: { description } }
    },
    // Available regions for the current record
    regions: {
      type: Array,
      default: () => []
      // Expected shape: [{ id, region_number, product }]
    },
    // Features to display
    features: {
      type: Array,
      default: () => []
      // Expected shape: [{ type, location, qualifiers }]
    },
    // Region boundaries (optional, for region-specific view)
    regionBoundaries: {
      type: Object,
      default: null
      // Expected shape: { start, end }
    },
    // PFAM color mapping
    pfamColorMap: {
      type: Object,
      default: () => ({})
      // Expected shape: { 'PF00001': '#FF0000', ... }
    },
    // Selected region ID (controlled from parent)
    selectedRegionId: {
      type: String,
      default: ""
    },
    // Data provider for fetching additional data
    dataProvider: {
      type: Object,
      default: null
    },
    // TFBS binding site hits
    tfbsHits: {
      type: Array,
      default: () => []
      // Expected shape: [{ name, start, species, link, description, consensus, confidence, strand, score, max_score }]
    },
    // TTA codon positions
    ttaCodons: {
      type: Array,
      default: () => []
      // Expected shape: [{ start, strand }]
    },
    // Resistance features
    resistanceFeatures: {
      type: Array,
      default: () => []
      // Expected shape: [{ locus_tag, query_id, reference_id, subfunctions, description, bitscore, evalue, query_start, query_end }]
    }
  },
  emits: [
    "region-changed",
    // Emitted when user selects a different region
    "annotation-clicked",
    // Emitted when user clicks an annotation
    "annotation-hovered",
    // Emitted when user hovers over an annotation
    "error"
    // Emitted when an error occurs
  ],
  setup(e, { emit: t, expose: n }) {
    const s = he(null), r = he(""), i = he(!1), o = he(""), l = he([]), a = he([]), u = he(!1), c = We(() => {
      if (!r.value || !e.regions.length)
        return "1";
      const b = e.regions.find((_) => _.id === r.value);
      return b ? String(b.region_number) : "1";
    });
    let f = null, p = {};
    const g = he(null), m = We(() => {
      var b;
      return ((b = g.value) == null ? void 0 : b.data) || null;
    });
    Ke(() => e.features, () => {
      e.features && e.features.length > 0 && x();
    }), Ke(() => e.tfbsHits, () => {
      e.features && e.features.length > 0 && x();
    }), Ke(() => e.ttaCodons, () => {
      e.features && e.features.length > 0 && x();
    }), Ke(() => e.resistanceFeatures, () => {
      e.features && e.features.length > 0 && x();
    }), Ke(() => e.selectedRegionId, (b) => {
      r.value = b;
    }, { immediate: !0 }), fr(() => {
      document.addEventListener("click", te), e.features && e.features.length > 0 && x();
    }), dr(() => {
      document.removeEventListener("click", te), f && (f.destroy(), f = null);
    });
    const x = async () => {
      try {
        if (o.value = "", !e.features || e.features.length === 0) {
          console.warn("No features provided");
          return;
        }
        console.log("Building viewer with", e.features.length, "features"), g.value = null, de(), console.log("Built tracks:", Object.keys(p));
        const b = Object.values(p).map((_) => ({
          id: _.id,
          label: _.label,
          annotationCount: _.annotations.length
        }));
        P(b), l.value = b, e.regionBoundaries ? a.value = b.filter(
          (_) => ["CDS"].includes(_.id) || _.id.includes("protocluster") || _.id.includes("PFAM_domain") || _.id.includes("cand_cluster")
        ).map((_) => _.id) : a.value = b.map((_) => _.id), console.log("Available tracks:", l.value), console.log("Selected tracks:", a.value), await rs(), console.log("Initializing viewer..."), H(), ie(), console.log("Viewer initialized and updated");
      } catch (b) {
        console.error("Error in rebuildViewer:", b), o.value = `Failed to build viewer: ${b.message}`, t("error", b);
      }
    }, A = () => {
      t("region-changed", r.value);
    }, H = () => {
      if (f && f.destroy(), !s.value) return;
      let b, _, R;
      if (e.regionBoundaries)
        b = e.regionBoundaries.start, _ = e.regionBoundaries.end, R = (_ - b) * 0.02, console.log("Using region boundaries:", b, "-", _);
      else {
        console.log("Calculating boundaries from", e.features.length, "features");
        const O = e.features.filter((j) => j.location).map((j) => {
          const re = j.location.match(/\[<?(\d+):>?(\d+)\]/);
          return re ? [parseInt(re[1]), parseInt(re[2])] : null;
        }).filter(Boolean).flat();
        console.log("Extracted positions:", O.slice(0, 10), O.length > 10 ? `... (${O.length} total)` : ""), O.length === 0 ? (console.warn("No valid positions found, using default domain"), b = 0, _ = 1e3) : (b = Math.min(...O), _ = Math.max(...O)), R = (_ - b) * 0.02, console.log("Calculated domain:", b - R, "to", _ + R);
      }
      console.log("Creating TrackViewer..."), f = new Vs({
        container: s.value,
        // width is not specified, so it will be responsive
        height: 400,
        domain: [b - R, _ + R],
        trackHeight: 40,
        showTrackLabels: !1,
        onAnnotationClick: (O, j) => {
          console.log("Clicked annotation:", O, "on track:", j), oe(O), t("annotation-clicked", { annotation: O, track: j });
        },
        onAnnotationHover: (O, j, re) => {
          t("annotation-hovered", { annotation: O, track: j, event: re });
        },
        onBackgroundClick: () => {
          g.value = null, X(), ie();
        }
      }), console.log("TrackViewer created successfully");
    }, q = (b, _, R) => {
      p[b] || (p[b] = {
        id: b,
        label: _,
        height: R,
        annotations: [],
        primitives: []
      }, b === "CDS" && p[b].primitives.push({
        id: "cds-baseline",
        trackId: b,
        type: "horizontal-line",
        class: "cds-baseline",
        fy: 0.5,
        stroke: "black",
        opacity: 1
      }));
    }, W = (b) => {
      if (!b) return null;
      const _ = b.match(/\[<?(\d+):>?(\d+)\](?:\(([+-])\))?/);
      if (!_) return null;
      const R = _[3] || null;
      return {
        start: parseInt(_[1]),
        end: parseInt(_[2]),
        strand: R,
        direction: R === "+" ? "right" : R === "-" ? "left" : "none"
      };
    }, I = (b) => b.toLowerCase().replace("_", " "), J = (b) => b.toLowerCase().replace(/[^a-z0-9]+/g, "-"), de = () => {
      if (p = {}, e.features.forEach((b) => {
        var re, Ae, Je, et, Mt, Jt, bt, ps, d, h, w, T, S, C, F, L, k, E, V, N, U, Y, ne, fe, ae, Ce, Ee, Be, He;
        if (!b.location) return;
        const _ = W(b.location);
        if (!_) {
          console.warn("Failed to parse location for feature:", b.type, b.location);
          return;
        }
        const R = [];
        R.push(Se(b.type));
        let O, j;
        switch (b.type) {
          case "cand_cluster":
            if ((((Ae = (re = b.qualifiers) == null ? void 0 : re.protoclusters) == null ? void 0 : Ae.length) || 1) < 2) break;
            const tt = ((et = (Je = b.qualifiers) == null ? void 0 : Je.candidate_cluster_number) == null ? void 0 : et[0]) || "unknown", Yt = ((Jt = (Mt = b.qualifiers) == null ? void 0 : Mt.kind) == null ? void 0 : Jt[0]) || "unknown";
            O = `cand_cluster-${tt}`, j = `Candidate Cluster ${tt}`, R.push(`candidate-${J(Yt)}`), q(O, j), p[O].annotations.push({
              id: `${b.type}-${tt}`,
              trackId: O,
              type: "box",
              heightFraction: 0.4,
              classes: R,
              label: `CC ${tt}: ${I(Yt)}`,
              labelPosition: "center",
              showLabel: "always",
              start: _.start,
              end: _.end,
              data: b
            });
            break;
          case "PFAM_domain":
            O = b.type, j = b.type, q(O, j);
            const ke = ((d = (ps = (bt = b.qualifiers) == null ? void 0 : bt.db_xref) == null ? void 0 : ps[0]) == null ? void 0 : d.replace("PFAM:", "")) || ((S = (T = (w = (h = b.qualifiers) == null ? void 0 : h.inference) == null ? void 0 : w[0]) == null ? void 0 : T.match(/PFAM:([^,\s]+)/)) == null ? void 0 : S[1]) || ((k = (L = (F = (C = b.qualifiers) == null ? void 0 : C.note) == null ? void 0 : F[0]) == null ? void 0 : L.match(/PF\d+/)) == null ? void 0 : k[0]), [Fe, gs] = ke ? ke.split(".") : ["", ""], Xt = Fe && e.pfamColorMap[Fe] ? e.pfamColorMap[Fe] : null, el = ((V = (E = b.qualifiers) == null ? void 0 : E.description) == null ? void 0 : V[0]) || "", tl = `${Fe} ${el}`.trim(), nl = {
              id: `${b.type}-${_.start}-${_.end}`,
              trackId: O,
              type: "box",
              classes: R,
              label: tl || "unknown",
              start: _.start,
              end: _.end,
              fill: Xt,
              stroke: Xt,
              data: b
            };
            p[O].annotations.push(nl);
            break;
          case "CDS":
            O = b.type, j = b.type, R.push(`gene-type-${((U = (N = b.qualifiers) == null ? void 0 : N.gene_kind) == null ? void 0 : U[0]) || "other"}`), q(O, j), p[O].annotations.push({
              id: `${b.type}-${_.start}-${_.end}`,
              trackId: O,
              type: "arrow",
              classes: R,
              label: v(b),
              start: _.start,
              end: _.end,
              direction: _.direction,
              stroke: "black",
              data: b
            });
            break;
          case "protocluster":
            const An = ((ne = (Y = b.qualifiers) == null ? void 0 : Y.protocluster_number) == null ? void 0 : ne[0]) || "unknown", sl = ((ae = (fe = b.qualifiers) == null ? void 0 : fe.category) == null ? void 0 : ae[0]) || "unknown", rl = ((Ee = (Ce = b.qualifiers) == null ? void 0 : Ce.product) == null ? void 0 : Ee[0]) || "unknown";
            R.push(sl), R.push(rl);
            const ms = W(((He = (Be = b.qualifiers) == null ? void 0 : Be.core_location) == null ? void 0 : He[0]) || null);
            O = `protocluster-track-${An}`, j = `Protocluster track ${An}`;
            for (let bs of Object.keys(p)) {
              if (!bs.startsWith("protocluster-track-")) continue;
              if (!p[bs].annotations.some((xr) => !(_.end < xr.start || _.start > xr.end))) {
                O = bs;
                break;
              }
            }
            q(O, j), p[O].annotations.push({
              id: `${b.type}-${An}`,
              trackId: O,
              type: "box",
              heightFraction: 0.3,
              classes: R,
              start: _.start,
              end: _.end,
              stroke: "none",
              opacity: 0.5,
              data: b
            }), ms && p[O].annotations.push({
              id: `${b.type}-${An}-core`,
              trackId: O,
              type: "box",
              heightFraction: 0.35,
              classes: [...R, "proto-core"],
              label: v(b),
              showLabel: "always",
              start: ms.start,
              end: ms.end,
              stroke: "black",
              data: b
            });
            break;
          default:
            O = b.type, j = b.type, q(O, j), p[O].annotations.push({
              id: `${b.type}-${_.start}-${_.end}`,
              trackId: O,
              type: "box",
              classes: R,
              label: v(b),
              start: _.start,
              end: _.end,
              data: b
            });
            break;
        }
      }), e.tfbsHits && e.tfbsHits.length > 0) {
        const b = "CDS";
        q(b, b), e.tfbsHits.forEach((_, R) => {
          _.confidence !== "strong" && _.confidence !== "medium" || p[b].annotations.push({
            id: `tfbs-${R}-${_.start}`,
            trackId: b,
            type: "pin",
            classes: ["tfbs-hit", `tfbs-${_.confidence}`],
            label: _.name,
            labelPosition: "above",
            showLabel: "hover",
            start: _.start,
            end: _.start,
            // Pins are at a single position
            fy: 0.5,
            // Middle of the track
            heightFraction: 0.5,
            opacity: 0.8,
            data: { ..._, _elementType: "Binding Site" }
          });
        });
      }
      if (e.ttaCodons && e.ttaCodons.length > 0) {
        const b = "CDS";
        q(b, b), e.ttaCodons.forEach((_, R) => {
          p[b].annotations.push({
            id: `tta-${R}-${_.start}`,
            trackId: b,
            type: "triangle",
            classes: ["tta-codon"],
            showLabel: "none",
            start: _.start,
            end: _.start,
            // Triangles are at a single position
            fy: 0.2,
            // Middle of the track
            heightFraction: 0.4,
            opacity: 0.9,
            data: { ..._, _elementType: "TTA Codon" }
          });
        });
      }
      if (e.resistanceFeatures && e.resistanceFeatures.length > 0) {
        const b = "CDS";
        q(b, b);
        const _ = {};
        e.features.forEach((R) => {
          var O, j;
          if (R.type === "CDS" && ((j = (O = R.qualifiers) == null ? void 0 : O.locus_tag) != null && j[0])) {
            const re = R.qualifiers.locus_tag[0], Ae = W(R.location);
            Ae && (_[re] = Ae);
          }
        }), e.resistanceFeatures.forEach((R, O) => {
          const j = R.locus_tag, re = _[j];
          console.log("Mapping resistance feature", R.reference_id, "to locus tag", j, "with location", re), re && p[b].annotations.push({
            id: `resistance-${O}-${j}`,
            trackId: b,
            type: "box",
            classes: ["resistance"],
            showLabel: "never",
            start: re.start,
            end: re.end,
            fy: 0.1,
            heightFraction: 0.2,
            fill: "#BBB",
            data: { ...R, _elementType: "Resistance" }
          });
        });
      }
    }, ie = () => {
      if (!f || !Object.keys(p).length) return;
      const b = [];
      l.value.forEach((j) => {
        a.value.includes(j.id) && p[j.id] && b.push(p[j.id]);
      });
      const _ = b.map((j) => ({
        id: j.id,
        label: j.label,
        height: j.height || void 0
      })), R = b.flatMap((j) => j.annotations), O = b.flatMap((j) => j.primitives);
      f.setData({ tracks: _, annotations: R, primitives: O });
    }, Se = (b) => `feature-${b.toLowerCase().replace(/[^a-z0-9]+/g, "-")}`, v = (b) => {
      var R, O, j, re, Ae;
      const _ = b.qualifiers || {};
      return (R = _.locus_tag) != null && R[0] ? _.locus_tag[0] : (O = _.gene) != null && O[0] ? _.gene[0] : (j = _.product) != null && j[0] ? _.product[0] : (re = _.description) != null && re[0] ? _.description[0] : (Ae = _.db_xref) != null && Ae[0] ? _.db_xref[0] : b.type || "Feature";
    }, P = (b) => {
      const _ = (R) => R.startsWith("cand_cluster") ? 1 : R.startsWith("protocluster") ? 2 : R.startsWith("CDS") ? 3 : R.startsWith("PFAM_domain") ? 4 : 5;
      return b.sort((R, O) => {
        const j = _(R.id), re = _(O.id);
        return j !== re ? j - re : R.id.localeCompare(O.id);
      });
    }, D = () => {
      u.value = !u.value;
    }, Q = (b) => {
      b.target.checked ? a.value = [...l.value.map((_) => _.id)] : a.value = [], ie();
    }, te = (b) => {
      b.target.closest(".multi-select-dropdown") || (u.value = !1);
    }, oe = (b, _) => {
      if (g.value = b, b.id.endsWith("-core")) {
        const R = p[b.trackId], O = b.id.replace("-core", "");
        if (R) {
          const j = R.annotations.find(
            (re) => re.id === O
          );
          j && (g.value = j);
        }
      }
      X(), ie();
    }, le = () => {
      g.value = null;
    }, z = (b) => {
      if (!g.value)
        return !0;
      if (b.id.endsWith("-core")) {
        const _ = p[b.trackId], R = b.id.replace("-core", "");
        if (_) {
          const O = _.annotations.find(
            (j) => j.id === R
          );
          if (O)
            return O.start >= g.value.start && O.end <= g.value.end;
        }
        return !1;
      }
      return b.start >= g.value.start && b.end <= g.value.end;
    }, X = () => {
      Object.values(p).forEach((b) => {
        b.id !== "CDS" && b.id !== "PFAM_domain" || (b.annotations.forEach((_) => {
          const R = z(_);
          _._originalOpacity === void 0 && (_._originalOpacity = _.opacity !== void 0 ? _.opacity : 1), R ? _.opacity = _._originalOpacity : _.opacity = _._originalOpacity * 0.5;
        }), b.id === "CDS" && b.primitives && b.primitives.forEach((_) => {
          _.id === "cds-baseline" && (_._originalOpacity === void 0 && (_._originalOpacity = _.opacity !== void 0 ? _.opacity : 1), g.value ? _.opacity = _._originalOpacity * 0.1 : _.opacity = _._originalOpacity);
        }));
      });
    };
    return n({
      clearViewer: () => {
        r.value = "", g.value = null, p = {}, l.value = [], a.value = [], o.value = "", s.value && (s.value.innerHTML = ""), f = null;
      },
      rebuildViewer: x
    }), {
      viewerContainer: s,
      selectedRegion: r,
      loading: i,
      error: o,
      availableTracks: l,
      selectedTracks: a,
      dropdownOpen: u,
      selectedElement: m,
      currentRegionNumber: c,
      onRegionChange: A,
      updateViewer: ie,
      toggleDropdown: D,
      toggleSelectAll: Q,
      clearSelectedElement: le
    };
  }
}, Vu = { class: "region-viewer-wrapper" }, Wu = {
  key: 0,
  class: "controls"
}, zu = ["value"], Ku = {
  key: 1,
  class: "no-regions-message"
}, Gu = {
  key: 2,
  class: "feature-controls"
}, Ju = { class: "multi-select-container" }, Yu = { class: "selected-display" }, Xu = { key: 0 }, Qu = { key: 1 }, Zu = { key: 2 }, ef = { class: "select-all-option" }, tf = ["checked", "indeterminate"], nf = ["value"], sf = {
  key: 1,
  class: "current-record-info"
}, rf = { class: "record-details" }, of = {
  key: 0,
  class: "description"
}, lf = {
  ref: "viewerContainer",
  class: "viewer-container"
}, af = {
  key: 2,
  class: "loading"
}, cf = {
  key: 3,
  class: "error"
}, uf = {
  key: 4,
  class: "feature-details-container"
};
function ff(e, t, n, s, r, i) {
  var a;
  const o = qn("FeatureDetails"), l = qn("SimpleDetails");
  return $(), B("div", Vu, [
    n.recordInfo ? ($(), B("div", Wu, [
      n.regions.length > 0 ? Ss(($(), B("select", {
        key: 0,
        "onUpdate:modelValue": t[0] || (t[0] = (u) => s.selectedRegion = u),
        onChange: t[1] || (t[1] = (...u) => s.onRegionChange && s.onRegionChange(...u)),
        class: "region-select"
      }, [
        t[7] || (t[7] = M("option", { value: "" }, "Show all features", -1)),
        ($(!0), B(_e, null, Xe(n.regions, (u) => ($(), B("option", {
          key: u.id,
          value: u.id
        }, " Region " + se(u.region_number) + " - " + se(u.product.join(", ")), 9, zu))), 128))
      ], 544)), [
        [xc, s.selectedRegion]
      ]) : be("", !0),
      n.regions.length === 0 && !s.loading ? ($(), B("div", Ku, " No regions found - showing all features for this record ")) : be("", !0),
      s.availableTracks.length > 0 ? ($(), B("div", Gu, [
        M("div", Ju, [
          M("div", {
            class: Lt(["multi-select-dropdown", { open: s.dropdownOpen }]),
            onClick: t[6] || (t[6] = (...u) => s.toggleDropdown && s.toggleDropdown(...u))
          }, [
            M("div", Yu, [
              s.selectedTracks.length === s.availableTracks.length ? ($(), B("span", Xu, " All tracks (" + se(s.selectedTracks.length) + ") ", 1)) : s.selectedTracks.length === 0 ? ($(), B("span", Qu, " No tracks selected ")) : ($(), B("span", Zu, se(s.selectedTracks.length) + " tracks selected ", 1)),
              t[8] || (t[8] = M("span", { class: "dropdown-arrow" }, "▼", -1))
            ]),
            s.dropdownOpen ? ($(), B("div", {
              key: 0,
              class: "dropdown-options",
              onClick: t[5] || (t[5] = To(() => {
              }, ["stop"]))
            }, [
              M("div", ef, [
                M("label", null, [
                  M("input", {
                    type: "checkbox",
                    checked: s.selectedTracks.length === s.availableTracks.length,
                    indeterminate: s.selectedTracks.length > 0 && s.selectedTracks.length < s.availableTracks.length,
                    onChange: t[2] || (t[2] = (...u) => s.toggleSelectAll && s.toggleSelectAll(...u))
                  }, null, 40, tf),
                  t[9] || (t[9] = St(" Select All ", -1))
                ])
              ]),
              t[10] || (t[10] = M("div", { class: "option-separator" }, null, -1)),
              ($(!0), B(_e, null, Xe(s.availableTracks, (u) => ($(), B("div", {
                key: u.id,
                class: "dropdown-option"
              }, [
                M("label", null, [
                  Ss(M("input", {
                    type: "checkbox",
                    value: u.id,
                    "onUpdate:modelValue": t[3] || (t[3] = (c) => s.selectedTracks = c),
                    onChange: t[4] || (t[4] = (...c) => s.updateViewer && s.updateViewer(...c))
                  }, null, 40, nf), [
                    [vc, s.selectedTracks]
                  ]),
                  St(" " + se(u.label) + " (" + se(u.annotationCount) + ") ", 1)
                ])
              ]))), 128))
            ])) : be("", !0)
          ], 2)
        ])
      ])) : be("", !0)
    ])) : be("", !0),
    n.recordInfo ? ($(), B("div", sf, [
      M("span", null, "Current Record: " + se(n.recordInfo.recordId) + " (" + se(n.recordInfo.filename) + ")", 1),
      M("div", rf, [
        (a = n.recordInfo.recordInfo) != null && a.description ? ($(), B("span", of, se(n.recordInfo.recordInfo.description), 1)) : be("", !0)
      ])
    ])) : be("", !0),
    Ss(M("div", lf, null, 512), [
      [oc, n.recordInfo]
    ]),
    s.loading ? ($(), B("div", af, " Loading region data... ")) : be("", !0),
    s.error ? ($(), B("div", cf, se(s.error), 1)) : be("", !0),
    n.recordInfo ? ($(), B("div", uf, [
      s.selectedElement && s.selectedElement.type && s.selectedElement.qualifiers ? ($(), mn(o, {
        key: 0,
        feature: s.selectedElement,
        "all-features": n.features,
        "data-provider": n.dataProvider,
        "record-info": n.recordInfo,
        "region-number": s.currentRegionNumber,
        onClose: s.clearSelectedElement
      }, null, 8, ["feature", "all-features", "data-provider", "record-info", "region-number", "onClose"])) : s.selectedElement ? ($(), mn(l, {
        key: 1,
        data: s.selectedElement,
        "element-type": s.selectedElement._elementType
      }, null, 8, ["data", "element-type"])) : be("", !0)
    ])) : be("", !0)
  ]);
}
const ko = /* @__PURE__ */ zt(qu, [["render", ff], ["__scopeId", "data-v-6b91567d"]]), df = {
  name: "RegionViewerContainer",
  components: {
    RegionViewer: ko
  },
  props: {
    // Data provider instance (optional, defaults to BGCViewerAPIProvider)
    dataProvider: {
      type: Object,
      default: null
    },
    // Record ID to load (can be changed dynamically)
    recordId: {
      type: String,
      default: ""
    },
    // Full record data from RecordListSelector (includes entryId, recordId, filename for API provider)
    recordData: {
      type: Object,
      default: null
    },
    // Initial region ID to select (optional)
    initialRegionId: {
      type: String,
      default: ""
    }
  },
  emits: [
    "region-changed",
    "annotation-clicked",
    "error"
  ],
  setup(e, { emit: t }) {
    const n = he(null), s = he(null), r = he([]), i = he([]), o = he(null), l = he({}), a = he([]), u = he([]), c = he([]), f = he(""), p = he(!1), g = he(""), m = () => {
      s.value = null, r.value = [], i.value = [], o.value = null, a.value = [], u.value = [], c.value = [], f.value = "";
    }, x = async () => {
      if (n.value)
        try {
          l.value = await n.value.getPfamColorMap(), console.log("Loaded PFAM colors for", Object.keys(l.value).length, "domains");
        } catch (v) {
          console.warn("Failed to load PFAM color mapping:", v.message);
        }
    }, A = async (v, P) => {
      if (n.value)
        try {
          const D = P.replace("region_", ""), Q = await n.value.getTFBSHits(v, D);
          a.value = Q.hits || [], console.log("Loaded", a.value.length, "TFBS binding sites for region", D);
        } catch (D) {
          console.warn("Failed to load TFBS hits:", D.message), a.value = [];
        }
    }, H = async (v) => {
      if (n.value)
        try {
          const P = await n.value.getTTACodons(v);
          u.value = P.codons || [], console.log("Loaded", u.value.length, "TTA codons for record", v);
        } catch (P) {
          console.warn("Failed to load TTA codons:", P.message), u.value = [];
        }
    }, q = async (v) => {
      if (n.value)
        try {
          const P = await n.value.getResistanceFeatures(v);
          c.value = P.features || [], console.log("Loaded", c.value.length, "resistance features for record", v);
        } catch (P) {
          console.warn("Failed to load resistance features:", P.message), c.value = [];
        }
    }, W = async (v) => {
      var P;
      p.value = !0, g.value = "", f.value = "";
      try {
        if (!n.value)
          throw new Error("No data provider available");
        const D = ((P = e.recordData) == null ? void 0 : P.entryId) || v, Q = await n.value.loadEntry(D);
        s.value = {
          recordId: Q.recordId,
          filename: Q.filename,
          recordInfo: Q.recordInfo
        };
        const te = await n.value.getRegions(s.value.recordId);
        r.value = te.regions || [], o.value = te.boundaries || null, await H(s.value.recordId), await q(s.value.recordId), r.value && r.value.length > 0 ? (f.value = r.value[0].id, await I(s.value.recordId, f.value)) : await J(s.value.recordId), p.value = !1;
      } catch (D) {
        console.error("Error loading record:", D), g.value = `Failed to load record: ${D.message}`, t("error", g.value), p.value = !1;
      }
    }, I = async (v, P) => {
      try {
        const D = await n.value.getRegionFeatures(v, P);
        i.value = D.features || [], o.value = D.region_boundaries || null, await A(v, P);
      } catch (D) {
        console.error("Error loading region features:", D), g.value = `Failed to load region features: ${D.message}`, t("error", g.value);
      }
    }, J = async (v) => {
      try {
        const P = await n.value.getRecordFeatures(v);
        i.value = P.features || [], o.value = null, a.value = [];
      } catch (P) {
        console.error("Error loading features:", P), g.value = `Failed to load features: ${P.message}`, t("error", g.value);
      }
    }, de = async (v) => {
      f.value = v, v ? await I(s.value.recordId, v) : await J(s.value.recordId), t("region-changed", v);
    }, ie = (v) => {
      t("annotation-clicked", v);
    }, Se = (v) => {
      t("error", v);
    };
    return fr(async () => {
      if (!e.dataProvider) {
        g.value = "No data provider specified. Please set the dataProvider property.", t("error", g.value);
        return;
      }
      n.value = e.dataProvider, await x();
    }), Ke(() => e.recordId, async (v, P) => {
      v && n.value ? await W(v) : !v && P !== void 0 && m();
    }, { immediate: !0 }), Ke(() => e.dataProvider, (v) => {
      v && (n.value = v, x());
    }), Ke(() => e.initialRegionId, (v) => {
      v && r.value.length > 0 && (f.value = v, s.value && I(s.value.recordId, v));
    }), {
      recordInfo: s,
      regions: r,
      features: i,
      regionBoundaries: o,
      pfamColorMap: l,
      tfbsHits: a,
      ttaCodons: u,
      resistanceFeatures: c,
      selectedRegionId: f,
      loading: p,
      error: g,
      handleRegionChanged: de,
      handleAnnotationClicked: ie,
      handleError: Se
    };
  }
}, hf = { class: "region-viewer-container" };
function pf(e, t, n, s, r, i) {
  const o = qn("RegionViewer");
  return $(), B("div", hf, [
    Te(o, {
      "record-info": s.recordInfo,
      regions: s.regions,
      features: s.features,
      "region-boundaries": s.regionBoundaries,
      "pfam-color-map": s.pfamColorMap,
      "selected-region-id": s.selectedRegionId,
      "data-provider": n.dataProvider,
      "tfbs-hits": s.tfbsHits,
      "tta-codons": s.ttaCodons,
      "resistance-features": s.resistanceFeatures,
      onRegionChanged: s.handleRegionChanged,
      onAnnotationClicked: s.handleAnnotationClicked,
      onError: s.handleError
    }, null, 8, ["record-info", "regions", "features", "region-boundaries", "pfam-color-map", "selected-region-id", "data-provider", "tfbs-hits", "tta-codons", "resistance-features", "onRegionChanged", "onAnnotationClicked", "onError"])
  ]);
}
const gf = /* @__PURE__ */ zt(df, [["render", pf], ["__scopeId", "data-v-e9236309"]]);
function Ro(e, t) {
  return function () {
    return e.apply(t, arguments);
  };
}
const { toString: mf } = Object.prototype, { getPrototypeOf: yr } = Object, { iterator: cs, toStringTag: Oo } = Symbol, us = /* @__PURE__ */ ((e) => (t) => {
  const n = mf.call(t);
  return e[n] || (e[n] = n.slice(8, -1).toLowerCase());
})(/* @__PURE__ */ Object.create(null)), Ze = (e) => (e = e.toLowerCase(), (t) => us(t) === e), fs = (e) => (t) => typeof t === e, { isArray: Kt } = Array, Vt = fs("undefined");
function Sn(e) {
  return e !== null && !Vt(e) && e.constructor !== null && !Vt(e.constructor) && $e(e.constructor.isBuffer) && e.constructor.isBuffer(e);
}
const Po = Ze("ArrayBuffer");
function bf(e) {
  let t;
  return typeof ArrayBuffer < "u" && ArrayBuffer.isView ? t = ArrayBuffer.isView(e) : t = e && e.buffer && Po(e.buffer), t;
}
const yf = fs("string"), $e = fs("function"), Lo = fs("number"), Cn = (e) => e !== null && typeof e == "object", _f = (e) => e === !0 || e === !1, Nn = (e) => {
  if (us(e) !== "object")
    return !1;
  const t = yr(e);
  return (t === null || t === Object.prototype || Object.getPrototypeOf(t) === null) && !(Oo in e) && !(cs in e);
}, wf = (e) => {
  if (!Cn(e) || Sn(e))
    return !1;
  try {
    return Object.keys(e).length === 0 && Object.getPrototypeOf(e) === Object.prototype;
  } catch {
    return !1;
  }
}, vf = Ze("Date"), xf = Ze("File"), Sf = Ze("Blob"), Cf = Ze("FileList"), Ef = (e) => Cn(e) && $e(e.pipe), Tf = (e) => {
  let t;
  return e && (typeof FormData == "function" && e instanceof FormData || $e(e.append) && ((t = us(e)) === "formdata" || // detect form-data instance
    t === "object" && $e(e.toString) && e.toString() === "[object FormData]"));
}, Af = Ze("URLSearchParams"), [kf, Rf, Of, Pf] = ["ReadableStream", "Request", "Response", "Headers"].map(Ze), Lf = (e) => e.trim ? e.trim() : e.replace(/^[\s\uFEFF\xA0]+|[\s\uFEFF\xA0]+$/g, "");
function En(e, t, { allOwnKeys: n = !1 } = {}) {
  if (e === null || typeof e > "u")
    return;
  let s, r;
  if (typeof e != "object" && (e = [e]), Kt(e))
    for (s = 0, r = e.length; s < r; s++)
      t.call(null, e[s], s, e);
  else {
    if (Sn(e))
      return;
    const i = n ? Object.getOwnPropertyNames(e) : Object.keys(e), o = i.length;
    let l;
    for (s = 0; s < o; s++)
      l = i[s], t.call(null, e[l], l, e);
  }
}
function Fo(e, t) {
  if (Sn(e))
    return null;
  t = t.toLowerCase();
  const n = Object.keys(e);
  let s = n.length, r;
  for (; s-- > 0;)
    if (r = n[s], t === r.toLowerCase())
      return r;
  return null;
}
const kt = typeof globalThis < "u" ? globalThis : typeof self < "u" ? self : typeof window < "u" ? window : global, Mo = (e) => !Vt(e) && e !== kt;
function Ws() {
  const { caseless: e, skipUndefined: t } = Mo(this) && this || {}, n = {}, s = (r, i) => {
    const o = e && Fo(n, i) || i;
    Nn(n[o]) && Nn(r) ? n[o] = Ws(n[o], r) : Nn(r) ? n[o] = Ws({}, r) : Kt(r) ? n[o] = r.slice() : (!t || !Vt(r)) && (n[o] = r);
  };
  for (let r = 0, i = arguments.length; r < i; r++)
    arguments[r] && En(arguments[r], s);
  return n;
}
const Ff = (e, t, n, { allOwnKeys: s } = {}) => (En(t, (r, i) => {
  n && $e(r) ? e[i] = Ro(r, n) : e[i] = r;
}, { allOwnKeys: s }), e), Mf = (e) => (e.charCodeAt(0) === 65279 && (e = e.slice(1)), e), If = (e, t, n, s) => {
  e.prototype = Object.create(t.prototype, s), e.prototype.constructor = e, Object.defineProperty(e, "super", {
    value: t.prototype
  }), n && Object.assign(e.prototype, n);
}, Nf = (e, t, n, s) => {
  let r, i, o;
  const l = {};
  if (t = t || {}, e == null) return t;
  do {
    for (r = Object.getOwnPropertyNames(e), i = r.length; i-- > 0;)
      o = r[i], (!s || s(o, e, t)) && !l[o] && (t[o] = e[o], l[o] = !0);
    e = n !== !1 && yr(e);
  } while (e && (!n || n(e, t)) && e !== Object.prototype);
  return t;
}, Df = (e, t, n) => {
  e = String(e), (n === void 0 || n > e.length) && (n = e.length), n -= t.length;
  const s = e.indexOf(t, n);
  return s !== -1 && s === n;
}, $f = (e) => {
  if (!e) return null;
  if (Kt(e)) return e;
  let t = e.length;
  if (!Lo(t)) return null;
  const n = new Array(t);
  for (; t-- > 0;)
    n[t] = e[t];
  return n;
}, jf = /* @__PURE__ */ ((e) => (t) => e && t instanceof e)(typeof Uint8Array < "u" && yr(Uint8Array)), Bf = (e, t) => {
  const s = (e && e[cs]).call(e);
  let r;
  for (; (r = s.next()) && !r.done;) {
    const i = r.value;
    t.call(e, i[0], i[1]);
  }
}, Hf = (e, t) => {
  let n;
  const s = [];
  for (; (n = e.exec(t)) !== null;)
    s.push(n);
  return s;
}, Uf = Ze("HTMLFormElement"), qf = (e) => e.toLowerCase().replace(
  /[-_\s]([a-z\d])(\w*)/g,
  function (n, s, r) {
    return s.toUpperCase() + r;
  }
), ni = (({ hasOwnProperty: e }) => (t, n) => e.call(t, n))(Object.prototype), Vf = Ze("RegExp"), Io = (e, t) => {
  const n = Object.getOwnPropertyDescriptors(e), s = {};
  En(n, (r, i) => {
    let o;
    (o = t(r, i, e)) !== !1 && (s[i] = o || r);
  }), Object.defineProperties(e, s);
}, Wf = (e) => {
  Io(e, (t, n) => {
    if ($e(e) && ["arguments", "caller", "callee"].indexOf(n) !== -1)
      return !1;
    const s = e[n];
    if ($e(s)) {
      if (t.enumerable = !1, "writable" in t) {
        t.writable = !1;
        return;
      }
      t.set || (t.set = () => {
        throw Error("Can not rewrite read-only method '" + n + "'");
      });
    }
  });
}, zf = (e, t) => {
  const n = {}, s = (r) => {
    r.forEach((i) => {
      n[i] = !0;
    });
  };
  return Kt(e) ? s(e) : s(String(e).split(t)), n;
}, Kf = () => {
}, Gf = (e, t) => e != null && Number.isFinite(e = +e) ? e : t;
function Jf(e) {
  return !!(e && $e(e.append) && e[Oo] === "FormData" && e[cs]);
}
const Yf = (e) => {
  const t = new Array(10), n = (s, r) => {
    if (Cn(s)) {
      if (t.indexOf(s) >= 0)
        return;
      if (Sn(s))
        return s;
      if (!("toJSON" in s)) {
        t[r] = s;
        const i = Kt(s) ? [] : {};
        return En(s, (o, l) => {
          const a = n(o, r + 1);
          !Vt(a) && (i[l] = a);
        }), t[r] = void 0, i;
      }
    }
    return s;
  };
  return n(e, 0);
}, Xf = Ze("AsyncFunction"), Qf = (e) => e && (Cn(e) || $e(e)) && $e(e.then) && $e(e.catch), No = ((e, t) => e ? setImmediate : t ? ((n, s) => (kt.addEventListener("message", ({ source: r, data: i }) => {
  r === kt && i === n && s.length && s.shift()();
}, !1), (r) => {
  s.push(r), kt.postMessage(n, "*");
}))(`axios@${Math.random()}`, []) : (n) => setTimeout(n))(
  typeof setImmediate == "function",
  $e(kt.postMessage)
), Zf = typeof queueMicrotask < "u" ? queueMicrotask.bind(kt) : typeof process < "u" && process.nextTick || No, ed = (e) => e != null && $e(e[cs]), y = {
  isArray: Kt,
  isArrayBuffer: Po,
  isBuffer: Sn,
  isFormData: Tf,
  isArrayBufferView: bf,
  isString: yf,
  isNumber: Lo,
  isBoolean: _f,
  isObject: Cn,
  isPlainObject: Nn,
  isEmptyObject: wf,
  isReadableStream: kf,
  isRequest: Rf,
  isResponse: Of,
  isHeaders: Pf,
  isUndefined: Vt,
  isDate: vf,
  isFile: xf,
  isBlob: Sf,
  isRegExp: Vf,
  isFunction: $e,
  isStream: Ef,
  isURLSearchParams: Af,
  isTypedArray: jf,
  isFileList: Cf,
  forEach: En,
  merge: Ws,
  extend: Ff,
  trim: Lf,
  stripBOM: Mf,
  inherits: If,
  toFlatObject: Nf,
  kindOf: us,
  kindOfTest: Ze,
  endsWith: Df,
  toArray: $f,
  forEachEntry: Bf,
  matchAll: Hf,
  isHTMLForm: Uf,
  hasOwnProperty: ni,
  hasOwnProp: ni,
  // an alias to avoid ESLint no-prototype-builtins detection
  reduceDescriptors: Io,
  freezeMethods: Wf,
  toObjectSet: zf,
  toCamelCase: qf,
  noop: Kf,
  toFiniteNumber: Gf,
  findKey: Fo,
  global: kt,
  isContextDefined: Mo,
  isSpecCompliantForm: Jf,
  toJSONObject: Yf,
  isAsyncFn: Xf,
  isThenable: Qf,
  setImmediate: No,
  asap: Zf,
  isIterable: ed
};
function ee(e, t, n, s, r) {
  Error.call(this), Error.captureStackTrace ? Error.captureStackTrace(this, this.constructor) : this.stack = new Error().stack, this.message = e, this.name = "AxiosError", t && (this.code = t), n && (this.config = n), s && (this.request = s), r && (this.response = r, this.status = r.status ? r.status : null);
}
y.inherits(ee, Error, {
  toJSON: function () {
    return {
      // Standard
      message: this.message,
      name: this.name,
      // Microsoft
      description: this.description,
      number: this.number,
      // Mozilla
      fileName: this.fileName,
      lineNumber: this.lineNumber,
      columnNumber: this.columnNumber,
      stack: this.stack,
      // Axios
      config: y.toJSONObject(this.config),
      code: this.code,
      status: this.status
    };
  }
});
const Do = ee.prototype, $o = {};
[
  "ERR_BAD_OPTION_VALUE",
  "ERR_BAD_OPTION",
  "ECONNABORTED",
  "ETIMEDOUT",
  "ERR_NETWORK",
  "ERR_FR_TOO_MANY_REDIRECTS",
  "ERR_DEPRECATED",
  "ERR_BAD_RESPONSE",
  "ERR_BAD_REQUEST",
  "ERR_CANCELED",
  "ERR_NOT_SUPPORT",
  "ERR_INVALID_URL"
  // eslint-disable-next-line func-names
].forEach((e) => {
  $o[e] = { value: e };
});
Object.defineProperties(ee, $o);
Object.defineProperty(Do, "isAxiosError", { value: !0 });
ee.from = (e, t, n, s, r, i) => {
  const o = Object.create(Do);
  y.toFlatObject(e, o, function (c) {
    return c !== Error.prototype;
  }, (u) => u !== "isAxiosError");
  const l = e && e.message ? e.message : "Error", a = t == null && e ? e.code : t;
  return ee.call(o, l, a, n, s, r), e && o.cause == null && Object.defineProperty(o, "cause", { value: e, configurable: !0 }), o.name = e && e.name || "Error", i && Object.assign(o, i), o;
};
const td = null;
function zs(e) {
  return y.isPlainObject(e) || y.isArray(e);
}
function jo(e) {
  return y.endsWith(e, "[]") ? e.slice(0, -2) : e;
}
function si(e, t, n) {
  return e ? e.concat(t).map(function (r, i) {
    return r = jo(r), !n && i ? "[" + r + "]" : r;
  }).join(n ? "." : "") : t;
}
function nd(e) {
  return y.isArray(e) && !e.some(zs);
}
const sd = y.toFlatObject(y, {}, null, function (t) {
  return /^is[A-Z]/.test(t);
});
function ds(e, t, n) {
  if (!y.isObject(e))
    throw new TypeError("target must be an object");
  t = t || new FormData(), n = y.toFlatObject(n, {
    metaTokens: !0,
    dots: !1,
    indexes: !1
  }, !1, function (x, A) {
    return !y.isUndefined(A[x]);
  });
  const s = n.metaTokens, r = n.visitor || c, i = n.dots, o = n.indexes, a = (n.Blob || typeof Blob < "u" && Blob) && y.isSpecCompliantForm(t);
  if (!y.isFunction(r))
    throw new TypeError("visitor must be a function");
  function u(m) {
    if (m === null) return "";
    if (y.isDate(m))
      return m.toISOString();
    if (y.isBoolean(m))
      return m.toString();
    if (!a && y.isBlob(m))
      throw new ee("Blob is not supported. Use a Buffer instead.");
    return y.isArrayBuffer(m) || y.isTypedArray(m) ? a && typeof Blob == "function" ? new Blob([m]) : Buffer.from(m) : m;
  }
  function c(m, x, A) {
    let H = m;
    if (m && !A && typeof m == "object") {
      if (y.endsWith(x, "{}"))
        x = s ? x : x.slice(0, -2), m = JSON.stringify(m);
      else if (y.isArray(m) && nd(m) || (y.isFileList(m) || y.endsWith(x, "[]")) && (H = y.toArray(m)))
        return x = jo(x), H.forEach(function (W, I) {
          !(y.isUndefined(W) || W === null) && t.append(
            // eslint-disable-next-line no-nested-ternary
            o === !0 ? si([x], I, i) : o === null ? x : x + "[]",
            u(W)
          );
        }), !1;
    }
    return zs(m) ? !0 : (t.append(si(A, x, i), u(m)), !1);
  }
  const f = [], p = Object.assign(sd, {
    defaultVisitor: c,
    convertValue: u,
    isVisitable: zs
  });
  function g(m, x) {
    if (!y.isUndefined(m)) {
      if (f.indexOf(m) !== -1)
        throw Error("Circular reference detected in " + x.join("."));
      f.push(m), y.forEach(m, function (H, q) {
        (!(y.isUndefined(H) || H === null) && r.call(
          t,
          H,
          y.isString(q) ? q.trim() : q,
          x,
          p
        )) === !0 && g(H, x ? x.concat(q) : [q]);
      }), f.pop();
    }
  }
  if (!y.isObject(e))
    throw new TypeError("data must be an object");
  return g(e), t;
}
function ri(e) {
  const t = {
    "!": "%21",
    "'": "%27",
    "(": "%28",
    ")": "%29",
    "~": "%7E",
    "%20": "+",
    "%00": "\0"
  };
  return encodeURIComponent(e).replace(/[!'()~]|%20|%00/g, function (s) {
    return t[s];
  });
}
function _r(e, t) {
  this._pairs = [], e && ds(e, this, t);
}
const Bo = _r.prototype;
Bo.append = function (t, n) {
  this._pairs.push([t, n]);
};
Bo.toString = function (t) {
  const n = t ? function (s) {
    return t.call(this, s, ri);
  } : ri;
  return this._pairs.map(function (r) {
    return n(r[0]) + "=" + n(r[1]);
  }, "").join("&");
};
function rd(e) {
  return encodeURIComponent(e).replace(/%3A/gi, ":").replace(/%24/g, "$").replace(/%2C/gi, ",").replace(/%20/g, "+");
}
function Ho(e, t, n) {
  if (!t)
    return e;
  const s = n && n.encode || rd;
  y.isFunction(n) && (n = {
    serialize: n
  });
  const r = n && n.serialize;
  let i;
  if (r ? i = r(t, n) : i = y.isURLSearchParams(t) ? t.toString() : new _r(t, n).toString(s), i) {
    const o = e.indexOf("#");
    o !== -1 && (e = e.slice(0, o)), e += (e.indexOf("?") === -1 ? "?" : "&") + i;
  }
  return e;
}
class ii {
  constructor() {
    this.handlers = [];
  }
  /**
   * Add a new interceptor to the stack
   *
   * @param {Function} fulfilled The function to handle `then` for a `Promise`
   * @param {Function} rejected The function to handle `reject` for a `Promise`
   *
   * @return {Number} An ID used to remove interceptor later
   */
  use(t, n, s) {
    return this.handlers.push({
      fulfilled: t,
      rejected: n,
      synchronous: s ? s.synchronous : !1,
      runWhen: s ? s.runWhen : null
    }), this.handlers.length - 1;
  }
  /**
   * Remove an interceptor from the stack
   *
   * @param {Number} id The ID that was returned by `use`
   *
   * @returns {void}
   */
  eject(t) {
    this.handlers[t] && (this.handlers[t] = null);
  }
  /**
   * Clear all interceptors from the stack
   *
   * @returns {void}
   */
  clear() {
    this.handlers && (this.handlers = []);
  }
  /**
   * Iterate over all the registered interceptors
   *
   * This method is particularly useful for skipping over any
   * interceptors that may have become `null` calling `eject`.
   *
   * @param {Function} fn The function to call for each interceptor
   *
   * @returns {void}
   */
  forEach(t) {
    y.forEach(this.handlers, function (s) {
      s !== null && t(s);
    });
  }
}
const Uo = {
  silentJSONParsing: !0,
  forcedJSONParsing: !0,
  clarifyTimeoutError: !1
}, id = typeof URLSearchParams < "u" ? URLSearchParams : _r, od = typeof FormData < "u" ? FormData : null, ld = typeof Blob < "u" ? Blob : null, ad = {
  isBrowser: !0,
  classes: {
    URLSearchParams: id,
    FormData: od,
    Blob: ld
  },
  protocols: ["http", "https", "file", "blob", "url", "data"]
}, wr = typeof window < "u" && typeof document < "u", Ks = typeof navigator == "object" && navigator || void 0, cd = wr && (!Ks || ["ReactNative", "NativeScript", "NS"].indexOf(Ks.product) < 0), ud = typeof WorkerGlobalScope < "u" && // eslint-disable-next-line no-undef
  self instanceof WorkerGlobalScope && typeof self.importScripts == "function", fd = wr && window.location.href || "http://localhost", dd = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
    __proto__: null,
    hasBrowserEnv: wr,
    hasStandardBrowserEnv: cd,
    hasStandardBrowserWebWorkerEnv: ud,
    navigator: Ks,
    origin: fd
  }, Symbol.toStringTag, { value: "Module" })), Pe = {
    ...dd,
    ...ad
  };
function hd(e, t) {
  return ds(e, new Pe.classes.URLSearchParams(), {
    visitor: function (n, s, r, i) {
      return Pe.isNode && y.isBuffer(n) ? (this.append(s, n.toString("base64")), !1) : i.defaultVisitor.apply(this, arguments);
    },
    ...t
  });
}
function pd(e) {
  return y.matchAll(/\w+|\[(\w*)]/g, e).map((t) => t[0] === "[]" ? "" : t[1] || t[0]);
}
function gd(e) {
  const t = {}, n = Object.keys(e);
  let s;
  const r = n.length;
  let i;
  for (s = 0; s < r; s++)
    i = n[s], t[i] = e[i];
  return t;
}
function qo(e) {
  function t(n, s, r, i) {
    let o = n[i++];
    if (o === "__proto__") return !0;
    const l = Number.isFinite(+o), a = i >= n.length;
    return o = !o && y.isArray(r) ? r.length : o, a ? (y.hasOwnProp(r, o) ? r[o] = [r[o], s] : r[o] = s, !l) : ((!r[o] || !y.isObject(r[o])) && (r[o] = []), t(n, s, r[o], i) && y.isArray(r[o]) && (r[o] = gd(r[o])), !l);
  }
  if (y.isFormData(e) && y.isFunction(e.entries)) {
    const n = {};
    return y.forEachEntry(e, (s, r) => {
      t(pd(s), r, n, 0);
    }), n;
  }
  return null;
}
function md(e, t, n) {
  if (y.isString(e))
    try {
      return (t || JSON.parse)(e), y.trim(e);
    } catch (s) {
      if (s.name !== "SyntaxError")
        throw s;
    }
  return (n || JSON.stringify)(e);
}
const Tn = {
  transitional: Uo,
  adapter: ["xhr", "http", "fetch"],
  transformRequest: [function (t, n) {
    const s = n.getContentType() || "", r = s.indexOf("application/json") > -1, i = y.isObject(t);
    if (i && y.isHTMLForm(t) && (t = new FormData(t)), y.isFormData(t))
      return r ? JSON.stringify(qo(t)) : t;
    if (y.isArrayBuffer(t) || y.isBuffer(t) || y.isStream(t) || y.isFile(t) || y.isBlob(t) || y.isReadableStream(t))
      return t;
    if (y.isArrayBufferView(t))
      return t.buffer;
    if (y.isURLSearchParams(t))
      return n.setContentType("application/x-www-form-urlencoded;charset=utf-8", !1), t.toString();
    let l;
    if (i) {
      if (s.indexOf("application/x-www-form-urlencoded") > -1)
        return hd(t, this.formSerializer).toString();
      if ((l = y.isFileList(t)) || s.indexOf("multipart/form-data") > -1) {
        const a = this.env && this.env.FormData;
        return ds(
          l ? { "files[]": t } : t,
          a && new a(),
          this.formSerializer
        );
      }
    }
    return i || r ? (n.setContentType("application/json", !1), md(t)) : t;
  }],
  transformResponse: [function (t) {
    const n = this.transitional || Tn.transitional, s = n && n.forcedJSONParsing, r = this.responseType === "json";
    if (y.isResponse(t) || y.isReadableStream(t))
      return t;
    if (t && y.isString(t) && (s && !this.responseType || r)) {
      const o = !(n && n.silentJSONParsing) && r;
      try {
        return JSON.parse(t, this.parseReviver);
      } catch (l) {
        if (o)
          throw l.name === "SyntaxError" ? ee.from(l, ee.ERR_BAD_RESPONSE, this, null, this.response) : l;
      }
    }
    return t;
  }],
  /**
   * A timeout in milliseconds to abort a request. If set to 0 (default) a
   * timeout is not created.
   */
  timeout: 0,
  xsrfCookieName: "XSRF-TOKEN",
  xsrfHeaderName: "X-XSRF-TOKEN",
  maxContentLength: -1,
  maxBodyLength: -1,
  env: {
    FormData: Pe.classes.FormData,
    Blob: Pe.classes.Blob
  },
  validateStatus: function (t) {
    return t >= 200 && t < 300;
  },
  headers: {
    common: {
      Accept: "application/json, text/plain, */*",
      "Content-Type": void 0
    }
  }
};
y.forEach(["delete", "get", "head", "post", "put", "patch"], (e) => {
  Tn.headers[e] = {};
});
const bd = y.toObjectSet([
  "age",
  "authorization",
  "content-length",
  "content-type",
  "etag",
  "expires",
  "from",
  "host",
  "if-modified-since",
  "if-unmodified-since",
  "last-modified",
  "location",
  "max-forwards",
  "proxy-authorization",
  "referer",
  "retry-after",
  "user-agent"
]), yd = (e) => {
  const t = {};
  let n, s, r;
  return e && e.split(`
`).forEach(function (o) {
    r = o.indexOf(":"), n = o.substring(0, r).trim().toLowerCase(), s = o.substring(r + 1).trim(), !(!n || t[n] && bd[n]) && (n === "set-cookie" ? t[n] ? t[n].push(s) : t[n] = [s] : t[n] = t[n] ? t[n] + ", " + s : s);
  }), t;
}, oi = Symbol("internals");
function nn(e) {
  return e && String(e).trim().toLowerCase();
}
function Dn(e) {
  return e === !1 || e == null ? e : y.isArray(e) ? e.map(Dn) : String(e);
}
function _d(e) {
  const t = /* @__PURE__ */ Object.create(null), n = /([^\s,;=]+)\s*(?:=\s*([^,;]+))?/g;
  let s;
  for (; s = n.exec(e);)
    t[s[1]] = s[2];
  return t;
}
const wd = (e) => /^[-_a-zA-Z0-9^`|~,!#$%&'*+.]+$/.test(e.trim());
function Os(e, t, n, s, r) {
  if (y.isFunction(s))
    return s.call(this, t, n);
  if (r && (t = n), !!y.isString(t)) {
    if (y.isString(s))
      return t.indexOf(s) !== -1;
    if (y.isRegExp(s))
      return s.test(t);
  }
}
function vd(e) {
  return e.trim().toLowerCase().replace(/([a-z\d])(\w*)/g, (t, n, s) => n.toUpperCase() + s);
}
function xd(e, t) {
  const n = y.toCamelCase(" " + t);
  ["get", "set", "has"].forEach((s) => {
    Object.defineProperty(e, s + n, {
      value: function (r, i, o) {
        return this[s].call(this, t, r, i, o);
      },
      configurable: !0
    });
  });
}
let je = class {
  constructor(t) {
    t && this.set(t);
  }
  set(t, n, s) {
    const r = this;
    function i(l, a, u) {
      const c = nn(a);
      if (!c)
        throw new Error("header name must be a non-empty string");
      const f = y.findKey(r, c);
      (!f || r[f] === void 0 || u === !0 || u === void 0 && r[f] !== !1) && (r[f || a] = Dn(l));
    }
    const o = (l, a) => y.forEach(l, (u, c) => i(u, c, a));
    if (y.isPlainObject(t) || t instanceof this.constructor)
      o(t, n);
    else if (y.isString(t) && (t = t.trim()) && !wd(t))
      o(yd(t), n);
    else if (y.isObject(t) && y.isIterable(t)) {
      let l = {}, a, u;
      for (const c of t) {
        if (!y.isArray(c))
          throw TypeError("Object iterator must return a key-value pair");
        l[u = c[0]] = (a = l[u]) ? y.isArray(a) ? [...a, c[1]] : [a, c[1]] : c[1];
      }
      o(l, n);
    } else
      t != null && i(n, t, s);
    return this;
  }
  get(t, n) {
    if (t = nn(t), t) {
      const s = y.findKey(this, t);
      if (s) {
        const r = this[s];
        if (!n)
          return r;
        if (n === !0)
          return _d(r);
        if (y.isFunction(n))
          return n.call(this, r, s);
        if (y.isRegExp(n))
          return n.exec(r);
        throw new TypeError("parser must be boolean|regexp|function");
      }
    }
  }
  has(t, n) {
    if (t = nn(t), t) {
      const s = y.findKey(this, t);
      return !!(s && this[s] !== void 0 && (!n || Os(this, this[s], s, n)));
    }
    return !1;
  }
  delete(t, n) {
    const s = this;
    let r = !1;
    function i(o) {
      if (o = nn(o), o) {
        const l = y.findKey(s, o);
        l && (!n || Os(s, s[l], l, n)) && (delete s[l], r = !0);
      }
    }
    return y.isArray(t) ? t.forEach(i) : i(t), r;
  }
  clear(t) {
    const n = Object.keys(this);
    let s = n.length, r = !1;
    for (; s--;) {
      const i = n[s];
      (!t || Os(this, this[i], i, t, !0)) && (delete this[i], r = !0);
    }
    return r;
  }
  normalize(t) {
    const n = this, s = {};
    return y.forEach(this, (r, i) => {
      const o = y.findKey(s, i);
      if (o) {
        n[o] = Dn(r), delete n[i];
        return;
      }
      const l = t ? vd(i) : String(i).trim();
      l !== i && delete n[i], n[l] = Dn(r), s[l] = !0;
    }), this;
  }
  concat(...t) {
    return this.constructor.concat(this, ...t);
  }
  toJSON(t) {
    const n = /* @__PURE__ */ Object.create(null);
    return y.forEach(this, (s, r) => {
      s != null && s !== !1 && (n[r] = t && y.isArray(s) ? s.join(", ") : s);
    }), n;
  }
  [Symbol.iterator]() {
    return Object.entries(this.toJSON())[Symbol.iterator]();
  }
  toString() {
    return Object.entries(this.toJSON()).map(([t, n]) => t + ": " + n).join(`
`);
  }
  getSetCookie() {
    return this.get("set-cookie") || [];
  }
  get [Symbol.toStringTag]() {
    return "AxiosHeaders";
  }
  static from(t) {
    return t instanceof this ? t : new this(t);
  }
  static concat(t, ...n) {
    const s = new this(t);
    return n.forEach((r) => s.set(r)), s;
  }
  static accessor(t) {
    const s = (this[oi] = this[oi] = {
      accessors: {}
    }).accessors, r = this.prototype;
    function i(o) {
      const l = nn(o);
      s[l] || (xd(r, o), s[l] = !0);
    }
    return y.isArray(t) ? t.forEach(i) : i(t), this;
  }
};
je.accessor(["Content-Type", "Content-Length", "Accept", "Accept-Encoding", "User-Agent", "Authorization"]);
y.reduceDescriptors(je.prototype, ({ value: e }, t) => {
  let n = t[0].toUpperCase() + t.slice(1);
  return {
    get: () => e,
    set(s) {
      this[n] = s;
    }
  };
});
y.freezeMethods(je);
function Ps(e, t) {
  const n = this || Tn, s = t || n, r = je.from(s.headers);
  let i = s.data;
  return y.forEach(e, function (l) {
    i = l.call(n, i, r.normalize(), t ? t.status : void 0);
  }), r.normalize(), i;
}
function Vo(e) {
  return !!(e && e.__CANCEL__);
}
function Gt(e, t, n) {
  ee.call(this, e ?? "canceled", ee.ERR_CANCELED, t, n), this.name = "CanceledError";
}
y.inherits(Gt, ee, {
  __CANCEL__: !0
});
function Wo(e, t, n) {
  const s = n.config.validateStatus;
  !n.status || !s || s(n.status) ? e(n) : t(new ee(
    "Request failed with status code " + n.status,
    [ee.ERR_BAD_REQUEST, ee.ERR_BAD_RESPONSE][Math.floor(n.status / 100) - 4],
    n.config,
    n.request,
    n
  ));
}
function Sd(e) {
  const t = /^([-+\w]{1,25})(:?\/\/|:)/.exec(e);
  return t && t[1] || "";
}
function Cd(e, t) {
  e = e || 10;
  const n = new Array(e), s = new Array(e);
  let r = 0, i = 0, o;
  return t = t !== void 0 ? t : 1e3, function (a) {
    const u = Date.now(), c = s[i];
    o || (o = u), n[r] = a, s[r] = u;
    let f = i, p = 0;
    for (; f !== r;)
      p += n[f++], f = f % e;
    if (r = (r + 1) % e, r === i && (i = (i + 1) % e), u - o < t)
      return;
    const g = c && u - c;
    return g ? Math.round(p * 1e3 / g) : void 0;
  };
}
function Ed(e, t) {
  let n = 0, s = 1e3 / t, r, i;
  const o = (u, c = Date.now()) => {
    n = c, r = null, i && (clearTimeout(i), i = null), e(...u);
  };
  return [(...u) => {
    const c = Date.now(), f = c - n;
    f >= s ? o(u, c) : (r = u, i || (i = setTimeout(() => {
      i = null, o(r);
    }, s - f)));
  }, () => r && o(r)];
}
const Yn = (e, t, n = 3) => {
  let s = 0;
  const r = Cd(50, 250);
  return Ed((i) => {
    const o = i.loaded, l = i.lengthComputable ? i.total : void 0, a = o - s, u = r(a), c = o <= l;
    s = o;
    const f = {
      loaded: o,
      total: l,
      progress: l ? o / l : void 0,
      bytes: a,
      rate: u || void 0,
      estimated: u && l && c ? (l - o) / u : void 0,
      event: i,
      lengthComputable: l != null,
      [t ? "download" : "upload"]: !0
    };
    e(f);
  }, n);
}, li = (e, t) => {
  const n = e != null;
  return [(s) => t[0]({
    lengthComputable: n,
    total: e,
    loaded: s
  }), t[1]];
}, ai = (e) => (...t) => y.asap(() => e(...t)), Td = Pe.hasStandardBrowserEnv ? /* @__PURE__ */ ((e, t) => (n) => (n = new URL(n, Pe.origin), e.protocol === n.protocol && e.host === n.host && (t || e.port === n.port)))(
  new URL(Pe.origin),
  Pe.navigator && /(msie|trident)/i.test(Pe.navigator.userAgent)
) : () => !0, Ad = Pe.hasStandardBrowserEnv ? (
  // Standard browser envs support document.cookie
  {
    write(e, t, n, s, r, i, o) {
      if (typeof document > "u") return;
      const l = [`${e}=${encodeURIComponent(t)}`];
      y.isNumber(n) && l.push(`expires=${new Date(n).toUTCString()}`), y.isString(s) && l.push(`path=${s}`), y.isString(r) && l.push(`domain=${r}`), i === !0 && l.push("secure"), y.isString(o) && l.push(`SameSite=${o}`), document.cookie = l.join("; ");
    },
    read(e) {
      if (typeof document > "u") return null;
      const t = document.cookie.match(new RegExp("(?:^|; )" + e + "=([^;]*)"));
      return t ? decodeURIComponent(t[1]) : null;
    },
    remove(e) {
      this.write(e, "", Date.now() - 864e5, "/");
    }
  }
) : (
  // Non-standard browser env (web workers, react-native) lack needed support.
  {
    write() {
    },
    read() {
      return null;
    },
    remove() {
    }
  }
);
function kd(e) {
  return /^([a-z][a-z\d+\-.]*:)?\/\//i.test(e);
}
function Rd(e, t) {
  return t ? e.replace(/\/?\/$/, "") + "/" + t.replace(/^\/+/, "") : e;
}
function zo(e, t, n) {
  let s = !kd(t);
  return e && (s || n == !1) ? Rd(e, t) : t;
}
const ci = (e) => e instanceof je ? { ...e } : e;
function Ft(e, t) {
  t = t || {};
  const n = {};
  function s(u, c, f, p) {
    return y.isPlainObject(u) && y.isPlainObject(c) ? y.merge.call({ caseless: p }, u, c) : y.isPlainObject(c) ? y.merge({}, c) : y.isArray(c) ? c.slice() : c;
  }
  function r(u, c, f, p) {
    if (y.isUndefined(c)) {
      if (!y.isUndefined(u))
        return s(void 0, u, f, p);
    } else return s(u, c, f, p);
  }
  function i(u, c) {
    if (!y.isUndefined(c))
      return s(void 0, c);
  }
  function o(u, c) {
    if (y.isUndefined(c)) {
      if (!y.isUndefined(u))
        return s(void 0, u);
    } else return s(void 0, c);
  }
  function l(u, c, f) {
    if (f in t)
      return s(u, c);
    if (f in e)
      return s(void 0, u);
  }
  const a = {
    url: i,
    method: i,
    data: i,
    baseURL: o,
    transformRequest: o,
    transformResponse: o,
    paramsSerializer: o,
    timeout: o,
    timeoutMessage: o,
    withCredentials: o,
    withXSRFToken: o,
    adapter: o,
    responseType: o,
    xsrfCookieName: o,
    xsrfHeaderName: o,
    onUploadProgress: o,
    onDownloadProgress: o,
    decompress: o,
    maxContentLength: o,
    maxBodyLength: o,
    beforeRedirect: o,
    transport: o,
    httpAgent: o,
    httpsAgent: o,
    cancelToken: o,
    socketPath: o,
    responseEncoding: o,
    validateStatus: l,
    headers: (u, c, f) => r(ci(u), ci(c), f, !0)
  };
  return y.forEach(Object.keys({ ...e, ...t }), function (c) {
    const f = a[c] || r, p = f(e[c], t[c], c);
    y.isUndefined(p) && f !== l || (n[c] = p);
  }), n;
}
const Ko = (e) => {
  const t = Ft({}, e);
  let { data: n, withXSRFToken: s, xsrfHeaderName: r, xsrfCookieName: i, headers: o, auth: l } = t;
  if (t.headers = o = je.from(o), t.url = Ho(zo(t.baseURL, t.url, t.allowAbsoluteUrls), e.params, e.paramsSerializer), l && o.set(
    "Authorization",
    "Basic " + btoa((l.username || "") + ":" + (l.password ? unescape(encodeURIComponent(l.password)) : ""))
  ), y.isFormData(n)) {
    if (Pe.hasStandardBrowserEnv || Pe.hasStandardBrowserWebWorkerEnv)
      o.setContentType(void 0);
    else if (y.isFunction(n.getHeaders)) {
      const a = n.getHeaders(), u = ["content-type", "content-length"];
      Object.entries(a).forEach(([c, f]) => {
        u.includes(c.toLowerCase()) && o.set(c, f);
      });
    }
  }
  if (Pe.hasStandardBrowserEnv && (s && y.isFunction(s) && (s = s(t)), s || s !== !1 && Td(t.url))) {
    const a = r && i && Ad.read(i);
    a && o.set(r, a);
  }
  return t;
}, Od = typeof XMLHttpRequest < "u", Pd = Od && function (e) {
  return new Promise(function (n, s) {
    const r = Ko(e);
    let i = r.data;
    const o = je.from(r.headers).normalize();
    let { responseType: l, onUploadProgress: a, onDownloadProgress: u } = r, c, f, p, g, m;
    function x() {
      g && g(), m && m(), r.cancelToken && r.cancelToken.unsubscribe(c), r.signal && r.signal.removeEventListener("abort", c);
    }
    let A = new XMLHttpRequest();
    A.open(r.method.toUpperCase(), r.url, !0), A.timeout = r.timeout;
    function H() {
      if (!A)
        return;
      const W = je.from(
        "getAllResponseHeaders" in A && A.getAllResponseHeaders()
      ), J = {
        data: !l || l === "text" || l === "json" ? A.responseText : A.response,
        status: A.status,
        statusText: A.statusText,
        headers: W,
        config: e,
        request: A
      };
      Wo(function (ie) {
        n(ie), x();
      }, function (ie) {
        s(ie), x();
      }, J), A = null;
    }
    "onloadend" in A ? A.onloadend = H : A.onreadystatechange = function () {
      !A || A.readyState !== 4 || A.status === 0 && !(A.responseURL && A.responseURL.indexOf("file:") === 0) || setTimeout(H);
    }, A.onabort = function () {
      A && (s(new ee("Request aborted", ee.ECONNABORTED, e, A)), A = null);
    }, A.onerror = function (I) {
      const J = I && I.message ? I.message : "Network Error", de = new ee(J, ee.ERR_NETWORK, e, A);
      de.event = I || null, s(de), A = null;
    }, A.ontimeout = function () {
      let I = r.timeout ? "timeout of " + r.timeout + "ms exceeded" : "timeout exceeded";
      const J = r.transitional || Uo;
      r.timeoutErrorMessage && (I = r.timeoutErrorMessage), s(new ee(
        I,
        J.clarifyTimeoutError ? ee.ETIMEDOUT : ee.ECONNABORTED,
        e,
        A
      )), A = null;
    }, i === void 0 && o.setContentType(null), "setRequestHeader" in A && y.forEach(o.toJSON(), function (I, J) {
      A.setRequestHeader(J, I);
    }), y.isUndefined(r.withCredentials) || (A.withCredentials = !!r.withCredentials), l && l !== "json" && (A.responseType = r.responseType), u && ([p, m] = Yn(u, !0), A.addEventListener("progress", p)), a && A.upload && ([f, g] = Yn(a), A.upload.addEventListener("progress", f), A.upload.addEventListener("loadend", g)), (r.cancelToken || r.signal) && (c = (W) => {
      A && (s(!W || W.type ? new Gt(null, e, A) : W), A.abort(), A = null);
    }, r.cancelToken && r.cancelToken.subscribe(c), r.signal && (r.signal.aborted ? c() : r.signal.addEventListener("abort", c)));
    const q = Sd(r.url);
    if (q && Pe.protocols.indexOf(q) === -1) {
      s(new ee("Unsupported protocol " + q + ":", ee.ERR_BAD_REQUEST, e));
      return;
    }
    A.send(i || null);
  });
}, Ld = (e, t) => {
  const { length: n } = e = e ? e.filter(Boolean) : [];
  if (t || n) {
    let s = new AbortController(), r;
    const i = function (u) {
      if (!r) {
        r = !0, l();
        const c = u instanceof Error ? u : this.reason;
        s.abort(c instanceof ee ? c : new Gt(c instanceof Error ? c.message : c));
      }
    };
    let o = t && setTimeout(() => {
      o = null, i(new ee(`timeout ${t} of ms exceeded`, ee.ETIMEDOUT));
    }, t);
    const l = () => {
      e && (o && clearTimeout(o), o = null, e.forEach((u) => {
        u.unsubscribe ? u.unsubscribe(i) : u.removeEventListener("abort", i);
      }), e = null);
    };
    e.forEach((u) => u.addEventListener("abort", i));
    const { signal: a } = s;
    return a.unsubscribe = () => y.asap(l), a;
  }
}, Fd = function* (e, t) {
  let n = e.byteLength;
  if (n < t) {
    yield e;
    return;
  }
  let s = 0, r;
  for (; s < n;)
    r = s + t, yield e.slice(s, r), s = r;
}, Md = async function* (e, t) {
  for await (const n of Id(e))
    yield* Fd(n, t);
}, Id = async function* (e) {
  if (e[Symbol.asyncIterator]) {
    yield* e;
    return;
  }
  const t = e.getReader();
  try {
    for (; ;) {
      const { done: n, value: s } = await t.read();
      if (n)
        break;
      yield s;
    }
  } finally {
    await t.cancel();
  }
}, ui = (e, t, n, s) => {
  const r = Md(e, t);
  let i = 0, o, l = (a) => {
    o || (o = !0, s && s(a));
  };
  return new ReadableStream({
    async pull(a) {
      try {
        const { done: u, value: c } = await r.next();
        if (u) {
          l(), a.close();
          return;
        }
        let f = c.byteLength;
        if (n) {
          let p = i += f;
          n(p);
        }
        a.enqueue(new Uint8Array(c));
      } catch (u) {
        throw l(u), u;
      }
    },
    cancel(a) {
      return l(a), r.return();
    }
  }, {
    highWaterMark: 2
  });
}, fi = 64 * 1024, { isFunction: Pn } = y, Nd = (({ Request: e, Response: t }) => ({
  Request: e,
  Response: t
}))(y.global), {
  ReadableStream: di,
  TextEncoder: hi
} = y.global, pi = (e, ...t) => {
  try {
    return !!e(...t);
  } catch {
    return !1;
  }
}, Dd = (e) => {
  e = y.merge.call({
    skipUndefined: !0
  }, Nd, e);
  const { fetch: t, Request: n, Response: s } = e, r = t ? Pn(t) : typeof fetch == "function", i = Pn(n), o = Pn(s);
  if (!r)
    return !1;
  const l = r && Pn(di), a = r && (typeof hi == "function" ? /* @__PURE__ */ ((m) => (x) => m.encode(x))(new hi()) : async (m) => new Uint8Array(await new n(m).arrayBuffer())), u = i && l && pi(() => {
    let m = !1;
    const x = new n(Pe.origin, {
      body: new di(),
      method: "POST",
      get duplex() {
        return m = !0, "half";
      }
    }).headers.has("Content-Type");
    return m && !x;
  }), c = o && l && pi(() => y.isReadableStream(new s("").body)), f = {
    stream: c && ((m) => m.body)
  };
  r && ["text", "arrayBuffer", "blob", "formData", "stream"].forEach((m) => {
    !f[m] && (f[m] = (x, A) => {
      let H = x && x[m];
      if (H)
        return H.call(x);
      throw new ee(`Response type '${m}' is not supported`, ee.ERR_NOT_SUPPORT, A);
    });
  });
  const p = async (m) => {
    if (m == null)
      return 0;
    if (y.isBlob(m))
      return m.size;
    if (y.isSpecCompliantForm(m))
      return (await new n(Pe.origin, {
        method: "POST",
        body: m
      }).arrayBuffer()).byteLength;
    if (y.isArrayBufferView(m) || y.isArrayBuffer(m))
      return m.byteLength;
    if (y.isURLSearchParams(m) && (m = m + ""), y.isString(m))
      return (await a(m)).byteLength;
  }, g = async (m, x) => {
    const A = y.toFiniteNumber(m.getContentLength());
    return A ?? p(x);
  };
  return async (m) => {
    let {
      url: x,
      method: A,
      data: H,
      signal: q,
      cancelToken: W,
      timeout: I,
      onDownloadProgress: J,
      onUploadProgress: de,
      responseType: ie,
      headers: Se,
      withCredentials: v = "same-origin",
      fetchOptions: P
    } = Ko(m), D = t || fetch;
    ie = ie ? (ie + "").toLowerCase() : "text";
    let Q = Ld([q, W && W.toAbortSignal()], I), te = null;
    const oe = Q && Q.unsubscribe && (() => {
      Q.unsubscribe();
    });
    let le;
    try {
      if (de && u && A !== "get" && A !== "head" && (le = await g(Se, H)) !== 0) {
        let R = new n(x, {
          method: "POST",
          body: H,
          duplex: "half"
        }), O;
        if (y.isFormData(H) && (O = R.headers.get("content-type")) && Se.setContentType(O), R.body) {
          const [j, re] = li(
            le,
            Yn(ai(de))
          );
          H = ui(R.body, fi, j, re);
        }
      }
      y.isString(v) || (v = v ? "include" : "omit");
      const z = i && "credentials" in n.prototype, X = {
        ...P,
        signal: Q,
        method: A.toUpperCase(),
        headers: Se.normalize().toJSON(),
        body: H,
        duplex: "half",
        credentials: z ? v : void 0
      };
      te = i && new n(x, X);
      let G = await (i ? D(te, P) : D(x, X));
      const b = c && (ie === "stream" || ie === "response");
      if (c && (J || b && oe)) {
        const R = {};
        ["status", "statusText", "headers"].forEach((Ae) => {
          R[Ae] = G[Ae];
        });
        const O = y.toFiniteNumber(G.headers.get("content-length")), [j, re] = J && li(
          O,
          Yn(ai(J), !0)
        ) || [];
        G = new s(
          ui(G.body, fi, j, () => {
            re && re(), oe && oe();
          }),
          R
        );
      }
      ie = ie || "text";
      let _ = await f[y.findKey(f, ie) || "text"](G, m);
      return !b && oe && oe(), await new Promise((R, O) => {
        Wo(R, O, {
          data: _,
          headers: je.from(G.headers),
          status: G.status,
          statusText: G.statusText,
          config: m,
          request: te
        });
      });
    } catch (z) {
      throw oe && oe(), z && z.name === "TypeError" && /Load failed|fetch/i.test(z.message) ? Object.assign(
        new ee("Network Error", ee.ERR_NETWORK, m, te),
        {
          cause: z.cause || z
        }
      ) : ee.from(z, z && z.code, m, te);
    }
  };
}, $d = /* @__PURE__ */ new Map(), Go = (e) => {
  let t = e && e.env || {};
  const { fetch: n, Request: s, Response: r } = t, i = [
    s,
    r,
    n
  ];
  let o = i.length, l = o, a, u, c = $d;
  for (; l--;)
    a = i[l], u = c.get(a), u === void 0 && c.set(a, u = l ? /* @__PURE__ */ new Map() : Dd(t)), c = u;
  return u;
};
Go();
const vr = {
  http: td,
  xhr: Pd,
  fetch: {
    get: Go
  }
};
y.forEach(vr, (e, t) => {
  if (e) {
    try {
      Object.defineProperty(e, "name", { value: t });
    } catch {
    }
    Object.defineProperty(e, "adapterName", { value: t });
  }
});
const gi = (e) => `- ${e}`, jd = (e) => y.isFunction(e) || e === null || e === !1;
function Bd(e, t) {
  e = y.isArray(e) ? e : [e];
  const { length: n } = e;
  let s, r;
  const i = {};
  for (let o = 0; o < n; o++) {
    s = e[o];
    let l;
    if (r = s, !jd(s) && (r = vr[(l = String(s)).toLowerCase()], r === void 0))
      throw new ee(`Unknown adapter '${l}'`);
    if (r && (y.isFunction(r) || (r = r.get(t))))
      break;
    i[l || "#" + o] = r;
  }
  if (!r) {
    const o = Object.entries(i).map(
      ([a, u]) => `adapter ${a} ` + (u === !1 ? "is not supported by the environment" : "is not available in the build")
    );
    let l = n ? o.length > 1 ? `since :
` + o.map(gi).join(`
`) : " " + gi(o[0]) : "as no adapter specified";
    throw new ee(
      "There is no suitable adapter to dispatch the request " + l,
      "ERR_NOT_SUPPORT"
    );
  }
  return r;
}
const Jo = {
  /**
   * Resolve an adapter from a list of adapter names or functions.
   * @type {Function}
   */
  getAdapter: Bd,
  /**
   * Exposes all known adapters
   * @type {Object<string, Function|Object>}
   */
  adapters: vr
};
function Ls(e) {
  if (e.cancelToken && e.cancelToken.throwIfRequested(), e.signal && e.signal.aborted)
    throw new Gt(null, e);
}
function mi(e) {
  return Ls(e), e.headers = je.from(e.headers), e.data = Ps.call(
    e,
    e.transformRequest
  ), ["post", "put", "patch"].indexOf(e.method) !== -1 && e.headers.setContentType("application/x-www-form-urlencoded", !1), Jo.getAdapter(e.adapter || Tn.adapter, e)(e).then(function (s) {
    return Ls(e), s.data = Ps.call(
      e,
      e.transformResponse,
      s
    ), s.headers = je.from(s.headers), s;
  }, function (s) {
    return Vo(s) || (Ls(e), s && s.response && (s.response.data = Ps.call(
      e,
      e.transformResponse,
      s.response
    ), s.response.headers = je.from(s.response.headers))), Promise.reject(s);
  });
}
const Yo = "1.13.2", hs = {};
["object", "boolean", "number", "function", "string", "symbol"].forEach((e, t) => {
  hs[e] = function (s) {
    return typeof s === e || "a" + (t < 1 ? "n " : " ") + e;
  };
});
const bi = {};
hs.transitional = function (t, n, s) {
  function r(i, o) {
    return "[Axios v" + Yo + "] Transitional option '" + i + "'" + o + (s ? ". " + s : "");
  }
  return (i, o, l) => {
    if (t === !1)
      throw new ee(
        r(o, " has been removed" + (n ? " in " + n : "")),
        ee.ERR_DEPRECATED
      );
    return n && !bi[o] && (bi[o] = !0, console.warn(
      r(
        o,
        " has been deprecated since v" + n + " and will be removed in the near future"
      )
    )), t ? t(i, o, l) : !0;
  };
};
hs.spelling = function (t) {
  return (n, s) => (console.warn(`${s} is likely a misspelling of ${t}`), !0);
};
function Hd(e, t, n) {
  if (typeof e != "object")
    throw new ee("options must be an object", ee.ERR_BAD_OPTION_VALUE);
  const s = Object.keys(e);
  let r = s.length;
  for (; r-- > 0;) {
    const i = s[r], o = t[i];
    if (o) {
      const l = e[i], a = l === void 0 || o(l, i, e);
      if (a !== !0)
        throw new ee("option " + i + " must be " + a, ee.ERR_BAD_OPTION_VALUE);
      continue;
    }
    if (n !== !0)
      throw new ee("Unknown option " + i, ee.ERR_BAD_OPTION);
  }
}
const $n = {
  assertOptions: Hd,
  validators: hs
}, st = $n.validators;
let Pt = class {
  constructor(t) {
    this.defaults = t || {}, this.interceptors = {
      request: new ii(),
      response: new ii()
    };
  }
  /**
   * Dispatch a request
   *
   * @param {String|Object} configOrUrl The config specific for this request (merged with this.defaults)
   * @param {?Object} config
   *
   * @returns {Promise} The Promise to be fulfilled
   */
  async request(t, n) {
    try {
      return await this._request(t, n);
    } catch (s) {
      if (s instanceof Error) {
        let r = {};
        Error.captureStackTrace ? Error.captureStackTrace(r) : r = new Error();
        const i = r.stack ? r.stack.replace(/^.+\n/, "") : "";
        try {
          s.stack ? i && !String(s.stack).endsWith(i.replace(/^.+\n.+\n/, "")) && (s.stack += `
` + i) : s.stack = i;
        } catch {
        }
      }
      throw s;
    }
  }
  _request(t, n) {
    typeof t == "string" ? (n = n || {}, n.url = t) : n = t || {}, n = Ft(this.defaults, n);
    const { transitional: s, paramsSerializer: r, headers: i } = n;
    s !== void 0 && $n.assertOptions(s, {
      silentJSONParsing: st.transitional(st.boolean),
      forcedJSONParsing: st.transitional(st.boolean),
      clarifyTimeoutError: st.transitional(st.boolean)
    }, !1), r != null && (y.isFunction(r) ? n.paramsSerializer = {
      serialize: r
    } : $n.assertOptions(r, {
      encode: st.function,
      serialize: st.function
    }, !0)), n.allowAbsoluteUrls !== void 0 || (this.defaults.allowAbsoluteUrls !== void 0 ? n.allowAbsoluteUrls = this.defaults.allowAbsoluteUrls : n.allowAbsoluteUrls = !0), $n.assertOptions(n, {
      baseUrl: st.spelling("baseURL"),
      withXsrfToken: st.spelling("withXSRFToken")
    }, !0), n.method = (n.method || this.defaults.method || "get").toLowerCase();
    let o = i && y.merge(
      i.common,
      i[n.method]
    );
    i && y.forEach(
      ["delete", "get", "head", "post", "put", "patch", "common"],
      (m) => {
        delete i[m];
      }
    ), n.headers = je.concat(o, i);
    const l = [];
    let a = !0;
    this.interceptors.request.forEach(function (x) {
      typeof x.runWhen == "function" && x.runWhen(n) === !1 || (a = a && x.synchronous, l.unshift(x.fulfilled, x.rejected));
    });
    const u = [];
    this.interceptors.response.forEach(function (x) {
      u.push(x.fulfilled, x.rejected);
    });
    let c, f = 0, p;
    if (!a) {
      const m = [mi.bind(this), void 0];
      for (m.unshift(...l), m.push(...u), p = m.length, c = Promise.resolve(n); f < p;)
        c = c.then(m[f++], m[f++]);
      return c;
    }
    p = l.length;
    let g = n;
    for (; f < p;) {
      const m = l[f++], x = l[f++];
      try {
        g = m(g);
      } catch (A) {
        x.call(this, A);
        break;
      }
    }
    try {
      c = mi.call(this, g);
    } catch (m) {
      return Promise.reject(m);
    }
    for (f = 0, p = u.length; f < p;)
      c = c.then(u[f++], u[f++]);
    return c;
  }
  getUri(t) {
    t = Ft(this.defaults, t);
    const n = zo(t.baseURL, t.url, t.allowAbsoluteUrls);
    return Ho(n, t.params, t.paramsSerializer);
  }
};
y.forEach(["delete", "get", "head", "options"], function (t) {
  Pt.prototype[t] = function (n, s) {
    return this.request(Ft(s || {}, {
      method: t,
      url: n,
      data: (s || {}).data
    }));
  };
});
y.forEach(["post", "put", "patch"], function (t) {
  function n(s) {
    return function (i, o, l) {
      return this.request(Ft(l || {}, {
        method: t,
        headers: s ? {
          "Content-Type": "multipart/form-data"
        } : {},
        url: i,
        data: o
      }));
    };
  }
  Pt.prototype[t] = n(), Pt.prototype[t + "Form"] = n(!0);
});
let Ud = class Xo {
  constructor(t) {
    if (typeof t != "function")
      throw new TypeError("executor must be a function.");
    let n;
    this.promise = new Promise(function (i) {
      n = i;
    });
    const s = this;
    this.promise.then((r) => {
      if (!s._listeners) return;
      let i = s._listeners.length;
      for (; i-- > 0;)
        s._listeners[i](r);
      s._listeners = null;
    }), this.promise.then = (r) => {
      let i;
      const o = new Promise((l) => {
        s.subscribe(l), i = l;
      }).then(r);
      return o.cancel = function () {
        s.unsubscribe(i);
      }, o;
    }, t(function (i, o, l) {
      s.reason || (s.reason = new Gt(i, o, l), n(s.reason));
    });
  }
  /**
   * Throws a `CanceledError` if cancellation has been requested.
   */
  throwIfRequested() {
    if (this.reason)
      throw this.reason;
  }
  /**
   * Subscribe to the cancel signal
   */
  subscribe(t) {
    if (this.reason) {
      t(this.reason);
      return;
    }
    this._listeners ? this._listeners.push(t) : this._listeners = [t];
  }
  /**
   * Unsubscribe from the cancel signal
   */
  unsubscribe(t) {
    if (!this._listeners)
      return;
    const n = this._listeners.indexOf(t);
    n !== -1 && this._listeners.splice(n, 1);
  }
  toAbortSignal() {
    const t = new AbortController(), n = (s) => {
      t.abort(s);
    };
    return this.subscribe(n), t.signal.unsubscribe = () => this.unsubscribe(n), t.signal;
  }
  /**
   * Returns an object that contains a new `CancelToken` and a function that, when called,
   * cancels the `CancelToken`.
   */
  static source() {
    let t;
    return {
      token: new Xo(function (r) {
        t = r;
      }),
      cancel: t
    };
  }
};
function qd(e) {
  return function (n) {
    return e.apply(null, n);
  };
}
function Vd(e) {
  return y.isObject(e) && e.isAxiosError === !0;
}
const Gs = {
  Continue: 100,
  SwitchingProtocols: 101,
  Processing: 102,
  EarlyHints: 103,
  Ok: 200,
  Created: 201,
  Accepted: 202,
  NonAuthoritativeInformation: 203,
  NoContent: 204,
  ResetContent: 205,
  PartialContent: 206,
  MultiStatus: 207,
  AlreadyReported: 208,
  ImUsed: 226,
  MultipleChoices: 300,
  MovedPermanently: 301,
  Found: 302,
  SeeOther: 303,
  NotModified: 304,
  UseProxy: 305,
  Unused: 306,
  TemporaryRedirect: 307,
  PermanentRedirect: 308,
  BadRequest: 400,
  Unauthorized: 401,
  PaymentRequired: 402,
  Forbidden: 403,
  NotFound: 404,
  MethodNotAllowed: 405,
  NotAcceptable: 406,
  ProxyAuthenticationRequired: 407,
  RequestTimeout: 408,
  Conflict: 409,
  Gone: 410,
  LengthRequired: 411,
  PreconditionFailed: 412,
  PayloadTooLarge: 413,
  UriTooLong: 414,
  UnsupportedMediaType: 415,
  RangeNotSatisfiable: 416,
  ExpectationFailed: 417,
  ImATeapot: 418,
  MisdirectedRequest: 421,
  UnprocessableEntity: 422,
  Locked: 423,
  FailedDependency: 424,
  TooEarly: 425,
  UpgradeRequired: 426,
  PreconditionRequired: 428,
  TooManyRequests: 429,
  RequestHeaderFieldsTooLarge: 431,
  UnavailableForLegalReasons: 451,
  InternalServerError: 500,
  NotImplemented: 501,
  BadGateway: 502,
  ServiceUnavailable: 503,
  GatewayTimeout: 504,
  HttpVersionNotSupported: 505,
  VariantAlsoNegotiates: 506,
  InsufficientStorage: 507,
  LoopDetected: 508,
  NotExtended: 510,
  NetworkAuthenticationRequired: 511,
  WebServerIsDown: 521,
  ConnectionTimedOut: 522,
  OriginIsUnreachable: 523,
  TimeoutOccurred: 524,
  SslHandshakeFailed: 525,
  InvalidSslCertificate: 526
};
Object.entries(Gs).forEach(([e, t]) => {
  Gs[t] = e;
});
function Qo(e) {
  const t = new Pt(e), n = Ro(Pt.prototype.request, t);
  return y.extend(n, Pt.prototype, t, { allOwnKeys: !0 }), y.extend(n, t, null, { allOwnKeys: !0 }), n.create = function (r) {
    return Qo(Ft(e, r));
  }, n;
}
const ve = Qo(Tn);
ve.Axios = Pt;
ve.CanceledError = Gt;
ve.CancelToken = Ud;
ve.isCancel = Vo;
ve.VERSION = Yo;
ve.toFormData = ds;
ve.AxiosError = ee;
ve.Cancel = ve.CanceledError;
ve.all = function (t) {
  return Promise.all(t);
};
ve.spread = qd;
ve.isAxiosError = Vd;
ve.mergeConfig = Ft;
ve.AxiosHeaders = je;
ve.formToJSON = (e) => qo(y.isHTMLForm(e) ? new FormData(e) : e);
ve.getAdapter = Jo.getAdapter;
ve.HttpStatusCode = Gs;
ve.default = ve;
const {
  Axios: Xd,
  AxiosError: Qd,
  CanceledError: Zd,
  isCancel: eh,
  CancelToken: th,
  VERSION: nh,
  all: sh,
  Cancel: rh,
  isAxiosError: ih,
  spread: oh,
  toFormData: lh,
  AxiosHeaders: ah,
  HttpStatusCode: ch,
  formToJSON: uh,
  getAdapter: fh,
  mergeConfig: dh
} = ve;
class Zo {
}
class hh extends Zo {
  constructor(t = {}) {
    super();
    const n = t.baseURL || "";
    this.axiosInstance = ve.create({ baseURL: n });
  }
  /**
   * Load an entry into the backend session and get its metadata
   */
  async loadEntry(t) {
    const n = await this.axiosInstance.post("/api/load-entry", {
      id: t
    });
    return {
      recordId: n.data.record_id,
      filename: n.data.filename,
      recordInfo: {
        description: n.data.record_info.description
      }
    };
  }
  /**
   * Get a list of available records
   */
  async getRecords() {
    throw new Error("getRecords() is not yet implemented in the API");
  }
  /**
   * Get regions for a specific record
   */
  async getRegions(t) {
    return (await this.axiosInstance.get(
      `/api/records/${t}/regions`
    )).data;
  }
  /**
   * Get all features for a record (no region filtering)
   */
  async getRecordFeatures(t) {
    return (await this.axiosInstance.get(
      `/api/records/${t}/features`
    )).data;
  }
  /**
   * Get features for a specific region within a record
   */
  async getRegionFeatures(t, n) {
    return (await this.axiosInstance.get(
      `/api/records/${t}/regions/${n}/features`
    )).data;
  }
  /**
   * Get PFAM domain color mapping
   */
  async getPfamColorMap() {
    const s = (await this.axiosInstance.get("/domain-colors.csv")).data.split(`
`), r = {};
    for (let i = 1; i < s.length; i++) {
      const o = s[i].trim();
      if (o) {
        const [l, a] = o.split(",");
        l && a && (r[l] = a);
      }
    }
    return r;
  }
  /**
   * Get MiBIG entries for a specific locus_tag
   */
  async getMiBIGEntries(t, n, s = "1") {
    return (await this.axiosInstance.get(
      `/api/records/${t}/mibig-entries/${n}`,
      { params: { region: s } }
    )).data;
  }
  /**
   * Get TFBS finder binding site hits for a specific region
   */
  async getTFBSHits(t, n = "1") {
    return (await this.axiosInstance.get(
      `/api/records/${t}/tfbs-hits`,
      { params: { region: n } }
    )).data;
  }
  /**
   * Get TTA codon positions for a record
   */
  async getTTACodons(t) {
    return (await this.axiosInstance.get(
      `/api/records/${t}/tta-codons`
    )).data;
  }
  /**
   * Get resistance features for a record
   */
  async getResistanceFeatures(t) {
    return (await this.axiosInstance.get(
      `/api/records/${t}/resistance`
    )).data;
  }
}
class ph extends Zo {
  constructor(t = {}) {
    super(), this.records = t.records || [], this.pfamColorMap = t.pfamColorMap || {};
  }
  /**
   * Load data from a JSON file
   * @param file - File object or URL to JSON file
   */
  async loadFromFile(t) {
    let n;
    if (t instanceof File) {
      const s = await t.text();
      n = JSON.parse(s);
    } else if (typeof t == "string")
      n = await (await fetch(t)).json();
    else
      throw new Error("Invalid file parameter");
    this.parseAntiSMASHData(n);
  }
  /**
   * Parse antiSMASH JSON data structure
   */
  parseAntiSMASHData(t) {
    t.records ? this.records = t.records : this.records = [t];
  }
  /**
   * Load an entry (no-op for JSON file provider as data is already in memory)
   */
  async loadEntry(t) {
    const n = t.includes(":") ? t.split(":", 2)[1] : t, s = this.records.find((r) => r.id === n);
    if (!s)
      throw new Error(`Record ${n} not found`);
    return {
      recordId: s.id,
      filename: t.includes(":") ? t.split(":", 2)[0] : "unknown",
      recordInfo: {
        description: s.description || ""
      }
    };
  }
  /**
   * Get a list of available records
   */
  async getRecords() {
    return this.records.map((t, n) => ({
      recordId: t.id || `record-${n}`,
      filename: t.filename || `record-${n}.json`,
      recordInfo: {
        description: t.description || t.definition || ""
      }
    }));
  }
  /**
   * Get regions for a specific record
   */
  async getRegions(t) {
    const n = this.findRecord(t);
    if (!n)
      throw new Error(`Record not found: ${t}`);
    const s = [];
    return (n.features || []).forEach((i, o) => {
      var l, a, u;
      if (i.type === "region") {
        const c = this.parseLocation(i.location);
        s.push({
          id: `region-${o}`,
          region_number: ((a = (l = i.qualifiers) == null ? void 0 : l.region_number) == null ? void 0 : a[0]) || o + 1,
          product: ((u = i.qualifiers) == null ? void 0 : u.product) || [],
          start: c == null ? void 0 : c.start,
          end: c == null ? void 0 : c.end
        });
      }
    }), { regions: s };
  }
  /**
   * Get all features for a record (no region filtering)
   */
  async getRecordFeatures(t) {
    const n = this.findRecord(t);
    if (!n)
      throw new Error(`Record not found: ${t}`);
    return {
      features: n.features || []
    };
  }
  /**
   * Get features for a specific region within a record
   */
  async getRegionFeatures(t, n) {
    const s = this.findRecord(t);
    if (!s)
      throw new Error(`Record not found: ${t}`);
    const r = s.features || [], i = parseInt(n.replace("region-", "")), o = r[i];
    if (!o || o.type !== "region")
      throw new Error(`Region not found: ${n}`);
    const l = this.parseLocation(o.location);
    if (!l)
      throw new Error(`Invalid region location: ${o.location}`);
    return {
      features: r.filter((u) => {
        const c = this.parseLocation(u.location);
        return c ? !(c.end < l.start || c.start > l.end) : !1;
      }),
      region_boundaries: {
        start: l.start,
        end: l.end
      }
    };
  }
  /**
   * Get PFAM domain color mapping
   */
  async getPfamColorMap() {
    return this.pfamColorMap;
  }
  /**
   * Set PFAM color map (for loading from external source)
   */
  setPfamColorMap(t) {
    this.pfamColorMap = t;
  }
  /**
   * Find a record by ID
   */
  findRecord(t) {
    return this.records.find(
      (n, s) => (n.id || `record-${s}`) === t
    );
  }
  /**
   * Get MiBIG entries for a specific locus_tag
   */
  async getMiBIGEntries(t, n, s = "1") {
    const r = this.findRecord(t);
    if (!r)
      throw new Error(`Record not found: ${t}`);
    const a = (((r.modules || {})["antismash.modules.clusterblast"] || {}).knowncluster || {}).mibig_entries || {};
    if (!a[s])
      throw new Error(`No MiBIG entries available for region ${s}`);
    const u = a[s][n] || [];
    if (u.length === 0)
      throw new Error(`No MiBIG entries found for locus_tag '${n}' in region ${s}`);
    const c = u.map((f) => ({
      mibig_protein: f[0],
      description: f[1],
      mibig_cluster: f[2],
      rank: f[3],
      mibig_product: f[4],
      percent_identity: f[5],
      blast_score: f[6],
      percent_coverage: f[7],
      evalue: f[8]
    }));
    return {
      record_id: t,
      locus_tag: n,
      region: s,
      count: c.length,
      entries: c
    };
  }
  /**
   * Get TFBS finder binding site hits for a specific region
   */
  async getTFBSHits(t, n = "1") {
    const s = this.findRecord(t);
    if (!s)
      throw new Error(`Record not found: ${t}`);
    const o = ((s.modules || {})["antismash.modules.tfbs_finder"] || {}).hits_by_region || {};
    if (!o[n])
      return {
        record_id: t,
        region: n,
        count: 0,
        hits: []
      };
    const l = o[n] || [];
    return {
      record_id: t,
      region: n,
      count: l.length,
      hits: l
    };
  }
  /**
   * Get TTA codon positions for a record
   */
  async getTTACodons(t) {
    const n = this.findRecord(t);
    if (!n)
      throw new Error(`Record not found: ${t}`);
    const i = ((n.modules || {})["antismash.modules.tta"] || {})["TTA codons"] || [];
    return {
      record_id: t,
      count: i.length,
      codons: i
    };
  }
  /**
   * Get resistance features for a record
   */
  async getResistanceFeatures(t) {
    const n = this.findRecord(t);
    if (!n)
      throw new Error(`Record not found: ${t}`);
    const l = ((((n.modules || {})["antismash.detection.genefunctions"] || {}).tools || {}).resist || {}).best_hits || {}, a = [];
    for (const [u, c] of Object.entries(l))
      a.push({
        locus_tag: u,
        ...c
      });
    return {
      record_id: t,
      count: a.length,
      features: a
    };
  }
  /**
   * Parse location string like "[164:2414](+)" or "[257:2393](+)"
   */
  parseLocation(t) {
    if (!t) return null;
    const n = t.match(/\[<?(\d+):>?(\d+)\](?:\(([+-])\))?/);
    return n ? {
      start: parseInt(n[1]),
      end: parseInt(n[2]),
      strand: n[3] || null
    } : null;
  }
}
const Wd = /* @__PURE__ */ Co(gf), zd = /* @__PURE__ */ Co(ko);
customElements.define("bgc-region-viewer-container", Wd);
customElements.define("bgc-region-viewer", zd);
export {
  hh as BGCViewerAPIProvider,
  ph as JSONFileProvider,
  zd as RegionViewer,
  Wd as RegionViewerContainer,
  Vs as TrackViewer
};
