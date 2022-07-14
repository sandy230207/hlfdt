#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG_NAME}/$1/" \
        -e "s/\${PEER_PORT}/$2/" \
        -e "s/\${CA_PEER_PORT}/$3/" \
        -e "s#\${PEER_PEM}#$PP#" \
        -e "s#\${CA_PEM}#$CP#" \
        -e "s/\${CA_PEER_PORT}/$3/" \
        -e "s#\${PEER_NAME}#$6#" \
        -e "s#\${CA_PEER_NAME}#$7#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG_NAME}/$1/" \
        -e "s/\${PEER_PORT}/$2/" \
        -e "s/\${CA_PEER_PORT}/$3/" \
        -e "s#\${PEER_PEM}#$PP#" \
        -e "s#\${PEER_NAME}#$6#" \
        -e "s#\${CA_PEER_NAME}#$7#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

export ORG_NAME=org1
export CA_PEER_PORT=7054
export PEER_PEM=organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
export CA_PEM=organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
export CA_PEER_NAME=ca.org1.example.com
echo "$(json_ccp $ORG_NAME 7051 $CA_PEER_PORT $PEER_PEM $CA_PEM peer1.org1.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org1.example.com/connection-org1-peer1-org1-example-com.json
echo "$(yaml_ccp $ORG_NAME 7051 $CA_PEER_PORT $PEER_PEM $CA_PEM peer1.org1.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org1.example.com/connection-org1-peer1-org1-example-com.yaml

echo "$(json_ccp $ORG_NAME 7151 $CA_PEER_PORT $PEER_PEM $CA_PEM peer2.org1.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org1.example.com/connection-org1-peer2-org1-example-com.json
echo "$(yaml_ccp $ORG_NAME 7151 $CA_PEER_PORT $PEER_PEM $CA_PEM peer2.org1.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org1.example.com/connection-org1-peer2-org1-example-com.yaml

echo "$(json_ccp $ORG_NAME 7251 $CA_PEER_PORT $PEER_PEM $CA_PEM peer3.org1.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org1.example.com/connection-org1-peer3-org1-example-com.json
echo "$(yaml_ccp $ORG_NAME 7251 $CA_PEER_PORT $PEER_PEM $CA_PEM peer3.org1.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org1.example.com/connection-org1-peer3-org1-example-com.yaml

export ORG_NAME=org2
export CA_PEER_PORT=8054
export PEER_PEM=organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
export CA_PEM=organizations/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem
export CA_PEER_NAME=ca.org2.example.com
echo "$(json_ccp $ORG_NAME 9051 $CA_PEER_PORT $PEER_PEM $CA_PEM peer1.org2.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org2.example.com/connection-org2-peer1-org2-example-com.json
echo "$(yaml_ccp $ORG_NAME 9051 $CA_PEER_PORT $PEER_PEM $CA_PEM peer1.org2.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org2.example.com/connection-org2-peer1-org2-example-com.yaml

echo "$(json_ccp $ORG_NAME 9151 $CA_PEER_PORT $PEER_PEM $CA_PEM peer2.org2.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org2.example.com/connection-org2-peer2-org2-example-com.json
echo "$(yaml_ccp $ORG_NAME 9151 $CA_PEER_PORT $PEER_PEM $CA_PEM peer2.org2.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org2.example.com/connection-org2-peer2-org2-example-com.yaml

echo "$(json_ccp $ORG_NAME 9251 $CA_PEER_PORT $PEER_PEM $CA_PEM peer3.org2.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org2.example.com/connection-org2-peer3-org2-example-com.json
echo "$(yaml_ccp $ORG_NAME 9251 $CA_PEER_PORT $PEER_PEM $CA_PEM peer3.org2.example.com $CA_PEER_NAME)" > organizations/peerOrganizations/org2.example.com/connection-org2-peer3-org2-example-com.yaml

