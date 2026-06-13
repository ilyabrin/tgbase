package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"tgbase/internal/fsm"

	"gopkg.in/telebot.v3"
)

// FSM states for the registration flow.
const (
	StateAskName = "register:ask_name"
	StateAskAge  = "register:ask_age"
)

// RegisterStart handles /register - kicks off the flow.
func RegisterStart(f *fsm.FSM) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if err := f.SetState(c, StateAskName); err != nil {
			return err
		}
		return c.Send("Привет! Как тебя зовут?")
	}
}

// RegisterAskName collects the name, moves to the next step.
func RegisterAskName(f *fsm.FSM) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		name := strings.TrimSpace(c.Text())
		if name == "" {
			return c.Send("Имя не может быть пустым. Попробуй ещё раз:")
		}

		// Store the name in telebot context storage for the next step.
		c.Set("reg_name", name)

		// telebot context is per-update, so persist via FSM data key in Redis.
		// For simplicity we encode name into the state value itself.
		if err := f.SetStateData(c, StateAskAge, name); err != nil {
			return err
		}
		return c.Send(fmt.Sprintf("Отлично, %s! Сколько тебе лет?", name))
	}
}

// RegisterAskAge collects the age, completes the flow.
func RegisterAskAge(f *fsm.FSM) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		age, err := strconv.Atoi(strings.TrimSpace(c.Text()))
		if err != nil || age < 1 || age > 150 {
			return c.Send("Введи корректный возраст (число от 1 до 150):")
		}

		name, _ := f.GetData(c)

		if err := f.ClearState(c); err != nil {
			return err
		}

		return c.Send(fmt.Sprintf(
			"Регистрация завершена!\nИмя: %s\nВозраст: %d",
			name, age,
		))
	}
}
