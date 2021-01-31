package leads

import (
	"errors"
	"regexp"
	"sync"

	"github.com/tw1nk/gochimp3"
)

// Custom errors
var (
	ErrInvalidEmail = errors.New("Email is invalid")
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type service struct {
	client *gochimp3.API
	listID string
	tokens []string
	lock   sync.Mutex
}

// New returns a new leads service
func New(apiKey, listID string) Service {
	client := gochimp3.New(apiKey)

	return &service{client: client, listID: listID}
}

func (s *service) Store(email string) error {
	if s.isEmailValid(email) == false {
		return ErrInvalidEmail
	}

	list, err := s.client.GetList(s.listID, nil)
	if err != nil {
		return err
	}

	req := &gochimp3.MemberRequest{
		EmailAddress: email,
		Status:       "subscribed",
	}

	if _, err := list.CreateMember(req); err != nil {
		return err
	}

	return nil
}

func (s *service) isEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	return emailRegex.MatchString(email)
}
