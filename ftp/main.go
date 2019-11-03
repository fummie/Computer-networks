package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"os"
	"time"
)

func main() {
	client, err := ftp.Dial("students.yss.su:21", ftp.DialWithTimeout(5*time.Second))

	if err != nil {
		panic(err)
	}

	if err := client.Login("ftpiu8", "3Ru7yOTA"); err != nil {
		panic(err)
	}

loop:
	for {
		var command string
		fmt.Println("enter the command")
		fmt.Scan(command)
		switch command {
		case "exit":
			break loop
		case "stor":
			myStor(client)
		}
	}

	if err := client.Quit(); err != nil {
		panic(err)
	}
}

func myStor(client *ftp.ServerConn) {
	var dest, inc string
	fmt.Println("enter destination path")
	fmt.Scan(dest)
	fmt.Println("enter file path")
	fmt.Scan(inc)

	file, err := os.Open(inc)
	if err != nil {
		panic(err)
	}

	if err := client.Stor(dest, file); err != nil {
		panic(err)
	}

	file.Close()
}
