package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

//go:generate plumb -output=plumb_gen.go

// A Phrase is raw greeting text. Nothing in the set produces one, so plumb
// treats it as an input to the generated build function.
type Phrase string

// A Message is a greeting ready to show.
type Message string

// NewMessage turns a Phrase into a Message.
//
//plumb:build
func NewMessage(phrase Phrase) Message {
	return Message(phrase)
}

// A Greeter greets with a Message, unless it woke up grumpy.
type Greeter struct {
	Message Message
	Grumpy  bool
}

// NewGreeter builds a Greeter. About half the time it is grumpy.
//
//plumb:build
func NewGreeter(m Message) Greeter {
	return Greeter{Message: m, Grumpy: rand.IntN(2) == 0}
}

// Greet returns the greeting to show.
func (g Greeter) Greet() Message {
	if g.Grumpy {
		return "Go away!"
	}
	return g.Message
}

// An Event is a greeting in progress.
type Event struct {
	Greeter Greeter
}

// NewEvent builds an Event. A grumpy greeter cancels it, so this can fail. That
// is why the generated build returns an error.
//
//plumb:build
func NewEvent(g Greeter) (Event, error) {
	if g.Grumpy {
		return Event{}, errors.New("the greeter is grumpy; no event today")
	}
	return Event{Greeter: g}, nil
}

// Start runs the event.
func (e Event) Start() {
	fmt.Println(e.Greeter.Greet())
}
