package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/Unknwon/goconfig"
)

var debug string

func quitOnError(err error) {
	if err != nil {
		log.Fatalln("INFO:", err)
	}
}

func worker(host string, port string, done chan bool) {
	c, err := net.Dial("tcp", host+":"+port)
	quitOnError(err)
	br := bufio.NewReader(c)
	got, err := br.ReadString('\n')
	quitOnError(err)
	got, err = br.ReadString('\n')
	quitOnError(err)
	got = strings.TrimSpace(got)
	gotInt, _ := strconv.Atoi(got)
	fmt.Fprintf(c, "%d\n", gotInt+1)
	log.Printf("INFO: Got %s, Sent %d.", got, gotInt+1)
	// got, err = br.ReadString('\n')
	c.Close()
	done <- true
}

func client(host, port string) {
	cnt := 100
	done := make(chan bool)
	for i := 0; i < cnt; i++ {
		go worker(host, port, done)
	}
	for i := 0; i < cnt; i++ {
		<-done
	}
}

func do(c net.Conn, counter int) {
	br := bufio.NewReader(c)
	got, err := br.ReadString('\n')
	quitOnError(err)
	got = strings.TrimSpace(got)
	gotInt, _ := strconv.Atoi(got)
	fmt.Fprintf(c, "%d\n", gotInt+1)
	log.Printf("INFO: Got %s, Sent %d.", got, gotInt+1)
	got, err = br.ReadString('\n')
}

func serve(c net.Conn, counter *int32) {
	fmt.Fprintln(c, "Welcome to the counter server.")
Start:
	cnt := atomic.LoadInt32(counter)
	fmt.Fprintf(c, "What is %d + 1?\n", cnt)
	got, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Printf("INFO: Read socket error. (%s)", err)
		c.Close()
		return
	}
	got = strings.TrimSpace(got)
	if got == "" {
		log.Printf("INFO: Got nothing, bye!")
		fmt.Fprintln(c, "Bye!")
		c.Close()
		return
	}
	gotInt, _ := strconv.Atoi(got)
	if gotInt == 0 || gotInt != int(cnt+1) {
		log.Printf("WARN: Sent %d, Got %s.", cnt, got)
		fmt.Fprintln(c, "Wrong.")
	} else if gotInt != int(cnt+1) {
		log.Printf("WARN: Sent %d, Got %s.", cnt, got)
		fmt.Fprintln(c, "Wrong.")
	} else {
		atomic.AddInt32(counter, 1)
		log.Printf("INFO: Sent %d, Got %d.", cnt, gotInt)
		fmt.Fprintln(c, "Correct.")
	}
	goto Start
}

func server(host, port string) {
	var counter int32 = 0
	l, err := net.Listen("tcp", host+":"+port)
	quitOnError(err)
	log.Printf("INFO: Listening on: %s:%s\n", host, port)
	for {
		c, err := l.Accept()
		quitOnError(err)
		go serve(c, &counter)
	}
}

func main() {
	cfg, _ := goconfig.LoadConfigFile("config.ini")
	mode := cfg.MustValue("main", "mode", "client")
	host := cfg.MustValue("main", "host", "localhost")
	port := cfg.MustValue("main", "port", "4000")
	logs := cfg.MustBool("main", "logs", false)
	if logs {
		fw, err := os.OpenFile("tcpecho.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err == nil {
			log.SetOutput(fw)
			log.SetFlags(log.Lshortfile | log.LstdFlags)
			defer fw.Close()
		} else {
			log.Print(err)
		}
	}

	if debug != "" {
		mode = "server"
	}

	if mode == "server" {
		server(host, port)
	} else {
		client(host, port)
	}
}
