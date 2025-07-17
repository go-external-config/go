package text_test

import (
	"testing"

	"github.com/madamovych/go/text"
	"github.com/stretchr/testify/require"
)

func Test_ExprProcessor_Process_DummyVariable(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.NewExprProcessor()
		processor.Define("name", "Mike")
		require.Equal(t, " Hello Mike! ", processor.Process(" Hello ${name}! "))
	})
}

func Test_ExprProcessor_Process_StructVariable(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.NewExprProcessor()
		processor.Define("person", map[string]any{"name": "Mike"})
		require.Equal(t, " Hello Mike! ", processor.Process(" Hello ${person.name}! "))
	})
}

func Test_ExprProcessor_Process_DummyExpression(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.NewExprProcessor()
		processor.Define("f", func(x, y int) int { return x + y })
		require.Equal(t, "2 + 2 = 4", processor.Process("2 + 2 = ${ f(2, 2) }"))
	})
}

func Test_ExprProcessor_Process_MultilineExpression(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := text.NewExprProcessor()
		processor.Define("age", 30)
		require.Equal(t, "John is of age 30", processor.Process(`$${
			let person = fromJSON('{"name": "John", "age": "${age}"}');
			person.name + " is of age " + person.age
		}$`))
	})
}
