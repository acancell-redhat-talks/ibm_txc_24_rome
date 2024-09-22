## NOTE

The manifests were modified by their original in https://github.com/go-micro/demo/tree/main/kubernetes-manifests.

In particular:
- modified `image` so that it uses _ppc64le_ ones
- set `TRACING_ENABLE` to _false_
- used _[Valkey](https://github.com/valkey-io/valkey)_ instead of _Redis_ and moved to a more fitting _StatefulSet_ instead of a _Deployment_
- used a _Route_ instead of a _Service_ of kind _LoadBalancer_