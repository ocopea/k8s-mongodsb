// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package mongo

import (
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type MongoMock struct {
	live bool
}

func NewMongoMock(live bool) *MongoMock {
	return &MongoMock{live: live}
}

func (m *MongoMock) MongoDump(mongoDBAddress string) (io.Reader, error) {
	return nil, nil
}

func (m *MongoMock) MongoRestore(mongoDbAddress string) (io.WriteCloser, *exec.Cmd, error) {
	return nil, nil, nil
}

func (m *MongoMock) VerifyMongoLive(mongoDBAddress string) error {
	if m.live {
		return nil
	} else {
		return errors.New("non responsive")
	}
}
