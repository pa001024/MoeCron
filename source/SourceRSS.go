package source

import (
	"encoding/xml"
	"github.com/pa001024/MoeCron/util"
	"net/http"
	"time"
)

// "format":"rss"
type FeedRSS struct {
	XMLName     xml.Name       `xml:"rss"`
	Id          string         `xml:"channel>id"`
	Title       string         `xml:"channel>title"`
	Updated     string         `xml:"channel>lastBuildDate"`
	Description string         `xml:"channel>description"`
	Generator   string         `xml:"channel>generator"`
	Language    string         `xml:"channel>language"`
	Item        []*FeedRSSItem `xml:"channel>item"`
}
type FeedRSSItem struct {
	Id          string `xml:"guid"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Updated     string `xml:"pubDate"`
	Description string `xml:"description"`
	Author      string `xml:"dc:creator"`
	Comments    string `xml:"comments"`
}

type SourceRSS struct { // RSS 实现接口ISource
	ISource
	Source

	FeedUrl string `json:"feed_url"` // http://www.mediawiki.org/wiki/Special:RecentChanges?feed=rss&namespace=0
}

func (this *SourceRSS) GetChan() <-chan []*FeedInfo {
	if this.C != nil {
		return this.C
	}
	chw := make(chan []*FeedInfo)
	t := time.NewTimer(time.Duration(this.Interval) * time.Second)
	go func() {
		<-t.C
		chw <- this.Get()
	}()
	return chw
}

func (this *SourceRSS) Get() (rst []*FeedInfo) {
	f := this.FetchFeed()
	if f == nil {
		return
	}
	last := this.LastUpdate
	fetched := 0
	rst = make([]*FeedInfo, 0)
	for _, v := range f.Item {
		if fetched >= this.Limit {
			break
		}
		d, err := time.Parse(time.RFC1123, v.Updated)
		if err != nil {
			util.Log("Time Parse Fail", err)
			continue
		}
		if d.Sub(last) <= 0 { // It means if feed.Updated >= this.LastUpdate
			break
		}
		if d.After(this.LastUpdate) {
			this.LastUpdate = d
		}
		fv := this.GetByFeedRSSItem(v)
		rst = append(rst, fv)
		fetched++
	}
	return
}
func (this *SourceRSS) GetByFeedRSSItem(v *FeedRSSItem) (rst *FeedInfo) {
	rst = &FeedInfo{
		Id:       v.Id,
		SourceId: this.Name,
		Title:    v.Title,
		Author:   v.Author,
		Link:     v.Link,
		Content:  v.Description,
	}
	return
}

func (this *SourceRSS) FetchFeed() (rst *FeedRSS) {
	res, err := http.Get(this.FeedUrl)
	if err != nil {
		util.Log("FetchFeed Fail")
		return
	}
	defer res.Body.Close()
	rst = &FeedRSS{}
	xml.NewDecoder(res.Body).Decode(rst)
	return
}
