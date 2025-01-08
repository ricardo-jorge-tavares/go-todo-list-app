package helpers

import "log"

func CheckError(err error) {
	// Handle errors
	if err != nil {
		log.Fatal(err)
	}
}
