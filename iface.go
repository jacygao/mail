package mail

// Mailer provides an interface to allow mocking and unit testing.
type Mailer interface {
	Send(to string, msg Message) error
}
