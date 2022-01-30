package main

import (
	utils "go-utilities/utils"
	"log"
)

var (
	kubeconfig = "/home/samjain/.kube/config"
	master     = ""
	namespace  = "argocd"
	podName    = ""
)

func main() {

	containerId, networkName, err := utils.NewEphemeralCreateTestFramework()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println([]string{containerId, networkName})
	errClean := utils.NewEphemeralCleanTestFramework(containerId, networkName)
	if errClean != nil {
		log.Fatal(errClean)
	}

}
