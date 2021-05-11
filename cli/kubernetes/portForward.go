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
)

type PortForward struct {
	stopChan chan struct{}
}

func NewPortForward(kubernetesProvider *Provider, namespace string, podName string, localPort uint16, podPort uint16, cancel context.CancelFunc) (*PortForward, error) {
	dialer := getHttpDialer(kubernetesProvider, namespace, podName)
	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, podPort)}, stopChan, readyChan, out, errOut)
	if err != nil {
		return nil, err
	}
	go func() {
		err = forwarder.ForwardPorts() // this is blocking
		if err != nil {
			fmt.Printf("kubernetes port-forwarding error: %s", err)
			cancel()
		}
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
	parsed, _ := url.Parse(kubernetesProvider.clientConfig.Host)
	serverURL := url.URL{Scheme: "https", Path: path, Host: parsed.Host}

	return spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)
}
