package main

type PubStream <-chan []byte

var subs = make([]chan<- []byte, 0, 10)
var pub PubStream

func messageSub() PubStream {
	newChan := make(chan []byte, 2)
	subs = append(subs, newChan)
	return newChan
}

func messagePub(in PubStream) {
	pub = in
}

func handlePubSub() {
	for pub == nil {
		// Hang around
	}

	for {
		msg, ok := <-pub

		// Shut down on message close
		if !ok {
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
