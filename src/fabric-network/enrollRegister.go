package fabricnetwork

import (
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"fmt"
	"strings"
)

func GenerateEnrollRegister(conf *config.Config) (string, error) {
	res := "#!/bin/bash\n\n"
	for _, org := range conf.Organizations {
		res = res + "function create" + org.Name + "() {\n"
		res = res + "  infoln \"Enrolling the CA admin\"\n"
		caName, caPort, err := findCA(conf.CertificateAuthorities, org.Name)
		if err != nil {
			return res, fmt.Errorf("error occur in finding CA: %v", err)
		}
		if org.Type == "peerOrg" {
			peers := org.CommittingPeers
			if len(peers) == 0 {
				peers = org.EndorsingPeers
			}
			host, err := utils.ExtractHost(peers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + "  mkdir -p organizations/peerOrganizations/" + host + "/\n\n"
			res = res + "  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/" + host + "/\n\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client enroll -u https://admin:adminpw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
			res = res + "  { set +x; } 2>/dev/null\n\n"
			res = res + "  echo 'NodeOUs:\n  Enable: true\n  ClientOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: client\n  PeerOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: peer\n  AdminOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: admin\n  OrdererOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/" + host + "/msp/config.yaml\n\n"
			for _, peer := range org.CommittingPeers {
				res = res + "  infoln \"Registering " + peer.Name + "\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client register --caname " + strings.Replace(caName, ".", "_", -1) + " --id.name " + peer.Name + " --id.secret " + peer.Name + "pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  infoln \"Registering " + peer.Name + "\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client register --caname " + strings.Replace(caName, ".", "_", -1) + " --id.name " + peer.Name + " --id.secret " + peer.Name + "pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			res = res + "  infoln \"Registering user\"\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client register --caname " + strings.Replace(caName, ".", "_", -1) + " --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
			res = res + "  { set +x; } 2>/dev/null\n\n"
			res = res + "  infoln \"Registering the org admin\"\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client register --caname " + strings.Replace(caName, ".", "_", -1) + " --id.name " + org.Name + "admin --id.secret " + org.Name + "adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
			res = res + "  { set +x; } 2>/dev/null\n\n"
			for _, peer := range org.CommittingPeers {
				res = res + "  infoln \"Registering the " + peer.Name + " msp\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client enroll -u https://" + peer.Name + ":" + peer.Name + "pw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/msp --csr.hosts " + peer.Name + " --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  infoln \"Registering the " + peer.Name + " msp\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client enroll -u https://" + peer.Name + ":" + peer.Name + "pw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/msp --csr.hosts " + peer.Name + " --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			for _, peer := range org.CommittingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/msp/config.yaml ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/msp/config.yaml\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/msp/config.yaml ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/msp/config.yaml\n"
			}
			res = res + "\n"
			for _, peer := range org.CommittingPeers {
				res = res + "  infoln \"Registering " + peer.Name + "-tls certificates\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client enroll -u https://" + peer.Name + ":" + peer.Name + "pw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls --enrollment.profile tls --csr.hosts " + peer.Name + " --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  infoln \"Registering " + peer.Name + "-tls certificates\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client enroll -u https://" + peer.Name + ":" + peer.Name + "pw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls --enrollment.profile tls --csr.hosts " + peer.Name + " --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			for _, peer := range org.CommittingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/ca.crt\n"
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/signcerts/* ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/server.crt\n"
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/keystore/* ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/server.key\n\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/ca.crt\n"
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/signcerts/* ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/server.crt\n"
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/keystore/* ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/server.key\n\n"
			}
			res = res + "  mkdir -p ${PWD}/organizations/peerOrganizations/" + host + "/msp/tlscacerts\n"
			for _, peer := range org.CommittingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/msp/tlscacerts/ca.crt\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/msp/tlscacerts/ca.crt\n"
			}
			res = res + "\n  mkdir -p ${PWD}/organizations/peerOrganizations/" + host + "/tlsca\n"
			for _, peer := range org.CommittingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/tlsca/tlsca." + host + "-cert.pem\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/tlsca/tlsca." + host + "-cert.pem\n"
			}
			res = res + "\n  mkdir -p ${PWD}/organizations/peerOrganizations/" + host + "/ca\n"
			for _, peer := range org.CommittingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/msp/cacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/ca/" + caName + "-cert.pem\n"
			}
			for _, peer := range org.EndorsingPeers {
				res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/peers/" + peer.Name + "/msp/cacerts/* ${PWD}/organizations/peerOrganizations/" + host + "/ca/" + caName + "-cert.pem\n"
			}
			res = res + "\n  infoln \"Generating the user msp\"\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client enroll -u https://user1:user1pw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/peerOrganizations/" + host + "/users/User1@" + host + "/msp --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
			res = res + "  { set +x; } 2>/dev/null\n\n"
			res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/msp/config.yaml ${PWD}/organizations/peerOrganizations/" + host + "/users/User1@" + host + "/msp/config.yaml\n\n"
			res = res + "  infoln \"Generating the org admin msp\"\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client enroll -u https://" + org.Name + "admin:" + org.Name + "adminpw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/peerOrganizations/" + host + "/users/Admin@" + host + "/msp --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n"
			res = res + "  { set +x; } 2>/dev/null\n\n"
			res = res + "  cp ${PWD}/organizations/peerOrganizations/" + host + "/msp/config.yaml ${PWD}/organizations/peerOrganizations/" + host + "/users/Admin@" + host + "/msp/config.yaml\n"
			res = res + "}\n\n"
		} else if org.Type == "orderOrg" {
			host, err := utils.ExtractHost(org.Peers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + "  mkdir -p organizations/ordererOrganizations/" + host + "\n\n"
			res = res + "  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/" + host + "\n\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client enroll -u https://admin:adminpw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n" // hardcode fabric-ca/org1
			res = res + "  { set +x; } 2>/dev/null\n\n"
			res = res + "  echo 'NodeOUs:\n  Enable: true\n  ClientOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: client\n  PeerOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: peer\n  AdminOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: admin\n  OrdererOUIdentifier:\n    Certificate: cacerts/localhost-" + caPort + "-" + strings.Replace(caName, ".", "_", -1) + ".pem\n    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/ordererOrganizations/" + host + "/msp/config.yaml\n\n"
			res = res + "  infoln \"Registering orderer\"\n"
			res = res + "  set -x\n"
			for _, peer := range org.Peers {
				res = res + "  infoln \"Registering " + peer.Name + "\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client register --caname " + strings.Replace(caName, ".", "_", -1) + " --id.name " + peer.Name + " --id.secret " + peer.Name + "pw --id.type orderer --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n" // hardcode fabric-ca/org1
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			res = res + "  infoln \"Registering the org admin\"\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client register --caname " + strings.Replace(caName, ".", "_", -1) + " --id.name " + org.Name + "admin --id.secret " + org.Name + "adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n" // hardcode fabric-ca/org1
			res = res + "  { set +x; } 2>/dev/null\n\n"
			for _, peer := range org.Peers {
				res = res + "  infoln \"Registering the " + peer.Name + " msp\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client enroll -u https://" + peer.Name + ":" + peer.Name + "pw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/msp --csr.hosts " + peer.Name + " --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n" // hardcode fabric-ca/org1
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			for _, peer := range org.Peers {
				res = res + "  cp ${PWD}/organizations/ordererOrganizations/" + host + "/msp/config.yaml ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/msp/config.yaml\n"
			}
			for _, peer := range org.Peers {
				res = res + "  infoln \"Registering " + peer.Name + "-tls certificates\"\n"
				res = res + "  set -x\n"
				res = res + "  fabric-ca-client enroll -u https://" + peer.Name + ":" + peer.Name + "pw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls --enrollment.profile tls --csr.hosts " + peer.Name + " --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n" // hardcode fabric-ca/org1
				res = res + "  { set +x; } 2>/dev/null\n\n"
			}
			for _, peer := range org.Peers {
				res = res + "  cp ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/ca.crt\n"
				res = res + "  cp ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/server.crt\n"
				res = res + "  cp ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/keystore/* ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/server.key\n\n"
			}
			for _, peer := range org.Peers {
				res = res + "  mkdir -p ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/msp/tlscacerts\n"
				res = res + "  cp ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/msp/tlscacerts/tlsca." + host + "-cert.pem\n"
			}
			res = res + "\n  mkdir -p ${PWD}/organizations/ordererOrganizations/" + host + "/msp/tlscacerts\n"
			for _, peer := range org.Peers {
				res = res + "  cp ${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + peer.Name + "/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/" + host + "/msp/tlscacerts/tlsca." + host + "-cert.pem\n"
			}
			res = res + "\n  infoln \"Generating the admin msp\"\n"
			res = res + "  set -x\n"
			res = res + "  fabric-ca-client enroll -u https://" + org.Name + "admin:" + org.Name + "adminpw@localhost:" + caPort + " --caname " + strings.Replace(caName, ".", "_", -1) + " -M ${PWD}/organizations/ordererOrganizations/" + host + "/users/Admin@" + host + "/msp --tls.certfiles ${PWD}/organizations/fabric-ca/" + org.Name + "/tls-cert.pem\n" // hardcode fabric-ca/org1
			res = res + "  { set +x; } 2>/dev/null\n\n"
			res = res + "  cp ${PWD}/organizations/ordererOrganizations/" + host + "/msp/config.yaml ${PWD}/organizations/ordererOrganizations/" + host + "/users/Admin@" + host + "/msp/config.yaml\n"
			res = res + "}"
		}
	}
	return res, nil
}

// findCA return the ca name, ca port
// according to org name
func findCA(cas []config.CertificateAuthority, orgName string) (string, string, error) {
	caName := ""
	caPort := ""
	for _, ca := range cas {
		for _, orgInCA := range ca.Organizations {
			if orgName == orgInCA {
				caName = ca.Name
				caPort = ca.Port
			}
		}
	}
	return caName, caPort, nil
}
