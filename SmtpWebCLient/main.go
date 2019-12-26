package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
)
//----------------------------------------------------------------------------------------------------------------------
var addr, port, sender, recipient, subject, message, full string
var auth smtp.Auth
var tlsConfig *tls.Config
var conn *tls.Conn
var client *smtp.Client
var err error
//----------------------------------------------------------------------------------------------------------------------
type Data struct {
	addr        string
	port       string
	sender      string
	passw   string
	answ string
}

var data Data
//----------------------------------------------------------------------------------------------------------------------
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
}
//----------------------------------------------------------------------------------------------------------------------
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	addr = r.FormValue("url")
	port = r.FormValue("port")
	sender = r.FormValue("login")
	passw := r.FormValue("password")
	//fmt.Println("0000000000000000 ", addr, port, sender, passw)
	if addr != "" && port != "" && sender != "" && passw != "" {
		data.addr = addr
		data.port = port
		data.sender = sender
		data.passw = passw
		//fmt.Println("---------------" + data.sender)

		auth = smtp.PlainAuth("", sender, string(passw), addr)

		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName: addr,
		}

		conn, err = tls.Dial("tcp", addr+":"+port, tlsConfig)
		if err != nil {
			log.Panic(err)
		}

		client, err = smtp.NewClient(conn, addr)
		if err != nil {
			log.Panic(err)
		}

		if err = client.Auth(auth); err != nil {
			log.Panic(err)
		}
	}

	http.ServeFile(w, r, "home.html")
	r.ParseForm()

	recipient = r.FormValue("recipient")
	subject = r.FormValue("subject")
	message = r.FormValue("message")
	if recipient != "" && subject != "" && message != "" {
		headers := make(map[string]string)
		headers["From"] = data.sender
		headers["To"] = recipient
		headers["Subject"] = subject

		for k, v := range headers {
			full += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		full += "\r\n" + message

		//fmt.Println("*******************" + data.sender)
		if err = client.Mail(data.sender); err != nil {
			log.Panic(err)
		}

		if err = client.Rcpt(recipient); err != nil {
			log.Panic(err)
		}

		writer, err := client.Data()
		if err != nil {
			log.Panic(err)
		}

		_, err = writer.Write([]byte(full))
		if err != nil {
			log.Panic(err)
		}

		err = writer.Close()
		if err != nil {
			log.Panic(err)
		}
	}
}
//----------------------------------------------------------------------------------------------------------------------
func main() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/home", HomeHandler)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("Liten and Serve: ", err)
	}
}