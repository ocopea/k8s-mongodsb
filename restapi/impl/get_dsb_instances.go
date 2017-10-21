// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"github.com/go-openapi/runtime/middleware"
	"ocopea/k8s-mongodsb/models"
	"ocopea/k8s-mongodsb/restapi/operations/dsb_web"
)

func DsbInstancesResponse() middleware.Responder {

	return dsb_web.NewGetServiceInstancesOK().WithPayload(
		[]*models.ServiceInstance{})
}
