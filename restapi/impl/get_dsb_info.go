// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"github.com/go-openapi/runtime/middleware"
	"ocopea/k8s-mongodsb/models"
	"ocopea/k8s-mongodsb/restapi/operations/dsb_web"
)

const DSB_NAME = "mongo-k8s-dsb"

func DsbInfoResponse() middleware.Responder {
	return dsb_web.NewGetDSBInfoOK().WithPayload(getMongoDsbInfo())
}

func getMongoDsbInfo() *models.DsbInfo {
	return &models.DsbInfo{
		Name:        DSB_NAME,
		Description: "mongodb DSB for kubernetes",
		Type:        "datasoruce",
		Plans: []*models.DsbPlan{
			{
				ID:          "standard",
				Name:        "standard",
				Description: "mongo container",
				DsbSettings: nil,
				CopyProtocols: []*models.DsbSupportedCopyProtocol{
					{
						CopyProtocol:        "ShpanRest",
						CopyProtocolVersion: "1.0",
					},
				},
				Protocols: []*models.DsbSupportedProtocol{
					{
						Protocol: "mongodb",
					},
				},
			},
		},
	}
}
