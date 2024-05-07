// Copyright (c) 2023 BVK Chaitanya

package subcmd

import (
	"context"
	"flag"
	"log"
	"testing"
)

type TestCmd struct {
	name  string
	flags *flag.FlagSet
	args  []string
}

func newTestCmd(name string) *TestCmd {
	return &TestCmd{
		name:  name,
		flags: flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

func (t *TestCmd) Command() (*flag.FlagSet, MainFunc) {
	return t.flags, MainFunc(func(_ context.Context, args []string) error {
		log.Println("running", t.name, "with args", args)
		t.args = args
		return nil
	})
}

func (t *TestCmd) CommandHelp() string {
	return `First line of help output is used as synopsis.
Rest of the text is displayed as documentation for the command.
`
}

func TestRun(t *testing.T) {
	ctx := context.Background()

	run := newTestCmd("run")
	background := run.flags.Bool("background", false, "set to run in background")

	jobsList := newTestCmd("list")
	jobsList.flags.String("format", "json", "list output format")
	jobsSummary := newTestCmd("summary")
	jobsSummary.flags.String("format", "json", "summary output format")
	jobs := Group("jobs", "manage jobs", jobsList, jobsSummary)

	jobPause := newTestCmd("pause")
	jobPause.flags.Duration("timeout", 0, "pause duration")
	jobResume := newTestCmd("resume")
	jobResume.flags.Duration("timeout", 0, "resume duration")
	jobCancel := newTestCmd("cancel")
	jobCancel.flags.Duration("after", 0, "cancellation delay")
	jobArchive := newTestCmd("archive")
	jobDelete := newTestCmd("delete")
	job := Group("job", "manage single job", jobPause, jobResume, jobCancel, jobArchive, jobDelete)

	dbGet := newTestCmd("get")
	dbSet := newTestCmd("set")
	dbDelete := newTestCmd("delete")
	dbScan := newTestCmd("scan")
	dbBackup := newTestCmd("backup")
	db := Group("db", "manage database", dbGet, dbSet, dbDelete, dbScan, dbBackup)

	cmds := []Command{run, jobs, job, db}

	{
		args := []string{"db", "scan", "db-scan-argument"}
		if err := Run(ctx, cmds, args); err != nil {
			t.Fatal(err)
		}
		if len(dbScan.args) != 1 || dbScan.args[0] != "db-scan-argument" {
			t.Fatalf("want `db-scan-argument`, got %v", dbScan.args)
		}
	}

	{
		args := []string{"run", "-background", "run-argument"}
		if err := Run(ctx, cmds, args); err != nil {
			t.Fatal(err)
		}
		if len(run.args) != 1 || run.args[0] != "run-argument" {
			t.Fatalf("want `run-argument`, got %v", run.args)
		}
		if *background == false {
			t.Fatalf("want true, got false")
		}
	}

	{
		args := []string{"-h"}
		if err := Run(ctx, cmds, args); err != nil {
			t.Fatal(err)
		}
	}

	{
		args := []string{"run", "-h"}
		if err := Run(ctx, cmds, args); err != nil {
			t.Fatal(err)
		}
	}
}
