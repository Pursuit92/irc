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
	"bytes"
	"fmt"
	"strings"
)

type Command struct {
	Prefix  string
	Command string
	Params  []string
}

type CmdErr struct {
	Cmd *Command
	Err error
}

type Message struct {
	Nick string
	User string
	Host string
	Body []string
}

func parseCommand(message string) (*Command, error) {
	var err error

	parts := strings.Split(message, " ")

	if len(parts) == 0 {
		return nil, err
	}

	i := 0

	var prefix string
	var cmd string

	if parts[i][0] != ':' {
		prefix = ""
	} else {
		prefix = parts[i][1:]
		i++
	}

	cmd = parts[i]
	i++
	var params []string = make([]string, len(parts[i:]))
	for j, v := range parts[i:] {
		if v[0] == ':' {
			// All but the first and last (colon and CR)
			params[j] = dropCRLF(strings.Join(parts[i+j:], " "))[1:]
			params = params[:j+1]
			//print(params[j])
			break
		} else if j == len(parts[i:])-1 {
			// All but the last (CR)
			params[j] = dropCRLF(v)
		} else {
			params[j] = v
		}
	}
	return &Command{prefix, cmd, params}, nil
}

func dropCRLF(s string) string {
	for i, v := range s {
		if v == '\r' || v == '\n' {
			return s[:i]
		}
	}
	return s
}

func (m Command) String() string {
	var prefix string
	var body string = ""
	var buf bytes.Buffer
	if len(m.Prefix) == 0 {
		prefix = ""
	} else {
		prefix = ":" + m.Prefix + " "
	}

	//log.Lprint("Params:")
	for i, v := range m.Params {
		//log.Lprintf("\t\t%s",v)
		if i == len(m.Params)-1 {
			body += ":" + v
		} else {
			body += v + " "
		}
	}

	fmt.Fprintf(&buf, "%s%s %s", prefix, m.Command, body)

	return buf.String()
}

func (c Command) Message() Message {
	n, u, h := splitFrom(c.Prefix)
	return Message{n, u, h, c.Params}
}

func splitFrom(from string) (nick, user, host string) {
	var i int
	var j int
	for i = 0; from[i] != '!'; i++ {
	}
	nick = from[:i]
	for j = i + 1; from[j] != '@'; j++ {
	}
	user = from[i+1 : j]
	host = from[j+1:]
	return
}
