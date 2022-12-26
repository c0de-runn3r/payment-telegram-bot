package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/c0de_runn3r/payments-telegram-bot/clients/telegram"
	storage "github.com/c0de_runn3r/payments-telegram-bot/files_storage"
	"github.com/c0de_runn3r/payments-telegram-bot/fsm"
)

func (p *Processor) doMessage(text string, chatID int, userID int64, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)
	currentState := fsm.SM.GetState(userID)
	if !CheckAdmin(userID) {
		switch text {
		case cmdStart:
			p.sendMessageWithReplyKB(chatID, msgHello, &StartKeyboard)
			usr := storage.DB.FindUser(userID)
			if usr.UserID == 0 {
				storage.DB.CreateNewUser(userID, chatID)
			}
		case cmdHelp, cmdQuestion:
			fsm.SM.SetState(userID, fsm.QuestionState)
			p.sendMessage(chatID, msgAskQuestion)
		case btnOneMonth:
			p.doInvoice(text, chatID, username)

		default:
			switch *currentState {
			case fsm.InitialState:
				p.sendMessageWithReplyKB(chatID, msgUnknownCommand, &StartKeyboard)
			case fsm.QuestionState:
				p.handleQuestion(text, chatID, userID, username)
				fsm.SM.SetState(userID, fsm.InitialState)
				p.sendMessageWithReplyKB(chatID, msgQuestionAccepted, &StartKeyboard)

			}
		}
	}
	if CheckAdmin(userID) {
		switch text {
		case cmdStart:
			p.sendMessageWithReplyKB(chatID, "Привіт\nЦе адмінська панель.", &AdminKeyboard)
		case cmdCurrentSubs:
			msg := "Активні підписки:\n\n"
			list := p.ListOfCurrentSubscribers()
			if len(list) > 0 {
				for i := 0; i < len(list); i++ {
					msg = msg + list[i] + "\n----------\n"
				}
			} else {
				msg = "Активних підписок зараз немає 😢"
			}
			p.sendMessageWithReplyKB(chatID, msg, &AdminKeyboard)
		default:
			switch *currentState {
			case fsm.InitialState:
				p.sendMessageWithReplyKB(chatID, msgUnknownCommand, &AdminKeyboard)
			case fsm.ReplyQuestionState:
				p.sendMessageWithReplyKB(fsm.CTX.Value("replyID").(int), text, &AdminKeyboard)
			}
		}
	}

	return nil
}

func (p *Processor) doCallbackQuerry(text string, chatID int, userID int64, username string) error {
	log.Printf("got new callback data '%s' from '%s'", text, username)
	switch text[0:3] {
	case "rpl":
		replyID, _ := strconv.Atoi(text[4:])
		fsm.CTX = context.WithValue(fsm.CTX, "replyID", replyID)
		fsm.SM.SetState(userID, fsm.ReplyQuestionState)
		p.sendMessage(chatID, msgSendAnswer)
	}
	return nil
}

func (p *Processor) doPreCheckoutQuery(invoiceID string, username string) error {
	log.Printf("got new precheckout data from '%s'", username)
	// check limitations - if there are product still avialable

	p.tg.AnswerPreCheckoutQuery(telegram.PreCheckoutParams{
		QueryID: invoiceID,
		OK:      true,
	})
	return nil
}

func (p *Processor) handleQuestion(text string, chatID int, userID int64, username string) {
	adminChat, _ := strconv.Atoi(os.Getenv("DEVELOPER_ID"))
	var msg string
	if username != "" {
		msg = fmt.Sprintf("Нове запитання від @%s:\n%s", username, text)
	} else {
		msg = fmt.Sprintf("Нове запитання (немає username):\n%s", text)
	}
	replyQuestionBut := &telegram.InlineKeyboardButton{
		Text:         "Відповісти",
		CallbackData: "rpl_" + strconv.Itoa(chatID),
	}
	replyQuestionKeyboard := &telegram.InlineKeyboardMarkup{
		Buttons: [][]telegram.InlineKeyboardButton{
			{*replyQuestionBut},
		},
	}
	p.sendMessageWithInlineKB(adminChat, msg, replyQuestionKeyboard)
}

func (p *Processor) doInvoice(text string, chatID int, username string) {
	var description string
	var label string
	var price int
	var payload string
	switch text {
	case btnOneMonth:
		description = "Підписка на 1 місяць"
		label = "Даша про харчування 🪄 - 1 місяць"
		price = 350 * 100
		payload = "1monthSub"
	}
	p.sendInvoice(chatID, username, "Даша про харчування 🪄", description, payload, label, price)
}

