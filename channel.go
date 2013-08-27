package irc

import (
	"fmt"
)

type Channel struct {
	Name string
	conn *Conn
	Expector
}

func (c *Conn) Join(channel string) (*Channel, error) {
	joinCmd := Command{Command: Join, Prefix: c.Nick, Params: []string{channel}}
	c.sendCommand(joinCmd)
	msgs, _ := Expect(c, Command{"", "PRIVMSG", []string{channel}})
	chanstruct := &Channel{channel, c, MakeExpector(msgs.Chan)}
	go handleExpects(chanstruct)
	return chanstruct, nil
}

func (c Channel) Write(b []byte) (int, error) {
	cmd := Command{Command: Privmsg, Prefix: c.conn.Nick, Params: []string{c.Name, string(b)}}
	return fmt.Fprintf(c.conn.conn, "%s\r\n", cmd.String())
}
