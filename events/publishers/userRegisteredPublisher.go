package publishers

import (
	"github.com/yonraz/gochat_auth/constants"
)

func (p *Publisher) UserRegistered(username string) error {
	body := map[string]string{"username": username}
	return p.Publish(constants.UserEventsExchange, constants.UserRegisteredKey, body)
}