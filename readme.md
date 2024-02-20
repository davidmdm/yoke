# Halloumi - Infrastructure-as-Code (IaC) Package Deployer for Kubernetes

<img src="./halloumi.jpg" width=600/>

## Overview

Halloumi is a Helm-inspired infrastructure-as-code (IaC) package deployer.

The philosophy behind Halloumi is that Kubernetes packages should be described via code. Programming environments have control flow, test frameworks, static typing, documentation, error management, and versioning. They are ideal for building contracts and enforcing them.

Halloumi deploys "platters" to Kubernetes (think helm charts or packages). A platter is a wasm executable that outputs the Kubernetes resources making up the package as JSON/YAML to stdout.

Halloumi embeds a pure-Go wasm runtime (wazero) and deploys your platter to Kubernetes. It keeps track of the different revisions for any given release and provides capabilities such as rollbacks and inspection. For ArgoCD compatibility, it can render the platter to the filesystem as raw YAML resources.

## Theme

Every Kubernetes related project needs a theme. Although K8 has historically inspired nautical themes, Halloumi is a slight departure from the norm as it tries to move away from the YAML centric world-view of Kubernetes. Therefore Halloumi has an aviation related theme. Core commands, for example, are named `takeoff`, `descent`, and `runway`. However their less whimsical aliases exist as well: `up / apply`, `down / rollback`, and `render / export`.

## Installation

From source:

```bash
go install github.com/davidmdm/halloumi/cmd/halloumi@latest
```

## Usage

### takeoff (deploy / up)

```bash
# deploy local halloumi platter
halloumi takeoff my-release ./release.wasm

# deploy local halloumi platter and pass arguments to it
halloumi takeoff my-release ./release.wasm -- [args...]

# deploy remote halloumi platter via http
halloumi takeoff my-release https://github.com/my_org/infra/releases/platter-v0.1.0.wasm -- [args...]
```

### descent (rollback / down)

TODO - under development

### runway (render / export)

TODO - under development

## Contributions

Contributions are welcome! If you encounter any issues or have suggestions, please open an issue on the Halloumi GitHub repository.

## License

This project is licensed under the MIT License.
