package kubeClient

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	yaml2 "sigs.k8s.io/yaml"
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

func (k *Kube) ResourceYaml(resource interface{}) {
	jsonData, _ := json.Marshal(resource)
	yamlData, _ := yaml2.JSONToYAML(jsonData)
	fmt.Println(string(yamlData))
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

// 获取指定最后几行的日志
func (k *Kube) PodGetLogTailLines(namespace, podName string, tailLines int64) *rest.Request {
	return k.Client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{TailLines: &tailLines})
}

// 获取实时日志
// 获取日志的范例
//reader, _ := req.Stream(context.Background())
//buf := make([]byte, 1024*10)
//for {
//_, err := reader.Read(buf)
//if err != nil{
//return
//}
//fmt.Println(string(buf))
//}
func (k *Kube) PodGetLogFollow(namespace, podName string) *rest.Request {
	var line int64 = 5
	return k.Client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{Follow: true, TailLines: &line})
}

func (k *Kube) ServiceList(namespace string) (*corev1.ServiceList, error) {
	return k.Client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
}

func (k *Kube) ServiceGet(namespace string, serviceName string) (*corev1.Service, error) {
	return k.Client.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
}
