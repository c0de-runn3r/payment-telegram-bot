package telegram

import (
	"github.com/c0de_runn3r/payments-telegram-bot/clients/telegram"
)

var (
	oneMonth        = &telegram.KeyboardButton{Text: btnOneMonth}
	threeMonths     = &telegram.KeyboardButton{Text: btnThreeMonths}
	sixMonths       = &telegram.KeyboardButton{Text: btnSixMonths}
	buySubscription = &telegram.KeyboardButton{Text: btnBuySubscription}
)

var StartKeyboard = telegram.ReplyKeyboardMarkup{
	Keyboard: [][]telegram.KeyboardButton{
		{*buySubscription},
	},
	ResizeKeyboard:  true,
	OneTimeKeyboard: true,
}

var PricesKeyboard = telegram.ReplyKeyboardMarkup{
	Keyboard: [][]telegram.KeyboardButton{
		{*oneMonth},
		{*threeMonths},
		{*sixMonths},
	},
	ResizeKeyboard:  true,
	OneTimeKeyboard: true,
}
