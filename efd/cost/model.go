package cost

import "github.com/mmcloughlin/ec3/internal/errutil"

type Model interface {
	Weight(Operation) float64
}

type Weights struct {
	I      float64
	M      float64
	S      float64
	Pow    float64
	ParamM float64
	Add    float64
	ConstM float64
}

func (w Weights) Weight(op Operation) float64 {
	switch operation := op.(type) {
	case Inv:
		return w.I
	case Mul:
		return w.M
	case Pow:
		if operation.IsSquare() {
			return w.S
		}
		return w.Pow
	case ParamMul:
		return w.ParamM
	case ConstMul:
		return w.ConstM
	case Add:
		return w.Add
	default:
		panic(errutil.UnexpectedType(operation))
	}
}
