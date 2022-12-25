package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"sigs.k8s.io/kind/pkg/cluster"
)

type TestKindCluster struct {
	Name           string
	kubeconfigFile string
	ClientSet      *kubernetes.Clientset
	CancelFunc     func()
}

func CreateCluster(name string, providerOptions []cluster.ProviderOption, createOptions []cluster.CreateOption) (*TestKindCluster, error) {
	ret := TestKindCluster{Name: name}

	kubeconfigPrefix := fmt.Sprintf("test-kind-%s-", name)
	file, err := os.CreateTemp(os.TempDir(), kubeconfigPrefix)
	if err != nil {
		log.Fatal(err)
	}
	ret.kubeconfigFile = file.Name()

	kind := cluster.NewProvider(providerOptions...)
	noKubeconfig := cluster.CreateWithKubeconfigPath(file.Name())
	createOptions = append(createOptions, noKubeconfig)
	log.Infof("starting to create %s cluster", name)
	err = kind.Create(name, createOptions...)
	log.Infof("cluster %s started", name)
	if err != nil {
		return &ret, err
	}

	ret.CancelFunc = func() {
		err := kind.Delete(name, file.Name())
		if err != nil {
			log.Errorf("Can't delete cluster: %v", err)
		}
		err = os.Remove(file.Name())
		if err != nil {
			log.Errorf("Can't delete kubeconfig file %s due to %v", file.Name(), err)
		}
	}

	_, err = ret.getClientSet(false)

	return &ret, err
}

func (tkc *TestKindCluster) getClientSet(internal bool) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", tkc.kubeconfigFile)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	tkc.ClientSet = clientSet
	return clientSet, err
}
