package sfox

type Side int

func (s Side) String() string {
	switch s {
	case Buy:
		return "buy"
	case Sell:
		return "sell"
	default:
		return ""
	}
}

const (
	_ Side = iota
	Buy
	Sell
)
