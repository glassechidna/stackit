package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/glassechidna/stackit/cmd/honey"
	"log"
	"os"
)

func awsSession(profile, region string) *session.Session {
	sessOpts := session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}

	if len(profile) > 0 {
		sessOpts.Profile = profile
	}

	if len(region) > 0 {
		sessOpts.Config.Region = aws.String(region)
	}

	if len(os.Getenv("STACKIT_AWS_VERBOSE")) > 0 {
		logger := log.New(os.Stderr, "", log.LstdFlags)
		sessOpts.Config.LogLevel = aws.LogLevel(aws.LogDebugWithHTTPBody)
		sessOpts.Config.Logger = aws.LoggerFunc(func(args ...interface{}) {
			logger.Println(args...)
		})
	}

	userAgentHandler := request.NamedHandler{
		Name: "stackit.UserAgentHandler",
		Fn:   request.MakeAddToUserAgentHandler("stackit", version),
	}

	sess, _ := session.NewSessionWithOptions(sessOpts)
	sess.Handlers.Build.PushFrontNamed(userAgentHandler)
	honey.TryAdd(sess)

	return sess
}
