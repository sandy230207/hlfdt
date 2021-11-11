package fabricnetwork

import (
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"fmt"
)

const pre = `#!/bin/bash

source scripts/utils.sh

CHANNEL_NAME=${1:-"mychannel"}
CC_NAME=${2}
CC_SRC_PATH=${3}
CC_SRC_LANGUAGE=${4}
CC_VERSION=${5:-"1.0"}
CC_SEQUENCE=${6:-"1"}
CC_INIT_FCN=${7:-"NA"}
CC_END_POLICY=${8:-"NA"}
CC_COLL_CONFIG=${9:-"NA"}
DELAY=${10:-"3"}
MAX_RETRY=${11:-"5"}
VERBOSE=${12:-"false"}

println "executing with the following"
println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
println "- CC_SRC_LANGUAGE: ${C_GREEN}${CC_SRC_LANGUAGE}${C_RESET}"
println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
println "- CC_END_POLICY: ${C_GREEN}${CC_END_POLICY}${C_RESET}"
println "- CC_COLL_CONFIG: ${C_GREEN}${CC_COLL_CONFIG}${C_RESET}"
println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
println "- DELAY: ${C_GREEN}${DELAY}${C_RESET}"
println "- MAX_RETRY: ${C_GREEN}${MAX_RETRY}${C_RESET}"
println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"

FABRIC_CFG_PATH=$PWD/config/

#User has not provided a name
if [ -z "$CC_NAME" ] || [ "$CC_NAME" = "NA" ]; then
  fatalln "No chaincode name was provided. Valid call example: ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go"

# User has not provided a path
elif [ -z "$CC_SRC_PATH" ] || [ "$CC_SRC_PATH" = "NA" ]; then
  fatalln "No chaincode path was provided. Valid call example: ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go"

# User has not provided a language
elif [ -z "$CC_SRC_LANGUAGE" ] || [ "$CC_SRC_LANGUAGE" = "NA" ]; then
  fatalln "No chaincode language was provided. Valid call example: ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go"

## Make sure that the path to the chaincode exists
elif [ ! -d "$CC_SRC_PATH" ]; then
  fatalln "Path to chaincode does not exist. Please provide different path."
fi

CC_SRC_LANGUAGE=$(echo "$CC_SRC_LANGUAGE" | tr [:upper:] [:lower:])

# do some language specific preparation to the chaincode before packaging
if [ "$CC_SRC_LANGUAGE" = "go" ]; then
  CC_RUNTIME_LANGUAGE=golang

  infoln "Vendoring Go dependencies at $CC_SRC_PATH"
  pushd $CC_SRC_PATH
  GO111MODULE=on go mod vendor
  popd
  successln "Finished vendoring Go dependencies"

elif [ "$CC_SRC_LANGUAGE" = "java" ]; then
  CC_RUNTIME_LANGUAGE=java

  infoln "Compiling Java code..."
  pushd $CC_SRC_PATH
  ./gradlew installDist
  popd
  successln "Finished compiling Java code"
  CC_SRC_PATH=$CC_SRC_PATH/build/install/$CC_NAME

elif [ "$CC_SRC_LANGUAGE" = "javascript" ]; then
  CC_RUNTIME_LANGUAGE=node

elif [ "$CC_SRC_LANGUAGE" = "typescript" ]; then
  CC_RUNTIME_LANGUAGE=node

  infoln "Compiling TypeScript code into JavaScript..."
  pushd $CC_SRC_PATH
  npm install
  npm run build
  popd
  successln "Finished compiling TypeScript code into JavaScript"

else
  fatalln "The chaincode language ${CC_SRC_LANGUAGE} is not supported by this script. Supported chaincode languages are: go, java, javascript, and typescript"
  exit 1
fi

INIT_REQUIRED="--init-required"
# check if the init fcn should be called
if [ "$CC_INIT_FCN" = "NA" ]; then
  INIT_REQUIRED=""
fi

if [ "$CC_END_POLICY" = "NA" ]; then
  CC_END_POLICY=""
else
  CC_END_POLICY="--signature-policy $CC_END_POLICY"
fi

if [ "$CC_COLL_CONFIG" = "NA" ]; then
  CC_COLL_CONFIG=""
else
  CC_COLL_CONFIG="--collections-config $CC_COLL_CONFIG"
fi

# import utils
. scripts/envVar.sh

`

