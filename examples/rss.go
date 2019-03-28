package main

import (
	"encoding/json"
	"fmt"
	"github.com/advancedlogic/easy/commons"
	. "github.com/advancedlogic/easy/easy"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/nats-io/go-nats"
	"io/ioutil"
	"net/http"
	"strings"
)

type RSSSource struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Folder  string   `json:"folder, omitempty"`
	Urls    []string `json:"urls, omitempty"`
	Timeout int      `json:"timeout"`
}

type RSS struct {
	interfaces.Service
}

func NewRSS() *RSS {
	return &RSS{}
}

func (r *RSS) Init(service interfaces.Service) error {
	r.Service = service
	return nil
}

func (r *RSS) Close() error { return nil }

func (r *RSS) Process(data interface{}) (interface{}, error) {
	var source RSSSource
	switch data.(type) {
	case RSSSource:
		source = data.(RSSSource)
	case string:

	}
	err := r.reload(source)
	if err != nil {
		return nil, err
	}
	rss := r.download(source)
	return rss, nil
}

func (r *RSS) download(source RSSSource) []string {
	fp := gofeed.NewParser()
	commons.Shuffle(source.Urls)
	feeds := make([]string, 0)
	for _, url := range source.Urls {
		feed, err := fp.ParseURL(url)
		if err != nil {
			r.Error(err)
			continue
		}

		bstr, err := json.Marshal(feed)
		if err != nil {
			r.Error(err)
			continue
		}
		feeds = append(feeds, string(bstr))
	}

	return feeds
}

func (r *RSS) reload(source RSSSource) error {
	files, err := ioutil.ReadDir(source.Folder)
	if err != nil {
		return err
	}
	urls := make([]string, 0)
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".rss") {
			blines, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", source.Folder, file.Name()))
			if err != nil {
				r.Error(err)
				continue
			}
			lines := strings.Split(string(blines), "\n")
			for _, line := range lines {
				urls = append(urls, strings.TrimSuffix(line, "\r\n"))
			}
		}
	}

	source.Urls = urls
	return nil
}

func main() {
	if microservice, err := Default("rss"); err == nil {

		rss := NewRSS()

		process := func(source RSSSource) error {
			feeds, err := rss.Process(source)
			if err != nil {
				return err
			}

			for _, feed := range feeds.([]string) {
				microservice.Info(feed)
			}
			return nil
		}

		if err := rss.POST("/api/v1/rss", func(c *gin.Context) {
			var rssSource RSSSource
			err = c.BindJSON(&rssSource)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"error": err.Error(),
				})
			}
			if err := process(rssSource); err != nil {
				c.JSON(http.StatusBadGateway, gin.H{
					"error": err.Error(),
				})
			}
		}); err != nil {
			microservice.Error(err)
		}

		if err := rss.Subscribe("rss", func(msg *nats.Msg) {
			var rssSource RSSSource
			err = json.Unmarshal(msg.Data, &rssSource)
			if err != nil {
				microservice.Error(err)
				return
			}
			if err := process(rssSource); err != nil {
				microservice.Error(err)
			}

		}); err != nil {
			microservice.Error(err)
		}

		microservice.Run()
	}
}
