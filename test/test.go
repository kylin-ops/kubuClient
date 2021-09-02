package main

import (
	"fmt"
	kubeClient "github.com/kylin-ops/kubuClient"
)

func main() {
	client, err := kubeClient.NewClient("config")
	if err != nil {
		panic(err)
	}
	deployment, _ := client.DeploymentGet("default", "nginx")
	fmt.Println(deployment.Name, deployment.Status, deployment.ClusterName)
	fmt.Println(deployment.Namespace)
}
