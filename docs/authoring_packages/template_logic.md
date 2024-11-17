# Template functions and logic

There are many powerful functions and control statements that can be used to transform data within templates, in order to produce the desired output.  This page will give you an overview of how you can use them to create your own templates.

- [Golang templating](#golang-templating)
- [Controlling whitespace](#controlling-whitespace)
- [Conditionals](#conditionals)
- [Loops](#loops)
- [Sprig functions](#sprig-functions)
- [Commonly useful template functions](#commonly-useful-template-functions)
  - [`index`](#index)
  - [`include`](#include)
  - [`indent` vs. `nindent`](#indent-vs-nindent)

## Golang templating

KPM uses Golang templating.  More information about defining and using templates can be found in the official [Golang template docs](https://golang.org/pkg/text/template/).

## Controlling whitespace

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

## Conditionals

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

## Loops

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

## Sprig functions

Sprig is a large library of useful template functions.  Sprig functions are available for use inside all templates in a template package, including the [interface definition](package_files.md#interfaceyaml), [helper templates](package_files.md#helpers) and [dependency definitions](package_files.md#dependencies).

- Documentation: <http://masterminds.github.io/sprig/>
- GitHub: <https://github.com/Masterminds/sprig>

## Commonly useful template functions

### `index`

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

### `include`

The `template` action can be used for inserting helper templates.

```yaml
# This will execute and insert the helper template's output as-is, and won't let you further manipulate it
{{ template "my helper template" . }}
```

However, since it does not return the helper template's output as a string, transforming it before inserting it is not possible.  The `include` function solves this problem by executing the helper template and then returning the result as a string (rather than immediately inserting it in the document as-is).

```yaml
# This will execute the helper template, and then allow you to do things like trim whitespace and add indentation (in this case 4 spaces of indentation)
{{ include "my helper template" . | trim | indent 4 }}

# You can even define your helper template as YAML, and then convert it into an object!
{{ index (include "my helper template" . | fromYaml) "myPropertyInsideHelperTemplate" }}
```

### `indent` vs. `nindent`

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
