# k8s-events-logger [![Docker](https://github.com/lescactus/k8s-events-logger/actions/workflows/docker.yml/badge.svg)](https://github.com/lescactus/k8s-events-logger/actions/workflows/docker.yml) [![Go](https://github.com/lescactus/k8s-events-logger/actions/workflows/go.yml/badge.svg)](https://github.com/lescactus/k8s-events-logger/actions/workflows/go.yml)

This repository contains a simple Kubernetes controller displaying namespaces events to stdout.

## Motivations

Kubernetes `Events` are objects showing what is happening inside a cluster, node, pod, or container. These objects are usually generated in response to changes that occur inside Kubernetes. The Kubernetes API Server enables all core components to create these events.
By default, `Events` are only retained one hour. The `--event-ttl` api server flag allow to change this value. However, increasing this value put pressure on the `etcd` cluster, hence it isn't recommended put a high value.

Since `Events` allow Ops teams and developers to help troubleshooting workloads when issues arise, it is interesting to keep these events messages for a longer time than the default TTL.

Enters `k8s-events-logger`: `k8s-events-logger` is meant to be running inside a Kubernetes cluster and  will read cluster events to display them on `stdout`. It can then be possible to aggregate the events logs into a log aggregation system (such as [ELK](https://www.elastic.co/what-is/elk-stack) for instance) for later retrival.

## Configuration

`k8s-events-logger` is a 12-factor app using [Viper](https://github.com/spf13/viper) as a configuration manager. It can read configuration from environment variables or from .env files.

### Available variables

* `OUTPUT`(default value: `console`). Define the output format of event logs. Available values are `console` or `json`.

* `NAMESPACES` (default value: `default`). Comma separated list of namespace(s) to watch events from. Examples: (`ns1,ns2`).

* `KUBECONFIG` (default value: ""). Path of a kube config. If empty, `k8s-events-logger` will create a in-cluster Kubernetes client and use the service account token to authenticate against the Kubernetes API.

## Building

### From source with go

You need a working [go](https://golang.org/doc/install) toolchain (It has been developped and tested with go 1.17 and go 1.18 and should work with go >= 1.17). Refer to the official documentation for more information (or from your Linux/Mac/Windows distribution documentation to install it from your favorite package manager).

```sh
# Clone this repository
git clone https://github.com/lescactus/k8s-events-logger.git && cd k8s-events-logger/

# Build from sources. Use the '-o' flag to change the compiled binary name
go build

# Default compiled binary is k8s-events-logger
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./k8s-events-logger
```

### From source with docker

If you don't have [go](https://golang.org/) installed but have docker, run the following command to build inside a docker container:

```sh
# Build from sources inside a docker container. Use the '-o' flag to change the compiled binary name
# Warning: the compiled binary belongs to root:root
docker run --rm -it -v "$PWD":/app -w /app golang:1.17 go build

# Default compiled binary is dict-go
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./k8s-events-logger
```

### From source with docker but built inside a docker image

If you don't want to pollute your computer with another program, this cli comes with its own docker image:

```sh
docker build -t k8s-events-logger .
```

## Installation

`k8s-events-logger` comes with its Kubernetes manifests. They are located in `deploy/k8s/`.

```
kubectl apply -f deploy/k8s/
```

Alternatively, you can use [skaffold](https://skaffold.dev/) to build and deploy in Kubernetes:

```
$ skaffold run
Generating tags...
 - k8s-events-logger -> k8s-events-logger:2022-06-28_13-31-43.916_CEST
Checking cache...
ERRO[0001] gcloud binary not found                      
 - k8s-events-logger: Found. Tagging
Starting test...
Tags used in deployment:
 - k8s-events-logger -> k8s-events-logger:fc2dd619f1ed0208172fb9cbca6b8fee158b1a59d211548d6f4499ba493b0ca3
Starting deploy...
 - deployment.apps/k8s-events-logger configured
 - clusterrolebinding.rbac.authorization.k8s.io/k8s-events-logger unchanged
 - serviceaccount/k8s-events-logger unchanged
Waiting for deployments to stabilize...
 - deployment/k8s-events-logger is ready.
Deployments stabilized in 3.091 seconds
You can also run [skaffold run --tail] to get the logs

```

## Examples

### Basic usage:
```
k8s-events-logger
2022/06/27 16:30:47 Starting k8s-events-logger. Output = console, Namespaces to watch [default]
W0627 16:30:48.387846       1 shared_informer.go:401] The sharedIndexInformer has started, run more than once is not allowed
Timestamp: 2022-06-27T15:50:56Z | Namespace: default | Type: Normal | Reason: Pulled | Object: Pod/redis-84df789599-6g8kt | Message: Container image "redis:6" already present on machine
Timestamp: 2022-06-27T15:50:54Z | Namespace: default | Type: Normal | Reason: Started | Object: Pod/redis-84df789599-rvbsx | Message: Started container redis
Timestamp: 2022-06-27T15:49:54Z | Namespace: default | Type: Normal | Reason: Created | Object: Pod/k8s-events-logger-5747f9b955-j59vr | Message: Created container k8s-events-logger
Timestamp: 2022-06-27T15:50:55Z | Namespace: default | Type: Normal | Reason: SuccessfulDelete | Object: ReplicaSet/k8s-events-logger-5965bf64db | Message: Deleted pod: k8s-events-logger-5965bf64db-w8lkv
Timestamp: 2022-06-27T15:48:11Z | Namespace: default | Type: Normal | Reason: Started | Object: Pod/k8s-events-logger-59c9f9c64-nvcsg | Message: Started container k8s-events-logger
Timestamp: 2022-06-27T16:30:44Z | Namespace: default | Type: Normal | Reason: Killing | Object: Pod/k8s-events-logger-5f86cd57d6-wx9zq | Message: Stopping container k8s-events-logger
Timestamp: 2022-06-27T16:30:44Z | Namespace: default | Type: Normal | Reason: SuccessfulDelete | Object: ReplicaSet/k8s-events-logger-5f86cd57d6 | Message: Deleted pod: k8s-events-logger-5f86cd57d6-wx9zq
Timestamp: 2022-06-27T16:30:47Z | Namespace: default | Type: Normal | Reason: Pulled | Object: Pod/k8s-events-logger-cbb5578d5-jrphs | Message: Container image "k8s-events-logger:fc2dd619f1ed0208172fb9cbca6b8fee158b1a59d211548d6f4499ba493b0ca3" already present on machine
...
```

### Json output
```
OUTPUT=json k8s-events-logger
2022/06/27 16:32:26 Starting k8s-events-logger. Output = json, Namespaces to watch [default]
W0627 16:32:26.684153       1 shared_informer.go:401] The sharedIndexInformer has started, run more than once is not allowed
{"metadata":{"name":"k8s-events-logger-5747f9b955-j59vr.16fc850593912325","namespace":"default","uid":"0672b400-0cdb-4474-9d96-19186071f977","resourceVersion":"10891","creationTimestamp":"2022-06-27T15:50:21Z","managedFields":[{"manager":"kubelet","operation":"Update","apiVersion":"v1","time":"2022-06-27T15:50:21Z","fieldsType":"FieldsV1","fieldsV1":{"f:count":{},"f:firstTimestamp":{},"f:involvedObject":{"f:apiVersion":{},"f:fieldPath":{},"f:kind":{},"f:name":{},"f:namespace":{},"f:resourceVersion":{},"f:uid":{}},"f:lastTimestamp":{},"f:message":{},"f:reason":{},"f:source":{"f:component":{},"f:host":{}},"f:type":{}}}]},"involvedObject":{"kind":"Pod","namespace":"default","name":"k8s-events-logger-5747f9b955-j59vr","uid":"9f22145c-401b-4491-9e38-e388cd69aca4","apiVersion":"v1","resourceVersion":"10820","fieldPath":"spec.containers{k8s-events-logger}"},"reason":"Killing","message":"Stopping container k8s-events-logger","source":{"component":"kubelet","host":"minikube"},"firstTimestamp":"2022-06-27T15:50:21Z","lastTimestamp":"2022-06-27T15:50:21Z","count":1,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}
{"metadata":{"name":"k8s-events-logger-5f86cd57d6.16fc8739a56de967","namespace":"default","uid":"a9c6494a-e388-44b3-900a-6188165297c9","resourceVersion":"12769","creationTimestamp":"2022-06-27T16:30:44Z","managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"v1","time":"2022-06-27T16:30:44Z","fieldsType":"FieldsV1","fieldsV1":{"f:count":{},"f:firstTimestamp":{},"f:involvedObject":{"f:apiVersion":{},"f:kind":{},"f:name":{},"f:namespace":{},"f:resourceVersion":{},"f:uid":{}},"f:lastTimestamp":{},"f:message":{},"f:reason":{},"f:source":{"f:component":{}},"f:type":{}}}]},"involvedObject":{"kind":"ReplicaSet","namespace":"default","name":"k8s-events-logger-5f86cd57d6","uid":"f9124107-fc5c-4b22-b14a-cf6058091819","apiVersion":"apps/v1","resourceVersion":"12765"},"reason":"SuccessfulDelete","message":"Deleted pod: k8s-events-logger-5f86cd57d6-wx9zq","source":{"component":"replicaset-controller"},"firstTimestamp":"2022-06-27T16:30:44Z","lastTimestamp":"2022-06-27T16:30:44Z","count":1,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}
{"metadata":{"name":"k8s-events-logger.16fc84e786ecdaa8","namespace":"default","uid":"2aa6f2a9-5d77-467c-8c81-1f20644f0997","resourceVersion":"10732","creationTimestamp":"2022-06-27T15:48:12Z","managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"v1","time":"2022-06-27T15:48:12Z","fieldsType":"FieldsV1","fieldsV1":{"f:count":{},"f:firstTimestamp":{},"f:involvedObject":{"f:apiVersion":{},"f:kind":{},"f:name":{},"f:namespace":{},"f:resourceVersion":{},"f:uid":{}},"f:lastTimestamp":{},"f:message":{},"f:reason":{},"f:source":{"f:component":{}},"f:type":{}}}]},"involvedObject":{"kind":"Deployment","namespace":"default","name":"k8s-events-logger","uid":"c505ed4b-19b5-49dd-bf33-8a0802117950","apiVersion":"apps/v1","resourceVersion":"10724"},"reason":"ScalingReplicaSet","message":"Scaled down replica set k8s-events-logger-689775574 to 0","source":{"component":"deployment-controller"},"firstTimestamp":"2022-06-27T15:48:12Z","lastTimestamp":"2022-06-27T15:48:12Z","count":1,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}
{"metadata":{"name":"redis-84df789599-rvbsx.16fc850d104f76f5","namespace":"default","uid":"891448a6-c6d1-46bd-b5bd-e36397f87632","resourceVersion":"10945","creationTimestamp":"2022-06-27T15:50:54Z","managedFields":[{"manager":"kubelet","operation":"Update","apiVersion":"v1","time":"2022-06-27T15:50:54Z","fieldsType":"FieldsV1","fieldsV1":{"f:count":{},"f:firstTimestamp":{},"f:involvedObject":{"f:apiVersion":{},"f:fieldPath":{},"f:kind":{},"f:name":{},"f:namespace":{},"f:resourceVersion":{},"f:uid":{}},"f:lastTimestamp":{},"f:message":{},"f:reason":{},"f:source":{"f:component":{},"f:host":{}},"f:type":{}}}]},"involvedObject":{"kind":"Pod","namespace":"default","name":"redis-84df789599-rvbsx","uid":"6a18ef3c-f662-448d-86ce-ae7616cdfa19","apiVersion":"v1","resourceVersion":"10928","fieldPath":"spec.containers{redis}"},"reason":"Pulled","message":"Container image \"redis:6\" already present on machine","source":{"component":"kubelet","host":"minikube"},"firstTimestamp":"2022-06-27T15:50:54Z","lastTimestamp":"2022-06-27T15:50:54Z","count":1,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}
{"metadata":{"name":"k8s-events-logger-5f86cd57d6.16fc850cdde0e1ce","namespace":"default","uid":"37353e51-9281-4c66-abce-6c6cf95fb88c","resourceVersion":"10929","creationTimestamp":"2022-06-27T15:50:53Z","managedFields":[{"manager":"kube-controller-manager","operation":"Update","apiVersion":"v1","time":"2022-06-27T15:50:53Z","fieldsType":"FieldsV1","fieldsV1":{"f:count":{},"f:firstTimestamp":{},"f:involvedObject":{"f:apiVersion":{},"f:kind":{},"f:name":{},"f:namespace":{},"f:resourceVersion":{},"f:uid":{}},"f:lastTimestamp":{},"f:message":{},"f:reason":{},"f:source":{"f:component":{}},"f:type":{}}}]},"involvedObject":{"kind":"ReplicaSet","namespace":"default","name":"k8s-events-logger-5f86cd57d6","uid":"f9124107-fc5c-4b22-b14a-cf6058091819","apiVersion":"apps/v1","resourceVersion":"10921"},"reason":"SuccessfulCreate","message":"Created pod: k8s-events-logger-5f86cd57d6-wx9zq","source":{"component":"replicaset-controller"},"firstTimestamp":"2022-06-27T15:50:53Z","lastTimestamp":"2022-06-27T15:50:53Z","count":1,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}
...
```

### Watch several namespaces events
```
NAMESPACES=default,ns1 k8s-events-logger
2022/06/27 16:33:45 Starting k8s-events-logger. Output = console, Namespaces to watch [default ns1]
W0627 16:33:45.883314       1 shared_informer.go:401] The sharedIndexInformer has started, run more than once is not allowed
Timestamp: 2022-06-27T15:50:21Z | Namespace: default | Type: Normal | Reason: Started | Object: Pod/k8s-events-logger-5965bf64db-w8lkv | Message: Started container k8s-events-logger
Timestamp: 2022-06-27T15:48:10Z | Namespace: default | Type: Normal | Reason: SuccessfulCreate | Object: ReplicaSet/k8s-events-logger-59c9f9c64 | Message: Created pod: k8s-events-logger-59c9f9c64-nvcsg
Timestamp: 2022-06-27T16:30:48Z | Namespace: default | Type: Normal | Reason: Killing | Object: Pod/k8s-events-logger-8688dc6587-pld4g | Message: Stopping container k8s-events-logger
Timestamp: 2022-06-27T15:50:55Z | Namespace: default | Type: Normal | Reason: Created | Object: Pod/redis-84df789599-q8g88 | Message: Created container redis
Timestamp: 2022-06-27T15:50:54Z | Namespace: default | Type: Normal | Reason: Created | Object: Pod/redis-84df789599-rvbsx | Message: Created container redis
Timestamp: 2022-06-27T16:33:25Z | Namespace: default | Type: Normal | Reason: Scheduled | Object: Pod/k8s-events-logger-565b7b9744-zzhgv | Message: Successfully assigned default/k8s-events-logger-565b7b9744-zzhgv to minikube
Timestamp: 2022-06-27T16:33:25Z | Namespace: default | Type: Normal | Reason: SuccessfulCreate | Object: ReplicaSet/k8s-events-logger-565b7b9744 | Message: Created pod: k8s-events-logger-565b7b9744-zzhgv
Timestamp: 2022-06-27T15:49:54Z | Namespace: default | Type: Normal | Reason: Started | Object: Pod/k8s-events-logger-5747f9b955-j59vr | Message: Started container k8s-events-logger
Timestamp: 2022-06-27T15:50:54Z | Namespace: default | Type: Normal | Reason: ScalingReplicaSet | Object: Deployment/redis | Message: Scaled down replica set redis-575c9b865f to 2
Timestamp: 2022-06-27T16:33:37Z | Namespace: ns1 | Type: Normal | Reason: Started | Object: Pod/redis-575c9b865f-zbgwx | Message: Started container redis
Timestamp: 2022-06-27T16:33:36Z | Namespace: ns1 | Type: Normal | Reason: ScalingReplicaSet | Object: Deployment/redis | Message: Scaled up replica set redis-575c9b865f to 1
```