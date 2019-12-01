package money

// Money represents an amount of money for a currency
//
// Money is any item or verifiable record that is generally accepted as payment
// for goods and services and repayment of debts
type Money struct {
	Amount   Decimal  `json:"amount"`
	Currency Currency `json:"currency"`
}

// MustParse is like Parse, but panics if the given amount or currency cannot
// be parsed. It simplifies safe initialisation of Money values.
func MustParse(amount, currency string) *Money {
	m, err := Parse(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// Parse parses amount which must contain a text representation of a
// floating-point number.
// The number of integers after the radix point (fraction) determines the
// mantissa precision.
//
//   e.g. 120.0 	-> Precision 1
//   e.g. 123.456	-> Precision 3
//
// It also validates the currency, which must represented in code as defined by
// the ISO 4217 format.
//
//   e.g. CHF 		-> Swiss franc
func Parse(amount, currency string) (*Money, error) {
	a, err := ParseDecimal(amount)
	if err != nil {
		return nil, err
	}
	c, err := ParseCurrency(currency)
	if err != nil {
		return nil, err
	}
	return &Money{
		Amount:   a,
		Currency: c,
	}, nil
}

// Equal tests whether y equal x. When the currency is different, it will
// always return false. Currency conversion is currently not supported.
func (x *Money) Equal(y *Money) bool {
	if x.Currency != y.Currency {
		return false
	}
	return x.Amount.Equal(y.Amount)
}

// Validate tests that both the decimal and the currency are valid
func (x *Money) Validate() error {
	if err := x.Currency.Validate(); err != nil {
		return err
	}
	return x.Amount.Validate()
}

// Add returns an amount set to the rounded sum x+y.
// The precision is set to the larger of x's or y's precision before the
// operation.
// Rounding is performed according to the default rounding mode
func Add(x, y *Money) *Money {
	z := Money{}
	return &z
}

// Sub returns an amount set to the rounded difference x-y.
// Precision, rounding, and accuracy reporting are as for Add.
// Sub panics with ErrNaN if x and y are infinities with equal
// signs.
func Sub(x, y *Money) *Money {
	z := Money{}
	return &z
}

// Mul sets z to the rounded product x*y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Mul panics with ErrNaN if one operand is zero and the other
// operand an infinity.
func Mul(x, y *Money) *Money {
	z := Money{}
	return &z
}

// Div sets z to the rounded quotient x/y and returns z.
// Precision, rounding, and accuracy reporting are as for Add.
// Quo panics with ErrNaN if both operands are zero or infinities.
func Div(x, y *Money) *Money {
	z := Money{}
	return &z
}
