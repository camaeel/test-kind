package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kind "sigs.k8s.io/kind/pkg/cluster"
	"testing"
	"time"
)

func TestCreateCluster(t *testing.T) {
	name := "test-123"
	testKindCluster, err := CreateCluster(name, []kind.ProviderOption{}, []kind.CreateOption{})
	defer testKindCluster.CancelFunc()
	require.NoError(t, err)
	time.Sleep(60 * time.Second)
	nodes, err := testKindCluster.ClientSet.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	require.NoError(t, err)
	assert.Len(t, nodes.Items, 1)
	assert.Equal(t, "Node", nodes.Items[0].Kind)
	fmt.Printf("Nodes: %v", nodes)
}
