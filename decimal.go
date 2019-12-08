package money

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/number"
)

// TODO: Replace with https://github.com/cockroachdb/apd
// TODO: Use currency to define the precision
// TODO: Do not require to define a precision

// divisionPrecision is the number of decimal places in the result when it
// doesn't divide exactly.
//
// Example:
//
//     d1 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3)
//     d1.String() // output: "0.6666666666666667"
//     d2 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(30000)
//     d2.String() // output: "0.0000666666666667"
//     d3 := decimal.NewFromFloat(20000).Div(decimal.NewFromFloat(3)
//     d3.String() // output: "6666.6666666666666667"
//     decimal.DivisionPrecision = 3
//     d4 := decimal.NewFromFloat(2).Div(decimal.NewFromFloat(3)
//     d4.String() // output: "0.667"
//
var divisionPrecision = 16

// marshalJSONWithoutQuotes should be set to true if you want the decimal to
// be JSON marshaled as a number, insteaddof as a string.
// WARNING: this is dangerous for decimals with many digits, since many JSON
// unmarshallers (ex: Javascript's) will unmarshal JSON numbers to IEEE 754
// double-precision floating point numbers, which means you can potentially
// silently lose precision.
var marshalJSONWithoutQuotes = false

// decSeparator is the decimal separator symbol
const decSeparator = "."

// allowedDecimalRunes contains the list of allowed runes that are not a decimal
// in a string representation
var allowedDecimalRunes = []rune{'+', '-', '.'}

const (
	// SignPositive is the number returned by Sign() when a decimal is positive
	SignPositive = 1
	// SignNeutral is the number returned by Sign() when a decimal is neutral (0)
	SignNeutral = 0
	// SignNegative is the number returned by Sign() when a decimal is negative
	SignNegative = -1
)

var (
	zero    = buildDecimal(0, 1)
	hundred = buildDecimal(100, 0)

	zeroInt = big.NewInt(0)
	oneInt  = big.NewInt(1)
	fiveInt = big.NewInt(5)
	tenInt  = big.NewInt(10)
)

var (
	// ErrInvalidDecimal indicates that the string is not a valid decimal
	ErrInvalidDecimal = errors.New("invalid decimal")
)

// Decimal represents a fixed-point decimal. It is immutable.
// number = value * 10 ^ exp
type Decimal struct {
	value big.Int
	exp   int32
}

func MustParseDecimal(value string) Decimal {
	d, err := ParseDecimal(value)
	if err != nil {
		panic(err)
	}
	return d
}

// ParseDecimal parses the value which must contain a text representation of a floating-point number.
// The number of integers after the radix point (fraction) determines the rounding precision.
//
//   e.g. 120.0 	-> Precision 1
//   e.g. 123.456	-> Precision 3
func ParseDecimal(value string) (Decimal, error) {
	var ints string
	var exp int64

	// Check format.
	// It avoids to parse valid big int values, such as:
	//  - exponents
	//  - infinity
	//  - base 2, 16, ...
	for _, c := range value {
		if unicode.IsDigit(c) {
			continue
		}

		var allowed bool
		for _, r := range allowedDecimalRunes {
			if c == r {
				allowed = true
				break
			}
		}

		if !allowed {
			return zero, ErrInvalidDecimal
		}
	}

	parts := strings.Split(value, decSeparator)
	switch len(parts) {
	case 1:
		ints = parts[0]
		exp = 0
	case 2:
		// strip the insignificant digits for more accurate comparisons.
		ints = parts[0] + parts[1]
		expInt := -len(parts[1])
		exp += int64(expInt)
	default:
		return zero, ErrInvalidDecimal
	}

	dValue := new(big.Int)
	if _, ok := dValue.SetString(ints, 10); !ok {
		return zero, ErrInvalidDecimal
	}
	if exp < math.MinInt32 || exp > math.MaxInt32 {
		return zero, ErrInvalidDecimal
	}

	return Decimal{
		value: *dValue,
		exp:   int32(exp),
	}, nil
}

// NewDecimal creates a Decimal from a float
//
// Example:
//
//     NewFromFloat(123.45678901234567).String() // output: "123.4567890123456"
//     NewFromFloat(.00000000000000001).String() // output: "0.00000000000000001"
//
// NOTE: errors occur on NaN, +/-inf
func NewDecimal(value float64) (Decimal, error) {
	floor := math.Floor(value)

	// fast path, where float is an int
	if floor == value && value <= math.MaxInt64 && value >= math.MinInt64 {
		return buildDecimal(int64(value), 0), nil
	}

	// TODO: Avoid the string conversion
	str := strconv.FormatFloat(value, 'f', -1, 64)
	dec, err := ParseDecimal(str)
	if err != nil {
		return zero, ErrInvalidDecimal
	}
	return dec, nil
}

