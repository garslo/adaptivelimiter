package adaptivelimiter

import "time"

func Differentiate(axis *PartitionedAxis) (float64, bool) {
	last2, ok := axis.LastN(2)
	if !ok {
		return 0, false
	}
	return (float64(last2[0]) - float64(last2[1])) /
		(time.Duration(axis.PartitionSize) * time.Nanosecond).Seconds(), true
}

type DerivativeTerm struct {
	Coefficient float64
}

func (me *DerivativeTerm) Value(axis *PartitionedAxis) (float64, bool) {
	d, ok := Differentiate(axis)
	if !ok {
		return 0, false
	}
	return me.Coefficient * d, true
}
