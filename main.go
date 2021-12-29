package main

import (
	"fmt"

	"Crawler/env"
	"Crawler/modules/finder"
	"Crawler/modules/rabbit"
)

func main() {
	massagesChan := make(chan rabbit.Message)
	RConn := rabbit.RConn{}
	RConn.InitRabbit()
	go RConn.Publisher(massagesChan)

	urls := make(chan string, 10)
	uniqUrls := make(map[string]int)
	f := finder.DataFinder{Urls: urls, UniqUrls: uniqUrls}
	f.Urls <- env.Url
	for i := 0; i < env.Treads; i++ {
		go f.Finder(massagesChan)
	}

	var ex string
	fmt.Scan(&ex)
}
