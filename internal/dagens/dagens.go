package dagens

import (
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	// Imported to be able to handle png file size
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/hossner/bdisplay/internal/basemodule"
	"github.com/hossner/bdisplay/settings"
)

const (
	moduleName = "dagens"
)

// SLResponseResult holds the response from SL (see SLResult)
type SLResponseResult struct {
	Created                 time.Time `json:"Created"`
	MainNews                bool      `json:"MainNews"`
	SortOrder               int       `json:"SortOrder"`
	Header                  string    `json:"Header"`
	Details                 string    `json:"Details"`
	Scope                   string    `json:"Scope"`
	DevCaseGid              int64     `json:"DevCaseGid"`
	DevMessageVersionNumber int       `json:"DevMessageVersionNumber"`
	ScopeElements           string    `json:"ScopeElements"`
	FromDateTime            string    `json:"FromDateTime"`
	UpToDateTime            string    `json:"UpToDateTime"`
	Updated                 time.Time `json:"Updated"`
}

// SLResult holds the response from SL
type SLResult struct {
	StatusCode    int                `json:"StatusCode"`
	Message       interface{}        `json:"Message"`
	ExecutionTime int                `json:"ExecutionTime"`
	Responses     []SLResponseResult `json:"ResponseData"`
}

// NamnsdagsResult holds the response from dryg.net
type NamnsdagsResult struct {
	Cachetid   string `json:"cachetid"`
	Version    string `json:"version"`
	URI        string `json:"uri"`
	Startdatum string `json:"startdatum"`
	Slutdatum  string `json:"slutdatum"`
	Dagar      []struct {
		Datum        string   `json:"datum"`
		Veckodag     string   `json:"veckodag"`
		ArbetsfriDag string   `json:"arbetsfri dag"`
		RDDag        string   `json:"röd dag"`
		Vecka        string   `json:"vecka"`
		DagIVecka    string   `json:"dag i vecka"`
		Namnsdag     []string `json:"namnsdag"`
		Flaggdag     string   `json:"flaggdag"`
	} `json:"dagar"`
}

// XKCDResult holds the result from xkcd.como
type XKCDResult struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

// Dagens holds the collected info for the module
type Dagens struct {
	basemodule.Module
	PageUpdated string
	Content     PageData
	FlagURL     template.HTML
}

// PageData is a struct as part of Dagens
type PageData struct {
	Namnsdag   string
	Flaggdag   bool
	DilbertURL template.HTML
	XKCDURL    template.HTML
	XKCDx      int
	XKCDy      int
	SLFel      []SL
}

// SL is a struct as part of PageData
type SL struct {
	Trafiklinje string
	Meddelande  []string
}

// New returns a new instance of the module
func New(cfg *settings.Settings, tpls *template.Template) *Dagens {
	return &Dagens{basemodule.Module{Settings: cfg, Name: moduleName, Template: tpls}, "", PageData{}, ""}
}

// HTTPHandler is the module's handler for the web page
func (d *Dagens) HTTPHandler(res http.ResponseWriter, req *http.Request) {
	if d.Template == nil {
		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	myFile := path.Join(d.Settings.ApplicationDir, d.Settings.ContentDir, d.Settings.Dagens.FileName)
	fi, err := os.Stat(myFile)
	if err != nil {
		d.Show500(res, err, "Error getting file info: "+myFile)
		return
	}

	d.PageUpdated = fmt.Sprintf("%s", fi.ModTime().Format("2006-01-02 15:04"))
	d.Content, err = getDagensData(myFile)
	if err != nil {
		d.Show500(res, err, "Error getting info from file: "+myFile)
		return
	}

	if d.Content.Namnsdag == "" {
		d.Content.Namnsdag = "Ingen... verkar det som..."
	}
	if d.Content.Flaggdag {
		d.FlagURL = "<img src=\"" + template.HTML(path.Join("/", d.Settings.ResourceDir, d.Settings.Dagens.FlagImg)) + "\" class=\"flag\"></img>"
	}
	if d.Content.DilbertURL != "" {
		d.Content.DilbertURL = "<img src=\"" + template.HTML(path.Join("/", d.Settings.ResourceDir, string(d.Content.DilbertURL))) + "\" class=\"dilbert\"></img>"
	}
	if d.Content.XKCDURL != "" {
		d.Content.XKCDURL = "<img src=\"" + template.HTML(path.Join("/", d.Settings.ResourceDir, string(d.Content.XKCDURL))) + "\" class=\"XKCD\"></img>"
		d.Content.XKCDx, d.Content.XKCDy = getImageDimension(d.Settings.ApplicationDir + path.Join("/", d.Settings.ResourceDir, d.Settings.Dagens.XKCDImg))
	}

	if err := d.Template.ExecuteTemplate(res, d.GetTemplateName(), d); err != nil {
		d.Show500(res, err, "Error executing template, using: "+d.GetTemplateName())
	}

}

func getDagensData(theFile string) (PageData, error) {
	var c PageData
	raw, err := ioutil.ReadFile(theFile)
	if err != nil {
		return c, fmt.Errorf("could not read the file: %s %v", theFile, err)
	}
	if err := json.Unmarshal(raw, &c); err != nil {
		return c, fmt.Errorf("could not parse the file: %s %v", theFile, err)
	}
	return c, nil
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Println(err.Error())
		return 750, 470
	}
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println(err.Error())
		return 750, 470
	}
	return img.Width + 10, img.Height + 35
}

