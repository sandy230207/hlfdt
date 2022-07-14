#!/bin/bash

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


createChannelTx() {
  set -x
  configtxgen -profile OrgsChannel -outputCreateChannelTx ./channel-artifacts/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME
  res=$?
  { set +x; } 2>/dev/null
  verifyResult $res "Failed to generate channel configuration transaction..."
}

createChannel() {
  setGlobals org1.example.com peer1.org1.example.com 7051

  local rc=1
  local COUNTER=1
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    peer channel create -o localhost:7050 -c $CHANNEL_NAME --ordererTLSHostnameOverride orderer.example.com -f ./channel-artifacts/${CHANNEL_NAME}.tx --outputBlock $BLOCKFILE --tls --cafile $ORDERER_CA >&log.txt
		
    res=$?
    { set +x; } 2>/dev/null
    let rc=$res
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  verifyResult $res "Channel creation failed"
}


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


setAnchorPeer() {
  HOST=$1
  PEER_NAME=$2
  PEER_PORT=$3
  docker exec cli ./scripts/setAnchorPeer.sh $HOST $CHANNEL_NAME $PEER_NAME $PEER_PORT
}


FABRIC_CFG_PATH=${PWD}/configtx

infoln "Generating channel create transaction '${CHANNEL_NAME}.tx'"
createChannelTx

FABRIC_CFG_PATH=$PWD/config/
BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"

infoln "Creating channel ${CHANNEL_NAME}"
createChannel org1.example.com peer1.org1.example.com 7051

successln "Channel '$CHANNEL_NAME' created"

infoln "Joining peer1.org1.example.com to the channel..."
joinChannel org1.example.com peer1.org1.example.com 7051

infoln "Joining peer2.org1.example.com to the channel..."
joinChannel org1.example.com peer2.org1.example.com 7151

infoln "Joining peer3.org1.example.com to the channel..."
joinChannel org1.example.com peer3.org1.example.com 7251

infoln "Joining peer1.org2.example.com to the channel..."
joinChannel org2.example.com peer1.org2.example.com 9051

infoln "Joining peer2.org2.example.com to the channel..."
joinChannel org2.example.com peer2.org2.example.com 9151

infoln "Joining peer3.org2.example.com to the channel..."
joinChannel org2.example.com peer3.org2.example.com 9251

infoln "Setting anchor peer for org1..."
setAnchorPeer org1.example.com peer1.org1.example.com 7051

infoln "Setting anchor peer for org2..."
setAnchorPeer org2.example.com peer1.org2.example.com 9051

successln "Channel '$CHANNEL_NAME' joined"