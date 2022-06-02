package command

import (
	"context"
	"errors"
	"reflect"
)

var (
	// Ierror - error interface
	Ierror = reflect.TypeOf((*error)(nil)).Elem()
)

const (
	aggregateID = "AggregateID"
)

// dispatch -
type dispatch struct {
	// map key - command
	m map[reflect.Type]content
}

// content -
type content struct {
	imple  reflect.Value
	method string
}

// Dispatch -
type Dispatch interface {
	// Handle -
	Handle(ctx context.Context, cmd Command) (event interface{}, err error)
}

// NewDispatch -
// set args to map
// using command as key
func NewDispatch(args ...interface{}) Dispatch {
	d := &dispatch{
		m: make(map[reflect.Type]content),
	}

	for x := 0; x < len(args); x++ {
		arg := args[x]
		t := reflect.TypeOf(arg)
		v := reflect.ValueOf(arg)

		for i := 0; i < t.NumMethod(); i++ {
			method := t.Method(i)

			fnType := method.Func.Type()

			// get args number
			numArgs := fnType.NumIn()

			argType := fnType.In(numArgs - 1)

			// set map
			d.m[argType] = content{
				imple:  v,
				method: method.Name,
			}
		}
	}

	return d
}

// Handle -
func (d *dispatch) Handle(ctx context.Context, cmd Command) (event interface{}, err error) {
	t, _ := cmd.Type()

	v, ok := d.m[t]
	if !ok {
		return nil, errors.New("no matching command")
	}

	setAggregateID(v.imple, aggregateID, cmd.AggregateID())

	// set value
	oV := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(cmd.Message()),
	}

	// call
	oValue := v.imple.MethodByName(v.method).Call(oV)

	// process result
	// find error content and set
	for _, oV := range oValue {
		i := oV.Interface()

		switch oV.Type().Implements(Ierror) {
		case true:
			err, ok = i.(error)
			if !ok {
				err = nil
			}

		default:
			event = i
		}
	}

	return
}

func setAggregateID(oV reflect.Value, fieldByName string, aggregateID string) {
	var ov reflect.Value

	switch oV.Kind() {
	case reflect.Ptr:
		ov = oV.Elem()
	}

	ov = ov.FieldByName(fieldByName)

	if ov.IsValid() && ov.CanSet() {
		ov.SetString(aggregateID)
	}
}
