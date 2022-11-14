package integration_test

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/kstiehl/index-bouncer/pkg/opensearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	oSearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v2/opensearchtransport"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ = Describe("opensearch", Ordered, func() {
	var openSearchContainer testcontainers.Container
	var openSearchIP string
	var client *http.Client
	var osClient *oSearch.Client

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
		osClient, err = oSearch.NewClient(oSearch.Config{
			Addresses:         []string{fmt.Sprintf("https://%s:9200", openSearchIP)},
			Username:          "admin",
			Password:          "admin",
			Transport:         client.Transport,
			EnableDebugLogger: false,
			Logger: &opensearchtransport.CurlLogger{
				EnableRequestBody: true,
				Output:            GinkgoWriter,
			},
		})

		Expect(err).ToNot(HaveOccurred())
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

	FIt("create datastream", func() {
		ctx := logr.NewContext(context.Background(), stdr.New(log.New(os.Stdout, "", log.Lshortfile)))

		existsRequest := opensearchapi.IndicesExistsIndexTemplateRequest{
			Name: "test",
		}

		res, err := existsRequest.Do(ctx, osClient)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.IsError()).To(BeTrue())

		err = opensearch.EnsureIndexTemplate(ctx, osClient, opensearch.DataStreamConfig{
			Name:          "test",
			IndexPatterns: []string{"test"},
		})
		Expect(err).ToNot(HaveOccurred())

		existsRequest = opensearchapi.IndicesExistsIndexTemplateRequest{
			Name: "test",
		}
		res, err = existsRequest.Do(ctx, osClient)
		defer res.Body.Close()
		Expect(err).ToNot(HaveOccurred())
		Expect(res.IsError()).To(BeFalse())
	})
})
