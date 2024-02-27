package halloumi

import (
	"fmt"

	"github.com/davidmdm/halloumi/internal/k8s"
)

type Client struct {
	k8s *k8s.Client
}

func FromKubeConfig(path string) (*Client, error) {
	client, err := k8s.NewClientFromKubeConfig(path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize k8s client: %w", err)
	}
	return &Client{client}, nil
}
