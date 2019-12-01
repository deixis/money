package money_test

import (
	"encoding/json"
	"testing"

	"github.com/deixis/money"
)

func TestMoney_Equal(t *testing.T) {
	t.Parallel()

	table := []struct {
		x      *money.Money
		y      *money.Money
		expect bool
	}{
		{
			x: money.MustParse("120.0", "CHF"),
			y: money.MustParse("120.0", "CHF"), expect: true},
		{
			x: money.MustParse("120.00", "CHF"),
			y: money.MustParse("120.0000", "CHF"), expect: true},
		{
			x: money.MustParse("120.00", "CHF"),
			y: money.MustParse("-120.0000", "CHF"), expect: false},
		{
			x: money.MustParse("0.0000", "CHF"),
			y: money.MustParse("0.0000", "CHF"), expect: true},
		{
			x: money.MustParse("-120.12", "CHF"),
			y: money.MustParse("-120.1234", "CHF"), expect: false},
	}

	for i, test := range table {
		res := test.x.Equal(test.y)
		if test.expect != res {
			t.Errorf("#%d - expect %t, but got %t", i, test.expect, res)
		}
	}
}

func TestMoney_Validate(t *testing.T) {
	t.Parallel()

	table := []struct {
		x      *money.Money
		expect error
	}{
		{x: money.MustParse("120.0", "CHF"), expect: nil},
		{x: money.MustParse("120.00", "CHF"), expect: nil},
		{x: money.MustParse("120.00", "CHF"), expect: nil},
		{x: money.MustParse("0.0000", "CHF"), expect: nil},
		{x: money.MustParse("-120.12", "CHF"), expect: nil},
		{x: &money.Money{}, expect: money.ErrInvalidCurrency},
	}

	for i, test := range table {
		res := test.x.Validate()
		if test.expect != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res)
		}
	}
}

func TestMoney_MarshalJSON(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  *money.Money
		expect string
	}{
		{
			input:  money.MustParse("120.0", "CHF"),
			expect: "{\"amount\":\"120.0\",\"currency\":\"CHF\"}"},
		{
			input:  money.MustParse("120.00", "CHF"),
			expect: "{\"amount\":\"120.00\",\"currency\":\"CHF\"}"},
		{
			input:  money.MustParse("120.0000", "CHF"),
			expect: "{\"amount\":\"120.0000\",\"currency\":\"CHF\"}"},
		{
			input:  money.MustParse("-120.00", "CHF"),
			expect: "{\"amount\":\"-120.00\",\"currency\":\"CHF\"}"},
		{
			input:  money.MustParse("0.00", "CHF"),
			expect: "{\"amount\":\"0.00\",\"currency\":\"CHF\"}"},
	}

	for i, test := range table {
		data, err := json.Marshal(test.input)
		if err != nil {
			t.Fatal(err)
		}

		if test.expect != string(data) {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, data)
		}
	}
}

func TestMoney_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	table := []struct {
		input *money.Money
	}{
		{input: money.MustParse("120.0", "CHF")},
		{input: money.MustParse("120.00", "CHF")},
		{input: money.MustParse("120.0000", "CHF")},
	}

	for i, test := range table {
		data, err := json.Marshal(test.input)
		if err != nil {
			t.Fatal(err)
		}

		res := &money.Money{}
		if err := json.Unmarshal(data, res); err != nil {
			t.Fatal(err)
		}

		if !test.input.Equal(res) {
			t.Errorf("#%d - expect %s, but got %s", i, test.input, res)
		}
	}
}
