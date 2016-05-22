package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gorhill/cronexpr"
)

var version = "dev-build"

type TinyCronJob struct {
	cmd      string
	args     []string
	schedule *cronexpr.Expression
	opts     TinyCronOpts
}

type TinyCronOpts struct {
	debug  bool
	daemon bool
	jobLog string
}

// parseExpression parses a cron schedule and command from a single string
//func parseExpression(s string) (string, *cronexpr.Expression, error) {
//}

func NewTinyCronJob(s string) (*TinyCronJob, error) {
	var expr string
	var cmdline []string
	parts := strings.Split(s, " ")

	if strings.HasPrefix(s, "@") {
		expr = parts[0]
		cmdline = parts[1:]
	} else {
		if len(parts) <= 7 {
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
		opts:     optsFromEnv(),
	}
	return job, nil
}

// Run an exec job, returning when completed
func (job *TinyCronJob) run() {
	exe := exec.Command(job.cmd, job.args...)
	exe.Stdout = os.Stdout
	exe.Stderr = os.Stderr
	if job.opts.debug {
		output("running job: %s %s", job.cmd, strings.Join(job.args, " "))
	}
	errHandler(exe.Run(), "job failed")
}

func (job *TinyCronJob) nap() {
	now := time.Now()
	nextRun := job.schedule.Next(now)
	timeDelta := nextRun.Sub(now)
	if job.opts.debug {
		output(fmt.Sprintf("next job scheduled for %s", nextRun))
	}
	time.Sleep(timeDelta)
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

func optsFromEnv() (opts TinyCronOpts) {
	if os.Getenv("TINYCRON_DEBUG") != "" {
		opts.debug = true
	}
	if os.Getenv("TINYCRON_DAEMON") != "" {
		opts.daemon = true
	}
	if os.Getenv("TINYCRON_JOBLOG") != "" {
		opts.jobLog = os.Getenv("TINYCRON_JOBLOG")
	}
	return opts
}

func parseArgs(args []string) {
	// var opts []string
	fmt.Println(args)

	//	opts := make([]string, 0)
	for i, arg := range args {
		fmt.Printf("%s - %s\n", i, arg)
		//		if strings.HasPrefix(arg, "-") {
		//			opts = append(opts, arg)
		//		} else {
		//			fmt.Println(opts)
		//			fmt.Println(args[i:])
		//			break
		//		}
	}
}

func usage() {
	fmt.Println("Usage: tinycron [expression] [command]")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		errHandler(fmt.Errorf("incorrect number of arguments"), "")
		usage()
	}

	switch {
	case os.Args[1] == "version":
		fmt.Printf("tinycron version %s\n", version)
		os.Exit(0)
	case os.Args[1] == "help":
		usage()
		fmt.Println("Normal")
	case len(os.Args) <= 2:
	}

	job, err := NewTinyCronJob(os.Args[1])
	exitOnErr(err, "error creating job")

	for _, s := range os.Args[2:] {
		job.args = append(job.args, s)
	}

	for {
		job.nap()
		go job.run()
	}
}
