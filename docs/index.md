# Mizu

## About

Debug and troubleshoot your microservices with an open source tool that enables you to view the complete API traffic inside of your Kubernetes cluster. 

Think of TCPDump and Chrome Dev Tools combined to see what’s going on inside your Kubernetes cluster.


### Key features
* Simple and powerful CLI
* No installation or code instrumentation required
* Decode and present all HTTP requests, REST and gRPC API calls in real time
* Works completely on premisses (on-prem)

## Quick start
Get your `mizu`:

- for **Mac** - 
```
curl -O https://static.up9.com/mizu/main/darwin.amd64/mizu && chmod 755 ./mizu
```

- for **Linux** - 
```
curl -O https://static.up9.com/mizu/main/linux.amd64/mizu && chmod 755 ./mizu
```


Run `mizu` and supply Kubernetes pod name to tap, for example:

```
mizu tap <podname> [-n <namespace>]
```

_Notes:_
You should have `kubectl` configured to run against your Kubernetes cluster.
Namespace is required in case the pod is in a different namespace than the default one.

After `mizu` starts and successfully connects to the specified pod, you point your browser [the traffic viewer web interface](http://localhost:8899/) which is available at [http://localhost:8899/](http://localhost:8899/) 

![Mizu Local Webapp](img/mizu-snapshot.png)


## Commands and command-line arguments
Usage and list of command-line arguments can be seen by running `mizu -h` or `mizu help`

### `tap`
Listen to the specified pod and display collected web traffic in the Web UI

```
Usage: mizu tap PODNAME [flags]

Flags:
  -p, --gui-port uint16     Provide a custom port for the web interface webserver (default 8899)
  -h, --help                help for tap
  -k, --kubeconfig string   Path to kubeconfig file
      --mizu-image string   Custom image for mizu collector (default "gcr.io/up9-docker-hub/mizu/develop:latest")
      --mizu-port uint16    Port which mizu cli will attempt to forward from the mizu collector pod (default 8899)
  -n, --namespace string    Namespace selector
```

### `version`
Display `mizu` version

### `help`
Display usage and help information

### `view`
Connect to running `mizu` listener and open traffic viewer UI in browser.
_not implemented yet_


### `fetch`
Connect to running `mizu` listener and download collected web traffic files.
_not implemented yet_


---
version: 04
