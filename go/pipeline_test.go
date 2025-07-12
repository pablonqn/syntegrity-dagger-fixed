package main

import (
	"context"
	"testing"

	"dagger.io/dagger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestPipeline(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Dagger Pipeline Suite")
}

var (
	ctx    context.Context
	client *dagger.Client
)

var _ = ginkgo.BeforeSuite(func() {
	var err error
	ctx = context.Background()
	client, err = dagger.Connect(ctx)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
})

var _ = ginkgo.AfterSuite(func() {
	if client != nil {
		gomega.Expect(client.Close()).To(gomega.Succeed())
	}
})

var _ = ginkgo.Describe("Dagger Pipeline", func() {
	var (
		pipeline *GoPipeline
		config   *Config
	)

	ginkgo.BeforeEach(func() {
		config = &Config{
			GoVersion: "1.21",
			Coverage:  80,
			Env:       "test",
			Branch:    "test",
		}
		pipeline = New(client, config)
	})

	ginkgo.Describe("Setup", func() {
		ginkgo.It("should setup the pipeline successfully", func() {
			err := pipeline.Setup(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Build", func() {
		ginkgo.It("should build successfully", func() {
			err := pipeline.Build(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Test", func() {
		ginkgo.It("should run tests successfully", func() {
			err := pipeline.Test(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Tag", func() {
		ginkgo.It("should tag successfully (no-op)", func() {
			err := pipeline.Tag(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Package", func() {
		ginkgo.It("should package successfully (no-op)", func() {
			err := pipeline.Package(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("Push", func() {
		ginkgo.It("should push successfully (no-op)", func() {
			err := pipeline.Push(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})
})
