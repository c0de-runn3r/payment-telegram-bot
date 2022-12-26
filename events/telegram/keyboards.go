package telegram

import (
	"github.com/c0de_runn3r/payments-telegram-bot/clients/telegram"
)

var (
	oneMonthBut = &telegram.KeyboardButton{Text: btnOneMonth}
	questionBut = &telegram.KeyboardButton{Text: cmdQuestion}
	checkSubs   = &telegram.KeyboardButton{Text: cmdCurrentSubs}
)

var StartKeyboard = telegram.ReplyKeyboardMarkup{
	Keyboard: [][]telegram.KeyboardButton{
		{*oneMonthBut},
		{*questionBut},
	},
	ResizeKeyboard: true,
}

var AdminKeyboard = telegram.ReplyKeyboardMarkup{
	Keyboard: [][]telegram.KeyboardButton{
		{*checkSubs},
	},
	ResizeKeyboard:  true,
	OneTimeKeyboard: true,
}
