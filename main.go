package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
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
	cfgFileName := flag.String("cfg", "config.json", "name of config file")
	flag.Parse()

	appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("could not determine application path: %v", err)
	}

	cfg, err := settings.New(*cfgFileName, appDir)
	if err != nil {
		log.Fatalf("could not create config: %v", err)
	}

	// Parse templates
	tpls := template.Must(template.ParseGlob(cfg.Templates()))

	// Handlers for all different modules
	om := ove.New(cfg, tpls)
	http.HandleFunc(om.GetModulePath(), om.HTTPHandler)

	mm := monitis.New(cfg, tpls)
	http.HandleFunc(mm.GetModulePath(), mm.HTTPHandler)

	pm := projects.New(cfg, tpls)
	http.HandleFunc(pm.GetModulePath(), pm.HTTPHandler)

	tm := tweets.New(cfg, tpls)
	http.HandleFunc(tm.GetModulePath(), tm.HTTPHandler)

	dm := dagens.New(cfg, tpls)
	http.HandleFunc(dm.GetModulePath(), dm.HTTPHandler)

	// File server for the resources in the resourceDir
	http.Handle("/"+cfg.ResourceDir+"/",
		http.StripPrefix("/"+cfg.ResourceDir+"/", http.FileServer(http.Dir(path.Join(appDir, cfg.ResourceDir)))))

	log.Fatalln(http.ListenAndServe(cfg.ServerHost+":"+cfg.ServerPort, nil))

}
