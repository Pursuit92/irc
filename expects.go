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

type Expectation struct {
	CommandMatcher
	ExpectChan
}

type Expectable interface {
	Expects() chan map[int]Expectation
	MsgOut() chan Command
}

type Expector struct {
	msgs    chan Command
	expects chan map[int]Expectation
}

func (e Expector) Expects() chan map[int]Expectation {
	return e.expects
}
func (e Expector) MsgOut() chan Command {
	return e.msgs
}

func MakeExpector(msgs chan Command) Expector {
	eChan := make(chan map[int]Expectation, 1)
	expects := map[int]Expectation{}
	eChan <- expects
	return Expector{msgs, eChan}
}

// Register a channel to receive messages of a certain type
func Expect(irc Expectable, cr Command) (ExpectChan, error) {
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
	eChan := irc.Expects()
	expects := <-eChan
	exists = true
	for exists {
		i = rgen.Intn(65534)
		_, exists = expects[i]
	}
	c := make(chan Command)
	match.Chan = c
	match.id = i
	expects[i] = match
	eChan <- expects
	return ExpectChan{i, c}, nil
}

func UnExpect(irc Expectable, e ExpectChan) error {
	eChan := irc.Expects()
	expects := <-eChan
	_, exists := expects[e.id]
	if !exists {
		eChan <- expects
		return IRCErr("Expect does not exist!")
	}
	delete(expects, e.id)
	eChan <- expects
	return nil
}

func handleExpects(c Expectable) {
	log.Printf("Starting Expect handler")
	msgOut := c.MsgOut()
	eChan := c.Expects()
	for {
		msg := <-msgOut
		//println("expect handler got message")
		expects := <-eChan
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
		eChan <- expects
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

	eChan := c.Expects()
	expects := <-eChan
	i := 65535
	ch := make(chan Command)
	match.Chan = ch
	match.id = i
	expects[i] = match
	eChan <- expects
	return ExpectChan{i, ch}
}
