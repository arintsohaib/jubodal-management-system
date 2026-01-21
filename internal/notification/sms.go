package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SMSService struct {
	APIKey   string
	URL      string
	SenderID string
}

func NewSMSService() *SMSService {
	return &SMSService{
		APIKey:   os.Getenv("SMS_API_KEY"),
		URL:      os.Getenv("SMS_API_URL"),
		SenderID: os.Getenv("SMS_SENDER_ID"),
	}
}

func (s *SMSService) Send(to string, message string) error {
	if s.APIKey == "" || s.URL == "" {
		fmt.Printf("[MOCK SMS] To: %s, Message: %s\n", to, message)
		return nil
	}

	payload := map[string]string{
		"api_key":   s.APIKey,
		"sender_id": s.SenderID,
		"number":    to,
		"message":   message,
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(s.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS gateway returned status: %d", resp.StatusCode)
	}

	return nil
}
