# Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
mvn clean install
docker build -t ocopea/mongo-k8s-dsb -f Dockerfile .
