// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/tylerb/graceful"

	"fmt"
	"ocopea/k8s-mongodsb/mongo"
	"ocopea/k8s-mongodsb/restapi/impl"
	"ocopea/k8s-mongodsb/restapi/operations"
	"ocopea/k8s-mongodsb/restapi/operations/dsb_web"
	k8sClient "ocopea/kubernetes/client"
	"os"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name k8s-mongodsb --spec ../dsb-swagger.yaml
var k8s *k8sClient.Client
var mongoClient *mongo.MongoClient

func configureFlags(api *operations.K8sMongodsbAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.K8sMongodsbAPI) http.Handler {
	var k8sURL string
	host, hb := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	port, pb := os.LookupEnv("KUBERNETES_SERVICE_PORT")
	k8sUserName := os.Getenv("K8S_USERNAME")
	k8sPassword := os.Getenv("K8S_PASSWORD")
	if hb && pb {
		k8sURL = "https://" + host + ":" + port
	} else {
		panic("could not find k8s host vars")
	}
	nazNS := os.Getenv("OCOPEA_NAMESPACE")
	if len(nazNS) == 0 {
		nazNS = "ocopea"
	}
	fmt.Printf("url %s\nnamespace:%s\n", k8sURL, nazNS)

	k8s, err := k8sClient.NewClient(
		k8sURL,
		nazNS,
		k8sUserName,
		k8sPassword,
		"/var/run/secrets/kubernetes.io/serviceaccount/token")

	if err != nil {
		panic(err)
	}

	mongoClient = &mongo.MongoClient{}

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.DsbWebCopyServiceInstanceHandler = dsb_web.CopyServiceInstanceHandlerFunc(func(params dsb_web.CopyServiceInstanceParams) middleware.Responder {
		return impl.CopyServiceInstance(k8s, mongoClient, params)
	})
	api.DsbWebCreateServiceInstanceHandler = dsb_web.CreateServiceInstanceHandlerFunc(func(params dsb_web.CreateServiceInstanceParams) middleware.Responder {
		return impl.CreateInstanceResponse(k8s, mongoClient, params)
	})
	api.DsbWebDeleteServiceInstanceHandler = dsb_web.DeleteServiceInstanceHandlerFunc(func(params dsb_web.DeleteServiceInstanceParams) middleware.Responder {
		return impl.DeleteInstanceResponse(k8s, params)
	})
	api.DsbWebGetDSBIconHandler = dsb_web.GetDSBIconHandlerFunc(func(params dsb_web.GetDSBIconParams) middleware.Responder {
		return middleware.NotImplemented("operation dsb_web.GetDSBIcon has not yet been implemented")
	})
	api.DsbWebGetDSBInfoHandler = dsb_web.GetDSBInfoHandlerFunc(func(params dsb_web.GetDSBInfoParams) middleware.Responder {
		return impl.DsbInfoResponse()
	})
	api.DsbWebGetServiceInstanceHandler = dsb_web.GetServiceInstanceHandlerFunc(func(params dsb_web.GetServiceInstanceParams) middleware.Responder {
		return impl.DsbGetServiceInstancesResponse(k8s, mongoClient, params)
	})
	api.DsbWebGetServiceInstancesHandler = dsb_web.GetServiceInstancesHandlerFunc(func(params dsb_web.GetServiceInstancesParams) middleware.Responder {
		return impl.DsbInstancesResponse()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
