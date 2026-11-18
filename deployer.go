package concourse

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/netip"

	"github.com/ChrisEineke/concourse/documents"
	docs "github.com/ChrisEineke/concourse/documents"
	"github.com/ChrisEineke/concourse/types"
	"github.com/chigopher/pathlib"
	"github.com/containerd/errdefs"
	"github.com/containerd/platforms"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/rogpeppe/rjson"
	"golang.org/x/sync/errgroup"
)

type Deployer struct {
	apiClient *client.Client
}

func NewDeployer(apiClient *client.Client) *Deployer {
	return &Deployer{apiClient}
}

func (d *Deployer) DispatchDocument(ctx context.Context, docFilePath *pathlib.Path) error {
	rawDocument, err := docFilePath.ReadFile()
	if err != nil {
		return err
	}

	headerDoc := docs.DocumentHeaderV1_0{}
	if err := rjson.Unmarshal(rawDocument, &headerDoc); err != nil {
		return err
	}

	switch headerDoc.Type {
	case documents.ContainerSetDocumentType:
		switch {
		case headerDoc.Version.Equal(documents.ContainerSetV_1_0_Version):
			containerSetDoc := &docs.ContainerSetDocumentV_1_0{}
			if err := rjson.Unmarshal(rawDocument, &containerSetDoc); err != nil {
				return err
			}
			return d.DeployContainerSet_v1_0(ctx, containerSetDoc)
		default:
			return ErrorUnsupportedDocumentVersion
		}
	default:
		return ErrorUnsupportedDocumentType
	}
}

func (d *Deployer) DeployContainerSet_v1_0(ctx context.Context, doc *docs.ContainerSetDocumentV_1_0) error {
	slog.Info("Deploying ContainerSet V1.0...", "name", doc.ContainerSet.Name)

	doc.ContainerSet.ID = rand.Text()[0:8]
	slog.Debug("Generated unique set ID", "name", doc.ContainerSet.Name, "id", doc.ContainerSet.ID)

	if err := d.pullImages(ctx, doc); err != nil {
		slog.Debug("Failed to deploy ContainerSet", "name", doc.ContainerSet.Name, "error", err)
		return err
	}
	if err := d.startAndCreateContainers(ctx, doc); err != nil {
		slog.Debug("Failed to deploy ContainerSet", "name", doc.ContainerSet.Name, "error", err)
		return err
	}
	slog.Debug("Deployed ContainerSet", "name", doc.ContainerSet.Name)
	return nil
}

// TODO: add different image pull policies
func (d *Deployer) pullImages(ctx context.Context, doc *docs.ContainerSetDocumentV_1_0) error {
	eg := errgroup.Group{}

	for _, ctrDef := range doc.ContainerSet.Containers {
		eg.Go(func() error {
			slog.Debug("Checking if container image is in local registry...", "image", ctrDef.Image)
			_, err := d.apiClient.ImageInspect(ctx, ctrDef.Image)
			if !errdefs.IsNotFound(err) {
				slog.Debug("Container image is in local registry", "image", ctrDef.Image)
				return nil
			} else {
				slog.Debug("Container is not in local registry", "image", ctrDef.Image)
			}

			slog.Debug("Pulling container image from remote registry...", "image", ctrDef.Image)
			res, err := d.apiClient.ImagePull(ctx, ctrDef.Image, client.ImagePullOptions{
				Platforms: []ocispec.Platform{platforms.DefaultSpec()},
			})
			if err != nil {
				slog.Debug("Failed to pull container into local registry", "image", ctrDef.Image, "error", err)
				return err
			}
			if err := res.Wait(ctx); err != nil {
				slog.Debug("Failed to pull container into local registry", "image", ctrDef.Image, "error", err)
				return err
			}
			if err := res.Close(); err != nil {
				slog.Debug("Failed to pull container into local registry", "image", ctrDef.Image, "error", err)
				return err
			}
			slog.Debug("Pulled container image into local registry", "image", ctrDef.Image)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (d *Deployer) startAndCreateContainers(ctx context.Context, doc *docs.ContainerSetDocumentV_1_0) error {
	for _, ctrDef := range doc.ContainerSet.Containers {
		slog.Debug("Creating container in ContainerSet...", "containerSet", doc.ContainerSet.Name, "container", ctrDef.Image)
		portMap, ok := d.portsToPortMap(ctrDef.Ports)
		if !ok {
			return fmt.Errorf("failed to convert port bindings slice to port map: %v", ctrDef.Ports)
		}
		res1, err := d.apiClient.ContainerCreate(ctx, client.ContainerCreateOptions{
			Config: &container.Config{
				AttachStderr: true,
				AttachStdin:  true,
				AttachStdout: true,
				Image:        ctrDef.Image,
				OpenStdin:    true,
				Tty:          true,
			},
			HostConfig: &container.HostConfig{
				AutoRemove: true,
				Annotations: map[string]string{
					"orchestrator":      "concourse",
					"containerSet.name": doc.ContainerSet.Name,
				},
				PortBindings: portMap,
				Resources: container.Resources{
					NanoCPUs:          ctrDef.Resources.Minimum.Cpu.Mul(types.Giga.Decimal).IntPart(),
					MemoryReservation: ctrDef.Resources.Minimum.Memory.IntPart(),
					Memory:            ctrDef.Resources.Maximum.Memory.IntPart(),
				},
			},
			NetworkingConfig: &network.NetworkingConfig{},
			Platform:         &ocispec.Platform{},
			Name:             fmt.Sprintf("%s-%s-%s", doc.ContainerSet.Name, doc.ContainerSet.ID, ctrDef.Name),
		})
		if err != nil {
			slog.Debug("Failed to create container in ContainerSet", "containerSet", doc.ContainerSet.Name, "container", ctrDef.Image, "error", err)
			return err
		}
		slog.Debug("Created container in ContainerSet", "containerSet", doc.ContainerSet.Name, "container", ctrDef.Image, "ID", res1.ID)

		slog.Debug("Starting container in ContainerSet...", "containerSet", doc.ContainerSet.Name, "ID", res1.ID)
		res2, err := d.apiClient.ContainerStart(ctx, res1.ID, client.ContainerStartOptions{})
		if err != nil {
			slog.Debug("Failed to start container in ContainerSet", "containerSet", doc.ContainerSet.Name, "ID", res1.ID, "error", err)
			return err
		}
		slog.Debug("Started container in ContainerSet", "containerSet", doc.ContainerSet.Name, "ID", res1.ID, "warnings", res2)
	}
	return nil
}

func (d *Deployer) portsToPortMap(portMappings []docs.Ports_1_0) (network.PortMap, bool) {
	result := network.PortMap{}
	for _, portMapping := range portMappings {
		protocol := network.TCP
		if portMapping.Protocol != "" {
			protocol = portMapping.Protocol
		}
		containerPort, ok := network.PortFrom(portMapping.ContainerPort, protocol)
		if !ok {
			return result, false
		}
		result[containerPort] = []network.PortBinding{
			{
				HostIP:   netip.AddrFrom4([4]byte{0, 0, 0, 0}),
				HostPort: fmt.Sprintf("%d", portMapping.HostPort),
			},
		}
	}
	return result, true
}
