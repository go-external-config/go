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
		PatternProcessor: *PatternProcessorOf(`\#\#\#\{(?P<complex>([^\$]|\$[^\{])*?)\}\#\#\#|\#\{(?P<expr>([^\$]|\$[^\{])*?)\}|\$\{(?P<prop>([^\$]|\$[^\{])*?)\}`),
		context:          make(map[string]any)}
	processor.OverrideResolve(processor.Resolve)
	return &processor
}

func (p *ExprProcessor) Resolve(match *lang.RegexpMatch,
	super func(*lang.RegexpMatch) string) string {

	prop := match.NamedGroup("prop")
	if prop.Present() {
		return fmt.Sprintf("%v", lang.OptionalOfEntry(p.context, prop.Value()).OrElse(match.Expr()))
	}

	expression := lang.FirstNonEmpty(match.NamedGroup("expr").OrElse(""), match.NamedGroup("complex").OrElse(""))
	return fmt.Sprintf("%v", lang.OptionalOfCommaErr(expr.Eval(expression, p.context)).
		OrElsePanic("Cannot evaluate '%s'", expression))
}

func (p *ExprProcessor) Define(key string, value any) {
	p.context[key] = value
}

func (p *ExprProcessor) Reset() {
	p.context = make(map[string]any)
}