func (p *Processor) processPayment(chatID int, userID int64, username string, paymentDetails *telegram.SuccessfulPayment) error {
	product := fetchPayload(*paymentDetails)
	storage.DB.NewPaymentRecord(paymentDetails.TotalAmount, userID, username, product, paymentDetails.TelegramPaymentID, paymentDetails.ProviderPaymentID)
	link := os.Getenv("INVITE_LINK")
	msg := msgThnxForPayment + link
	p.sendMessageWithReplyKB(chatID, msg, &StartKeyboard)
	storage.DB.UpdateSubscriptionTime(userID, chatID, product)
	return nil
}

func (p *Processor) doJoinRequest(userID int64) error {
	user := storage.DB.FindUser(userID)
	isValid, _ := storage.DB.CheckSubscription(userID)
	if isValid {
		log.Println("approving chat joing request")
		if err := p.tg.ApproveChatJoinRequest(userID); err != nil {
			return err
		}
		p.sendMessage(user.ChatID, msgRequestApproved)
	}
	if !isValid {
		log.Println("declining chat joing request")
		p.sendMessage(user.ChatID, msgRequestDeclined)
	}
	return nil
}

func (p *Processor) sendInvoice(chatID int, username, title, description, payload, label string, price int) {
	params := telegram.InvoiceParams{
		ChatID:        chatID,
		Title:         title,
		Description:   description,
		Payload:       payload,
		ProviderToken: os.Getenv("PAYMENT_TOKEN"),
		Currency:      "UAH",
		Prices: &[]telegram.LabeledPrice{{
			Label:  label,
			Amount: price,
		}},
	}
	p.tg.SendInvoice(params)
}

func (p *Processor) sendMessage(chatID int, message string) {
	p.tg.SendMessage(telegram.MessageParams{
		ChatID: chatID,
		Text:   message,
	})
}

func (p *Processor) sendMessageWithReplyKB(chatID int, message string, keyboard *telegram.ReplyKeyboardMarkup) {
	p.tg.SendMessage(telegram.MessageParams{
		ChatID:        chatID,
		Text:          message,
		KeyboardReply: keyboard,
	})
}

func (p *Processor) sendMessageWithInlineKB(chatID int, message string, keyboard *telegram.InlineKeyboardMarkup) {
	p.tg.SendMessage(telegram.MessageParams{
		ChatID:         chatID,
		Text:           message,
		KeyboardInline: keyboard,
	})
}

func fetchPayload(paymentDetails telegram.SuccessfulPayment) (product string) {
	payload := strings.Split(paymentDetails.Payload, " ")
	return payload[0] // product
}

func (p *Processor) EveryHourCheck() {
	for {
		log.Println("checking user's subscriptions")
		users := storage.DB.GetAllUsers()
		for i := 0; i < len(users); i++ {
			isValid, timeTill := storage.DB.CheckSubscription(users[i].UserID)
			if !isValid {
				if users[i].WarningMessage != "ended" {
					log.Printf("removing user [%v] from channel", users[i].UserID)
					p.tg.BanUser(users[i].UserID)
					p.tg.UnbanUser(users[i].UserID)
					p.sendMessageWithReplyKB(users[i].ChatID, msgSubscriptionEnded, &StartKeyboard)
					storage.DB.Table("users").Where("user_id = ?", users[i].UserID).Updates(storage.User{WarningMessage: "ended"})
					continue
				}
			}
			if timeTill.Before(time.Now().UTC().AddDate(0, 0, 3)) { // less than 3 days
				if users[i].WarningMessage != "ended" && users[i].WarningMessage != "3days" {
					p.sendMessageWithReplyKB(users[i].ChatID, msgUpdateSubscription, &StartKeyboard)
					storage.DB.Table("users").Where("user_id = ?", users[i].UserID).Updates(storage.User{WarningMessage: "3days"})
					continue
				}
			}
		}
		time.Sleep(time.Hour * 1)
	}
}

func (p *Processor) ListOfCurrentSubscribers() []string {
	list := make([]string, 0)
	users := storage.DB.GetAllUsers()
	for i := 0; i < len(users); i++ {
		valid, till := storage.DB.CheckSubscription(users[i].UserID)
		if valid {
			userData, _ := p.tg.GetUser(strconv.Itoa(users[i].ChatID))
			if userData.Result.Username == "" {
				userStr := fmt.Sprintf("%s %s - до %s", userData.Result.FirstName, userData.Result.LastName, till.Local().Format("02-01-2006"))
				list = append(list, userStr)
			} else {
				userStr := fmt.Sprintf("%s %s [@%s] - до %s", userData.Result.FirstName, userData.Result.LastName, userData.Result.Username, till.Local().Format("02-01-2006"))
				list = append(list, userStr)
			}
		}
	}
	return list

}

func CheckAdmin(userID int64) bool {
	admin := os.Getenv("ADMIN_ID")
	devloper := os.Getenv("DEVELOPER_ID")
	if strconv.FormatInt(userID, 10) == admin || strconv.FormatInt(userID, 10) == devloper {
		return true
	}
	return false
}
