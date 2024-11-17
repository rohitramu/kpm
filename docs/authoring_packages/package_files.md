# Files in a template package

## `package.yaml`

The package information file defines the name and version of a package.  For example, a `package.yaml` file for version `1.0.0` of the package `kpmtool/helloworld` would look like this:

```yaml
name: kpmtool/helloworld
version: 1.0.0
```

The name of the package may only contain lowercase letters, numbers, forward slashes, underscores, dashes and dots. Also, it must start with a lowercase letter.
The recommended convention for naming packages is to use dots for separating segments (i.e. creating a heirarchy), and using underscores to separate words inside a segment:

```sh
my_username/my_organization.my_product.my_package_name
```

The version must be in the format `"major.minor.revision"`.  Leading zeros are not permitted in the `major`, `minor` or `revision` segments of the version string, however a segment may be just `0` (zero).  The zero version (`0.0.0`) is not allowed.

NOTE: This file will not be evaluated as a template.

## `parameters.yaml`

A package author must provide a set of default parameters which will be used whenever a user does not provide a parameter.  The default parameters file is also a great place to document each parameter using comments.

NOTE: This file will not be evaluated as a template.

For example, here is a sample default parameters file:

```yaml
# The user's name
name:
  first: "Mr."
  last: 'FooBar'

# Whether or not to include the custom object
isCustom: true

# Provide your custom object
custom-object:
  hello: world!
  other:
  - is
  - my
  - true
  - custom: object

# Choose some colors
colors:
- red
- green
- blue
```

## `interface.yaml`

The interface is a YAML template which defines what parameters the package requires in order to correctly generate output.  Parameters which are provided by the user are used as the input to this interface.  The resulting YAML is then used as the input to all other templates in the package.

If a value is not provided by the user for a parameter, the default value will be used.  Default values are defined in the [parameters](#parametersyaml) file.

The values in the interface may be defined as "constants" (i.e. hardcoded values) which can be referenced in templates.  However, for more complex string values (e.g. multi-line strings, strings with special characters, etc.), [helper templates](#helpers) should be preferred.  See the `constantGreeting` property or the `favorite-things` list below for examples of how hardcoded values may be defined.

Here is an example interface definition which can accept the parameters from the [example](#parametersyaml) above:

```yaml
username: {{ .name.first }} {{ .name.last }}
someColors:
{{- range .colors }}
- {{ . }}
{{- end }}
constantGreeting: Hello, World!
favorite-things:
- Music
- Photography
- Car racing
- Basketball
- Software!
isCustom: {{ .isCustom }}
customObj: {{- index . "custom-object" | toYaml | nindent 2 }}
```

Since the interface definition is itself a template, all of the normal template functions are available to be used.  See the [templates](#templates) section for more details.

Once the interface is executed with the provided parameters, it is combined with the package definition to become the input to all other templates.  Here is an example of what is provided to all templates if we ran the `kpmktool/example` package with the above interface file and the parameters from the [parameters example](#parametersyaml):

```yaml
package:
  name: kpmtool/example
  version: 1.0.0
values:
  username: Mr. FooBar
  someColors:
  - red
  - green
  - blue
  constantGreeting: Hello, World!
  favorite-things:
  - Music
  - Photography
  - Car racing
  - Basketball
  - Software!
  isCustom: true
  customObj:
    hello: world!
    other:
    - is
    - my
    - true
    - custom: object
```

This is why all templates in the package (other than the interface) need to reference the ".values" object to get the values supplied by the interface.

## `templates/`

Files in the templates directory are the text templates which will be used to generate the output files.  These can be used to generate any text format with any filename.

Here is an example of a template file which generates a text file by using values provided by the [interface example](#interfaceyaml) above:

```sh
Hello, {{ .values.username }}!

Are you enjoying KPM?  Let me know if you have any issues: https://github.com/rohitramu/kpm/issues

You can iterate over an array like this:
{{- range .values.someColors }}
 - {{ . }}
{{- end }}

Wondering how to access values that have special characters in their names?  Here is a list of my favorite things:
{{- range $thing := (index .values "favorite-things") }}
 - {{ $thing }}
{{- end }}

Maybe you'd like to indent some text?  Here you go (indented by 2 spaces): {{- include "nesting helper" . | nindent 2 }}

If you want to see an example of a template which generates yaml output, take a look at the "object.yaml" file in the templates folder!

Hopefully these examples gave you an idea of how you can use KPM to create template packages which generate output files in any format, potentially using fairly complex logic :).
```

And here is one more example which outputs a YAML file:

```yaml
# object.yaml
configuration:
  package-name: {{ .package.name }}
  colors: {{- include "colors helper" . | trim | nindent 2 }}
  your-name: {{ .values.username }}
  {{- if .values.isCustom }}
  custom:
    your-object: {{- .values.customObj | toYaml | trim | nindent 6 }}
  {{- end }}
```

## `helpers/`

Helper templates (a.k.a. "named templates" or "partial templates") allow the sharing of logic between templates in the package.  If you find yourself copying and pasting parts of templates, defining helper templates will allow you to simplify your templates and reduce the likelihood of copy-paste errors.

Helper templates can be inserted by using either the [`template` action or the `include` function](template_logic.md#include).

A helper template file may contain any number of helper templates, and must have the extension ".tpl".  Here is an example helper template file:

```sh
{{- define "my helper template" -}}
This is my helper template!  My username is {{ .values.username }}.
{{- end -}}


{{- define "my yaml colors helper" -}}
colors:
{{- range .values.colors }}
- {{ . }}
{{- end }}
{{- end -}}
```

These helpers may be used in other templates like this:

```yaml
message: {{ include "my helper template" . }}
list: {{- include "my yaml colors helper" . | trim | nindent 2 }}
```

## `dependencies/`

Dependency definitions are references to other template packages.  A package dependency must contain both the package information (i.e. package name and version) and the parameters to send to that package.

```yaml
# First, define the package reference
package:
  name: kpmtool/helloworld
  version: 1.0.0

# Next, specify the parameters to send to that package
parameters:
  myName: Rohit
```

Since dependency definitions are templates themselves, the parameters to send to the referenced package can be quite flexible.  The package reference itself can change based on the parameters provided to the parent template!

```yaml
package:
{{- if .values.sayHello }}
  name: kpmtool/helloworld
  version: 1.0.0
{{- else }}
  name: kpmtool/example
  version: {{ .values.exampleVersion }}
{{- end }}

parameters:
{{- if not .values.useDefaults }}
  colors: {{- .values.colors | toYaml | nindent 2 }}
{{- end }}
```
