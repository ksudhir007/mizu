# Mizu

## About

Mizu is a tool that helps developers troubleshoot their web applications developed and running in Kubernetes environment. It allows user to instantly tap into any Kubernetes pod, collect and present its traffic in a developer-friendly way.

### Key features
* simple and powerful CLI
* no installation required
* instantly capture HTTP requests sent to a given pod
* decode and present any HTTP requests, REST and gRPC API calls.

_Note: TBD_

## Quick start
Download your `mizu`:

* for Mac run `curl -O https://static.up9.com/mizu/main/darwin.amd64/mizu && chmod 755 ./mizu`
* for Linux run `curl -O https://static.up9.com/mizu/main/linux.amd64/mizu && chmod 755 ./mizu`


Run `mizu` and supply Kubernetes pod name to tap, for example:
```shell
mizu tap podname
```

_Notes:_ you should have `kubectl` configured to run against your Kubernetes cluster.

After `mizu` starts and successfully connects to the specified pod, you point your browser [the traffic viewer web interface](http://localhost:8899/) which is available at [http://localhost:8899/](http://localhost:8899/) 



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
_not implemented yet_
Connect to running `mizu` listener and open traffic viewer UI in browser.


### `fetch`
_not implemented yet_
Connect to running `mizu` listener and download collected web traffic files.


