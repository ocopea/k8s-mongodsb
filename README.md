# Ocopea Kubernetes mongodsb

**Ocopea MongoDB Data Service Broker implementation for running MongoDB on Kubernetes**

For more information about Ocopea for Kubernetes, visit the 
Ocopea [Kubernetes repo](https://github.com/ocopea/kubernetes).

This DSB supports a single plan named "standard" that runs the Mongo Docker image on Kubernetes and supports 
taking a service copy using `mongodump` and `mongorestore`.


# How to build

* clone the repo under GOPATH/src/ocopea/
* run the `buildImage.sh` script

# How to run

The easiest and the recommended way to run the mongodsb is by using the Ocopea 
[Kubernetes deployer](https://github.com/ocopea/kubernetes)
using the `deploy-mongodsb` command.
for example for running the DSB on minikube into a namespace called "testing":

```
$ go run deployer.go deploy-mongodsb -namespace=testing -local-cluster-ip=$(minikube ip)
```

If you wish to run the container yourself using Kubernetes deployment all you have to do is make sure you supply 
Kubernetes credentials if required using both environment variables: `K8S_USERNAME, K8S_PASSWORD`.

Once running you can attach the DSB to an Ocopea site using the site config screen in the UI, or directly via the Ocopea API.

# Tests

The project uses the go-swagger server code generator when building the Docker image.
In order to build the project, debug or run the tests locally you'll have to run 
the swagger tool locally. The swagger generator version we use is 0.10.0. 
Find installation instructions [here](https://goswagger.io/).

```
$ swagger generate server -f dsb-swagger.yaml -A k8s-mongodsb
```

To run the unit test:

```
$ cd restapi/impl
$ go test
```

To run end to end tests on a Kubernetes cluster use the 
[Kubernetes tester project](https://github.com/ocopea/kubernetes/tree/master/tests).


## Contribution

* [Contributing to Ocopea](https://github.com/ocopea/documentation/blob/master/docs/contributing.md)
* [Ocopea Developer Guidelines](https://github.com/ocopea/documentation/blob/master/docs/guidelines.md)

