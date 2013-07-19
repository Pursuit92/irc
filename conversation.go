package irc

import "io"
type Conversation interface {
	Expectable
	io.Writer
}
