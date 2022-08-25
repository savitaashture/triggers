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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	//logger "github.com/sirupsen/logrus"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"os"
	//"os/signal"
	"sync"
	//"syscall"

	//certresources "knative.dev/pkg/webhook/certificates/resources"
	"log"
	"net"
	"net/http"
	"time"
	"crypto/tls"

	triggersclientset "github.com/tektoncd/triggers/pkg/client/clientset/versioned"
	clusterinterceptorsinformer "github.com/tektoncd/triggers/pkg/client/injection/informers/triggers/v1alpha1/clusterinterceptor"
	"github.com/tektoncd/triggers/pkg/interceptors"
	"github.com/tektoncd/triggers/pkg/interceptors/server"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/labels"
	kubeclientset "k8s.io/client-go/kubernetes"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

const (
	// HTTPSPort is the port where interceptor service listens on
	HTTPSPort    = 8443
	readTimeout  = 5 * time.Second
	writeTimeout = 20 * time.Second
	idleTimeout  = 60 * time.Second
)

type keypairReloader struct {
	certMu   sync.RWMutex
	cert     *tls.Certificate
	certPath string
	keyPath  string
}

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
	if err != nil {
		logger.Errorf("failed to create new Clientset for the given config: %v", err)
		return
	}

	service, err := server.NewWithCoreInterceptors(interceptors.NewKubeClientSecretGetter(kubeclient.Get(ctx).CoreV1(), 1024, 5*time.Second), logger)
	if err != nil {
		logger.Errorf("failed to initialize core interceptors: %s", err)
		return
	}
	startInformer()

	mux := http.NewServeMux()
	mux.Handle("/", service)
	mux.HandleFunc("/ready", handler)

	//keyFile, certFile, serverKey, serverCert, caCert, err := server.CreateCerts(ctx, kubeClient.CoreV1(), logger)
	keyFile, certFile, _, serverCert, caCert, err := server.CreateCerts(ctx, kubeClient.CoreV1(), logger, time.Minute)
	if err != nil {
		return
	}

	tc, err := triggersclientset.NewForConfig(cfg)
	if err != nil {
		return
	}

	//ticker1 := time.NewTicker(time.Minute)
	//quit1 := make(chan struct{})
	//var c = []tls.Certificate{}
	//var (
	//	cert1 *x509.Certificate
	//	er  error
	//	srv *http.Server
	//	//certrrr x509.Certificate
	//)
	//go func() {
	//	for {
	//		select {
	//		case <-ticker1.C:
	//			// Check the expiration date of the certificate to see if it needs to be updated
	//			roots := x509.NewCertPool()
	//			ok := roots.AppendCertsFromPEM(caCert)
	//			if !ok {
	//				logger.Error("failed to parse root certificate")
	//			}
	//			block, _ := pem.Decode(serverCert)
	//			if block == nil {
	//				logger.Error("failed to parse certificate PEM")
	//			} else {
	//				cert1, er = x509.ParseCertificate(block.Bytes)
	//				if err != nil {
	//					logger.Errorf("failed to parse certificate: %v", er.Error())
	//				}
	//			}
	//
	//			opts := x509.VerifyOptions{
	//				Roots: roots,
	//			}
	//
	//			if _, er := cert1.Verify(opts); er != nil {
	//				logger.Errorf("failed to verify certificate: %v", er.Error())
	//
	//				//keyFile1, certFile1, _, _, caCert1, er := server.CreateCerts(ctx, kubeClient.CoreV1(), logger)
	//				keyFile1, certFile1, _, _, caCert1, er := server.CreateCerts(ctx, kubeClient.CoreV1(), logger)
	//				if er != nil {
	//					logger.Errorf(er.Error())
	//
	//				}
	//
	//				certrrr, er := tls.LoadX509KeyPair(certFile1, keyFile1)
	//				fmt.Println("what is error of LoadX509KeyPair", certrrr, "&*************&&&", er)
	//				srv = &http.Server{
	//					Addr: fmt.Sprintf(":%d", HTTPSPort),
	//					BaseContext: func(listener net.Listener) context.Context {
	//						return ctx
	//					},
	//					ReadTimeout:  readTimeout,
	//					WriteTimeout: writeTimeout,
	//					IdleTimeout:  idleTimeout,
	//					Handler:      mux,
	//					TLSConfig: &tls.Config{
	//						Certificates: append(c, certrrr),
	//					},
	//				}
	//
	//				fmt.Println("the server data herer", srv.TLSConfig.Certificates)
	//
	//				if er := listAndUpdateClusterInterceptorCRD(ctx, tc, service, caCert1, true); er != nil {
	//					logger.Errorf(er.Error())
	//				}
	//
	//			}
	//		case <-quit1:
	//			ticker1.Stop()
	//			return
	//		}
	//	}
	//}()

