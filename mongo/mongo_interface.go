// Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
package mongo

import (
	"io"
	"os/exec"
)

type MongoInterface interface {
	MongoDump(mongoDBAddress string) (io.Reader, error)
	MongoRestore(mongoDbAddress string) (io.WriteCloser, *exec.Cmd, error)
	VerifyMongoLive(mongoDBAddress string) error
}
