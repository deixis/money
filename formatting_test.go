package money_test

import (
	"testing"

	"github.com/deixis/money"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestMoney_Format(t *testing.T) {
	t.Parallel()

	iso := &money.Formatter{
		CurrencyFormater: money.FormatterISO,
		Rounding:         money.RoundingStandard,
	}
	symbol := &money.Formatter{
		CurrencyFormater: money.FormatterSymbol,
		Rounding:         money.RoundingStandard,
	}
	narrowSymbol := &money.Formatter{
		CurrencyFormater: money.FormatterNarrowSymbol,
		Rounding:         money.RoundingStandard,
	}

	table := []struct {
		input     *money.Money
		formatter *money.Formatter
		lang      language.Tag
		expect    string
	}{
		{
			input:     money.MustParse("120.0", "CHF"),
			formatter: iso,
			lang:      language.English,
			expect:    "CHF 120.00",
		},
		{
			input:     money.MustParse("120.05", "CHF"),
			formatter: iso,
			lang:      language.English,
			expect:    "CHF 120.05",
		},
		{
			input:     money.MustParse("120.01", "CHF"),
			formatter: symbol,
			lang:      language.English,
			expect:    "CHF 120.01",
		},
		{
			input:     money.MustParse("120.01", "CHF"),
			formatter: narrowSymbol,
			lang:      language.English,
			expect:    "CHF 120.01",
		},
		{
			input:     money.MustParse("120.01", "CHF"),
			formatter: iso,
			lang:      language.English,
			expect:    "CHF 120.01",
		},
		{
			input:     money.MustParse("120.001", "CHF"),
			formatter: iso,
			lang:      language.English,
			expect:    "CHF 120.00",
		},
		{
			input:     money.MustParse("10000.001", "CHF"),
			formatter: iso,
			lang:      language.English,
			expect:    "CHF 10000.00",
		},
		{
			input:     money.MustParse("10000.001", "CHF"),
			lang:      language.AmericanEnglish,
			formatter: iso,
			expect:    "CHF 10000.00",
		},
		{
			input:     money.MustParse("1000000.001", "CHF"),
			formatter: iso,
			lang:      language.BritishEnglish,
			expect:    "CHF 1000000.00",
		},
		{
			input:     money.MustParse("1000000.001", "CHF"),
			formatter: iso,
			lang:      language.French,
			expect:    "CHF 1000000.00",
		},
		{
			input:     money.MustParse("1000000.001", "CHF"),
			formatter: iso,
			lang:      language.German,
			expect:    "CHF 1000000.00",
		},
		{
			input:     money.MustParse("1000000.001", "CHF"),
			formatter: iso,
			lang:      language.Chinese,
			expect:    "CHF 1000000.00",
		},
		{
			input:     money.MustParse("1000000.001", "CHF"),
			formatter: symbol,
			lang:      language.Burmese,
			expect:    "CHF 1000000.00",
		},
		{
			input:     money.MustParse("1.01", "EUR"),
			formatter: symbol,
			lang:      language.English,
			expect:    "€ 1.01",
		},
		{
			input:     money.MustParse("1.0", "EUR"),
			formatter: symbol,
			lang:      language.English,
			expect:    "€ 1.00",
		},
		{
			input:     money.MustParse("1.0", "EUR"),
			formatter: narrowSymbol,
			lang:      language.English,
			expect:    "€ 1.00",
		},
		{
			input:     money.MustParse("1.001", "EUR"),
			formatter: iso,
			lang:      language.English,
			expect:    "EUR 1.00",
		},
		{
			input:     money.MustParse("-100.001", "EUR"),
			formatter: iso,
			lang:      language.English,
			expect:    "EUR -100.00",
		},
		{
			input:     money.MustParse("-100.001", "EUR"),
			formatter: iso,
			lang:      language.Chinese,
			expect:    "EUR -100.00",
		},
		{
			input:     money.MustParse("-100.009", "CNY"),
			formatter: iso,
			lang:      language.Chinese,
			expect:    "CNY -100.01",
		},
		{
			input:     money.MustParse("-100.009", "JPY"),
			formatter: iso,
			lang:      language.Chinese,
			expect:    "JPY -100",
		},
	}

	for i, test := range table {
		input := test.formatter.Wrap(test.input)

		p := message.NewPrinter(test.lang)
		res := p.Sprintf("%f", input)

		if test.expect != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res)
		}
	}
}
