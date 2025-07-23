package env

type PropertySource interface {
	Name() string
	HasProperty(key string) bool
	Property(key string) string
	Properties() map[string]string
}
