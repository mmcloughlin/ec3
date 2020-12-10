# ec3

Elliptic Curve Cryptography Compiler.

`ec3` is an _incomplete experiment_ that set out with the goal of
auto-generating Go elliptic curve implementations for any parameter choices,
with performance as close as possible to hand-crafted implementations.

In its current state, `ec3` can generate a full implementation of the NIST
P-256 elliptic curve, however much more work is required to reach performance
parity with the standard library implementation. It's unpolished and comes
with no guarantees: `ec3` is shared as-is in the hopes that the approach is
interesting to others, or some subcomponents are independently useful. In
particular:

* The [`addchain`](https://github.com/mmcloughlin/addchain) project started
  life as part of `ec3`. This library generates short addition chains required
  for optimization of finite field exponentiation. The results rival or
  sometimes beat the best known hand-optimized chains.
* The `efd` package embeds the [Explicit-Formulas
  Database](https://hyperelliptic.org/EFD) and provides Go libraries for
  manipulating formulae in op3 format.

If you are interested in the `ec3` approach, you should check out
[ECCKiila](https://gitlab.com/nisec/ecckiila) which used
[`fiat-crypto`](https://github.com/mit-plv/fiat-crypto) for finite field
operations. Work is also ongoing to [import `fiat-crypto` implementations in
the Go standard library](https://golang.org/issue/40171).
