package env_test

import (
	"testing"
	"time"

	"github.com/go-external-config/go/env"
	"github.com/stretchr/testify/require"
)

func Test_ExprProcessor_Process_DummyVariable(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		propertySource := env.MapPropertySourceOf("map")
		env.SetActiveProfiles("").
			WithPropertySource(propertySource)

		processor := env.ExprProcessorOf(false)
		require.Equal(t, " Hello ${name}! ", processor.Process(" Hello ${name}! "))
		require.Equal(t, " Hello Unknown! ", processor.Process(" Hello ${name:Unknown}! "))
		propertySource.SetProperty("name", "Mike")
		require.Equal(t, " Hello Mike! ", processor.Process(" Hello ${name}! "))
	})
}

func Test_ExprProcessor_Process_VariableNotDefined(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					require.Equal(t, "Cannot resolve property ${name}", x)
				}
			}
		}()
		env.SetActiveProfiles("")

		processor := env.ExprProcessorOf(true)
		processor.Process(" Hello ${name}! ")
		require.Fail(t, "panic expected")
	})
}

func Test_ExprProcessor_Process_ComplexVariable(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		propertySource := env.MapPropertySourceOf("map")
		env.SetActiveProfiles("").
			WithPropertySource(propertySource)
		propertySource.SetProperty("property", "name")

		processor := env.ExprProcessorOf(false)
		processor.Define("person", map[string]any{"name": "Mike"})
		require.Equal(t, " Hello Mike! ", processor.Process(" Hello #{person.${property}}! "))
	})
}

func Test_ExprProcessor_Process_DummyExpression(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := env.ExprProcessorOf(false)
		processor.Define("f", func(x, y int) int { return x + y })
		require.Equal(t, 4, processor.Process("#{ f(2, 2) }"))
		require.Equal(t, "2 + 2 = 4", processor.Process("2 + 2 = #{ f(2, 2) }"))
		require.Equal(t, []string{"dev", "test"}, processor.Process("#{split('dev,test', ',')}"))
	})
}

func Test_ExprProcessor_Process_TimeConstants(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		processor := env.ExprProcessorOf(false)
		require.Equal(t, 10*time.Millisecond, processor.Process("#{10 * time.Millisecond}"))
		require.Equal(t, 2*time.Hour, processor.Process("#{2 * time.Hour}"))
		require.Equal(t, 24*time.Hour, processor.Process("#{time.Day}"))
	})
}

func Test_ExprProcessor_Process_ComplexExpression(t *testing.T) {
	t.Run("should substitute variable", func(t *testing.T) {
		propertySource := env.MapPropertySourceOf("map")
		env.SetActiveProfiles("").
			WithPropertySource(propertySource)
		processor := env.ExprProcessorOf(false)
		propertySource.SetProperty("age", "30")
		require.Equal(t, "John is of age 30", processor.Process(`###{
			let person = fromJSON('{"name": "John", "age": "${age}"}');
			person.name + " is of age " + person.age
		}###`))
	})
}
