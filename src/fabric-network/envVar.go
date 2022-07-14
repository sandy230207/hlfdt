package fabricnetwork

import (
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"fmt"
)

const envPreImport = `#!/bin/bash

. scripts/utils.sh

export CORE_PEER_TLS_ENABLED=true

`

const setGlobals1 = `
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
`

const setGlobalsCLI = `
setGlobalsCLI() {
  setGlobals $1 $2 $3
  PEER_NAME=$2
  PEER_PORT=$3
  export CORE_PEER_ADDRESS=${PEER_NAME}:${PEER_PORT}
}

`

const parsePeerConnectionParameters = `
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

`

const verifyResult = `
verifyResult() {
  if [ $1 -ne 0 ]; then
      fatalln "$2"
  fi
}
`

func GenerateEnvVar(conf *config.Config) (string, error) {
	res := envPreImport
	for _, org := range conf.Organizations {
		if org.Type == "orderOrg" {
			host, err := utils.ExtractHost(org.Peers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + "export ORDERER_CA=${PWD}/organizations/ordererOrganizations/" + host + "/orderers/" + org.Peers[0].Name + "/msp/tlscacerts/tlsca." + host + "-cert.pem\n"
			break
		}
	}
	res = res + setGlobals1
	res = res + setGlobalsCLI
	res = res + parsePeerConnectionParameters
	res = res + verifyResult

	return res, nil
}