// MinDecimal returns the smallest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Min(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Min with 0 arguments.
func MinDecimal(first Decimal, rest ...Decimal) Decimal {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) < 0 {
			ans = item
		}
	}
	return ans
}

// MaxDecimal returns the largest Decimal that was passed in the arguments.
//
// To call this function with an array, you must do:
//
//     Max(arr[0], arr[1:]...)
//
// This makes it harder to accidentally call Max with 0 arguments.
func MaxDecimal(first Decimal, rest ...Decimal) Decimal {
	ans := first
	for _, item := range rest {
		if item.Cmp(ans) > 0 {
			ans = item
		}
	}
	return ans
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	d2Value := new(big.Int).Abs(&d.value)
	return Decimal{
		value: *d2Value,
		exp:   d.exp,
	}
}

// Add returns d + d2.
func (d Decimal) Add(d2 Decimal) Decimal {
	baseScale := min(d.exp, d2.exp)
	rd := d.rescale(baseScale)
	rd2 := d2.rescale(baseScale)

	d3Value := new(big.Int).Add(&rd.value, &rd2.value)
	return Decimal{
		value: *d3Value,
		exp:   baseScale,
	}
}

// Sub returns d - d2.
func (d Decimal) Sub(d2 Decimal) Decimal {
	baseScale := min(d.exp, d2.exp)
	rd := d.rescale(baseScale)
	rd2 := d2.rescale(baseScale)

	d3Value := new(big.Int).Sub(&rd.value, &rd2.value)
	return Decimal{
		value: *d3Value,
		exp:   baseScale,
	}
}

// Mul returns d * d2.
func (d Decimal) Mul(d2 Decimal) Decimal {
	expInt64 := int64(d.exp) + int64(d2.exp)
	if expInt64 > math.MaxInt32 || expInt64 < math.MinInt32 {
		// better to panic than give incorrect results, as
		// Decimals are usually used for money
		panic(fmt.Sprintf("exponent %v overflows an int32!", expInt64))
	}

	d3Value := new(big.Int).Mul(&d.value, &d2.value)
	return Decimal{
		value: *d3Value,
		exp:   int32(expInt64),
	}
}

// Div returns d / d2. If it doesn't divide exactly, the result will have
// DivisionPrecision digits after the decimal point.
func (d Decimal) Div(d2 Decimal) Decimal {
	return d.divRound(d2, int32(divisionPrecision))
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	val := new(big.Int).Neg(&d.value)
	return Decimal{
		value: *val,
		exp:   d.exp,
	}
}

// Mod returns d % d2.
func (d Decimal) Mod(d2 Decimal) Decimal {
	quo := d.Div(d2).Truncate(0)
	return d.Sub(d2.Mul(quo))
}

// Round rounds the decimal to places decimal places.
// If places < 0, it will round the integer part to the nearest 10^(-places).
//
// Example:
//
// 	   NewFromFloat(5.45).Round(1).String() // output: "5.5"
// 	   NewFromFloat(545).Round(-1).String() // output: "550"
//
func (d Decimal) Round(places int32) Decimal {
	// truncate to places + 1
	ret := d.rescale(-places - 1)

	// add sign(d) * 0.5
	if ret.value.Sign() == SignNegative {
		ret.value.Sub(&ret.value, fiveInt)
	} else {
		ret.value.Add(&ret.value, fiveInt)
	}

	// floor for positive numbers, ceil for negative numbers
	_, m := ret.value.DivMod(&ret.value, tenInt, new(big.Int))
	ret.exp++
	if ret.value.Sign() == SignNegative && m.Cmp(zeroInt) != 0 {
		ret.value.Add(&ret.value, oneInt)
	}

	return ret
}

// RoundUp rounds the decimal up to the given precision instead of to the nearest even
//
//	e.g.:
// 	3.1416 -> f(3) = 3.142
// 	3.1416 -> f(2) = 3.15
//
func (d Decimal) RoundUp(precision int32) Decimal {
	if d.Round(precision).Equal(d) {
		return d
	}

	halfPrecision := buildDecimal(5, -precision-1)
	return d.Add(halfPrecision).Round(precision)
}

