package ove

import (
	"html/template"
	"net/http"

	"github.com/hossner/bdisplay/internal/basemodule"
	"github.com/hossner/bdisplay/settings"
)

const (
	moduleName = "ove"
)

// Ove is the base struct for the module
type Ove struct {
	basemodule.Module
}

// New returns a new instance of the module
func New(cfg *settings.Settings, tpls *template.Template) *Ove {
	return &Ove{basemodule.Module{Settings: cfg, Name: moduleName, Template: tpls}}
}

// ModuleID is use when parsing the template
func (o *Ove) ModuleID() string {
	return o.Settings.Ove.ModuleID
}

// MonitisURL is used then parsing the template
func (o *Ove) MonitisURL() string {
	return o.Settings.Ove.MonitisURL
}

// HTTPHandler is the handler for the web page
func (o *Ove) HTTPHandler(res http.ResponseWriter, req *http.Request) {
	if o.Template == nil {
		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := o.Template.ExecuteTemplate(res, o.GetTemplateName(), o); err != nil {
		o.Show500(res, err, "Error executing template, using: "+o.GetTemplateName())
	}
}
