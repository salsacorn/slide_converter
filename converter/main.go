package main

import (
	"fmt"
	"os"
)

func main() {
	//sqs_listqueus()
	//createQueues()
	/*
		send_msgs := []string{`{"id":4, "object_key":"1d5837d3cbda46b3fe902d62cb79624f"}`,
			`{"id":5, "object_key":"e473585337ecde0a9071a7d55b8c6e72"}`,
			`{"id":6, "object_key":"fbbac927f1037c2f7382c190fd68be7b"}`}
		for _, msg := range send_msgs {
			sendMessage(msg)
		}
	*/
	fmt.Println("Reciveing Message from SQS")
	msgs, err := reciveMessage()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	for _, msg := range msgs {
		parsedMsg, _ := PerseMessage(msg)
		fmt.Printf("[Message Detail] id:%s, object_key:%s\n", parsedMsg.Id, parsedMsg.Object_key)
		err = DownloadFile(parsedMsg.Object_key)
		if err != nil {
			fmt.Println("Error", err)
		}

		if err := deleteMessage(msg); err != nil {
			fmt.Println(err)
		}
	}
}
