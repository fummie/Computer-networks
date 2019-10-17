package main

import (
	"github.com/RealJK/rss-parser-go"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
type PortalUrls struct {
    Relation map[string]string
}

func(portalUrls *PortalUrls) init() {
	portalUrls.Relation = make(map[string]string)
}

func(portalUrls *PortalUrls) add(name string, url string) (ok bool) {
	if _, ok = portalUrls.Relation[name]; !ok {
		portalUrls.Relation[name] = url
	}
	return !ok
}

func(portalUrls *PortalUrls) getUrl(source string) (url string, ok bool) {
	url, ok = portalUrls.Relation[source]
	return
}
//----------------------------------------------------------------------------------------------------------------------
type NewsBlock struct {
	Title string
	Link string
	Image string
	Description string
	Guid string

	ImgAvailable bool
}

func (block *NewsBlock) parseDescription() {
	var open, close int

	open = strings.Index(block.Description, "<img")
	if open != -1 {
		close = strings.Index(block.Description[open : ], ">")
		block.Image = block.Description[open + 10 : close]
		block.Description = block.Description[ : open] + block.Description[close + 2 : ]
		block.ImgAvailable = true
	}

	open = strings.Index(block.Description, "<a href")
	if open != -1 {
		close = strings.Index(block.Description[open : ], ">")
		block.Description = block.Description[ : open] + block.Description[open + close + 1: ]

		close = strings.Index(block.Description, "</a>")
		block.Description = block.Description[ : close] + block.Description[close + 4 : ]

		close = strings.Index(block.Description, "</br>")
		block.Description = block.Description[ : close] + block.Description[close + 5 : ]
	}

}
//----------------------------------------------------------------------------------------------------------------------
type Data struct {
	Source string
	Urls   PortalUrls

	NewsBlocks []NewsBlock
}

func (data *Data) addUrl(name string, url string) (ok bool) {
	ok = data.Urls.add(name, url)
	data.NewsBlocks = make([]NewsBlock, 0)
	return
}

func (data *Data) init() {
	data.Urls.init()
	data.NewsBlocks = make([]NewsBlock, 0)
}

func (data *Data) getUrl() (url string, ok bool) {
	url, ok = data.Urls.getUrl(data.Source)
	return
}

func (data *Data) changeSource(source string) (ok bool) {
	if source != "" {
		data.Source = source
		return true
	} else {
		return false
	}

}

func (data *Data) addNewsBlock(newBlock NewsBlock) {
	data.NewsBlocks = append(data.NewsBlocks, newBlock)
}

func (data *Data) resetNewsBlocks() {
	data.NewsBlocks = make([]NewsBlock, 0)
}

var data Data
//----------------------------------------------------------------------------------------------------------------------
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	data.changeSource(r.FormValue("source"))
	data.resetNewsBlocks()

	url, ok := data.getUrl()
	if ok {
		rssObject, err := rss.ParseRSS(url)
		if err != nil {
			var temp NewsBlock
			for v := range rssObject.Channel.Items {
				item := rssObject.Channel.Items[v]
				temp.Title = item.Title
				temp.Link = item.Link
				temp.Description = item.Description
				temp.Guid = item.Guid.Value
				temp.parseDescription()
				data.addNewsBlock(temp)
			}
		}
	}

	tmpl, _ := template.ParseFiles("index.html")
	tmpl.Execute(w, data)
}

func main() {
	data.init()
	data.addUrl("blagnews", "http://blagnews.ru/rss_vk.xml")
	data.addUrl("rssboard", "http://www.rssboard.org/files/sample-rss-2.xml")
	data.addUrl("lenta", "https://lenta.ru/rss")
	data.addUrl("mail", "https://news.mail.ru/rss/90/")
	data.addUrl("technolog", "http://technolog.edu.ru/index.php?option=com_k2&view=itemlist&layout=category&task=category&id=8&lang=ru&format=feed")
	data.addUrl("vz", "https://vz.ru/rss.xml")
	data.addUrl("appa", "http://news.ap-pa.ru/rss.xml")

	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")),)
	http.Handle("/static/", staticHandler)
	err := http.ListenAndServe("localhost:9022", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("Liten and Serve: ", err)
	}
}
