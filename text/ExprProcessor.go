package text

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/madamovych/go/lang"
)

type ExprProcessor struct {
	PatternProcessor
	context map[string]any
}

func NewExprProcessor() *ExprProcessor {
	processor := ExprProcessor{
		PatternProcessor: *PatternProcessorOf("\\$\\$\\{(?P<long>([^\\$]|\\$[^\\{])*?)\\}\\$|\\$\\{(?P<short>([^\\$]|\\$[^\\{])*?)\\}"),
		context:          make(map[string]any)}
	processor.OverrideResolve(processor.Resolve)
	return &processor
}

func (p *ExprProcessor) Resolve(match *lang.RegexpMatch,
	super func(*lang.RegexpMatch) string) string {

	expression := lang.FirstNonEmpty(match.NamedGroup("long"), match.NamedGroup("short"))
	value, ok := p.context[expression]
	if ok {
		return fmt.Sprintf("%v", value)
	}
	return fmt.Sprintf("%v", lang.OptionalOfCommaErr(expr.Eval(expression, p.context)).
		OrElsePanic(fmt.Sprintf("Cannot evaluate '%s'", expression)))
}

func (p *ExprProcessor) Define(key string, value any) {
	p.context[key] = value
}

func (p *ExprProcessor) Reset() {
	p.context = make(map[string]any)
}
