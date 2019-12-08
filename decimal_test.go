package money_test

import (
	"testing"

	"github.com/deixis/money"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type decPair struct {
	X string
	Y string
}

func TestParseDecimal(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
		err    error
	}{
		{input: "120.0", expect: 120},
		{input: "120.00", expect: 120},
		{input: "120.000", expect: 120},
		{input: "120.12", expect: 120.12},
		{input: "120.125", expect: 120.125},
		{input: "120.123456789", expect: 120.123456789},
		{input: "0.00000001", expect: 0.00000001},           // 1 satoshi
		{input: "17950000000000.0", expect: 17950000000000}, // GDP USA
		{input: "3.141592653589793", expect: 3.141592653589793},
		{input: ".1111111111111110", expect: 0.1111111111111110},
		{input: ".1111111111111111", expect: 0.1111111111111111},
		{input: ".1111111111111119", expect: 0.1111111111111119},
		{input: "0.001", expect: 0.001},
		{input: "0.002", expect: 0.002},
		{input: "0.003", expect: 0.003},
		{input: "0.004", expect: 0.004},
		{input: "0.005", expect: 0.005},
		{input: "0.006", expect: 0.006},
		{input: "0.007", expect: 0.007},
		{input: "0.008", expect: 0.008},
		{input: "0.009", expect: 0.009},
		{input: "0.00", expect: 0},
		{input: "-0.00", expect: -0},
		{input: "-1.00", expect: -1.0},
		{input: "0", expect: 0},
		{input: "120", expect: 120},
		{input: "yyy", err: money.ErrInvalidDecimal},
		{input: "yyy.yyy", err: money.ErrInvalidDecimal},
		{input: "0x1.fffffffffffffp1023", err: money.ErrInvalidDecimal},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			if test.err != err {
				t.Errorf("#%d - expect error %s, but got %s - %s", i, test.err, err, test.input)
			}
			continue
		}
		if test.expect != dec.Float64() {
			t.Errorf("#%d - expect %f, but got %f - %s", i, test.expect, dec.Float64(), test.input)
		}
	}
}

func TestNewDecimal(t *testing.T) {
	t.Parallel()

	table := []struct {
		input float64
	}{
		{input: 120.0},
		{input: 120.12},
		{input: 120.125},
		{input: 0.00000001},
		{input: 17950000000000.0},
		{input: 0.0},
		{input: -1.0},
	}

	for i, test := range table {
		dec, err := money.NewDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.Float64()

		if test.input != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.input, res)
		}
	}
}

func TestMinDecimal(t *testing.T) {
	t.Parallel()

	table := []struct {
		input decPair
	}{
		{input: decPair{X: "1.0", Y: "1.0"}},
		{input: decPair{X: "-1.0", Y: "1.0"}},
		{input: decPair{X: "-1.0", Y: "-1.0"}},
		{input: decPair{X: "0.0", Y: "0.0"}},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		z := money.MinDecimal(x, y)
		expect := x.Float64()
		res := z.Float64()
		if expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, expect, res)
		}
	}
}

func TestMaxDecimal(t *testing.T) {
	t.Parallel()

	table := []struct {
		input decPair
	}{
		{input: decPair{X: "1.0", Y: "1.0"}},
		{input: decPair{X: "-1.0", Y: "1.0"}},
		{input: decPair{X: "-1.0", Y: "-1.0"}},
		{input: decPair{X: "0.0", Y: "0.0"}},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		z := money.MaxDecimal(x, y)
		expect := y.Float64()
		res := z.Float64()
		if expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, expect, res)
		}
	}
}

func TestDecimal_String(t *testing.T) {
	t.Parallel()

	table := []struct {
		input string
	}{
		{input: "120.0"},
		{input: "120.00"},
		{input: "120.000"},
		{input: "0.00"},
		{input: "-1.00"},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.String()

		if test.input != res {
			t.Errorf("#%d - expect `%s`, but got `%s`", i, test.input, res)
		}
	}

	d := money.Decimal{}
	got := d.String()
	expect := "0.0"
	if expect != got {
		t.Errorf("expect `%s`, but got `%s`", expect, got)
	}
}

