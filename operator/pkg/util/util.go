package util

import "fmt"

func EndpointURL(name, namespace string, port int, path string) string {
	// "http://"+foundDataset.Name+"."+foundDataset.Namespace+".svc.cluster.local"+":"+strconv.Itoa(int(foundDataset.Spec.Port))+foundDataset.Spec.Path)
	return fmt.Sprintf("http://%s.%s.svc:%d%s", name, namespace, port, path)
}
