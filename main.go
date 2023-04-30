package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct{}

// initialize the server
func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	// buf := make([]byte, 2048)
	buf := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		// n, err := conn.Read(buf)
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		// file := buf[:n]
		// fmt.Println(file)
		// panic("should panic !!") // will only work after copying is ended
		fmt.Println(buf.Bytes())
		fmt.Printf("recieved %d bytes over the network\n", n)
	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	// keep copying into buf until eof is reached
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	// n, err := conn.Write(file)
	if err != nil {
		return err
	}
	fmt.Printf("Written %d bytes over the network\n", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(2 * time.Second)
		sendFile(2000)
		fmt.Printf("this routine has ended")
	}()
	server := &FileServer{}
	server.start()
}
