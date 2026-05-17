package main

import "log"

func order() {
	log.Println("Start order...")
	sendEmail()
	sendSMS()
	log.Println("Success order.")
}

func sendEmail() {
	log.Println("Start sendEmail...")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Notes: hotfix...")
		}
	}()
	panic("sendEmail error!")
	log.Println("Success sendEmail.")
}

func sendSMS() {
	log.Println("Start sendSMS...")
	log.Println("Success sendSMS.")
}

func main() {
	log.Println("Process API...")
	order()
	log.Println("End API.")
}
