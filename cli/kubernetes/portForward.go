package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"strings"
)

type PortForward struct {
	stopChan chan struct{}
}

func NewPortForward(kubernetesProvider *Provider, namespace string, podName string, localPort uint16, podPort uint16, ctx context.Context) (*PortForward, error) {
	retries := 0
	dialer := getHttpDialer(kubernetesProvider, namespace, podName)
	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, podPort)}, stopChan, readyChan, out, errOut)
	if err != nil {
		return nil, err
	}
	go func() {
		for ctx.Err() != nil && retries < 5 {
			err = forwarder.ForwardPorts() // this is blocking
			if err != nil {
				retries += 1
				fmt.Printf("kubernetes port-forwarding error: %s", err)
			}
		}
		fmt.Printf("Stopping to retry port-forward")
	}()
	return &PortForward{stopChan: stopChan}, nil
}

func (portForward *PortForward) Stop() {
	close(portForward.stopChan)
}

func getHttpDialer(kubernetesProvider *Provider, namespace string, podName string) httpstream.Dialer {
	roundTripper, upgrader, err := spdy.RoundTripperFor(&kubernetesProvider.clientConfig)
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, podName)
	hostIP := strings.TrimLeft(kubernetesProvider.clientConfig.Host, "htps:/")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}

	return spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)
}
