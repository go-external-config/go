package env

type PropertySource interface {
	Properties() map[string]any
}
