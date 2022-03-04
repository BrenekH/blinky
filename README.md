# Blinky

Simple, all in one Pacman repository hosting server software.

<!-- TODO: Insert more badges here -->
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/BrenekH/blinky)
[![Go Report Card](https://goreportcard.com/badge/github.com/BrenekH/blinky)](https://goreportcard.com/report/github.com/BrenekH/blinky)
[![Blinky CI/CD](https://github.com/BrenekH/blinky/actions/workflows/blinky-ci-cd.yaml/badge.svg)](https://github.com/BrenekH/blinky/actions/workflows/blinky-ci-cd.yaml)

## Document Conventions

This Git repository contains both the Blinky server and the `blinky` CLI tool.
In order to keep track of which one is being talked about, the server will be referred to as either Blinky, the server, or `blinkyd`, depending on the context.

Likewise, the CLI tool will be referred to as the CLI, or as simply `blinky`.

## Installation

### Source

To install the Blinky components from source, the [Go compiler](https://go.dev/dl) needs to be installed.

#### Server

`go install github.com/BrenekH/blinky/cmd/blinkyd@latest`

#### CLI

`go install github.com/BrenekH/blinky/cmd/blinky@latest`

<!-- ?Perhaps talk about installing shell completions? -->

### Package Managers

Currently, Blinky is not available in any package managers.

## Usage

The following usage instructions are just the basics of how to use Blinky.
For a full run down, please visit the [Blinky wiki](https://github.com/BrenekH/blinky/wiki).

### Server

`blinkyd` supports 3 methods of configuration, command-line arguments, environment variables, and a TOML config file located at either `/etc/blinky/config.toml` or `$XDG_CONFIG_HOME/blinky/config.toml`.
The names are the same across each method, but with differing capitalization and word separators.
This document will only use command-line args, but a full table can be found in the [wiki](https://github.com/BrenekH/blinky/wiki).

#### Setting up repository paths

Blinky uses the same syntax as the Linux PATH variable to specify where the repository's files should be stored.
The required folders will be created if they do not exist.
The last element of the path is used as the name of the repo.
For example, the following command will create 3 repos: `repo1`, `repo2`, and `repo3`.

`blinkyd --repo-path "/mnt/storage1/repo1:/mnt/storage1/repo2:/opt/repo3"`

#### Change the port

Using `--http-port` the port Blinky uses for the HTTP server can be changed.

#### Protecting the API

A username and password can be set so that only those who know the username and password can manage packages.

The username is set by using `--api-uname` and the password is set with `--api-passwd`.

#### Getting help from the terminal

Using `blinkyd --help` will output a usage text and exit the application.

### CLI

`blinky` has 4 basic commands: `login`, `logout`, `upload`, and `remove`.

`login` is used to save the login credentials for a server.

`logout` removes the saved credentials for a server.

`upload` uploads new/updated packages to a repository hosted on a server.

`remove` deletes packages from a repository on a server.

**Examples:**

```text
$ blinky login --default https://blinky.example.com
Username: user
Password:

$ blinky upload custom_repo my_package.pkg.tar.zst
...

$ blinky remove custom_repo my_package
...

$ blinky logout https://blinky.example.com
...
```

More detailed usage instructions can be found by running `blinky --help`.

## Recommended Security Practices

### HTTPS

Because Blinky uses the HTTP Basic Authentication standard, it is recommended to use an HTTPS certificate with a reverse proxy to encrypt the conversation between Blinky and clients.

### Passwords in Shell History

While `blinky` does support providing a password as a CLI flag, it is not a recommended method to do so.

Instead, the interactive prompt should be used for entering passwords where possible.

## Known Issues

### Only x86_64 is supported

Blinky is built with only the x86_64 architecture in mind, since that is all that Arch Linux officially supports.
If you require support for other architectures, please reach out either in the [issues](https://github.com/BrenekH/blinky/issues) or [discussions](https://github.com/BrenekH/blinky/discussions) and we can figure out how to make it work.

## Frequently Asked Questions

### What is wrong with `repo-add`?

`repo-add` requires direct access to the directory that is being used to host the database and all of the package files.
While this is fine for a local repository, building a hosted version that can be updated from anywhere starts to become more complicated than using an API.
Instead of opening SSH, setting up Samba, or configuring a VPN, one can use the HTTP setup already required for both managing and consuming the repository.

### Why is it called Blinky?

According to Pac-Man lore, each of the four ghosts from the original game has a unique name.
The American versions of these names are Inky, Blinky, Pinky, and Clyde and because servers often have blinky lights, Blinky was chosen as the name.

## License

Blinky is licensed under the GNU General Public License version 3, a copy of which can be found in the [LICENSE](LICENSE) file.
