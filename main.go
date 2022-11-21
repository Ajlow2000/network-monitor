package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/tatsushid/go-fastping"
	"gopkg.in/gomail.v2"
)

var (
	Info_Level  *log.Logger
	Event_Level *log.Logger
	Error_Level *log.Logger
)

func init() {
	file, err := os.OpenFile("network-monitor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Info_Level = log.New(file, "INFO : ", log.Ldate|log.Ltime)
	Event_Level = log.New(file, "EVENT: ", log.Ldate|log.Ltime)
	Error_Level = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
}

type Event struct {
	Loc   string
	Desc  string
	Start time.Time
	End   time.Time
}

func ping(target string) bool {
	result := false

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", target)
	if err != nil {
		Error_Level.Printf(err.Error())
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		result = true
	}
	err = p.Run()
	if err != nil {
		Error_Level.Printf(err.Error())
	}
	return result
}

func emailNotify(e Event) {
	type prettyEvent struct {
		Loc      string
		Desc     string
		Start    string
		End      string
		Duration string
	}

	pe := prettyEvent{
		e.Loc,
		e.Desc,
		e.Start.Format(time.RFC1123),
		e.End.Format(time.RFC1123),
		fmt.Sprintf("%f", e.End.Sub(e.Start).Minutes()),
	}

	var body bytes.Buffer
	t, err := template.ParseFiles("email.html")
	if err != nil {
		Error_Level.Printf("Error opening email.html template - " + err.Error())
		return
	}
	t.Execute(&body, pe)

	msg := gomail.NewMessage()
	msg.SetHeader("From", os.Getenv("NM_API_EMAIL"))
	msg.SetHeader("To", os.Getenv("NM_RECIPIENT_EMAIL"))
	msg.SetHeader("Subject", e.Desc+" at "+getSSID())
	msg.SetBody("text/html", body.String())
	msg.Attach("network-monitor.log")

	n := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("NM_API_EMAIL"), os.Getenv("NM_API_EMAIL_PASSWORD"))

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		Error_Level.Printf(err.Error())
	}
	log.Println("Email Sent to " + os.Getenv("NM_RECIPIENT_EMAIL"))
}

func getSSID() string {
	return os.Getenv("NM_SSID")
	// cmd := exec.Command("nmcli", "connection", "show", "--active")

	// out, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := cmd.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	// scan := bufio.NewScanner(out)
	// for scan.Scan() {
	// 	line := scan.Text()
	// 	if !strings.Contains(line, "wifi") {
	// 		continue
	// 	}

	// 	parts := strings.SplitN(line, " ", 2)
	// 	return parts[0]
	// }

	// if err := scan.Err(); err != nil {
	// 	Error_Level.Fatal(err)
	// }

	// if err := cmd.Wait(); err != nil {
	// 	Error_Level.Fatal(err)
	// }
	// return ""
}

func monitor(interval int) {
	target := "8.8.8.8"

	var downtime_start time.Time
	var downtime_end time.Time

	log.Printf("Beginning Network-Monitor")
	Event_Level.Printf("Beginning Network-Monitor")
	for {
		result := ping(target)
		if !result {
			Info_Level.Printf("%v Unreachable", target)
			log.Printf("%v Unreachable", target)
			if downtime_start.IsZero() {
				downtime_start = time.Now()
				Event_Level.Printf("Outage Detected")
				log.Printf("Outage Detected")
			}
		} else {
			Info_Level.Printf("%v Received", target)
			log.Printf("%v Received", target)
			if !downtime_start.IsZero() && downtime_end.IsZero() {
				downtime_end = time.Now()
				duration := downtime_end.Sub(downtime_start)

				e := Event{getSSID(), "Outage Resolved", downtime_start, downtime_end}
				Event_Level.Printf(e.Desc+" - Duration (seconds): %v", duration.Seconds())
				log.Printf(e.Desc+" - Duration (seconds): %v", duration.Seconds())
				emailNotify(e)

				downtime_start = time.Time{}
				downtime_end = time.Time{}
			}
		}
		time.Sleep(time.Minute * time.Duration(interval))
	}
}

func main() {
	interval, err := strconv.Atoi(os.Getenv("NM_PING_INTERVAL"))
	if err != nil {
		panic(err)
	}

	monitor(interval)
}
