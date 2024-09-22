## CREDITS

This project was adapted from _https://github.com/go-micro/demo_ (which is itself based on the upstream _https://github.com/GoogleCloudPlatform/microservices-demo_).

A copy of the original files can be also found here in the folder _onlineboutique-micro_ORIGINAL_

## PORT TO _PPC64LE_

All container images were adapted and rebuilt to _ppc64le_. They can be found on: _https://quay.io/repository/acancell-redhat-talks/ibm_txc_24_rome/onlineboutique-micro-ppc64le/\<service\>_

For an howto about how to adapted and rebuilt such images to _ppc64le_, see PREPARATION.md

## IMPLEMENT PERSISTENCE

The original _cartservice_ (https://github.com/go-micro/demo/tree/main/src/cartservice) requires an underlying Redis pod, however **it still saves data only on memory**.
I refactored the code, so that it actually stores data on the database (also using Valkey instead of Redis).

For further info, see folder _onlineboutique-micro_cartservice-persist_

## PORT TO _OCP_

Several changes were made to update the K8s manifests to a more recent K8s/OCP release or to simplify the manifests for an deployment onprem, for example:
- add a _ServiceAccount_ with the required permissions
- use _[Valkey](https://github.com/valkey-io/valkey)_ instead of _Redis_ and move to a more fitting _StatefulSet_ instead of a _Deployment_
- use a _Route_ instead of a _Service_ of kind _LoadBalancer_

For an howto about the prerequisites to apply on OCP, see PREPARATION.md

## DEPLOY ON _OCP_

- `cd kubernetes-manifests`
- `oc new-project onlineboutique-micro-ppc64le`
- `oc apply -f authz.yaml`
- `oc adm policy add-scc-to-user anyuid -z go-micro -n onlineboutique-micro-ppc64le`
- `for manifest in *; do oc apply -f $manifest; done`
- `oc expose service frontend`
