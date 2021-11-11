#!/bin/bash

. scripts/envVar.sh
. scripts/configUpdate.sh


createAnchorPeerUpdate() {
  echo CHANNEL_NAME=$CHANNEL_NAME  CORE_PEER_LOCALMSPID=${CORE_PEER_LOCALMSPID} PEER_NAME=$PEER_NAME PEER_PORT=$PEER_PORT HOST=${HOST}
  infoln "Fetching channel config for channel $CHANNEL_NAME"
  fetchChannelConfig $HOST $CHANNEL_NAME ${CORE_PEER_LOCALMSPID}config.json $PEER_NAME $PEER_PORT
  
  infoln "Generating anchor peer update transaction for $HOST on channel $CHANNEL_NAME"

  set -x
  jq '.channel_group.groups.Application.groups.'${CORE_PEER_LOCALMSPID}'.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "'${PEER_NAME}'","port": '${PEER_PORT}'}]},"version": "0"}}' ${CORE_PEER_LOCALMSPID}config.json > ${CORE_PEER_LOCALMSPID}modified_config.json
  { set +x; } 2>/dev/null
  
  createConfigUpdate ${CHANNEL_NAME} ${CORE_PEER_LOCALMSPID}config.json ${CORE_PEER_LOCALMSPID}modified_config.json ${CORE_PEER_LOCALMSPID}anchors.tx
}

updateAnchorPeer() {
peer channel update -o orderer.example.com:7050 --ordererTLSHostnameOverride orderer.example.com -c $CHANNEL_NAME -f ${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile $ORDERER_CA >&log.txt


  res=$?
  cat log.txt
  verifyResult $res "Anchor peer update failed"
  successln "Anchor peer set for org '$CORE_PEER_LOCALMSPID' on channel '$CHANNEL_NAME'"
}

ORG=$1
HOST=$1
CHANNEL_NAME=$2
PEER_NAME=$3
PEER_PORT=$4

setGlobalsCLI $ORG $PEER_NAME $PEER_PORT

createAnchorPeerUpdate

updateAnchorPeer 
