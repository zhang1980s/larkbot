package dao

import (
	"case-refresh/config"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/sirupsen/logrus"
)

func getSecretValue(ctx context.Context, secretArn string) (string, error) {

	cfg, err := awsconfig.LoadDefaultConfig(ctx)

	if err != nil {
		logrus.Errorf("failed to load AWS config: %s", err)
		return "", err
	}

	svc := secretsmanager.NewFromConfig(cfg)

	output, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretArn),
	})

	if err != nil {
		logrus.Errorf("failed to get secret value: %s", err)
		return "", err
	}

	secretValue := aws.ToString(output.SecretString)

	return secretValue, nil

}

func GetAppID() (string, error) {
	return getSecretValue(context.Background(), config.Conf.AppIDARN)
}
func GetAPPSecret() (string, error) {
	return getSecretValue(context.Background(), config.Conf.AppSecretARN)
}
