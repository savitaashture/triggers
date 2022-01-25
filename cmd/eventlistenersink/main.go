/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"github.com/tektoncd/triggers/pkg/adapter"
	dynamicClientset "github.com/tektoncd/triggers/pkg/client/dynamic/clientset"
	"github.com/tektoncd/triggers/pkg/client/dynamic/clientset/tekton"
	"github.com/tektoncd/triggers/pkg/sink"
	"k8s.io/client-go/dynamic"
	evadapter "knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/clients/dynamicclient"
	"knative.dev/pkg/signals"
	"log"
)

const (
	// EventListenerLogKey is the name of the logger for the eventlistener cmd
	EventListenerLogKey = "eventlistener"
)

func main() {
	ctx := signals.NewContext()

	cfg := injection.ParseAndGetRESTConfigOrDie()

	dc := dynamic.NewForConfigOrDie(cfg)
	dc = dynamicClientset.New(tekton.WithClient(dc))
	ctx = context.WithValue(ctx, dynamicclient.Key{}, dc)

	// Set up ctx with the set of things based on the
	// dynamic client we've set up above.
	ctx = injection.Dynamic.SetupDynamic(ctx)

	sinkArgs, err := sink.GetArgs()
	if err != nil {
		log.Fatal(err.Error())
	}
	sinkClients, err := sink.ConfigureClients(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	recorder, err := sink.NewRecorder()
	if err != nil {
		log.Fatal(err.Error())
	}

	if !sinkArgs.IsMultiNS {
		ctx = injection.WithNamespaceScope(ctx, sinkArgs.ElNamespace)
	}

	//ctx, _ = injection.EnableInjectionOrDie(ctx, cfg)
	//
	//logger, atomicLevel := sharedmain.SetupLoggerOrDie(ctx, EventListenerLogKey)
	//defer flush(logger)
	//ctx = logging.WithLogger(ctx, logger)
	//cmw := sharedmain.SetupConfigMapWatchOrDie(ctx, logger)
	//store := config.NewStore(logging.FromContext(ctx).Named("config-store"))
	//store.WatchConfigs(cmw)
	//sharedmain.WatchLoggingConfigOrDie(ctx, cmw, logger, atomicLevel, EventListenerLogKey)
	// Get mww object from knative pkg and try to do watch here
	//  use envFrom and env as an envirnotment

	evadapter.MainWithContext(ctx, EventListenerLogKey, adapter.NewEnvConfig, adapter.New(sinkArgs, sinkClients, recorder))
}

//func flush(logger *zap.SugaredLogger) {
//	logger.Sync()
//	metrics.FlushExporter()
//}