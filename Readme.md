Simple motor calculator based on http://www.chiefdelphi.com/media/papers/2432 but rewritten in Go for cross platform use.

You can use the tool with `mcalc-[os] [motor name] [amps|torque] [value]`

So, for example, to compute the motor values for a CIM drawing 30A:
```
mcalc-osx CIM amps 30

torque(N*m)	rpm	amps	output(W)	heat(W)	eff
0.51		4196	30.00	225.32		134.68	62.59
```
