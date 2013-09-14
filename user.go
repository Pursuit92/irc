package irc

import "fmt"

type IRCUser struct {
	Nick     string
	Host     string
	Server   string
	Name     string
	RealName string
}

func (u IRCUser) String() string {
	return fmt.Sprintf("%s!%s@%s", u.Nick, u.Name, u.Host)
}
