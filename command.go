// Copyright (c) 2023 BVK Chaitanya

// Package subcmd implements subcommand parsing mechanism for the command-line.
//
// Users can define commands and group them together into a parent subcommand,
// etc. to from subcommand hierarchies of arbitrary depths.
//
// Commands can define flags using `flag.FlagSet` objects from the Go standard
// library. Commands can also provide detailed documentation which is optional.
//
// A few special top-level commands "help", "flags", and "commands" are added
// automatically for documentation. More detailed documentation is collected
// through the optional `interface{ CommandHelp() string }` method on the
// Command objects.
//
// # EXAMPLE
//
//	type runCmd struct {
//		background  bool
//		port        int
//		ip          string
//		secretsPath string
//		dataDir     string
//	}
//
//	// Run implements the `main` method for "run" subcommand.
//	func (r *runCmd) Run(ctx context.Context, args []string) error {
//	 ...
//		if len(p.dataDir) == 0 {
//			p.dataDir = filepath.Join(os.Getenv("HOME"), ".data")
//		}
//		...
//		return nil
//	}
//
//	// Command implements the subcmd.Command interface.
//	func (r *runCmd) Command() (*flag.FlagSet, MainFunc) {
//		fset := flag.NewFlagSet("run", flag.ContinueOnError)
//		fset.BoolVar(&r.background, "background", false, "runs the daemon in background")
//		fset.IntVar(&r.port, "port", 10000, "TCP port number for the daemon")
//		fset.StringVar(&r.ip, "ip", "0.0.0.0", "TCP ip address for the daemon")
//		fset.StringVar(&r.secretsPath, "secrets-file", "", "path to credentials file")
//		fset.StringVar(&r.dataDir, "data-dir", "", "path to the data directory")
//	  return fset, MainFunc(r.Run)
//	}
package subcmd

import (
	"context"
	"flag"
	"os"
)

// MainFunc defines the signature for `main` function for a subcommand.
type MainFunc func(ctx context.Context, args []string) error

// Command interface defines the requirements for Command objects.
type Command interface {
	// Command returns the main function for a subcommand and it's command-line
	// flags. Returned `flag.FlagSet` must have a non-empty name which is taken
	// as the subcommand name.
	//
	// NOTE: This method is called just once per subcommand, so implementations
	// can return a new `flag.FlagSet` object.
	Command() (*flag.FlagSet, MainFunc)
}

// Group creates a parent command with the given subcommands nested under it's
// name.
func Group(name, description string, cmds ...Command) Command {
	return &cmdGroup{
		flags:    flag.NewFlagSet(name, flag.ContinueOnError),
		subcmds:  cmds,
		synopsis: description,
	}
}

// Run parses command-line arguments from `args` into flags and subcommands and
// selects the most appropriate subcommand to execute from `cmds`. Global
// command-line flags from `flag.CommandLine` are also processed on the way to
// resolving the final subcommand.
func Run(ctx context.Context, cmds []Command, args []string) error {
	if cmds == nil {
		return os.ErrInvalid
	}
	root := cmdGroup{
		flags:   flag.CommandLine,
		subcmds: cmds,
	}
	return root.run(ctx, args)
}
