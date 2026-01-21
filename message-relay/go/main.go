package main

import (
	"flag"
	"log"

	"github.com/Scotiacon-Tech/libs/message-relay/go/lib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	keyArg := flag.String("key", "", "Auth Key")
	serviceArg := flag.String("service", "", "Service Name")
	fromArg := flag.String("from", "", "Message From")
	toArg := flag.String("to", "", "Message Recipient")
	subjectArg := flag.String("subject", "", "Message Subject")
	bodyArg := flag.String("body", "", "Message Body")

	flag.Parse()

	client := lib.NewClient()

	req := client.NewSendRequest()
	req.From = *fromArg
	req.To = *toArg
	req.Subject = *subjectArg
	req.Body = *bodyArg

	res, err := client.SendMessage(*keyArg, *serviceArg, req)

	if res == nil {
		log.Print(err)
	} else {
		log.Printf("Key: %s", res.NewKey)
		log.Printf("MessageID: %s", res.MessageID)
	}
}
