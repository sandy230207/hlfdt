package fabricnetwork

import (
	"fabric-tool/src/config"
)

const updatePreImport = `#!/bin/bash

. scripts/envVar.sh
`
const fetchChannelConfig1 = `
fetchChannelConfig() {
  HOST=$1
  ORG=$HOST
  CHANNEL=$2
  OUTPUT=$3
  PEER_NAME=$4
  PEER_PORT=$5
  
  setGlobals $HOST $PEER_NAME $PEER_PORT
  
  infoln "Fetching the most recent configuration block for the channel"
  set -x

`

const fetchChannelConfig2 = `
  { set +x; } 2>/dev/null

  infoln "Decoding config block to JSON and isolating config to ${OUTPUT}"
  set -x
  configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >"${OUTPUT}"
  { set +x; } 2>/dev/null
}

`

const createConfigUpdate = `
createConfigUpdate() {
  CHANNEL=$1
  ORIGINAL=$2
  MODIFIED=$3
  OUTPUT=$4
  
  set -x
  configtxlator proto_encode --input "${ORIGINAL}" --type common.Config >original_config.pb
  configtxlator proto_encode --input "${MODIFIED}" --type common.Config >modified_config.pb
  configtxlator compute_update --channel_id "${CHANNEL}" --original original_config.pb --updated modified_config.pb >config_update.pb
  configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json
  configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >"${OUTPUT}"
  { set +x; } 2>/dev/null
}

`
const signConfigtxAsPeerOrg = `
signConfigtxAsPeerOrg() {
  HOST=$1
  ORG=$HOST
  CONFIGTXFILE=$2
  PEER_NAME=$3
  PEER_PORT=$4

  setGlobals $HOST $PEER_NAME $PEER_PORT
  set -x
  peer channel signconfigtx -f "${CONFIGTXFILE}"
  { set +x; } 2>/dev/null
}

`

func GenerateConfigUpdate(conf *config.Config) (string, error) {
	res := updatePreImport
	res = res + fetchChannelConfig1
	for _, org := range conf.Organizations {
		if org.Type == "orderOrg" {
			res = res + "  peer channel fetch config config_block.pb -o " + org.Peers[0].Name + ":" + org.Peers[0].Port + " --ordererTLSHostnameOverride " + org.Peers[0].Name + " -c $CHANNEL --tls --cafile $ORDERER_CA\n"
		}
	}
	res = res + fetchChannelConfig2
	res = res + createConfigUpdate
	res = res + signConfigtxAsPeerOrg
	return res, nil
}