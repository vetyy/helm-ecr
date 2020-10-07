# helm-ecr
![Helm3 supported](https://img.shields.io/badge/Helm%203-supported-green)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/vetyy/helm-ecr.svg)](https://github.com/vetyy/helm-ecr/releases)

Helm plugin that allows installing Helm charts from [Amazon ECR](https://aws.amazon.com/ecr/) stored as OCI Artifacts.

> :warning: **Notice**
> This is not an official plugin and does not use Helm's new experimental API.
> Main motivation for this plugin was to be able to install charts stored in ECR using existing tools
> like [Flux Helm Operator](https://github.com/fluxcd/helm-operator) until better options are available.

Plugin currently supports only Helm v3.

This plugin was motivated by [helm-s3](https://github.com/hypnoglow/helm-s3) plugin.

## Install


Install latest version:

    $ helm plugin install https://github.com/vetyy/helm-ecr.git

Install a specific release version:

    $ helm plugin install https://github.com/vetyy/helm-ecr.git --version 0.1.2

To use the plugin, you do not need any special dependencies. The installer will
download versioned release with prebuilt binary from [github releases](https://github.com/vetyy/helm-ecr/releases).

## Overview

Plugin provides a new made up protocol registered as `ecr://`.
There are no additional commands provided by the plugin and it integrates only with native Helm commands.

Because ECR repository basically stores only a single chart, but multiple versions,
we must provide a chart with a name derived from repository url.
We decided to use the last part of the repository url as chart name and all attached image tags will be the chart versions.

Given `ecr://aws_account_id.dkr.ecr.region.amazonaws.com/namespace/NAME`.
The `NAME` will be the name of the chart. See usage below for more details.

## Usage

Add Helm repo:

    $ helm repo add my-ecr ecr://aws_account_id.dkr.ecr.region.amazonaws.com/namespace/app
    "my-ecr" has been added to your repositories

Discover chart versions:

    $ helm search repo my-ecr -l
    NAME	                CHART        VERSION	APP VERSION	DESCRIPTION
    my-ecr/app	        0.1.1
    my-ecr/app	        0.1.0
    my-ecr/app	        my-image-tag

Install chart:

    $ helm install my-app my-ecr/app

Alternatively you can install chart without adding the Helm repo,
but you must specify the version (image tag) as the last part of the chart name.

    $ helm install my-app ecr://aws_project_id.dkr.ecr.region.amazonaws.com/namespace/app/0.1.0

## Helm Operator

You can also use this plugin with Flux Helm Operator and their Kubernetes CRD `HelmRelease`.

```yaml
---
apiVersion: helm.fluxcd.io/v1
kind: HelmRelease
metadata:
  labels:
    app: my-app
  name: my-app
spec:
  chart:
    name: app
    repository: ecr://aws_project_id.dkr.ecr.region.amazonaws.com/namespace/app
    version: 0.1.0
  values:
    some_key: some_value
```

## Uninstall

    $ helm plugin remove ecr

## Development

On regular plugin installation, helm triggers post-install hook
that downloads prebuilt versioned release of the plugin binary and installs it.
To disable this behavior, you need to pass `HELM_PLUGIN_NO_INSTALL_HOOK=true` to the installer:

    $ HELM_PLUGIN_NO_INSTALL_HOOK=true helm plugin install https://github.com/vetyy/helm-ecr.git
    Development mode: not downloading versioned release.
    Installed plugin: ecr

Next, you may want to ensure if you have all prerequisites to build the plugin from source:

    cd ~/.helm/plugins/helm-ecr
    make deps build

If you see no messages - build was successful. Try to run some helm commands
that involve the plugin, or jump straight into plugin development.

## License

[MIT](LICENSE)
