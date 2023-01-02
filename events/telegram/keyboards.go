package telegram

import (
	"strconv"

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

func makeConfirmPaymentKeyboard(userID int64) *telegram.InlineKeyboardMarkup {
	confirmPaymentBut := &telegram.InlineKeyboardButton{
		Text:         "Підтвердити оплату ✅",
		CallbackData: "sub_" + strconv.FormatInt(userID, 10),
	}
	confirmPaymentKeyboard := &telegram.InlineKeyboardMarkup{
		Buttons: [][]telegram.InlineKeyboardButton{
			{*confirmPaymentBut},
		},
	}
	return confirmPaymentKeyboard
}

func makeReplyQuestionKeyboard(chatID int) *telegram.InlineKeyboardMarkup {
	replyQuestionBut := &telegram.InlineKeyboardButton{
		Text:         "Відповісти",
		CallbackData: "rpl_" + strconv.Itoa(chatID),
	}
	replyQuestionKeyboard := &telegram.InlineKeyboardMarkup{
		Buttons: [][]telegram.InlineKeyboardButton{
			{*replyQuestionBut},
		},
	}
	return replyQuestionKeyboard
}
