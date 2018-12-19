package projects

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/hossner/bdisplay/internal/basemodule"
	"github.com/hossner/bdisplay/settings"
)

const (
	moduleName = "projs"
)

type project struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Budget string `json:"budget"`
	Cost   string `json:"cost"`
	Lead   string `json:"lead"`
}

// Projects is the base struct for the module
type Projects struct {
	basemodule.Module
	PageUpdated string
	NrOfProjs   string
	ProjectList []project
}

// New returns a new instance of the Projects struct
func New(cfg *settings.Settings, tpls *template.Template) *Projects {
	return &Projects{basemodule.Module{Settings: cfg, Name: moduleName, Template: tpls}, "", "", nil}
}

// HTTPHandler is the handler for the web page for the module
func (p *Projects) HTTPHandler(res http.ResponseWriter, req *http.Request) {
	if p.Template == nil {
		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")

	myFile := path.Join(p.Settings.ApplicationDir, p.Settings.ContentDir, p.Settings.Projects.FileName)
	fi, err := os.Stat(myFile)
	if err != nil {
		p.Show500(res, err, "Error getting file info: "+myFile)
		return
	}

	p.PageUpdated = fmt.Sprintf("%s", fi.ModTime().Format("2006-01-02 15:04"))
	p.ProjectList, err = getProjects(myFile)
	if err != nil {
		p.Show500(res, err, "Error reading project file: "+myFile)
		return
	}
	p.NrOfProjs = fmt.Sprintf("%d st", len(p.ProjectList))
	if err = p.Template.ExecuteTemplate(res, p.GetTemplateName(), p); err != nil {
		p.Show500(res, err, "Error executing template, using: "+p.GetTemplateName())
	}
}

func getProjects(theFile string) ([]project, error) {
	var c []project
	raw, err := ioutil.ReadFile(theFile)
	if err != nil {
		return nil, fmt.Errorf("could not read projects from file: %v", err)
	}
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("could not parse projects from file: %v", err)
	}
	return c, nil
}
