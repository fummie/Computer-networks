package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"html/template"
	"log"
	"net"
	"net/http"
)

//----------------------------------------------------------------------------------------------------------------------
type Data struct {
	URL      string
	Port     string
	Login    string
	Password string
}

var data Data

//----------------------------------------------------------------------------------------------------------------------
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

//----------------------------------------------------------------------------------------------------------------------
var session *ssh.Session
var client *ssh.Client
var err error
var s []byte
var flag bool

//----------------------------------------------------------------------------------------------------------------------
func SshHandler(w http.ResponseWriter, r *http.Request) {
	//http.ServeFile(w, r, "ssh.html")
	file, _ := template.ParseFiles("ssh.html")
	r.ParseForm()
	url := r.FormValue("url")
	port := r.FormValue("port")
	login := r.FormValue("login")
	passw := r.FormValue("password")
	//s = make([]byte, 0)
	if url != "" && port != "" && login != "" && passw != "" {
		config := &ssh.ClientConfig{
			User: login,
			Auth: []ssh.AuthMethod{
				ssh.Password(passw),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		addr := fmt.Sprintf("%s:%s", url, port)
		client, err = ssh.Dial("tcp", addr, config)
		if err != nil {
			panic(err)
		}

		file.Execute(w, string(s))
	}

	req := r.FormValue("req")

	if r.FormValue("button") == "button" {

		session, err = client.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()

		s, err = session.Output(req)
		if err != nil {
			panic(err)
		}

		file.Execute(w, string(s))
	}
}

//----------------------------------------------------------------------------------------------------------------------
func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/ssh", SshHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("Liten and Serve: ", err)
	}
	session.Close()
}
