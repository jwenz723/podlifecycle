# Podlifecycle

This repo contains a few applications I have used to learn about the Kubernetes pod lifecycle. The code contained within
[server](server) and [cli](cli) are discussed in my blog post [here](https://medium.com/@jwenz723/deploy-kubernetes-grpc-workloads-with-zero-down-time-3585c146f74f).

The code in [server2](server2) was used to understand the order in which the preStop hook, liveness probe, and readiness probe
are executed by Kubernetes, as well as how these APIs are impacted when a blocking request is being handled by the http server.