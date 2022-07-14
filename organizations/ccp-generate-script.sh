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

