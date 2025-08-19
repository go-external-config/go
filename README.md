# go-external-config

[![Go Reference](https://pkg.go.dev/badge/github.com/go-external-config/go.svg)](https://pkg.go.dev/github.com/go-external-config/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-external-config/go)](https://goreportcard.com/report/github.com/go-external-config/go)
[![Release](https://img.shields.io/github/v/release/go-external-config/go)](https://github.com/go-external-config/go/releases)

go-external-config lets you externalize your configuration so that you can work with the same application code in different environments. You can use a variety of external configuration sources including properties files, YAML files, environment variables, and command-line arguments.

Property values can be injected directly into your code by using the `env.Value` function, or be bound to structured objects through `env.ConfigurationProperties`

go-external-config uses a very particular `PropertySource` order that is designed to allow sensible overriding of values. Later property sources can override the values defined in earlier ones.

Sources are considered in the following order:  
1. Config data (such as application.properties files).
2. OS environment variables.
3. Command line arguments.

Config data files are considered in the following order:  
1. Application properties (application.properties and YAML variants).
2. Profile-specific application properties (application-{profile}.properties and YAML variants).

> It is recommended to stick with one format for your entire application. If you have configuration files with both .properties and YAML format in the same location, .properties takes precedence.

To provide a concrete example, suppose you develop a component that uses a name property, as shown in the following example:

	name := env.Value[string]("${name}")

## Accessing Command Line Properties

By default, go-external-config converts any command line option arguments (that is, arguments starting with `--` (or `-`), such as `--server.port=9000`) to a property and adds them to the `Environment`. As mentioned previously, command line properties always take precedence over file-based property sources.

## External Application Properties

go-external-config will automatically find and load `application.properties` and `application.yaml` files from the following locations when your application starts:  
1. The current directory
2. The `config/` subdirectory in the current directory

The list is ordered by precedence (with values from lower items overriding earlier ones). Documents from the loaded files are added as `PropertySource` instances to the  `Environment`.

If you do not like application as the configuration file name, you can switch to another file name by specifying a `config.name` property (or `CONFIG_NAME` environment variable). For example, to look for `myproject.properties` and `myproject.yaml` files you can run your application as follows:

	go run ./cmd/myproject/ --config.name=myproject

You can also refer to an explicit location by using the `config.location` property (or `CONFIG_LOCATION` environment variable). This property accepts a comma-separated list of one or more locations to check.

The following example shows how to specify two distinct files:  

	go run ./cmd/myproject/ --config.location=\
		default.properties,\
		override.properties

> Locations are optional and application will not fail if they do not exist.

If `config.location` contains directories (as opposed to files), they should end in `/`. At runtime they will be appended with the names generated from `config.name` before being loaded. Files specified in `config.location` are imported directly.

> Both directory and file location values are also expanded to check for profile-specific files. For example, if you have a `config.location` of `myconfig.properties`, you will also find appropriate `myconfig-<profile>.properties` files are loaded.

In most situations, each `config.location` item you add will reference a single file or directory. Locations are processed in the order that they are defined and later ones can override the values of earlier ones.

If you have a complex location setup, and you use profile-specific configuration files, you may need to provide further hints so that go-external-config knows how they should be grouped. A location group is a collection of locations that are all considered at the same level. Items within a location group should be separated with `;`. See the example in the Profile Specific Files section for more details.

Locations configured by using `config.location` replace the default locations. If you prefer to add additional locations, rather than replacing them, you can use `config.additional-location` property (or `CONFIG_ADDITIONAL_LOCATION` environment variable). Properties loaded from additional locations can override those in the default locations.

This search ordering lets you specify default values in one configuration file and then selectively override those values in another. You can provide default values for your application in `application.properties` (or whatever other basename you choose with `config.name`) in one of the default locations. These default values can then be overridden at runtime with a different file located in one of the custom locations.

## Profile Specific Files

As well as `application` property files, go-external-config will also attempt to load profile-specific files using the naming convention `application-{profile}`. For example, if your application activates a profile named `prod` and uses YAML files, then both `application.yaml` and `application-prod.yaml` will be considered.

Profile-specific properties are loaded from the same locations as standard `application.properties`, with profile-specific files always overriding the non-specific ones. If several profiles are specified, a last-wins strategy applies. For example, if profiles `prod,live` are specified by the `profiles.active` property (or `PROFILES_ACTIVE` environment variable), values in `application-prod.properties` can be overridden by those in `application-live.properties`.

The last-wins strategy applies at the location group level. A `config.location` of `cfg/,ext/` will not have the same override rules as `cfg/;ext/`.
For example, continuing our prod,live example above, we might have the following files:  

	cfg
		application-live.properties
	ext
		application-live.properties
		application-prod.properties

When we have a `config.location` of `cfg/,ext/` we process all `cfg` files before all `ext` files:  

	cfg/application-live.properties
	ext/application-prod.properties
	ext/application-live.properties

When we have `cfg/;ext/` instead (with a `;` delimiter) we process `cfg` and `ext` at the same level:  

	ext/application-prod.properties
	cfg/application-live.properties
	ext/application-live.properties

## Using Environment Variables

When running applications on a cloud platform (such as Kubernetes) you often need to read config values that the platform supplies. Assume there’s an environment variable called `CLUSTER`:

	my.name=Service1
	my.cluster=${CLUSTER}

## Property Placeholders

The values in `application.properties` and `application.yaml` are filtered through the existing `Environment` when they are used, so you can refer back to previously defined values (for example, from environment variables). The standard `${name}` property-placeholder syntax can be used anywhere within a value. Property placeholders can also specify a default value using a `:` to separate the default value from the property name, for example `${name:default}`.

The use of placeholders with and without defaults is shown in the following example:  

	app.name=MyApp
	app.description=${app.name} is an application written by ${username:Unknown}

Assuming that the `username` property has not been set elsewhere, `app.description` will have the value `MyApp is an application written by Unknown`.

You can also use this technique to create “short” variants of existing properties. Some people like to use (for example) `--port=9000` instead of `--server.port=9000` to set configuration properties on the command line. You can enable this behavior by using placeholders in `application.properties`, as shown in the following example:  

	server.port=${port:8080}

## Encrypting Properties

go-external-config does not provide any built-in support for encrypting property values, however, it does provide the hook point necessary to modify values contained in the Environment. See [Base64PropertySource](https://github.com/go-external-config/go/blob/main/env/Base64PropertySource.go) as an example. The same way you can load property from external location, for example AWS Systems Manager Parameter Store etc.

## Working With YAML

YAML is a superset of JSON and, as such, is a convenient format for specifying hierarchical configuration data. The go-external-config automatically supports YAML as an alternative to properties. 

YAML documents need to be converted from their hierarchical format to a flat structure that can be used with the `Environment`. For example, consider the following YAML document:  

	environments:
		dev:
			url: "https://dev.example.com"
			name: "Developer Setup"
		prod:
			url: "https://another.example.com"
			name: "My Cool App"

In order to access these properties from the `Environment`, they would be flattened as follows:  

	environments.dev.url=https://dev.example.com
	environments.dev.name=Developer Setup
	environments.prod.url=https://another.example.com
	environments.prod.name=My Cool App

Likewise, YAML lists also need to be flattened. They are represented as property keys with [index] dereferencers. For example, consider the following YAML:  

	my:
		servers:
		- "dev.example.com"
		- "another.example.com"

The preceding example would be transformed into these properties:  

	my.servers[0]=dev.example.com
	my.servers[1]=another.example.com

## Base64 encoding

[Base64PropertySource](https://github.com/go-external-config/go/blob/main/env/Base64PropertySource.go) (available by default) is useful for decoding property values in Base64 format as shown in the following example:  

	my.secret=base64:c2FjcmVkIHNlY3JldA==

## Configuration Properties

Using the `env.Value[string]("${property}")` to inject configuration properties can sometimes be cumbersome, especially if you are working with multiple properties or your data is hierarchical in nature. go-external-config provides an alternative method of working with properties that lets strongly typed fields govern and validate the configuration of your application. It is possible to bind struct properties as shown in the following example:

	var db struct {
		Host  string
		Port  int
		Username string
		Password string
	}
	
	env.ConfigurationProperties("db", &db)

> Host value will be looked-up in `db.Host` and `db.host` properties

## Properties Conversion

go-external-config attempts to coerce the external application properties to the right type when it binds the `env.Value[type]()` or the `env.ConfigurationProperties()`.

## Expression Language

go-external-config provides support for [expr-lang](https://github.com/expr-lang/expr). Consider the following example:  

	servers=host1,host2,host3

 	servers := env.Value[[]string]("#{split('${servers}', ',')}")

## Profiles

go-external-config provide a way to segregate parts of your application configuration and make it be available only in certain environments. Any `Bean` can be created with `Profile` to limit when it is loaded, as shown in the following example ([go-beans](https://github.com/go-beans/go) dependency required):

	ioc.Bean[http.Client]().Profile("prod").Factory(func() *http.Client {
		return &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
	}).Register()

You can use a `profiles.active` `Environment` property to specify which profiles are active. You can specify the property in any of the ways described earlier in this chapter. For example, you could include it in your `application.properties`, as shown in the following example:

	profiles.active=dev,hsqldb

If no profile is active, a default profile is enabled.

The `profiles.active` property follows the same ordering rules as other properties. The highest `PropertySource` wins. This means that you can specify active profiles in `application.properties` and then replace them by using the command line switch.

## Programmatically Setting Profiles

You can programmatically set active profiles by calling `env.SetActiveProfiles("...")` before your application runs. This can be useful for tests to mock `Bean`s or other scenarious.

## Credits

[Spring Externalized Configuration](https://docs.spring.io/spring-boot/reference/features/external-config.html)

## Installation

```bash
go get github.com/go-external-config/go
```
### See also
[github.com/go-beans/go](https://github.com/go-beans/go)
