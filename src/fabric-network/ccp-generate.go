package fabricnetwork

import (
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"fmt"
	"io/ioutil"
	"strings"
)

func GenerateCCPGenerate(filename string, conf *config.Config) (string, error) {
	res := ""
	ccpFunction, err := ioutil.ReadFile(filename)
	if err != nil {
		return res, fmt.Errorf("error occur in read ccp script template file: %v", err)
	}
	res = string(ccpFunction)
	for _, org := range conf.Organizations {
		if org.Type == "peerOrg" {
			peers := org.CommittingPeers
			if len(peers) == 0 {
				peers = org.EndorsingPeers
			}
			host, err := utils.ExtractHost(peers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			caPeerPort := ""
			caName := ""
			caPem := "organizations/peerOrganizations/" + host + "/ca/"
			for _, ca := range conf.CertificateAuthorities {
				for _, orgInCA := range ca.Organizations {
					if org.Name == orgInCA {
						caPeerPort = ca.Port
						caName = ca.Name
						break
					}
				}
			}
			caPem = caPem + caName + "-cert.pem"
			peerPem := "organizations/peerOrganizations/" + host + "/tlsca/tlsca." + host + "-cert.pem"
			res = res + "export ORG_NAME=" + org.Name + "\n"
			res = res + "export CA_PEER_PORT=" + caPeerPort + "\n"
			res = res + "export PEER_PEM=" + peerPem + "\n"
			res = res + "export CA_PEM=" + caPem + "\n"
			res = res + "export CA_PEER_NAME=" + caName + "\n"
			for _, peer := range org.CommittingPeers {
				res = res + "echo \"$(json_ccp $ORG_NAME " + peer.Port + " $CA_PEER_PORT $PEER_PEM $CA_PEM " + peer.Name + " $CA_PEER_NAME)\" > organizations/peerOrganizations/" + host + "/connection-" + org.Name + "-" + strings.Replace(peer.Name, ".", "-", -1) + ".json\n"
				res = res + "echo \"$(yaml_ccp $ORG_NAME " + peer.Port + " $CA_PEER_PORT $PEER_PEM $CA_PEM " + peer.Name + " $CA_PEER_NAME)\" > organizations/peerOrganizations/" + host + "/connection-" + org.Name + "-" + strings.Replace(peer.Name, ".", "-", -1) + ".yaml\n\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "echo \"$(json_ccp $ORG_NAME " + peer.Port + " $CA_PEER_PORT $PEER_PEM $CA_PEM " + peer.Name + " $CA_PEER_NAME)\" > organizations/peerOrganizations/" + host + "/connection-" + org.Name + "-" + strings.Replace(peer.Name, ".", "-", -1) + ".json\n"
				res = res + "echo \"$(yaml_ccp $ORG_NAME " + peer.Port + " $CA_PEER_PORT $PEER_PEM $CA_PEM " + peer.Name + " $CA_PEER_NAME)\" > organizations/peerOrganizations/" + host + "/connection-" + org.Name + "-" + strings.Replace(peer.Name, ".", "-", -1) + ".yaml\n\n"
			}
		}
	}
	return res, nil
}
