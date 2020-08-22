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

func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer
	region := "us-east-1"
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{Region: &region},
	}))
	queue := "common-queue"
	result, err := GetQueueURL(sess, &queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return Response{StatusCode: 500}, err
	}
	//queueUrl := result.QueueUrl
	queueUrl := "https://sqs.us-east-1.amazonaws.com/221863034380/common-queue"
	fmt.Println("queueUrl is: "+ queueUrl)
	fmt.Println("Credentials are: ")
	fmt.Println(sess.Config.Credentials.Get())
	fmt.Println("queue url is: "+ *result.QueueUrl)
	err = SendMsg(sess, result.QueueUrl)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return Response{StatusCode: 500}, err
	}

	fmt.Println("Sent message to queue ")
	body, err := json.Marshal(map[string]interface{}{
		"message": "Event successfully pushed to queue",
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
			"X-MyCompany-Func-Reply": "world-handler",
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

	return result, nil
}

func SendMsg(sess *session.Session, queueURL *string) error {
	// Create an SQS service client
	// snippet-start:[sqs.go.send_message.call]
	svc := sqs.New(sess)

	_, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("The Whistler"),
			},
			"Author": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("John Grisham"),
			},
			"WeeksOn": &sqs.MessageAttributeValue{
				DataType:    aws.String("Number"),
				StringValue: aws.String("6"),
			},
		},
		MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
		QueueUrl:    queueURL,
	})
	// snippet-end:[sqs.go.send_message.call]
	if err != nil {
		return err
	}

	return nil
}


func main() {
	lambda.Start(Handler)
}


