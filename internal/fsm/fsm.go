// Package fsm provides conversation state management for Telegram bots.
//
// Typical usage:
//
//	storage := fsm.NewRedisStorage(redisClient, fsm.WithTTL(3600))
//	f := fsm.New(storage)
//
//	b.Handle("/start", func(c telebot.Context) error {
//	    f.SetState(c, "ask_name")
//	    return c.Send("What's your name?")
//	})
//
//	b.Handle(telebot.OnText, f.Route(
//	    fsm.On("ask_name", func(c telebot.Context) error {
//	        f.SetState(c, "ask_age")
//	        return c.Send("How old are you?")
//	    }),
//	    fsm.On("ask_age", func(c telebot.Context) error {
//	        f.ClearState(c)
//	        return c.Send("Done!")
//	    }),
//	))
package fsm

import (
	"context"

	"gopkg.in/telebot.v3"
)

// State is a named step in the conversation flow.
type State = string

// None is the zero state — user has no active flow.
const None State = ""

// FSM routes Telegram messages based on per-user conversation state.
type FSM struct {
	storage  Storage
	fallback telebot.HandlerFunc
}

// New creates an FSM backed by the given storage.
func New(storage Storage) *FSM {
	return &FSM{storage: storage}
}

// Fallback sets a handler called when the user's state matches no registered step.
// Without a fallback, unmatched messages are silently ignored.
func (f *FSM) Fallback(h telebot.HandlerFunc) *FSM {
	f.fallback = h
	return f
}

// Step pairs a state name with its handler.
type Step struct {
	State   State
	Handler telebot.HandlerFunc
}

// On creates a Step for use in Route.
func On(state State, h telebot.HandlerFunc) Step {
	return Step{State: state, Handler: h}
}

// Route returns a telebot.HandlerFunc that dispatches to the matching Step
// based on the sender's current state. Unmatched states call the fallback if set.
//
// Typical registration:
//
//	b.Handle(telebot.OnText, f.Route(
//	    fsm.On("ask_name", nameHandler),
//	    fsm.On("ask_age",  ageHandler),
//	))
func (f *FSM) Route(steps ...Step) telebot.HandlerFunc {
	table := make(map[State]telebot.HandlerFunc, len(steps))
	for _, s := range steps {
		table[s.State] = s.Handler
	}

	return func(c telebot.Context) error {
		state, err := f.storage.Get(context.Background(), c.Sender().ID)
		if err != nil {
			return err
		}
		if h, ok := table[state]; ok {
			return h(c)
		}
		if f.fallback != nil {
			return f.fallback(c)
		}
		return nil
	}
}

// SetState transitions the sender to a new state.
func (f *FSM) SetState(c telebot.Context, state State) error {
	return f.storage.Set(context.Background(), c.Sender().ID, state)
}

// GetState returns the sender's current state.
func (f *FSM) GetState(c telebot.Context) (State, error) {
	return f.storage.Get(context.Background(), c.Sender().ID)
}

// ClearState removes the sender's state, returning them to None.
func (f *FSM) ClearState(c telebot.Context) error {
	return f.storage.Clear(context.Background(), c.Sender().ID)
}
