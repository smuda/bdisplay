package basemodule

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/hossner/bdisplay/settings"
)

// Module is the base struct for the different "modules", or pages
type Module struct {
	Settings *settings.Settings
	Name     string
	Template *template.Template
}

// HTTPHandler must be shadowd by all implemented modules
func (m *Module) HTTPHandler(res http.ResponseWriter, req *http.Request) {
	return
}

// GetModulePath returns the end point for the module
func (m *Module) GetModulePath() string {
	return "/" + m.Name
}

// RedirectString returns the HTML tag for the page to redirect to the next module
func (m *Module) RedirectString() template.HTML {
	return m.Settings.GetRedirectString(m.Name)
}

// CSSFile returns the correct path to the CSS used in the pages
func (m *Module) CSSFile() string {
	return path.Join("/", m.Settings.ResourceDir, m.Settings.CSSFile)
}

// GetTemplateName returns a correctly formatted path to the HTML template for the module
func (m *Module) GetTemplateName() string {
	return m.Name + m.Settings.TemplatesExt
}

// Show500 shows and logs an error message
func (m *Module) Show500(res http.ResponseWriter, err error, msg string) {
	log.Println(msg)
	log.Println(err.Error())
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte("<H1>That's a 500</H1>Sorry, something bad happened!<br>More info in the server logs..."))
}

// Update updates the modules JSON file, where applied. Modules must shadow this func where applied
func (m *Module) Update() error {
	return nil
}
