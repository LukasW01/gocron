package cron

import (
	"gocron/src/util"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

type Status struct {
	ExitStatus int
	Stdout     string
	Stderr     string
	Time       string
	Pid        int
}

type Process struct {
	Running  map[string]*Status
	Status   *Status
	Schedule string
}

var Proc Process

func Execute(command string, args []string) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	run := new(Status)
	run.Time = time.Now().Format(time.RFC3339)
	run.Pid = cmd.Process.Pid
	Proc.Running[strconv.Itoa(run.Pid)] = run
	go util.Output(&run.Stdout, stdout, run.Pid)
	go util.Output(&run.Stderr, stderr, run.Pid)

	if err := cmd.Wait(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// so set the error code to tremporary value
			run.ExitStatus = 127
			if status, ok := exit.Sys().(syscall.WaitStatus); ok {
				run.ExitStatus = status.ExitStatus()
				log.Printf("%d Exit Status: %d", run.Pid, run.ExitStatus)
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	delete(Proc.Running, strconv.Itoa(run.Pid))
	Proc.Status = run

	log.Println(run.Pid, "cmd:", command, strings.Join(args, " "))
}

func Create(schedule string, command string, args []string) (cr *cron.Cron, wgr *sync.WaitGroup) {
	Proc = Process{map[string]*Status{}, &Status{}, schedule}

	wg := &sync.WaitGroup{}
	c := cron.New(
		cron.WithParser(
			cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
			),
		),
	)

	log.Println("new cron:", schedule)
	_, err := c.AddFunc(schedule, func() {
		wg.Add(1)
		Execute(command, args)
		wg.Done()
	})
	if err != nil {
		log.Printf("error adding cron function: %v", err)
	}

	return c, wg
}

func Start(c *cron.Cron) {
	c.Start()
}

func Stop(c *cron.Cron, wg *sync.WaitGroup) {
	log.Println("Stopping")
	c.Stop()
	log.Println("Waiting")
	wg.Wait()
	log.Println("Exiting")
	os.Exit(0)
}

func RunJobs(c *cron.Cron) {
	for _, e := range c.Entries() {
		e.Job.Run()
	}
}
