// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package mongo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type MongoClient struct {
}

func (m *MongoClient) MongoDump(mongoDBAddress string) (io.Reader, error) {
	fmt.Println("Running mongodump")

	mongoDBName := "" //Empty means all DBs

	cmdName := "mongodump"
	cmdArgs := []string{"-v", "--host", mongoDBAddress, "--archive"}
	if mongoDBName != "" {
		cmdArgs = append(cmdArgs, "--db=\""+mongoDBName+"\"")
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	fmt.Println(cmd)
	cmd.Stderr = os.Stdout
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	//	reader := bufio.NewReader(cmdReader)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmdReader, nil
	/*
		err = cmd.Wait()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			panic(err)
		}
	*/
}

func (m *MongoClient) MongoRestore(mongoDbAddress string) (io.WriteCloser, *exec.Cmd, error) {
	log.Printf("Running mongorestore for %s", mongoDbAddress)

	mongoDBName := "" //Empty means all DBs

	cmdName := "mongorestore"
	cmdArgs := []string{"-v", "--host", mongoDbAddress, "--archive"}
	if mongoDBName != "" {
		cmdArgs = append(cmdArgs, "--db=\""+mongoDBName+"\"")
	}

	cmd := exec.Command(cmdName, cmdArgs...)
	fmt.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmdWriter, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return cmdWriter, cmd, nil
}

func (m *MongoClient) VerifyMongoLive(mongoDBAddress string) error {
	fmt.Println("Verifying mongo live")

	cmdName := "mongoshell"
	cmdArgs := []string{"--verbose", "--host", mongoDBAddress, "--eval", "db.version()"}

	cmd := exec.Command(cmdName, cmdArgs...)
	fmt.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	return cmd.Run()
}
