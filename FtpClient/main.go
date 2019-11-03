package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	var server, login, password string

	fmt.Print("Enter the host url: ")
	if _, err := fmt.Scan(&server); err != nil {
		panic(err)
	}

	client, err := ftp.Dial(server+":21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to " + server)

	fmt.Print("Enter login: ")
	if _, err := fmt.Scan(&login); err != nil {
		panic(err)
	}

	fmt.Print("Enter password: ")
	if _, err := fmt.Scan(&password); err != nil {
		panic(err)
	}

	if err := client.Login(login, password); err != nil {
		panic(err)
	}

	var command string
	args := make([]string, 0)

loop:
	for {
		fmt.Print(login + "@" + server + ":$ ")

		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}
		line = strings.TrimSuffix(line, "\n")
		temp := strings.Split(line, " ")
		command = temp[0]
		args = temp[1:]

		if command != "" {
			switch command {
			case "exit":
				break loop
			case "stor":
				myStor(client, args)
			case "retr":
				myRetr(client, args)
			case "makedir":
				myMakeDir(client, args)
			case "delete":
				myDelete(client, args)
			case "list":
				myList(client, args)
			case "help":
				fmt.Println("The commands are:\n" +
					"\tstor <destination path> <file path>\tUpload the file\n" +
					"\tretr <file path>\t\t\tDownload the file\n" +
					"\tmakedir <directory path>\t\tMake the directory\n" +
					"\tdelete <file path>\t\t\tDelete the file\n" +
					"\tlist <directory path>\t\t\tThe directory content\n" +
					"\texit\t\t\t\t\tClose the connection")
			default:
				fmt.Println(command + ": Unknown command\nRun 'help' for usage")
			}
		}
	}

	if err := client.Quit(); err != nil {
		panic(err)
	}
}

func myStor(client *ftp.ServerConn, args []string) {
	var dest, source, name string

	if len(args) != 2 {
		fmt.Println("Invalid usage of 'stor'\nRun 'help'")
	} else {
		dest = args[0]
		source = args[1]
		name = source[strings.LastIndex(source, "/")+1:]
		dest = strings.TrimSuffix(dest, "/")

		content, err := ioutil.ReadFile(source)
		if err != nil {
			panic(err)
		}

		reader := bytes.NewReader(content)

		if err := client.Stor(dest+"/"+name, reader); err != nil {
			panic(err)
		}

		fmt.Println("File stored")
	}
}

func myRetr(client *ftp.ServerConn, args []string) {
	var dest, name, downloads string

	if len(args) != 1 {
		fmt.Println("Invalid usage of 'retr'\nRun 'help'")
	} else {
		dest = args[0]
		name = dest[strings.LastIndex(dest, "/"):]

		resp, err := client.Retr(dest)
		if err != nil {
			panic(err)
		}

		buf, err := ioutil.ReadAll(resp)
		if err != nil {
			panic(err)
		}

		downloads, err = os.UserHomeDir()
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

		fmt.Println("File downloaded")
	}
}

func myMakeDir(client *ftp.ServerConn, args []string) {
	var path string

	if len(args) != 1 {
		fmt.Println("Invalid usage of 'makedir'\nRun 'help'")
	} else {
		path = args[0]
		path = strings.TrimSuffix(path, "/")

		if err := client.MakeDir(path); err != nil {
			panic(err)
		}

		fmt.Println("Directory created")
	}
}

func myDelete(client *ftp.ServerConn, args []string) {
	var path string

	if len(args) != 1 {
		fmt.Println("Invalid usage of 'delete'\nRun 'help'")
	} else {
		path = args[0]

		if err := client.Delete(path); err != nil {
			panic(err)
		}

		fmt.Println("File deleted")
	}
}

func myList(client *ftp.ServerConn, args []string) {
	var path string

	if len(args) != 1 {
		fmt.Println("Invalid usage of 'list'\nRun 'help'")
	} else {
		path = args[0]
		path = strings.TrimSuffix(path, "/")

		entries, err := client.List(path)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			fmt.Println("\t" + entry.Name)
		}
	}
}
