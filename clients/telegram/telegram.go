package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest("getUpdates", q)
	if err != nil {
		return nil, err
	}
	var res UpdateResponce

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (c *Client) SendMessage(params MessageParams) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(params.ChatID))
	q.Add("text", params.Text)
	if params.KeyboardReply != nil {
		jsonKeyboard, err := json.Marshal(params.KeyboardReply)
		if err != nil {
			return err
		}
		q.Add("reply_markup", string(jsonKeyboard))
	}
	if params.KeyboardInline != nil {
		jsonKeyboard, err := json.Marshal(params.KeyboardInline)
		if err != nil {
			return err
		}
		q.Add("reply_markup", string(jsonKeyboard))
	}
	_, err := c.doRequest("sendMessage", q)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func (c *Client) SendInvoice(params InvoiceParams) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(params.ChatID))
	q.Add("title", params.Title)
	q.Add("description", params.Description)
	q.Add("payload", params.Payload)
	q.Add("provider_token", params.ProviderToken)
	q.Add("currency", params.Currency)
	data, err := json.Marshal(params.Prices)
	if err != nil {
		return err
	}
	q.Add("prices", string(data))
	// q.Add("photo_url", params.PhotoURL)

	_, err = c.doRequest("sendInvoice", q)
	if err != nil {
		return fmt.Errorf("can't send invoice: %w", err)
	}
	return nil
}

func (c *Client) AnswerPreCheckoutQuery(params PreCheckoutParams) error {
	q := url.Values{}
	q.Add("pre_checkout_query_id", params.QueryID)
	q.Add("ok", strconv.FormatBool(params.OK))
	q.Add("error_message", params.ErrorMessage)

	_, err := c.doRequest("answerPreCheckoutQuery", q)
	if err != nil {
		return fmt.Errorf("can't answer precheckout query: %w", err)
	}
	return nil
}

func (c *Client) CreateInviteLink(linkName string) error {
	q := url.Values{}
	q.Add("chat_id", os.Getenv("CHANNEL_ID"))
	q.Add("name", linkName)
	q.Add("creates_join_request", strconv.FormatBool(true))

	_, err := c.doRequest("createChatInviteLink", q)
	if err != nil {
		return fmt.Errorf("can't create invite link: %w", err)
	}
	return nil
}

func (c *Client) ApproveChatJoinRequest(userID int64) error {
	q := url.Values{}
	q.Add("chat_id", os.Getenv("CHANNEL_ID"))
	q.Add("user_id", strconv.FormatInt(userID, 10))

	_, err := c.doRequest("approveChatJoinRequest", q)
	if err != nil {
		return fmt.Errorf("can't approve chat join request: %w", err)
	}
	return nil
}

func (c *Client) BanUser(userID int64) error {
	q := url.Values{}
	q.Add("chat_id", os.Getenv("CHANNEL_ID"))
	q.Add("user_id", strconv.FormatInt(userID, 10))

	_, err := c.doRequest("banChatMember", q)
	if err != nil {
		return fmt.Errorf("can't ban user: %w", err)
	}
	return nil
}

func (c *Client) UnbanUser(userID int64) error {
	q := url.Values{}
	q.Add("chat_id", os.Getenv("CHANNEL_ID"))
	q.Add("user_id", strconv.FormatInt(userID, 10))

	_, err := c.doRequest("unbanChatMember", q)
	if err != nil {
		return fmt.Errorf("can't unban user: %w", err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	return body, nil
}