func TestDecimal_Abs(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
	}{
		{input: "1.0", expect: 1.0},
		{input: "0.5", expect: 0.5},
		{input: "0.", expect: 0.0},
		{input: "-0.0", expect: 0.0},
		{input: "-0.5", expect: 0.5},
		{input: "-1.0", expect: 1.0},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.Abs().Float64()

		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Add(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  decPair
		expect float64
	}{
		{input: decPair{X: "1.0", Y: "1.0"}, expect: 2.0},
		{input: decPair{X: "1.0", Y: "-1.0"}, expect: 0.0},
		{input: decPair{X: "-1.0", Y: "1.0"}, expect: 0.0},
		{input: decPair{X: "-1.0", Y: "-1.0"}, expect: -2.0},
		{input: decPair{X: "1.0", Y: "0.0001"}, expect: 1.0001},
		{input: decPair{X: "2454495034.0", Y: "3451204593.0"}, expect: 5905699627},
		{input: decPair{X: "24544.95034", Y: "0.3451204593"}, expect: 24545.2954604593},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Add(y).Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Sub(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  decPair
		expect float64
	}{
		{input: decPair{X: "1.0", Y: "1.0"}, expect: 0.0},
		{input: decPair{X: "1.0", Y: "-1.0"}, expect: 2.0},
		{input: decPair{X: "-1.0", Y: "1.0"}, expect: -2.0},
		{input: decPair{X: "-1.0", Y: "-1.0"}, expect: 0.0},
		{input: decPair{X: "1.0", Y: "0.0001"}, expect: 0.9999},
		{input: decPair{X: "2454495034.0", Y: "3451204593.0"}, expect: -996709559},
		{input: decPair{X: "24544.95034", Y: "0.3451204593"}, expect: 24544.6052195407},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Sub(y).Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Mul(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  decPair
		expect float64
	}{
		{input: decPair{X: "1.0", Y: "1.0"}, expect: 1.0},
		{input: decPair{X: "1.0", Y: "-1.0"}, expect: -1.0},
		{input: decPair{X: "-1.0", Y: "1.0"}, expect: -1.0},
		{input: decPair{X: "-1.0", Y: "-1.0"}, expect: 1.0},
		{input: decPair{X: "1.0", Y: "0.0001"}, expect: 0.0001},
		{input: decPair{X: "2454495034.0", Y: "3451204593.0"}, expect: 8470964534836491162},
		{input: decPair{X: "24544.95034", Y: "0.3451204593"}, expect: 8470.964534836491162},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Mul(y).Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Div(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  decPair
		expect float64
	}{
		{input: decPair{X: "1.0", Y: "1.0"}, expect: 1.0},
		{input: decPair{X: "1.0", Y: "-1.0"}, expect: -1.0},
		{input: decPair{X: "-1.0", Y: "1.0"}, expect: -1.0},
		{input: decPair{X: "-1.0", Y: "-1.0"}, expect: 1.0},
		{input: decPair{X: "100.0", Y: "1.08"}, expect: 92.5925925925925926},
		{input: decPair{X: "1.0", Y: "0.0001"}, expect: 10000},
		{input: decPair{X: "1023427554493.0", Y: "43432632.0"}, expect: 23563.5628642767953828}, // rounded
		{input: decPair{X: "10234274355545544493.0", Y: "-3.0"}, expect: -3411424785181848164.3333333333333333},
		{input: decPair{X: "-4612301402398.4753343454", Y: "23.5"}, expect: -196268144782.9138440146978723},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Div(y).Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Neg(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
	}{
		{input: "1.0", expect: -1.0},
		{input: "0.5", expect: -0.5},
		{input: "0.", expect: 0.0},
		{input: "-0.0", expect: 0.0},
		{input: "-0.5", expect: 0.5},
		{input: "-1.0", expect: 1.0},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Neg().Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Mod(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  decPair
		expect float64
	}{
		{input: decPair{X: "10.0", Y: "3.0"}, expect: 1.0},
		{input: decPair{X: "-10.0", Y: "3.0"}, expect: -1.0},
		{input: decPair{X: "10.0", Y: "-3.0"}, expect: 1.0},
		{input: decPair{X: "-10.0", Y: "-3.0"}, expect: -1.0},
		{input: decPair{X: "0.1", Y: "0.1"}, expect: 0.0},
		{input: decPair{X: "0.0", Y: "2.0"}, expect: 0.0},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Mod(y).Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Round(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
		prec   int32
	}{
		{input: "120.0", prec: 1, expect: 120},
		{input: "120.00", prec: 2, expect: 120},
		{input: "120.000", prec: 3, expect: 120},
		{input: "0.00000001", prec: 8, expect: 0.00000001},
		{input: "17950000000000.0", prec: 1, expect: 17950000000000},
		{input: "3.141592653589793", prec: 2, expect: 3.14},
		{input: "0.11115", prec: 2, expect: 0.11},
		{input: "0.11115", prec: 3, expect: 0.111},
		{input: "0.11115", prec: 4, expect: 0.1112},
		{input: "0.101", prec: 2, expect: 0.10},
		{input: "0.102", prec: 2, expect: 0.10},
		{input: "0.103", prec: 2, expect: 0.10},
		{input: "0.104", prec: 2, expect: 0.10},
		{input: "0.105", prec: 2, expect: 0.11},
		{input: "0.106", prec: 2, expect: 0.11},
		{input: "0.107", prec: 2, expect: 0.11},
		{input: "0.108", prec: 2, expect: 0.11},
		{input: "0.109", prec: 2, expect: 0.11},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.Round(test.prec).Float64()

		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_RoundUp(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
		prec   int32
	}{
		{input: "120.0", prec: 1, expect: 120},
		{input: "120.00", prec: 2, expect: 120},
		{input: "120.000", prec: 3, expect: 120},
		{input: "0.00000001", prec: 8, expect: 0.00000001},
		{input: "17950000000000.0", prec: 1, expect: 17950000000000},
		{input: "3.141592653589793", prec: 2, expect: 3.15},
		{input: "0.11115", prec: 2, expect: 0.12},
		{input: "0.11115", prec: 3, expect: 0.112},
		{input: "0.11115", prec: 4, expect: 0.1112},
		{input: "0.101", prec: 2, expect: 0.11},
		{input: "0.102", prec: 2, expect: 0.11},
		{input: "0.103", prec: 2, expect: 0.11},
		{input: "0.104", prec: 2, expect: 0.11},
		{input: "0.105", prec: 2, expect: 0.11},
		{input: "0.106", prec: 2, expect: 0.11},
		{input: "0.107", prec: 2, expect: 0.11},
		{input: "0.108", prec: 2, expect: 0.11},
		{input: "0.109", prec: 2, expect: 0.11},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.RoundUp(test.prec).Float64()

		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_RoundDown(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
		prec   int32
	}{
		{input: "120.0", prec: 1, expect: 120},
		{input: "120.00", prec: 2, expect: 120},
		{input: "120.000", prec: 3, expect: 120},
		{input: "0.00000001", prec: 8, expect: 0.00000001},
		{input: "17950000000000.0", prec: 1, expect: 17950000000000},
		{input: "3.141592653589793", prec: 2, expect: 3.14},
		{input: "0.11115", prec: 2, expect: 0.11},
		{input: "0.11115", prec: 3, expect: 0.111},
		{input: "0.11115", prec: 4, expect: 0.1111},
		{input: "0.101", prec: 2, expect: 0.10},
		{input: "0.102", prec: 2, expect: 0.10},
		{input: "0.103", prec: 2, expect: 0.10},
		{input: "0.104", prec: 2, expect: 0.10},
		{input: "0.105", prec: 2, expect: 0.10},
		{input: "0.106", prec: 2, expect: 0.10},
		{input: "0.107", prec: 2, expect: 0.10},
		{input: "0.108", prec: 2, expect: 0.10},
		{input: "0.109", prec: 2, expect: 0.10},
		{input: "1.88", prec: 0, expect: 1.00},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.RoundDown(test.prec).Float64()

		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_RoundNearest(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
		unit   float64
	}{
		{input: "120.0", unit: 0.05, expect: 120},
		{input: "120.00", unit: 0.05, expect: 120},
		{input: "120.000", unit: 0.05, expect: 120},
		{input: "0.00000001", unit: 0.05, expect: 0.00},
		{input: "17950000000000.0", unit: 0.05, expect: 17950000000000},
		{input: "3.141592653589793", unit: 0.05, expect: 3.15},
		{input: "0.11112", unit: 0.0001, expect: 0.1111},
		{input: "0.11113", unit: 0.0001, expect: 0.1111},
		{input: "0.11115", unit: 0.0001, expect: 0.1112},
		{input: "0.10", unit: 0.05, expect: 0.10},
		{input: "0.11", unit: 0.05, expect: 0.10},
		{input: "0.12", unit: 0.05, expect: 0.10},
		{input: "0.13", unit: 0.05, expect: 0.15},
		{input: "0.14", unit: 0.05, expect: 0.15},
		{input: "0.15", unit: 0.05, expect: 0.15},
		{input: "0.16", unit: 0.05, expect: 0.15},
		{input: "0.17", unit: 0.05, expect: 0.15},
		{input: "0.18", unit: 0.05, expect: 0.20},
		{input: "0.19", unit: 0.05, expect: 0.20},
		{input: "-0.10", unit: 0.05, expect: -0.10},
		{input: "-0.11", unit: 0.05, expect: -0.10},
		{input: "-0.12", unit: 0.05, expect: -0.10},
		{input: "-0.13", unit: 0.05, expect: -0.15},
		{input: "-0.14", unit: 0.05, expect: -0.15},
		{input: "-0.15", unit: 0.05, expect: -0.15},
		{input: "-0.16", unit: 0.05, expect: -0.15},
		{input: "-0.17", unit: 0.05, expect: -0.15},
		{input: "-0.18", unit: 0.05, expect: -0.20},
		{input: "-0.19", unit: 0.05, expect: -0.20},
		{input: "0.10", unit: 0.01, expect: 0.10},
		{input: "0.11", unit: 0.01, expect: 0.11},
		{input: "0.12", unit: 0.01, expect: 0.12},
		{input: "0.13", unit: 0.01, expect: 0.13},
		{input: "0.14", unit: 0.01, expect: 0.14},
		{input: "0.15", unit: 0.01, expect: 0.15},
		{input: "0.16", unit: 0.01, expect: 0.16},
		{input: "0.17", unit: 0.01, expect: 0.17},
		{input: "0.18", unit: 0.01, expect: 0.18},
		{input: "0.19", unit: 0.01, expect: 0.19},
		{input: "1.75", unit: 1.00, expect: 2.00},
		{input: "1.5", unit: 1.00, expect: 2.00},
		{input: "1.49", unit: 1.00, expect: 1.00},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		unit, err := money.NewDecimal(test.unit)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.RoundNearest(unit).Float64()

		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f - %s", i, test.expect, res, test.input)
		}
	}
}

func TestDecimal_Truncate(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
		prec   int32
	}{
		{input: "120.0", prec: 1, expect: 120},
		{input: "120.00", prec: 2, expect: 120},
		{input: "120.000", prec: 3, expect: 120},
		{input: "0.00000001", prec: 8, expect: 0.00000001},
		{input: "17950000000000.0", prec: 1, expect: 17950000000000},
		{input: "3.141592653589793", prec: 2, expect: 3.14},
		{input: "0.11115", prec: 2, expect: 0.11},
		{input: "0.11115", prec: 3, expect: 0.111},
		{input: "0.11115", prec: 4, expect: 0.1111},
		{input: "0.101", prec: 2, expect: 0.10},
		{input: "0.102", prec: 2, expect: 0.10},
		{input: "0.103", prec: 2, expect: 0.10},
		{input: "0.104", prec: 2, expect: 0.10},
		{input: "0.105", prec: 2, expect: 0.10},
		{input: "0.106", prec: 2, expect: 0.10},
		{input: "0.107", prec: 2, expect: 0.10},
		{input: "0.108", prec: 2, expect: 0.10},
		{input: "0.109", prec: 2, expect: 0.10},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.Truncate(test.prec).Float64()

		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Floor(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
	}{
		{input: "120.0", expect: 120},
		{input: "120.00", expect: 120},
		{input: "120.000", expect: 120},
		{input: "0.0", expect: 0.0},
		{input: "17950000000000.0", expect: 17950000000000},
		{input: "3.141592653589793", expect: 3.0},
		{input: "0.11115", expect: 0.0},
		{input: "0.11115", expect: 0.0},
		{input: "0.11115", expect: 0.0},
		{input: "1.1", expect: 1.0},
		{input: "1.2", expect: 1.0},
		{input: "1.3", expect: 1.0},
		{input: "1.4", expect: 1.0},
		{input: "1.5", expect: 1.0},
		{input: "1.6", expect: 1.0},
		{input: "1.7", expect: 1.0},
		{input: "1.8", expect: 1.0},
		{input: "1.9", expect: 1.0},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.Floor().Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Ceil(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect float64
	}{
		{input: "120.0", expect: 120},
		{input: "120.00", expect: 120},
		{input: "120.000", expect: 120},
		{input: "0.0", expect: 0.0},
		{input: "17950000000000.0", expect: 17950000000000},
		{input: "3.141592653589793", expect: 4.0},
		{input: "0.11115", expect: 1.0},
		{input: "0.11115", expect: 1.0},
		{input: "0.11115", expect: 1.0},
		{input: "1.1", expect: 2.0},
		{input: "1.2", expect: 2.0},
		{input: "1.3", expect: 2.0},
		{input: "1.4", expect: 2.0},
		{input: "1.5", expect: 2.0},
		{input: "1.6", expect: 2.0},
		{input: "1.7", expect: 2.0},
		{input: "1.8", expect: 2.0},
		{input: "1.9", expect: 2.0},
	}

	for i, test := range table {
		dec, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}
		res := dec.Ceil().Float64()
		if test.expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, test.expect, res)
		}
	}
}

func TestDecimal_Cmp(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  decPair
		expect int
	}{
		{input: decPair{X: "1.0", Y: "1.0"}, expect: 0},
		{input: decPair{X: "-1.0", Y: "1.0"}, expect: -1},
		{input: decPair{X: "-1.0", Y: "-1.0"}, expect: 0},
		{input: decPair{X: "0.0", Y: "0.0"}, expect: 0},
		{input: decPair{X: "0.0", Y: "0.0000"}, expect: 0},
		{input: decPair{X: "1.0", Y: "-1.0"}, expect: 1},
		{input: decPair{X: "1.11", Y: "1.112"}, expect: -1},
		{input: decPair{X: "1.112", Y: "1.11"}, expect: 1},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Cmp(y)
		if test.expect != res {
			t.Errorf("#%d - expect %d, but got %d", i, test.expect, res)
		}
	}
}

func TestDecimal_Equal(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  decPair
		expect bool
	}{
		{input: decPair{X: "1.0", Y: "1.0"}, expect: true},
		{input: decPair{X: "-1.0", Y: "1.0"}, expect: false},
		{input: decPair{X: "-1.0", Y: "-1.0"}, expect: true},
		{input: decPair{X: "0.0", Y: "0.0"}, expect: true},
		{input: decPair{X: "1.0", Y: "-1.0"}, expect: false},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input.X)
		if err != nil {
			t.Fatal(err)
		}
		y, err := money.ParseDecimal(test.input.Y)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Equal(y)
		if test.expect != res {
			t.Errorf("#%d - expect %t, but got %t", i, test.expect, res)
		}
	}
}

func TestDecimal_Formatter(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  money.Decimal
		lang   language.Tag
		expect string
	}{
		{
			input:  money.MustParseDecimal("1000000.0"),
			lang:   language.English,
			expect: "1,000,000.0",
		},
		{
			input:  money.MustParseDecimal("1000000.0"),
			lang:   language.French,
			expect: "1 000 000,0",
		},
		{
			input:  money.MustParseDecimal("1000000.0"),
			lang:   language.German,
			expect: "1.000.000,0",
		},
		{
			input:  money.MustParseDecimal("-20.033"),
			lang:   language.English,
			expect: "-20.033",
		},
		{
			input:  money.MustParseDecimal("-20.033"),
			lang:   language.French,
			expect: "-20,033",
		},
		{
			input:  money.MustParseDecimal("-20.033"),
			lang:   language.German,
			expect: "-20,033",
		},
		{
			input:  money.MustParseDecimal("100.0"),
			lang:   language.English,
			expect: "100.0",
		},
		{
			input:  money.MustParseDecimal("33.33"),
			lang:   language.English,
			expect: "33.33",
		},
		{
			input:  money.MustParseDecimal("33.33"),
			lang:   language.French,
			expect: "33,33",
		},
		{
			input:  money.MustParseDecimal("33.33"),
			lang:   language.German,
			expect: "33,33",
		},
		{
			input:  money.MustParseDecimal("1.0"),
			lang:   language.English,
			expect: "1.0",
		},
		{
			input:  money.MustParseDecimal("0.5"),
			lang:   language.English,
			expect: "0.5",
		},
		{
			input:  money.MustParseDecimal("7.70"),
			lang:   language.English,
			expect: "7.70",
		},
	}

	for i, test := range table {
		p := message.NewPrinter(test.lang)
		f := test.input.Formatter()
		res := p.Sprint(f)

		if test.expect != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res)
		}
	}
}

