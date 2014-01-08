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

type CommandMatcher struct {
	Prefix, Command *regexp.Regexp
	Params          []*regexp.Regexp
}

type ExpectChan struct {
	id   int
	Chan chan Command
}

type Expectation struct {
	CommandMatcher
	ExpectChan
}

type Expectable interface {
	Expects() syncmap.Map
	MsgOut() chan Command
}

type Expector struct {
	msgs    chan Command
	expects syncmap.Map
}

func (e Expector) Expects() syncmap.Map {
	return e.expects
}
func (e Expector) MsgOut() chan Command {
	return e.msgs
}

// Turns a channel into an Expector
func MakeExpector(msgs chan Command) Expector {
	eMap := syncmap.New()
	return Expector{msgs, eMap}
}

// Register a channel to receive messages matching a specific pattern
func Expect(irc Expectable, cr Command) (ExpectChan, error) {
	log.Out.Printf(3,"Registering Expect for %s\n", cr.String())
	var exists bool
	var i int
	var match Expectation
	var err error
	match.Params = make([]*regexp.Regexp, len(cr.Params))
	if len(cr.Prefix) == 0 {
		match.Prefix = regexp.MustCompile(`.*`)
	} else {
		match.Prefix, err = regexp.Compile(cr.Prefix)
		if err != nil {
			return ExpectChan{}, err
		}
	}
	if len(cr.Command) == 0 {
		match.Command = regexp.MustCompile(`.*`)
	} else {
		match.Command, err = regexp.Compile(cr.Command)
		if err != nil {
			return ExpectChan{}, err
		}
	}
	for i, v := range cr.Params {
		match.Params[i], err = regexp.Compile(v)
		if err != nil {
			return ExpectChan{}, err
		}
	}
	eMap := irc.Expects()
	exists = true
	for exists {
		i = rgen.Intn(65534)
		_, exists = eMap.Get(i)
	}
	c := make(chan Command)
	match.Chan = c
	match.id = i
	eMap.Set(i,match)
	log.Out.Printf(3,"Expect id: %d\n", i)
	return ExpectChan{i, c}, nil
}

func UnExpect(irc Expectable, e ExpectChan) {
	log.Out.Printf(3,"Removing expect with id %d", e.id)
	eMap := irc.Expects()
	eMap.Delete(e.id)
}

func handleExpects(c Expectable) {
	log.Out.Printf(3,"Starting Expect handler")
	msgOut := c.MsgOut()
	eMap := c.Expects()
	for msg := range msgOut {
		sent := false
		//println("expect handler got message")
		//log.Out.Printf("Testing message: %s",msg.String())
		for _, v := range eMap.Map() {
			w := v.(Expectation)
			if matchCommand(msg, w.CommandMatcher) {
				log.Out.Printf(3,"Sending message to Expect channel with id %d: %s", w.id, msg.String())
				w.Chan <- msg
				sent = true
			}
		}
		def, ok := eMap.Get(65535)
		if ok && !sent {
			d := def.(Expectation)
			log.Out.Print(3,"Sending message to default channel")
			d.Chan <- msg
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

func compileCmdMatch(cr Command) (CommandMatcher, error) {
	var match CommandMatcher
	var err error
	match.Params = make([]*regexp.Regexp, len(cr.Params))
	if len(cr.Prefix) == 0 {
		match.Prefix = regexp.MustCompilePOSIX(`.*`)
	} else {
		match.Prefix, err = regexp.CompilePOSIX(cr.Prefix)
		if err != nil {
			return CommandMatcher{}, err
		}
	}
	if len(cr.Command) == 0 {
		match.Command = regexp.MustCompilePOSIX(`.*`)
	} else {
		match.Command, err = regexp.CompilePOSIX(cr.Command)
		if err != nil {
			return CommandMatcher{}, err
		}
	}
	for i, v := range cr.Params {
		match.Params[i], err = regexp.CompilePOSIX(v)
		if err != nil {
			return CommandMatcher{}, err
		}
	}
	return match, nil
}

func DefaultExpect(c Expectable) ExpectChan {
	var match Expectation

	eMap := c.Expects()
	i := 65535
	ch := make(chan Command)
	match.Chan = ch
	match.id = i
	eMap.Set(i,match)
	return ExpectChan{i, ch}
}