const packageChaincode = `
packageChaincode() {
  set -x
  peer lifecycle chaincode package ${CC_NAME}.tar.gz --path ${CC_SRC_PATH} --lang ${CC_RUNTIME_LANGUAGE} --label ${CC_NAME}_${CC_VERSION} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode packaging has failed"
  successln "Chaincode is packaged"
}

`

const installChaincode = `
installChaincode() {
  HOST=$1
  ORG=$HOST
  PEER_NAME=$2
  PEER_PORT=$3
  setGlobals $HOST $PEER_NAME $PEER_PORT
  set -x
  peer lifecycle chaincode install ${CC_NAME}.tar.gz >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode installation on ${PEER_NAME} has failed"
  successln "Chaincode is installed on ${PEER_NAME}"
}
`

const queryInstalled = `
queryInstalled() {
  HOST=$1
  ORG=$HOST
  PEER_NAME=$2
  PEER_PORT=$3
  setGlobals $HOST $PEER_NAME $PEER_PORT
  set -x
  peer lifecycle chaincode queryinstalled >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
  verifyResult $res "Query installed on ${PEER_NAME} has failed"
  successln "Query installed successful on ${PEER_NAME} on channel"
}
`

const approveForMyOrg1 = `
approveForMyOrg() {
  HOST=$1
  ORG=$HOST
  PEER_NAME=$2
  PEER_PORT=$3
  setGlobals $HOST $PEER_NAME $PEER_PORT
  set -x

`

const approveForMyOrg2 = `
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode definition approved on ${PEER_NAME} on channel '$CHANNEL_NAME' failed"
  successln "Chaincode definition approved on ${PEER_NAME} on channel '$CHANNEL_NAME'"
}

`

const checkCommitReadiness = `
checkCommitReadiness() {
  HOST=$1
  ORG=$HOST
  shift 1
  PEER_NAME=$2
  PEER_PORT=$3
  setGlobals $HOST $PEER_NAME $PEER_PORT
  infoln "Checking the commit readiness of the chaincode definition on ${PEER_NAME} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to check the commit readiness of the chaincode definition on ${PEER_NAME}, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} --output json >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=0
    for var in "$@"; do
      grep "$var" log.txt &>/dev/null || let rc=1
    done
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
      infoln "Checking the commit readiness of the chaincode definition successful on ${PEER_NAME} on channel '$CHANNEL_NAME'"
  else
      fatalln "After $MAX_RETRY attempts, Check commit readiness result on ${PEER_NAME} is INVALID!"
  fi
}

`

const commitChaincodeDefinition1 = `
commitChaincodeDefinition() {
  parsePeerConnectionParameters $@
  res=$?
  verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
  set -x

`
const commitChaincodeDefinition2 = `
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode definition commit failed on ${PEER_NAME} on channel '$CHANNEL_NAME' failed"
   successln "Chaincode definition committed on channel '$CHANNEL_NAME'"
}

`

const queryCommitted = `
queryCommitted() {
  HOST=$1
  ORG=$HOST
  PEER_NAME=$2
  PEER_PORT=$3
  setGlobals $HOST $PEER_NAME $PEER_PORT
  EXPECTED_RESULT="Version: ${CC_VERSION}, Sequence: ${CC_SEQUENCE}, Endorsement Plugin: escc, Validation Plugin: vscc"
  infoln "Querying chaincode definition on ${PEER_NAME} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to Query committed status on ${PEER_NAME}, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME} >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    test $res -eq 0 && VALUE=$(cat log.txt | grep -o '^Version: '$CC_VERSION', Sequence: [0-9]*, Endorsement Plugin: escc, Validation Plugin: vscc')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
      successln "Query chaincode definition successful on ${PEER_NAME} on channel '$CHANNEL_NAME'"
  else
      fatalln "After $MAX_RETRY attempts, Query chaincode definition result on ${PEER_NAME} is INVALID!"
  fi
}
`

const chaincodeInvokeInit1 = `
chaincodeInvokeInit() {
  parsePeerConnectionParameters $@
  res=$?
  verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
  set -x
  fcn_call='{"function":"'${CC_INIT_FCN}'","Args":[]}'
  infoln "invoke fcn call:${fcn_call}"

`

const chaincodeInvokeInit2 = `
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Invoke execution on $PEERS failed "
  successln "Invoke transaction successful on $PEERS on channel '$CHANNEL_NAME'"
}

`

