package utils

import (
	// utils "github.com/redhat-appstudio/managed-gitops/utilities/load-test/loadtest"
	// metrics "k8s.io/apimachinery/pkg/apis/metrics/v1"
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

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
}
