package services

import (
	"context"
	"errors"
	"fmt"
	"sample-web/configs"
	"sample-web/utils"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

const (
	twilioChanelSMS = "sms"
)

const (
	twilioStatusApproved = "approved"
)

type OTPService interface {
	SendOTP(ctx context.Context, phoneNumber string) (string, error)
	VerifyOTP(ctx context.Context, phoneNumber string, code string) (bool, error)
}

type twilioService struct {
	client *twilio.RestClient
	config configs.TwilioConfig
}

func (t *twilioService) SendOTP(ctx context.Context, phoneNumber string) (string, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "twilioService.SendOTP")
	defer span.End()

	log.Info(spanCtx, fmt.Sprintf("SendOTP Request receieved for phone number %s", phoneNumber))

	params := &verify.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel(twilioChanelSMS)

	log.Info(spanCtx, fmt.Sprintf("SendingOTP to %s", phoneNumber))

	resp, err := t.client.VerifyV2.CreateVerification(t.config.ServiceSID, params)

	if err != nil {
		log.Error(spanCtx, err.Error())
		return "", err
	}

	verificationSid := *resp.Sid
	verificationStatus := *resp.Status

	log.Info(spanCtx, fmt.Sprintf("OTP Sent to %s with verification_sid as %s and status %s", phoneNumber, verificationSid, verificationStatus))

	return verificationSid, nil
}

func (t *twilioService) VerifyOTP(ctx context.Context, phoneNumber string, code string) (bool, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "twilioService.VerifyOTP")
	defer span.End()

	log.Info(spanCtx, fmt.Sprintf("VerifyOTP request received for phone number %s", phoneNumber))

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	log.Info(spanCtx, "Verifying OTP")

	resp, err := t.client.VerifyV2.CreateVerificationCheck(t.config.ServiceSID, params)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error Verifying OTP with error %s", err.Error()))
		return false, err
	}

	verificationSid := *resp.Sid
	verificationStatus := *resp.Status

	if *resp.Status != twilioStatusApproved {
		log.Error(spanCtx, fmt.Sprintf("OTP Verification failed for %s with verification_sid as %s and status %s", phoneNumber, verificationSid, verificationStatus))
		return false, errors.New("verification failed")
	}

	log.Info(spanCtx, fmt.Sprintf("OTP Verification Succeeded for %s with verification_sid as %s and status %s", phoneNumber, verificationSid, verificationStatus))

	return true, nil
}

func NewTwilioClient(config configs.TwilioConfig) OTPService {
	client := twilio.NewRestClientWithParams(
		twilio.ClientParams{
			Username: config.AccountSID,
			Password: config.AuthToken,
		},
	)
	return &twilioService{
		client: client,
		config: config,
	}
}
