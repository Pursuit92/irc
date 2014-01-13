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
	"github.com/Pursuit92/LeveledLogger/log"
	"github.com/Pursuit92/syncmap"
	"regexp"
)

// CommandMatcher mirrors the structure of a Command, but with regular expressions
type CommandMatcher struct {
	Prefix, Command *regexp.Regexp
	Params          []*regexp.Regexp
}

type ExpectChan struct {
	id   int
	Chan chan CmdErr
}

type Expectation struct {
	CommandMatcher
	ExpectChan
}

type Expector struct {
	msgs    chan CmdErr
	// Map of ints to Expectations
	expects syncmap.Map
}

// Turns a channel into an Expector
func MakeExpector(msgs chan CmdErr) Expector {
	eMap := syncmap.New()
	exp := Expector{msgs, eMap}
	go exp.handleExpects()
	return exp
}

func CompileMatcher(cmdRE Command) (match CommandMatcher, err error) {
	match.Params = make([]*regexp.Regexp, len(cmdRE.Params))
	if len(cmdRE.Prefix) == 0 {
		match.Prefix = regexp.MustCompile(`.*`)
	} else {
		match.Prefix, err = regexp.Compile(cmdRE.Prefix)
		if err != nil {
			return match, err
		}
	}
	if len(cmdRE.Command) == 0 {
		match.Command = regexp.MustCompile(`.*`)
	} else {
		match.Command, err = regexp.Compile(cmdRE.Command)
		if err != nil {
			return match, err
		}
	}
	for i, v := range cmdRE.Params {
		if len(v) == 0 {
			match.Params[i] = regexp.MustCompile(`.*`)
		} else {
			match.Params[i], err = regexp.Compile(v)
			if err != nil {
				return match, err
			}
		}
	}
	return match, err
}

// Register a channel to receive messages matching a specific pattern
func (exp Expector) Expect(cr Command) (eChan ExpectChan, err error) {
	log.Out.Printf(3,"Registering Expect for %s\n", cr.String())
	var exists bool
	var i int
	var match Expectation
	match.CommandMatcher, err = CompileMatcher(cr)
	if err != nil {
		return eChan, err
	}
	eMap := exp.expects
	exists = true
	for exists {
		i = rgen.Intn(65534)
		_, exists = eMap.Get(i)
	}
	c := make(chan CmdErr)
	eChan.Chan = c
	eChan.id = i
	match.ExpectChan = eChan
	eMap.Set(i,match)
	log.Out.Printf(3,"Expect id: %d\n", i)
	return eChan, nil
}

func (exp Expector) UnExpect(e ExpectChan) {
	log.Out.Printf(3,"Removing expect with id %d", e.id)
	eMap := exp.expects
	close(e.Chan)
	eMap.Delete(e.id)
}

func (c Expector) handleExpects() {
	log.Out.Printf(3,"Starting Expect handler")
	msgOut := c.msgs
	eMap := c.expects
	for msg := range msgOut {
		//println("expect handler got message")
		//log.Out.Printf("Testing message: %s",msg.String())
		for _, v := range eMap.Map() {
			w := v.(Expectation)
			if msg.Err != nil {
				w.Chan <- msg
				close(w.Chan)
				continue
			}
			if matchCommand(msg.Cmd, w.CommandMatcher) {
				log.Out.Printf(3,"Sending message to Expect channel with id %d: %s", w.id, msg.Cmd.String())
				w.Chan <- msg
			}
		}
	}
}

func matchCommand(com Command, mat CommandMatcher) bool {
	if len(com.Params) < len(mat.Params) {
		return false
	}
	if mat.Prefix.MatchString(com.Prefix) == false {
		return false
	}
	if mat.Command.MatchString(com.Command) == false {
		return false
	}
	for i, v := range mat.Params {
		if v.MatchString(com.Params[i]) == false {
			return false
		}
	}
	return true
}

func (c Expector) DefaultExpect() ExpectChan {
	var match Expectation

	match.CommandMatcher, _ = CompileMatcher(Command{})

	eMap := c.expects
	i := 65535
	ch := make(chan CmdErr)
	match.Chan = ch
	match.id = i
	eMap.Set(i,match)
	return ExpectChan{i, ch}
}
