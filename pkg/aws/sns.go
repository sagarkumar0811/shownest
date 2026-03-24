package aws

import (
	"context"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snsTypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
)

type SNSClient struct {
	client   *sns.Client
	mockMode bool
}

func NewSNSClient(cfg awssdk.Config, mockMode bool) *SNSClient {
	var client *sns.Client
	if !mockMode {
		client = sns.NewFromConfig(cfg)
	}
	return &SNSClient{client: client, mockMode: mockMode}
}

func (s *SNSClient) SendSMS(ctx context.Context, phone, message string) error {
	if s.mockMode {
		logger.WithContext(ctx).Info("[SNS MOCK] SMS not sent",
			zap.String("phone", phone),
			zap.String("message", message),
		)
		return nil
	}

	_, err := s.client.Publish(ctx, &sns.PublishInput{
		PhoneNumber: awssdk.String(phone),
		Message:     awssdk.String(message),
		MessageAttributes: map[string]snsTypes.MessageAttributeValue{
			"AWS.SNS.SMS.SMSType": {
				DataType:    awssdk.String("String"),
				StringValue: awssdk.String("Transactional"),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("sns: send sms to %s: %w", phone, err)
	}
	return nil
}
