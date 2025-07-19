package text_test

import (
	"testing"

	"github.com/madamovych/go/text"
	"github.com/stretchr/testify/require"
)

func Test_ExprProcessor_Process_DummyVariable(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.ExprProcessorOf(false)
		require.Equal(t, " Hello ${name}! ", processor.Process(" Hello ${name}! "))
		processor.Define("name", "Mike")
		require.Equal(t, " Hello Mike! ", processor.Process(" Hello ${name}! "))
	})
}

func Test_ExprProcessor_Process_VariableNotDefined(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					require.Equal(t, "Cannot resolve ${name}", x)
				}
			}
		}()
		processor := text.ExprProcessorOf(true)
		processor.Process(" Hello ${name}! ")
		require.Fail(t, "panic expected")
	})
}

func Test_ExprProcessor_Process_ComplexVariable(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.ExprProcessorOf(false)
		processor.Define("person", map[string]any{"name": "Mike"})
		processor.Define("property", "name")
		require.Equal(t, " Hello Mike! ", processor.Process(" Hello #{person.${property}}! "))
	})
}

func Test_ExprProcessor_Process_DummyExpression(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.ExprProcessorOf(false)
		processor.Define("f", func(x, y int) int { return x + y })
		require.Equal(t, "4", processor.Process("#{ f(2, 2) }"))
		require.Equal(t, "2 + 2 = 4", processor.Process("2 + 2 = #{ f(2, 2) }"))
	})
}

func Test_ExprProcessor_Process_ComplexExpression(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.ExprProcessorOf(false)
		processor.Define("age", 30)
		require.Equal(t, "John is of age 30", processor.Process(`###{
			let person = fromJSON('{"name": "John", "age": "${age}"}');
			person.name + " is of age " + person.age
		}###`))
	})
}
