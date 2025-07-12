package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"dagger.io/dagger"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var (
	ctx    context.Context
	client *dagger.Client
)

func TestPipelineSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Pipeline Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	ginkgo.DeferCleanup(cancel)

	var err error
	client, err = dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
})

var _ = ginkgo.AfterSuite(func() {
	if client != nil {
		err := client.Close()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
})
