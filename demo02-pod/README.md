### Install the `oc` CLI

Follow the related OpenShift [documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/4.16/html/cli_tools/openshift-cli-oc#cli-about-cli_cli-developer-commands). 

### Login into the cluster

Several options available. Follow the related OpenShift [documentation](https://docs.redhat.com/en/documentation/openshift_container_platform/4.16/html/installing/installing-on-ibm-power#cli-logging-in-kubeadmin_installing-ibm-power). 

### Verify login to the cluster is successful

```
oc whoami --show-context
```


### Create a new project on the OpenShift cluster

```
oc new-project demo-project
```

### Apply the _yaml_ file to create a Pod

```
cd demo02-pod
oc apply -f pod.yaml
```

### Check the Pod


```
oc get pod demo-pod
oc get pod demo-pod -o wide
oc get pod demo-pod -o yaml
```

### Delete the Pod

```
oc delete pod demo-pod
```

### Recreate the Pod

```
oc apply -f pod.yaml
```

### Expose the Pod to the Internet

```
cat service.yaml
oc apply -f service.yaml
oc get service demo-service -o yaml
```
```
oc expose service demo-service
oc get route demo-service -o yaml
curl http://demo-service-demo-project.apps.parma.edu.ihost.com/
```