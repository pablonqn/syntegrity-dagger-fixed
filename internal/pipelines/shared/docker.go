package shared

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// DockerDeployer encapsula la lógica para construir y publicar una imagen Docker.
type DockerDeployer struct {
	Client       *dagger.Client
	ImageName    string            // ejemplo: registry.gitlab.com/mi-org/mi-proyecto/servicio
	Source       *dagger.Directory // directorio con Dockerfile
	RegistryUser string            // normalmente "gitlab-ci-token"
	RegistryPass *dagger.Secret    // CI_JOB_TOKEN o GITLAB_PAT como secreto
	Tag          string            // ejemplo: "latest", "v1.2.3", "sha"
}

// NewDockerDeployer crea una instancia del deployer, asegurando que la contraseña se maneje como secreto.
func NewDockerDeployer(client *dagger.Client, imageName string, src *dagger.Directory, tag, user, pass string) *DockerDeployer {
	return &DockerDeployer{
		Client:       client,
		ImageName:    imageName,
		Source:       src,
		RegistryUser: user,
		RegistryPass: client.SetSecret("registry-pass", pass),
		Tag:          tag,
	}
}

// BuildAndPush construye la imagen y la publica en la registry.
func (d *DockerDeployer) BuildAndPush(ctx context.Context) error {
	fmt.Printf("🐳 Construyendo y publicando %s:%s...\n", d.ImageName, d.Tag)

	image := d.Client.Container().
		Build(d.Source).
		WithRegistryAuth(d.ImageName, d.RegistryUser, d.RegistryPass)

	ref := fmt.Sprintf("%s:%s", d.ImageName, d.Tag)
	_, err := image.Publish(ctx, ref)
	if err != nil {
		return fmt.Errorf("❌ error al publicar la imagen: %w", err)
	}

	fmt.Printf("✅ Imagen publicada: %s\n", ref)
	return nil
}
