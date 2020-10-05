package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() (*session.Session, error) {
	sess, err := Session(AssumeRoleTokenProvider(StderrTokenProvider))
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// SessionOption is an option for session.
type SessionOption func(*session.Options)

// AssumeRoleTokenProvider is an option for setting custom assume role token provider.
func AssumeRoleTokenProvider(provider func() (string, error)) SessionOption {
	return func(options *session.Options) {
		options.AssumeRoleTokenProvider = provider
	}
}

func Session(opts ...SessionOption) (*session.Session, error) {
	so := session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: StderrTokenProvider,
	}

	for _, opt := range opts {
		opt(&so)
	}

	return session.NewSessionWithOptions(so)
}

// StderrTokenProvider implements token provider for AWS SDK.
func StderrTokenProvider() (string, error) {
	var v string
	fmt.Fprintf(os.Stderr, "Assume Role MFA token code: ")
	_, err := fmt.Fscanln(os.Stderr, &v)
	return v, err
}
