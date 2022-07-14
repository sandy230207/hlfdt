#!/bin/bash

function createorg1() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/org1.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca_org1_example_com --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca_org1_example_com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca_org1_example_com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca_org1_example_com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca_org1_example_com.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml

  infoln "Registering peer1.org1.example.com"
  set -x
  fabric-ca-client register --caname ca_org1_example_com --id.name peer1.org1.example.com --id.secret peer1.org1.example.compw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer2.org1.example.com"
  set -x
  fabric-ca-client register --caname ca_org1_example_com --id.name peer2.org1.example.com --id.secret peer2.org1.example.compw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer3.org1.example.com"
  set -x
  fabric-ca-client register --caname ca_org1_example_com --id.name peer3.org1.example.com --id.secret peer3.org1.example.compw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca_org1_example_com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca_org1_example_com --id.name org1admin --id.secret org1adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the peer1.org1.example.com msp"
  set -x
  fabric-ca-client enroll -u https://peer1.org1.example.com:peer1.org1.example.compw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp --csr.hosts peer1.org1.example.com --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the peer2.org1.example.com msp"
  set -x
  fabric-ca-client enroll -u https://peer2.org1.example.com:peer2.org1.example.compw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/msp --csr.hosts peer2.org1.example.com --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the peer3.org1.example.com msp"
  set -x
  fabric-ca-client enroll -u https://peer3.org1.example.com:peer3.org1.example.compw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/msp --csr.hosts peer3.org1.example.com --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp/config.yaml
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/msp/config.yaml
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/msp/config.yaml

  infoln "Registering peer1.org1.example.com-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer1.org1.example.com:peer1.org1.example.compw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls --enrollment.profile tls --csr.hosts peer1.org1.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer2.org1.example.com-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer2.org1.example.com:peer2.org1.example.compw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls --enrollment.profile tls --csr.hosts peer2.org1.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer3.org1.example.com-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer3.org1.example.com:peer3.org1.example.compw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls --enrollment.profile tls --csr.hosts peer3.org1.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/server.key

  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/server.key

  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/server.key

  mkdir -p ${PWD}/organizations/peerOrganizations/org1.example.com/msp/tlscacerts
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/ca.crt

  mkdir -p ${PWD}/organizations/peerOrganizations/org1.example.com/tlsca
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem

  mkdir -p ${PWD}/organizations/peerOrganizations/org1.example.com/ca
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer2.org1.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer3.org1.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/config.yaml

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://org1admin:org1adminpw@localhost:7054 --caname ca_org1_example_com -M ${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org1/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/config.yaml
}

