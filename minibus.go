package minibus

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	Default string = ""
)

type Client struct {
	workDir string
	pid     []byte
	out     net.Conn
	cxd     bool
}

func New(workDir string) *Client {
	var err error
	if workDir == Default {
		workDir, err = os.UserCacheDir()
		if err != nil {
			panic(err)
		}
	}
	pid := []byte(strconv.Itoa(os.Getpid()))
	c := Client{workDir, pid, nil, false}
	return &c
}

func (t *Client) Send(channel string, msg []byte) error {
	if !t.cxd {
		// conenct
		fmt.Println(filepath.Join(t.workDir, "minibus", "minibus"))
		c, err := net.Dial("unixgram", filepath.Join(t.workDir, "minibus", "minibus"))
		if err != nil {
			return err
		}
		c.Write([]byte("OPEN"))
		t.out = c
		t.cxd = true
	}
	fmt.Println("sending:", fmt.Sprintf("%s:%s", channel, msg))
	_, err := t.out.Write([]byte(fmt.Sprintf("%s:%s", channel, msg)))
	if err != nil {
		return err
	}
	return nil
}

func (t Client) OpenChannel(channel string) (chan []byte, chan bool, error) {
	ch := make(chan []byte)
	closer := make(chan bool)

	go func() {
		sockPath := filepath.Join(t.workDir, "minibus", fmt.Sprintf("%s-%s", t.pid, channel))
		fmt.Println("open channel: ", sockPath)
		addr, err := net.ResolveUnixAddr("unixgram", sockPath)
		if err != nil {
			return
		}
		// listen on the socket
		conn, err := net.ListenUnixgram("unixgram", addr)
		if err != nil {
			return
		}
		// scan the connection and push to the channel
		scanner := bufio.NewScanner(conn)
		fmt.Println("Begin scanning..")
		for scanner.Scan() {
			fmt.Println("scanner has data")
			msg := scanner.Text()
			fmt.Println("gotfrom  bus:", msg)
			ch <- []byte(msg)
		}
		fmt.Println("stopped scanning")
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			close(ch)
			return
		}
		// close the connection if the closer channel is closed
		// this will cause the scanner to end and exit
		go func() {
			<-closer
			conn.Close()
		}()

		defer func() {

			fmt.Println("Closing connection..")
			// close socket when we finish
			conn.Close()
			// remove the socket file like a good citizen
			err := syscall.Unlink(sockPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Unlink()", err)
			}
		}()
	}()

	return ch, closer, nil
}
