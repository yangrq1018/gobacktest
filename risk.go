package gobacktest

// RiskHandler is the basic interface for accessing risks of a portfolio
type RiskHandler interface {
	EvaluateOrder(OrderEvent, DataEvent, map[string]Position) (OrderEvent, error)
}

// Risk is a basic risk handler implementation
type Risk struct {
}

// EvaluateOrder handles the risk of an order, refines or cancel it
func (r *Risk) EvaluateOrder(o OrderEvent, data DataEvent, positions map[string]Position) (OrderEvent, error) {
	// simple implementation, just gives the received order back
	// no risk management
	return o, nil
}
