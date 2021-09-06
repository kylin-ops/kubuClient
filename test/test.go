package main

import (
	"context"
	"fmt"
	kubeClient "github.com/kylin-ops/kubuClient"
)

func main() {
	client, err := kubeClient.NewClient("config")
	if err != nil {
		panic(err)
	}
	// req := client.PodGetLogTailLines("default", "nginx-f89759699-f6rw4", 10)
	req := client.PodGetLogFollow("default", "nginx-f89759699-f6rw4")
	reader, _ := req.Stream(context.Background())
	buf := make([]byte, 1024)
	for {
		_, err := reader.Read(buf)
		if err != nil {
			return
		}
		fmt.Println(string(buf))
	}

	//for _, item :=range items.Items {
	//	client.ResourceYaml(item)
	//	fmt.Println()
	//	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	//}
}
