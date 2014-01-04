/*
 *  irc: IRC client library in Go
 *  Copyright (C) 2013  Joshua Chase <jcjoshuachase@gmail.com>
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License along
 *  with this program; if not, write to the Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package irc

// BUG(josh) Need to replace the wonky chan map[] hack with real writelocks

import (
	"bufio"
	"fmt"
	"github.com/Pursuit92/LeveledLogger/log"
	"github.com/Pursuit92/syncmap"
	"math/rand"
	"net"
	"time"
)

var rgen *rand.Rand

func init() {
	rgen = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Conn struct {
	Host     string
	Nicks    []string
	Nick     string
	Name     string
	RealName string
	//PingInterval int
	msgOut  chan Command
	conn    net.Conn
	expects syncmap.Map
}

func (c Conn) Expects() syncmap.Map {
	return c.expects
}

func (c Conn) MsgOut() chan Command {
	return c.msgOut
}

type IRCErr string

func (i IRCErr) Error() string {
	return string(i)
}

func DialIRC(host string, nicks []string, name, realname string /*, pingint int*/) (*Conn, error) {
	ircConn := Conn{
		Host: host,
		Nicks: nicks,
		Nick: "",
		Name: name,
		RealName: realname,
		expects: syncmap.New(),
	}
	log.Out.Printf(2,"Connecting to %s...", host)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	ircConn.conn = conn
	log.Out.Printf(2,"Connected! Performing setup...")

	ircConn.msgOut = make(chan Command, 16)

	go ircConn.recvCommands()
	go handleExpects(ircConn)
	go ircConn.pongsGalore()

	return &ircConn, nil
}

func (c Conn) sendCommand(m Command) error {
	msg := m.String()
	log.Out.Printf(2,"Sending Command: %s", msg)
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
	log.Out.Printf(2,"Starting message reciever")
	buffered := bufio.NewReader(c.conn)
	for {
		msg, err := c.recvCommand(buffered)
		if err != nil {
			log.Fatal(err)
		} else {
			c.msgOut <- *msg
		}
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
	pong := Command{Prefix: "", Command: Pong}
	pings, _ := Expect(c, Command{Command: Ping})
	for {
		ping := <-pings.Chan
		pong.Params = []string{ping.Params[0]}
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
	log.Out.Printf(2,"Attempting to register nick")
	userMsg := Command{Command: User,
		Params: []string{c.Name, "0", "*", c.RealName}}
	welcomeChan, _ := Expect(c, Command{Command: RplWelcome})
	errChan, _ := Expect(c, Command{Command: ErrNicknameinuse})
	defer UnExpect(c, welcomeChan)
	defer UnExpect(c, errChan)

	c.sendCommand(userMsg)

	var err error = nil

	var success bool = false

	nickMsg := Command{Command: Nick}

	var ret Command
	for _, v := range c.Nicks {
		nickMsg.Params = []string{v}
		c.sendCommand(nickMsg)
		log.Out.Printf(2,"Waiting for response...")
		select {
		case resp := <-welcomeChan.Chan:
			log.Out.Printf(2,"Received welcome message: %s", resp.String())
			//println(resp.String())
			ret = resp
			success = true
		case errmsg := <-errChan.Chan:
			log.Out.Printf(2,"Received error message: %s", errmsg.String())
		}
		if success {
			log.Out.Printf(2,"Done registering")
			c.Nick = v
			// Don't really need pings
			// go c.pingsGalore()
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
	c.conn.Close()
}
