package fabricnetwork

import (
	"fabric-tool/src/config"
	"fabric-tool/src/utils"
	"fmt"
)

const createChannelPreImport = `#!/bin/bash

. scripts/envVar.sh
. scripts/utils.sh

CHANNEL_NAME="$1"
DELAY="$2"
MAX_RETRY="$3"
VERBOSE="$4"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${MAX_RETRY:="5"}
: ${VERBOSE:="false"}

if [ ! -d "channel-artifacts" ]; then
  mkdir channel-artifacts
fi

`
const createChannelTx = `
createChannelTx() {
  set -x
  configtxgen -profile OrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
  res=$?
  { set +x; } 2>/dev/null
  verifyResult $res "Failed to generate channel configuration transaction..."
}

`

const createChannel1 = `

  local rc=1
  local COUNTER=1
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
`

const createChannel2 = `		
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  verifyResult $res "Channel creation failed"
}

`

const joinChannel = `
joinChannel() {
  FABRIC_CFG_PATH=$PWD/config/
  HOST=$1
  PEER_NAME=$2
  PEER_PORT=$3
  setGlobals $HOST $PEER_NAME $PEER_PORT
    local rc=1
    local COUNTER=1

    while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    peer channel join -b $BLOCKFILE >&log.txt
    res=$?
    { set +x; } 2>/dev/null
      let rc=$res
      COUNTER=$(expr $COUNTER + 1)
    done
    cat log.txt
    verifyResult $res "After $MAX_RETRY attempts, ${PEER_NAME} has failed to join channel '$CHANNEL_NAME' "
}

`

const setAnchorPeer = `
setAnchorPeer() {
  HOST=$1
  PEER_NAME=$2
  PEER_PORT=$3
  docker exec cli ./scripts/setAnchorPeer.sh $HOST $CHANNEL_NAME $PEER_NAME $PEER_PORT
}

`

const createChannelMain = `
FABRIC_CFG_PATH=${PWD}/configtx-${CHANNEL_NAME}

infoln "Generating channel create transaction '${CHANNEL_NAME}.tx'"
createChannelTx

FABRIC_CFG_PATH=$PWD/config/
BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"

infoln "Creating channel ${CHANNEL_NAME}"
`

func GenerateCreateChannel(conf *config.Config, channelName string) (string, error) {
	res := createChannelPreImport
	res = res + createChannelTx

	res = res + "createChannel() {\n"
	res = res + "  setGlobals $1 $2 $3\n"
	res = res + createChannel1
	for _, org := range conf.Organizations {
		if org.Type == "orderOrg" {
			res = res + "    peer channel create -o localhost:" + org.Peers[0].Port + " -c $CHANNEL_NAME --ordererTLSHostnameOverride " + org.Peers[0].Name + " -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock $BLOCKFILE --tls --cafile $ORDERER_CA >&log.txt\n"
			break
		}
	}
	res = res + createChannel2
	res = res + joinChannel
	res = res + setAnchorPeer
	res = res + createChannelMain
	for _, org := range conf.Organizations {
		isChannelCreated := false
		for _, channel := range org.Channels {
			if org.Type == "peerOrg" && channel == channelName && !isChannelCreated {
				peers := org.CommittingPeers
				if len(peers) == 0 {
					peers = org.EndorsingPeers
				}
				host, err := utils.ExtractHost(peers[0].Name, 1)
				if err != nil {
					return res, fmt.Errorf("error occur in extracting host: %v", err)
				}
				res = res + "createChannel " + host + " " + peers[0].Name + " " + peers[0].Port + "\n\n"
				isChannelCreated = true
			}
		}
		if isChannelCreated {
			break
		}
	}
	res = res + "successln \"Channel '$CHANNEL_NAME' created\"\n\n"
	for _, org := range conf.Organizations {
		for _, channel := range org.Channels {
			if org.Type == "peerOrg" && channel == channelName {
				for _, peer := range org.CommittingPeers {
					host, err := utils.ExtractHost(peer.Name, 1)
					if err != nil {
						return res, fmt.Errorf("error occur in extracting host: %v", err)
					}
					res = res + "infoln \"Joining " + peer.Name + " to the channel...\"\n"
					res = res + "joinChannel " + host + " " + peer.Name + " " + peer.Port + "\n\n"
				}
				for _, peer := range org.EndorsingPeers {
					host, err := utils.ExtractHost(peer.Name, 1)
					if err != nil {
						return res, fmt.Errorf("error occur in extracting host: %v", err)
					}
					res = res + "infoln \"Joining " + peer.Name + " to the channel...\"\n"
					res = res + "joinChannel " + host + " " + peer.Name + " " + peer.Port + "\n\n"
				}
			}
		}
	}
	for _, org := range conf.Organizations {
		for _, channel := range org.Channels {
			if org.Type == "peerOrg" && channel == channelName {
				isAnchorPeerBeSet := false
				for _, peer := range org.CommittingPeers {
					if isAnchorPeerBeSet {
						break
					}
					host, err := utils.ExtractHost(peer.Name, 1)
					if err != nil {
						return res, fmt.Errorf("error occur in extracting host: %v", err)
					}
					res = res + "infoln \"Setting anchor peer for " + org.Name + "...\"\n"
					res = res + "setAnchorPeer " + host + " " + peer.Name + " " + peer.Port + "\n\n"
					isAnchorPeerBeSet = true
				}
				for _, peer := range org.EndorsingPeers {
					if isAnchorPeerBeSet {
						break
					}
					host, err := utils.ExtractHost(peer.Name, 1)
					if err != nil {
						return res, fmt.Errorf("error occur in extracting host: %v", err)
					}
					res = res + "infoln \"Setting anchor peer for " + org.Name + "...\"\n"
					res = res + "setAnchorPeer " + host + " " + peer.Name + " " + peer.Port + "\n\n"
					isAnchorPeerBeSet = true
				}
			}
		}
	}
	res = res + "successln \"Channel '$CHANNEL_NAME' joined\""
	return res, nil
}
