package main

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
type Data struct {
	Client *ftp.ServerConn
	Authorised bool
	Command    bool
	URL        string
	Port       string
	Login      string
	Password   string
	List []string
	Dest string
	Source string
}

var data Data
//----------------------------------------------------------------------------------------------------------------------
func StartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("StartHandler")
	http.ServeFile(w, r, "start.html")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	//if err := r.ParseForm(); err != nil {
	//	log.Fatal(err)
	//}

	if !data.Authorised {
		url := r.FormValue("url")
		port := r.FormValue("port")
		login := r.FormValue("login")
		passw := r.FormValue("password")

		if url != "" && port != "" && login != "" && passw != "" {
			data.Client, _ = ftp.Dial(url+":"+port, ftp.DialWithTimeout(5*time.Second))
			//if err != nil {
			//	log.Println(err)
			//}

			if err := data.Client.Login(login, passw); err != nil {
				log.Println(err)
			}

			data.Authorised = true
		}
	}

	http.ServeFile(w, r, "home.html")
}

func StorHandler(w http.ResponseWriter, r *http.Request) {
	data.Dest = r.FormValue("dest")
	data.Source = r.FormValue("source")
	name := data.Source[strings.LastIndex(data.Source, "/")+1:]
	data.Dest = strings.TrimSuffix(data.Dest, "/")


	content, err := ioutil.ReadFile(data.Source)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(content)

	if err := data.Client.Stor(data.Dest+"/"+name, reader); err != nil {
		log.Fatal(err)
	}

	tmpl, _ := template.ParseFiles("stor.html")
	tmpl.Execute(w, data)
}

func RetrHandler(w http.ResponseWriter, r *http.Request) {
	data.Source = r.FormValue("dest")
	name := data.Source[strings.LastIndex(data.Source, "/"):]

	resp, err := data.Client.Retr(data.Source)
	if err != nil {
		panic(err)
	}

	buf, err := ioutil.ReadAll(resp)
	if err != nil {
		panic(err)
	}

	downloads, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	downloads += "/Downloads"

	file, err := os.Create(downloads + "/" + name)
	if err != nil {
		panic(err)
	}

	if _, err := file.Write(buf); err != nil {
		panic(err)
	}

	if err = resp.Close(); err != nil {
		panic(err)
	}

	tmpl, _ := template.ParseFiles("retr.html")
	tmpl.Execute(w, data)
}

func MakedirHandler(w http.ResponseWriter, r *http.Request) {
	data.Dest = r.FormValue("dest")
	data.Dest = strings.TrimSuffix(data.Dest, "/")

	if err := data.Client.MakeDir(data.Dest); err != nil {
		panic(err)
	}

	tmpl, _ := template.ParseFiles("makedir.html")
	tmpl.Execute(w, data)
}

func RemovedirHandler(w http.ResponseWriter, r *http.Request) {
	data.Dest = r.FormValue("dest")
	data.Dest = strings.TrimSuffix(data.Dest, "/")

	if err := data.Client.RemoveDir(data.Dest); err != nil {
		panic(err)
	}

	tmpl, _ := template.ParseFiles("removedir.html")
	tmpl.Execute(w, data)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	data.Dest = r.FormValue("dest")

	if err := data.Client.Delete(data.Dest); err != nil {
		panic(err)
	}

	tmpl, _ := template.ParseFiles("delete.html")
	tmpl.Execute(w, data)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	data.Dest = r.FormValue("dest")
	data.Dest = strings.TrimSuffix(data.Dest, "/")

	entries, err := data.Client.List(data.Dest)
	if err != nil {
		panic(err)
	}

	data.List = make([]string, 0)

	for _, entry := range entries {
		data.List = append(data.List, "\t" + entry.Name)
	}

	tmpl, _ := template.ParseFiles("list.html")
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", StartHandler)
	http.HandleFunc("/home", HomeHandler)
	http.HandleFunc("/stor", StorHandler)
	http.HandleFunc("/retr", RetrHandler)
	http.HandleFunc("/makedir", MakedirHandler)
	http.HandleFunc("/removedir", RemovedirHandler)
	http.HandleFunc("/delete", DeleteHandler)
	http.HandleFunc("/list", ListHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("Liten and Serve: ", err)
	}
}
