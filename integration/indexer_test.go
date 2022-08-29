package integration_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ = Describe("Indexer", Ordered, func() {
	const esPassword = "jrjgkrejgkrejgkjredswjr8594839uwe09"
	var elasticContainer testcontainers.Container
	var elasticIP string
	var client *http.Client

	BeforeAll(func() {
		req := testcontainers.ContainerRequest{
			Image: "docker.elastic.co/elasticsearch/elasticsearch:8.4.0",
			Env: map[string]string{
				"ELASTIC_PASSWORD": esPassword,
			},
			WaitingFor:   wait.ForLog("Cluster health status changed from [YELLOW] to [GREEN]"),
			ExposedPorts: []string{"9200/tcp"},
		}

		ctx := context.Background()
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		Expect(err).ToNot(HaveOccurred())

		ip, err := container.ContainerIP(ctx)

		Expect(err).ToNot(HaveOccurred())
		Expect(ip).ToNot(BeEmpty())
		elasticContainer = container
		elasticIP = ip

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	})

	AfterAll(func() {
		elasticContainer.Terminate(context.Background())
	})

	It("this and that", func() {

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("https://%s:9200/index/_doc/test", elasticIP), strings.NewReader(`{"foo": "bar"}`))
		req.Header.Add("Content-Type", "application/json")

		Expect(err).ToNot(HaveOccurred())

		req.SetBasicAuth("elastic", esPassword)
		resp, err := client.Do(req)

		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))
	})
})
