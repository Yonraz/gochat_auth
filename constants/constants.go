package constants

type Queues string
type RoutingKey string
type Exchange string

const (
	UserEventsExchange Exchange = "UserEventsExchange"
)

const (
	UserRegisteredKey RoutingKey = "user.registered"
	UserLoggedInKey   RoutingKey = "user.logged.in"
	UserSignedoutKey  RoutingKey = "user.signed.out"
)

const (
	UserRegistrationQueue Queues = "UserRegistrationQueue"
	UserLoginQueue        Queues = "UserLoginQueue"
	UserSignoutQueue      Queues = "UserSignoutQueue"
)