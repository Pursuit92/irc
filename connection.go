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

import (
	"bufio"
	"fmt"
	"github.com/Pursuit92/LeveledLogger/log"
	"math/rand"
	"net"
	"time"
	"github.com/Pursuit92/pubsub"
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
	conn    net.Conn
	pubsub.Publisher
}


func DialIRC(host string, nicks []string, name, realname string) (*Conn, error) {
	ircConn := Conn{
		Host: host,
		Nicks: nicks,
		Nick: "",
		Name: name,
		RealName: realname,
	}
	log.Out.Lprintf(2,"Connecting to %s...", host)
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	ircConn.conn = conn
	log.Out.Lprintf(2,"Connected! Performing setup...")

	msgOut := make(chan pubsub.Matchable, 16)

	go recvCommands(ircConn.conn,msgOut)

	ircConn.Publisher = pubsub.MakePublisher(msgOut)

	go ircConn.pongsGalore()

	return &ircConn, nil
}

func (c Conn) Send(m Command) error {
	msg := m.String()
	log.Out.Lprintf(2,"Sending Command: %s", msg)
	_, err := fmt.Fprintf(c.conn, "%s\r\n", msg)

	if err == nil {
		return err
	} else {
		return nil
	}
}

// receive a single command
func recvCommand(buffered *bufio.Reader) (cmd *Command, err error) {
	message, err := buffered.ReadString(0x0a)
	if err == nil {
		cmd, err = parseCommand(message)
	}
	return cmd, err
}

// Sends commands to the msgOut channel till an error is encountered
func recvCommands(read net.Conn, msgOut chan pubsub.Matchable) (err error) {
	// This is the only thing that should be writing to the channel.
	// Close it when we're done.
	defer close(msgOut)

	var cmd *Command
	log.Out.Lprintf(2,"Starting message reciever")
	buffered := bufio.NewReader(read)
	for err == nil {
		cmd, err = recvCommand(buffered)
		msgOut <- CmdErr{cmd,err}
	}
	return err
}

// Watches for Pings and responds with Pongs. Pretty simple.
func (c Conn) pongsGalore() {
	pong := Command{Prefix: "", Command: Pong}
	pings, _ := c.Subscribe(Command{Command: Ping})
	defer c.UnSubscribe(pings)
	for match := range pings.Chan {
		ping := match.(CmdErr)
		if ping.Err == nil {
			pong.Params = []string{ping.Cmd.Params[0]}
			c.Send(pong)
		}
	}
}

// Registers your Nick on the IRC server. Iterates through the slice of nicks in
// (c *Conn) and returns an error if all result in an error
func (c *Conn) Register() (cmd Command, err error) {
	// Default error
	err = UserTaken
	log.Out.Lprintf(2,"Attempting to register nick")
	userMsg := Command{Command: User,
		Params: []string{c.Name, "0", "*", c.RealName}}

	// Set up all of the expectations for the messages in the exchange
	welcomeChan, _ := c.Subscribe(Command{Command: RplWelcome})
	errChan, _ := c.Subscribe(Command{Command: ErrNicknameinuse})
	defer c.UnSubscribe(welcomeChan)
	defer c.UnSubscribe(errChan)

	sendErr := c.Send(userMsg)
	if sendErr != nil {
		return cmd, sendErr
	}

	nickMsg := Command{Command: Nick}

	for _, v := range c.Nicks {
		nickMsg.Params = []string{v}
		c.Send(nickMsg)
		log.Out.Lprintf(2,"Waiting for response...")
		select {
		case match := <-welcomeChan.Chan:
			cmd := match.(CmdErr)
			if cmd.Err != nil {
				err = cmd.Err
				break
			}
			log.Out.Lprintf(2,"Received welcome message: %s", cmd.Cmd.String())
			log.Out.Lprintf(2,"Done registering")
			c.Nick = v
			err = nil
			break
		case match := <-errChan.Chan:
			errmsg := match.(CmdErr)
			if errmsg.Err != nil {
				err = errmsg.Err
				break
			}
			log.Out.Lprintf(2,"Received error message: %s", errmsg.Cmd.String())
		}
	}

	return cmd, err
}

// Sends the Quit message to the IRC server and closes the connection
func (c Conn) Quit() {
	msg := Command{Command: Quit, Params: []string{"Leaving"}}
	c.Send(msg)
	c.conn.Close()
}
