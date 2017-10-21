// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"ocopea/k8s-mongodsb/models"
	"ocopea/k8s-mongodsb/mongo"
	"ocopea/kubernetes/client"
	"ocopea/kubernetes/client/v1"
	"testing"
)

// Create mongoDsb service Instance copy
func TestCopyServiceInstance(t *testing.T) {

	instanceId := "mongo1"
	kClient := &client.ClientMock{
		MockTestService: func(serviceName string) (bool, *v1.Service, error) {
			return true,
				&v1.Service{
					ObjectMeta: v1.ObjectMeta{
						Name: serviceName,
					},
					Spec: v1.ServiceSpec{
						Type:      v1.ServiceTypeClusterIP,
						ClusterIP: "10.11.12.13",
					},
				}, nil
		},
	}

	copyId := "copy1"
	err, _ := createMongoCopy(
		kClient,
		mongo.NewMongoMock(true),
		instanceId,
		&models.CopyServiceInstance{
			CopyID:   &copyId,
			CopyType: "dump",
			CopyRepoCredentials: map[string]string{
				"url":    "http://crb-data-url/repo1/copy1",
				"repoId": "repo1",
			},
		},
	)

	if err == nil {
		t.Fatal("service copy should fail due to invalid cr url")
	}

}
