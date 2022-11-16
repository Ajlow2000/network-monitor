package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

var (
	Info_Level  *log.Logger
	Event_Level *log.Logger
)

func init() {
	file, err := os.OpenFile("network-monitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Info_Level = log.New(file, "INFO : ", log.Ldate|log.Ltime)
	Event_Level = log.New(file, "EVENT: ", log.Ldate|log.Ltime)
}

const wanTarget = "8.8.8.8"
const INTERVAL = 5

type Event struct {
	Desc  string
	Start time.Time
	End   time.Time
}

func ping(target string) bool {
	result := false

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		// fmt.Printf("IP Addr: %s received, RTT: %v\n", addr.String(), rtt)
		result = true
	}
	p.OnIdle = func() {
		// fmt.Println("finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func main() {
	var downtime_start time.Time
	var downtime_end time.Time

	log.Printf("Beginning monitor")
	for {
		result := ping(wanTarget)
		if !result {
			Info_Level.Printf("%v Unreachable", wanTarget)
			if downtime_start.IsZero() {
				downtime_start = time.Now()
				Event_Level.Printf("Outage Detected")
			}
		} else {
			Info_Level.Printf("%v Received", wanTarget)
			if !downtime_start.IsZero() && downtime_end.IsZero() {
				downtime_end = time.Now()
				duration := downtime_end.Sub(downtime_start)

				Event_Level.Printf("Outage Resolved - Duration (seconds): %v", duration.Seconds())

				downtime_start = time.Time{}
				downtime_end = time.Time{}
			}
		}
		time.Sleep(time.Second * INTERVAL)
	}
}
