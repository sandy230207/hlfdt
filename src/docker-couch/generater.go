package dockercouch

import (
	"fabric-tool/src/config"
	"fmt"
	"strings"
)

func ConvertConf(conf *config.Config) (Config, error) {
	services := make(map[string]interface{})
	for _, v := range conf.Organizations {
		if v.Type == "peerOrg" {
			committingPeers, err := extractPeers(v.CommittingPeers)
			if err != nil {
				return Config{}, fmt.Errorf("error occur in extracting committing peers: %v", err)
			}
			endorsingPeers, err := extractPeers(v.EndorsingPeers)
			if err != nil {
				return Config{}, fmt.Errorf("error occur in extracting endorsing peers: %v", err)
			}
			committingPeers = mergeMaps(committingPeers, endorsingPeers)
			services = mergeMaps(committingPeers, services)
		}
	}
	couchConfig := Config{
		Version:  "2",
		Networks: Network{},
		Services: services,
	}
	return couchConfig, nil
}

func extractPeers(peers []config.Peer) (map[string]interface{}, error) {
	services := make(map[string]interface{})
	for _, peer := range peers {
		peerName := strings.Replace(peer.Name, ".", "_", -1)
		dbName := "couchdb_" + peerName
		services[dbName] = DBContainer{
			Image: "couchdb:3.1.1",
			Environment: []string{
				"COUCHDB_USER=admin",
				"COUCHDB_PASSWORD=adminpw",
			},
			Ports: []string{
				peer.DBPort + ":" + "5984",
			},
			ContainerName: dbName,
			Networks:      []string{"test"},
		}
		services[peerName] = Container{
			Environment: []string{
				"CORE_LEDGER_STATE_STATEDATABASE=CouchDB",
				"CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=" + dbName + ":5984",
			},
			DependsOn: []string{
				dbName,
			},
		}
	}
	return services, nil
}

func mergeMaps(map1 map[string]interface{}, map2 map[string]interface{}) map[string]interface{} {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}
