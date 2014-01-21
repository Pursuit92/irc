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
	"regexp"
	"github.com/Pursuit92/pubsub"
)

// CommandMatcher mirrors the structure of a Command, but with regular expressions
type CommandMatcher struct {
	Prefix, Command *regexp.Regexp
	Params          []*regexp.Regexp
}

func (cmd Command) MakeMatcher() (pubsub.Matcher,error) {
	var err error
	match := CommandMatcher{}
	match.Params = make([]*regexp.Regexp, len(cmd.Params))
	if len(cmd.Prefix) == 0 {
		match.Prefix = regexp.MustCompile(`.*`)
	} else {
		match.Prefix, err = regexp.Compile(cmd.Prefix)
		if err != nil {
			return match, err
		}
	}
	if len(cmd.Command) == 0 {
		match.Command = regexp.MustCompile(`.*`)
	} else {
		match.Command, err = regexp.Compile(cmd.Command)
		if err != nil {
			return match, err
		}
	}
	for i, v := range cmd.Params {
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

func (ce CmdErr) MakeMatcher() (pubsub.Matcher, error) {
	m, e := ce.Cmd.MakeMatcher()
	return m,e
}

func (mat CommandMatcher) Match(check pubsub.Matchable) bool {
	var com *Command
	var c,ce bool
	cmderr, ce := check.(CmdErr)
	if ce {
		if cmderr.Err != nil {
			return true
		}
		com = cmderr.Cmd
	println("com is",com)
	} else {
		cmd, c := check.(Command)
		if c {
			com = &cmd
	println("com is",com)
		}
	}
	if c || ce {
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
	return false
}
