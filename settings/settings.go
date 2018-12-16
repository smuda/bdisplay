package settings

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"path"
	"strconv"
)

// Settings holds all the settings for the application.
// Use the method GetSettings to get a new instance, or use
// ReadFromFile after having set the file name.
type Settings struct {
	ApplicationDir  string
	DoRedirect      bool     `json:"doRedirect"`
	Modules         []string `json:"modules"`
	RedirectTimeout []int    `json:"redirectTimeout"`
	ServerPort      string   `json:"serverPort"`
	ServerHost      string   `json:"serverHost"`
	ResourceDir     string   `json:"resourceDir"`
	TemplatesDir    string   `json:"templatesDir"`
	ContentDir      string   `json:"contentDir"`
	TemplatesExt    string   `json:"templatesExt"`
	CSSFile         string   `json:"cssFile"`
	BDUpdate        struct {
		LogFile string `json:"logFile"`
	} `json:"bdupdate"`
	Projects struct {
		FileName string `json:"fileName"`
	} `json:"projects"`
	Tweets struct {
		FileName       string `json:"fileName"`
		ConsumerKey    string `json:"consumerKey"`
		ConsumerSecret string `json:"consumerSecret"`
		AccessToken    string `json:"accessToken"`
		AccessSecret   string `json:"accessSecret"`
		QueryString    string `json:"queryString"`
		QueryLang      string `json:"queryLang"`
	} `json:"tweets"`
	Dagens struct {
		FlagImg     string `json:"flagImg"`
		DilbertImg  string `json:"dilbertImg"`
		XKCDImg     string `json:"xkcdImg"`
		FileName    string `json:"fileName"`
		SLKey       string `json:"slKey"`
		DilbertURL  string `json:"dilbertUrl"`
		NamesdayURL string `json:"namesdayUrl"`
		XKCDURL     string `json:"xkcdUrl"`
		SLURL       string `json:"slUrl"`
	} `json:"dagens"`
	Ove struct {
		MonitisURL string `json:"monitisUrl"`
		ModuleID   string `json:"moduleId"`
	} `json:"ove"`
	Monitis struct {
		MonitisURL string   `json:"monitisUrl"`
		ModuleID   []string `json:"moduleId"`
	} `json:"monitis"`
}

// New returns a new instance of a Settings struct.
// It takes a file name to read the JSON settings file from.
func New(fn, appDir string) (*Settings, error) {
	if !path.IsAbs(fn) {
		fn = path.Join(appDir, fn)
	}
	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", fn, err)
	}
	news := new(Settings)
	err = json.Unmarshal(raw, news)
	if err != nil {
		return nil, fmt.Errorf("could not parse settings file: %v", err)
	}
	news.ApplicationDir = appDir
	return news, nil
}

// GetRedirectString returns the HTML string, for the specified module, to place in the header
// of the HTML page to automatically switch to the next module page.
func (s *Settings) GetRedirectString(modName string) template.HTML {
	if s.DoRedirect {
		i, url := s.getRedirectTime(modName)
		if i >= 0 {
			return template.HTML("<meta http-equiv=\"refresh\" content=\"" + strconv.Itoa(i) + ";URL='http://" + s.ServerHost + ":" + s.ServerPort + "/" + url + "'\">")
		}
	}
	return ""
}

// Templates returns a path with files to parse as templates
func (s *Settings) Templates() string {
	if path.IsAbs(s.TemplatesDir) {
		return path.Join(s.TemplatesDir, "*"+s.TemplatesExt)
	}
	return path.Join(s.ApplicationDir, s.TemplatesDir, "*"+s.TemplatesExt)
}

func (s *Settings) getRedirectTime(val string) (int, string) {
	ro := len(s.Modules)
	if (ro != len(s.RedirectTimeout)) || ro <= 0 {
		return -1, ""
	}
	for i, v := range s.Modules {
		if val == v {
			var x int
			if i == ro-1 {
				x = 0
			} else {
				x = i + 1
			}
			return s.RedirectTimeout[i], s.Modules[x]
		}
	}
	return -1, ""
}
