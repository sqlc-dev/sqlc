package postgres

func args(a ...int) map[int]struct{} {
	m := map[int]struct{}{}
	for _, arg := range a {
		m[arg] = struct{}{}
	}
	return m
}

var Functions = map[string]map[int]struct{}{

	// https://www.postgresql.org/docs/current/functions-math.html
	"abs":     args(1), // (x) 	(same as input) 	absolute value 	abs(-17.4) 	17.4
	"cbrt":    args(1), // (dp) 	dp 	cube root 	cbrt(27.0) 	3
	"ceil":    args(1), // (dp or numeric) 	(same as input) 	nearest integer greater than or equal to argument 	ceil(-42.8) 	-42
	"ceiling": args(1), // (dp or numeric) 	(same as input) 	nearest integer greater than or equal to argument (same as ceil) 	ceiling(-95.3) 	-95
	"degrees": args(1), // (dp) 	dp 	radians to degrees 	degrees(0.5) 	28.6478897565412
	"div":     args(2), // (y numeric, x numeric) 	numeric 	integer quotient of y/x 	div(9,4) 	2
	"exp":     args(1), // (dp or numeric) 	(same as input) 	exponential 	exp(1.0) 	2.71828182845905
	"floor":   args(1), // (dp or numeric) 	(same as input) 	nearest integer less than or equal to argument 	floor(-42.8) 	-43
	"ln":      args(1), // (dp or numeric) 	(same as input) 	natural logarithm 	ln(2.0) 	0.693147180559945
	"log":     args(1, 2),
	// (dp or numeric) 	(same as input) 	base 10 logarithm 	log(100.0) 	2
	// (b numeric, x numeric) 	numeric 	logarithm to base b 	log(2.0, 64.0) 	6.0000000000
	"mod":   args(2), // (y, x) 	(same as argument types) 	remainder of y/x 	mod(9,4) 	1
	"pi":    args(0), // () 	dp 	“π” constant 	pi() 	3.14159265358979
	"power": args(2),
	// (a dp, b dp) 	dp 	a raised to the power of b 	power(9.0, 3.0) 	729
	// power(a numeric, b numeric) 	numeric 	a raised to the power of b 	power(9.0, 3.0) 	729
	"radians": args(1),    // (dp) 	dp 	degrees to radians 	radians(45.0) 	0.785398163397448
	"round":   args(1, 2), // (dp or numeric) 	(same as input) 	round to nearest integer 	round(42.4) 	42
	// round(v numeric, s int) 	numeric 	round to s decimal places 	round(42.4382, 2) 	42.44
	"scale": args(1),    // (numeric) 	integer 	scale of the argument (the number of decimal digits in the fractional part) 	scale(8.41) 	2
	"sign":  args(1),    // (dp or numeric) 	(same as input) 	sign of the argument (-1, 0, +1) 	sign(-8.4) 	-1
	"sqrt":  args(1),    // (dp or numeric) 	(same as input) 	square root 	sqrt(2.0) 	1.4142135623731
	"trunc": args(1, 2), // (dp or numeric) 	(same as input) 	truncate toward zero 	trunc(42.8) 	42
	// trunc(v numeric, s int) 	numeric 	truncate to s decimal places 	trunc(42.4382, 2) 	42.43
	"width_bucket": args(2, 4), // (operand dp, b1 dp, b2 dp, count int) 	int 	return the bucket number to which operand would be assigned in a histogram having count equal-width buckets spanning the range b1 to b2; returns 0 or count+1 for an input outside the range 	width_bucket(5.35, 0.024, 10.06, 5) 	3
	// width_bucket(operand numeric, b1 numeric, b2 numeric, count int) 	int 	return the bucket number to which operand would be assigned in a histogram having count equal-width buckets spanning the range b1 to b2; returns 0 or count+1 for an input outside the range 	width_bucket(5.35, 0.024, 10.06, 5) 	3
	// width_bucket(operand anyelement, thresholds anyarray) 	int 	return the bucket number to which operand would be assigned given an array listing the lower bounds of the buckets; returns 0 for an input less than the first lower bound; the thresholds array must be sorted, smallest first, or unexpected results will be obtained 	width_bucket(now(), array['yesterday', 'today', 'tomorrow']::timestamptz[]) 	2

	// https://www.postgresql.org/docs/current/functions-datetime.html
	"now": args(0),

	// https://www.postgresql.org/docs/current/functions-aggregate.html
	"count": args(0, 1),
}