// RoundDown rounds the decimal down to the given precision instead of to the nearest even
//
//	e.g.:
// 	3.1416 -> f(3) = 3.142
// 	3.1416 -> f(2) = 3.15
//
func (d Decimal) RoundDown(precision int32) Decimal {
	if d.Round(precision).Equal(d) {
		return d
	}

	halfPrecision := buildDecimal(-5, -precision-1)
	return d.Add(halfPrecision).Round(precision)
}

// RoundNearest rounds the decimal to the nearest unit
//
//	e.g.:
// 	3.1216 -> f(0.05) = 3.10
// 	3.1416 -> f(0.05) = 3.15
//
func (d Decimal) RoundNearest(unit Decimal) Decimal {
	// First round to the unit precision
	rounded := d.Round(-unit.exp)

	// Then move to the nearest unit
	remainder := rounded.Mod(unit)

	// Inverse signs if it is a negative decimal
	var cmp = 1
	if rounded.Sign() == SignNegative {
		cmp = -1
		unit = unit.Neg()
	}

	// Round up
	if remainder.Cmp(unit.Div(buildDecimal(2, 0))) == cmp {
		return rounded.Add(unit.Sub(remainder))
	}
	// Round down
	return rounded.Sub(remainder)
}

// Truncate truncates off digits from the number, without rounding.
//
// NOTE: precision is the last digit that will not be truncated (must be >= 0).
//
// Example:
//
//     decimal.NewFromString("123.456").Truncate(2).String() // "123.45"
//
func (d Decimal) Truncate(precision int32) Decimal {
	if precision >= 0 && -precision > d.exp {
		return d.rescale(-precision)
	}
	return d
}

// Floor returns the nearest integer value less than or equal to d.
func (d Decimal) Floor() Decimal {
	exp := big.NewInt(10)

	// must negate after casting to prevent int32 overflow
	exp.Exp(exp, big.NewInt(-int64(d.exp)), nil)

	z := new(big.Int).Div(&d.value, exp)
	return Decimal{value: *z, exp: 0}
}

// Ceil returns the nearest integer value greater than or equal to d.
func (d Decimal) Ceil() Decimal {
	exp := big.NewInt(10)

	// must negate after casting to prevent int32 overflow
	exp.Exp(exp, big.NewInt(-int64(d.exp)), nil)

	z, m := new(big.Int).DivMod(&d.value, exp, new(big.Int))
	if m.Cmp(zeroInt) != 0 {
		z.Add(z, oneInt)
	}
	return Decimal{value: *z, exp: 0}
}

// Cmp compares the numbers represented by d and d2 and returns:
//
//     -1 if d <  d2
//      0 if d == d2
//     +1 if d >  d2
//
func (d Decimal) Cmp(d2 Decimal) int {
	if d.exp == d2.exp {
		return d.value.Cmp(&d2.value)
	}

	// Ensure both decimals are on the same scale
	baseExp := min(d.exp, d2.exp)
	var rd, rd2 Decimal
	if d.exp != baseExp {
		rd = d.rescale(baseExp)
		rd2 = d2
	} else if d2.exp != baseExp {
		rd = d
		rd2 = d2.rescale(baseExp)
	}

	return rd.value.Cmp(&rd2.value)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 Decimal) bool {
	return d.Cmp(d2) == 0
}

// IsZero reports whether d represents the zero value
func (d Decimal) IsZero() bool {
	return d.Cmp(zero) == 0
}

// Sign returns:
//
//	-1 if d <  0
//	 0 if d == 0
//	+1 if d >  0
//
func (d Decimal) Sign() int {
	return d.value.Sign()
}

// Exponent returns the exponent, or scale component of the decimal.
func (d Decimal) Exponent() int32 {
	return d.exp
}

// Coefficient returns the coefficient of the decimal. It is scaled by 10^Exponent()
func (d Decimal) Coefficient() big.Int {
	return d.value
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	scaledD := d.rescale(0)
	return scaledD.value.Int64()
}

// Rat returns a rational number representation of the decimal.
func (d Decimal) Rat() *big.Rat {
	if d.exp <= 0 {
		// must negate after casting to prevent int32 overflow
		denom := new(big.Int).Exp(tenInt, big.NewInt(-int64(d.exp)), nil)
		return new(big.Rat).SetFrac(&d.value, denom)
	}

	mul := new(big.Int).Exp(tenInt, big.NewInt(int64(d.exp)), nil)
	num := new(big.Int).Mul(&d.value, mul)
	return new(big.Rat).SetFrac(num, oneInt)
}

