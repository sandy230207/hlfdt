package dockerca

import (
	"fabric-tool/src/config"
	"strings"
)

func ConvertConf(conf *config.Config) (Config, error) {
	services := map[string]Container{}
	for _, v := range conf.CertificateAuthorities {
		ca := strings.Replace(v.Name, ".", "_", -1)
		volumes := []string{}
		// for _, org := range v.Organizations {
		volumes = append(volumes, "../organizations/fabric-ca/"+v.Organizations[0]+":/etc/hyperledger/fabric-ca-server")
		// }
		services[ca] = Container{
			Image: "hyperledger/fabric-ca:$IMAGE_TAG",
			Environment: []string{
				"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
				"FABRIC_CA_SERVER_CA_NAME=" + strings.Replace(v.Name, ".", "_", -1),
				"FABRIC_CA_SERVER_TLS_ENABLED=true",
				"FABRIC_CA_SERVER_PORT=" + v.Port,
			},
			Ports: []string{
				v.Port + ":" + v.Port,
			},
			Command:       "sh -c 'fabric-ca-server start -b admin:adminpw -d'",
			Volumes:       volumes,
			ContainerName: ca,
			Networks:      []string{"test"},
		}
	}
	caConfig := Config{
		Version:  "2",
		Networks: Network{},
		Services: services,
	}
	return caConfig, nil
}
