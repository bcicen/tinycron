package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/gorhill/cronexpr"
)

var version = "dev-build"

func output(s string, vars ...string) {
	msg := fmt.Sprintf(s, vars)
	fmt.Printf("[%s] %s\n", color.CyanString("tinycron"), msg)
}

func errHandler(err error, msg string) {
	if err != nil {
		if msg == "" {
			output(err.Error())
		} else {
			output(fmt.Sprintf("%s %s", color.RedString(msg), err.Error()))
		}
	}
}

func exitOnErr(err error, msg string) {
	if err != nil {
		errHandler(err, msg)
		os.Exit(1)
	}
}

type execCmd struct {
	cmd  string
	args []string
}

// Run an exec job, returning when completed
func runJob(cmdline []string) {
	cmd, args := cmdline[0], cmdline[1:]
	job := exec.Command(cmd, args...)
	job.Stdout = os.Stdout
	job.Stderr = os.Stderr
	output("running job: %s", cmd)
	errHandler(job.Run(), "job failed")
}

func nap(schedule *cronexpr.Expression) {
	now := time.Now()
	nextRun := schedule.Next(now)
	timeDelta := nextRun.Sub(now)
	output(fmt.Sprintf("next job scheduled for %s", nextRun))
	time.Sleep(timeDelta)
}

func main() {
	app := cli.NewApp()
	app.Name = "tinycron"
	app.Usage = "a very small replacement for cron"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "foreground, f",
			Usage: "Keep tinycron running in the foreground",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debug output",
		},
	}

	app.Action = func(c *cli.Context) {
		if len(c.Args()) < 2 {
			cli.ShowAppHelp(c)
			errHandler(fmt.Errorf("incorrect number of arguments"), "")
			os.Exit(1)
		}
		schedule, err := cronexpr.Parse(c.Args()[0])
		exitOnErr(err, "erroring parsing schedule")

		cmdline := c.Args()[1:]

		for {
			nap(schedule)
			go runJob(cmdline)
		}
	}

	app.Run(os.Args)
}
