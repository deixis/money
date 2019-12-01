package money

import (
	"golang.org/x/text/currency"
)

// RoundingMode defines the rounding Mode to apply
type RoundingMode string

const (
	// Down rounds down to the previous increment
	// e.g. decimal: 1.49 increment: 0.1 result: 1.4
	RoundDown RoundingMode = "down"
	// Up rounds up to the next increment
	// e.g. decimal: 1.41 increment: 0.1 result: 1.5
	RoundUp RoundingMode = "up"
	// ToNearest rounds to the nearest increment
	// e.g. decimal: 1.45 increment: 0.1 result: 1.5
	RoundToNearest RoundingMode = "to_nearest"
)

// RoundingKind defines a rounding standard for currencies
type RoundingKind string

var (
	// RoundingStandard defines standard rounding and formatting for currencies.
	RoundingStandard RoundingKind = "standard"

	// RoundingCash defines rounding and formatting standards for cash
	// transactions.
	RoundingCash RoundingKind = "cash"

	// RoundingAccounting defines rounding and formatting standards for
	// accounting.
	RoundingAccounting RoundingKind = "accounting"
)

func (k RoundingKind) kind() currency.Kind {
	switch k {
	case RoundingStandard:
		return currency.Standard
	case RoundingCash:
		return currency.Cash
	case RoundingAccounting:
		return currency.Accounting
	}
	return currency.Standard
}

// Round rounds the given amount from the given unit
func Round(x Decimal, unit Decimal, mode RoundingMode) Decimal {
	prec := unit.Exponent() * -1

	switch mode {
	case RoundDown:
		rounded := x.RoundDown(prec)
		return rounded.Sub(rounded.Mod(unit)).Truncate(prec)
	case RoundUp:
		rounded := x.RoundUp(prec)
		return rounded.Add(rounded.Mod(unit)).Truncate(prec)
	case RoundToNearest:
		return x.RoundNearest(unit).Truncate(prec)
	}
	return Decimal{}
}
