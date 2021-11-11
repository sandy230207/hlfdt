package fabricnetwork

import (
	"fabric-tool/src/config"
)

const createOrgs = `#!/bin/bash
export PATH=${PWD}/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx
export VERBOSE=false

. scripts/utils.sh

function clearContainers() {
  CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    infoln "No containers available for deletion"
  else
    docker rm -f $CONTAINER_IDS
  fi
}

function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*/) {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    infoln "No images available for deletion"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

NONWORKING_VERSIONS="^1\.0\. ^1\.1\. ^1\.2\. ^1\.3\. ^1\.4\."

function checkPrereqs() {
  peer version > /dev/null 2>&1
  if [[ $? -ne 0 || ! -d "./config" ]]; then
    errorln "Peer binary and configuration files not found.."
    errorln
    errorln "Follow the instructions in the Fabric docs to install the Fabric Binaries:"
    errorln "https://hyperledger-fabric.readthedocs.io/en/latest/install.html"
    exit 1
  fi
  
  LOCAL_VERSION=$(peer version | sed -ne 's/ Version: //p')
  DOCKER_IMAGE_VERSION=$(docker run --rm hyperledger/fabric-tools:$IMAGETAG peer version | sed -ne 's/ Version: //p' | head -1)

  infoln "LOCAL_VERSION=$LOCAL_VERSION"
  infoln "DOCKER_IMAGE_VERSION=$DOCKER_IMAGE_VERSION"

  if [ "$LOCAL_VERSION" != "$DOCKER_IMAGE_VERSION" ]; then
    warnln "Local fabric binaries and docker images are out of  sync. This may cause problems."
  fi

  for UNSUPPORTED_VERSION in $NONWORKING_VERSIONS; do
    infoln "$LOCAL_VERSION" | grep -q $UNSUPPORTED_VERSION
    if [ $? -eq 0 ]; then
      fatalln "Local Fabric binary version of $LOCAL_VERSION does not match the versions supported by the test network."
    fi

    infoln "$DOCKER_IMAGE_VERSION" | grep -q $UNSUPPORTED_VERSION
    if [ $? -eq 0 ]; then
      fatalln "Fabric Docker image version of $DOCKER_IMAGE_VERSION does not match the versions supported by the test network."
    fi
  done

  if [ "$CRYPTO" == "Certificate Authorities" ]; then

    fabric-ca-client version > /dev/null 2>&1
    if [[ $? -ne 0 ]]; then
      errorln "fabric-ca-client binary not found.."
      errorln
      errorln "Follow the instructions in the Fabric docs to install the Fabric Binaries:"
      errorln "https://hyperledger-fabric.readthedocs.io/en/latest/install.html"
      exit 1
    fi
    CA_LOCAL_VERSION=$(fabric-ca-client version | sed -ne 's/ Version: //p')
    CA_DOCKER_IMAGE_VERSION=$(docker run --rm hyperledger/fabric-ca:$CA_IMAGETAG fabric-ca-client version | sed -ne 's/ Version: //p' | head -1)
    infoln "CA_LOCAL_VERSION=$CA_LOCAL_VERSION"
    infoln "CA_DOCKER_IMAGE_VERSION=$CA_DOCKER_IMAGE_VERSION"

    if [ "$CA_LOCAL_VERSION" != "$CA_DOCKER_IMAGE_VERSION" ]; then
      warnln "Local fabric-ca binaries and docker images are out of sync. This may cause problems."
    fi
  fi
}

function createOrgs() {
  if [ -d "organizations/peerOrganizations" ]; then
    rm -Rf organizations/peerOrganizations && rm -Rf organizations/ordererOrganizations
  fi

  if [ "$CRYPTO" == "Certificate Authorities" ]; then
    infoln "Generating certificates using Fabric CA"

    IMAGE_TAG=${CA_IMAGETAG} docker-compose -f $COMPOSE_FILE_CA up -d 2>&1

    . organizations/fabric-ca/registerEnroll.sh

  while :
    do
    if [ ! -f "organizations/fabric-ca/org1/tls-cert.pem" ]; then
      sleep 1
    else
      break
    fi
    done

`

