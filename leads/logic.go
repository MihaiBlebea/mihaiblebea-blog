package leads

import (
	"fmt"
	"strings"

	"github.com/tw1nk/gochimp3"
)

type service struct {
	client *gochimp3.API
	listID string
}

// New returns a new leads service
func New(apiKey, listID string) Service {
	client := gochimp3.New(apiKey)
	return &service{client, listID}
}

func (s *service) Store(email string) error {
	list, err := s.client.GetList(s.listID, nil)
	if err != nil {
		return err
	}

	fmt.Println(list)

	req := &gochimp3.MemberRequest{
		EmailAddress: email,
		Status:       "subscribed",
	}

	if _, err := list.CreateMember(req); err != nil {
		return err
	}

	return nil
}

func (s *service) formatName(name string) string {
	return strings.Title(name)
}

func (s *service) formatEmail(email string) string {
	return strings.ToLower(email)
}
