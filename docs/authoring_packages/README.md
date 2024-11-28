# Authoring a template package

## Directory structure

A template package must have the following directory structure, but not every file is required (a `*` indicates a required file). Note that items with trailing forward slashes (`/`) are directories.

```txt
<template_dir>
|
+-- package.yaml *
+-- interface.yaml *
+-- parameters.yaml *
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

## Testing your template package

### Pack your template package

To make your package available to use locally, run the "pack" subcommand:

```sh
kpm pack /path/to/package/root
```

This adds the package to your local KPM repository.

If you are in the root directory of the package, you can just run the following (note the dot - it represents the current directory):

```sh
kpm pack .
```

The pack command will then add your package to the local KPM repository, using the name and version specified in the [package definition](package_files.md#packageyaml) file.

It will be shown by the ["list" subcommand](../README.md#list-the-locally-available-template-packages), and can be executed with the ["run" subcommand](../README.md#execute-a-template-package):

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
