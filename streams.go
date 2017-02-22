package main

var subs = make([]chan<- string, 0, 10)
var pub <-chan string

func messageSub() <-chan string {
	newChan := make(chan string, 2)
	subs = append(subs, newChan)
	return newChan
}

func messagePub(in <-chan string) {
	pub = in
}

func handlePubSub() {
	for pub == nil {
		// Hang around
	}

	for {
		msg := <-pub

		// Shut down on message close
		if msg == "" {
			log.Info("Shutting down pub/sub due to channel closure")
			for _, c := range subs {
				close(c)
			}
			pub = nil
			break
		}

		for i, out := range subs {
			select {
			case out <- msg:
			default:
				// Can't write to channel? Destroy it
				close(out)
				subs[i] = subs[len(subs)-1]
				subs = subs[:len(subs)-1]
			}
		}

	}
}
