package main

import (
	"fmt"
	"time"

	"github.com/EZCampusDevs/firepit/data"
)

/*

This is a test to see if we can control thread states using a channel

In this example, we have a goroutine updating the id and printing . when it does

The state is then changed a few times

In order for this to be threadsafe the run function is the only thing which should modify the stuct

Eveything else can read them safely when it is paused

If members of the struct need to be modified, they should be done in the run function using a channel

*/

type Room struct {
	ID string

	setId chan string
	state chan byte
}

func NewRoom(name string) *Room {
	return &Room{
		setId: make(chan string),
		state: make(chan byte),
	}
}

func (r *Room) run() {
	defer fmt.Printf("Thread is dead\n")

	for {
		select {
		case id := <-r.setId:
			r.ID = id
			fmt.Print(".")
		case s := <-r.state:

			fmt.Printf("Setting state to %s\n", data.ChannelStateToString(s))

			switch s {
			case data.CHAN__DEAD:
				return
			case data.CHAN__RUNNING:
				continue

			case data.CHAN__PAUSED:
			default:
				break
			}

		lock:
			for {

				s = <-r.state
				fmt.Printf("Setting state to %s\n", data.ChannelStateToString(s))
				switch s {
				case data.CHAN__PAUSED:
					break

				case data.CHAN__DEAD:
					return

				case data.CHAN__RUNNING:
					break lock
				}
			}
		}
	}
}

func main() {

	r := NewRoom("this is the room name")

	go r.run()
	go func() {

		for i := 0; ; i++ {
			r.setId <- fmt.Sprintf("%d", i)

			time.Sleep(200 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)

	// test that this keeps it running
	r.state <- data.CHAN__RUNNING

	time.Sleep(1 * time.Second)

	// pause it
	r.state <- data.CHAN__PAUSED

	time.Sleep(1 * time.Second)

	// make sure this keeps it paused
	r.state <- data.CHAN__PAUSED

	time.Sleep(1 * time.Second)

	// read the id, this should be thread safe, since we are reading
	fmt.Printf("\nThe name is %s\n", r.ID)

	time.Sleep(2 * time.Second)

	// resume the thread
	r.state <- data.CHAN__RUNNING

	time.Sleep(2 * time.Second)

	// kill the thread
	r.state <- data.CHAN__DEAD

	time.Sleep(5 * time.Second)

}
