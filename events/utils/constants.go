package utils

import "github.com/yonraz/gochat_auth/constants"

var Brokers = []string{"kafka-srv:9092"}
var Topics = []constants.Topic{constants.UserLoggedIn, constants.UserRegistered, constants.UserSignedout}