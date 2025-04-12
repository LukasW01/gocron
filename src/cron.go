package cron

import (
	"io"
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
	ExitStatus   int
	Stdout       string
	Stderr       string
	ExitTime     string
	Pid          int
	StartingTime string
}

type Proc struct {
	Running  map[string]*Status
	Status      *Status
	Schedule string
}

var proc Proc

func Output(out *string, src io.ReadCloser, pid int) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if n != 0 {
			s := string(buf[:n])
			*out = *out + s
			log.Printf("%d: %v", pid, s)
		}
		if err != nil {
			break
		}
	}
}

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
	run.StartingTime = time.Now().Format(time.RFC3339)
	run.Pid = cmd.Process.Pid
	proc.Running[strconv.Itoa(run.Pid)] = run

	go Output(&run.Stdout, stdout, run.Pid)
	go Output(&run.Stderr, stderr, run.Pid)

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

	run.ExitTime = time.Now().Format(time.RFC3339)

	delete(proc.Running, strconv.Itoa(run.Pid))
	//run.Pid = 0
	proc.Status = run

	log.Println(run.Pid, "cmd:", command, strings.Join(args, " "))
}

func Create(schedule string, command string, args []string) (cr *cron.Cron, wgr *sync.WaitGroup) {
	proc = Proc{map[string]*Status{}, &Status{}, schedule}

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
