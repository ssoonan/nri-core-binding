package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/containerd/nri/pkg/api"
	"github.com/containerd/nri/pkg/stub"
)

type CoreBindingPlugin struct{}

func (p *CoreBindingPlugin) Configure(ctx context.Context, config, runtime, version string) (stub.EventMask, error) {
	log.Printf("CoreBindingPlugin configured: config=%s, runtime=%s, version=%s", config, runtime, version)
	// Returning 0 here is fine; the stub discovers handlers automatically.
	return 0, nil
}

func (p *CoreBindingPlugin) Synchronize(ctx context.Context, pods []*api.PodSandbox, containers []*api.Container) ([]*api.ContainerUpdate, error) {
	log.Printf("Synchronizing %d pods and %d containers", len(pods), len(containers))
	return nil, nil
}

func (p *CoreBindingPlugin) RunPodSandbox(ctx context.Context, pod *api.PodSandbox) error {
	log.Printf("Running pod sandbox: %s", pod.Id)
	return nil
}

func (p *CoreBindingPlugin) StopPodSandbox(ctx context.Context, pod *api.PodSandbox) error {
	log.Printf("Stopping pod sandbox: %s", pod.Id)
	return nil
}

func (p *CoreBindingPlugin) RemovePodSandbox(ctx context.Context, pod *api.PodSandbox) error {
	log.Printf("Removing pod sandbox: %s", pod.Id)
	return nil
}

func (p *CoreBindingPlugin) CreateContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container) (*api.ContainerAdjustment, []*api.ContainerUpdate, error) {
	log.Printf("Creating container: %s in pod %s", container.Id, pod.Id)

	// Read desired core binding from pod annotations.
	cpus := getCPUSetFromAnnotations(pod.Annotations)
	if cpus == "" {
		log.Printf("No core-id/cpuset annotation found for pod %s", pod.Id)
		return nil, nil, nil
	}

	// Build an OCI spec adjustment to set cpuset.cpus via LinuxResources.
	// NRI runtime will apply this to the container's OCI config before creation.
	adjust := &api.ContainerAdjustment{
		Linux: &api.LinuxContainerAdjustment{
			Resources: &api.LinuxResources{
				Cpu: &api.LinuxCPU{
					Cpus: cpus,
				},
			},
		},
	}

	log.Printf("Requesting OCI cpuset binding for container %s to CPUs=%s", container.Id, cpus)
	return adjust, nil, nil
}

func (p *CoreBindingPlugin) PostCreateContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container) error {
	log.Printf("Post-creating container: %s", container.Id)
	return nil
}

func (p *CoreBindingPlugin) StartContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container) error {
	log.Printf("Starting container: %s", container.Id)
	return nil
}

func (p *CoreBindingPlugin) PostStartContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container) error {
	log.Printf("Post-starting container: %s", container.Id)
	return nil
}

func (p *CoreBindingPlugin) UpdateContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container, resources *api.LinuxResources) ([]*api.ContainerUpdate, error) {
	log.Printf("Updating container: %s", container.Id)
	return nil, nil
}

func (p *CoreBindingPlugin) PostUpdateContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container) error {
	log.Printf("Post-updating container: %s", container.Id)
	return nil
}

func (p *CoreBindingPlugin) StopContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container) error {
	log.Printf("Stopping container: %s", container.Id)
	return nil
}

func (p *CoreBindingPlugin) RemoveContainer(ctx context.Context, pod *api.PodSandbox, container *api.Container) error {
	log.Printf("Removing container: %s", container.Id)
	return nil
}

// getCPUSetFromAnnotations extracts a cpuset string from known annotation keys.
// Accepts either a single core id (e.g., "3") or a cpuset list/range (e.g., "1,3-5").
func getCPUSetFromAnnotations(ann map[string]string) string {
	if ann == nil {
		return ""
	}
	keys := []string{
		"scheduling.waters2019.io/core-id",
		"core-id",
		"cpuset.cpus",
	}
	var raw string
	for _, k := range keys {
		if v, ok := ann[k]; ok {
			raw = strings.TrimSpace(v)
			if raw != "" {
				break
			}
		}
	}
	
	if raw == "" {
		return ""
	}
	// If strictly digits, normalize to that single CPU.
	if isDigits(raw) {
		return raw
	}
	// Allow only digits, comma, hyphen.
	for _, ch := range raw {
		if !(ch >= '0' && ch <= '9') && ch != ',' && ch != '-' {
			return ""
		}
	}
	// Basic cleanup: collapse spaces.
	raw = strings.ReplaceAll(raw, " ", "")
	return raw
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func main() {
	p := &CoreBindingPlugin{}
	s, err := stub.New(p)
	if err != nil {
		log.Fatalf("Failed to create NRI plugin stub: %v", err)
	}
	if err := s.Run(context.Background()); err != nil {
		log.Printf("Plugin exited (%v)", err)
		os.Exit(1)
	}
}
