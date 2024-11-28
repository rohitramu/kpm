# KPM

KPM is a command line tool which modularizes the process of generating text files.  It was initially developed to generate configuration files for Kubernetes as a more robust alternative to Helm Charts, however it can be used in any situation that requires text file generation.

- [Template packages](#template-packages)
- [Setup](#setup)
  - [Installation](#installation)
  - [Command line usage](#command-line-usage)
  - [Version information](#version-information)

## Template packages

KPM uses the concept of "template packages" to generate files.  A template package is simply a collection of individual text file templates.  This collection of templates can be thought of as a program which can be executed.  A template package accepts inputs in the form of parameters, and produces outputs in the form of generated files.

See [the docs](docs/README.md) for more information about how to [use template packages](docs/using_packages/README.md) and [author your own template packages](docs/authoring_packages/README.md).  There are also [examples](docs/examples/README.md) that you can use as a starting point for authoring your own template packages.

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

These values will be used to install the appropriate binary.

#### From the GitHub website

Download the KPM executable from the [Releases](https://github.com/rohitramu/kpm/releases) tab.  Put this executable somewhere on your PATH so it is available from anywhere on your machine.

#### Linux

##### Add `~/bin` to your PATH environment variable

If `~/bin` is not already on your PATH, add it:

```sh
export PATH="$PATH:~/bin"
```

To make this permanent, add it to your `~/.bashrc` file:

```sh
echo "export PATH=\"\$PATH:~/bin\"" >> ~/.bashrc
```

##### Download and install KPM

Note that the file needs to be made executable after download.

```sh
wget -O ~/bin/kpm "https://github.com/rohitramu/kpm/releases/latest/download/kpm_${os}_${arch}"
chmod 777 ~/bin/kpm
```

For example, if `os=linux` and `arch=amd64`:

```sh
wget -O ~/bin/kpm "https://github.com/rohitramu/kpm/releases/latest/download/kpm_linux_amd64"
chmod 777 ~/bin/kpm
```

#### Using PowerShell on Windows

```powershell
Invoke-WebRequest -OutFile "${Env:ProgramFiles}/kpm.exe" "https://github.com/rohitramu/kpm/releases/latest/download/kpm_${os}_${arch}.exe"
```

For example, if `os=windows` and `arch=amd64`:

```powershell
Invoke-WebRequest -OutFile "${Env:ProgramFiles}/kpm.exe" "https://github.com/rohitramu/kpm/releases/latest/download/kpm_windows_amd64.exe"
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
