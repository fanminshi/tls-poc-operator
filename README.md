# tls-poc-operator

The tls-poc-operator tests the TLS utility protype defined in [tls util](https://github.com/fanminshi/operator-sdk/tree/tls_util_design/pkg/util/tlsutil).

## Overview

The tlc-poc-operator deploys a [simple-server](https://github.com/fanminshi/simple-server) which is a server that serves a "Hello World" static html page and a simple-client that retrieves the "Hello world" page from the server. The connection between server client is secured via mutal TLS. The tlc-poc-operator also creates the necessary TLS assets, service, and deployment manifests to deploy both the simple-server and the simple-client.

## Quick Start

```sh
# Download the tls-poc-operator
$ mkdir $GOPATH/src/github.com/fanminshi/
$ git clone https://github.com/fanminshi/simple-server.git
$ cd simple-server
# Setup the vendor Dependences
$ dep ensure -v
# Create the CRD and Custom Resouce.
$ kubectl create -f deploy/crd.yaml
$ kubectl create -f deploy/cr.yaml
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
simple-client   1         1         1            0           20s
simple-server   1         1         1            1           21s
# Verify that the svc is up.
# The simple-server-service is used to access the simple-server.
# The simple-client-service is a headless service for the client pod. It is used in
# Clinet Cert's SAN field for the the server to verify the identify of the client.
$ kubectl get svc
simple-client-service   ClusterIP   None             <none>        8080/TCP    21h
simple-server-service   ClusterIP   10.105.237.227   <none>        8080/TCP    21h
# Once client and server are deployed, verify that the client is able to get Hello World page from the server.
$ kubectl logs -f simple-client-586dc44756-hdprr
2018/07/31 19:19:43 <!DOCTYPE html>
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
