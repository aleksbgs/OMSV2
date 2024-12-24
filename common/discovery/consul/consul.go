package consul

import (
	"context"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"strconv"
	"strings"
)

type Registry struct {
	client *consul.Client
}

func (r *Registry) Register(ctx context.Context, instanceId, serviceName, hostPort string) error {
	log.Printf("Register service: %s", serviceName)

	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid host port %s", hostPort)
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	host := parts[0]
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceId,
		Address: host,
		Port:    port,
		Name:    serviceName,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceId,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceId, serviceName string) error {
	log.Printf("Deregister service: %s", serviceName)
	return r.client.Agent().CheckDeregister(instanceId)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	log.Printf("Discover service: %s", serviceName)
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var instances []string
	for _, entry := range entries {
		instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}
	return instances, nil
}

func (r *Registry) HealthCheck(instanceID, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", consul.HealthPassing)
}

func NewRegistry(addr string, serviceName string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client}, nil
}
