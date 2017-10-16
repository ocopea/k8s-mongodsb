// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package impl

import (
	"errors"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"ocopea/k8s-mongodsb/models"
	"ocopea/k8s-mongodsb/mongo"
	"ocopea/k8s-mongodsb/restapi/operations/dsb_web"
	k8sClient "ocopea/kubernetes/client"
)

func CopyServiceInstance(
	k8s k8sClient.ClientInterface,
	mongoClient mongo.MongoInterface,
	params dsb_web.CopyServiceInstanceParams) middleware.Responder {

	err, response := createMongoCopy(k8s, mongoClient, params.InstanceID, params.CopyDetails)
	if err != nil {
		log.Printf("Failed creating copy for instance id %s - %s\n", params.InstanceID, err.Error())
		return formatCopyServiceError(500, err.Error())
	}

	log.Printf("copy created successfully for %s with %d", response.CopyID, response.Status)
	return dsb_web.NewCopyServiceInstanceOK().WithPayload(&models.CopyServiceInstanceResponse{
		CopyID:        *params.CopyDetails.CopyID,
		Status:        0,
		StatusMessage: "yey",
	})
}

func formatCopyServiceError(errCode int, message string) middleware.Responder {
	errCodeVar := int32(errCode)
	return dsb_web.NewCopyServiceInstanceDefault(errCode).WithPayload(
		&models.Error{Code: &errCodeVar, Message: &message},
	)
}

func createMongoCopy(
	k8sClient k8sClient.ClientInterface,
	mongoClient mongo.MongoInterface,
	instanceId string,
	copyDetails *models.CopyServiceInstance) (error, *models.CopyServiceInstanceResponse) {

	log.Printf("Creating copy for %s\n", instanceId)

	mongoBindings, err := getBindingInfoForInstance(k8sClient, mongoClient, instanceId)
	if err != nil {
		return fmt.Errorf("Failed getting bind details for %s: %s", instanceId, err.Error()), nil
	}

	dumpReader, err := mongoClient.MongoDump(*mongoBindings.BindingPorts[0].Destination + ":" + fmt.Sprint(*mongoBindings.BindingPorts[0].Port))
	if err != nil {
		return fmt.Errorf("Failed getting mongo dump reader for %s: %s", instanceId, err.Error()), nil
	}

	copyRepoUrl, err := getCopyRepUrl(copyDetails)
	if err != nil {
		return err, nil
	}

	postCopyUrl := fmt.Sprintf("%s/repositories/%s/copies/%s/data", copyRepoUrl, copyDetails.CopyRepoCredentials["repoId"], *copyDetails.CopyID)
	log.Printf("Posting copy to %s\n", postCopyUrl)
	req, err := http.NewRequest("POST", postCopyUrl, dumpReader)
	if err != nil {
		return err, nil
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("dsb", DSB_NAME)
	req.Header.Set("copyTimestamp", fmt.Sprintf("%d", copyDetails.CopyTime))
	req.Header.Set("facility", "mongodump")

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"Failed execute request for posting copy %s to CRB with url %s: %s", *copyDetails.CopyID, copyRepoUrl, err.Error()), nil
	}
	log.Printf("After Do for copy %s to CRB URL %s, status is %s(%d)\n", *copyDetails.CopyID, postCopyUrl, resp.Status, resp.StatusCode)

	// Validating http return code
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"Failed execute request for posting copy %s to CRB with url %s: - received status %s", *copyDetails.CopyID, copyRepoUrl, resp.Status), nil
	}

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Failed reading post copy response from crb url %s for copyId %s - %s", copyRepoUrl, copyDetails.CopyID, err.Error()), nil
	}
	log.Printf("response %s", string(all))

	return nil, &models.CopyServiceInstanceResponse{
		CopyID:        *copyDetails.CopyID,
		Status:        0,
		StatusMessage: "yey",
	}

}
func getCopyRepUrl(copyDetails *models.CopyServiceInstance) (string, error) {
	shpanRestCopyUrl := copyDetails.CopyRepoCredentials["url"]
	if shpanRestCopyUrl == "" {
		return "", errors.New("Invalid copy repo credentails format. missing \"url\"")
	} else {
		return shpanRestCopyUrl, nil
	}

}
