#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${PEERPORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEER}#$6#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${PEERPORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${PEER}#$6#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORG=1

P0PORT=7051
P1PORT=7151
P2PORT=7251
P3PORT=7351
P4PORT=7451
P5PORT=7551
P6PORT=7651
P7PORT=7751
P8PORT=7851
P9PORT=7951

PEER0=0
PEER1=1
PEER2=2
PEER3=3
PEER4=4
PEER5=5
PEER6=6
PEER7=7
PEER8=8
PEER9=9

CAPORT=7054
PEERPEM=organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
CAPEM=organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $PEER0)" > organizations/peerOrganizations/org1.example.com/connection-org1-0.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $PEER0)" > organizations/peerOrganizations/org1.example.com/connection-org1-0.yaml

echo "$(json_ccp $ORG $P1PORT $CAPORT $PEERPEM $CAPEM $PEER1)" > organizations/peerOrganizations/org1.example.com/connection-org1-1.json
echo "$(yaml_ccp $ORG $P1PORT $CAPORT $PEERPEM $CAPEM $PEER1)" > organizations/peerOrganizations/org1.example.com/connection-org1-1.yaml

echo "$(json_ccp $ORG $P2PORT $CAPORT $PEERPEM $CAPEM $PEER2)" > organizations/peerOrganizations/org1.example.com/connection-org1-2.json
echo "$(yaml_ccp $ORG $P2PORT $CAPORT $PEERPEM $CAPEM $PEER2)" > organizations/peerOrganizations/org1.example.com/connection-org1-2.yaml

echo "$(json_ccp $ORG $P3PORT $CAPORT $PEERPEM $CAPEM $PEER3)" > organizations/peerOrganizations/org1.example.com/connection-org1-3.json
echo "$(yaml_ccp $ORG $P3PORT $CAPORT $PEERPEM $CAPEM $PEER3)" > organizations/peerOrganizations/org1.example.com/connection-org1-3.yaml

echo "$(json_ccp $ORG $P4PORT $CAPORT $PEERPEM $CAPEM $PEER4)" > organizations/peerOrganizations/org1.example.com/connection-org1-4.json
echo "$(yaml_ccp $ORG $P4PORT $CAPORT $PEERPEM $CAPEM $PEER4)" > organizations/peerOrganizations/org1.example.com/connection-org1-4.yaml

echo "$(json_ccp $ORG $P5PORT $CAPORT $PEERPEM $CAPEM $PEER5)" > organizations/peerOrganizations/org1.example.com/connection-org1-5.json
echo "$(yaml_ccp $ORG $P5PORT $CAPORT $PEERPEM $CAPEM $PEER5)" > organizations/peerOrganizations/org1.example.com/connection-org1-5.yaml

echo "$(json_ccp $ORG $P6PORT $CAPORT $PEERPEM $CAPEM $PEER6)" > organizations/peerOrganizations/org1.example.com/connection-org1-6.json
echo "$(yaml_ccp $ORG $P6PORT $CAPORT $PEERPEM $CAPEM $PEER6)" > organizations/peerOrganizations/org1.example.com/connection-org1-6.yaml

echo "$(json_ccp $ORG $P7PORT $CAPORT $PEERPEM $CAPEM $PEER7)" > organizations/peerOrganizations/org1.example.com/connection-org1-7.json
echo "$(yaml_ccp $ORG $P7PORT $CAPORT $PEERPEM $CAPEM $PEER7)" > organizations/peerOrganizations/org1.example.com/connection-org1-7.yaml

echo "$(json_ccp $ORG $P8PORT $CAPORT $PEERPEM $CAPEM $PEER8)" > organizations/peerOrganizations/org1.example.com/connection-org1-8.json
echo "$(yaml_ccp $ORG $P8PORT $CAPORT $PEERPEM $CAPEM $PEER8)" > organizations/peerOrganizations/org1.example.com/connection-org1-8.yaml

echo "$(json_ccp $ORG $P9PORT $CAPORT $PEERPEM $CAPEM $PEER9)" > organizations/peerOrganizations/org1.example.com/connection-org1-9.json
echo "$(yaml_ccp $ORG $P9PORT $CAPORT $PEERPEM $CAPEM $PEER9)" > organizations/peerOrganizations/org1.example.com/connection-org1-9.yaml

# ORG=2
# P0PORT=9051
# CAPORT=8054
# PEERPEM=organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
# CAPEM=organizations/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem

# echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org2.example.com/connection-org2.json
# echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org2.example.com/connection-org2.yaml
