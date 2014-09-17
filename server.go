package main

import (
	"bufio"
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

func quitOnError(err error) {
	if err != nil {
		log.Fatalln("Error:", err.Error())
	}
}

func serve(c net.Conn, counter *int32) {
	fmt.Fprintln(c, "Welcome to the counter server.")
	cnt := atomic.LoadInt32(counter)
	fmt.Fprintf(c, "%d\n", cnt)
	got, err := bufio.NewReader(c).ReadString('\n')
	quitOnError(err)
	got = strings.TrimSpace(got)
	gotInt, _ := strconv.Atoi(got)
	if gotInt == 0 || gotInt != int(cnt+1) {
		log.Printf("WARN: Sent %d, Got '%s'.", cnt, got)
		fmt.Fprintln(c, "Wrong.")
	} else if gotInt != int(cnt+1) {
		log.Printf("WARN: Sent %d, Got '%s'.", cnt, got)
		fmt.Fprintln(c, "Wrong.")
	} else {
		atomic.AddInt32(counter, 1)
		log.Printf("INFO: Sent %d, Got %d.", cnt, gotInt)
		fmt.Fprintln(c, "Correct.")
	}
	c.Close()
}

func main() {
	cfg, _ := goconfig.LoadConfigFile("config.ini")
	host := cfg.MustValue("server", "host", "localhost")
	port := cfg.MustValue("server", "port", "4000")
	logs := cfg.MustBool("server", "logs", false)
	if logs {
		fw, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE, 0600)
		if err == nil {
			log.SetOutput(fw)
		} else {
			log.Println(err)
		}
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)

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
