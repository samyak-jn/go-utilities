package main

import (
	"fmt"
	utils "go-utilities/utils"
)

var (
	kubeconfig = "/home/samjain/.kube/config"
	master     = ""
	namespace  = "argocd"
	podName    = ""
)

func main() {
	utils.Get_pod_info(kubeconfig, master, namespace, podName)
	// utils.Kubectl_apply()var jsonMap map[string]interface{}
	fmt.Print(utils.PodRestartcount("argocd"))
}
