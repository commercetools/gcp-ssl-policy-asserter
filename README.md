This program asserts SSLPolicies for all HTTPSProxies in a
Google load balancer project.

## Configuration

| Environment Variable | Function |
| ---------- | ------- |
| SSL_POLICY_NAME  | Given a name will assert a Policy with that name exists  |
| GOOGLE_PROJECT   | Google Project _ID_  to manage HTTPSProxies for |
| GOOGLE_APPLICATION_CREDENTIALS | Path to Google Auth file. More info [here](https://cloud.google.com/docs/authentication/getting-started) |
| MIN_TLS_VERSION  | Minimum TLS version (default `TLS_1_2`). Only supports the higher than TLS 1.1 version   |
| SSL_PROFILE      | SSL Profile type (default `MODERN`). Only supports the `COMPATIBLE`/`MODERN`/`RESTRICTED`  |

| YAML Property | Function |
| --------- | --------- |
| ignoreProxies[] | If an HTTPSProxy uses a URLMap within this list the SSLPolicy will not be asserted |

## Build and Deploy

Refer to the Makefile. It has all the commands detailed and variables set.

To publish a new version you should only have to do:

`make publish`

Deployed via k8s-manifest. Chart is located [here](https://github.com/commercetools/k8s-manifests/tree/master/charts/ssl-asserter).


## IAM Permissions

```
# sslPolicy permissions
compute.sslPolicies.create
compute.sslPolicies.get
compute.sslPolicies.list
compute.sslPolicies.listAvailableFeatures
compute.sslPolicies.use

#  httpsProxies
compute.targetHttpsProxies.list
compute.targetHttpsProxies.setSslPolicy

# operations (to view long running operation status)
# SSLPolicy creation is one of these, but it creates quickly.
# Could be useful for errors though.
compute.globalOperations.get
compute.globalOperations.list

# project permissions
resourcemanager.projects.get
```
