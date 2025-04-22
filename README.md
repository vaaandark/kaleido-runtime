# kaleido-runtime

Container runtime for pod live migration.

## How to build

Build the binary, and then the binary will be located at `./_output/bin/runm`:

```shell
make build
```

## How to deploy

Install `runm` binary to `/usr/local/bin/runm`.

Refer to [containerd-config.5.toml](https://github.com/containerd/containerd/blob/main/docs/man/containerd-config.toml.5.md). Add the following lines to `/etc/containerd/config.toml`:

```toml
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runm]
          runtime_type = "io.containerd.runc.v2"
          pod_annotations = [ "kaleido.io/*",]
          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runm.options]
            BinaryName = "/usr/local/bin/runm"
```

Apply runtime resources to the cluster:

```shell
kubectl apply -f ./deploy/runtimeclass.yaml
```

## How to use

If a pod has annotation `kaleido.io/source-pod` and runtimeclass `kaleido-runtime`, the containers of the pod will be created from directory `/var/lib/criu/checkpoints/<pod-uid>/<container-name>`.

![demo](./images/demo.gif)
