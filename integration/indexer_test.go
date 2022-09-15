package integration_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ = Describe("Indexer", Ordered, func() {
	var openSearchContainer testcontainers.Container
	var openSearchIP string
	var client *http.Client

	BeforeAll(func() {
		req := testcontainers.ContainerRequest{
			Image: "opensearchproject/opensearch:2.3.0",
			Env: map[string]string{
				"discovery.type": "single-node",
			},
			WaitingFor:   wait.ForLog(".opendistro_security is used as internal security index."),
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
		openSearchContainer = container
		openSearchIP = ip

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	})

	AfterAll(func() {
		openSearchContainer.Terminate(context.Background())
	})

	It("That indexing works", func() {
		req, err := http.NewRequest(
			http.MethodPut,
			fmt.Sprintf("https://%s:9200/index/_doc/test", openSearchIP),
			strings.NewReader(`{"foo": "bar"}`),
		)
		req.Header.Add("Content-Type", "application/json")

		Expect(err).ToNot(HaveOccurred())

		req.SetBasicAuth("admin", "admin")
		resp, err := client.Do(req)

		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusCreated))

		req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s:9200/index/_source/test", openSearchIP), nil)
		Expect(err).ToNot(HaveOccurred())

		req.SetBasicAuth("admin", "admin")

		resp, err = client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()

		var responseBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		Expect(err).ToNot(HaveOccurred())

		Expect("bar").To(BeEquivalentTo(responseBody["foo"]))
	})
})
