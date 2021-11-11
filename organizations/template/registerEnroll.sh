function createOrg() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/${ORG}.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/${ORG}.example.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:${PORT} --caname ca-${ORG} --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-${PORT}-ca-${ORG}.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-${PORT}-ca-${ORG}.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-${PORT}-ca-${ORG}.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-${PORT}-ca-${ORG}.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/${ORG}.example.com/msp/config.yaml

  infoln "Registering peer0"
  set -x
  fabric-ca-client register --caname ca-${ORG} --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca-${ORG} --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca-${ORG} --id.name ${ORG}admin --id.secret ${ORG}adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Generating the peer0 msp"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:${PORT} --caname ca-${ORG} -M ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/msp --csr.hosts peer0.${ORG}.example.com --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/msp/config.yaml

  infoln "Generating the peer0-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:${PORT} --caname ca-${ORG} -M ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls --enrollment.profile tls --csr.hosts peer0.${ORG}.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/server.key

  mkdir -p ${PWD}/organizations/peerOrganizations/${ORG}.example.com/msp/tlscacerts
  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/${ORG}.example.com/msp/tlscacerts/ca.crt

  mkdir -p ${PWD}/organizations/peerOrganizations/${ORG}.example.com/tlsca
  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/${ORG}.example.com/tlsca/tlsca.${ORG}.example.com-cert.pem

  mkdir -p ${PWD}/organizations/peerOrganizations/${ORG}.example.com/ca
  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/peers/peer0.${ORG}.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/${ORG}.example.com/ca/ca.${ORG}.example.com-cert.pem

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:${PORT} --caname ca-${ORG} -M ${PWD}/organizations/peerOrganizations/${ORG}.example.com/users/User1@${ORG}.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/${ORG}.example.com/users/User1@${ORG}.example.com/msp/config.yaml

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://${ORG}admin:${ORG}adminpw@localhost:${PORT} --caname ca-${ORG} -M ${PWD}/organizations/peerOrganizations/${ORG}.example.com/users/Admin@${ORG}.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/${ORG}/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/${ORG}.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/${ORG}.example.com/users/Admin@${ORG}.example.com/msp/config.yaml
}