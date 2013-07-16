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

type Message struct {
	From string
	To   string
	body string
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

	//log.Print("Params:")
	for i, v := range m.Params {
		//log.Printf("\t\t%s",v)
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
	return Message{c.Prefix, c.Params[0], strings.Join(c.Params, " ")}
}
