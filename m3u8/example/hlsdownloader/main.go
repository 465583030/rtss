package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/groupcache/lru"
	"github.com/osrtss/rtss/m3u8"
)

var client = &http.Client{}

func doRequest(c *http.Client, req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", opts.UserAgent)
	resp, err := c.Do(req)
	return resp, err
}

// Downloader hlsdownloader struct
type Downloader struct {
	URI           string
	totalDuration time.Duration
}

func getPlaylist(urlStr string, recTime time.Duration, useLocalTime bool, dlc chan *Downloader) {
	startTime := time.Now()
	var recDuration time.Duration = time.Duration(0)
	cache := lru.New(1024)
	playlistUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalf("M3U8 url parse failed, err: %v", err)
	}
	for {
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			log.Fatalf("Create HTTP request failed, err: %v", err)
		}
		resp, err := doRequest(client, req)
		if err != nil {
			log.Errorf("Do HTTP request failed, err: %v", err)
			time.Sleep(time.Duration(3) * time.Second)
		}
		playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
		if err != nil {
			log.Fatalf("M3U8 decode failed, err: %v", err)
		}
		resp.Body.Close()
		if listType == m3u8.MEDIA {
			mpl := playlist.(*m3u8.MediaPlaylist)
			for _, v := range mpl.Segments {
				if v != nil {
					var msURI string
					if strings.HasPrefix(v.URI, "http") {
						msURI, err = url.QueryUnescape(v.URI)
						if err != nil {
							log.Fatalf("URL query unescape failed, err: %v", err)
						}
					} else {
						msUrl, err := playlistUrl.Parse(v.URI)
						if err != nil {
							log.Errorf("URL parse failed, err: %v", err)
							continue
						}
						msURI, err = url.QueryUnescape(msUrl.String())
						if err != nil {
							log.Fatalf("URL query unescape failed, err: %v", err)
						}
					}
					_, hit := cache.Get(msURI)
					if !hit {
						cache.Add(msURI, nil)
						if useLocalTime {
							recDuration = time.Now().Sub(startTime)
						} else {
							recDuration += time.Duration(int64(v.Duration * 1000000000))
						}
						dlc <- &Downloader{msURI, recDuration}
					}
					if recTime != 0 && recDuration != 0 && recDuration >= recTime {
						close(dlc)
						return
					}
				}
			}
			if mpl.Closed {
				close(dlc)
				return
			} else {
				time.Sleep(time.Duration(int64(mpl.TargetDuration * 1000000000)))
			}
		} else {
			log.Fatalln("Not a valid media playlist")
		}
	}
}

func downloadSegment(fn string, dlc chan *Downloader, recTime time.Duration) {
	out, err := os.Create(fn)
	if err != nil {
		log.Fatalf("Output file create failed, err: %v", err)
	}
	defer out.Close()
	for v := range dlc {
		req, err := http.NewRequest("GET", v.URI, nil)
		if err != nil {
			log.Fatalf("Create HTTP request failed, err: %v", err)
		}
		resp, err := doRequest(client, req)
		if err != nil {
			log.Errorf("Do HTTP request failed, err: %v", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Errorf("HTTP request status code not 200, code: %v", resp.StatusCode)
			continue
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalf("IO copy failed, err: %v", err)
		}
		resp.Body.Close()
		log.Infof("Downloaded %v\n", v.URI)
		if recTime != 0 {
			log.Infof("Recorded %v of %v\n", v.totalDuration, recTime)
		} else {
			log.Infof("Recorded %v\n", v.totalDuration)
		}
	}
}

func main() {
	if len(opts.M3U8) == 0 || len(opts.Output) == 0 {
		log.Fatalln("MU38 or output path is empty")
	}

	if !strings.HasPrefix(opts.M3U8, "http") {
		log.Fatalln("M3U8 playlist must begin with http/https.")
	}

	duration, err := time.ParseDuration(opts.Duration)
	if err != nil {
		log.Fatalln("Wrong duration.")
	}

	msChan := make(chan *Downloader, 1024)
	go getPlaylist(opts.M3U8, duration, opts.UseLocalTime, msChan)
	downloadSegment(opts.Output, msChan, duration)
}
