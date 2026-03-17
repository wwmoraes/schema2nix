package nix

import (
	"fmt"
)

var (
	_ intLimit = (*intUnlimited)(nil)
	_ intLimit = (*intMinimum)(nil)
	_ intLimit = (*intMaximum)(nil)
	_ intLimit = (*intBetween)(nil)
)

// intLimit describes an integer value limits. They may be unlimited, have a
// minimum, a maximum or both.
type intLimit interface {
	fmt.Stringer
}

type (
	intUnlimited struct{}
	intMinimum   struct{ Value int64 }
	intMaximum   struct{ Value int64 }
	intBetween   struct{ Minimum, Maximum int64 }
)

func (intUnlimited) String() string {
	return "int"
}

func (minimum intMinimum) String() string {
	switch minimum.Value {
	case 0:
		return "ints.unsigned"
	case 1:
		return "ints.positive"
	default:
		return fmt.Sprintf("addCheck int (x: x >= %d)", minimum.Value)
	}
}

func (maximum intMaximum) String() string {
	return fmt.Sprintf("addCheck int (x: x <= %d)", maximum.Value)
}

func (between intBetween) String() string {
	return fmt.Sprintf("ints.between %d %d", between.Minimum, between.Maximum)
}
