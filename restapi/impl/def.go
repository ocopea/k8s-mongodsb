// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"github.com/go-openapi/runtime"
	"net/http"
	"ocopea/k8s-mongodsb/models"
)

type ErrorResponse interface {
	SetPayload(payload *models.Error)
	WriteResponse(http.ResponseWriter, runtime.Producer)
}

func getServiceNameFromInstanceId(instanceId string) string {
	l := len(instanceId)
	if l > 22 {
		return "m-" + instanceId[l-22:l]
	} else {
		return "m-" + instanceId
	}
}
