// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"log"
	"ocopea/k8s-mongodsb/models"
	"ocopea/k8s-mongodsb/mongo"
	"ocopea/k8s-mongodsb/restapi/operations/dsb_web"
	k8sClient "ocopea/kubernetes/client"
	"ocopea/kubernetes/client/v1"
)

type MongoBindingInfo struct {
	Server string `json:"server"`
	Port   string `json:"port"`
}

func DsbGetServiceInstancesResponse(
	k8s k8sClient.ClientInterface,
	mongoClient mongo.MongoInterface,
	params dsb_web.GetServiceInstanceParams) middleware.Responder {

	d, err := getBindingInfoForInstance(k8s, mongoClient, params.InstanceID)
	if err != nil {
		return formatGetServiceError(500, err.Error())
	}

	return dsb_web.NewGetServiceInstanceOK().WithPayload(d)

}

func formatGetServiceError(errCode int, message string) middleware.Responder {
	errCodeVar := int32(errCode)
	return dsb_web.NewGetServiceInstanceDefault(errCode).WithPayload(
		&models.Error{Code: &errCodeVar, Message: &message},
	)
}

func getBindingInfoForInstance(k8s k8sClient.ClientInterface, mongoClient mongo.MongoInterface, instanceId string) (*models.ServiceInstanceDetails, error) {
	serviceName := getServiceNameFromInstanceId(instanceId)

	isReady, svc, err := k8s.TestService(serviceName)
	if err != nil {
		return nil, err
	}

	// Still creating
	if !isReady {
		return &models.ServiceInstanceDetails{
			InstanceID: instanceId,
			State:      "CREATING",
		}, nil

	}

	var port int32 = 27017
	var server string
	bindingInfo := make(map[string]string)
	bindingInfo["port"] = fmt.Sprintf("%d", port)
	if svc.Spec.Type == v1.ServiceTypeClusterIP {
		server = svc.Spec.ClusterIP
		bindingInfo["server"] = server
	} else {
		return nil, fmt.Errorf("Unsupported k8s service type %s for service %s", svc.Spec.Type, serviceName)
	}

	// Verify mongo is alive
	mongoAddr := fmt.Sprintf("%s:%d", server, port)
	fmt.Printf("TESTING for mongo to start accepting connections %s", mongoAddr)
	err = mongoClient.VerifyMongoLive(mongoAddr)
	if err != nil {
		log.Printf("mongo not live yet %s on %s - %s", serviceName, mongoAddr, err.Error())
		// Assuming still creating
		return &models.ServiceInstanceDetails{
			InstanceID: instanceId,
			State:      "CREATING",
		}, nil
	}

	p := "tcp"
	return &models.ServiceInstanceDetails{
		InstanceID: instanceId,
		State:      "RUNNING",
		Binding:    bindingInfo,
		BindingPorts: []*models.BindingPort{
			{
				Protocol:    &p,
				Destination: &server,
				Port:        &port},
		},
		Size:        500,
		StorageType: "Kubernetes Temp Volume",
	}, nil

}
