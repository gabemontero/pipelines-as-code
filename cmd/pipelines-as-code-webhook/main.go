package main

import (
	"context"
	"log"
	"os"

	validation_webhook "github.com/openshift-pipelines/pipelines-as-code/pkg/webhook"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"
	"knative.dev/pkg/webhook"
	"knative.dev/pkg/webhook/certificates"

	"github.com/openshift-pipelines/pipelines-as-code/pkg/params"
)

func main() {
	ctx := signals.NewContext()
	// set up client/informer overrides for kcp
	run := params.New()
	err := run.Clients.NewClients(ctx, &run.Info)
	if err != nil {
		log.Fatal("failed to init clients : ", err)
	}
	run.Informers.Clients = run.Clients
	run.Informers.NewInformers(ctx)

	serviceName := os.Getenv("WEBHOOK_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "pipelines-as-code-webhook"
	}
	secretName := os.Getenv("WEBHOOK_SECRET_NAME")
	if secretName == "" {
		secretName = "pipelines-as-code-webhook-certs"
	}
	// Set up a signal context with our webhook options
	ctx = webhook.WithOptions(ctx, webhook.Options{
		ServiceName: serviceName,
		Port:        8443,
		SecretName:  secretName,
	})

	sharedmain.WebhookMainWithConfig(ctx, "pipelines-as-code-webhook",
		injection.ParseAndGetRESTConfigOrDie(),
		certificates.NewController,
		newValidationAdmissionController,
	)
}

func newValidationAdmissionController(ctx context.Context, _ configmap.Watcher) *controller.Impl {
	return validation_webhook.NewAdmissionController(ctx,

		// Name of the resource webhook.
		"validation.pipelinesascode.tekton.dev",

		// The path on which to serve the webhook.
		"/validate",

		// A function that infuses the context passed to Validate/SetDefaults with custom metadata.
		func(ctx context.Context) context.Context {
			return ctx
		},

		// Whether to disallow unknown fields.
		true,
	)
}
