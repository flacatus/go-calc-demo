package e2e

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("E2E Server Tests", func() {
	routeURL := os.Getenv("CONTAINER_ROUTE_URL")

	ginkgo.It("should return 5 when 2 and 3 are added", func() {
		// Create an HTTP client that skips TLS verification
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		// Get the route URL from the environment
		resp, err := client.Get(routeURL + "/add?a=2&b=3")
		time.Sleep(5 * time.Second)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		expectedResponse := `{"valueA":2,"valueB":3,"result":5}`
		actualResponse := strings.TrimSpace(string(body)) // Trim any extra whitespace

		fmt.Println("Expected:", expectedResponse)
		fmt.Println("Actual:", actualResponse)

		// Assert the response body and status code
		gomega.Expect(actualResponse).To(gomega.Equal(expectedResponse))
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))
	})

	ginkgo.It("should return 6 when 2 and 3 are multiplied", func() {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		resp, err := client.Get(routeURL + "/mul?a=2&b=3")
		time.Sleep(5 * time.Second)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		expectedResponse := `{"valueA":2,"valueB":3,"result":6}`
		actualResponse := strings.TrimSpace(string(body))

		fmt.Println("Expected:", expectedResponse)
		fmt.Println("Actual:", actualResponse)

		// Assert the response body and status code
		gomega.Expect(actualResponse).To(gomega.Equal(expectedResponse))
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))
	})
})
