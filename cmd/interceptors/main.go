/*
Copyright 2020 The Tekton Authors

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
	"io/ioutil"
	kubeclientset "k8s.io/client-go/kubernetes"

	//"crypto/tls"
	"fmt"
	//"strings"

	//"io"
	//"io/ioutil"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	//	"bufio"

	//"strings"

	"log"
	"net"
	"net/http"
	//kubeclientset "k8s.io/client-go/kubernetes"
	"os"
	"time"

	"github.com/tektoncd/triggers/pkg/interceptors/server"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	secretInformer "knative.dev/pkg/client/injection/kube/informers/core/v1/secret"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
	certresources "knative.dev/pkg/webhook/certificates/resources"
	//"os/exec"
)

const (
	// Port is the port that the port that interceptor service listens on
	HTTPSPort    = 8443
	readTimeout  = 5 * time.Second
	writeTimeout = 20 * time.Second
	idleTimeout  = 60 * time.Second

	//oneWeek = 7 * 24 * time.Hour
	oneWeek = 10 * time.Minute
	oneDay  = 24 * time.Hour

	keyFile  = "/tmp/server-key.pem"
	certFile = "/tmp/server-cert.pem"
)

func main() {
	// set up signals so we handle the first shutdown signal gracefully
	ctx := signals.NewContext()

	cfg := injection.ParseAndGetRESTConfigOrDie()

	ctx, startInformer := injection.EnableInjectionOrDie(ctx, cfg)

	zap, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %s", err)
	}
	logger := zap.Sugar()
	ctx = logging.WithLogger(ctx, logger)
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("failed to sync the logger: %s", err)
		}
	}()

	kubeClient, err := kubeclientset.NewForConfig(cfg)

	secretLister := secretInformer.Get(ctx).Lister()
	service, err := server.NewWithCoreInterceptors(secretLister, logger)
	if err != nil {
		log.Printf("failed to initialize core interceptors: %s", err)
		return
	}
	startInformer()

	mux := http.NewServeMux()
	mux.Handle("/", service)
	mux.HandleFunc("/ready", handler)

	name := os.Getenv("SVC_NAME")
	namespace := os.Getenv("SYSTEM_NAMESPACE")
	fmt.Println("namespace is", namespace)

	secret, err := secretLister.Secrets(namespace).Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// The secret should be created explicitly by a higher-level system
			// that's responsible for install/updates.  We simply populate the
			// secret information.
			fmt.Println("not fund")
			log.Printf("secret %s is missing", name)
			logger.Infof("secret %s is missing", name)
			return
		}
		logger.Errorf("error accessing certificate secret %q: %v", name, err)
		return
	}

	fmt.Println("secretNamesecretNamesecretName", secret.Data)

	//if _, haskey := secret.Data[certresources.ServerKey]; !haskey {
	//	fmt.Println("UUUUUUUUU")
	//	logger.Infof("Certificate secret %q is missing key %q", name, certresources.ServerKey)
	//} else if _, haskey := secret.Data[certresources.ServerCert]; !haskey {
	//	logger.Infof("Certificate secret %q is missing key %q", name, certresources.ServerCert)
	//} else if _, haskey := secret.Data[certresources.CACert]; !haskey {
	//	logger.Infof("Certificate secret %q is missing key %q", name, certresources.CACert)
	//} else {
	//	fmt.Println("UUUUUUUU11111111111111U")
	//	// Check the expiration date of the certificate to see if it needs to be updated
	//	cert, err := tls.X509KeyPair(secret.Data[certresources.ServerCert], secret.Data[certresources.ServerKey])
	//	if err != nil {
	//		logger.Warnw("Error creating pem from certificate and key", err.Error())
	//	} else {
	//		certData, err := x509.ParseCertificate(cert.Certificate[0])
	//		if err != nil {
	//			logger.Errorw("Error parsing certificate", err.Error())
	//		} else if time.Now().Add(oneDay).Before(certData.NotAfter) {
	//			return
	//		}
	//	}
	//}
	//fmt.Println("22222222222222222222222222222222222")

	// Don't modify the informer copy.
	secret = secret.DeepCopy()

	serverKey, serverCert, caCert, err := certresources.CreateCerts(ctx, name, namespace, time.Now().Add(oneWeek))
	fmt.Println("errrrrrrr value", err)
	if err != nil {
		fmt.Println("errororororor in createCerts", err)
		return
	}
	newSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			certresources.ServerKey:  serverKey,
			certresources.ServerCert: serverCert,
			certresources.CACert:     caCert,
		},
	}
	secret.Data = newSecret.Data
	//fmt.Println("newSecretnewSecretnewSecret", newSecret.Data)
	fmt.Println("kubeClientkubeClientkubeClient", kubeClient)
	updatedSecret, err := kubeClient.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return
	}
	fmt.Println("eeeeeeeeeeeeeeeeeeeeeeeeeeee", err)
	fmt.Println("updated secret value", updatedSecret.Data)

	err = ioutil.WriteFile(keyFile, serverKey, 0600)
	fmt.Println("writing file is errorororororor", err)
	fmt.Println("IIIIIIIIIIIIIIIIIIIIIII111 crt", err)
	err = ioutil.WriteFile(certFile, serverCert, 0644)
	fmt.Println("writing file is errorororororor", err)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", HTTPSPort),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      mux,
	}
	logger.Infof("Listen and serve on port %d", HTTPSPort)
	if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
		logger.Fatalf("failed to start interceptors service: %v", err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
