// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"errors"
	"ocopea/k8s-mongodsb/models"
	"ocopea/k8s-mongodsb/mongo"
	"ocopea/kubernetes/client"
	"ocopea/kubernetes/client/v1"
	"testing"
)

// Create mongoDsb Instance
func TestCreateDsbInstance(t *testing.T) {

	instanceId := "mongo1"
	kClient := &client.ClientMock{
		MockCreateReplicationController: func(rc *v1.ReplicationController, force bool) (*v1.ReplicationController, error) {

			// Validating replication controller validity
			if rc.Name != "m-mongo1" {
				t.Fatalf("Invalid k8s entity name %s\n", rc.Name)
			}

			if len(rc.Spec.Template.Spec.Containers) != 1 {
				t.Fatal("Containers count must be 1")
			}

			if rc.Spec.Template.Spec.Containers[0].Image != "mongo" {
				t.Fatalf("Invalid container image used  %s\n", rc.Spec.Template.Spec.Containers[0].Image)
			}

			return rc, nil
		},
		MockCreateService: func(svc *v1.Service, force bool) (*v1.Service, error) {
			// Validating service validity
			if svc.Name != "m-mongo1" {
				t.Fatalf("Invalid k8s entity name %s", svc.Name)
			}

			return svc, nil
		},
	}

	err := createMongoService(
		kClient,
		mongo.NewMongoMock(true),
		&models.CreateServiceInstance{
			InstanceID: &instanceId,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateDsbInstanceFailsDueToK8SFailureCreatingReplicationController(t *testing.T) {

	kClient := &client.ClientMock{
		MockCreateReplicationController: func(rc *v1.ReplicationController, force bool) (*v1.ReplicationController, error) {
			return rc, errors.New("bye")
		},
		MockCreateService: func(svc *v1.Service, force bool) (*v1.Service, error) {
			return svc, nil
		},
	}

	instanceId := "mongo1"
	err := createMongoService(
		kClient,
		mongo.NewMongoMock(true),
		&models.CreateServiceInstance{
			InstanceID: &instanceId,
		},
	)

	if err == nil {
		t.Fatal("Create Mongo service should have failed with an error")
	}

}

func TestCreateDsbInstanceFailsDueToK8SFailureCreatingService(t *testing.T) {

	kClient := &client.ClientMock{
		MockCreateReplicationController: func(rc *v1.ReplicationController, force bool) (*v1.ReplicationController, error) {
			return rc, nil
		},
		MockCreateService: func(svc *v1.Service, force bool) (*v1.Service, error) {
			return svc, errors.New("bye")
		},
	}

	instanceId := "mongo1"
	err := createMongoService(
		kClient,
		mongo.NewMongoMock(true),
		&models.CreateServiceInstance{
			InstanceID: &instanceId,
		},
	)

	if err == nil {
		t.Fatal("Create Mongo service should have failed with an error")
	}

}
