package honey

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/trace"
	"os"
)

func RootContext() (context.Context, func()) {
	ctx := context.Background()
	if !isConfigured {
		return ctx, func() {}
	}

	ctx, span := beeline.StartSpan(ctx, "stackit")
	return ctx, span.Send
}

var isConfigured bool

func init() {
	writeKey := os.Getenv("HONEYCOMB_WRITEKEY")
	dataset := os.Getenv("HONEYCOMB_DATASET")
	if len(writeKey) == 0 || len(dataset) == 0 {
		return
	}

	beeline.Init(beeline.Config{
		WriteKey: writeKey,
		Dataset:  dataset,
	})

	isConfigured = true
}

func TryAdd(sess *session.Session) {
	if !isConfigured {
		return
	}

	sess.Handlers.Build.PushFrontNamed(request.NamedHandler{
		Name: "honeycomb.build",
		Fn:   honeycombStartAws,
	})
	sess.Handlers.Complete.PushFrontNamed(request.NamedHandler{
		Name: "honeycomb.complete",
		Fn:   honeycombCompleteAws,
	})
}

func honeycombStartAws(r *request.Request) {
	name := fmt.Sprintf("%s %s", r.ClientInfo.ServiceName, r.Operation.Name)
	ctx, span := beeline.StartSpan(r.Context(), name)
	span.AddField("aws.service", r.ClientInfo.ServiceName)
	span.AddField("aws.action", r.Operation.Name)
	span.AddField("aws.input", r.Params)
	r.SetContext(ctx)
}

func honeycombCompleteAws(r *request.Request) {
	span := trace.GetSpanFromContext(r.Context())
	span.AddField("aws.output", r.Data)

	if awsErr, ok := r.Error.(awserr.Error); ok {
		span.AddField("aws.error.code", awsErr.Code())
		span.AddField("aws.error.message", awsErr.Message())
	}
	span.Send()
}
