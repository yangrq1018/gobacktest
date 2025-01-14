package gobacktest

// CommissionHandler is the basic interface for executing orders
type CommissionHandler interface {
	Calculate(f Fill) (float64, error)
}

// FixedCommission is a commission handler implementation which returns a fixed price commission
type FixedCommission struct {
	Commission float64
}

// Calculate calculates the commission of the trade
func (c *FixedCommission) Calculate(f Fill) (float64, error) {
	// no trade value, no commision
	if f.qty == 0 || f.price == 0 {
		return 0, nil
	}
	return c.Commission, nil
}

// TresholdFixedCommission is a commission handler implementation which returns a fixed price commission
// if the value of the trade is above a set treshold
type TresholdFixedCommission struct {
	Commission float64
	MinValue   float64
}

// Calculate calculates the commission of the trade
func (c *TresholdFixedCommission) Calculate(f Fill) (float64, error) {
	// no trade value, no commision
	if f.qty == 0 || f.price == 0 {
		return 0, nil
	}
	// minimum value of trade below treshold
	if c.MinValue > (float64(f.qty) * f.price) {
		return 0, nil
	}

	return c.Commission, nil
}

// PercentageCommission is a commission handler implementation which returns a percentage price commission
// calculated of the value of the trade
type PercentageCommission struct {
	Commission float64
}

// Calculate calculates the commission of the trade
func (c *PercentageCommission) Calculate(f Fill) (float64, error) {
	// no trade value, no commision
	if f.qty == 0 || f.price == 0 {
		return 0, nil
	}

	commission := float64(f.qty) * f.price * c.Commission

	return commission, nil
}

// ValueCommission is a commission handler implementation which returns a percentage price commission
// calculated of the value of the trade, if the value of the trade is within a given commission span
type ValueCommission struct {
	Commission    float64
	MinCommission float64
	MaxCommission float64
}

// Calculate calculates the commission of the trade
func (c *ValueCommission) Calculate(f Fill) (float64, error) {
	// no trade value, no commision
	if f.qty == 0 || f.price == 0 {
		return 0, nil
	}

	// value of trade below minimum commission
	if c.MinCommission > (float64(f.qty) * f.price * c.Commission) {
		return c.MinCommission, nil
	}

	// value of trade above maximum commission
	if c.MaxCommission < (float64(f.qty) * f.price * c.Commission) {
		return c.MaxCommission, nil
	}

	commission := float64(f.qty) * f.price * c.Commission

	return commission, nil
}
