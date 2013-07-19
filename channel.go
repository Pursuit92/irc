package irc

import (
	"fmt"
)

type Channel struct {
	Name	string
	conn *Conn
	msgOut chan Command
	expects chan map[int] Expectation
}

func (c Channel) Expects() chan map[int] Expectation {
	return c.expects
}

func (c Channel) MsgOut() chan Command {
	return c.msgOut
}

func (c *Conn) Join(channel string) (*Channel, error) {
	joinCmd := Command{Command: Join, Prefix: c.Nick, Params: []string{channel}}
	c.sendCommand(joinCmd)
	msgs,_ := Expect(c,Command{"","PRIVMSG",[]string{channel}})
	eChan := make(chan map[int] Expectation, 1)
	expects := map[int] Expectation{}
	eChan <-expects
	chanstruct := &Channel{channel,c,msgs.Chan,eChan}
	go handleExpects(chanstruct)
	return chanstruct, nil

}

func (c Channel) Write(b []byte) (int,error) {
	cmd := Command{Command: Privmsg,Prefix: c.conn.Nick,Params: []string{c.Name,string(b)}}
	return fmt.Fprintf(c.conn.conn,"%s\r\n",cmd.String())
}

func (c Channel) Expect(from string, body string) (ExpectChan,error) {
	return Expect(c.conn,Command{from,Privmsg,[]string{c.Name,body}})
}
