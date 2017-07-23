package main

import "log"

func checkErr(err error) {
	if err != nil {
		log.Println("Warning:", err)
	}
}

func checkCritErr(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v\n", err)
	}
}
