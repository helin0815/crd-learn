package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	// 引入自定义资源类型

	hltestv1 "github.com/helin0815/crd-learn/api/v1" // 替换为你的实际路径
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(可选) 绝对路径到kubeconfig文件")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "绝对路径到kubeconfig文件")
	}
	flag.Parse()

	// 使用给定的kubeconfig文件配置集群信息
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 创建客户端集
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 自定义资源客户端
	hltestClient, err := hltestv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 设置自定义资源参数
	hltestInstance := &hltestv1.HlTest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-hltest",
			Namespace: "default", // 更改为你的命名空间
		},
		Spec: hltestv1.HlTestSpec{
			User: "test",
		},
	}

	// 创建自定义资源
	_, err = hltestClient.ExampleV1().HlTests(hltestInstance.Namespace).Create(context.TODO(), hltestInstance, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) {
		fmt.Println("资源已存在，尝试更新...")
		_, updateErr := hltestClient.ExampleV1().HlTests(hltestInstance.Namespace).Update(context.TODO(), hltestInstance, metav1.UpdateOptions{})
		if updateErr != nil {
			panic(updateErr.Error())
		}
		fmt.Println("资源已更新")
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("资源创建成功")
	}
}
