package main

import (
	"flag"
	"gocron/src"
	"gocron/src/http"
	"gocron/src/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	build    string
	version  string
	help     = flag.Bool("h", false, "display usage")
	port     = flag.String("p", "18080", "bind healthcheck to a specific port, set to 0 to not open HTTP port at all")
	schedule = flag.String("s", "* * * * *", "schedule the task the cron style")
	initrun  = flag.Bool("i", false, "run one first time on start after initialization")
)

func main() {
	flagArgs, execArgs := util.SplitArgs()
	os.Args = flagArgs

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	log.Println("Running version:", version, "|", "Build:", build)

	c, wg := cron.Create(*schedule, execArgs[0], execArgs[1:])
	go cron.Start(c)
	if *port != "0" {
		go http.Server(*port)
	}
	if *initrun {
		go cron.RunJobs(c)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	cron.Stop(c, wg)
}
