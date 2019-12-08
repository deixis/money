package money

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/text/currency"
)

// TODO: Refactor to use directly golang.org/x/text/currency

var (
	// ErrInvalidCurrency indicates that the string is not a valid currency as defined by ISO 4217
	ErrInvalidCurrency = errors.New("invalid currency")
	// ErrUnsupportedCurrency indicates that the currency is not supported
	ErrUnsupportedCurrency = errors.New("unsupported currency")
)

// Currency is represented in code as defined by the ISO 4217 format.
//
// Examples:
//   * Swiss franc - CHF
//   * United States dollar - USD
//   * Euro - EUR
//   * Polish złoty - PLN
//   * Bitcoin - XBT
type Currency string

// MustParseCurrency is like ParseCurency, but panics if the given currency
// unit cannot be parsed. It simplifies safe initialisation of Currency values.
func MustParseCurrency(s string) Currency {
	c, err := ParseCurrency(s)
	if err != nil {
		panic(err)
	}
	return c
}

// ParseCurrency parses a 3-letter ISO 4217 currency code. It returns an error
// if it is not well-formed or not a recognised currency code.
//
// Examples:
//   * CHF
//   * XBT
func ParseCurrency(s string) (Currency, error) {
	// Sanitise
	s = strings.TrimSpace(strings.ToUpper(s))

	if _, ok := unoficialCurrencies.Load(s); ok {
		return Currency(s), nil
	}
	u, err := currency.ParseISO(s)
	if err != nil {
		return nullCurrency, ErrInvalidCurrency
	}
	return Currency(u.String()), nil
}

// Scale returns the standard currency scale
func (c Currency) Scale() int {
	scale, _ := currency.Kind(RoundingStandard.kind()).Rounding(*c.currency())
	return scale
}

// RoundUnit returns a rounding unit for the given kind
func (c Currency) RoundUnit(kind RoundingKind) Decimal {
	// Get rounding for the currency
	scale, inc := currency.Kind(kind.kind()).Rounding(*c.currency())
	return buildDecimal(int64(inc), int32(scale*-1))
}

// String returns the ISO 4217 representation of a currency (e.g. CHF)
func (c Currency) String() string {
	return string(c)
}

// Validate returns whether the currency is valid
func (c Currency) Validate() error {
	_, err := ParseCurrency(string(c))
	return err
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Currency) UnmarshalJSON(data []byte) error {
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		s := string(data[1 : len(data)-1])
		currency, err := ParseCurrency(s)
		if err != nil {
			return fmt.Errorf("Error parsing money/currency '%s': %s", s, err)
		}
		*c = currency
		return nil
	}

	// Accept empty data. The Validate function should be used to make sure it
	// is valid
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (c Currency) MarshalJSON() ([]byte, error) {
	return []byte("\"" + c + "\""), nil
}

// GobEncode implements the gob.GobEncoder interface for gob serialization.
func (c Currency) GobEncode() ([]byte, error) {
	return []byte(c), nil
}

// GobDecode implements the gob.GobDecoder interface for gob serialization.
func (c *Currency) GobDecode(data []byte) error {
	cur, err := ParseCurrency(string(data))
	if err != nil {
		return err
	}
	*c = cur
	return nil
}

func (c *Currency) currency() *currency.Unit {
	cur, err := currency.ParseISO(c.String())
	if err != nil {
		panic("invalid currency unit: " + err.Error())
	}
	return &cur
}

const (
	nullCurrency Currency = ""
)

var unoficialCurrencies = sync.Map{}

// RegisterUnoficialCurrency registers a currency code that is not a valid
// ISO 4217 currency code.
//
// This can be used for crypto currency codes, such as ETH, DAI, USDC, ...
func RegisterUnoficialCurrency(code string) {
	code = strings.TrimSpace(strings.ToUpper(code))

	if _, err := ParseCurrency(code); err == ErrInvalidCurrency {
		unoficialCurrencies.Store(code, true)
	}
}
