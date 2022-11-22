package helper

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchtransport"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func GetOpenSearch(ctx context.Context) (testcontainers.Container, *opensearch.Client) {
	req := testcontainers.ContainerRequest{
		Image: "opensearchproject/opensearch:2.3.0",
		Env: map[string]string{
			"discovery.type": "single-node",
		},
		WaitingFor:   wait.ForLog(".opendistro_security is used as internal security index."),
		ExposedPorts: []string{"9200/tcp"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	Expect(err).ToNot(HaveOccurred())

	Expect(err).ToNot(HaveOccurred())
	Expect(container.ContainerIP(ctx)).ToNot(BeEmpty())
	return container, GetOpenSearchClient(ctx, container)
}

func GetOpenSearchClient(ctx context.Context, container testcontainers.Container) *opensearch.Client {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	containerIP, err := container.ContainerIP(ctx)
	Expect(err).ToNot(HaveOccurred())
	osClient, err := opensearch.NewClient(opensearch.Config{
		Addresses:         []string{fmt.Sprintf("https://%s:9200", containerIP)},
		Username:          "admin",
		Password:          "admin",
		Transport:         client.Transport,
		EnableDebugLogger: false,
		Logger: &opensearchtransport.CurlLogger{
			EnableRequestBody: true,
			Output:            GinkgoWriter,
		},
	})
	return osClient
}
