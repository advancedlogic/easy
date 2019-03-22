package main

import (
	"encoding/json"
	"fmt"
	"github.com/advancedlogic/easy/commons"
	"github.com/advancedlogic/easy/interfaces"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
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
			r.Error(err.Error())
			continue
		}

		bstr, err := json.Marshal(feed)
		if err != nil {
			r.Error(err.Error())
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
				r.Error(err.Error())
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
