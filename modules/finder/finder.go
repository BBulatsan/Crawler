package finder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"time"

	"Crawler/env"
	"Crawler/modules/errors"
	"Crawler/modules/rabbit"
)

type DataFinder struct {
	html     string
	Urls     chan string
	UniqUrls map[string]int
	mutex    sync.Mutex
}

func (df *DataFinder) download(url string) {
	res, err := http.Get(url)
	if err != nil {
		errors.HandlerError(err)
	}
	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		errors.HandlerError(err)
	}
	err = res.Body.Close()
	if err != nil {
		errors.HandlerError(err)
	}
	df.mutex.Lock()
	df.html = string(html)
	df.mutex.Unlock()
	time.Sleep(env.Delay * time.Millisecond)
}

func (df *DataFinder) finderUrls() {
	findHref := regexp.MustCompile("<a.*?href=\"(https://ff.ua/.*?)\"")
	df.mutex.Lock()
	allHref := findHref.FindAllStringSubmatch(df.html, -1)
	df.mutex.Unlock()
	for _, i := range allHref {
		n := checkNesting(i[1])
		if env.NestingEnable == true {
			_, isOK := df.UniqUrls[i[1]]
			if n <= env.Nesting && isOK == false {
				df.Urls <- i[1]
				df.mutex.Lock()
				df.UniqUrls[i[1]] = n
				fmt.Println(df.UniqUrls, len(df.UniqUrls))
				df.mutex.Unlock()
			} else {
				continue
			}
		} else {
			if _, isOK := df.UniqUrls[i[1]]; isOK == false {
				df.Urls <- i[1]
				df.UniqUrls[i[1]] = n
			}
		}
	}
}

func checkNesting(url string) int {
	find := regexp.MustCompile("/[^/]")
	return len(find.FindAllIndex([]byte(url), -1))
}

func (df *DataFinder) Finder(ch chan rabbit.Message) {
	for u := range df.Urls {
		df.download(u)
		df.finderUrls()
		ms := rabbit.Message{Rk: "*.urls.#", Mess: u}
		ch <- ms
	}
}