// Float64 returns the nearest float64 value for d
func (d Decimal) Float64() float64 {
	f, _ := d.Rat().Float64()
	return f
}

func (d Decimal) String() string {
	if d.exp >= 0 {
		v := d.rescale(0).value
		intPart := v.String()

		var number bytes.Buffer
		var prec int
		if d.roundPrec() > 0 {
			prec = int(d.roundPrec())
		} else {
			prec = 1
		}
		number.WriteString(intPart)
		number.WriteString(decSeparator)
		number.WriteString(strings.Repeat("0", prec))
		return number.String()
	}

	abs := new(big.Int).Abs(&d.value)
	str := abs.String()

	// this cast to int will cause bugs if d.exp == INT_MIN it is a 32-bit machine
	var intPart, fractionalPart string
	dExpInt := int(d.exp)
	if len(str) > -dExpInt {
		intPart = str[:len(str)+dExpInt]
		fractionalPart = str[len(str)+dExpInt:]
	} else {
		intPart = "0"
		num0s := -dExpInt - len(str)
		fractionalPart = strings.Repeat("0", num0s) + str
	}

	var number bytes.Buffer
	number.WriteString(intPart)

	if len(fractionalPart) > 0 {
		number.WriteString(decSeparator)
		number.WriteString(fractionalPart)
	}

	if d.value.Sign() == SignNegative {
		return "-" + number.String()
	}
	return number.String()
}

// Formatter returns a language/currency-specific formatter for a
// floating point decimal
func (d *Decimal) Formatter(scale ...int) number.Formatter {
	var s int
	if len(scale) > 0 {
		s = scale[0]
	} else {
		s = int(d.roundPrec())
	}

	return number.Decimal(
		d.Float64(),
		number.Scale(s),
	)
}

// PercentFormatter returns a language-specific formatter for a percent
func (d *Decimal) PercentFormatter() number.Formatter {
	return number.Percent(
		d.Div(hundred).Float64(),
		number.MaxFractionDigits(int(d.roundPrec()+2)), // +2 because div by hundred
	)
}

