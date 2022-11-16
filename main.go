package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
	"gopkg.in/gomail.v2"
	// requires iwgetid util
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

<<<<<<< HEAD
func notify(e Event) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "ajlow2000.api@gmail.com")
	msg.SetHeader("To", "junkmail00310@gmail.com")
	msg.SetHeader("Subject", e.Desc)
	msg.SetBody("text/html", "<b>This is the body of the mail</b>")
	// msg.Attach("/home/User/cat.jpg")

	n := gomail.NewDialer("smtp.gmail.com", 587, "ajlow2000.api@gmail.com", "urtpnabjocusjdwe")

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}

}

func main() {
	// fmt.Println("wifi name ", wifiname.WifiName())
	var downtime_start time.Time
	var downtime_end time.Time

=======
func main() {
	var downtime_start time.Time
	var downtime_end time.Time

>>>>>>> main
	log.Printf("Beginning monitor")
	for {
		result := ping(wanTarget)
		if !result {
			Info_Level.Printf("%v Unreachable", wanTarget)
			if downtime_start.IsZero() {
				downtime_start = time.Now()
<<<<<<< HEAD
				e := Event{"Outage Detected", downtime_start, downtime_end}
				Event_Level.Printf(e.Desc)
				notify(e)
=======
				Event_Level.Printf("Outage Detected")
>>>>>>> main
			}
		} else {
			Info_Level.Printf("%v Received", wanTarget)
			if !downtime_start.IsZero() && downtime_end.IsZero() {
				downtime_end = time.Now()
				duration := downtime_end.Sub(downtime_start)

<<<<<<< HEAD
				e := Event{"Outage Resolved", downtime_start, downtime_end}
				Event_Level.Printf(e.Desc+" - Duration (seconds): %v", duration.Seconds())
				notify(e)
=======
				Event_Level.Printf("Outage Resolved - Duration (seconds): %v", duration.Seconds())
>>>>>>> main

				downtime_start = time.Time{}
				downtime_end = time.Time{}
			}
		}
		time.Sleep(time.Second * INTERVAL)
	}
}
