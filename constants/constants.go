package constants

type Queues string
type Topic string
type Exchange string

const (
	UserEventsExchange Exchange = "UserEventsExchange"
)

const (
	UserRegistered Topic = "user.registered"
	UserLoggedIn   Topic = "user.logged.in"
	UserSignedout  Topic = "user.signed.out"
)

const (
	UserRegistrationQueue Queues = "UserRegistrationQueue"
	UserLoginQueue        Queues = "UserLoginQueue"
	UserSignoutQueue      Queues = "UserSignoutQueue"
)