package transformations

import (
	"github.com/Knetic/govaluate"

	errors "github.com/capeprivacy/cape/partyerrors"
)

type Conditional struct {
	expression *govaluate.EvaluableExpression
}

func NewConditional(expression string) (*Conditional, error) {
	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return nil, err
	}

	return &Conditional{expression: exp}, nil
}

func (c *Conditional) Evaluate(params map[string]interface{}) (bool, error) {
	res, err := c.expression.Evaluate(params)
	if err != nil {
		return false, err
	}

	shouldFilter, ok := res.(bool)
	if !ok {
		return false, errors.New(EvaluateBoolOnly, "Conditional expressions should only evaluate to booleans")
	}

	return shouldFilter, nil
}

func (c *Conditional) Vars() []string {
	return c.expression.Vars()
}
