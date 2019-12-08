package money_test

import (
	"testing"

	"github.com/deixis/money"
)

func TestParseCurrency(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect money.Currency
		err    error
	}{
		{input: "CHF", expect: "CHF"},
		{input: "  chf  ", expect: "CHF"},
		{input: "cHf ", expect: "CHF"},
		{input: "USD", expect: "USD"},
		{input: "XXX", expect: money.Currency("XXX")},
	}

	for i, test := range table {
		res, err := money.ParseCurrency(test.input)
		if err != nil {
			if test.err != err {
				t.Errorf("#%d - expect error %s, but got %s - %s", i, test.err, err, test.input)
			}
			continue
		}
		if test.expect != res {
			t.Errorf("#%d - expect %s, but got %s - %s", i, test.expect, res, test.input)
		}
	}
}

func TestCurency_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	table := []struct {
		input money.Currency
	}{
		{input: "CHF"},
		{input: "USD"},
		{input: "CNY"},
	}

	for i, test := range table {
		data, err := test.input.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		var res money.Currency
		if err := res.UnmarshalJSON(data); err != nil {
			t.Fatal(err)
		}

		if test.input != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.input, res)
		}
	}
}

func TestCurency_GobEncode(t *testing.T) {
	t.Parallel()

	table := []struct {
		input money.Currency
	}{
		{input: "CHF"},
		{input: "USD"},
		{input: "CNY"},
	}

	for i, test := range table {
		data, err := test.input.GobEncode()
		if err != nil {
			t.Fatal(err)
		}

		var res money.Currency
		if err := res.GobDecode(data); err != nil {
			t.Fatal(err)
		}

		if test.input != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.input, res)
		}
	}
}

func Test_UnoficialCurrency(t *testing.T) {
	t.Parallel()

	table := []struct {
		input money.Currency
	}{
		{input: "ETH"},
		{input: "USDC"},
		{input: "DAI"},
	}

	for i, test := range table {
		_, err := money.ParseCurrency(test.input.String())
		if err != money.ErrInvalidCurrency {
			t.Fatalf("#%d - expect unoficial currency to fail when not registered", i)
		}

		money.RegisterUnoficialCurrency(test.input.String())

		_, err = money.ParseCurrency(test.input.String())
		if err != nil {
			t.Errorf("#%d - expect unoficial currency to be valid when registered, but got %s", i, err)
		}
	}
}
