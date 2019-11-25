package main

import (
	"fmt"
	"github.com/sparrc/go-ping"
	//"os"
	"strconv"
)

func main() {
	pinger, err := ping.NewPinger(/*os.Args[1]*/"www.ubuntu.com")
	if err != nil {
		panic(err)
	}

	pinger.Count, err = strconv.Atoi(/*os.Args[2]*/"3")
	if err != nil {
		panic(err)
	}

	stats := pinger.Statistics()

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.OnFinish(stats)

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()

}
