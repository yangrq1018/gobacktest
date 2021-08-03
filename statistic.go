package gobacktest

import (
	"errors"
	"fmt"
	"math"
	"time"

	"gonum.org/v1/gonum/stat"
)

// StatisticHandler is a basic statistic interface
type StatisticHandler interface {
	EventTracker
	TransactionTracker
	StatisticPrinter
	Reseter
	StatisticUpdater
	Resulter
}

// EventTracker is responsible for all event tracking during a backtest
type EventTracker interface {
	TrackEvent(EventHandler)
	Events() []EventHandler
}

// TransactionTracker is responsible for all transaction tracking during a backtest
type TransactionTracker interface {
	TrackTransaction(FillEvent)
	Transactions() []FillEvent
}

// StatisticPrinter handles printing of the statistics to screen
type StatisticPrinter interface {
	PrintResult()
}

// StatisticUpdater handles the updateing of the statistics
type StatisticUpdater interface {
	Update(DataEvent, PortfolioHandler)
}

// Resulter bundles all methods which return the results of the backtest
type Resulter interface {
	TotalEquityReturn() (float64, error)
	MaxDrawdown() float64
	MaxDrawdownTime() time.Time
	MaxDrawdownDuration() time.Duration
	SharpRatio(float64) float64
	SortinoRatio(float64) float64
}

// Statistic is a basic test statistic, which holds simple lists of historic events
type Statistic struct {
	eventHistory       []EventHandler
	transactionHistory []FillEvent
	equity             []EquityPoint
	high               EquityPoint
	low                EquityPoint
}

type EquityPoint struct {
	Timestamp    time.Time
	Equity       float64
	EquityReturn float64
	Drawdown     float64
}

func (s *Statistic) Equity() []EquityPoint {
	return s.equity
}

// Update the complete statistics to a given data event.
func (s *Statistic) Update(d DataEvent, p PortfolioHandler) {
	// create new equity point based on current data timestamp and portfolio value
	e := EquityPoint{}
	e.Timestamp = d.Time()
	e.Equity = p.Value()

	// calc equity return for current equity point
	if len(s.equity) > 0 {
		e = s.calcEquityReturn(e)
	}

	// calc drawdown for current equity point
	if len(s.equity) > 0 {
		e = s.calcDrawdown(e)
	}

	// set high and low equity point
	if e.Equity >= s.high.Equity {
		s.high = e
	}
	if e.Equity <= s.low.Equity {
		s.low = e
	}

	// append new equity point
	s.equity = append(s.equity, e)
}

// TrackEvent tracks an event
func (s *Statistic) TrackEvent(e EventHandler) {
	s.eventHistory = append(s.eventHistory, e)
}

// Events returns the complete events history
func (s Statistic) Events() []EventHandler {
	return s.eventHistory
}

// TrackTransaction tracks a transaction aka a fill event
func (s *Statistic) TrackTransaction(f FillEvent) {
	s.transactionHistory = append(s.transactionHistory, f)
}

// Transactions returns the complete events history
func (s Statistic) Transactions() []FillEvent {
	return s.transactionHistory
}

// Reset the statistic to a clean state
func (s *Statistic) Reset() error {
	s.eventHistory = nil
	s.transactionHistory = nil
	s.equity = nil
	s.high = EquityPoint{}
	s.low = EquityPoint{}
	return nil
}

// PrintResult prints the backtest statistics to the screen
func (s Statistic) PrintResult() {
	fmt.Println("Printing backtest results:")
	fmt.Printf("Counted %d total events.\n", len(s.Events()))

	fmt.Printf("Counted %d total transactions:\n", len(s.Transactions()))
	for k, v := range s.Transactions() {
		fmt.Printf("%d. Transaction: %v Action: %v Price: %f Qty: %d\n", k+1, v.Time().Format("2006-01-02"), v.Direction(), v.Price(), v.Qty())
	}
}

// TotalEquityReturn calculates the the total return on the first and last equity point
func (s Statistic) TotalEquityReturn() (r float64, err error) {
	firstEquityPoint, ok := s.firstEquityPoint()
	if !ok {
		return r, errors.New("could not calculate totalEquityReturn, no equity points found")
	}
	firstEquity := firstEquityPoint.Equity

	lastEquityPoint, _ := s.lastEquityPoint()
	// if !ok {
	// 	return r, errors.New("could not calculate totalEquityReturn, no last equity point")
	// }
	lastEquity := lastEquityPoint.Equity

	totalEquityReturn := (lastEquity - firstEquity) / firstEquity
	total := math.Round(totalEquityReturn*math.Pow10(DP)) / math.Pow10(DP)
	return total, nil
}

