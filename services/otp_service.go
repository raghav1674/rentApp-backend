package services

import (
	"context"
	"fmt"
	"sample-web/clients"
	"sample-web/configs"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

const (
	twilioChannelSMS          = "sms"
	twilioTotalRetries        = 3
	twilioRetryKeySuffix      = "twilio_invalid_attempts"
	twilioExpiryTimeInMinutes = 10
)

const (
	twilioStatusApproved = "approved"
)

type OTPService interface {
	SendOTP(ctx context.Context, phoneNumber string) (string, error)
	VerifyOTP(ctx context.Context, phoneNumber string, code string) (bool, error)
}

type twilioOtpService struct {
	client      *twilio.RestClient
	config      configs.TwilioConfig
	redisClient *clients.RedisClient
	expiry      time.Duration
}

func NewTwilioOTPService(cfg configs.TwilioConfig, redisClient *clients.RedisClient) OTPService {
	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSID,
		Password: cfg.AuthToken,
	})
	return &twilioOtpService{
		client:      twilioClient,
		config:      cfg,
		redisClient: redisClient,
		expiry:      twilioExpiryTimeInMinutes * time.Minute,
	}
}

func (s *twilioOtpService) SendOTP(ctx context.Context, phoneNumber string) (string, error) {
	retryKey := s.buildRetryKey(phoneNumber)

	if blocked, ttl := s.isBlocked(ctx, retryKey); blocked {
		return "", &PhoneNumberBlockedError{PhoneNumber: phoneNumber, RetryAfter: ttl}
	}

	params := &verify.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel(twilioChannelSMS)

	resp, err := s.client.VerifyV2.CreateVerification(s.config.ServiceSID, params)
	if err != nil {
		s.incrementRetry(ctx, retryKey)
		return "", err
	}

	return *resp.Sid, nil
}

func (s *twilioOtpService) VerifyOTP(ctx context.Context, phoneNumber, code string) (bool, error) {
	retryKey := s.buildRetryKey(phoneNumber)

	if blocked, ttl := s.isBlocked(ctx, retryKey); blocked {
		return false, &PhoneNumberBlockedError{PhoneNumber: phoneNumber, RetryAfter: ttl}
	}

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := s.client.VerifyV2.CreateVerificationCheck(s.config.ServiceSID, params)
	if err != nil {
		s.incrementRetry(ctx, retryKey)
		return false, err
	}

	if resp.Status == nil || *resp.Status != twilioStatusApproved {
		count := s.incrementRetry(ctx, retryKey)
		if count >= totalRetries {
			return false, &PhoneNumberBlockedError{PhoneNumber: phoneNumber, RetryAfter: s.expiry}
		}
		remaining := totalRetries - count
		return false, fmt.Errorf("invalid OTP, %d attempt(s) remaining", remaining)
	}

	// Clear retry key on success
	_ = s.redisClient.Client.Del(ctx, retryKey).Err()
	return true, nil
}

func (s *twilioOtpService) buildRetryKey(phone string) string {
	return fmt.Sprintf("otp:%s:%s", phone, twilioRetryKeySuffix)
}

func (s *twilioOtpService) isBlocked(ctx context.Context, retryKey string) (bool, time.Duration) {
	count, err := s.getRetryCount(ctx, retryKey)
	if err != nil {
		return false, 0
	}
	if count >= twilioTotalRetries {
		ttl, _ := s.redisClient.Client.TTL(ctx, retryKey).Result()
		return true, ttl
	}
	return false, 0
}

func (s *twilioOtpService) getRetryCount(ctx context.Context, retryKey string) (int, error) {
	countStr, err := s.redisClient.Client.Get(ctx, retryKey).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(countStr)
}

func (s *twilioOtpService) incrementRetry(ctx context.Context, retryKey string) int {
	pipe := s.redisClient.Client.TxPipeline()
	count := pipe.Incr(ctx, retryKey)
	pipe.Expire(ctx, retryKey, s.expiry)
	pipe.Exec(ctx)
	c, _ := count.Result()
	return int(c)
}
