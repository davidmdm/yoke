package resource

type Deployment Resource[DeploymentSpec]

type DeploymentSpec struct {
	Replicas int32           `json:"replicas"`
	Selector Selector        `json:"selector"`
	Template PodTemplateSpec `json:"template"`
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

type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int    `json:"containerPort"`
}
type Container struct {
	Name    string          `json:"name"`
	Image   string          `json:"image"`
	Command []string        `json:"command,omitempty"`
	Ports   []ContainerPort `json:"ports,omitempty"`
}
