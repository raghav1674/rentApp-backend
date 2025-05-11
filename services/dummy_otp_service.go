package services

import (
	"context"
	"fmt"
	"math/rand"
	"sample-web/clients"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type PhoneNumberBlockedError struct {
	PhoneNumber string
	RetryAfter  time.Duration
}

func (e *PhoneNumberBlockedError) Error() string {
	return fmt.Sprintf("phone number %s is blocked, retry after %v", e.PhoneNumber, e.RetryAfter)
}

const (
	expiryTimeInMinutes = 1
	totalRetries        = 3
	retryKeySuffix      = "invalid_attempts"
)

type dummyOtpService struct {
	redisClient *clients.RedisClient
	expiry      time.Duration
}

func NewDummyOTPService(redisClient *clients.RedisClient) OTPService {
	return &dummyOtpService{
		redisClient: redisClient,
		expiry:      expiryTimeInMinutes * time.Minute,
	}
}

// SendOTP generates and stores OTP, enforcing retry limit.
func (s *dummyOtpService) SendOTP(ctx context.Context, phoneNumber string) (string, error) {
	retryKey := s.buildRetryKey(phoneNumber)

	// Check for retry limit
	if blocked, ttl := s.isBlocked(ctx, retryKey); blocked {
		return "", &PhoneNumberBlockedError{PhoneNumber: phoneNumber, RetryAfter: ttl}
	}

	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	if err := s.redisClient.Client.Set(ctx, phoneNumber, otp, s.expiry).Err(); err != nil {
		return "", err
	}
	return otp, nil
}

// VerifyOTP checks the provided OTP and updates retry state.
func (s *dummyOtpService) VerifyOTP(ctx context.Context, phoneNumber, code string) (bool, error) {

	retryKey := s.buildRetryKey(phoneNumber)

	isBlocked, ttl := s.isBlocked(ctx, retryKey)
	if isBlocked {
		return false, &PhoneNumberBlockedError{PhoneNumber: phoneNumber, RetryAfter: ttl}
	}

	otp, err := s.redisClient.Client.Get(ctx, phoneNumber).Result()
	if err != nil {
		if err == redis.Nil {
			return false, fmt.Errorf("OTP not found or expired")
		}
		return false, err
	}

	if otp != code {
		pipe := s.redisClient.Client.TxPipeline()
		pipe.Incr(ctx, retryKey)
		pipe.Expire(ctx, retryKey, s.expiry)
		_, _ = pipe.Exec(ctx)

		remaining := totalRetries - s.getRetryCount(ctx, retryKey)
		if remaining <= 0 {
			return false, &PhoneNumberBlockedError{PhoneNumber: phoneNumber, RetryAfter: s.expiry}
		}

		return false, fmt.Errorf("invalid OTP, %d attempt(s) remaining", remaining)
	}

	pipe := s.redisClient.Client.TxPipeline()
	pipe.Del(ctx, phoneNumber, retryKey)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

// buildRetryKey builds a Redis key for tracking attempts.
func (s *dummyOtpService) buildRetryKey(phone string) string {
	return fmt.Sprintf("otp:%s:%s", phone, retryKeySuffix)
}

// isBlocked checks if a phone number is rate-limited.
func (s *dummyOtpService) isBlocked(ctx context.Context, retryKey string) (bool, time.Duration) {
	countStr := s.redisClient.Client.Get(ctx, retryKey).Val()
	count, _ := strconv.Atoi(countStr)

	if count >= totalRetries {
		ttl, _ := s.redisClient.Client.TTL(ctx, retryKey).Result()
		return true, ttl
	}
	return false, 0
}

// getRetryCount returns the current retry count for a phone number.
func (s *dummyOtpService) getRetryCount(ctx context.Context, retryKey string) int {
	countStr := s.redisClient.Client.Get(ctx, retryKey).Val()
	count, _ := strconv.Atoi(countStr)
	return count
}
