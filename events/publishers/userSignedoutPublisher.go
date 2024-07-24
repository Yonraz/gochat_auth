package publishers

import "github.com/yonraz/gochat_auth/constants"

func (p *Publisher) UserSignedout(username string) error {
	body := map[string]string{"username": username}
	return p.Publish(constants.UserEventsExchange, constants.UserSignedoutKey, body)
}