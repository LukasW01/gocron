package main

import (
	"flag"
	"fmt"
	"gocron/src"
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
	flagArgs, execArgs := SplitArgs()
	os.Args = flagArgs

	flag.Parse()
	if *help {
		fmt.Println("Usage of", os.Args[0], "[ OPTIONS ] -- [ COMMAND ]", "(build", build, ")")
		flag.PrintDefaults()
		os.Exit(0)
	}
	log.Println("Running version:", version)

	c, wg := cron.Create(*schedule, execArgs[0], execArgs[1:])
	go cron.Start(c)
	if *port != "0" {
		go cron.HttpServer(*port)
	}
	if *initrun {
		go cron.RunJobs(c)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	cron.Stop(c, wg)
}

func SplitArgs() (flagArgs []string, execArgs []string) {
	for idx, arg := range os.Args {
		if arg == "--" {
			return os.Args[:idx], os.Args[idx+1:]
		}
	}
	return os.Args, nil
}
