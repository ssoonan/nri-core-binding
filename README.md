# NRI Core Binding Plugin

This is a Node Resource Interface (NRI) plugin for containerd that enables CPU core binding for containers based on Kubernetes pod annotations.

## What it does

The plugin reads pod annotations such as `scheduling.waters2019.io/core-id`, `core-id`, or `cpuset.cpus` and adjusts the container's OCI spec to bind the container to the specified CPU cores. This allows fine-grained control over CPU affinity for workloads.

## Building

To build the plugin, run the build script:

```bash
./build-nri.sh
```

This will compile the Go code and move the binary to `/opt/nri/plugins/`.

## Setup

To use this plugin, you need to configure containerd to enable NRI and point to the plugin directory.

For k3s, modify the containerd configuration at `/var/lib/rancher/k3s/agent/etc/containerd/config.toml.tmpl` and add the following section:

```toml
[plugins."io.containerd.nri.v1.nri"]
  disable = false
  plugin_path = "/opt/nri/plugins/"
  plugin_config_path = "/etc/nri/conf.d"
  socket_path = "/var/run/nri/nri.sock"
```

After updating the config, restart containerd or k3s for the changes to take effect.
