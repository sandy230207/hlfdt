#!/bin/bash

. scripts/utils.sh

export CORE_PEER_TLS_ENABLED=true

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

setGlobals() {
  HOST=$1
  arrHost=(${HOST//./ })
  ORG=${arrHost[0]}
  PEER_NAME=$2
  PEER_PORT=$3
  infoln "Using organization peer ${PEER_NAME}"
  export CORE_PEER_LOCALMSPID=${ORG}
  export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/${HOST}/peers/${PEER_NAME}/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/${HOST}/users/Admin@${HOST}/msp
  export CORE_PEER_ADDRESS=localhost:${PEER_PORT}

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

setGlobalsCLI() {
  setGlobals $1 $2 $3
  PEER_NAME=$2
  PEER_PORT=$3
  export CORE_PEER_ADDRESS=${PEER_NAME}:${PEER_PORT}
}


parsePeerConnectionParameters() {
  PEER_CONN_PARMS=""
  PEERS=""
  while [ "$#" -gt 0 ]; do
    setGlobals $1 $2 $3
    HOST=$1
    PEER_NAME=$2
    PEER_PORT=$3

    PEERS="$PEERS $PEER_NAME"
    PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $CORE_PEER_ADDRESS"

    TLSINFO=$(eval echo "--tlsRootCertFiles \${PWD}/organizations/peerOrganizations/${HOST}/peers/${PEER_NAME}/tls/ca.crt")
    PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"

    shift
    shift
    shift
  done

  PEERS="$(echo -e "$PEERS" | sed -e 's/^[[:space:]]*//')"
  }


verifyResult() {
  if [ $1 -ne 0 ]; then
      fatalln "$2"
  fi
}
