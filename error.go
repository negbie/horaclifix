package main

import "log"

func checkErr(err error) {
	if err != nil {
		log.Printf("%s warning: %s", *name, err)
	}
}

func checkCritErr(err error) {
	if err != nil {
		log.Printf("%s critical: %s", *name, err)
	}
}
