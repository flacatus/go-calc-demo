package e2e

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("E2E Server Tests", func() {
	ginkgo.It("should return 5 when 2 and 3 are added", func() {
		resp, err := http.Get("http://localhost:8080/add?a=2&b=3")
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())

		expectedResponse := `{"valueA":2,"valueB":3,"result":5}`
		actualResponse := strings.TrimSpace(string(body)) // Trim any extra whitespace

		fmt.Println("Expected:", expectedResponse)
		fmt.Println("Actual:", actualResponse)

		gomega.Expect(actualResponse).To(gomega.Equal(expectedResponse))
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK))
	})
})
