package telegram

import (
	"errors"
	"fmt"

	"github.com/c0de_runn3r/payments-telegram-bot/clients/telegram"
	"github.com/c0de_runn3r/payments-telegram-bot/events"
)

type Processor struct {
	tg     *telegram.Client
	offset int
}

type Meta struct {
	ChatID         int
	Username       string
	UserID         int64
	InvoiceID      string
	PaymentDetails *telegram.SuccessfulPayment
}

var (
	ErrUknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client) *Processor {
	return &Processor{
		tg: client,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events %w", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		if u.Message == nil && u.CallbackQuery == nil && u.PreCheckoutQuery == nil && u.JoinRequest == nil {
			continue
		}
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.CallbackQuery:
		return p.processCallbackQuery(event)
	case events.PreCheckoutQuery:
		return p.processPreCheckoutQuery(event)
	case events.JoinRequest:
		return p.processJoinRequest(event)
	default:
		return fmt.Errorf("can't process message %w", ErrUknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message %w", err)
	}

	photo, _ := photo(event)

	if meta.PaymentDetails != nil {
		if err := p.processPayment(meta.ChatID, meta.UserID, meta.Username, meta.PaymentDetails); err != nil {
			return fmt.Errorf("can't process payment %w", err)
		}
	} else if event.Photo != nil {
		if err := p.processPhoto(meta.ChatID, meta.UserID, meta.Username, photo); err != nil {
			return fmt.Errorf("can't process photo %w", err)
		}
	} else {
		if err := p.doMessage(event.Text, meta.ChatID, meta.UserID, meta.Username); err != nil {
			return fmt.Errorf("can't process message %w", err)
		}
	}
	return nil
}

func (p *Processor) processCallbackQuery(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process callback %w", err)
	}
	if err := p.doCallbackQuerry(event.Text, event.ID, meta.ChatID, meta.UserID, meta.Username); err != nil {
		return fmt.Errorf("can't process callback %w", err)
	}
	return nil
}

func (p *Processor) processJoinRequest(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process join request %w", err)
	}
	if err := p.doJoinRequest(meta.UserID); err != nil {
		return fmt.Errorf("can't process join request %w", err)
	}
	return nil
}

func (p *Processor) processPreCheckoutQuery(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process precheckout event %w", err)
	}
	if err := p.doPreCheckoutQuery(meta.InvoiceID, meta.Username); err != nil {
		return fmt.Errorf("can't process precheckout event %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta %w", ErrUnknownMetaType)
	}
	return res, nil
}

func photo(event events.Event) (telegram.Photo, bool) {
	res, ok := event.Photo.(telegram.Photo)
	return res, ok
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
		ID:   fetchMessageID(upd),
	}

	if updType == events.Message {
		if upd.Message.SuccessfulPayment != nil {
			res.Meta = Meta{
				ChatID:         upd.Message.Chat.ID,
				Username:       upd.Message.From.Username,
				UserID:         upd.Message.From.ID,
				PaymentDetails: upd.Message.SuccessfulPayment,
			}
		} else if upd.Message.Photo != nil {
			res.Meta = Meta{
				ChatID:   upd.Message.Chat.ID,
				Username: upd.Message.From.Username,
				UserID:   upd.Message.From.ID,
			}
			res.Photo = telegram.Photo{
				ID:       upd.Message.Photo[0].ID,
				UniqueID: upd.Message.Photo[0].UniqueID,
				Width:    upd.Message.Photo[0].Width,
				Height:   upd.Message.Photo[0].Height,
			}
		} else {
			res.Meta = Meta{
				ChatID:   upd.Message.Chat.ID,
				Username: upd.Message.From.Username,
				UserID:   upd.Message.From.ID,
			}
		}
	}
	if updType == events.CallbackQuery {
		res.Meta = Meta{
			ChatID:   upd.CallbackQuery.Message.Chat.ID,
			Username: upd.CallbackQuery.From.Username,
			UserID:   upd.CallbackQuery.From.ID,
		}
	}
	if updType == events.PreCheckoutQuery {
		res.Meta = Meta{
			Username:  upd.PreCheckoutQuery.From.Username,
			InvoiceID: upd.PreCheckoutQuery.ID,
		}
	}
	if updType == events.JoinRequest {
		res.Meta = Meta{
			UserID: upd.JoinRequest.From.ID,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Data
	}
	return ""
}

func fetchMessageID(upd telegram.Update) int {
	if upd.Message != nil {
		return upd.Message.MessageID
	}
	if upd.CallbackQuery != nil {
		return upd.CallbackQuery.Message.MessageID
	}
	return -1
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}
	if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}
	if upd.PreCheckoutQuery != nil {
		return events.PreCheckoutQuery
	}
	if upd.JoinRequest != nil {
		return events.JoinRequest
	}
	return events.Unknown

}
