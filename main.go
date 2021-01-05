package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type PostBody struct {
	Page int
}

type Jobs struct {
	Title  string
	Salary string
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	title := os.Args[1]
	if title == " " {
		log.Fatal("Please enter job title")
		os.Exit(1)
	}

	log.Info("Search started ...")
	handle(title)
	log.Info("Search finished")
}

func handle(jobTitle string) {
	var page = 1
	var jobs []Jobs

	for {
		log.Infof("Page %d was processed", page)
		time.Sleep(time.Millisecond * 10)
		client := PostBody{
			Page: page,
		}

		resp := client.makeRequest("https://www.rabota.md/ru/all")
		parsedJobs := parseJobs(resp.Body)
		if len(parsedJobs) == 0 {
			break
		}

		for _, item := range parsedJobs {
			if strings.Contains(item.Title, jobTitle) {
				jobs = append(jobs, item)
			}
		}

		page++
	}

	file, _ := json.MarshalIndent(jobs, "", " ")
	_ = ioutil.WriteFile("jobs.json", file, 0644)
}

func (postBody *PostBody) makeRequest(urlAddress string) *http.Response {
	urlValues := url.Values{
		"page": {strconv.Itoa(postBody.Page)},
	}

	resp, err := http.PostForm(urlAddress, urlValues)

	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func parseJobs(r io.ReadCloser) []Jobs {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}

	var jobs []Jobs

	doc.Find("div.preview").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find(".vacancy-title").Text()
		salary := selection.Find("span.span_salary").Text()
		jobs = append(jobs, Jobs{
			Title:  title,
			Salary: salary,
		})
	})

	return jobs
}
