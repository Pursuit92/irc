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
	"fmt"
	"regexp"
	"strings"
	"github.com/Pursuit92/pubsub"
)

type Channel struct {
	Name string
	conn *Conn
	pubsub.Publisher
}

func (c *Conn) Join(channel string) (chanstruct *Channel, err error) {
	joinCmd := Command{Command: Join, Prefix: c.Nick, Params: []string{channel}}
	err = c.Send(joinCmd)
	if err == nil {
		msgs, _ := c.Subscribe(Command{"", "", []string{channel}})
		chanstruct = &Channel{channel, c, pubsub.MakePublisher(msgs.Chan)}
		join, _ := chanstruct.Subscribe(Command{"", "JOIN", []string{}})
		defer chanstruct.UnSubscribe(join)
		<-join.Chan
	}
	return chanstruct, err
}

func (c Channel) Write(b []byte) (int, error) {
	cmd := Command{Command: Privmsg, Prefix: c.conn.Nick, Params: []string{c.Name, string(b)}}
	return fmt.Fprintf(c.conn.conn, "%s\r\n", cmd.String())
}

func parseWhoReply(cmd *Command) IRCUser {
	if cmd.Command != RplWhoreply {
		return IRCUser{}
	}
	whoreplReg := regexp.MustCompile(`^(?P<name>[^ ]+) (?P<host>[^ ]+) (?P<server>[^ ]+) (?P<nick>[^ ]+) (?:[^ ]+ ){2}(?P<realname>.*)$`)
	var user IRCUser
	cmdStr := strings.Join(cmd.Params[2:], " ")
	user.Nick = whoreplReg.ReplaceAllString(cmdStr, "${nick}")
	user.Host = whoreplReg.ReplaceAllString(cmdStr, "${host}")
	user.Server = whoreplReg.ReplaceAllString(cmdStr, "${server}")
	user.Name = whoreplReg.ReplaceAllString(cmdStr, "${name}")
	user.RealName = whoreplReg.ReplaceAllString(cmdStr, "${realname}")
	return user
}

func (c Channel) GetUsers() (map[string]IRCUser,error) {
	users := make(map[string]IRCUser)
	userMsgs, _ := c.conn.Subscribe(Command{"", RplWhoreply, []string{}})
	userEnd, _ := c.conn.Subscribe(Command{"", RplEndofwho, []string{}})
	defer c.conn.UnSubscribe(userMsgs)
	defer c.conn.UnSubscribe(userEnd)
	whoCmd := Command{Command: Who, Params: []string{c.Name}}
	c.conn.Send(whoCmd)
	// BUG(Josh) What if stuff happens before the userEnd Chan talks?
	for {
		select {
		case match := <-userMsgs.Chan:
			msg := match.(CmdErr)
			if msg.Err != nil {
				return nil, msg.Err
			}
			user := parseWhoReply(msg.Cmd)
			if user.Nick != "" {
				users[user.Nick] = user
			}
		case <-userEnd.Chan:
			return users,nil
		}
	}
}
