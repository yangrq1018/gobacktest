package gobacktest

// Direction defines which direction a signal indicates
type Direction int

// different types of order directions
const (
	// BOT Buy
	BOT Direction = iota // 0
	// SLD Sell
	SLD
	// HLD Hold
	HLD
	// EXT Exit
	EXT
)

func (d Direction) String() string {
	switch d {
	case BOT:
		return "BUY"
	case SLD:
		return "SELL"
	case HLD:
		return "HOLD"
	case EXT:
		return "EXIT"
	default:
		return ""
	}
}

// Signal declares a basic signal event
type Signal struct {
	Event
	direction Direction // long, short, exit or hold
}

// Direction returns the Direction of a Signal
func (s Signal) Direction() Direction {
	return s.direction
}

// SetDirection sets the Directions field of a Signal
func (s *Signal) SetDirection(dir Direction) {
	s.direction = dir
}
