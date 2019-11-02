FABRIC_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer

echo $FABRIC_PATH

export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=$FABRIC_PATH/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer1.org1.example.com:7051

peer chaincode invoke -o orderer0.example.com:7050 --tls --cafile $FABRIC_PATH/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C ch1 -n fabticket -c '{"function":"queryUserTickets","Args":["owen1994"]}'
