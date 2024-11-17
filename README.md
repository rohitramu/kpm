# KPM

- [What is KPM?](#what-is-kpm)
- [Setup](#setup)
  - [Installation](#installation)
  - [Command line usage](#command-line-usage)
  - [Version information](#version-information)
  - [Golang templating](#golang-templating)
- [Template packages](#template-packages)
  - [List the locally available template packages](#list-the-locally-available-template-packages)
  - [Create a parameters file](#create-a-parameters-file)
  - [View a template package's default parameters file](#view-a-template-packages-default-parameters-file)
  - [Execute a template package](#execute-a-template-package)
- [Authoring a template package](#authoring-a-template-package)
  - [Directory structure](#directory-structure)
  - [`package.yaml`](#packageyaml)
  - [`parameters.yaml`](#parametersyaml)
  - [`interface.yaml`](#interfaceyaml)
  - [`templates/`](#templates)
  - [`helpers/`](#helpers)
  - [`dependencies/`](#dependencies)
  - [Template functions and logic](#template-functions-and-logic)
- [Testing your package locally](#testing-your-package-locally)
- [Pack your template package](#pack-your-template-package)
  - [Unpack a template package](#unpack-a-template-package)

## What is KPM?

KPM is a command line tool which simplifies and modularizes the process of generating text files.  It was initially developed to generate configuration files for Kubernetes as an alternative to Helm Charts, however it can be used in any situation that  requires text file generation.

## Setup

### Installation

First, choose and remember the appropriate operating system string, `${os}`, for your system:

- `darwin`
- `linux`
- `windows`

Also, choose and remember the appropriate architecture string, `${arch}`, for your system:

- `amd64`
- `arm64`
- `386`

These values should will be used to install the appropriate binary.

#### From the GitHub website

Download the KPM executable from the [Releases](https://github.com/rohitramu/kpm/releases) tab.  Add this executable to your PATH environment variable so it is available from anywhere on your machine.

#### Using `wget` on Linux or MacOS

```sh
wget -P /usr/bin -O kpm "https://github.com/rohitramu/kpm/releases/latest/kpm_${os}_${arch}"
```

#### Using PowerShell on Windows

```powershell
Invoke-WebRequest -OutFile "${Env:ProgramFiles}" "https://github.com/rohitramu/kpm/releases/latest/kpm_${os}_${arch}"
```

### Command line usage

To see the list of available subcommands:

```sh
kpm -h
```

To see the usage pattern for any subcommand, use the `-h` flag:

```sh
kpm new -h
```

### Version information

To get the version information for the KPM binary, run:

```sh
kpm version
```

### Golang templating

KPM uses Golang templating.  More information about defining and using templates can be found in the official [Golang template docs](https://golang.org/pkg/text/template/).

## Template packages

A template package is simply a collection of templates.  This collection of templates can be thought of as a program or function which can be executed.  A template package accepts inputs in the form of parameters, and produces outputs in the form of generated files.

Usage of a template package typically consists of these steps:

1. [Create a parameters file](#create-a-parameters-file) which can be used as input to the template package.
1. [Run the template package](#execute-a-template-package) with your parameters file.
1. View your generated files!

### List the locally available template packages

To see the list of all template packages in the local KPM repository, run the "list" subcommand:

```sh
kpm ls
```

### Create a parameters file

A parameters file is just a YAML file containing the input values to a template package.  The structure of this parameters file is determined by the [interface](#interfaceyaml) of the package that you would like to execute.

Looking at the [default parameters file](#parametersyaml) is a great way to understand the expected structure of your parameters file, because:

- All parameters must have default values.
- Package authors are encouraged to add documentation about their package in this file (as comments).

To see the default parameters file for a package, use the ["view" subcommand](#view-a-template-packages-default-parameters-file).

### View a template package's default parameters file

The default parameters file is provided by package authors to specify default values for parameters that the user's parameters file does not set.  The default parameters file is also used to document the template package.  Most importantly, explanations of how to set each parameter should be included in this file by the package author.

Users may find it useful to view the default parameters file in order to gain a better understanding of how the template package works.  The "view" command can print out the default parameters file for any template package which is available in your local KPM repository:

```sh
kpm view kpmtool/example -v 1.0.0
```

If a version is not specified, the highest available version which is in the local KPM repository (i.e. one that has already been [packed](#pack-your-template-package)) will be used.

### Execute a template package

Execute a template package with the "run" subcommand:

```sh
// Pull the package if you haven't already
kpm pull kpmtool/example -v 1.0.0

// Run the package with default parameters
kpm run kpmtool/example -v 1.0.0

// Run the package with custom parameters
kpm run kpmtool/example -v 1.0.0 -f my_params.yaml
```

If a version is not specified, the highest available version which is in the local KPM repository (i.e. one that has already been [packed](#pack-your-template-package)) will be used.

If an output directory is not specified with the `--output-dir` flag, files will be generated in `<current directory>/.kpm_generated/<output name>`.

If an output name is not specified with the `--output-name` flag, `<package name>-<package version>` will be used as the output name.

## Authoring a template package

### Directory structure

A template package must have the following directory structure, but not every file is required (items with trailing slashes are directories):

```txt
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

### `package.yaml`

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

NOTE: This file cannot be a template, and must only contain concrete values.

### `parameters.yaml`

A package author must provide a set of default parameters which will be used whenever a user does not provide a parameter.  The default parameters file is also a great place to document each parameter using comments.

NOTE: This file cannot be a template, and must only contain concrete values.

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

### `interface.yaml`

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

### `templates/`

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

### `helpers/`

Helper templates (a.k.a. "named templates" or "partial templates") allow the definition of templates which are useful in other templates in the package.  If you find yourself copying and pasting parts of templates, defining helper templates will allow you to simplify your templates and reduce the likelihood of copy-paste errors.

Helper templates can be inserted by using either the [`template` action or the `include` function](#include).

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

### `dependencies/`

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

### Template functions and logic

Template functions and control structures are used to transform data within templates, in order to produce the desired output.

#### Controlling whitespace

Controlling whitespace in certain types of files can be very important.

One example is YAML, where properties are defined on objects by indenting them.  In this case, producing incorrect indentation would lead to an incorrect definition of objects.

Thanks to [Golang templating](#golang-templating), whitespace may be controlled using special syntax when inserting placeholders in templates:

```yaml
# This will result in the string "hello world" being indented
    {{ "hello" }} {{ "world" }}

# This will result in the string "hello world" being left-justified (i.e. not indented at all)
    {{- "hello" }} {{ "world" }}

# This will result in the string "hello" being indented, and "world" being inserted without
# the spaces or new line after "hello" (i.e. it would become the string "    helloworld")
    {{ "hello" -}}
    {{ "world" }}
```

The dash (`-`) character after the start of the placeholder (`{{`) or before the end of the placeholder (`}}`) tells the renderer to remove any whitespace to the left or right before inserting the value.  This trims all whitespace, including **spaces, tabs and new lines**.

#### Conditionals

Conditional generation of output is done using if-else statements:

```sh
{{- if .values.show -}}
This text will only appear if the ".values.show" boolean property is set to true.
{{- end -}}
```

Multiple conditions can be checked, and you can also specify default behavior if all conditions fail:

```sh
{{- if (eq .values.color "green") -}}
Go!
{{- else if (eq .values.color "red") -}}
Stop!
{{- else -}}
Caution...
{{- end -}}
```

The "with" action is useful when you need to check for whether a property has been defined/supplied:

```sh
{{- with .values.myProperty }}
This will only appear if ".values.myProperty" has been set.  Its value is: {{ . }}
{{- end }}
```

The value retrieved by "with" may also be assigned to a variable:

```sh
{{- with $myProperty := .values.myProperty }}
Congratulations, you set your variable to "{{ $myProperty }}"!
{{- end }}
```

#### Loops

We can iterate over items in an array with the "range" action:

```sh
{{- range .values.myArray }}
- {{ . }}
{{- end }}
```

Each item in the array can be assigned to a variable as well.  This is useful when the items are objects rather than value types:

```sh
{{- range $myObject := .values.myObjectArray }}
- I have both {{ $myObject.property1 }} and {{ $myObject.property2 }}!
{{- end }}
```

#### Sprig functions

Sprig is a large library of useful template functions.  All of these functions are available for use inside all templates in a template package, including the [interface definition](#interfaceyaml), [helper templates](#helpers) and [dependency definitions](#dependencies).

- Documentation: <http://masterminds.github.io/sprig/>
- GitHub: <https://github.com/Masterminds/sprig>

#### Other template functions

##### `index`

It is very difficult to reference parameters which have special characters in their name, for example `my-property`.  The dash character (`-`) will break the "dot notation", so the following statement will fail:

```yaml
# This fails because dot notation does not handle special characters, e.g. dashes
{{ .values.my-property }}
```

The "index" function comes to the rescue!

```yaml
# We can reference any property, even if it has a weird name
{{ index .values "my-property" }}
```

The index function actually takes a list of keys, so this can be used to access a property at any depth:

```yaml
# We can reference properties at any depth
{{ index .values.myObject "my-weirdly-named-object" "inner-property" }}
```

The index function returns the property as-is, so you can continue to manipulate it afterwards:

```yaml
# We can continue manipulating the result object from the index function
{{ (index .values.myObject "my-weirdly-named-object") | toYaml | trim | indent 2 }}
```

##### `include`

The `template` action can be used for inserting helper templates.

```yaml
# This will execute and insert the helper template's output as-is, and won't let you further manipulate it
{{ template "my helper template" . }}
```

However, since it does not return the helper template's output as a string, transforming it before inserting it is not possible.  The `include` function solves this problem by executing the helper template and then returning the result as a string (rather than immediately inserting it in the document as-is).

```yaml
# This will execute the helper template, and then allow you to do things like trim whitespace and add indentation (in this case 4 spaces of indentation)
{{ include "my helper template" . | trim | indent 4 }}

# You can even define your helper template as a YAML object, and then convert it into an object inside the template!
{{ index (include "my helper template" . | toYaml) "myProperty" }}
```

##### `indent` vs. `nindent`

The `indent` function is used to indent all lines of a string by the given number of spaces.  This is very useful in files where whitespace and indenting is important, however it can lead to templates which are difficult to read, since the placeholder needs to be left-justified and placed on a new line:

```yaml
myObject:
  myNestedObject:
{{ include "my helper template" . | indent 4 }}
  anotherNestedObject:
    foo: bar
```

The [Sprig library](#sprig-functions) provides a `nindent` function which adds the new line for you before indenting the string.  This makes it much easier to read the template definition:

```yaml
myObject:
  myNestedObject: {{- include "my helper template" . | nindent 4 }}
  anotherNestedObject:
    foo: bar
```

## Testing your package locally

## Pack your template package

To make your package available to use locally, run the "pack" subcommand:

```sh
kpm pack /path/to/package/root
```

This adds the package to your local KPM repository.

If you are in the root directory of the package, you can just run the following (note the dot - it represents the current directory):

```sh
kpm pack .
```

The pack command will then add your package to the local KPM repository, using the name and version specified in the [package definition](#packageyaml) file.

It will be shown by the ["list" subcommand](#list-the-locally-available-template-packages), and can be executed with the ["run" subcommand](#execute-a-template-package):

```sh
kpm ls
kpm run username/my.package -v 0.1.0
```

### Unpack a template package

Unpacking (i.e. extracting) a template package to your file system can be useful to inspect its inner workings.

Unpack a template package by running the "unpack" subcommand:

```sh
kpm unpack kpmtool/example -v 1.0.0
```

The files in the package will be copied to the directory `<output directory>/<export name>`.

If an output directory is not provided with the `--output-dir` flag, `<current directory>/.kpm_exported` will be used.

If an export name is not provided with the `--export-name` flag, `<package name>-<package version>` will be used.

If a version is not specified, the highest available version which is in the local KPM repository (i.e. one that has already been [packed](#pack-your-template-package)) will be used.
