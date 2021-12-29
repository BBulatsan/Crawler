package main

import (
	"time"

	"Crawler/env"
	"Crawler/modules/finder"
	"Crawler/modules/rabbit"
)

func main() {
	massagesChan := make(chan rabbit.Message)
	RConn := rabbit.RConn{}
	RConn.InitRabbit()
	go RConn.Publisher(massagesChan)

	urls := make(map[string]int)
	urls[env.Url] = 1
	uniqUrls := make(map[string]int)
	f := finder.DataFinder{Urls: urls, UniqUrls: uniqUrls}
	for i := 0; i < env.Treads; i++ {
		go f.Finder(massagesChan)
	}

	for {
		time.Sleep(2 * env.Delay * time.Millisecond)
		if len(f.Urls) == 0 {
			break
		}
	}
}
