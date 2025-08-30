package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func TryConnection(host string, port uint) error {
	var (
		addr    string
		conn    net.Conn
		timeOut time.Duration
		err     error
	)

	timeOut = time.Millisecond * 500
	addr = fmt.Sprintf("%s:%d", host, port)

	conn, err = net.DialTimeout("tcp", addr, timeOut)

	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

func worker(host string, ports chan uint,wg *sync.WaitGroup) {
	defer wg.Done()
	for p := range ports {
		if TryConnection(host, p) == nil {
			fmt.Printf("port is up %d\n", p)
		}
	}
}

func Master(ports chan<- uint, s, e, reqps uint) {
	ticker := time.NewTicker(time.Second / time.Duration(reqps))
	defer ticker.Stop()

	for port := s; port <= e; port++ {
		<-ticker.C // wait for tick (controls rate)
		ports <- port
	}
	close(ports) // close after sending all ports
}

func main() {

	var (
		host             string
		startingPort     uint
		lastPort         uint
		requestPerSecond uint
	)

	flag.StringVar(&host, "host", "locolhost", "host to scan")
	flag.UintVar(&startingPort, "start_port", 1, "starting port Number")
	flag.UintVar(&lastPort, "last_port", 3000, "starting port Number")
	flag.UintVar(&requestPerSecond, "rps", 100, "request per second")
	flag.Parse()
	// args host start_port end_port requestPerSecond


	fmt.Printf("start scanning : %s\n", host)
	var (
		ports       chan uint
		nWorker     int
		wg          sync.WaitGroup
	)

	nWorker = 10
	ports = make(chan uint, 500)


	for i := 0; i < nWorker; i+= 1 {
		wg.Add(1)
		go worker(host, ports, &wg)
	}

	Master(ports, startingPort, lastPort, requestPerSecond)
	wg.Wait()

}
