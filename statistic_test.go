package gobacktest

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestTrackEvent(t *testing.T) {
	var testCases = []struct {
		msg     string
		stat    Statistic
		event   EventHandler
		expStat Statistic
	}{
		{"testing simple event",
			Statistic{},
			&Bar{Close: 10},
			Statistic{
				eventHistory: []EventHandler{
					&Bar{Close: 10},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.stat.TrackEvent(tc.event)
		if !reflect.DeepEqual(tc.stat, tc.expStat) {
			t.Errorf("%v TrackEvent(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expStat, tc.stat)
		}
	}
}

func TestEvents(t *testing.T) {
	var testCases = []struct {
		msg       string
		stat      Statistic
		expEvents []EventHandler
	}{
		{"testing single event",
			Statistic{
				eventHistory: []EventHandler{
					&Bar{Close: 10},
				},
			},
			[]EventHandler{
				&Bar{Close: 10},
			},
		},
		{"testing multiple events",
			Statistic{
				eventHistory: []EventHandler{
					&Bar{Close: 10},
					&Bar{Close: 11},
					&Bar{Close: 9},
				},
			},
			[]EventHandler{
				&Bar{Close: 10},
				&Bar{Close: 11},
				&Bar{Close: 9},
			},
		},
		{"testing nil events",
			Statistic{},
			nil,
		},
	}

	for _, tc := range testCases {
		events := tc.stat.Events()
		if !reflect.DeepEqual(events, tc.expEvents) {
			t.Errorf("%v Events(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expEvents, events)
		}
	}
}

func TestTrackTransaction(t *testing.T) {
	var testCases = []struct {
		msg     string
		stat    Statistic
		fill    FillEvent
		expStat Statistic
	}{
		{"testing simple fill",
			Statistic{},
			&Fill{direction: BOT, qty: 100},
			Statistic{
				transactionHistory: []FillEvent{
					&Fill{direction: BOT, qty: 100},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc.stat.TrackTransaction(tc.fill)
		if !reflect.DeepEqual(tc.stat, tc.expStat) {
			t.Errorf("%v TrackTransaction(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expStat, tc.stat)
		}
	}
}

func TestTransactions(t *testing.T) {
	var testCases = []struct {
		msg             string
		stat            Statistic
		expTransactions []FillEvent
	}{
		{"testing single fill",
			Statistic{
				transactionHistory: []FillEvent{
					&Fill{direction: BOT, qty: 100},
				},
			},
			[]FillEvent{
				&Fill{direction: BOT, qty: 100},
			},
		},
		{"testing multiple fill events",
			Statistic{
				transactionHistory: []FillEvent{
					&Fill{direction: BOT, qty: 100},
					&Fill{direction: SLD, qty: 100},
					&Fill{direction: BOT, qty: 50},
				},
			},
			[]FillEvent{
				&Fill{direction: BOT, qty: 100},
				&Fill{direction: SLD, qty: 100},
				&Fill{direction: BOT, qty: 50},
			},
		},
		{"testing nil fill Events",
			Statistic{},
			nil,
		},
	}

	for _, tc := range testCases {
		transactions := tc.stat.Transactions()
		if !reflect.DeepEqual(transactions, tc.expTransactions) {
			t.Errorf("%v Transactions(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expTransactions, transactions)
		}
	}
}

func TestResetStatistic(t *testing.T) {
	var testCases = []struct {
		msg     string
		stat    Statistic
		expStat Statistic
	}{
		{"testing full statistic",
			Statistic{
				eventHistory: []EventHandler{
					&Bar{Close: 10},
					&Bar{Close: 15},
					&Signal{direction: BOT},
					&Order{direction: BOT},
					&Fill{qty: 100},
				},
				transactionHistory: []FillEvent{
					&Fill{qty: 100},
				},
				equity: []EquityPoint{
					{Equity: 100},
					{Equity: 90},
				},
				high: EquityPoint{Equity: 100},
				low:  EquityPoint{Equity: 90},
			},
			Statistic{},
		},
		{"testing empty statistic",
			Statistic{
				eventHistory:       []EventHandler{},
				transactionHistory: []FillEvent{},
				equity:             []EquityPoint{},
				high:               EquityPoint{},
				low:                EquityPoint{},
			},
			Statistic{},
		},
	}

	for _, tc := range testCases {
		tc.stat.Reset()
		if !reflect.DeepEqual(tc.stat, tc.expStat) {
			t.Errorf("%v Reset(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expStat, tc.stat)
		}
	}
}

func TestTotalEquityReturn(t *testing.T) {
	// set up test cases for getting first equity point
	var testCases = []struct {
		msg       string
		stat      Statistic
		expReturn float64
		expErr    error
	}{
		{"testing for multiple entryPoints",
			Statistic{
				equity: []EquityPoint{
					{Equity: 100, EquityReturn: 0},
					{Equity: 120, EquityReturn: 0.2},
				},
			},
			0.2,
			nil},
		{"testing for multiple entryPoints with same value",
			Statistic{
				equity: []EquityPoint{
					{Equity: 100, EquityReturn: 0},
					{Equity: 100, EquityReturn: 0},
				},
			},
			0,
			nil},
		{"testing for last entryPoints with 0 equity",
			Statistic{
				equity: []EquityPoint{
					{Equity: 100, EquityReturn: 0.1},
					{Equity: 0, EquityReturn: -1},
				},
			},
			-1,
			nil},
		{"testing for nil entryPoints",
			Statistic{},
			0,
			errors.New("could not calculate totalEquityReturn, no equity points found")},
	}

	for _, tc := range testCases {
		r, err := tc.stat.TotalEquityReturn()
		if (r != tc.expReturn) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v TotalEquityReturn(): \nexpected %#v %#v, \nactual   %#v %#v",
				tc.msg, tc.expReturn, tc.expErr, r, err)
		}
	}
}
func TestGetEquityPoint(t *testing.T) {
	var statCases = map[string]Statistic{
		"multiple": {
			equity: []EquityPoint{
				{Equity: 100, EquityReturn: 0.1},
				{Equity: 110, EquityReturn: 0.2},
				{Equity: 120, EquityReturn: 0.3},
			},
		},
		"single": {
			equity: []EquityPoint{
				{Equity: 150, EquityReturn: 0.25},
			},
		},
		"empty": {
			equity: []EquityPoint{},
		},
	}

	// define test cases struct
	type testCase struct {
		msg   string
		stat  Statistic
		expEP EquityPoint
		expOk bool
	}

	// set up test cases for getting first equity point
	var testCasesFirst = []testCase{
		{"testing first for multiple entryPoints",
			statCases["multiple"],
			EquityPoint{Equity: 100, EquityReturn: 0.1},
			true},
		{"testing first for single entryPoints",
			statCases["single"],
			EquityPoint{Equity: 150, EquityReturn: 0.25},
			true},
		{"testing first for nil entryPoints",
			statCases["empty"],
			EquityPoint{},
			false},
	}

	for _, tc := range testCasesFirst {
		ep, ok := tc.stat.firstEquityPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (ok != tc.expOk) {
			t.Errorf("%v firstEquityPoint(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expEP, tc.expOk, ep, ok)
		}
	}

	// set up test cases for getting last equity point
	var testCasesLast = []testCase{
		{"testing last for multiple entryPoints",
			statCases["multiple"],
			EquityPoint{Equity: 120, EquityReturn: 0.3},
			true},
		{"testing last for single entryPoints",
			statCases["single"],
			EquityPoint{Equity: 150, EquityReturn: 0.25},
			true},
		{"testing last for nil entryPoints",
			statCases["empty"],
			EquityPoint{},
			false},
	}

	for _, tc := range testCasesLast {
		ep, ok := tc.stat.lastEquityPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (ok != tc.expOk) {
			t.Errorf("%v firstEquityPoint(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expEP, tc.expOk, ep, ok)
		}
	}
}

func TestCalcEquityReturn(t *testing.T) {
	var testCases = []struct {
		msg   string
		stat  Statistic
		ep    EquityPoint
		expEP EquityPoint
	}{
		{"testing equity return with single equity points",
			Statistic{
				equity: []EquityPoint{
					{Equity: 100},
				},
			},
			EquityPoint{Equity: 90},
			EquityPoint{
				Equity:       90,
				EquityReturn: -0.1,
			},
		},
		{"testing equity return with multiple equity points",
			Statistic{
				equity: []EquityPoint{
					{Equity: 100},
					{Equity: 90},
					{Equity: 110},
				},
			},
			EquityPoint{Equity: 100},
			EquityPoint{
				Equity:       100,
				EquityReturn: -0.0909,
			},
		},
		{"testing equity return with single equity points but 0 equity",
			Statistic{
				equity: []EquityPoint{
					{Equity: 0},
				},
			},
			EquityPoint{Equity: 100},
			EquityPoint{
				Equity:       100,
				EquityReturn: 1,
			},
		},
		{"testing equity return with nil equity points",
			Statistic{
				equity: []EquityPoint{},
			},
			EquityPoint{Equity: 100},
			EquityPoint{
				Equity:   100,
				Drawdown: 0,
			},
		},
	}

	for _, tc := range testCases {
		ep := tc.stat.calcEquityReturn(tc.ep)
		if !reflect.DeepEqual(ep, tc.expEP) {
			t.Errorf("%v calcEquityReturn(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.ep, tc.expEP, ep)
		}
	}

}

func TestCalcDrawdown(t *testing.T) {
	var testCases = []struct {
		msg   string
		stat  Statistic
		ep    EquityPoint
		expEP EquityPoint
	}{
		{"testing drawdown with simple high EquityPoint",
			Statistic{
				high: EquityPoint{Equity: 100},
			},
			EquityPoint{Equity: 90},
			EquityPoint{
				Equity:   90,
				Drawdown: -0.1,
			},
		},
		{"testing drawdown with simple high EquityPoint equal equity",
			Statistic{
				high: EquityPoint{Equity: 100},
			},
			EquityPoint{Equity: 100},
			EquityPoint{
				Equity:   100,
				Drawdown: 0,
			},
		},
		{"testing drawdown with simple high EquityPoint lower equity",
			Statistic{
				high: EquityPoint{Equity: 90},
			},
			EquityPoint{Equity: 100},
			EquityPoint{
				Equity:   100,
				Drawdown: 0,
			},
		},
		{"testing drawdown with empty high EquityPoint",
			Statistic{},
			EquityPoint{Equity: 100},
			EquityPoint{
				Equity:   100,
				Drawdown: 0,
			},
		},
	}

	for _, tc := range testCases {
		ep := tc.stat.calcDrawdown(tc.ep)
		if !reflect.DeepEqual(ep, tc.expEP) {
			t.Errorf("%v calcDrawdown(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.ep, tc.expEP, ep)
		}
	}

}

func TestMaxDrawdown(t *testing.T) {
	var time1, _ = time.Parse("2006-01-02", "2017-09-25")
	var time2, _ = time.Parse("2006-01-02", "2017-09-26")
	var time3, _ = time.Parse("2006-01-02", "2017-09-27")
	var time4, _ = time.Parse("2006-01-02", "2017-09-28")
	var time5, _ = time.Parse("2006-01-02", "2017-09-29")

	// set up test cases for getting the max drawdown point
	var testCases = []struct {
		msg     string
		stat    Statistic
		expEP   EquityPoint
		expInt  int
		expMax  float64
		expTime time.Time
		expDur  float64 // duration in hours
	}{
		{"testing maxdrawdown for multiple entryPoints",
			Statistic{
				equity: []EquityPoint{
					{Timestamp: time1, Equity: 100, Drawdown: 0},
					{Timestamp: time2, Equity: 110, Drawdown: 0},
					{Timestamp: time3, Equity: 105, Drawdown: -0.0455},
					{Timestamp: time4, Equity: 95, Drawdown: -0.1364},
					{Timestamp: time5, Drawdown: 0},
				},
			},
			EquityPoint{Timestamp: time4, Equity: 95, Drawdown: -0.1364},
			3,
			-0.1364,
			time4,
			48,
		},
		{"testing maxdrawdown for single entryPoints",
			Statistic{
				equity: []EquityPoint{
					{Timestamp: time1, Equity: 100, Drawdown: 0},
				},
			},
			EquityPoint{Timestamp: time1, Equity: 100, Drawdown: 0},
			0,
			0,
			time1,
			0,
		},
		{"testing maxdrawdown for nil entryPoints",
			Statistic{},
			EquityPoint{},
			0,
			0,
			time.Time{},
			0,
		},
	}

	// testing for max drawdown equity point
	for _, tc := range testCases {
		i, ep := tc.stat.maxDrawdownPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (i != tc.expInt) {
			t.Errorf("%v maxDrawdownPoint(): \nexpected %d %#v, \nactual   %d %#v",
				tc.msg, tc.expInt, tc.expEP, i, ep)
		}
	}

	// testing for max drawdown value
	for _, tc := range testCases {
		max := tc.stat.MaxDrawdown()
		if !reflect.DeepEqual(max, tc.expMax) {
			t.Errorf("%v MaxDrawdown(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expMax, max)
		}
	}

	// testing for max drawdown time
	for _, tc := range testCases {
		time := tc.stat.MaxDrawdownTime()
		if !reflect.DeepEqual(time, tc.expTime) {
			t.Errorf("%v MaxDrawdownTime(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expTime, time)
		}
	}

	// testing for max drawdown duration
	for _, tc := range testCases {
		duration := tc.stat.MaxDrawdownDuration()
		if !reflect.DeepEqual(duration.Hours(), tc.expDur) {
			t.Errorf("%v MaxDrawdownDuration(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expDur, duration.Hours())
		}
	}
}

func TestSharpRatio(t *testing.T) {
	var testCases = []struct {
		msg      string
		stat     Statistic
		riskfree float64
		expSharp float64
	}{
		{"testing simple positiv sharp ratio",
			Statistic{
				equity: []EquityPoint{
					{EquityReturn: 1},
					{EquityReturn: 2},
					{EquityReturn: 3},
				},
			},
			0,
			2},
		{"testing simple zero sharp ratio",
			Statistic{
				equity: []EquityPoint{
					{EquityReturn: -1},
					{EquityReturn: 0},
					{EquityReturn: 1},
				},
			},
			0,
			0},
	}

	for _, tc := range testCases {
		sharp := tc.stat.SharpRatio(tc.riskfree)
		if sharp != tc.expSharp {
			t.Errorf("%v SharpRatio(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.riskfree, tc.expSharp, sharp)
		}
	}
}

func TestSortinoRatio(t *testing.T) {
	var testCases = []struct {
		msg        string
		stat       Statistic
		riskfree   float64
		expSortino float64
	}{
		{"testing simple sortino ratio",
			Statistic{
				equity: []EquityPoint{
					{EquityReturn: -3},
					{EquityReturn: -2},
					{EquityReturn: -1},
					{EquityReturn: 0},
					{EquityReturn: 1},
				},
			},
			0,
			-1},
	}

	for _, tc := range testCases {
		sortino := tc.stat.SortinoRatio(tc.riskfree)
		if sortino != tc.expSortino {
			t.Errorf("%v SortinoRatio(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.riskfree, tc.expSortino, sortino)
		}
	}
}
