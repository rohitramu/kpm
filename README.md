# Overview
KPM is a command line tool which attempts to simplify and modularize the process of generating text files.

# Setup
## Prerequisites
In order to [push](#push-package) and [pull](#pull-package) template packages, you must have installed and configured [Docker](https://docs.docker.com/install/).

The default Docker registry used is `docker.io` (i.e. Docker Hub).  Ensure that you have a Docker Hub account and have credentials by running `docker login docker.io`.

To specify a different Docker registry, set the `--docker-registry` flag on subcommands which interact with a Docker registry (e.g. [`push`](#push-package), [`pull`](#pull-package), [`run`](#run-package)).  E.g. to log in and run a package from your own Docker registry:
```
docker login my.registry.com
kpm run mynamespace/mypackage -v 2.0.0 --docker-registry my.registry.com
```

## Installation
Download the KPM executable from the [Releases](https://github.com/rohitramu/kpm/releases) tab.

## Command line usage
To see the list of available commands:
```
kpm help
```

To see the usage pattern for any subcommand, use the `-h` flag:
```
kpm run -h
```


# Template packages
A template package is simply a collection of templates.  This collection of templates can be thought of as a program or function which can be executed.  A template package accepts inputs in the form of parameters, and produces outputs in the form of generated files.

Usage of a template package typically consists of these steps:
 1. [Pull a template package](#pull-package) from a Docker registry to make it available locally.
 2. [Create a parameters file](#user-parameters) which can be used as input to the template package.
 3. [Run the template package](#run-package) with your parameters file.
 4. View your generated files!

KPM uses Golang templating.  More information about defining templates can be found here: https://golang.org/pkg/text/template/

## <a name="pull-package"></a>Pull a template package from a Docker registry
Pull a template package from a Docker registry with the "pull" subcommand.  For example, try running the following command:
```
kpm pull kpmtool/example -v 1.0.0
```

## <a name="list-packages"></a>List the locally available template packages
To see the list of all template packages in the local KPM repository, run the "list" subcommand:
```
kpm ls
```

## <a name="user-parameters"></a>Create a parameters file
A parameters file is just a YAML file containing the input values to a template package.  The structure of this parameters file is determined by the author of the package that you would like to execute.  Take a look at the "parameters.yaml" file in the root of the package directory for a good example of the structure that the package expects.  Package authors are also encouraged to add documentation by way of comments into this file in order to assist with understanding correct usage of the package.

## <a name="unpack-package"></a>Unpack a template package
Unpack a template package by running the "unpack" subcommand:
```
kpm unpack rohitramu/kpm.example /path/to/output/folder
```

If a version is not specified, the highest available version which is in the local KPM repository (i.e. one that has already been [pulled](#pull-package) or [packed](#pack-package)) will be used.

If an output directory is not specified, files will be copied to `<current directory>/.kpm_exported/<package full name>`.

## <a name="run-package"></a>Execute a template package
Execute a template package with the "run" subcommand:
```
kpm run rohitramu/kpm.example
```

If a version is not specified like in the above example, the highest available version which is in the local KPM repository (i.e. one that has already been [pulled](#pull-package) or [packed](#pack-package)) will be used.

If a version is specified and the package cannot be found in the local KPM repository, an attempt will be made to pull the package from the remote Docker registry.  By default, this is `docker.io`.  A different Docker registry can be specified setting the `--docker-registry` flag.

If an output directory is not specified, files will be generated in `<current directory>/.kpm_generated/<output name>`.

If an output name is not specified, the package's full name will be used as the output name.


# <a name="authoring"></a>Authoring a template package

## Directory structure
A template package must have the following directory structure (items with trailing slashes are directories):
```
<working dir>
|
+-- package.yaml
+-- interface.yaml
+-- parameters.yaml
|
+-- templates/
|  |
|  +-- my_template.txt
|  +-- other template.text
|  \-- someConfig.yaml
|
+-- helpers/
|  |
|  +-- my_helpers.tpl
|  +-- other helpers.tpl
|  \-- extraHelpers.tpl
|
+-- dependencies/
|  |
|  +-- dependency1.yaml
|  +-- dependency_2.yaml
|  \-- other-dependency.yaml
\
```

## package.yaml
The package information file defines the name and version of a package.  For example, a `package.yaml` file for version `1.0.0` of the package `kpmtool/helloworld` would look like this:
```yaml
name: kpmtool/helloworld
version: 1.0.0
```

The name of the package should include the namespace that will be used when pushing it to a Docker registry.  It may only contain lowercase letters, numbers, forward slashes and dots. Also, it must start with a lowercase letter.

The version must be in the format `"major.minor.revision"`.  Leading zeros are not permitted in the `major`, `minor` or `revision` segments of the version string, however a segment may be just `0` (zero).  The zero version (`0.0.0`) is not allowed.

## <a name="default-parameters"></a>parameters.yaml
A package author must provide a set of default parameters which will be used whenever a user does not provide a parameter.  The default parameters file is also a great place to document each parameter using comments.

NOTE: This is the only file in the package which cannot be a template, and must only contain concrete values.

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

## <a name="interface"></a>interface.yaml
The interface is a YAML template which defines what parameters the package requires in order to correctly generate output.  Parameters which are provided by the user are used as the input to this interface.  The resulting YAML is then used as the input to all other templates in the package.

If a value is not provided by the user for a parameter, the default value will be used.  Default values are defined in the [parameters](default-parameters) file.

The values in the interface may be defined as "constants" (i.e. hardcoded values) which can be referenced in templates.  However, for more complex string values (e.g. multi-line strings, strings with special characters, etc.), [helper templates](#helpers) should be preferred.  See the `constantGreeting` property or the `favorite-things` list below for examples of how hardcoded values may be defined.

Here is an example interface definition which can accept the parameters from the [example above](#default-parameters):
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
- Software!
isCustom: {{ .isCustom }}
customObj: {{- index . "custom-object" | toYaml | nindent 2 }}
```

Since the interface definition is itself a template, all of the normal template functions are available to be used.  See the [templates](#templates) section for more details.

## <a name="templates"></a>templates/
Files in the templates directory are the text templates which will be used to generate the output files.  These can be used to generate any text format with any filename.

Here is an example of a template file which generates a text file by using values provided by the [interface example](#interface) above:
```
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

## <a name="helpers"></a>helpers/
Helper templates (a.k.a. "named templates" or "partial templates") allow the definition of templates which are useful in other templates in the package.  If you find yourself copying any pasting parts of templates, defining helper templates will allow you to simplify your templates and reduce the likelihood of copy-paste errors.

Helper templates can be inserted by using either the [`template` action or `include` function](#include-function).

A helper template file may contain any number of helper templates, and must have the extension ".tpl".  Here is an example helper template file:
```
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

## <a name="dependencies"></a>dependencies/
Dependency definitions are references to other template packages.  A package dependency must contain both the package information (i.e. package name and version) and the parameters to send to that package.

```yaml
# First, define the package reference
package:
  name: kpmtool/helloworld
  version: 1.0.0

# Next, specify the parameters to send to that package
parameters:
  myName: World
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

## <a name="template-functions-and-logic"></a>Template functions and logic
Template functions and control structures are used to transform data within templates, in order to produce the desired output.

### Conditionals
Conditional generation of output is done using if-else statements:
```
{{- if .values.show -}}
This text will only appear if the ".values.show" boolean property is set to true.
{{- end -}}
```

Multiple conditions can also be checked, and you can also specify default behavior if all conditions fail:
```
{{- if (eq .values.color "green") -}}
Go!
{{- else if (eq .values.color "red") -}}
Stop!
{{- else -}}
Caution...
{{- end -}}
```

The "with" action is useful when you need to check for whether a property has been defined/supplied:
```
{{- with .values.myProperty }}
This will only appear if ".values.myProperty" has been set.  Its value is: {{ . }}
{{- end }}
```

The value retrieved by "with" may also be assigned to a variable:
```
{{- with $myProperty := .values.myProperty }}
Congratulations, you set your variable to "{{ $myProperty }}"!
{{- end }}
```

### Loops
We can iterate over items in an array with the "range" action:
```
{{- range .values.myArray }}
- {{ . }}
{{- end }}
```

Each item in the array can be assigned to a variable as well.  This is useful when the items are objects rather than value types:
```
{{- range $myObject := .values.myObjectArray }}
- I have both {{ $myObject.property1 }} and {{ $myObject.property2 }}!
{{- end }}
```

### Sprig functions
Sprig is a large library of useful template functions.  All of these functions are available for use inside all templates in a template package, including the [interface definition](#interface), [helper templates](#helpers) and [dependency definitions](#dependencies)).
- Documentation: http://masterminds.github.io/sprig/
- GitHub: https://github.com/Masterminds/sprig

### Other template functions
#### index
It is very difficult to reference parameters which have special characters in their name, for example `my-property`.  The dash character (`-`) will break the "dot" notation, so the following statement will fail:
```yaml
# This fails because dot notation does not handle special characters like dashes
{{ .values.my-property }}
```

#### <a name="include-function"></a>include
The `template` action can be used for inserting helper templates.  However, since it does not return the helper template as a string, transforming it before inserting it is not possible.  The `include` function solves this problem by executing the helper template and then returning it as a string (rather than immediately printing it in the document as-is).

```yaml
# This will insert the helper template as-is, and won't let you run the output through string functions
{{- template "my helper template" . }}

# This will execute the helper template, and then allow you to do things like trim whitespace and add indentation (in this case 4 spaces of indentation)
{{- include "my helper template" . | trim | indent 4 }}

# You can even define your helper template as a YAML object, and then convert it into an object inside the template!
{{- index (include "my helper template" . | toYaml) "myProperty" }}
```

#### indent vs. nindent
The `indent` function is used to indent all lines of a string by the given number of spaces.  This is very useful in files where whitespace and indenting is important, however it can lead to templates which are difficult to read, since the placeholder needs to be left-justified and placed on a new line:
```yaml
myObject:
  myNestedObject:
{{ include "my helper template" . | toYaml | indent 4 }}
  anotherNestedObject:
    foo: bar
```

The Sprig library provides a `nindent` function which adds the new line for you before indenting the string.  This makes it much easier to read the template definition:
```yaml
myObject:
  myNestedObject: {{- include "my helper template" . | toYaml | nindent 4 }}
  anotherNestedObject:
    foo: bar
```

# Testing your package locally

## <a name="pack-package"></a>Pack your template package
To make your package available to use locally, run the "pack" subcommand:
```
kpm pack /path/to/package/root
```

If you are in the root directory of the package, you can just run the following (note the dot):
```
kpm pack .
```

It can then be executed as normal with the ["run" command](#run-package), and will be shown by the ["list" command](list-packages).


# Sharing your template package

## <a name="push-package"></a>Push your template package to a Docker registry
1. If you have not already done so, make your package locally available by running the ["pack" command](#pack-package).  Verify that it really is locally available with the ["list" command](#list-packages).

2. To push your package to a Docker registry, ensure that you have the credentials to do so by first logging in.  For example, to get credentials for your Docker Hub account, run:
```
docker login docker.io
```

3. Run the "push" subcommand to push the template package to the Docker registry:
```
kpm push rohitramu/kpm.example
```
