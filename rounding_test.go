package money_test

import (
	"testing"

	"github.com/deixis/money"
)

func TestMoney_Round(t *testing.T) {
	t.Parallel()

	table := []struct {
		input      *money.Money
		standard   *money.Money
		cash       *money.Money
		accounting *money.Money
	}{
		// CHF
		{
			input:      money.MustParse("120.0", "CHF"),
			standard:   money.MustParse("120.0", "CHF"),
			cash:       money.MustParse("120.0", "CHF"),
			accounting: money.MustParse("120.0", "CHF"),
		},
		{
			input:      money.MustParse("120.01", "CHF"),
			standard:   money.MustParse("120.01", "CHF"),
			cash:       money.MustParse("120.00", "CHF"),
			accounting: money.MustParse("120.01", "CHF"),
		},
		{
			input:      money.MustParse("120.03", "CHF"),
			standard:   money.MustParse("120.03", "CHF"),
			cash:       money.MustParse("120.05", "CHF"),
			accounting: money.MustParse("120.03", "CHF"),
		},
		{
			input:      money.MustParse("120.08", "CHF"),
			standard:   money.MustParse("120.08", "CHF"),
			cash:       money.MustParse("120.10", "CHF"),
			accounting: money.MustParse("120.08", "CHF"),
		},
		{
			input:      money.MustParse("120.001", "CHF"),
			standard:   money.MustParse("120.00", "CHF"),
			cash:       money.MustParse("120.00", "CHF"),
			accounting: money.MustParse("120.00", "CHF"),
		},
		{
			input:      money.MustParse("120.009", "CHF"),
			standard:   money.MustParse("120.01", "CHF"),
			cash:       money.MustParse("120.00", "CHF"),
			accounting: money.MustParse("120.01", "CHF"),
		},
		// EUR
		{
			input:      money.MustParse("120.0", "EUR"),
			standard:   money.MustParse("120.0", "EUR"),
			cash:       money.MustParse("120.0", "EUR"),
			accounting: money.MustParse("120.0", "EUR"),
		},
		{
			input:      money.MustParse("120.01", "EUR"),
			standard:   money.MustParse("120.01", "EUR"),
			cash:       money.MustParse("120.01", "EUR"),
			accounting: money.MustParse("120.01", "EUR"),
		},
		{
			input:      money.MustParse("120.03", "EUR"),
			standard:   money.MustParse("120.03", "EUR"),
			cash:       money.MustParse("120.03", "EUR"),
			accounting: money.MustParse("120.03", "EUR"),
		},
		{
			input:      money.MustParse("120.08", "EUR"),
			standard:   money.MustParse("120.08", "EUR"),
			cash:       money.MustParse("120.08", "EUR"),
			accounting: money.MustParse("120.08", "EUR"),
		},
		{
			input:      money.MustParse("120.001", "EUR"),
			standard:   money.MustParse("120.00", "EUR"),
			cash:       money.MustParse("120.00", "EUR"),
			accounting: money.MustParse("120.00", "EUR"),
		},
		{
			input:      money.MustParse("120.009", "EUR"),
			standard:   money.MustParse("120.01", "EUR"),
			cash:       money.MustParse("120.01", "EUR"),
			accounting: money.MustParse("120.01", "EUR"),
		},
		// JPY
		{
			input:      money.MustParse("120.0", "JPY"),
			standard:   money.MustParse("120.0", "JPY"),
			cash:       money.MustParse("120.0", "JPY"),
			accounting: money.MustParse("120.0", "JPY"),
		},
		{
			input:      money.MustParse("120.01", "JPY"),
			standard:   money.MustParse("120.0", "JPY"),
			cash:       money.MustParse("120.0", "JPY"),
			accounting: money.MustParse("120.0", "JPY"),
		},
		{
			input:      money.MustParse("120.03", "JPY"),
			standard:   money.MustParse("120.0", "JPY"),
			cash:       money.MustParse("120.0", "JPY"),
			accounting: money.MustParse("120.0", "JPY"),
		},
		{
			input:      money.MustParse("120.08", "JPY"),
			standard:   money.MustParse("120.0", "JPY"),
			cash:       money.MustParse("120.0", "JPY"),
			accounting: money.MustParse("120.0", "JPY"),
		},
		{
			input:      money.MustParse("120.001", "JPY"),
			standard:   money.MustParse("120.0", "JPY"),
			cash:       money.MustParse("120.0", "JPY"),
			accounting: money.MustParse("120.0", "JPY"),
		},
		{
			input:      money.MustParse("120.009", "JPY"),
			standard:   money.MustParse("120.0", "JPY"),
			cash:       money.MustParse("120.0", "JPY"),
			accounting: money.MustParse("120.0", "JPY"),
		},
		{
			input:      money.MustParse("120.75", "JPY"),
			standard:   money.MustParse("121.0", "JPY"),
			cash:       money.MustParse("121.0", "JPY"),
			accounting: money.MustParse("121.0", "JPY"),
		},
	}

	for i, test := range table {
		unit := test.input.Currency.RoundUnit(money.RoundingStandard)
		standard := money.Round(test.input.Amount, unit, money.RoundToNearest)
		if !test.standard.Amount.Equal(standard) {
			t.Errorf("#%d - expect rounding standard to return %s, but got %s",
				i, test.standard, standard,
			)
		}

		unit = test.input.Currency.RoundUnit(money.RoundingCash)
		cash := money.Round(test.input.Amount, unit, money.RoundToNearest)
		if !test.cash.Amount.Equal(cash) {
			t.Errorf("#%d - expect rounding cash to return %s, but got %s",
				i, test.cash, cash,
			)
		}

		unit = test.input.Currency.RoundUnit(money.RoundingAccounting)
		accounting := money.Round(test.input.Amount, unit, money.RoundToNearest)
		if !test.accounting.Amount.Equal(accounting) {
			t.Errorf("#%d - expect rounding accounting to return %s, but got %s",
				i, test.accounting, accounting,
			)
		}
	}
}
