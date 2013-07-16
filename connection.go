package irc

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var rgen *rand.Rand

func init() {
	rgen = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Conn struct {
	Host         string
	Port         int
	Nicks        []string
	Nick         string
	Name         string
	RealName     string
	//PingInterval int
	msgOut       chan Command
	conn         net.Conn
	expects      chan map[int]Expect
}

type IRCErr string

func (i IRCErr) Error() string {
	return string(i)
}

func DialIRC(host string, port int, nicks []string, name, realname string/*, pingint int*/) (*Conn, error) {
	ircConn := Conn{host, port, nicks, "", name, realname, /*pingint,*/ nil, nil, nil}
	log.Printf("Connecting to %s:%d...", host, port)
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	ircConn.conn = conn

	// create expects map and the default expect
	expects := map[int]Expect{}

	// create mutex so goroutines play nice
	ircConn.expects = make(chan map[int]Expect, 1)
	ircConn.expects <- expects

	ircConn.msgOut = make(chan Command, 16)

	go ircConn.recvCommands()
	go ircConn.handleExpects()

	return &ircConn, nil
}

func (c Conn) sendCommand(m Command) error {
	msg := m.String()
	log.Printf("Sending Command: %s", msg)
	_, err := fmt.Fprintf(c.conn, "%s\r\n", msg)

	if err == nil {
		return err
	} else {
		return nil
	}
}

func (c Conn) recvCommand(buffered *bufio.Reader) (*Command, error) {
	var message string
	message, err := buffered.ReadString(0x0a)
	//log.Printf("Received: %s",message)
	if err != nil {
		return nil, err
	}
	msg, err := parseCommand(message)
	if err == nil {
		//log.Printf("Received Command: %s",msg)
	}
	return msg, err
}

// Needs error handling
func (c Conn) recvCommands() {
	log.Printf("Starting message reciever")
	buffered := bufio.NewReader(c.conn)
	for {
		msg, _ := c.recvCommand(buffered)
		c.msgOut <- *msg
	}
}

func (c Conn) Send(m Command) error {
	err := c.sendCommand(m)
	if err != nil {
		return err
	}
	return nil
}

func (c Conn) pongsGalore() {
	pong := Command{Prefix: "", Command: Pong, Params: []string{c.Nick, c.Host}}
	pings, _ := c.Expect(Command{Command: Ping})
	for {
		ping := <-pings.Chan
		pong.Params = []string{c.Nick, ping.Params[0]}
		c.sendCommand(pong)
	}
}
/*
func (c Conn) pingsGalore() {
	// This needs to be made fault-tolerant
	ping := Command{Prefix: "", Command: Ping, Params: []string{c.Host}}

	pongs, _ := c.Expect(Command{Command: Pong})
	for {
		c.sendCommand(ping)
		repl :=  <-pongs.Chan
		fmt.Println(repl.String())
		time.Sleep(time.Duration(c.PingInterval) * time.Second)
	}
}
*/

func (c *Conn) Register() (Command, error) {
	log.Printf("Attempting to register nick")
	userMsg := Command{Command: User,
		Params: []string{c.Name, "0", "*", c.RealName}}
	welcomeChan, _ := c.Expect(Command{Command: RplWelcome})
	errChan, _ := c.Expect(Command{Command: ErrNicknameinuse})
	defer c.UnExpect(welcomeChan)
	defer c.UnExpect(welcomeChan)

	c.sendCommand(userMsg)

	var err error = nil

	var success bool = false

	nickMsg := Command{Command: Nick}

	var ret Command
	for _, v := range c.Nicks {
		nickMsg.Params = []string{v}
		c.sendCommand(nickMsg)
		log.Printf("Waiting for response...")
		select {
		case resp := <-welcomeChan.Chan:
			log.Printf("Received welcome message: %s", resp.String())
			//println(resp.String())
			ret = resp
			success = true
		case errmsg := <-errChan.Chan:
			log.Printf("Received error message: %s", errmsg.String())
		}
		if success {
			log.Printf("Done registering")
			c.Nick = v
			// Don't really need pings
			// go c.pingsGalore()
			go c.pongsGalore()
			break
		}
	}
	if !success {
		err = IRCErr("All Nicks in use!")
	}
	return ret, err
}

func (c Conn) Quit() {
	msg := Command{Command: Quit, Params: []string{"Leaving"}}
	c.sendCommand(msg)
}
