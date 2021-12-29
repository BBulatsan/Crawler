package finder

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"Crawler/env"
	"Crawler/modules/errors"
	"Crawler/modules/rabbit"
)

type DataFinder struct {
	html     string
	Urls     map[string]int
	UniqUrls map[string]int
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
	df.html = string(html)
	time.Sleep(env.Delay * time.Millisecond)
}

func (df *DataFinder) finderUrls() {
	findHref := regexp.MustCompile("<a.*?href=\"(https://ff.ua/.*?)\"")
	allHref := findHref.FindAllStringSubmatch(df.html, -1)
	for _, i := range allHref {
		n := checkNesting(i[1])
		if env.NestingEnable == true {
			_, isOK := df.UniqUrls[i[1]]
			if n <= env.Nesting && isOK == false {
				df.Urls[i[1]] = n
				df.UniqUrls[i[1]] = n
				fmt.Println(df.Urls)
			} else {
				continue
			}
		} else {
			if _, isOK := df.UniqUrls[i[1]]; isOK == false {
				df.Urls[i[1]] = n
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
	for {
		if len(df.Urls) != 0 {
			for u, _ := range df.Urls {
				df.download(u)
				df.finderUrls()
				delete(df.Urls, u)
				ms := rabbit.Message{Rk: "*.urls.#", Mess: u}
				ch <- ms
			}
		} else {
			fmt.Println("Messages send to rabbit: ", len(df.UniqUrls))
			break
		}
	}
}