// Validate returns whether the currency is valid
func (d Decimal) Validate() error {
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		s := string(data[1 : len(data)-1])
		decimal, err := ParseDecimal(s)
		if err != nil {
			return fmt.Errorf("Error parsing money/decimal '%s': %s", s, err)
		}
		*d = decimal
		return nil
	}

	// Accept empty data. The Validate function should be used to make sure it
	// is valid
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface. As a string representation
// is already used when encoding to text, this method stores that string as []byte
func (d *Decimal) UnmarshalBinary(data []byte) error {
	// Extract the exponent
	d.exp = int32(binary.BigEndian.Uint32(data[:4]))

	// Extract the value
	d.value = big.Int{}
	return d.value.GobDecode(data[4:])
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (d Decimal) MarshalBinary() (data []byte, err error) {
	// Write the exponent first since it's a fixed size
	v1 := make([]byte, 4)
	binary.BigEndian.PutUint32(v1, uint32(d.exp))

	// Add the value
	var v2 []byte
	if v2, err = d.value.GobEncode(); err != nil {
		return data, err
	}

	// Return the byte array
	data = append(v1, v2...)
	return data, err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (d *Decimal) UnmarshalText(text []byte) error {
	str := string(text)

	dec, err := ParseDecimal(str)
	*d = dec
	if err != nil {
		return fmt.Errorf("Error decoding string '%s': %s", str, err)
	}

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d Decimal) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (d Decimal) GobEncode() ([]byte, error) {
	return d.MarshalBinary()
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (d *Decimal) GobDecode(data []byte) error {
	return d.UnmarshalBinary(data)
}

func (d Decimal) DeepCopy(dst interface{}) error {
	switch dst := dst.(type) {
	case *Decimal:
		dst.value = d.value
		dst.exp = d.exp
		return nil
	case Decimal:
		dst.value = d.value
		dst.exp = d.exp
		return nil
	}
	return fmt.Errorf("Decimal deep copy on an unknown type %T", dst)
}

func buildDecimal(value int64, exp int32) Decimal {
	return Decimal{
		value: *big.NewInt(value),
		exp:   exp,
	}
}

// Pow returns d to the power d2
func (d Decimal) Pow(d2 Decimal) Decimal {
	var temp Decimal
	if d2.IntPart() == 0 {
		x, err := NewDecimal(1)
		if err != nil {
			panic(err)
		}
		return x
	}

	x, err := NewDecimal(2)
	if err != nil {
		panic(err)
	}
	temp = d.Pow(d2.Div(x))
	if d2.IntPart()%2 == 0 {
		return temp.Mul(temp)
	}
	if d2.IntPart() > 0 {
		return temp.Mul(temp).Mul(d)
	}
	return temp.Mul(temp).Div(d)
}

// divRound divides and rounds to a given precision
// i.e. to an integer multiple of 10^(-precision)
//   for a positive quotient digit 5 is rounded up, away from 0
//   if the quotient is negative then digit 5 is rounded down, away from 0
// Note that precision<0 is allowed as input.
func (d Decimal) divRound(d2 Decimal, precision int32) Decimal {
	// QuoRem already checks initialization
	q, r := d.quoRem(d2, precision)
	// the actual rounding decision is based on comparing r*10^precision and d2/2
	// instead compare 2 r 10 ^precision and d2
	var rv2 big.Int
	rv2.Abs(&r.value)
	rv2.Lsh(&rv2, 1)
	// now rv2 = abs(r.value) * 2
	r2 := Decimal{value: rv2, exp: r.exp + precision}
	// r2 is now 2 * r * 10 ^ precision
	var c = r2.Cmp(d2.Abs())

	if c < 0 {
		return q
	}

	if d.value.Sign()*d2.value.Sign() < SignNegative {
		return q.Sub(buildDecimal(1, -precision))
	}

	return q.Add(buildDecimal(1, -precision))
}

// quoRem does divsion with remainder
// d.QuoRem(d2,precision) returns quotient q and remainder r such that
//   d = d2 * q + r, q an integer multiple of 10^(-precision)
//   0 <= r < abs(d2) * 10 ^(-precision) if d>=0
//   0 >= r > -abs(d2) * 10 ^(-precision) if d<0
// Note that precision<0 is allowed as input.
func (d Decimal) quoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	if d2.value.Sign() == SignNeutral {
		panic("decimal division by 0")
	}
	scale := -precision
	e := int64(d.exp - d2.exp - scale)
	if e > math.MaxInt32 || e < math.MinInt32 {
		panic("overflow in decimal QuoRem")
	}
	var aa, bb, expo big.Int
	var scalerest int32
	// d = a 10^ea
	// d2 = b 10^eb
	if e < 0 {
		aa = d.value
		expo.SetInt64(-e)
		bb.Exp(tenInt, &expo, nil)
		bb.Mul(&d2.value, &bb)
		scalerest = d.exp
		// now aa = a
		//     bb = b 10^(scale + eb - ea)
	} else {
		expo.SetInt64(e)
		aa.Exp(tenInt, &expo, nil)
		aa.Mul(&d.value, &aa)
		bb = d2.value
		scalerest = scale + d2.exp
		// now aa = a ^ (ea - eb - scale)
		//     bb = b
	}
	var q, r big.Int
	q.QuoRem(&aa, &bb, &r)
	dq := Decimal{value: q, exp: scale}
	dr := Decimal{value: r, exp: scalerest}
	return dq, dr
}

// rescale returns a rescaled version of the decimal. Returned
// decimal may be less precise if the given exponent is bigger
// than the initial exponent of the Decimal.
// NOTE: this will truncate, NOT round
//
// Example:
//
// 	d := New(12345, -4)
//	d2 := d.rescale(-1)
//	d3 := d2.rescale(-4)
//	println(d1)
//	println(d2)
//	println(d3)
//
// Output:
//
//	1.2345
//	1.2
//	1.2000
//
func (d Decimal) rescale(exp int32) Decimal {
	// must convert exps to float64 before - to prevent overflow
	diff := math.Abs(float64(exp) - float64(d.exp))
	value := new(big.Int).Set(&d.value)

	expScale := new(big.Int).Exp(tenInt, big.NewInt(int64(diff)), nil)
	if exp > d.exp {
		value = value.Quo(value, expScale)
	} else if exp < d.exp {
		value = value.Mul(value, expScale)
	}

	return Decimal{
		value: *value,
		exp:   exp,
	}
}

func (d *Decimal) roundPrec() uint {
	if d.exp < 0 {
		return uint(d.exp * -1)
	}
	return 0
}

func min(x, y int32) int32 {
	if x >= y {
		return y
	}
	return x
}
