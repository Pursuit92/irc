package irc

import (
	"fmt"
)

type Channel struct {
	Name	string
	conn *Conn
}

func (c *Conn) Join(channel string) (*Channel, error) {
	joinCmd := Command{Command: Join, Prefix: c.Nick, Params: []string{channel}}
	c.sendCommand(joinCmd)
	return &Channel{channel,c}, nil

}

func (c Channel) Write(b []byte) (int,error) {
	cmd := Command{Command: Privmsg,Prefix: c.conn.Nick,Params: []string{c.Name,string(b)}}
	return fmt.Fprintf(c.conn.conn,"%s\r\n",cmd.String())
}

func (c Channel) Expect(from string, body string) (ExpectChan,error) {
	return c.conn.Expect(Command{from,Privmsg,[]string{c.Name,body}})
}
