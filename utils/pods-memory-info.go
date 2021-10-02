package utils

import (
	// utils "github.com/redhat-appstudio/managed-gitops/utilities/load-test/loadtest"
	// metrics "k8s.io/apimachinery/pkg/apis/metrics/v1"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func PodRestartcount(namespace string) [][]string {
	var allrestartInfo [][]string

	out, err := exec.Command("kubectl", "get", "pods", "-n", namespace).Output()
	if err != nil {
		log.Fatal(err)
	}

	res := string(out)

	for index, i := range strings.Split(res, "\n") {
		var restartInfo []string

		if index != 0 && index < len(strings.Split(res, "\n"))-1 {
			restartInfo = append(restartInfo, strings.Fields(i)[0], strings.Fields(i)[3])
			allrestartInfo = append(allrestartInfo, restartInfo)
		}
	}

	return allrestartInfo
}

func Get_pod_info(kubeconfig string, master string, namespace string, podName string) {
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		panic(err)
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	podMetrics, err := mc.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// podMetrics, err := mc.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// To get memory info of the specific pod passed as an argument
	// podMetrics_, err := mc.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%+v", podMetrics.Items)

	for _, elements := range podMetrics.Items {
		fmt.Println(elements.Containers[0].Name, elements.Containers[0].Usage["cpu"].ToUnstructured(), elements.Containers[0].Usage["memory"].ToUnstructured())
	}
	fmt.Println("----------------------")

	// fmt.Println(strings.Split(string(out), " "))

	// for _, i := range strings.Split(string(out), "\n") {
	// 	for _, j := range strings.Split(strings.Trim(i, " "), " ") {
	// 		fmt.Println(strings.Trim(j, " "))
	// 	}
	// }

	//

}
