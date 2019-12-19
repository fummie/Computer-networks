package main

import (
	"fmt"
	"github.com/sparrc/go-ping"
	"log"
	"time"
)

func pin(addr string, count int, iter int) {
	stop := time.Duration(iter * 5000)
	time.Sleep(time.Microsecond * stop)
	pinger, err := ping.NewPinger(addr)
	pinger.SetPrivileged(true)
	if err != nil {
		log.Fatal(err)
	}

	pinger.Count = count

	pinger.Run()
	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
}

func main() {
	addr := /*os.Args[1]*/ "www.google.com"
	count := /*strconv.Atoi(os.Args[2])*/ 10
	size := /*strconv.Atoi(os.Args[3])*/ 3
	for i := 0; i < size; i++ {
		go pin(addr, count, i)
	}
	pin(addr, count, size)
}
