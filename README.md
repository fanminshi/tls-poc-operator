# tls-poc-operator

The tls-poc-operator tests the TLS utility protype defined in [tls util](https://github.com/fanminshi/operator-sdk/tree/tls_util_design/pkg/util/tlsutil).

## Overview

The tlc-poc-operator deploys a [simple-server](https://github.com/fanminshi/simple-server) which is a server that serves a "Hello World" static html page via TLS as Kubernetes pod. The tlc-poc-operator creates the necessary TLS assets, service, and deployment manifests to deploy the simple-server.

## Quick Start

```sh
# Download the tls-poc-operator
$ mkdir $GOPATH/src/github.com/fanminshi/
$ git clone https://github.com/fanminshi/simple-server.git
$ cd simple-server
# Setup the vendor Dependences
$ dep ensure -v
# Create the CRD
$ kubectl create -f deploy/crd.yaml
# Run the operator locally
$ OPERATOR_NAME=app-operator operator-sdk up local --namespace=default
INFO[0000] Go Version: go1.10
INFO[0000] Go OS/Arch: darwin/amd64
INFO[0000] operator-sdk Version: 0.0.5+git
INFO[0000] Metrics service app-operator created
INFO[0000] Watching security.example.com/v1alpha1, Security, default, 5
# Verify that the deployment is ready
$ kubectl get deploy
NAME            DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
simple-server   1         1         1            1           19m
# Verify that the svc is up. The simple-server-service is used to access the simple-server.
$ kubectl get svc
simple-server-service   ClusterIP   10.96.91.119    <none>        443/TCP     20m
# Deploy a busy box with curl command in order to access simple-server.
$ kubectl run curl --image=radial/busyboxplus:curl -i --tty
# Once gain access, access the secured simple-server-service using curl.
[ root@curl-545bbf5f9c-f7xtj:/ ]$ curl -k https://simple-server-service
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>hello world</title>
  </head>
  <body>
    <h1>hello world</h1>
  </body>
</html>
```
