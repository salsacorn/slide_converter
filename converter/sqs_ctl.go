package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

const (
	AWS_REGION = "ap-northeast-1"
	QUEUE_URL  = "https://sqs.ap-northeast-1.amazonaws.com/985295721336/slidehub-convert02"
)

type Queue struct {
	Client sqsiface.SQSAPI
	URL    string
}

// Message is a concrete representation of the SQS message
type Message struct {
	Id         int    `json:"id"`
	Object_key string `json:"object_key"`
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func createQueues() {
	sess := session.Must(session.NewSession())
	svc := sqs.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)
	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String("slidehub-convert01"),
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String("60"),
			"MessageRetentionPeriod": aws.String("86400"),
		},
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *result.QueueUrl)
}

func reciveMessage() ([]*sqs.Message, error) {
	sess := session.Must(session.NewSession())
	svc := sqs.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

	resp, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(QUEUE_URL),
		MaxNumberOfMessages: aws.Int64(5),
		VisibilityTimeout:   aws.Int64(36000), // 10 hours
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get messages, %v", err)
	}

	if len(resp.Messages) == 0 {
		return nil, fmt.Errorf("failed to get resp.messages, %v", err)
	}

	/*
		msgs := make([]Message, len(resp.Messages))
		for i, msg := range resp.Messages {
			parsedMsg := Message{}
			if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedMsg); err != nil {
				return nil, fmt.Errorf("failed to unmarshal message, %v", err)
			}

			msgs[i] = parsedMsg
		}
	*/
	return resp.Messages, nil
}

func PerseMessage(msg *sqs.Message) (Message, error) {
	parsedMsg := Message{}
	if err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), &parsedMsg); err != nil {
		// ToDo: nil返却について調べる
		return parsedMsg, fmt.Errorf("failed to unmarshal message, %v", err)
	}
	return parsedMsg, nil
}

func deleteMessage(msg *sqs.Message) error {
	sess := session.Must(session.NewSession())
	svc := sqs.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(QUEUE_URL),
		ReceiptHandle: aws.String(*msg.ReceiptHandle),
	})

	if err != nil {
		fmt.Println("Delete Error", err)
		return err
	}

	fmt.Println("Message Deleted")
	return nil
}

func sendMessage(msg string) {
	sess := session.Must(session.NewSession())
	svc := sqs.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

	// URL to our queue
	result, err := svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:     aws.String(QUEUE_URL),
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(msg),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *result.MessageId)
}