const createConsortium = `
function createConsortium() {
  which configtxgen
  if [ "$?" -ne 0 ]; then
    fatalln "configtxgen tool not found."
  fi
  
  infoln "Generating Orderer Genesis block"
  
  set -x

`

const networkUp = `
function networkUp() {
  checkPrereqs
  #docker volume rm $(docker volume ls)
  
  if [ ! -d "organizations/peerOrganizations" ]; then
    createOrgs
    createConsortium
  fi
  
  COMPOSE_FILES="-f ${COMPOSE_FILE_BASE}"
  
  if [ "${DATABASE}" == "couchdb" ]; then
    COMPOSE_FILES="${COMPOSE_FILES} -f ${COMPOSE_FILE_COUCH}"
  fi
  
  IMAGE_TAG=$IMAGETAG docker-compose ${COMPOSE_FILES} up -d 2>&1
  
  docker ps -a
  if [ $? -ne 0 ]; then
    fatalln "Unable to start network"
  fi
}

`

const createChannel = `
function createChannel() {  
  if [ ! -d "organizations/peerOrganizations" ]; then
    infoln "Bringing up network"
    networkUp
  fi
  
  scripts/createChannel.sh $CHANNEL_NAME $CLI_DELAY $MAX_RETRY $VERBOSE
}

`

const deployCC = `
function deployCC() {
  ./scripts/deployCC.sh $CHANNEL_NAME $CC_NAME $CC_SRC_PATH $CC_SRC_LANGUAGE $CC_VERSION $CC_SEQUENCE $CC_INIT_FCN $CC_END_POLICY $CC_COLL_CONFIG $CLI_DELAY $MAX_RETRY $VERBOSE
  
  if [ $? -ne 0 ]; then
    fatalln "Deploying chaincode failed"
  fi
}

`

const networkDown = `
function networkDown() {
  docker-compose -f $COMPOSE_FILE_BASE -f $COMPOSE_FILE_COUCH -f $COMPOSE_FILE_CA down --volumes --remove-orphans
  # docker-compose -f $COMPOSE_FILE_COUCH_ORG3 -f $COMPOSE_FILE_ORG3 down --volumes --remove-orphans
  if [ "$MODE" != "restart" ]; then
    clearContainers
    removeUnwantedImages
    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf system-genesis-block/*.block organizations/peerOrganizations organizations/ordererOrganizations'

`

const envVar = `
OS_ARCH=$(echo "$(uname -s | tr '[:upper:]' '[:lower:]' | sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
CRYPTO="cryptogen"
MAX_RETRY=5
CLI_DELAY=3
CHANNEL_NAME="mychannel"
CC_NAME="NA"
CC_SRC_PATH="NA"
CC_END_POLICY="NA"
CC_COLL_CONFIG="NA"
CC_INIT_FCN="NA"
COMPOSE_FILE_BASE=docker/docker-compose-test-net.yaml
COMPOSE_FILE_COUCH=docker/docker-compose-couch.yaml
COMPOSE_FILE_CA=docker/docker-compose-ca.yaml

CC_SRC_LANGUAGE="NA"
CC_VERSION="1.0"
CC_SEQUENCE=1
IMAGETAG="latest"
CA_IMAGETAG="latest"
DATABASE="leveldb"

`

const parseMode = `
if [[ $# -lt 1 ]] ; then
  printHelp
  exit 0
else
  MODE=$1
  shift
fi

`
const parseCreateChannelSubcommand = `
if [[ $# -ge 1 ]] ; then
  key="$1"
  if [[ "$key" == "createChannel" ]]; then
    export MODE="createChannel"
    shift
  fi
fi

`

