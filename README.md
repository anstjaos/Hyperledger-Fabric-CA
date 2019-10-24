<<<<<<< HEAD
<<<<<<< HEAD
# Hyperledger-Fabric-CA
Multi Orderer with Kafka, Zookeeper
=======
# Hyperledger-Fabric-CA
>>>>>>> 78c108e11f823fcca07b2bca04e579258f11b8c0
=======
# Hyperledger Fabric Network 구축

## 사전작업
1. 로컬에서 네트워크 구축하기 위해서 virtualBox에 Ubuntu를 설치
2. node 개수만큼 vm 생성. (Cryptogen : 7개, fabric-ca : 7개)
Memory : 2GB 수동 디스크 할당 : 15GB
3. 각각 vm에 ip 설정
4. /etc/hosts 등록
peer0 10.0.1.11
peer1 10.0.1.12
peer2 10.0.1.21
peer3 10.0.1.22
orderer 10.0.1.31
client 10.0.1.41
fabric-ca 10.0.1.100

5. ~/.profile 환경변수 추가

export GOROOT=/usr/local/go
export gitlabPATH=$HOME/p1034_swing
export GOPATH=$gitlabPATH/saficket/
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
export FABRIC_HOME=$gitlabPATH/saficket/src/github.com/hyperledger/fabric
export PATH=$GOPATH/bin:$GOROOT/bin:$FABRIC_HOME/.build/bin:$PATH
export FABRIC_CFG_PATH=$GOPATH/cryptonet/$HOSTNAME
export FABRIC_CA_HOME=$GOPATH/src/github.com/hyperledger/fabric-ca
export PATH=$FABRIC_CA_HOME/bin:$PATH


## 간략한 구축 설명
구축하는 과정에서 msp를 생성하는 과정이 있는데 이미 생성되어 있기 때문에
따로 msp를 생성하는 과정을 하지 않아도 됩니다.

Cryptogen기반 네트워크와 Fabric-ca기반 네트워크의 차이는 msp 분배과정에 있습니다.

Cryptogen기반 네트워크는 command를 이용해 네트워크 전체 msp를 생성 후 모든 
peer와 orderer, client에게 분배하는 방식입니다.

Fabric-ca기반 네트워크는 Fabric-ca server를 구동하여 server에 msp를 요청하여
각 노드에 분배하는 방식입니다.

Saficket 프로젝트에서는 fabric-ca server를 사용하지 않는 것으로 계획을 변경해서
곧 수정할 예정입니다.
### Cryptogen기반 네트워크
cryptonet폴더는 Cryptogen을 기반으로 구축했으며, node는 peer0, peer1, peer2, 
peer3, orderer, client로 구성되어 있습니다.

cryptonet 폴더에서 진행합니다.

1. 각 peer 노드에서 runPeer.sh에 경로 수정 후 sudo ./runPeer.sh
2. orderer 노드에서 docker-compose up 로 kafka-zookeeper 구동하고나서
그 상태로 새 orderer vm 창 띄운 다음 runOrderer.sh에 경로 수정 후 
sudo ./runOrderer.sh
3. client 노드에서 create-channel.sh에 경로 수정 후 sudo ./create-channel.sh
4. 각 peer 노드에서 ./join-channl-peeer(NUM).sh

구축이 완료된다.

### Fabric-ca기반 네트워크
net 폴더는 fabric-ca를 기반으로 구축했으며, node는 peer0, peer1, peer2, 
peer3, orderer, Fabric-ca로 구성되어 있습니다.

순서는 Cryptogen기반 네트워크와 동일합니다.
폴더는 net폴더에서 진행합니다.

>>>>>>> fc4fc348523f0e9c36e53c554d124f6f01b7b0f1