func TestDecimal_PercentFormatter(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  money.Decimal
		lang   language.Tag
		expect string
	}{
		{
			input:  money.MustParseDecimal("1000000.0"),
			lang:   language.English,
			expect: "1,000,000%",
		},
		{
			input:  money.MustParseDecimal("1000000.0"),
			lang:   language.French,
			expect: "1 000 000 %",
		},
		{
			input:  money.MustParseDecimal("1000000.0"),
			lang:   language.German,
			expect: "1.000.000 %",
		},
		{
			input:  money.MustParseDecimal("-20.033"),
			lang:   language.English,
			expect: "-20.033%",
		},
		{
			input:  money.MustParseDecimal("-20.033"),
			lang:   language.French,
			expect: "-20,033 %",
		},
		{
			input:  money.MustParseDecimal("-20.033"),
			lang:   language.German,
			expect: "-20,033 %",
		},
		{
			input:  money.MustParseDecimal("100.0"),
			lang:   language.English,
			expect: "100%",
		},
		{
			input:  money.MustParseDecimal("33.33"),
			lang:   language.English,
			expect: "33.33%",
		},
		{
			input:  money.MustParseDecimal("33.33"),
			lang:   language.French,
			expect: "33,33 %",
		},
		{
			input:  money.MustParseDecimal("33.33"),
			lang:   language.German,
			expect: "33,33 %",
		},
		{
			input:  money.MustParseDecimal("3.1415926535897"),
			lang:   language.German,
			expect: "3,1415926535897 %",
		},
		{
			input:  money.MustParseDecimal("1.0"),
			lang:   language.English,
			expect: "1%",
		},
		{
			input:  money.MustParseDecimal("0.5"),
			lang:   language.English,
			expect: "0.5%",
		},
		{
			input:  money.MustParseDecimal("7.70"),
			lang:   language.English,
			expect: "7.7%",
		},
	}

	for i, test := range table {
		p := message.NewPrinter(test.lang)
		f := test.input.PercentFormatter()
		res := p.Sprint(f)

		if test.expect != res {
			t.Errorf("#%d - expect %s, but got %s", i, test.expect, res)
		}
	}
}

