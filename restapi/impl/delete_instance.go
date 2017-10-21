// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"github.com/go-openapi/runtime/middleware"
	"log"
	"net/http"
	"ocopea/k8s-mongodsb/models"
	"ocopea/k8s-mongodsb/restapi/operations/dsb_web"
	k8sClient "ocopea/kubernetes/client"
)

func DeleteInstanceResponse(
	k8s k8sClient.ClientInterface,
	params dsb_web.DeleteServiceInstanceParams) middleware.Responder {

	// Validating request
	instanceId := params.InstanceID
	if instanceId == "" {
		return formatDeleteServiceError(403, "missing instanceId")
	}

	appUniqueName := getServiceNameFromInstanceId(instanceId)

	// Returning 404 if none exist
	exists, _ := k8s.CheckServiceExists(appUniqueName)
	if !exists {
		return formatDeleteServiceError(http.StatusNotFound, "instance id "+instanceId+" not found")
	}

	err := deleteMongoService(k8s, appUniqueName)
	if err != nil {
		log.Printf("delete instance resulted in error %s\n", err.Error())
		return formatDeleteServiceError(500, err.Error())
	}

	return dsb_web.NewDeleteServiceInstanceOK().WithPayload(
		&models.ServiceInstance{
			InstanceID: instanceId,
		})
}

func formatDeleteServiceError(errCode int, message string) middleware.Responder {
	errCodeVar := int32(errCode)
	return dsb_web.NewDeleteServiceInstanceDefault(errCode).WithPayload(
		&models.Error{Code: &errCodeVar, Message: &message},
	)
}

func deleteMongoService(k8s k8sClient.ClientInterface, appUniqueName string) error {

	// In order to delete a service we need to delete both replication controller and service
	err := k8s.DeleteReplicationController(appUniqueName)
	if err != nil {
		return err
	}

	err = k8s.DeleteService(appUniqueName)
	if err != nil {
		return err
	}

	return nil

}