//	checkCertValidity(serverCert, caCert, logger)

	if err := listAndUpdateClusterInterceptorCRD(ctx, tc, service, caCert, false); err != nil {
		return
	}
	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := listAndUpdateClusterInterceptorCRD(ctx, tc, service, caCert, false); err != nil {
					return
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	//cert, err := tls.LoadX509KeyPair(certFile, keyFile)
//	fmt.Println("what is error of LoadX509KeyPair", cert, "&*************&&&", err)

	//var c = []tls.Certificate{}

	kpr, err := NewKeypairReloader(ctx, certFile, keyFile, caCert, serverCert,kubeClient.CoreV1(), logger, tc, service)
	fmt.Println("what is errororor herere", err)
	if err != nil {
		log.Fatal(err)
	}


	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", HTTPSPort),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      mux,
		TLSConfig: &tls.Config{
			GetCertificate: func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				fmt.Println("11111111111")
				//kpr.certMu.RLock()
				fmt.Println("22222222222222222222222222")
				//defer kpr.certMu.RUnlock()
				return kpr.cert, nil
			},
		},
		//TLSConfig: &tls.Config{
		//	Certificates: append(c, cert),
		//},
	}
	//fmt.Println("PPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPP",  kpr.GetCertificateFunc())
	//srv.TLSConfig.GetCertificate = kpr.GetCertificateFunc()
	//srv.TLSConfig.GetCertificate = kpr.cert
	fmt.Println("HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHh")
	//fmt.Println("what is server certificate datattatatata********", srv.TLSConfig.Certificates)
	logger.Infof("Listen and serve on port %d", HTTPSPort)
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		logger.Fatalf("failed to start interceptors service: %v", err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func listAndUpdateClusterInterceptorCRD(ctx context.Context, tc *triggersclientset.Clientset, service *server.Server, caCert []byte, certIsInvalid bool) error {
	clusterInterceptorList, err := clusterinterceptorsinformer.Get(ctx).Lister().List(labels.NewSelector())
	if err != nil {
		return err
	}

	if err := service.UpdateCRDWithCaCert(ctx, tc.TriggersV1alpha1(), clusterInterceptorList, caCert, certIsInvalid); err != nil {
		fmt.Println("is it erroring oit",err)
		return err
	}
	return nil
}

func NewKeypairReloader(ctx context.Context, certPath, keyPath string, caCert, serverCert []byte, coreV1Interface corev1.CoreV1Interface,
	logger *zap.SugaredLogger, tc *triggersclientset.Clientset, service *server.Server) (*keypairReloader, error) {
	result := &keypairReloader{
		certPath: certPath,
		keyPath:  keyPath,
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	result.cert = &cert
	if err := result.maybeReload(); err != nil {
		log.Printf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
	}
	fmt.Println("for the first time")
	//go func() {
	//
	//	c := make(chan os.Signal, 1)
	//	signal.Notify(c, syscall.SIGHUP)
	//	for range c {
	//		log.Printf("Received SIGHUP, reloading TLS certificate and key from %q and %q", certPath, keyPath)
	//		if err := result.maybeReload(); err != nil {
	//			log.Printf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
	//		}
	//	}
	//}()
	ticker1 := time.NewTicker(time.Minute)
	quit1 := make(chan struct{})
	//var c = []tls.Certificate{}
	var (
		cert1 *x509.Certificate
		//err  error
		//srv *http.Server
		//certrrr x509.Certificate
	)
	go func() {
		for {
			select {
			case <-ticker1.C:
				// Check the expiration date of the certificate to see if it needs to be updated
				roots := x509.NewCertPool()
				ok := roots.AppendCertsFromPEM(caCert)
				if !ok {
					logger.Error("failed to parse root certificate")
				}
				block, _ := pem.Decode(serverCert)
				if block == nil {
					logger.Error("failed to parse certificate PEM")
				} else {
					cert1, err = x509.ParseCertificate(block.Bytes)
					if err != nil {
						logger.Errorf("failed to parse certificate: %v", err.Error())
					}
				}

				opts := x509.VerifyOptions{
					Roots: roots,
				}

				if _, err := cert1.Verify(opts); err != nil {
					logger.Errorf("failed to verify certificate: %v", err.Error())

					//keyFile1, certFile1, _, _, caCert1, er := server.CreateCerts(ctx, kubeClient.CoreV1(), logger)
					keyFile1, certFile1, _, _, caCert1, err := server.CreateCerts(ctx, coreV1Interface, logger, 3*time.Minute)
					fmt.Println("createCert eror", err)
					if err != nil {
						logger.Errorf(err.Error())

					}

					result = &keypairReloader{
						certPath: certFile1,
						keyPath:  keyFile1,
					}
					cert, err := tls.LoadX509KeyPair(certFile1, keyFile1)
					if err != nil {
						fmt.Println("errororor inside go func of loop", err)
					}
					result.cert = &cert
					if err := result.maybeReload(); err != nil {
						fmt.Println("does reolad success")
						log.Printf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
					}

					if er := listAndUpdateClusterInterceptorCRD(ctx, tc, service, caCert1, true); er != nil {
						fmt.Println("how cme erroror", er)
						logger.Errorf(er.Error())
					}

				}
			case <-quit1:
				ticker1.Stop()
				return
			}
		}
	}()
	return result, nil
}

func (kpr *keypairReloader) maybeReload() error {
	newCert, err := tls.LoadX509KeyPair(kpr.certPath, kpr.keyPath)
	if err != nil {
		return err
	}
	kpr.certMu.Lock()
	defer kpr.certMu.Unlock()
	kpr.cert = &newCert
	return nil
}

func (kpr *keypairReloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	fmt.Println("oooooooooooooooooooo", kpr.certMu)
	fmt.Println("oooooooooooooooooooo111111111111", kpr.certMu.RLock)
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		fmt.Println("11111111111")
		//kpr.certMu.RLock()
		fmt.Println("22222222222222222222222222")
		//defer kpr.certMu.RUnlock()
		return kpr.cert, nil
	}
}
