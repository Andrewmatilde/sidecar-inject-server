# README

Sidecar-Inject is a simple server to inject containers into Pod By `Sidecar` which defined by `CRD`

It's easy to use and also offer container selector ability.

## Build Image

`Sidecar-inject` is managed by `go mod`

run `build.sh` to build docker image

## How to Install

you can install sidecar-inject by `Helm`

First install `sidecar-inejct-init` 
```yaml
helm install install/helm/sidecar-inject-init --name sidecar-init --namespace sidecar-system 
```

Then install `sidecar-inject-server`
```yaml
helm install install/helm/sidecar-inject-server --name sidecar-server --namespace sidecar-system --values install/helm/sidecar-inject-server/values.yaml 
```

## How to test

Label the `namespace` and create a `deployment` to test sidecar-inject-server

```
kubectl label namespace sidecar-system sidecar-injector=enabled
kubectl apply -f example/
```
And Here is the result
```
$ kubectl get sidecar
NAME                CREATED AT
test-sidecarset     55m
test-sidecarset-2   55m

$ kubectl get pod
NAME                                                    READY   STATUS        RESTARTS   AGE
sidecar-server-sidecar-inject-server-6b9dff7dd6-7b5nq   1/1     Running       0          39m
sleep-578649fc85-xcg2r                                  3/3     Running       0          20s
```

## Sidecar

`Sidecar` is defined by `CRD`. It shows which containers developers want to inject Into Pod When they are creating by `containers`
and it also offer the selector condition ability by `selector`

Here is one example for `Sidecar`. All Pods will be injected Containers defined in `containers` if their annotations contain `app` key and its value is `sidecar1`
```yaml
apiVersion: yisaer.github.io/v1alpha1
kind: Sidecar
metadata:
  name: test-sidecarset
spec:
  selector:
    matchLabels:
      app: sidecar1
  containers:
    - name: sidecar1
      image: centos:7
      command: ["sleep", "999d"]
    - name: sidecar2
      image: centos:7
      command: ["sleep", "999d"]
```