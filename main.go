package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"text/template"
)

var wg sync.WaitGroup
var client = &http.Client{}
var noFollowClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
var data = map[string]int8{
	"One":   1,
	"Two":   2,
	"Three": 3,
}

func print(values ...interface{}) {
	fmt.Println(values...)
}

func worry(err error) {
	if err != nil {
		panic(err)
	}
}

func guid() string {
	item, _ := func() (s string, err error) {
		b := make([]byte, 8)
		_, err = rand.Read(b)
		worry(err)
		return fmt.Sprintf("%x", b), nil
	}()
	return item
}

func format(html string, context map[string]string) string {
	tmpl, err := template.New("_format").Parse(html)
	worry(err)
	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, context)
	worry(err)
	return buf.String()
}

func get(base string) string {
	req, err := http.NewRequest("GET", base, nil)
	worry(err)
	req.Header.Add("Accept", "application/json")
	req.AddCookie(&http.Cookie{Name: "csrftoken", Value: "🤷‍♂️"})
	chrome := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"
	req.Header.Add("User-agent", chrome)
	q := req.URL.Query()
	q.Add("key", "val")
	req.URL.RawQuery = q.Encode()
	res, _ := client.Do(req)
	b, err := ioutil.ReadAll(res.Body)
	worry(err)
	defer res.Body.Close()
	print(res.Header.Get("Server"))
	print(res.Header.Get("Set-Cookie"))
	print(res.Header.Get("Content-Type"))
	print(res.Header.Get("Last-Modified"))
	print(res.StatusCode)
	return string(b)
}

func main() {

	// Iterate over data structure.
	for k, v := range data {
		print("key", k)
		print("value", v)
	}

	// Template formatting.
	html := "{{.One}}, {{.Two}} and {{.Three}}"
	context := map[string]string{
		"One":   "1",
		"Two":   "2",
		"Three": "3",
	}
	print(format(html, context))

	// GUIDs are cool.
	print(guid())

	// Make a bunch of GET requests in parallel.
	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			print(i, len(get("https://danmasq.com")))
			wg.Done()
		}(i)
	}
	wg.Wait()

	print("Done.")
}