func TestDecimal_IsZero(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect bool
	}{
		{input: "-1.0", expect: false},
		{input: "0.0", expect: true},
		{input: "0.0000", expect: true},
		{input: "1.0", expect: false},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}

		res := x.IsZero()
		if test.expect != res {
			t.Errorf("#%d - expect %t, but got %t", i, test.expect, res)
		}
	}
}

func TestDecimal_Sign(t *testing.T) {
	t.Parallel()

	table := []struct {
		input  string
		expect int
	}{
		{input: "1.0", expect: 1},
		{input: "-1.0", expect: -1},
		{input: "0.0", expect: 0},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}

		res := x.Sign()
		if test.expect != res {
			t.Errorf("#%d - expect %d, but got %d", i, test.expect, res)
		}
	}
}

func TestDecimal_JSON(t *testing.T) {
	t.Parallel()

	table := []struct {
		input string
	}{
		{input: "1.0"},
		{input: "-1.0"},
		{input: "0.0"},
		{input: "0.00000001"},
		{input: "17950000000000.0"},
		{input: "3.141592653589793"},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}

		data, err := x.MarshalJSON()
		if err != nil {
			t.Fatal("cannot marshal JSON", err)
		}

		y := money.Decimal{}
		if err := y.UnmarshalJSON(data); err != nil {
			t.Fatal("cannot unmarshal JSON", err)
		}

		expect := x.Float64()
		res := y.Float64()
		if expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, expect, res)
		}
	}
}

func TestDecimal_Gob(t *testing.T) {
	t.Parallel()

	table := []struct {
		input string
	}{
		{input: "1.0"},
		{input: "-1.0"},
		{input: "0.0"},
		{input: "0.00000001"},
		{input: "17950000000000.0"},
		{input: "3.141592653589793"},
	}

	for i, test := range table {
		x, err := money.ParseDecimal(test.input)
		if err != nil {
			t.Fatal(err)
		}

		data, err := x.GobEncode()
		if err != nil {
			t.Fatal("cannot gob encode", err)
		}

		y := money.Decimal{}
		if err := y.GobDecode(data); err != nil {
			t.Fatal("cannot gob decode", err)
		}

		expect := x.Float64()
		res := y.Float64()
		if expect != res {
			t.Errorf("#%d - expect %f, but got %f", i, expect, res)
		}
	}
}
