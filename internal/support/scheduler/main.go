/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package scheduler

import (
	"context"
	"flag"

	"github.com/edgexfoundry/edgex-go"
	"github.com/edgexfoundry/edgex-go/internal"
	"github.com/edgexfoundry/edgex-go/internal/pkg/bootstrap/handlers/database"
	"github.com/edgexfoundry/edgex-go/internal/pkg/telemetry"
	"github.com/edgexfoundry/edgex-go/internal/pkg/usage"
	"github.com/edgexfoundry/edgex-go/internal/support/scheduler/config"
	"github.com/edgexfoundry/edgex-go/internal/support/scheduler/container"

	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/handlers/httpserver"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/handlers/message"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/handlers/secret"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/handlers/testing"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/interfaces"
	"github.com/edgexfoundry/go-mod-bootstrap/bootstrap/startup"
	"github.com/edgexfoundry/go-mod-bootstrap/di"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"

	"github.com/gorilla/mux"
)

func Main(ctx context.Context, cancel context.CancelFunc, router *mux.Router, readyStream chan<- bool) {
	startupTimer := startup.NewStartUpTimer(internal.BootRetrySecondsDefault, internal.BootTimeoutSecondsDefault)
	var useRegistry bool
	var configDir, profileDir string

	flag.BoolVar(&useRegistry, "registry", false, "Indicates the service should use Registry.")
	flag.BoolVar(&useRegistry, "r", false, "Indicates the service should use Registry.")
	flag.StringVar(&profileDir, "profile", "", "Specify a profile other than default.")
	flag.StringVar(&profileDir, "p", "", "Specify a profile other than default.")
	flag.StringVar(&configDir, "confdir", "", "Specify local configuration directory")

	flag.Usage = usage.HelpCallback
	flag.Parse()

	configuration := &config.ConfigurationStruct{}
	dic := di.NewContainer(di.ServiceConstructorMap{
		container.ConfigurationName: func(get di.Get) interface{} {
			return configuration
		},
	})

	httpServer := httpserver.NewBootstrap(router)

	bootstrap.Run(
		ctx,
		cancel,
		configDir,
		profileDir,
		internal.ConfigFileName,
		useRegistry,
		clients.SupportSchedulerServiceKey,
		configuration,
		startupTimer,
		dic,
		[]interfaces.BootstrapHandler{
			secret.NewSecret().BootstrapHandler,
			database.NewDatabase(httpServer, configuration, readyStream != nil).BootstrapHandler,
			NewBootstrap(router).BootstrapHandler,
			telemetry.BootstrapHandler,
			httpServer.BootstrapHandler,
			message.NewBootstrap(clients.SupportSchedulerServiceKey, edgex.Version).BootstrapHandler,
			testing.NewBootstrap(httpServer, readyStream).BootstrapHandler,
		})
}
