package text

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/parser"
	"github.com/madamovych/go/lang"
)

type ExprProcessor struct {
	PatternProcessor
	context map[string]any
	strict  bool
}

func ExprProcessorOf(strict bool) *ExprProcessor {
	processor := ExprProcessor{
		PatternProcessor: *PatternProcessorOf(`\#\#\#\{(?P<complex>([^\$#]\{|[^\{])*?)\}\#\#\#|\#\{(?P<expr>([^\$#]\{|[^\{])*?)\}|\$\{(?P<prop>([^\$#]\{|[^\{])*?)\}`),
		context:          make(map[string]any),
		strict:           strict}
	processor.OverrideResolve(processor.Resolve)
	return &processor
}

func (p *ExprProcessor) Resolve(match *lang.RegexpMatch,
	super func(*lang.RegexpMatch) string) (resolved string) {
	if !p.strict {
		defer func() {
			if recover() != nil {
				resolved = match.Expr()
			}
		}()
	}
	prop := match.NamedGroup("prop")
	if prop.Present() {
		resolved = fmt.Sprintf("%v", lang.OptionalOfEntry(p.context, prop.Value()).OrElsePanic("Cannot resolve %s", match.Expr()))
	} else {
		expression := lang.FirstNonEmpty(match.NamedGroup("expr").OrElse(""), match.NamedGroup("complex").OrElse(""))
		resolved = fmt.Sprintf("%v", lang.OptionalOfNilable(p.eval(expression, p.context)).OrElsePanic("Cannot resolve %s", match.Expr()))
	}
	// fmt.Printf("ExprProcessor: %s -> %s\n", match.Expr(), resolved)
	return resolved
}

func (p *ExprProcessor) Define(key string, value any) {
	p.context[key] = value
}

func (p *ExprProcessor) Reset() {
	p.context = make(map[string]any)
}

func (p *ExprProcessor) eval(input string, env any) any {
	config := conf.CreateNew()
	config.Strict = true
	tree := lang.OptionalOfCommaErr(parser.Parse(input)).OrElsePanic("Cannot parse expression")
	program := lang.OptionalOfCommaErr(compiler.Compile(tree, config)).OrElsePanic("Cannot compile expression")
	output := lang.OptionalOfCommaErr(expr.Run(program, env)).OrElsePanic("Cannot evaluate expression")
	return output
}

func (p *ExprProcessor) SetStrict(strict bool) {
	p.strict = strict
}
