package model_generator

import (
	"github.com/bborbe/kubernetes_tools/config"
	"github.com/bborbe/kubernetes_tools/model"
)

const K8S_DEFAULT_VERSION = "v1.3.5"

func GenerateModel(config *config.Cluster) *model.Cluster {
	cluster := new(model.Cluster)

	cluster.UpdateRebootStrategy = config.UpdateRebootStrategy
	if len(cluster.UpdateRebootStrategy) == 0 {
		cluster.UpdateRebootStrategy = "etcd-lock"
	}
	cluster.Version = config.KubernetesVersion
	if len(cluster.Version) == 0 {
		cluster.Version = K8S_DEFAULT_VERSION
	}
	cluster.Region = config.Region
	//	cluster.LvmVolumeGroup = config.LvmVolumeGroup

	for _, configHost := range config.Hosts {
		//c.Host = cluster.Host
		//c.Bridge = cluster.Bridge
		//c.ApiServerPublicIp = cluster.ApiServerPublicIp
		//c.Network = cluster.Network
		//c.Gateway = valueOf(cluster.Gateway, fmt.Sprintf("%s.1", cluster.Network))
		//c.Dns = valueOf(cluster.Dns, fmt.Sprintf("%s.1", cluster.Network))

		host := model.Host{}

		counter := 0
		for _, configNode := range configHost.Nodes {
			for i := 0; i < configNode.Amount; i++ {

				if configNode.Storage && configNode.Nfsd {
					panic("storage and nfsd at the same time is currently not supported")
				}

				//name := generateNodeName(n, i)
				node := model.Node{
					//Name:        name,
					//Ip:  fmt.Sprintf("%s.%d", configHost.KubernetesNetwork, counter + 10),
					//Mac: fmt.Sprintf("%s%02x", configHost.KubernetesNetwork, counter + 10),
					//VolumeName:  fmt.Sprintf("%s%s", cluster.VolumePrefix, name),
					//VmName:      fmt.Sprintf("%s%s", cluster.VmPrefix, name),
					Etcd:        configNode.Etcd,
					Worker:      configNode.Worker,
					Master:      configNode.Master,
					Storage:     configNode.Storage,
					Nfsd:        configNode.Nfsd,
					Cores:       configNode.Cores,
					Memory:      configNode.Memory,
					NfsSize:     configNode.NfsSize,
					StorageSize: configNode.StorageSize,
					RootSize:    valueOfSize(configNode.RootSize, "10G"),
					DockerSize:  valueOfSize(configNode.DockerSize, "10G"),
					KubeletSize: valueOfSize(configNode.KubeletSize, "10G"),
				}
				host.Nodes = append(host.Nodes, node)
				counter++
			}
		}
		cluster.Hosts = append(cluster.Hosts, host)
	}

	return cluster
}

func valueOfSize(size model.Size, defaultSize model.Size) model.Size {
	if len(size) == 0 {
		return defaultSize
	}
	return size
}

//func generateNodeName(node config.Node, number int) string {
//	if node.Amount == 1 {
//		return node.Name
//	} else {
//		return fmt.Sprintf("%s%d", node.Name, number)
//	}
//}

func valueOf(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
