# PREPARATION

## BUILD MICROSERVICES IMAGES TO PPC64LE

- Download repo

`git clone https://github.com/go-micro/demo.git`

- Work around dockerhub download quota, by using an alternative public registry

`sed -i 's%FROM golang:1.19.0-alpine%FROM public.ecr.aws/docker/library/golang:1.19.0-alpine%g' demo/src/cartservice/Dockerfile` 

- Fix error _Liveness probe failed: exec /bin/grpc_health_probe: exec format error_

`sed -i 's%grpc_health_probe-linux-amd64%grpc_health_probe-linux-ppc64le%g' demo/src/cartservice/Dockerfile` 

- **On a PPC64LE host**, build the image

`podman build -t quay.io/acancell-redhat-talks/ibm_txc_24_rome/onlineboutique-micro-ppc64le/cartservice:latest demo/src/cartservice`

- Login into _quay.io_

`podman login quay.io`

- Push the image

`podman push quay.io/acancell-redhat-talks/ibm_txc_24_rome/onlineboutique-micro-ppc64le/cartservice:latest`

- Connect to _quay.io_ GUI and make the repo public

https://quay.io/acancell-redhat-talks/ibm_txc_24_rome/onlineboutique-micro-ppc64le/cartservice > Setting > Repository Visibility > Public

- **Repeat for each service in `demo/src/`**

## BUILD VALKEY TO PPC64LE

Valkey (https://github.com/valkey-io/valkey) is an open source fork of Redis. Its offical container code is at: https://github.com/valkey-io/valkey-container

**NOTE** All following steps are to be made on a PPC64LE host

- Download repo

`git clone https://github.com/valkey-io/valkey-container.git`

- Work around dockerhub download quota, by using an alternative public registry

`sed -i 's%FROM alpine:{{ .alpine.version }}%FROM public.ecr.aws/docker/library/alpine:{{ .alpine.version }}%g' valkey-container/Dockerfile.template`

- Install the prerequisite tool `bashbrew` for the script `update.sh`

```
mkdir ~/.local/bin
wget https://doi-janky.infosiftr.net/job/bashbrew/job/master/lastSuccessfulBuild/artifact/bashbrew-ppc64le
mv bashbrew-ppc64le bashbrew
chmod +x bashbrew
mv bashbrew ~/.local/bin/
```

- Generate the final Dockefile

`valkey-container/update.sh`

- Build the image

`podman build -t quay.io/acancell-redhat-talks/ibm_txc_24_rome/onlineboutique-micro-ppc64le/valkey:latest valkey-container/7.2/alpine`

- Login into _quay.io_

`podman login quay.io`

- Push the image

`podman push quay.io/acancell-redhat-talks/ibm_txc_24_rome/onlineboutique-micro-ppc64le/valkey:latest`

- Connect to _quay.io_ GUI and make the repo public:

https://quay.io/acancell-redhat-talks/ibm_txc_24_rome/onlineboutique-micro-ppc64le/valkey > Setting > Repository Visibility > Public

## PREPARE AUTHN ON OCP

- Add a _ServiceAccoun_ and related authz

```
$ cat << EOF | oc apply -f -
> apiVersion: v1
kind: ServiceAccount
metadata:
  name: go-micro
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: go-micro-registry
rules:
  - apiGroups:
  	- ""
	resources:
  	- pods
	verbs:
  	- list
  	- patch
  	- watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: go-micro-registry
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: go-micro-registry
subjects:
  - kind: ServiceAccount
	name: go-micro
EOF
```

It will create:

```
serviceaccount/go-micro created
clusterrole.rbac.authorization.k8s.io/go-micro-registry created
rolebinding.rbac.authorization.k8s.io/go-micro-registry created
```

- To solve the error `unable to validate against any security context constraint: [provider "anyuid": Forbidden: not usable by user or serviceaccount ...` run the following:

`oc adm policy add-scc-to-user anyuid -z go-micro -n <project>`