const parseFlags = `
while [[ $# -ge 1 ]] ; do
  key="$1"
  case $key in
  -h )
    printHelp $MODE
    exit 0
    ;;
  -c )
    CHANNEL_NAME="$2"
    shift
    ;;
  -ca )
    CRYPTO="Certificate Authorities"
    ;;
  -r )
    MAX_RETRY="$2"
    shift
    ;;
  -d )
    CLI_DELAY="$2"
    shift
    ;;
  -s )
    DATABASE="$2"
    shift
    ;;
  -ccl )
    CC_SRC_LANGUAGE="$2"
    shift
    ;;
  -ccn )
    CC_NAME="$2"
    shift
    ;;
  -ccv )
    CC_VERSION="$2"
    shift
    ;;
  -ccs )
    CC_SEQUENCE="$2"
    shift
    ;;
  -ccp )
    CC_SRC_PATH="$2"
    shift
    ;;
  -ccep )
    CC_END_POLICY="$2"
    shift
    ;;
  -cccg )
    CC_COLL_CONFIG="$2"
    shift
    ;;
  -cci )
    CC_INIT_FCN="$2"
    shift
    ;;
  -i )
    IMAGETAG="$2"
    shift
    ;;
  -cai )
    CA_IMAGETAG="$2"
    shift
    ;;
  -verbose )
    VERBOSE=true
    shift
    ;;
  * )
    errorln "Unknown flag: $key"
    printHelp
    exit 1
    ;;
  esac
  shift
done

`

const cryptoMode = `
if [ ! -d "./organizations/peerOrganizations" ]; then
  CRYPTO_MODE="with crypto from '${CRYPTO}'"
else
  CRYPTO_MODE=""
fi

`

const parseOperationMode = `
if [ "$MODE" == "up" ]; then
  infoln "Starting nodes with CLI timeout of '${MAX_RETRY}' tries and CLI delay of '${CLI_DELAY}' seconds and using database '${DATABASE}' ${CRYPTO_MODE}"
elif [ "$MODE" == "createChannel" ]; then
  infoln "Creating channel '${CHANNEL_NAME}'."
  infoln "If network is not up, starting nodes with CLI timeout of '${MAX_RETRY}' tries and CLI delay of '${CLI_DELAY}' seconds and using database '${DATABASE} ${CRYPTO_MODE}"
elif [ "$MODE" == "down" ]; then
  infoln "Stopping network"
elif [ "$MODE" == "restart" ]; then
  infoln "Restarting network"
elif [ "$MODE" == "deployCC" ]; then
  infoln "deploying chaincode on channel '${CHANNEL_NAME}'"
else
  printHelp
  exit 1
fi

if [ "${MODE}" == "up" ]; then
  networkUp
elif [ "${MODE}" == "createChannel" ]; then
  createChannel
elif [ "${MODE}" == "deployCC" ]; then
  deployCC
elif [ "${MODE}" == "down" ]; then
  networkDown
else
  printHelp
  exit 1
fi
`

func GenerateNetwork(conf *config.Config) (string, error) {
	res := createOrgs
	for _, org := range conf.Organizations {
		res = res + "    infoln \"Creating " + org.Name + " Identities\"\n"
		res = res + "    create" + org.Name + "\n"
	}
	res = res + "  fi\n"
	res = res + "  infoln \"Generating CCP files\"\n"
	res = res + "  ./organizations/ccp-generate.sh\n"
	res = res + "}\n\n"

	res = res + createConsortium
	res = res + "  configtxgen -profile OrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block\n"
	res = res + "  res=$?\n"
	res = res + "  { set +x; } 2>/dev/null\n"
	res = res + "  if [ $res -ne 0 ]; then\n"
	res = res + "    fatalln \"Failed to generate orderer genesis block...\"\n"
	res = res + "  fi\n"
	res = res + "}\n\n"

	res = res + networkUp

	res = res + createChannel

	res = res + deployCC

	res = res + networkDown
	for _, org := range conf.Organizations {
		res = res + "    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf organizations/fabric-ca/" + org.Name + "/msp organizations/fabric-ca/" + org.Name + "/tls-cert.pem organizations/fabric-ca/" + org.Name + "/ca-cert.pem organizations/fabric-ca/" + org.Name + "/IssuerPublicKey organizations/fabric-ca/" + org.Name + "/IssuerRevocationPublicKey organizations/fabric-ca/" + org.Name + "/fabric-ca-server.db'\n"
	}
	res = res + "    docker run --rm -v $(pwd):/data busybox sh -c 'cd /data && rm -rf channel-artifacts log.txt *.tar.gz'\n"
	res = res + "  fi\n"
	res = res + "}\n"

	res = res + envVar

	res = res + parseMode
	res = res + parseCreateChannelSubcommand
	res = res + parseFlags
	res = res + cryptoMode
	res = res + parseOperationMode

	return res, nil
}
