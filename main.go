package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/gorhill/cronexpr"
)

var version = "dev-build"

type TinyCronJob struct {
	cmd      string
	args     []string
	schedule *cronexpr.Expression
	debug    bool
}

func output(msg string, vars ...interface{}) {
	if len(vars) > 0 {
		msg = fmt.Sprintf(msg, vars...)
	}
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

// Run an exec job, returning when completed
func (job *TinyCronJob) run() {
	exe := exec.Command(job.cmd, job.args...)
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr
	if job.debug {
		output("running job: %s %s", job.cmd, strings.Join(job.args, " "))
	}
	errHandler(exe.Run(), "job failed")
}

func (job *TinyCronJob) nap() {
	now := time.Now()
	nextRun := job.schedule.Next(now)
	timeDelta := nextRun.Sub(now)
	if job.debug {
		output(fmt.Sprintf("next job scheduled for %s", nextRun))
	}
	time.Sleep(timeDelta)
}

func NewTinyCronJob(s string) (*TinyCronJob, error) {
	var expr string
	var cmdline []string
	parts := strings.Split(s, " ")

	if strings.HasPrefix(s, "@") {
		expr = parts[0]
		cmdline = parts[1:]
	} else {
		if len(parts) < 8 {
			return nil, fmt.Errorf("incomplete cron expression")
		}
		expr = strings.Join(parts[0:7], " ")
		cmdline = parts[7:]
	}
	schedule, err := cronexpr.Parse(expr)
	if err != nil {
		return nil, err
	}
	job := &TinyCronJob{
		cmd:      cmdline[0],
		args:     cmdline[1:],
		schedule: schedule,
	}
	return job, nil
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

		job, err := NewTinyCronJob(c.Args()[0])
		exitOnErr(err, "error creating job")

		for _, s := range c.Args()[1:] {
			job.args = append(job.args, s)
		}

		job.debug = c.Bool("debug")

		for {
			job.nap()
			go job.run()
		}
	}

	app.Run(os.Args)
}
