# ec3

Elliptic Curve Cryptography Compiler.

`ec3` is an _incomplete experiment_ that set out with the goal of
auto-generating Go elliptic curve implementations for any parameter choices,
with performance as close as possible to hand-crafted implementations.

In its current state, `ec3` can generate a full implementation of the NIST
P-256 elliptic curve, however much more work is required to reach performance
parity with the standard library implementation. The project is public to
share work-in-progress, and in the hopes that some subcomponents may be
useful to others.
