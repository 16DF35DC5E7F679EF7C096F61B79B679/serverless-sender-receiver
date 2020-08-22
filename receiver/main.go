package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Receive(sess *session.Session, queueURL *string) (*sqs.ReceiveMessageOutput, error){

	svc := sqs.New(sess)


	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(60), // 60 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})
	if err != nil {
		fmt.Printf("Error %e \n", err)
		return nil, err
	}
	fmt.Println("Message is: ")
	fmt.Println(result.GoString())
	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
		return result, nil
	} else {
		fmt.Printf("Success: %+v\n", result.Messages)
		return nil, fmt.Errorf("Empty string detected in message ")
	}
}

func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	queue := "common-queue"
	result, err := GetQueueURL(sess, &queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return Response{StatusCode: 500}, err
	}

	queueURL := result.QueueUrl
	message, err := Receive(sess, queueURL)
	if err != nil {
		fmt.Println("Got an error receiving the message:")
		fmt.Println(err)
		return Response{StatusCode: 500}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": message.GoString(),
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "receiver-handler",
		},
	}

	return resp, nil
}

func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("url is : "+ result.GoString())
	return result, nil
}


func main() {
	lambda.Start(Handler)
}


