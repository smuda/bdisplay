package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/hossner/bdisplay/internal/dagens"
	"github.com/hossner/bdisplay/internal/monitis"
	"github.com/hossner/bdisplay/internal/ove"
	"github.com/hossner/bdisplay/internal/projects"
	"github.com/hossner/bdisplay/internal/tweets"
	"github.com/hossner/bdisplay/settings"
)

func main() {
	cfgFileName := flag.String("cfg", "config.json", "configuration file name")
	appDir := flag.String("dir", "", "base directory for bdisplay")
	flag.Parse()

	if appDir == nil {
		baseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatalf("could not determine application path: %v", err)
		}
		appDir = &baseDir
	}

	settings, err := settings.New(*cfgFileName, *appDir)
	if err != nil {
		log.Fatalf("could not read config file %s: %v", *cfgFileName, err)
	}

	logFile, err := startLogging(settings.ApplicationDir, settings.BDUpdate.LogFile)
	if err != nil {
		log.Fatalf("could not start logging: %v", err)
	} else {
		defer logFile.Close()
	}

	ds := dagens.New(settings, nil)
	ms := monitis.New(settings, nil)
	ss := ove.New(settings, nil)
	ps := projects.New(settings, nil)
	ts := tweets.New(settings, nil)

	if err := ds.Update(); err != nil {
		log.Println(err.Error())
	}
	if err := ms.Update(); err != nil {
		log.Println(err.Error())
	}
	if err := ss.Update(); err != nil {
		log.Println(err.Error())
	}
	if err := ps.Update(); err != nil {
		log.Println(err.Error())
	}
	if err := ts.Update(); err != nil {
		log.Println(err.Error())
	}

	log.Println("Closing log...")
}

func startLogging(appDir, fn string) (*os.File, error) {
	if !path.IsAbs(fn) {
		fn = path.Join(appDir, fn)
	}
	logFile, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(logFile)
		log.Println("Starting log...")
	}
	return logFile, err
}
