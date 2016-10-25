package adaptivelimiter

type Term interface {
	Value(*PartitionedAxis) (float64, bool)
}

func NewPIDController(axis *PartitionedAxis, p, i, d float64, history int) *PIDController {
	return &PIDController{
		axis,
		&ProportionalTerm{p},
		&IntegralTerm{i, history},
		&DerivativeTerm{d},
	}
}

type PIDController struct {
	ErrorAxis    *PartitionedAxis
	Proportional Term
	Derivative   Term
	Integral     Term
}

func (me *PIDController) Update(errorValue float64) {
	me.ErrorAxis.Inc(errorValue)
}

func (me *PIDController) Delta() (float64, bool) {
	p, p_ok := me.Proportional.Value(me.ErrorAxis)
	d, d_ok := me.Derivative.Value(me.ErrorAxis)
	i, i_ok := me.Integral.Value(me.ErrorAxis)
	if p_ok && d_ok && i_ok {
		return p + d + i, true
	}
	return 0, false
}
