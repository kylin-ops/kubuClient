package kubeClient

import (
	"context"
	"encoding/json"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(configPath string) (*Kube, error) {
	kube := &Kube{}
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}
	configSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	kube.Client = configSet
	return kube, nil
}

var ctx = context.Background()

type Kube struct {
	configPath string
	Client     *kubernetes.Clientset
}

//
func resourceToStruct(resource string, obj interface{}) error {
	d, err := yaml.ToJSON([]byte(resource))
	if err != nil {
		return err
	}
	return json.Unmarshal(d, obj)
}

func (k *Kube) DeploymentGet(namespace, departmentName string) (*v1.Deployment, error) {
	return k.Client.AppsV1().Deployments(namespace).Get(ctx, departmentName, metav1.GetOptions{})
}

func (k *Kube) DeploymentList(namespace string) (*v1.DeploymentList, error) {
	return k.Client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
}

func (k *Kube) DeploymentCreate(namespace, resource string) (*v1.Deployment, error) {
	var deployment v1.Deployment
	if err := resourceToStruct(resource, &deployment); err != nil {
		return nil, err
	}
	return k.Client.AppsV1().Deployments(namespace).Create(ctx, &deployment, metav1.CreateOptions{})
}

func (k *Kube) DeploymentUpdate(namespace, resource string) (*v1.Deployment, error) {
	var deployment v1.Deployment
	if err := resourceToStruct(resource, &deployment); err != nil {
		return nil, err
	}
	return k.Client.AppsV1().Deployments(namespace).Update(ctx, &deployment, metav1.UpdateOptions{})
}

func (k *Kube) DeploymentApply(namespace, resource string) (*v1.Deployment, error) {
	var deployment appsv1.DeploymentApplyConfiguration
	if err := resourceToStruct(resource, &deployment); err != nil {
		return nil, err
	}
	return k.Client.AppsV1().Deployments(namespace).Apply(ctx, &deployment, metav1.ApplyOptions{})
}

func (k *Kube) DeploymentDelete(namespace, deploymentName string) error {
	return k.Client.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, metav1.DeleteOptions{})
}

func (k *Kube) PodList(namespace string) (*corev1.PodList, error) {
	return k.Client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
}

func (k *Kube) PodGet(namespace, podName string) (*corev1.Pod, error) {
	return k.Client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
}

func (k *Kube) PodLog(namespace, podName string) *rest.Request {
	return k.Client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})
}
