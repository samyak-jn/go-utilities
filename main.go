package main

import (
	utils "go-utilities/utils"
)

var (
	kubeconfig = "/home/samjain/.kube/config"
	master     = ""
	namespace  = "argocd"
	podName    = "argocd-server-86dcc9f88f-b64d6"
)

func main() {
	utils.Get_pod_info(kubeconfig, master, namespace, podName)
	// utils.Kubectl_apply()var jsonMap map[string]interface{}

}
