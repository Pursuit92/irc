package irc

import (
	"log"
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

type Expect struct {
	CommandMatcher
	ExpectChan
}

// Register a channel to receive messages of a certain type
func (irc Conn) Expect(cr Command) (ExpectChan, error) {
	var exists bool
	var i int
	var match Expect
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
	expects := <-irc.expects
	exists = true
	for exists {
		i = rgen.Intn(65534)
		_, exists = expects[i]
	}
	c := make(chan Command)
	match.Chan = c
	match.id = i
	expects[i] = match
	irc.expects <- expects
	return ExpectChan{i, c}, nil
}

func (irc Conn) UnExpect(e ExpectChan) error {
	expects := <-irc.expects
	_, exists := expects[e.id]
	if !exists {
		irc.expects <- expects
		return IRCErr("Expect does not exist!")
	}
	delete(expects, e.id)
	irc.expects <- expects
	return nil
}

func (c Conn) handleExpects() {
	log.Printf("Starting Expect handler")
	for {
		msg := <-c.msgOut
		//println("expect handler got message")
		expects := <-c.expects
		for _, v := range expects {
			if matchCommand(msg, v.CommandMatcher) {
				log.Print("Sending message to Expect channel")
				v.Chan <- msg
			}
		}
		d, ok := expects[65535]
		if ok {
			log.Print("Sending message to default channel")
			d.Chan <- msg
		}
		c.expects <- expects
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

func (c Conn) DefaultExpect() ExpectChan {
	var match Expect

	expects := <-c.expects
	i := 65535
	ch := make(chan Command)
	match.Chan = ch
	match.id = i
	expects[i] = match
	c.expects <- expects
	return ExpectChan{i, ch}
}
