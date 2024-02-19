package resource

type Metadata struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Resource[T any] struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       T        `json:"spec"`
}

type Deployment Resource[DeploymentSpec]

type DeploymentSpec struct {
	Replicas int32           `json:"replicas"`
	Selector Selector        `json:"selector"`
	Template PodTemplateSpec `json:"template"`
}

type Selector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

type PodTemplateSpec struct {
	Metadata TemplateMetadata `json:"metadata"`
	Spec     PodSpec          `json:"spec"`
}

type TemplateMetadata struct {
	Labels map[string]string `json:"labels"`
}

type PodSpec struct {
	Containers []Container `json:"containers"`
}

type Container struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command,omitempty"`
	Ports   []struct {
		Name          string `json:"name"`
		ContainerPort int    `json:"containerPort"`
	} `json:"ports,omitempty"`
}
