package main

import "log"

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func noerr[T any](r T, err error) T {
	must(err)
	return r
}