function createorg2() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/peerOrganizations/org2.example.com/

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org2.example.com/

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:8054 --caname ca_org2_example_com --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-8054-ca_org2_example_com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-8054-ca_org2_example_com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-8054-ca_org2_example_com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-8054-ca_org2_example_com.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml

  infoln "Registering peer1.org2.example.com"
  set -x
  fabric-ca-client register --caname ca_org2_example_com --id.name peer1.org2.example.com --id.secret peer1.org2.example.compw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer2.org2.example.com"
  set -x
  fabric-ca-client register --caname ca_org2_example_com --id.name peer2.org2.example.com --id.secret peer2.org2.example.compw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer3.org2.example.com"
  set -x
  fabric-ca-client register --caname ca_org2_example_com --id.name peer3.org2.example.com --id.secret peer3.org2.example.compw --id.type peer --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering user"
  set -x
  fabric-ca-client register --caname ca_org2_example_com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca_org2_example_com --id.name org2admin --id.secret org2adminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the peer1.org2.example.com msp"
  set -x
  fabric-ca-client enroll -u https://peer1.org2.example.com:peer1.org2.example.compw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp --csr.hosts peer1.org2.example.com --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the peer2.org2.example.com msp"
  set -x
  fabric-ca-client enroll -u https://peer2.org2.example.com:peer2.org2.example.compw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/msp --csr.hosts peer2.org2.example.com --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the peer3.org2.example.com msp"
  set -x
  fabric-ca-client enroll -u https://peer3.org2.example.com:peer3.org2.example.compw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/msp --csr.hosts peer3.org2.example.com --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp/config.yaml
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/msp/config.yaml
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/msp/config.yaml

  infoln "Registering peer1.org2.example.com-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer1.org2.example.com:peer1.org2.example.compw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls --enrollment.profile tls --csr.hosts peer1.org2.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer2.org2.example.com-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer2.org2.example.com:peer2.org2.example.compw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls --enrollment.profile tls --csr.hosts peer2.org2.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering peer3.org2.example.com-tls certificates"
  set -x
  fabric-ca-client enroll -u https://peer3.org2.example.com:peer3.org2.example.compw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls --enrollment.profile tls --csr.hosts peer3.org2.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/server.key

  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/server.key

  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/signcerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/server.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/keystore/* ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/server.key

  mkdir -p ${PWD}/organizations/peerOrganizations/org2.example.com/msp/tlscacerts
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/msp/tlscacerts/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/msp/tlscacerts/ca.crt
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/msp/tlscacerts/ca.crt

  mkdir -p ${PWD}/organizations/peerOrganizations/org2.example.com/tlsca
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/tls/tlscacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem

  mkdir -p ${PWD}/organizations/peerOrganizations/org2.example.com/ca
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer2.org2.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem
  cp ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer3.org2.example.com/msp/cacerts/* ${PWD}/organizations/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem

  infoln "Generating the user msp"
  set -x
  fabric-ca-client enroll -u https://user1:user1pw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/config.yaml

  infoln "Generating the org admin msp"
  set -x
  fabric-ca-client enroll -u https://org2admin:org2adminpw@localhost:8054 --caname ca_org2_example_com -M ${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/org2/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/config.yaml
}

function createorderOrg() {
  infoln "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/example.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/example.com

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca_orderer_example_com --tls.certfiles ${PWD}/organizations/fabric-ca/orderOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca_orderer_example_com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca_orderer_example_com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca_orderer_example_com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca_orderer_example_com.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml

  infoln "Registering orderer"
  set -x
  infoln "Registering orderer.example.com"
  set -x
  fabric-ca-client register --caname ca_orderer_example_com --id.name orderer.example.com --id.secret orderer.example.compw --id.type orderer --tls.certfiles ${PWD}/organizations/fabric-ca/orderOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the org admin"
  set -x
  fabric-ca-client register --caname ca_orderer_example_com --id.name orderOrgadmin --id.secret orderOrgadminpw --id.type admin --tls.certfiles ${PWD}/organizations/fabric-ca/orderOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  infoln "Registering the orderer.example.com msp"
  set -x
  fabric-ca-client enroll -u https://orderer.example.com:orderer.example.compw@localhost:9054 --caname ca_orderer_example_com -M ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/orderOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/config.yaml
  infoln "Registering orderer.example.com-tls certificates"
  set -x
  fabric-ca-client enroll -u https://orderer.example.com:orderer.example.compw@localhost:9054 --caname ca_orderer_example_com -M ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls --enrollment.profile tls --csr.hosts orderer.example.com --csr.hosts localhost --tls.certfiles ${PWD}/organizations/fabric-ca/orderOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
  cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/signcerts/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.crt
  cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/keystore/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/server.key

  mkdir -p ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  mkdir -p ${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts
  cp ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/tlscacerts/* ${PWD}/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem

  infoln "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://orderOrgadmin:orderOrgadminpw@localhost:9054 --caname ca_orderer_example_com -M ${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp --tls.certfiles ${PWD}/organizations/fabric-ca/orderOrg/tls-cert.pem
  { set +x; } 2>/dev/null

  cp ${PWD}/organizations/ordererOrganizations/example.com/msp/config.yaml ${PWD}/organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/config.yaml
}