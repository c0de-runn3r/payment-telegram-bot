package telegram

type UpdateResponce struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID               int               `json:"update_id"`
	Message          *IncomingMessage  `json:"message"`
	CallbackQuery    *CallbackQuery    `json:"callback_query"`
	PreCheckoutQuery *PreCheckoutQuery `json:"pre_checkout_query"`
	JoinRequest      *JoinRequest      `json:"chat_join_request"`
}

type IncomingMessage struct {
	Text              string             `json:"text"`
	From              User               `json:"from"`
	Chat              Chat               `json:"chat"`
	Invoice           *Invoice           `json:"invoice"`
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment"`
}

type CallbackQuery struct {
	ID      int              `json:"update_id"`
	From    User             `json:"from"`
	Message *IncomingMessage `json:"message"`
	Data    string           `json:"data"`
	Chat    string           `json:"chat_instance"`
}

type JoinRequest struct {
	Chat Chat `json:"chat"`
	From User `json:"from"`
	Date int  `json:"date"`
}

type PreCheckoutQuery struct {
	ID             string `json:"id"`
	From           User   `json:"from"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
	InvoicePayload string `json:"invoice_payload"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

type SuccessfulPayment struct {
	Currency          string `json:"currency"`
	TotalAmount       int    `json:"total_amount"`
	Payload           string `json:"invoice_payload"`
	TelegramPaymentID string `json:"telegram_payment_charge_id"`
	ProviderPaymentID string `json:"provider_payment_charge_id"`
}

type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard        bool               `json:"resize_keyboard"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard"`
	InputFieldPlaceholder string             `json:"input_field_placeholder"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type InlineKeyboardMarkup struct {
	Buttons [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
	Url          string `json:"url"`
}

type MessageParams struct {
	ChatID         int
	Text           string
	KeyboardReply  *ReplyKeyboardMarkup
	KeyboardInline *InlineKeyboardMarkup
}

type InvoiceParams struct {
	ChatID        int
	Title         string
	Description   string
	Payload       string
	ProviderToken string
	Currency      string
	Prices        *[]LabeledPrice //JSON-serialized list of components
	// PhotoURL string
}

type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

type PreCheckoutParams struct {
	QueryID      string
	OK           bool
	ErrorMessage string
}
