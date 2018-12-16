package monitis

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/hossner/bdisplay/internal/basemodule"
	"github.com/hossner/bdisplay/settings"
)

const (
	moduleName = "monitis"
)

// Monitis is the base struct for the module
type Monitis struct {
	basemodule.Module
}

// New returns a new instance of Monitis
func New(cfg *settings.Settings, tpls *template.Template) *Monitis {
	return &Monitis{basemodule.Module{Settings: cfg, Name: moduleName, Template: tpls}}
}

// ModuleID is used when parsing the template
func (m *Monitis) ModuleID(nr string) string {
	i, err := strconv.Atoi(nr)
	if err != nil {
		return ""
	}
	if i >= len(m.Settings.Monitis.ModuleID) {
		return ""
	}
	return m.Settings.Monitis.ModuleID[i]
}

// MonitisURL is used when parsing the template
func (m *Monitis) MonitisURL() string {
	return m.Settings.Ove.MonitisURL
}

// HTTPHandler is the handler for the web page
func (m *Monitis) HTTPHandler(res http.ResponseWriter, req *http.Request) {
	if m.Template == nil {
		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := m.Template.ExecuteTemplate(res, m.GetTemplateName(), m); err != nil {
		m.Show500(res, err, "Error executing template, using: "+m.GetTemplateName())
	}
}
