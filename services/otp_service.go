package services

import (
	"context"
	"errors"
	"sample-web/configs"
	"sample-web/utils"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	_, span := utils.Tracer().Start(ctx, "twilioService.SendOTP")
	defer span.End()

	span.AddEvent("SendOTPRequestReceived", trace.WithAttributes(
		attribute.String("phone_number", phoneNumber),
	))

	params := &verify.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel(twilioChanelSMS)

	span.AddEvent("SendingOTP")

	resp, err := t.client.VerifyV2.CreateVerification(t.config.ServiceSID, params)

	if err != nil {
		span.RecordError(err)
		return "", err
	}

	verification_sid := *resp.Sid
	verification_status := *resp.Status

	if !*resp.Valid {
		span.AddEvent("InvalidPhoneNumber", trace.WithAttributes(
			attribute.String("phone_number", phoneNumber),
			attribute.String("status", verification_status),
			attribute.String("verification_sid", verification_sid),
		))
		span.RecordError(errors.New("invalid phone number"))
		return "", errors.New("invalid phone number")
	}

	span.AddEvent("OTPSent", trace.WithAttributes(
		attribute.String("verification_sid", verification_sid),
		attribute.String("status", verification_status),
	))
	return verification_sid, nil
}

func (t *twilioService) VerifyOTP(ctx context.Context, phoneNumber string, code string) (bool, error) {

	_, span := utils.Tracer().Start(ctx, "twilioService.VerifyOTP")
	defer span.End()
	span.AddEvent("VerifyOTPRequestReceived", trace.WithAttributes(
		attribute.String("phone_number", phoneNumber),
	))

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	span.AddEvent("VerifyingOTP")

	resp, err := t.client.VerifyV2.CreateVerificationCheck(t.config.ServiceSID, params)

	if err != nil {
		span.AddEvent("ErrorVerifyingOTP")
		span.RecordError(err)
		return false, err
	}

	verification_sid := *resp.Sid
	verification_status := *resp.Status

	if *resp.Status != twilioStatusApproved {
		span.AddEvent("VerificationFailed", trace.WithAttributes(
			attribute.String("verification_sid", verification_sid),
			attribute.String("status", verification_status),
		))
		return false, errors.New("verification failed")
	}

	span.AddEvent("VerificationSuccessful", trace.WithAttributes(
		attribute.String("verification_sid", verification_sid),
		attribute.String("status", verification_status),
	))

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
