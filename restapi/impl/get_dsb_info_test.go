// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"ocopea/k8s-mongodsb/models"
)

import (
	"reflect"
	"testing"
)

var expectedDsbInfo = &models.DsbInfo{
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

// DSB info call matches expectation
func TestGetDsbInfo(t *testing.T) {
	dsbInfoResult := getMongoDsbInfo()

	if !reflect.DeepEqual(dsbInfoResult, expectedDsbInfo) {
		t.Errorf("invalid dsb Info returned: %v, want %v", dsbInfoResult, expectedDsbInfo)
	}
}
