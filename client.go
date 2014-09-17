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
)

func quitOnError(err error) {
	if err != nil {
		log.Fatalln("INFO:", err)
	}
}

func client(host string, port string, done chan bool) {
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

func main() {
	cfg, _ := goconfig.LoadConfigFile("config.ini")
	host := cfg.MustValue("client", "host", "localhost")
	port := cfg.MustValue("client", "port", "4000")
	logs := cfg.MustBool("client", "logs", true)
	if logs {
		fw, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE, 0600)
		if err == nil {
			log.SetOutput(fw)
		} else {
			log.Println(err)
		}
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cnt := 100
	done := make(chan bool)
	for i := 0; i < cnt; i++ {
		go client(host, port, done)
	}
	for i := 0; i < cnt; i++ {
		<-done
	}
}
