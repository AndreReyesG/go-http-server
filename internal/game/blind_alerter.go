package game

import (
	"fmt"
	"os"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}

// NOTE: Remember that any type can implement an interface, not just structs.
// If you are making a library that exposes an interface with one function
// defined it is a common idiom to also expose a MyInterfaceFunc type.
type BlindAlerterFunc func(duration time.Duration, amount int)

// This type will be a func which will also implement your interface.
// That way users of your interface have the option to implement your interface
// with just a function; rather than having to create an empty struct type.
func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

// We then create the function StdOutAlerter which has the same signature as
// the function and just use time.AfterFunc to schedule it to print to os.Stdout.
func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
