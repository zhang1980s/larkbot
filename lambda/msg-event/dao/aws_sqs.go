package dao

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

func SendContentToSQS(ctx context.Context, content string, msgid string) error {

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		logrus.Errorf("failed to load AWS config: %s", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	messageBody, err := json.Marshal(map[string]string{
		"content": content,
	})

	if err != nil {
		logrus.Errorf("Error marshalling Q messages %s", err)
		return err
	}

	input := &sqs.SendMessageInput{
		QueueUrl:       aws.String(os.Getenv("SQS_URL")),
		MessageGroupId: aws.String(msgid),
		MessageBody:    aws.String(string(messageBody)),
	}

	_, err = sqsClient.SendMessage(context.TODO(), input)

	if err != nil {
		logrus.Errorf("Error sending Q content message to SQS %s", err)
		return err
	}

	return nil
}
