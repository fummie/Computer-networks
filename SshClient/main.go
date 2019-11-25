package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

var (
	user = flag.String("u", "iu9_31_22", "User name")
	//pk   = flag.String("pk", defaultKeyPath(), "Private key file")
	password = flag.String("pass", "Fm4Irnzv", "Password")
	host = flag.String("h", "localhost", "Host")
	port = flag.Int("p", 3000, "Port")
)

func defaultKeyPath() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

func main() {
	flag.Parse()
	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error{
			return nil
		},
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		panic(err)
	}

	for{
		in := bufio.NewReader(os.Stdin)
		str, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}

		if str == "EXIT"{
			break
		} else {
			session, err := client.NewSession()
			if err != nil {
				panic(err)
			}

			stdout, err := session.StdoutPipe()
			if err != nil{
				panic(err)
			}
			go io.Copy(os.Stdout, stdout)

			stderr, err := session.StderrPipe()
			if err != nil{
				panic(err)
			}

			go io.Copy(os.Stderr, stderr)

			session.Run(str)
			session.Close()
		}
	}
}