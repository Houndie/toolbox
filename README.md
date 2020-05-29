Toolbox
=======

**A tool vendoring assistant**

Why Toolbox?
------------

Have you ever used a code generation tool, only for git to tell you that you just changed 74 files in your repository?  It's likely that you and your teammates are using a different version of the tool!

When multiple people use the same code base, you'll soon find that it's exceedingly useful for your whole team to have the same versions of the same tools.  This helps get your team up and running quickly, eliminates meaning churn when generating files, and prevents errors when behavior changes between versions.

Enter Toolbox.

Toolbox is a small tool meant to assist you with installing tools for your project.  The go team's recommendation is to use go modules to vendor your tools, and, in fact, this is what toolbox does under the hood.  However, using go to vendor your tools can be fiddly and prone to user-error.  Toolbox leverages the `go` binary (and optionally, `goimports`) to automate the most common use cases for you.

Instillation
-------------

`go get github.com/houndie/toolbox`

How do I use it?
----------------

Toolbox has five main commands:

* `$ toolbox add <toolname> [version]` Downloads and installs `<toolname>`.  A version may optionally be provided, otherwise it will attempt to find the latest version.  This command is also used to upgrade/downgrade tools.
* `$ toolbox remove <toolname>` Removes `<toolname>` from being tracked by your project.  Also attempts to uninstall the installed binary.
* `$ toolbox sync` Downloads and installs any missing tools.  If a tool is up-to-date, no action is taken for that tool.
* `$ toolbox do -- <command>` Runs a command in an environment where the tools managed by toolbox are available.  The dash (`--`) is optional, and is used to denote that flags will belong to the subcommand (otherwise, toolbox will attempt to parse flags itself).
* `$ toolbox list` Lists all saved tools and their options

Example
-------

```
# Update your repository, and get any tool changes
git pull
toolbox sync

# Add the stringer tool to the project
toolbox add golang.org/x/tools/cmd/stringer

# Use the stringer tool
toolbox do -- stringer -type=Pill

# Change stringer to a different version
toolbox add golang.org/x/tools/cmd/stringer v0.4

# Remove the tool from the project
toolbox remove golang.org/x/tools/cmd/stringer
```

Toolbox creates one file and one directory, the names of which default to `tools.go` and `_tools` respectively.  `tools.go` contains a list of the tools managed by toolbox, and should be checked into source control.  `_tools` contains the downloaded tool binaries, and does not need to be checked into source control.

How does it work?
-----------------

*Warning, boring implementation details ahead*

In short, this tool automates the best practices put forth by the go team [here](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module) and [here](https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md).

Whenever you add a tool to your repository, toolbox calls `go get` with a modified `GOBIN` to download the tool into the `_tools` directory.  A project-local directory is used so that the versions of tools that affect one project, do not affect another project.  `go get` automatically records the version of the tool being used in your `go.mod` file, which is used later when syncing tools.  Toolbox also adds your tool to the import list in `tools.go`, the whole purpose of which is to store a list of tools so that toolbox and go modules knows what tools you're using.  The tools are recorded in an import list in a `.go` file so that go modules knows that this a dependency of your project, and won't remove them from `go.mod` on your next `go mod tidy`.  `tools.go` contains a build tag so that it is never actually built into your project.  If goimports is installed, it will also be called on the `tools.go` file, so that you don't accidentally automatically change it just by viewing it in a text editor.

Whenever you remove a tool from your repository, toolbox rewrites `tools.go` to no longer reference your tool (as before, `goimports` is used if available).  It then calls `go mod tidy` to remove references to it in `go.mod` as well.  Toolbox also attempts to find and delete the tool from the `_tools` directory.

Calling `toolbox sync` causes toolbox to loop through all known tools in `tools.go` and call `go install` for each one of them.  `go install` is preferred here, as it does not take a version argument, and instead references `go.mod` to see what version to install.  `go install` also differs from `go get` in that it does not call out to the internet if your tool is up-to-date.

Finally `toolbox do` simply edits your system `PATH` environment variable to include the `_tools` directory.  It also sets `GOBIN` so that if your tool installs more tools, they will also end up in the `_tools` folder for this project. It then runs the passed command in this new environment.

Configuration
-------------

The executables used for `go` and `goimports` can both be specified with flags.  Similarly, you can also customize the names for the `tools.go` and `_tools` folder.

`toolbox` is built on [viper](https://github.com/spf13/viper) and [cobra](https://github.com/spf13/cobra), and therefore most commandline flags can also be configured in a configuration file as well. By default, toolbox looks for the configuration file `.toolbox` with either an `ini`, `json`, `yaml`, or `toml` extension.  You can also specify a configuration file on the commandline.

Library
-------

Toolbox can also be used as a library, with import path `github.com/Houndie/toolbox/pkg/toolbox`.  This is useful if you want to use toolbox in a [magefile](https://github.com/magefile/mage), or another go scripting system.

Thanks
------

* Thanks for the go team for building a great tooling system making this much less work than it could've been
* Thanks to [viper](https://github.com/spf13/viper) and [cobra](https://github.com/spf13/cobra) for making a great framework for creating executables.
* Finally, this project is HEAVILY inspired by [retool](https://github.com/twitchtv/retool) which did tool vendoring in a pre-module world.  If you're also working in a pre-module world, definitely check out retool as a vendoring system.