// MaxDrawdown returns the maximum draw down value in percent.
func (s Statistic) MaxDrawdown() float64 {
	_, ep := s.maxDrawdownPoint()
	return ep.Drawdown
}

// MaxDrawdownTime returns the time of the maximum draw down value.
func (s Statistic) MaxDrawdownTime() time.Time {
	_, ep := s.maxDrawdownPoint()
	return ep.Timestamp
}

// MaxDrawdownDuration returns the maximum draw down value in percent
func (s Statistic) MaxDrawdownDuration() (d time.Duration) {
	i, ep := s.maxDrawdownPoint()

	if len(s.equity) == 0 {
		return d
	}

	// walk the equity slice up to find a higher equity point
	maxPoint := EquityPoint{}
	for index := i; index >= 0; index-- {
		if s.equity[index].Equity > maxPoint.Equity {
			maxPoint = s.equity[index]
		}
	}

	d = ep.Timestamp.Sub(maxPoint.Timestamp)
	return d
}

// SharpRatio returns the Sharp ratio compared to a risk free benchmark return.
func (s *Statistic) SharpRatio(riskfree float64) float64 {
	var equityReturns = make([]float64, len(s.equity))

	for i, v := range s.equity {
		equityReturns[i] = v.EquityReturn
	}
	mean, stddev := stat.MeanStdDev(equityReturns, nil)

	sharp := (mean - riskfree) / stddev
	return sharp
}

// SortinoRatio returns the Sortino ratio compared to a risk free benchmark return.
func (s *Statistic) SortinoRatio(riskfree float64) float64 {
	var equityReturns = make([]float64, len(s.equity))

	for i, v := range s.equity {
		equityReturns[i] = v.EquityReturn
	}
	mean := stat.Mean(equityReturns, nil)

	// sortino uses the stddev of only the negativ returns
	var negReturns = []float64{}
	for _, v := range equityReturns {
		if v < 0 {
			negReturns = append(negReturns, v)
		}
	}
	stdDev := stat.StdDev(negReturns, nil)

	sortino := (mean - riskfree) / stdDev
	return sortino
}

// returns the first EquityPoint
func (s Statistic) firstEquityPoint() (ep EquityPoint, ok bool) {
	if len(s.equity) <= 0 {
		return ep, false
	}
	ep = s.equity[0]

	return ep, true
}

// returns the last EquityPoint
func (s Statistic) lastEquityPoint() (ep EquityPoint, ok bool) {
	if len(s.equity) <= 0 {
		return ep, false
	}
	ep = s.equity[len(s.equity)-1]

	return ep, true
}

// calculates the equity return of an equity point relativ to the last equity point
func (s Statistic) calcEquityReturn(e EquityPoint) EquityPoint {
	last, ok := s.lastEquityPoint()
	// no equity point before the current
	if !ok {
		e.EquityReturn = 0
		return e
	}

	lastEquity := last.Equity
	currentEquity := e.Equity

	// last equity point has 0 equity
	if lastEquity == 0 {
		e.EquityReturn = 1
		return e
	}

	equityReturn := (currentEquity - lastEquity) / lastEquity
	e.EquityReturn = math.Round(equityReturn*math.Pow10(DP)) / math.Pow10(DP)

	return e
}

// calculates the drawdown of an equity point relativ to the latest high of the statistic handler
func (s Statistic) calcDrawdown(e EquityPoint) EquityPoint {
	if s.high.Equity == 0 {
		e.Drawdown = 0
		return e
	}

	lastHigh := s.high.Equity
	equity := e.Equity

	if equity >= lastHigh {
		e.Drawdown = 0
		return e
	}

	drawdown := (equity - lastHigh) / lastHigh
	e.Drawdown = math.Round(drawdown*math.Pow10(DP)) / math.Pow10(DP)

	return e
}

// returns the equity point with the maximum drawdown
func (s Statistic) maxDrawdownPoint() (i int, ep EquityPoint) {
	if len(s.equity) == 0 {
		return 0, ep
	}

	var maxDrawdown = 0.0
	var index = 0

	for i, ep := range s.equity {
		if ep.Drawdown < maxDrawdown {
			maxDrawdown = ep.Drawdown
			index = i
		}
	}

	return index, s.equity[index]
}
