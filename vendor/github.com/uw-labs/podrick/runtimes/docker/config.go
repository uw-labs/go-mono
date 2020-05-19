package docker

import (
	"strings"

	ct "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	units "github.com/docker/go-units"
	"github.com/uw-labs/podrick"
)

func createConfig(conf *podrick.ContainerConfig) (*ct.Config, *ct.HostConfig, *network.NetworkingConfig) {
	dc := &ct.Config{
		Image:        conf.Repo + ":" + conf.Tag,
		Env:          conf.Env,
		Cmd:          conf.Cmd,
		ExposedPorts: nat.PortSet{nat.Port(conf.Port): struct{}{}},
	}
	for _, p := range conf.ExtraPorts {
		dc.ExposedPorts[nat.Port(p)] = struct{}{}
	}
	if conf.Entrypoint != nil {
		dc.Entrypoint = strings.Split(*conf.Entrypoint, " ")
	}

	hc := &ct.HostConfig{
		PublishAllPorts: true,
	}
	for _, ulimit := range conf.Ulimits {
		hc.Ulimits = append(hc.Ulimits, &units.Ulimit{
			Name: ulimit.Name,
			Hard: ulimit.Hard,
			Soft: ulimit.Soft,
		})
	}

	nc := &network.NetworkingConfig{}
	return dc, hc, nc
}
