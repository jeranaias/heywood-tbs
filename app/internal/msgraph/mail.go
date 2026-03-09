package msgraph

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"heywood-tbs/internal/models"
)

// MailService reads Outlook mail via Microsoft Graph.
type MailService struct {
	client *Client
}

// NewMailService creates a mail service backed by the given Graph client.
func NewMailService(client *Client) *MailService {
	return &MailService{client: client}
}

// graphMessage is the Microsoft Graph message response shape.
type graphMessage struct {
	ID          string `json:"id"`
	Subject     string `json:"subject"`
	BodyPreview string `json:"bodyPreview"`
	ReceivedAt  string `json:"receivedDateTime"`
	IsRead      bool   `json:"isRead"`
	HasAttach   bool   `json:"hasAttachments"`
	From        struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"from"`
}

// GetMailSummary retrieves recent emails for a user.
// Returns up to 10 recent messages with subject, sender, preview.
func (s *MailService) GetMailSummary(userID string, unreadOnly bool) ([]models.MailSummary, error) {
	if !s.client.IsConfigured() {
		return nil, nil
	}

	path := fmt.Sprintf("/users/%s/messages", userID)
	params := map[string]string{
		"$select":  "id,subject,bodyPreview,receivedDateTime,isRead,hasAttachments,from",
		"$orderby": "receivedDateTime desc",
		"$top":     "10",
	}
	if unreadOnly {
		params["$filter"] = "isRead eq false"
	}

	body, err := s.client.Get(path, params)
	if err != nil {
		slog.Error("graph mail query failed", "user", userID, "error", err)
		return nil, err
	}

	var resp struct {
		Value []graphMessage `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse mail response: %w", err)
	}

	mails := make([]models.MailSummary, 0, len(resp.Value))
	for _, gm := range resp.Value {
		from := gm.From.EmailAddress.Name
		if from == "" {
			from = gm.From.EmailAddress.Address
		}
		mails = append(mails, models.MailSummary{
			ID:        gm.ID,
			Subject:   gm.Subject,
			From:      from,
			Preview:   gm.BodyPreview,
			Received:  gm.ReceivedAt,
			IsRead:    gm.IsRead,
			HasAttach: gm.HasAttach,
		})
	}

	return mails, nil
}

// UnreadCount returns the number of unread messages for a user.
func (s *MailService) UnreadCount(userID string) (int, error) {
	if !s.client.IsConfigured() {
		return 0, nil
	}

	path := fmt.Sprintf("/users/%s/mailFolders/inbox", userID)
	params := map[string]string{
		"$select": "unreadItemCount",
	}

	body, err := s.client.Get(path, params)
	if err != nil {
		return 0, err
	}

	var resp struct {
		UnreadItemCount int `json:"unreadItemCount"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}

	return resp.UnreadItemCount, nil
}

// SendMail sends an email via Microsoft Graph on behalf of a user.
// Uses POST /users/{id}/sendMail
func (s *MailService) SendMail(userID string, to []string, subject, body string) error {
	if !s.client.IsConfigured() {
		return fmt.Errorf("graph client not configured")
	}

	recipients := make([]map[string]interface{}, len(to))
	for i, addr := range to {
		recipients[i] = map[string]interface{}{
			"emailAddress": map[string]string{"address": addr},
		}
	}

	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"subject":      subject,
			"body":         map[string]string{"contentType": "Text", "content": body},
			"toRecipients": recipients,
		},
	}

	path := fmt.Sprintf("/users/%s/sendMail", userID)
	_, err := s.client.Post(path, payload)
	return err
}

// ReplyToMail replies to a message via Microsoft Graph.
// Uses POST /users/{id}/messages/{messageId}/reply
func (s *MailService) ReplyToMail(userID, messageID, comment string) error {
	if !s.client.IsConfigured() {
		return fmt.Errorf("graph client not configured")
	}

	payload := map[string]interface{}{
		"comment": comment,
	}

	path := fmt.Sprintf("/users/%s/messages/%s/reply", userID, messageID)
	_, err := s.client.Post(path, payload)
	return err
}
