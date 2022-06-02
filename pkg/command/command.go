package command

import "reflect"

// The command is only a request, and thus may be refused.
// The event is a fact from the past.

// Types -
type Types Command

// Command -
type Command interface {
	// AggregateID -
	AggregateID() string

	// CommandType -
	Type() (reflect.Type, string)

	// Message -
	Message() interface{}
}

// descriptor -
type descriptor struct {
	id      string
	message interface{}
}

// New -
func New(aggregateID string, message interface{}) Command {
	return &descriptor{
		id:      aggregateID,
		message: message,
	}
}

// Type -
func (d *descriptor) Type() (reflect.Type, string) {
	t := reflect.TypeOf(d.message)

	return t, t.Name()
}

// AggregateID -
func (d *descriptor) AggregateID() string {
	return d.id
}

// Command -
func (d *descriptor) Message() interface{} {
	return d.message
}
