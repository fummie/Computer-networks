package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"os"
	"time"
)

func main() {
	var server, login, password string

	fmt.Print("Enter the host url: ")
	if _, err := fmt.Scanf("%s", &server); err != nil {
		panic(err)
	}

	client, err := ftp.Dial(server + ":21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to " + server)

	fmt.Print("Enter login: ")
	if _, err := fmt.Scanf("%s", &login); err != nil {
		panic(err)
	}

	fmt.Print("Enter password: ")
	if _, err := fmt.Scanf("%s", &password); err != nil {
		panic(err)
	}


	if err := client.Login(login, password); err != nil {
		panic(err)
	}

	fmt.Print(login + "@" + server + "$ ")

loop:
	for {
		var command string
		if _, err := fmt.Scanf("%s", &command); err != nil {
			panic(err)
		}
		switch command {
		case "exit":
			break loop
		case "stor":
			myStor(client)
		case "retr":
			myRetr(client)
		case "makedir":
			myMakeDir(client)
		case "delete":
			myDelete(client)
		case "list":
			myList(client)
		case "help":
			fmt.Println("the commands are:\n" +
				"\tstor\tupload a file\n" +
				"\tretr\tdownload a file\n" +
				"\tmakedir\tmake a directory\n" +
				"\tdelete\tdelete a file\n" +
				"\tlist\tdirectory content\n" +
				"\texit\tclose connection")
		default:
			fmt.Println(command + ": unknown command\nrun 'help' for usage")
		}
	}

	if err := client.Quit(); err != nil {
		panic(err)
	}
}

func myStor(client *ftp.ServerConn) {
	var dest, inc string
	fmt.Println("enter destination path")
	if _, err := fmt.Scanf("%s", &dest); err != nil {
		panic(err)
	}
	fmt.Println("enter file path")
	if _, err := fmt.Scanf("%s", &inc); err != nil {
		panic(err)
	}

	file, err := os.Open(inc)
	if err != nil {
		panic(err)
	}

	if err := client.Stor(dest, file); err != nil {
		panic(err)
	}

	if err = file.Close(); err != nil {
		panic(err)
	}

	fmt.Println("file stored")
}

func myRetr(client *ftp.ServerConn) {
	var inc, dest, name, full string
	fmt.Println("enter file path")
	if _, err := fmt.Scanf("%s", &inc); err != nil {
		panic(err)
	}
	fmt.Println("enter destination path")
	if _, err := fmt.Scanf("%s", &dest); err != nil {
		panic(err)
	}
	fmt.Println("enter name for the new file")
	if _, err := fmt.Scanf("%s", &name); err != nil {
		panic(err)
	}
	full = dest + name

	resp, err := client.Retr(inc)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 0)
	if _, err = resp.Read(buf); err != nil {
		panic(err)
	}

	file, err := os.Create(full)
	if err != nil {
		panic(err)
	}

	if _, err := file.Write(buf); err != nil {
		panic(err)
	}

	if err = resp.Close(); err != nil {
		panic(err)
	}

	fmt.Println("file downloaded")
}

func myMakeDir(client *ftp.ServerConn) {
	var path, name, full string
	fmt.Println("enter the path")
	if _, err := fmt.Scanf("%s", &path); err != nil {
		panic(err)
	}
	fmt.Println("enter name for the new directory")
	if _, err := fmt.Scanf("%s", &name); err != nil {
		panic(err)
	}
	full = path + name

	if err := client.MakeDir(full); err != nil {
		panic(err)
	}

	fmt.Println("directory created")
}

func myDelete(client *ftp.ServerConn) {
	var path string
	fmt.Println("enter file path")
	if _, err := fmt.Scanf("%s", &path); err != nil {
		panic(err)
	}

	if err := client.Delete(path); err != nil {
		panic(err)
	}

	fmt.Println("file deleted")
}

func myList(client *ftp.ServerConn) {
	var path string
	fmt.Println("enter the path")
	if _, err := fmt.Scanf("%s", &path); err != nil {
		panic(err)
	}

	entries, err := client.List(path)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		fmt.Println("\t" + entry.Name)
	}
}