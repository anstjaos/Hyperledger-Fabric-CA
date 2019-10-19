#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $8)
    local CP=$(one_line_pem $9)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0IP}/$2/" \
        -e "s/\${P0PORT}/$3/" \
        -e "s/\${P1IP}/$4/" \
        -e "s/\${P1PORT}/$5/" \
        -e "s/\${CAIP}/$6/" \
        -e "s/\${CAPORT}/$7/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ccp-template.json 
}

function yaml_ccp {
    local PP=$(one_line_pem $8)
    local CP=$(one_line_pem $9)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0IP}/$2/" \
        -e "s/\${P0PORT}/$3/" \
        -e "s/\${P1IP}/$4/" \
        -e "s/\${P1PORT}/$5/" \
        -e "s/\${CAIP}/$6/" \
        -e "s/\${CAPORT}/$7/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ccp-template.yaml | sed -e $'s/\\\\n/\\\n        /g'
}

ORG=0
P0IP=15.164.140.245
P0PORT=7051
P1IP=13.209.61.41
P1PORT=7051
CAIP=13.209.128.90
CAPORT=7054
PEERPEM=crypto-config/peerOrganizations/org0.example.com/tlsca/tlsca.org0.example.com-cert.pem
CAPEM=crypto-config/peerOrganizations/org0.example.com/ca/ca.org0.example.com-cert.pem

echo "$(json_ccp $ORG $P0IP $P0PORT $P1IP $P1PORT $CAIP $CAPORT $PEERPEM $CAPEM)" > connection-org0.json
echo "$(yaml_ccp $ORG $P0IP $P0PORT $P1IP $P1PORT $CAIP $CAPORT $PEERPEM $CAPEM)" > connection-org0.yaml

ORG=1
POIP=15.164.153.125
P0PORT=7051
P1IP=15.164.121.13
P1PORT=7051
CAIP=13.209.128.90
CAPORT=8054
PEERPEM=crypto-config/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
CAPEM=crypto-config/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem

echo "$(json_ccp $ORG $P0IP $P0PORT $P1IP $P1PORT $CAIP $CAPORT $PEERPEM $CAPEM)" > connection-org1.json
echo "$(yaml_ccp $ORG $P0IP $P0PORT $P1IP $P1PORT $CAIP $CAPORT $PEERPEM $CAPEM)" > connection-org1.yaml
