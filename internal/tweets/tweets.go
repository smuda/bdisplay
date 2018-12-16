package tweets

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

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	moduleName = "tweets"
)

type tweet struct {
	ID   string
	Date string
	Text string
}

// Tweets is the base struct for the module
type Tweets struct {
	basemodule.Module
	PageUpdated string
	TweetList   []tweet
}

// New returns a new instance of the Tweets struct
func New(cfg *settings.Settings, tpls *template.Template) *Tweets {
	return &Tweets{basemodule.Module{Settings: cfg, Name: moduleName, Template: tpls}, "", nil}
}

// HTTPHandler is the handler for the web page
func (t *Tweets) HTTPHandler(res http.ResponseWriter, req *http.Request) {
	if t.Template == nil {
		return
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	myFile := path.Join(t.Settings.ApplicationDir, t.Settings.ContentDir, t.Settings.Tweets.FileName)
	fi, err := os.Stat(myFile)
	if err != nil {
		t.Show500(res, err, "Error getting file info: "+myFile)
		return
	}
	t.PageUpdated = fmt.Sprintf("%s", fi.ModTime().Format("2006-01-02 15:04"))
	t.TweetList, err = getTweets(myFile)
	if err != nil {
		t.Show500(res, err, "error getting tweets")
	}
	if err = t.Template.ExecuteTemplate(res, t.GetTemplateName(), t); err != nil {
		t.Show500(res, err, "Error executing template, using: "+t.GetTemplateName())
	}
}

// Update updates the modules JSON file
func (t *Tweets) Update() error {
	var myTweets []tweet
	config := oauth1.NewConfig(t.Settings.Tweets.ConsumerKey, t.Settings.Tweets.ConsumerSecret)
	token := oauth1.NewToken(t.Settings.Tweets.AccessToken, t.Settings.Tweets.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
		Query: t.Settings.Tweets.QueryString,
		Lang:  t.Settings.Tweets.QueryLang,
	})
	if err != nil {
		return fmt.Errorf("failed searching Twitter: %v", err)
	}
	for _, twt := range search.Statuses {
		tm, _ := twt.CreatedAtTime()
		myTweets = append(myTweets, tweet{
			ID:   twt.User.Name,
			Date: fmt.Sprintf("%s", tm.Format("2006-01-02 15:04")),
			Text: twt.Text,
		})
	}

	jsonTweets, err := json.Marshal(myTweets)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(t.Settings.ApplicationDir, t.Settings.ContentDir, t.Settings.Tweets.FileName), jsonTweets, 0644)
}

func getTweets(theFile string) ([]tweet, error) {
	var c []tweet
	raw, err := ioutil.ReadFile(theFile)
	if err != nil {
		return c, fmt.Errorf("could not read tweet file: %v", err)
	}
	if err = json.Unmarshal(raw, &c); err != nil {
		return c, fmt.Errorf("could not parse tweet file: %v", err)
	}
	return c, nil
}
