# Using template packages

- [What is a template package?](#what-is-a-template-package)
- [List the locally available template packages](#list-the-locally-available-template-packages)
- [Create a parameters file](#create-a-parameters-file)
- [The `view` subcommand](#the-view-subcommand)
- [Execute a template package](#execute-a-template-package)

## What is a template package?

A template package is simply a collection of templates.  This collection of templates can be thought of as a program or function which can be executed.  A template package accepts inputs in the form of parameters, and produces outputs in the form of generated files.

Usage of a template package typically consists of these steps:

1. [Create a parameters file](#create-a-parameters-file) which can be used as input to the template package.
1. [Run the template package](#execute-a-template-package) with your parameters file.
1. View your generated files!

## List the locally available template packages

To see the list of all template packages in the local KPM repository, run the "list" subcommand:

```sh
kpm ls
```

## Create a parameters file

A parameters file is just a YAML file containing the input values to a template package.  The structure of this parameters file is determined by the [interface](../authoring_packages/package_files.md#interfaceyaml) of the package that you would like to execute.

Looking at the [default parameters file](../authoring_packages/package_files.md#parametersyaml) is a great way to understand the expected structure of your parameters file, because:

- All parameters must have default values.
- Package authors are encouraged to add documentation about their package in this file (as comments).

To see the default parameters file for a package, use the ["view" subcommand](#the-view-subcommand).

## The `view` subcommand

The default parameters file is provided by package authors to specify default values for parameters that the user's parameters file does not set.  The default parameters file is also used to document the template package.  Most importantly, explanations of how to set each parameter should be included in this file by the package author.

Users may find it useful to view the default parameters file in order to gain a better understanding of how the template package works.  The "view" command can print out the default parameters file for any template package which is available in your local KPM repository:

```sh
kpm view kpmtool/example -v 1.0.0
```

If a version is not specified, the highest available version which is in the local KPM repository (i.e. one that has already been [packed](../authoring_packages/README.md#pack-your-template-package)) will be used.

## Execute a template package

Execute a template package with the "run" subcommand:

```sh
// Pull the package if you haven't already
kpm pull kpmtool/example -v 1.0.0

// Run the package with default parameters
kpm run kpmtool/example -v 1.0.0

// Run the package with custom parameters
kpm run kpmtool/example -v 1.0.0 -f my_params.yaml
```

If a version is not specified, the highest available version which is in the local KPM repository (i.e. one that has already been [packed](../authoring_packages/README.md#pack-your-template-package)) will be used.

If an output directory is not specified with the `--output-dir` flag, files will be generated in `<current directory>/.kpm_generated/<output name>`.

If an output name is not specified with the `--output-name` flag, `<package name>-<package version>` will be used as the output name.
