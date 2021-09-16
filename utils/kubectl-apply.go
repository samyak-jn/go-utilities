package utils

import (
	"fmt"
	"log"
	"os/exec"
)

func Kubectl_apply() {
	prg := "kubectl apply -n %s -f %s"
	namespace := "argocd"
	URL := "https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml"
	s := fmt.Sprintf(prg, namespace, URL)

	// To print which command is running
	fmt.Println("Running: ", s)

	// To get the output of the command
	out, err := exec.Command("kubectl", "apply", "-n", namespace, "-f", URL).Output()
	if err != nil {
		log.Fatal(err)
	}

	// To actually run the command (runs in background)
	cmd_run := exec.Command("kubectl", "apply", "-n", namespace, "-f", URL)
	err_run := cmd_run.Run()

	if err_run != nil {
		log.Fatal(err_run)
	}
	fmt.Println(string(out), "Command Run Successful!")
}
