package irc

import (
	"github.com/Pursuit92/LeveledLogger/log"
)

func SetLogLevel(n int) {
	log.Out.SetLevel(n)
}
