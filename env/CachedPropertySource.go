package env

import (
	"fmt"
	"strings"

	"github.com/go-errr/go/err"
	"github.com/go-jang/go/util/concurrent"
)

const CACHED_VALUE_PREFIX = "cached:"

// CachedPropertySource resolves and caches property values prefixed with "cached:".
//
// The wrapped value is resolved only once per property key for the lifetime of
// the Environment. Subsequent lookups return the cached result.
//
// This is useful for generators and dynamic placeholders whose values should be
// stable across the application run.
//
// For example:
//
//	node.id=cached:${random.uuid}
//
// The first lookup of "node.id" resolves "${random.uuid}" and stores
// the generated UUID. All subsequent lookups of "node.id" return the
// same UUID.
//
// This is a common requirement for values such as application instance IDs,
// correlation prefixes, node identifiers, or ephemeral secrets that should be
// generated once at startup and reused consistently throughout the application.
//
// Properties without the "cached:" prefix retain their original behavior. For
// example:
//
//	request.id=${random.uuid}
//
// resolves "${random.uuid}" on every lookup, producing a new UUID each time.
// This is useful for values that are expected to be unique, such as request,
// message, or task identifiers.
type CachedPropertySource struct {
	cachedProperties *concurrent.HashMap[string, string]
}

func NewCachedPropertySource() *CachedPropertySource {
	return &CachedPropertySource{
		cachedProperties: concurrent.NewHashMap[string, string]()}
}

func (this *CachedPropertySource) Name() string {
	return "CachedPropertySource"
}

func (this *CachedPropertySource) HasProperty(key string) bool {
	if this.cachedProperties.ContainsKey(key) {
		return true
	}
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			return strings.HasPrefix(source.Property(key), CACHED_VALUE_PREFIX)
		}
	}
	return false
}

func (this *CachedPropertySource) Property(key string) string {
	if this.cachedProperties.ContainsKey(key) {
		return this.cachedProperties.Get(key)
	}
	for _, source := range environment.PropertySources() {
		if source.Properties() != nil && source.HasProperty(key) {
			value := source.Property(key)[len(CACHED_VALUE_PREFIX):]
			resolved := environment.ResolveRequiredPlaceholders(value)
			return this.cachedProperties.PutIfAbsent(key, fmt.Sprint(resolved))
		}
	}
	panic(err.NewIllegalArgumentException("No value present for " + key))
}

func (this *CachedPropertySource) Properties() map[string]string {
	return nil
}
