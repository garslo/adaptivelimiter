package adaptivelimiter

func Integrate(axis *PartitionedAxis, count int) (float64, bool) {
	data, ok := axis.LastN(count)
	if !ok {
		return 0, false
	}
	sum := float64(0)
	for _, d := range data {
		sum += d
	}
	return sum, true
}

type IntegralTerm struct {
	Coefficient float64
	History     int
}

func (me *IntegralTerm) Value(axis *PartitionedAxis) (float64, bool) {
	integral, ok := Integrate(axis, me.History)
	if !ok {
		return 0, false
	}
	return me.Coefficient * integral, true
}
