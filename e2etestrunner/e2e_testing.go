// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2etestrunner

import (
	"context"
	"log"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-e2e-testing/e2etestrunner/testclient"
)

type ApplyPersistent struct {
	AutoApprove bool `arg:"--auto-approve" default:"false" help:"Approve without prompting. Default is false."`
}

type CmdWithProjectId struct {
	ProjectID string `arg:"required,--project-id,env:PROJECT_ID" help:"GCP project id/name"`
}

type CmdWithImage struct {
	Image string `arg:"required" help:"docker container image to deploy and test"`
}

type LocalCmd struct {
	CmdWithImage

	Port string `default:"8000"`

	// Needed when running without a metadata server for credentials
	GoogleApplicationCredentials string `arg:"--google-application-credentials,env:GOOGLE_APPLICATION_CREDENTIALS" help:"Path to google credentials key file to mount into test server container"`

	// May be needed when running this binary in a container
	Network string `help:"Docker network to use when starting the container, optional"`

	ContainerUser string `arg:"--container-user" help:"Optional user to use when running the container"`
}

type GceCmd struct {
	CmdWithImage
}

type GkeCmd struct {
	CmdWithImage
}

type GaeCmd struct {
	CmdWithImage

	Runtime string `arg:"required" help:"The language runtime for the instrumented test server, used in naming the service"`
}

type GaeStandardCmd struct {
	Runtime    string `arg:"required" help:"The language runtime for the instrumented test server, used in naming the service"`
	AppSource  string `arg:"required" help:"The full path of the zip file that contains the source code to run in GAE"`
	Entrypoint string `help:"Optional entrypoint to control how GAE starts the application. See https://cloud.google.com/appengine/docs/standard/reference/app-yaml#entrypoint"`
}

type CloudRunCmd struct {
	CmdWithImage
}

type CloudFunctionsGen2Cmd struct {
	// Needed to configure which language will the function instance support
	Runtime string `arg:"required" help:"Configure the language runtime environment for CloudFunction"`

	EntryPoint string `arg:"required" help:"The entry point for the code that will run in this CloudFunction"`

	FunctionSource string `arg:"required" help:"The full path of the zip file that contains the code source that needs to run within the CloudFunction"`
}

type Args struct {
	// This subcommand is a special case, it doesn't run any tests. It just
	// applies the persistent resources which are used across tests. See
	// tf/persistent/README.md for details on what is in there.
	ApplyPersistent *ApplyPersistent `arg:"subcommand:apply-persistent" help:"Terraform apply the resources in tf/persistent and exit (does not run tests)."`

	Local              *LocalCmd              `arg:"subcommand:local" help:"Deploy the test server locally with docker and execute tests"`
	Gke                *GkeCmd                `arg:"subcommand:gke" help:"Deploy the test server on GKE and execute tests"`
	Gce                *GceCmd                `arg:"subcommand:gce" help:"Deploy the test server on GCE and execute tests"`
	Gae                *GaeCmd                `arg:"subcommand:gae" help:"Deploy the test server on GAE and execute tests"`
	GaeStandard        *GaeStandardCmd        `arg:"subcommand:gae-standard" help:"Deploy the test server on GAE standard and execute tests"`
	CloudRun           *CloudRunCmd           `arg:"subcommand:cloud-run" help:"Deploy the test server on Cloud Run and execute tests"`
	CloudFunctionsGen2 *CloudFunctionsGen2Cmd `arg:"subcommand:cloud-functions-gen2" help:"Deploy the test server on Cloud Function (2nd Gen) and execute tests"`

	CmdWithProjectId
	GoTestFlags        string        `help:"go test flags to pass through, e.g. --gotestflags='-test.v'"`
	HealthCheckTimeout time.Duration `arg:"--health-check-timeout" help:"A duration (e.g. 5m) to wait for the test server health check. Default is 2m." default:"2m"`

	// This is used in a new terraform workspace's name and in the GCP resources
	// we create. Pass the GCB build ID in CI to get the build id formatted into
	// resources created for debugging. If not provided, we generate a hex
	// string.
	TestRunID string `arg:"--test-run-id,env:TEST_RUN_ID" help:"Optional test run id to use to partition terraform resources"`
}

type Cleanup func()
type SetupFunc func(
	context.Context,
	*Args,
	*log.Logger,
) (*testclient.Client, Cleanup, error)

func NoopCleanup() {}
