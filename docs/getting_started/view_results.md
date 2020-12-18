# Evaluate the results

This guide illustrates how to view the results of the benchmark.

## Viewing the Kibana dashboards

In order to view the dashboards you should have installed the Logging
subsystem of KOBE.

After all pods are in Running state Kibana dashboards can be accessed at 
```
http://<NODE-IP>:<NODEPORT>/app/kibana#/dashboard/
``` 
where `<NODE-IP>` the IP of any of the Kubernetes cluster nodes and `<NODEPORT>`
the result of `kubectl get -o jsonpath="{.spec.ports[0].nodePort}" services
kibana-kibana`.


