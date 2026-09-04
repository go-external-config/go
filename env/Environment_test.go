package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvironmentProfilesInclude(test *testing.T) {
	dir := test.TempDir()
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-prod.yaml"), []byte(`
profiles:
  include: kubernetes
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-kubernetes.yaml"), []byte(`
test:
  value: kubernetes
`), 0644))

	environment := Environment{
		profiles:              []string{"prod"},
		includedProfiles:      make([]string, 0),
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	location := filepath.ToSlash(dir) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"prod", "kubernetes"}, environment.activeProfiles())
	require.Equal(test, "kubernetes", environment.property("test.value"))
}

func TestEnvironmentProfilesIncludeAcrossLocations(test *testing.T) {
	dir1 := test.TempDir()
	dir2 := test.TempDir()

	require.NoError(test, os.WriteFile(filepath.Join(dir1, "application-kubernetes.yaml"), []byte(`
test:
  value: kubernetes
`), 0644))

	require.NoError(test, os.WriteFile(filepath.Join(dir2, "application-prod.yaml"), []byte(`
profiles:
  include: kubernetes
`), 0644))

	environment := Environment{
		profiles:              []string{"prod"},
		includedProfiles:      make([]string, 0),
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	location := filepath.ToSlash(dir1) + "/," + filepath.ToSlash(dir2) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"prod", "kubernetes"}, environment.activeProfiles())
	require.Equal(test, "kubernetes", environment.property("test.value"))
}

func TestEnvironmentProfilesIncludeTransitive(test *testing.T) {
	dir := test.TempDir()

	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-prod.yaml"), []byte(`
profiles:
  include: kubernetes
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-kubernetes.yaml"), []byte(`
profiles:
  include: aws
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-aws.yaml"), []byte(`
test:
  value: aws
`), 0644))

	environment := Environment{
		profiles:              []string{"prod"},
		includedProfiles:      make([]string, 0),
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	location := filepath.ToSlash(dir) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"prod", "kubernetes", "aws"}, environment.activeProfiles())
	require.Equal(test, "aws", environment.property("test.value"))
}

func TestEnvironmentProfilesActiveFromConfigurationWithExternalInclude(test *testing.T) {
	dir := test.TempDir()

	require.NoError(test, os.WriteFile(filepath.Join(dir, "application.yaml"), []byte(`
profiles:
  active: prod
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-prod.yaml"), []byte(`
test:
  active: prod
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-kubernetes.yaml"), []byte(`
test:
  include: kubernetes
`), 0644))

	environment := Environment{
		profiles:              make([]string, 0),
		includedProfiles:      []string{"kubernetes"},
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	location := filepath.ToSlash(dir) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"prod", "kubernetes"}, environment.activeProfiles())
	require.Equal(test, "prod", environment.property("test.active"))
	require.Equal(test, "kubernetes", environment.property("test.include"))
}

func TestEnvironmentProfilesIncludeDeduplicatesActiveProfiles(test *testing.T) {
	dir := test.TempDir()

	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-prod.yaml"), []byte(`
profiles:
  include: prod,kubernetes
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-kubernetes.yaml"), []byte(`
test:
  value: kubernetes
`), 0644))

	environment := Environment{
		profiles:              []string{"prod"},
		includedProfiles:      make([]string, 0),
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	location := filepath.ToSlash(dir) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"prod", "kubernetes"}, environment.activeProfiles())
	require.Equal(test, "kubernetes", environment.property("test.value"))
}

func TestEnvironmentIncludedProfileOverridesActiveProfile(test *testing.T) {
	dir := test.TempDir()

	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-prod.yaml"), []byte(`
profiles:
  include: kubernetes
test:
  value: prod
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-kubernetes.yaml"), []byte(`
test:
  value: kubernetes
`), 0644))

	environment := Environment{
		profiles:              []string{"prod"},
		includedProfiles:      make([]string, 0),
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	location := filepath.ToSlash(dir) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"prod", "kubernetes"}, environment.activeProfiles())
	require.Equal(test, "kubernetes", environment.property("test.value"))
}

func TestEnvironmentMatchesIncludedProfile(test *testing.T) {
	environment := Environment{
		profiles:         []string{"prod"},
		includedProfiles: []string{"kubernetes"},
	}

	require.True(test, environment.matchesProfiles("prod"))
	require.True(test, environment.matchesProfiles("kubernetes"))
	require.True(test, environment.matchesProfiles("prod & kubernetes"))
	require.False(test, environment.matchesProfiles("dev"))
}

func TestEnvironmentProfilesIncludeFromImport(test *testing.T) {
	dir := test.TempDir()

	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-prod.yaml"), []byte(`
config:
  import: profiles.yaml
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "profiles.yaml"), []byte(`
profiles:
  include: kubernetes
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-kubernetes.yaml"), []byte(`
test:
  value: kubernetes
`), 0644))

	environment := Environment{
		profiles:              []string{"prod"},
		includedProfiles:      make([]string, 0),
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	location := filepath.ToSlash(dir) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"prod", "kubernetes"}, environment.activeProfiles())
	require.Equal(test, "kubernetes", environment.property("test.value"))
}

func TestEnvironmentProfilesIncludeExternalPrecedence(test *testing.T) {
	dir := test.TempDir()

	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-env.yaml"), []byte(`
test:
  value: env
`), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-cli.yaml"), []byte(`
test:
  value: cli
`), 0644))

	environment := Environment{
		profiles:              make([]string, 0),
		includedProfiles:      make([]string, 0),
		paramsPropertySource:  MapPropertySourceOf("Application parameters"),
		environPropertySource: MapPropertySourceOf("Environment variables"),
		sources:               make([]PropertySource, 0),
		exprProcessor:         NewExprProcessor(false),
		strictExprProcessor:   NewExprProcessor(true),
	}

	environment.environPropertySource.SetProperty("PROFILES_INCLUDE", "env")
	environment.paramsPropertySource.SetProperty("profiles.include", "cli")

	environment.includeProfiles(environment.environPropertySource.properties["PROFILES_INCLUDE"])
	environment.includeProfiles(environment.paramsPropertySource.properties["profiles.include"])

	location := filepath.ToSlash(dir) + "/"
	environment.discoverProfiles(location, "application")
	environment.sources = environment.sources[:0]
	environment.loadConfigurations(location, "application")

	require.Equal(test, []string{"env", "cli"}, environment.activeProfiles())
	require.Equal(test, "cli", environment.property("test.value"))
}
