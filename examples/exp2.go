package main

import (
	"flag"
	"fmt"
	"koding/newkite/kite"
	"koding/newkite/protocol"
	"math/rand"
	"time"
)

var port = flag.String("port", "", "port to bind itself")

func main() {
	flag.Parse()

	options := &kite.Options{
		Kitename:    "application",
		Version:     "0.0.1",
		Port:        *port,
		Region:      "localhost",
		Environment: "development",
		Username:    "devrim",
	}

	k := kite.New(options)
	k.Start()

	query := protocol.KontrolQuery{
		Username:    "devrim",
		Environment: "development",
		Name:        "mathworker",
	}

	// To demonstrate we can receive notifications matcing to our query.
	onEvent := func(e *protocol.KiteEvent) {
		fmt.Printf("--- kite event: %#v\n", e)
	}

	go func() {
		err := k.Kontrol.WatchKites(query, onEvent)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// .. or just get the current kites and dial for one
	kites, err := k.Kontrol.GetKites(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	mathWorker := kites[0]
	err = mathWorker.Dial()
	if err != nil {
		fmt.Println("Cannot connect to remote mathworker")
		return
	}

	squareOf := func(i int) {
		response, err := mathWorker.Call("square", i)
		if err != nil {
			fmt.Println(err)
			return
		}

		result := response.MustFloat64()
		fmt.Printf("input: %d  rpc result: %f\n", i, result)
	}

	for {
		squareOf(rand.Intn(10))
		time.Sleep(time.Second)
	}
}
