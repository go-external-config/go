package env

import (
	"fmt"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/parser"
	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/optional"
	"github.com/go-external-config/go/util/regex"
	"github.com/go-external-config/go/util/text"
)

type ExprProcessor struct {
	text.PatternProcessor
	context map[string]any
	strict  bool
}

func ExprProcessorOf(strict bool) *ExprProcessor {
	processor := ExprProcessor{
		PatternProcessor: *text.PatternProcessorOf(`\#\#\#\{(?P<complex>([^\$#]\{|[^\{])*?)\}\#\#\#|\#\{(?P<expr>([^\$#]\{|[^\{])*?)\}|\$\{(?P<prop>([^\$#:]\{|[^\{\}:])*)(:(?P<defaultValue>([^\$#]\{|[^\{])*?))?\}`),
		context:          make(map[string]any),
		strict:           strict}
	processor.OverrideResolve(processor.Resolve)
	processor.context["time"] = map[string]any{
		"Nanosecond":  time.Nanosecond,
		"Microsecond": time.Microsecond,
		"Millisecond": time.Millisecond,
		"Second":      time.Second,
		"Minute":      time.Minute,
		"Hour":        time.Hour,
		"Day":         24 * time.Hour,
	}
	return &processor
}

func (p *ExprProcessor) Resolve(match *regex.Match,
	super func(*regex.Match) any) (resolved any) {
	if !p.strict {
		defer func() {
			if recover() != nil {
				resolved = match.Expr()
			}
		}()
	}
	prop := match.NamedGroup("prop")
	if prop.Present() {
		resolvedValue := Instance().lookupRawProperty(prop.Value())
		defaultValue := match.NamedGroup("defaultValue")
		if resolvedValue.Present() {
			resolved = fmt.Sprint(resolvedValue.Value())
		} else if defaultValue.Present() {
			resolved = defaultValue.Value()
		} else {
			panic(fmt.Sprintf("Cannot resolve property %s", match.Expr()))
		}
	} else {
		expression := lang.FirstNonEmpty(match.NamedGroup("expr").OrElse(""), match.NamedGroup("complex").OrElse(""))
		resolved = optional.OfNilable(p.eval(expression, p.context)).OrElsePanic("Cannot evaluate expression %s", match.Expr())
	}
	// slog.Debug(fmt.Sprintf("ExprProcessor: %s -> %s\n", match.Expr(), resolved))
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
	tree := optional.OfCommaErr(parser.Parse(input)).OrElsePanic("Cannot parse expression")
	program := optional.OfCommaErr(compiler.Compile(tree, config)).OrElsePanic("Cannot compile expression")
	output := optional.OfCommaErr(expr.Run(program, env)).OrElsePanic("Cannot evaluate expression")
	return output
}

func (p *ExprProcessor) SetStrict(strict bool) {
	p.strict = strict
}
