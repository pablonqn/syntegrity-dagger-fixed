package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.com/syntegrity/syntegrity-infra/internal/config"

	"dagger.io/dagger"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines"
	dockerGO "gitlab.com/syntegrity/syntegrity-infra/internal/pipelines/docker-go"
	goKit "gitlab.com/syntegrity/syntegrity-infra/internal/pipelines/go-kit"
)

type stepHandler struct {
	name     string
	execute  func(context.Context, pipelines.Pipeline) error
	requires func(pipelines.Pipeline) bool
	skip     func(pipelines.Config) bool
}

var stepHandlers = map[string]stepHandler{
	"setup": {
		name: "setup",
		execute: func(ctx context.Context, p pipelines.Pipeline) error {
			return p.(pipelines.Setupper).Setup(ctx)
		},
		requires: func(p pipelines.Pipeline) bool {
			_, ok := p.(pipelines.Setupper)
			return ok
		},
		skip: func(_ pipelines.Config) bool {
			return false
		},
	},
	"build": {
		name: "build",
		execute: func(ctx context.Context, p pipelines.Pipeline) error {
			return p.(pipelines.Builder).Build(ctx)
		},
		requires: func(p pipelines.Pipeline) bool {
			_, ok := p.(pipelines.Builder)
			return ok
		},
		skip: func(cfg pipelines.Config) bool {
			return cfg.OnlyTest
		},
	},
	"test": {
		name: "test",
		execute: func(ctx context.Context, p pipelines.Pipeline) error {
			return p.(pipelines.Tester).Test(ctx)
		},
		requires: func(p pipelines.Pipeline) bool {
			_, ok := p.(pipelines.Tester)
			return ok
		},
		skip: func(cfg pipelines.Config) bool {
			return cfg.OnlyBuild
		},
	},
	"tag": {
		name: "tag",
		execute: func(ctx context.Context, p pipelines.Pipeline) error {
			return p.(pipelines.Tagger).Tag(ctx)
		},
		requires: func(p pipelines.Pipeline) bool {
			_, ok := p.(pipelines.Tagger)
			return ok
		},
		skip: func(_ pipelines.Config) bool {
			return false
		},
	},
	"package": {
		name: "package",
		execute: func(ctx context.Context, p pipelines.Pipeline) error {
			return p.(pipelines.Packager).Package(ctx)
		},
		requires: func(p pipelines.Pipeline) bool {
			_, ok := p.(pipelines.Packager)
			return ok
		},
		skip: func(_ pipelines.Config) bool {
			return false
		},
	},
	"push": {
		name: "push",
		execute: func(ctx context.Context, p pipelines.Pipeline) error {
			return p.(pipelines.Pusher).Push(ctx)
		},
		requires: func(p pipelines.Pipeline) bool {
			_, ok := p.(pipelines.Pusher)
			return ok
		},
		skip: func(cfg pipelines.Config) bool {
			return cfg.SkipPush
		},
	},
}

func executeStep(ctx context.Context, pipe pipelines.Pipeline, cfg pipelines.Config, step string) error {
	handler, ok := stepHandlers[step]
	if !ok {
		return fmt.Errorf("invalid step: %s", step)
	}

	if !handler.requires(pipe) {
		return fmt.Errorf("pipeline does not support step: %s", step)
	}

	if handler.skip(cfg) {
		return nil
	}

	fmt.Printf("➡️  step: %s\n", handler.name)
	if err := handler.execute(ctx, pipe); err != nil {
		return fmt.Errorf("%s failed: %w", handler.name, err)
	}

	return nil
}

func executePipeline(ctx context.Context, pipe pipelines.Pipeline, cfg pipelines.Config) {
	steps := []string{"setup", "build", "test", "tag", "package", "push"}
	for _, step := range steps {
		if err := executeStep(ctx, pipe, cfg, step); err != nil {
			log.Fatalf("❌ %v", err)
		}
	}
	fmt.Println("✅ pipeline completed.")
}

func executeSingleStep(ctx context.Context, pipe pipelines.Pipeline, cfg pipelines.Config, step string) {
	if step != "setup" {
		if err := executeStep(ctx, pipe, cfg, "setup"); err != nil {
			log.Fatalf("❌ %v", err)
		}
	}

	if err := executeStep(ctx, pipe, cfg, step); err != nil {
		log.Fatalf("❌ %v", err)
	}
	fmt.Printf("✅ Step '%s' ejecutado correctamente.\n", step)
}

func main() {
	defaultGitAuth := "ssh"
	if os.Getenv("CI_JOB_TOKEN") != "" {
		defaultGitAuth = "https"
	}

	var (
		pipelineName = flag.String("pipeline", "go-kit", "Name of the pipeline to be executed")
		coverage     = flag.Float64("coverage", 90, "Minimum coverage percentage required (in: 90 for 90%)")
		branch       = flag.String("branch", "develop", "branch name")
		env          = flag.String("env", "dev", "Environment: dev, staging, prod")

		skipPush  = flag.Bool("skip-push", false, "Skip image push")
		onlyBuild = flag.Bool("only-build", false, "Run build only")
		onlyTest  = flag.Bool("only-test", false, "Run only tests")
		verbose   = flag.Bool("verbose", false, "Verbose mode")
		gitRef    = flag.String("git-ref", "main", "Branch name (default: main)")
		step      = flag.String("step", "", "Individual pipeline step to execute (setup, build, test, tag, package, push)")
		gitAuth   = flag.String("git-auth", defaultGitAuth, "Git authentication method: ssh or https")
	)

	flag.Parse()

	fmt.Println("Pipeline:", *pipelineName)
	fmt.Println("Git ref:", *gitRef)
	fmt.Println("Cobertura:", *coverage)

	ctx := context.Background()

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatalf("❌ error connecting to Dagger: %v", err)
	}
	defer func(client *dagger.Client) {
		err = client.Close()
		if err != nil {
			log.Printf("❌ error closing Dagger: %v", err)
		}
	}(client)

	pipelinesCfg := pipelines.NewConfig(
		pipelines.WithEnv(*env),
		pipelines.WithSkipPush(*skipPush),
		pipelines.WithOnlyBuild(*onlyBuild),
		pipelines.WithOnlyTest(*onlyTest),
		pipelines.WithVerbose(*verbose),
		pipelines.WithGitGitProtocol(*gitAuth),
		pipelines.WithCoverage(*coverage),
		pipelines.WithBranch(*branch),
		pipelines.WithRegistry(cfg.RegistryBaseURL, cfg.RegistryPass),
		pipelines.WithRegistryUser(cfg.RegistryUser),
		pipelines.WithGoVersion(cfg.GoVersion),
	)

	registry := pipelines.NewRegistry()
	registry.Register("go-kit", goKit.New)
	registry.Register("docker-go", dockerGO.New)

	p, err := registry.Get(*pipelineName, client, pipelinesCfg)
	if err != nil {
		log.Printf("❌ pipeline does not exit: %v", err)
		return
	}
	fmt.Printf("🚀 running pipeline: %s (%s)\n", p.Name(), pipelinesCfg.Env)

	if *step != "" {
		executeSingleStep(ctx, p, pipelinesCfg, *step)
	} else {
		executePipeline(ctx, p, pipelinesCfg)
	}
}
