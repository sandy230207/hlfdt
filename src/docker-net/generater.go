package dockernet

import (
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"fmt"
)

func ConvertConf(conf *config.Config) (Config, error) {
	services := make(map[string]Container)
	peerNames := make(map[string]string)
	peers := []string{}
	for _, v := range conf.Organizations {
		committingPeers, committingPeerNames, committers, err := extractPeers(v.Name, v.Type, v.CommittingPeers)
		if err != nil {
			return Config{}, fmt.Errorf("error occur in extracting committing peers: %v", err)
		}
		endorsingPeers, endorsingPeerNames, endorsers, err := extractPeers(v.Name, v.Type, v.EndorsingPeers)
		if err != nil {
			return Config{}, fmt.Errorf("error occur in extracting endorsing peers: %v", err)
		}
		orderingPeers, orderingPeerNames, orderers, err := extractPeers(v.Name, v.Type, v.Peers)
		if err != nil {
			return Config{}, fmt.Errorf("error occur in extracting ordering peers: %v", err)
		}
		committingPeers = mergeMaps(committingPeers, endorsingPeers)
		committingPeers = mergeMaps(committingPeers, orderingPeers)
		services = mergeMaps(committingPeers, services)
		committingPeerNames = mergeStringMaps(committingPeerNames, endorsingPeerNames)
		committingPeerNames = mergeStringMaps(committingPeerNames, orderingPeerNames)
		peerNames = mergeStringMaps(committingPeerNames, peerNames)
		committers = append(committers, endorsers...)
		committers = append(committers, orderers...)
		peers = append(peers, committers...)
	}
	services["cli"] = Container{
		ContainerName: "cli",
		Image:         "hyperledger/fabric-tools:$IMAGE_TAG",
		Tty:           true,
		StdinOpen:     true,
		Environment: []string{
			"GOPATH=/opt/gopath",
			"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
			"FABRIC_LOGGING_SPEC=INFO",
			"CORE_PEER_TLS_ENABLED=true",
			"CORE_PEER_PROFILE_ENABLED=true",
			"CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt",
			"CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key",
			"CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt",
		},
		WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric/peer",
		Command:    "/bin/bash",
		Volumes: []string{
			"/var/run/:/host/var/run/",
			"../organizations:/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations",
			"../scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/",
		},
		DependsOn: peers,
		Networks:  []string{"test"},
	}
	netConfig := Config{
		Version:  "2",
		Volumes:  peerNames,
		Networks: Network{},
		Services: services,
	}
	return netConfig, nil
}

func extractPeers(orgName string, orgType string, peers []config.Peer) (map[string]Container, map[string]string, []string, error) {
	services := make(map[string]Container)
	peerNames := make(map[string]string)
	peerOrgPeers := []string{}
	for _, peer := range peers {
		peerNames[peer.Name] = ""
		host, err := utils.ExtractHost(peer.Name, 1)
		if err != nil {
			return services, peerNames, peerOrgPeers, fmt.Errorf("error occur in extracting host: %v", err)
		}
		if orgType == "peerOrg" {
			peerOrgPeers = append(peerOrgPeers, peer.Name)
			services[peer.Name] = Container{
				Image: "hyperledger/fabric-peer:$IMAGE_TAG",
				Environment: []string{
					"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
					"CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=${COMPOSE_PROJECT_NAME}_test",
					"FABRIC_LOGGING_SPEC=INFO",
					"CORE_PEER_TLS_ENABLED=true",
					"CORE_PEER_PROFILE_ENABLED=true",
					"CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt",
					"CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key",
					"CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt",
					"CORE_PEER_ID=" + peer.Name,
					"CORE_PEER_ADDRESS=" + peer.Name + ":" + peer.Port,
					"CORE_PEER_LISTENADDRESS=0.0.0.0:" + peer.Port,
					"CORE_PEER_CHAINCODEADDRESS=" + peer.Name + ":" + peer.DBPort,
					"CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:" + peer.DBPort,
					"CORE_PEER_GOSSIP_EXTERNALENDPOINT=" + peer.Name + ":" + peer.Port,
					"CORE_PEER_GOSSIP_BOOTSTRAP=" + peer.Name + ":" + peer.Port,
					"CORE_PEER_LOCALMSPID=" + orgName,
				},
				Volumes: []string{
					"/var/run/docker.sock:/host/var/run/docker.sock",
					"../organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/msp:/etc/hyperledger/fabric/msp",
					"../organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls:/etc/hyperledger/fabric/tls",
					peer.Name + ":/var/hyperledger/production",
				},
				WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric/peer",
				Command:    "peer node start",
				Ports: []string{
					peer.Port + ":" + peer.Port,
				},
				ContainerName: peer.Name,
				Networks:      []string{"test"},
			}
		} else if orgType == "orderOrg" {
			services[peer.Name] = Container{
				Image: "hyperledger/fabric-orderer:$IMAGE_TAG",
				Environment: []string{
					"FABRIC_LOGGING_SPEC=INFO",
					"ORDERER_GENERAL_LISTENADDRESS=0.0.0.0",
					"ORDERER_GENERAL_LISTENPORT=" + peer.Port,
					"ORDERER_GENERAL_GENESISMETHOD=file",
					"ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block",
					"ORDERER_GENERAL_LOCALMSPID=" + orgName,
					"ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp",
					"ORDERER_GENERAL_TLS_ENABLED=true",
					"ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
					"ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
					"ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
					"ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR=1",
					"ORDERER_KAFKA_VERBOSE=true",
					"ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
					"ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
					"ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
				},
				Volumes: []string{
					"../system-genesis-block/genesis.block:/var/hyperledger/orderer/orderer.genesis.block",
					"../organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/msp:/var/hyperledger/orderer/msp",
					"../organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/:/var/hyperledger/orderer/tls",
					peer.Name + ":/var/hyperledger/production/orderer",
				},
				WorkingDir: "/opt/gopath/src/github.com/hyperledger/fabric",
				Command:    "orderer",
				Ports: []string{
					peer.Port + ":" + peer.Port,
				},
				ContainerName: peer.Name,
				Networks:      []string{"test"},
			}
		}
	}
	return services, peerNames, peerOrgPeers, nil
}

func mergeMaps(map1 map[string]Container, map2 map[string]Container) map[string]Container {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}

func mergeStringMaps(map1 map[string]string, map2 map[string]string) map[string]string {
	for k, v := range map2 {
		map1[k] = v
	}
	return map1
}
