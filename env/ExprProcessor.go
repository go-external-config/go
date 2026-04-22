package env

import (
	"fmt"
	"runtime"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/compiler"
	"github.com/expr-lang/expr/conf"
	"github.com/expr-lang/expr/parser"
	"github.com/go-errr/go/err"
	"github.com/go-external-config/go/lang"
	"github.com/go-external-config/go/util/optional"
	"github.com/go-external-config/go/util/regex"
)

// See expr-lang: https://expr-lang.org/docs/language-definition
type ExprProcessor struct {
	regex.PatternProcessor
	context map[string]any
	strict  bool
}

func ExprProcessorOf(strict bool) *ExprProcessor {
	processor := ExprProcessor{
		PatternProcessor: *regex.PatternProcessorOf(`\#\#\#\{(?P<complex>([^\$#]\{|[^\{])*?)\}\#\#\#|\#\{(?P<expr>([^\$#]\{|[^\{])*?)\}|\$\{(?P<prop>([^\$#:]\{|[^\{\}:])*)(:(?P<defaultValue>([^\$#]\{|[^\{])*?))?\}`),
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
		"Week":        7 * 24 * time.Hour,
	}
	processor.context["size"] = map[string]any{
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}
	processor.context["runtime"] = map[string]any{
		"NumCPU": runtime.NumCPU(),
	}
	return &processor
}

func (this *ExprProcessor) Resolve(match *regex.Match,
	super func(*regex.Match) any) (resolved any) {
	if !this.strict {
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
			panic(err.NewRuntimeException(fmt.Sprintf("Cannot resolve property %s", match.Expr())))
		}
	} else {
		expression := lang.FirstNonEmpty(match.NamedGroup("expr").OrElse(""), match.NamedGroup("complex").OrElse(""))
		resolved = optional.OfNilable(this.eval(expression, this.context)).OrElsePanic("Cannot evaluate expression %s", match.Expr())
	}
	// fmt.Printf("ExprProcessor: %s -> %s\n", match.Expr(), resolved)
	return resolved
}

func (this *ExprProcessor) Define(key string, value any) {
	this.context[key] = value
}

func (this *ExprProcessor) Reset() {
	this.context = make(map[string]any)
}

func (this *ExprProcessor) eval(input string, env any) any {
	config := conf.CreateNew()
	config.Strict = true
	tree := optional.OfCommaErr(parser.Parse(input)).OrElsePanic("Cannot parse expression")
	program := optional.OfCommaErr(compiler.Compile(tree, config)).OrElsePanic("Cannot compile expression")
	output := optional.OfCommaErr(expr.Run(program, env)).OrElsePanic("Cannot evaluate expression")
	return output
}

func (this *ExprProcessor) SetStrict(strict bool) {
	this.strict = strict
}
