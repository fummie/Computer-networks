package main


import (
	"fmt"
	"golang.org/x/net/trace"
	"log"
	"net/http"
)

type Fetcher struct {
	domain string
	events trace.EventLog
}

func NewFetcher(domain string) *Fetcher {
	return &Fetcher{
		domain,
		trace.NewEventLog("mypkg.Fetcher", domain),
	}
}

func (f *Fetcher) Fetch(path string) (string, error){
	resp, err := http.Get("http://" + f.domain + "/" + path)
	if err != nil{
		f.events.Errorf("Get(%q) = %v", path, err)
	}
	f.events.Printf("Get(%q) = %s", path, resp.Status)
	return fmt.Sprintf("Get(%q) = %s", path, resp.Status), nil
}

func (f *Fetcher) Close() error {
	f.events.Finish()
	return nil
}

var fetch *Fetcher

func fooHandler(w http.ResponseWriter, req *http.Request) {
	trace.Traces(w, req)
	tr := trace.New("Trace", "vk.com")
	defer tr.Finish()

	fetch = NewFetcher("vk.com")
	str, err := fetch.Fetch("")
	if err != nil {
		panic(err)
	}
	//fmt.Println(str)
	tr.LazyPrintf("some event %q happened", str)
}

func main()  {
	defer fetch.Close()
	http.HandleFunc("/", fooHandler)         // установим роутер
	err := http.ListenAndServe("localhost:3000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}