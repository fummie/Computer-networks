package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"os"
	"strings"
)
func main() {
	ssh.Handle(func(session ssh.Session) {
		fmt.Println(session.User())
		//dir, _ := os.Getwd()
		//term := terminal.NewTerminal(session, session.User() + "@" + "localhost" + ":" + dir + "$ ")
			dir, _ := os.Getwd()
			term := terminal.NewTerminal(session, session.User() + "@" + "localhost" + ":" + dir + "$ ")
			line := session.RawCommand()
			line = strings.TrimSuffix(line, "\n")
			temp := strings.Split(line, " ")
			command := temp[0]
			fmt.Println(line)
			args := temp[1:]
			if command == "ls"{
				if len(args) == 1 {
					dir += string(os.PathSeparator) + args[1]
				}
				files, err := ioutil.ReadDir(dir)
				if err != nil{
					panic(err)
				}
				for _, file := range files{
					term.Write(append([]byte(file.Name()), '\n'))
				}
			}
			if command == "mkdir"{
				if len(args) != 1 {
					term.Write([]byte("Invalid usage of mkdir"))
				} else {
					os.Mkdir(args[0], 0777)
				}
			}
			if command == "rmdir"{
				if len(args) != 1 {
					term.Write([]byte("Invalid usage of rmdir"))
				} else {
					os.Remove(args[0])
				}
			}
		log.Println("terminal closed")
	})
	log.Fatal(ssh.ListenAndServe("localhost:3000", nil, ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool{
		return pass == "Fm4Irnzv"
		}),
	))
}
