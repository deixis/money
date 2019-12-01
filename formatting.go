package money

import (
	"fmt"

	"golang.org/x/text/currency"
)

// CurrencyFormatter decorates a given number with formatting options.
type CurrencyFormatter = currency.Formatter

var (
	// FormatterNarrowSymbol usess narrow symbols. Overrides Symbol, if present.
	FormatterNarrowSymbol = currency.NarrowSymbol

	// FormatterSymbol uses symbols instead of ISO codes, when available.
	FormatterSymbol = currency.Symbol

	// FormatterISO uses ISO code as symbol.
	FormatterISO = currency.ISO
)

// Formatter formats Money to its string representation
type Formatter struct {
	CurrencyFormater CurrencyFormatter
	Rounding         RoundingKind
}

// Wrap decorates x with the formating preferences
func (f *Formatter) Wrap(x *Money) fmt.Formatter {
	fn := f.CurrencyFormater.Default(
		*x.Currency.currency(),
	).Kind(
		currency.Kind(f.Rounding.kind()),
	)
	return fn(x.Amount.Float64())
}

// DecimalFormatter formats Decimal to its string representation
type DecimalFormatter struct {
	CurrencyFormater CurrencyFormatter
	Currency         Currency
	Rounding         RoundingKind
}

// Wrap decorates x with the formating preferences
func (f *DecimalFormatter) Wrap(x *Decimal) fmt.Formatter {
	fn := f.CurrencyFormater.Default(
		*f.Currency.currency(),
	).Kind(
		currency.Kind(f.Rounding.kind()),
	)
	return fn(x.Float64())
}