const chaincodeQuery = `
chaincodeQuery() {
  HOST=$1
  ORG=$HOST
  PEER_NAME=$2
  PEER_PORT=$3
  setGlobals $HOST $PEER_NAME $PEER_PORT
  infoln "Querying on ${PEER_NAME} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to Query ${PEER_NAME}, Retry after $DELAY seconds."
    set -x
    peer chaincode query -C $CHANNEL_NAME -n ${CC_NAME} -c '{"Args":["queryAllCars"]}' >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
      successln "Query successful on ${PEER_NAME} on channel '$CHANNEL_NAME'"
  else
      fatalln "After $MAX_RETRY attempts, Query result on ${PEER_NAME} is INVALID!"
  fi
}

`

func GenerateDeployCC(conf *config.Config) (string, error) {
	res := pre
	res = res + packageChaincode
	res = res + installChaincode
	res = res + queryInstalled
	res = res + approveForMyOrg1
	for _, org := range conf.Organizations {
		if org.Type == "orderOrg" {
			for _, peer := range org.Peers {
				res = res + "  peer lifecycle chaincode approveformyorg -o localhost:" + peer.Port + " --ordererTLSHostnameOverride " + peer.Name + " --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt\n"
			}
		}
	}
	res = res + approveForMyOrg2
	res = res + checkCommitReadiness
	res = res + commitChaincodeDefinition1
	for _, org := range conf.Organizations {
		if org.Type == "orderOrg" {
			for _, peer := range org.Peers {
				res = res + "  peer lifecycle chaincode commit -o localhost:" + peer.Port + " --ordererTLSHostnameOverride " + peer.Name + " --tls --cafile $ORDERER_CA --channelID $CHANNEL_NAME --name ${CC_NAME} $PEER_CONN_PARMS --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt\n"
			}
		}
	}
	res = res + commitChaincodeDefinition2
	res = res + queryCommitted
	res = res + chaincodeInvokeInit1
	for _, org := range conf.Organizations {
		if org.Type == "orderOrg" {
			for _, peer := range org.Peers {
				res = res + "  peer chaincode invoke -o localhost:" + peer.Port + " --ordererTLSHostnameOverride " + peer.Name + " --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} $PEER_CONN_PARMS --isInit -c ${fcn_call} >&log.txt\n"
			}
		}
	}
	res = res + chaincodeInvokeInit2
	res = res + chaincodeQuery
	res = res + "packageChaincode\n"

	for _, org := range conf.Organizations {
		for _, peer := range org.EndorsingPeers {
			host, err := utils.ExtractHost(peer.Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + "infoln \"Installing chaincode on " + peer.Name + "...\"\n"
			res = res + "installChaincode " + host + " " + peer.Name + " " + peer.Port + "\n\n"
			res = res + "queryInstalled " + host + " " + peer.Name + " " + peer.Port + "\n\n"
		}
		if len(org.EndorsingPeers) > 0 {
			host, err := utils.ExtractHost(org.EndorsingPeers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + "approveForMyOrg " + host + " " + org.EndorsingPeers[0].Name + " " + org.EndorsingPeers[0].Port + "\n\n"
		}
	}
	res = res + "commitChaincodeDefinition "
	for _, org := range conf.Organizations {
		if len(org.EndorsingPeers) > 0 {
			host, err := utils.ExtractHost(org.EndorsingPeers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + host + " " + org.EndorsingPeers[0].Name + " " + org.EndorsingPeers[0].Port + " "
		}
	}
	res = res + "\n\n"
	for _, org := range conf.Organizations {
		if len(org.EndorsingPeers) > 0 {
			host, err := utils.ExtractHost(org.EndorsingPeers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + "queryCommitted " + host + " " + org.EndorsingPeers[0].Name + " " + org.EndorsingPeers[0].Port + "\n\n"
		}
	}
	res = res + "if [ \"$CC_INIT_FCN\" = \"NA\" ]; then\n"
	res = res + "  infoln \"Chaincode initialization is not required\"\n"
	res = res + "else\n"
	res = res + "  chaincodeInvokeInit "
	for _, org := range conf.Organizations {
		if len(org.EndorsingPeers) > 0 {
			host, err := utils.ExtractHost(org.EndorsingPeers[0].Name, 1)
			if err != nil {
				return res, fmt.Errorf("error occur in extracting host: %v", err)
			}
			res = res + host + " " + org.EndorsingPeers[0].Name + " " + org.EndorsingPeers[0].Port + " "
		}
	}
	res = res + "\nfi\n\n"
	res = res + "exit 0"

	return res, nil
}
