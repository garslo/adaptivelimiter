package adaptivelimiter

type ProportionalTerm struct {
	Coefficient float64
}

func (me *ProportionalTerm) Value(axis *PartitionedAxis) (float64, bool) {
	last, ok := axis.LastN(1)
	if !ok {
		return 0, false
	}
	return me.Coefficient * last[0], true
}
