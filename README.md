# Blinky

Simple, all in one Pacman repository hosting server software.

<!-- TODO: Insert badges here -->

## Conventions

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
<!-- TODO -->

### Server

<!-- Repo file setup -->
<!-- Env var/command line -->

### CLI

<!-- Basic usage -->

More detailed usage instructions can be found by running `blinky help`.

## Security

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