// Update updates the modules JSON file.
// In this case, the JSON file is created by other means
func (d *Dagens) Update() error {
	var err error
	// Namesday and flag day...
	year, month, day := time.Now().Date()
	d.Content.Namnsdag, d.Content.Flaggdag, err = getNamnsdag(fmt.Sprintf(d.Settings.Dagens.NamesdayURL, year, month, day))
	if err != nil {
		log.Printf("failed to get namesday info: %v", err)
	}

	// Dilbert strip...
	if err := getCartoon(fmt.Sprintf(d.Settings.Dagens.DilbertURL, time.Now().Format("20060102")), path.Join(d.Settings.ApplicationDir, d.Settings.ResourceDir, d.Settings.Dagens.DilbertImg)); err != nil {
		log.Printf("failed getting Dilbert image: %v", err)
	}
	d.Content.DilbertURL = template.HTML(d.Settings.Dagens.DilbertImg)

	// XKCD strip...
	if err := getXKCD(d.Settings.Dagens.XKCDURL,
		path.Join(d.Settings.ApplicationDir, d.Settings.ResourceDir, d.Settings.Dagens.XKCDImg)); err != nil {
		log.Printf("failed getting XKCD image: %v", err)
	}
	d.Content.XKCDURL = template.HTML(d.Settings.Dagens.XKCDImg)

	// SL errors...
	d.Content.SLFel = getSLFel(d.Settings.Dagens.SLKey, d.Settings.Dagens.SLURL)
	jsonDagens, err := json.Marshal(d.Content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(d.Settings.ApplicationDir, d.Settings.ContentDir, d.Settings.Dagens.FileName), jsonDagens, 0644)
}

func getNamnsdag(myURL string) (string, bool, error) {
	myClient := &http.Client{Timeout: 10 * time.Second}
	var nd NamnsdagsResult
	res, err := myClient.Get(myURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to get URL %s: %v", myURL, err)
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&nd)
	if err != nil {
		return "", false, fmt.Errorf("failed to JSON decode response from %s: %v", myURL, err)
	}
	return strings.Join(nd.Dagar[0].Namnsdag, ", "), (nd.Dagar[0].Flaggdag != ""), nil
}

func getCartoon(myURL, fn string) error {
	response, err := http.Get(myURL)
	if err != nil {
		return fmt.Errorf("could not get image from URL %s: %v", myURL, err)
	}
	defer response.Body.Close()
	file, err := os.Create(fn)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("could not create file %s: %v", fn, err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("could not save downloaded image to file %s: %v", fn, err)
	}
	return nil
}

func getXKCD(myURL, fn string) error {
	myClient := &http.Client{Timeout: 10 * time.Second}
	var xr XKCDResult
	res, err := myClient.Get(myURL)
	if err != nil {
		return fmt.Errorf("could not get URL %s: %v", myURL, err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&xr); err != nil {
		return fmt.Errorf("could not JSON decode response from %s: %v", myURL, err)
	}
	if err := getCartoon(xr.Img, fn); err != nil {
		return fmt.Errorf("failed to retrieve cartoon image: %v", err)
	}
	return nil
}

func getSLFel(slKey, slURL string) []SL {
	myClient := &http.Client{Timeout: 10 * time.Second}
	var sl SLResult
	var ret []SL
	year, month, day := time.Now().Date()
	res, err := myClient.Get(fmt.Sprintf(slURL+"?key="+slKey+"&transportMode=metro,train&fromDate=%d-%02d-%02d&toDate=%d-%02d-%02d", year, month, day, year, month, day))
	if err != nil {
		log.Println(err.Error())
		return ret
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&sl)
	if err != nil {
		log.Println(err.Error())
		return ret
	}
	sort.Slice(sl.Responses, func(i int, j int) bool { return sl.Responses[i].ScopeElements > sl.Responses[j].ScopeElements })
	var str = ""
	var tl *SL
	for _, v := range sl.Responses {
		if v.ScopeElements != str {
			if str != "" {
				ret = append(ret, *tl)
			}
			tl = new(SL)
			tl.Trafiklinje = v.ScopeElements
			str = v.ScopeElements
		}
		tl.Meddelande = append(tl.Meddelande, v.Header)
	}
	return ret
}
