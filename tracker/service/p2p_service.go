package service

import
(
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/omnilaboratory/obd/bean"
	cfg "github.com/omnilaboratory/obd/tracker/config"
	"github.com/omnilaboratory/obd/tracker/tkrpc"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const protocolIdForScanObd = "obd/forScanObd/1.0.1"
const obdRendezvousString = "obd meet at tracker"
const trackerRendezvousString = "tracker meet here"

var ctx = context.Background()
var kademliaDHT *dht.IpfsDHT
var routingDiscovery *discovery.RoutingDiscovery
var hostNode host.Host

func StartP2PNode() {

	nodeId := int64(binary.BigEndian.Uint64([]byte(GetTrackerNodeId())))
	r := rand.New(rand.NewSource(nodeId))
	prvKey, _, err := crypto.GenerateECDSAKeyPair(r)
	if err != nil {
		log.Println(err)
		return
	}

	sourceMultiAddr, _ := ma.NewMultiaddr("/ip4/0.0.0.0/tcp/" + strconv.Itoa(cfg.P2P_sourcePort))
	hostNode, err = libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
		libp2p.EnableRelay(circuit.OptHop),
	)
	if err != nil {
		log.Println(err)
		return
	}

	hostNode.SetStreamHandler(bean.ProtocolIdForUserState, handleUserStateStream)
	hostNode.SetStreamHandler(bean.ProtocolIdForChannelInfoChange, handleChannelStream)

	cfg.P2pLocalAddress = fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", cfg.P2P_hostIp, cfg.P2P_sourcePort, hostNode.ID().Pretty())
	log.Println("local p2p node address: ", cfg.P2pLocalAddress)

	kademliaDHT, err = dht.New(ctx, hostNode, dht.Mode(dht.ModeServer))

	if err != nil {
		log.Println(err)
		return
	}

	err = kademliaDHT.Bootstrap(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	routingDiscovery = discovery.NewRoutingDiscovery(kademliaDHT)

	startSchedule()
}

func startSchedule() {
	announceSelf()
	scanNodes()
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				log.Println("timer 1m", t)
				announceSelf()
				scanNodes()
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				log.Println("timer 10s", t)
				go updateLockChannel()
			}
		}
	}()
}

func announceSelf() {
	needAnnounceSelf := false
	if len(cfg.BootstrapPeers) > 0 {
		var wg sync.WaitGroup
		for _, peerAddr := range cfg.BootstrapPeers {
			peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
			if peerInfo.ID == hostNode.ID() {
				continue
			}
			log.Println("connecting to bootstrap node: ", peerAddr)
			wg.Add(1)
			needAnnounceSelf = true
			go func() {
				defer wg.Done()
				err := hostNode.Connect(ctx, *peerInfo)
				if err != nil {
					log.Println(err, peerInfo)
				} else {
					log.Println("connect to bootstrap node ", *peerInfo)
				}
			}()
		}
		wg.Wait()
	}

	if needAnnounceSelf {
		log.Println("announce self")
		discovery.Advertise(ctx, routingDiscovery, trackerRendezvousString)
	}
}

func scanNodes() {
	log.Println("scanNodes")
	peerChan, err := routingDiscovery.FindPeers(ctx, obdRendezvousString)
	if err != nil {
		log.Println(err)
		return
	}

	for node := range peerChan {
		log.Println("target p2p node id", node.ID)
		if node.ID == hostNode.ID() {
			continue
		}

		//和tracker直接连接的obd，不需要同步数据 因为用户登录，都会直接向tracker报告登录状态
		if obdOnlineNodesMap[node.ID.Pretty()] != nil {
			continue
		}

		log.Println("begin to connect ", node.ID, node.Addrs)
		err = hostNode.Connect(ctx, node)
		if err == nil {
			stream, err := hostNode.NewStream(ctx, node.ID, protocolIdForScanObd)
			if err == nil {
				go handleScanStream(stream)
			}
		} else {
			delete(userOnlineOfOtherObdMap, node.ID.Pretty())
		}
	}
}

var userOnlineOfOtherObdMap = make(map[string]map[string]string)

func handleScanStream(stream network.Stream) {
	log.Println("begin scan obd channel")
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	str, err := rw.ReadString('~')
	if err != nil {
		return
	}

	if str == "" {
		return
	}
	if str != "" {
		str = strings.TrimSuffix(str, "~")
		data := make(map[string]string)
		err = json.Unmarshal([]byte(str), &data)
		if err == nil {
			//online user
			delete(userOnlineOfOtherObdMap, stream.Conn().RemotePeer().Pretty())
			if _, ok := data["userInfo"]; ok == true {
				log.Println(data["userInfo"])
				userInfo := make(map[string]string)
				err = json.Unmarshal([]byte(data["userInfo"]), &userInfo)
				if err == nil {
					userOnlineOfOtherObdMap[stream.Conn().RemotePeer().Pretty()] = userInfo
				}
			}
			//channel
			if _, ok := data["channelInfo"]; ok == true {
				_ = ChannelService.updateChannelInfo(data["obdP2pNodeId"], data["channelInfo"])
			}
		}
	}
	_ = stream.Close()
}

func handleUserStateStream(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	str, err := rw.ReadString('~')
	if err != nil {
		return
	}
	if str == "" {
		return
	}
	if str != "" {
		str = strings.TrimSuffix(str, "~")
		log.Println("handleUserStateStream", str)
		params := strings.Split(str, "_")
		if len(params) > 1 {
			if _, ok := userOnlineOfOtherObdMap[params[0]]; ok == true {
				if _, ok = userOnlineOfOtherObdMap[params[0]][params[1]]; ok == true {
					delete(userOnlineOfOtherObdMap[params[0]], params[1])
				} else {
					if userOnlineOfOtherObdMap[params[0]] == nil {
						userOnlineOfOtherObdMap[params[0]] = make(map[string]string)
					}
					userOnlineOfOtherObdMap[params[0]][params[1]] = params[2]
				}
			} else {
				if userOnlineOfOtherObdMap[params[0]] == nil {
					userOnlineOfOtherObdMap[params[0]] = make(map[string]string)
				}
				userOnlineOfOtherObdMap[params[0]][params[1]] = params[2]
			}
		}
	}
	_ = stream.Close()
}

func handleChannelStream(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	str, err := rw.ReadString('~')
	if err != nil {
		return
	}
	if str == "" {
		return
	}
	if str != "" {
		str = strings.TrimSuffix(str, "~")
		log.Println("handleChannelStream", str)
		if len(str) > 0 {
			_ = ChannelService.updateChannelInfo(stream.Conn().RemotePeer().Pretty(), str)
		}
	}
	_ = stream.Close()
}

func SendChannelLockInfoToObd(channelId, userId, obdP2pNodeId string, isUnlock bool) bool {
	findID, err := peer.Decode(obdP2pNodeId)
	if err == nil {
		findPeer, err := kademliaDHT.FindPeer(ctx, findID)
		if err == nil {
			protoId:=bean.ProtocolIdForLockChannel
			if isUnlock {
				protoId=bean.ProtocolIdForUnlockChannel
			}
			stream, err := hostNode.NewStream(ctx, findPeer.ID, protocol.ID( protoId))
			if err == nil {
				rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
				request := bean.TrackerLockChannelRequest{UserId: userId, ChannelId: channelId}
				marshal, _ := json.Marshal(request)
				_, _ = rw.WriteString(string(marshal) + "~")
				err = rw.Flush()
				if err == nil {
					str, err := rw.ReadString('~')
					if err != nil {
						return false
					}
					if str == "" {
						return false
					}
					if str != "" {
						str = strings.TrimSuffix(str, "~")
						log.Println("OnSendChannelLockInfoToObd", str)
						if str == "1" {
							return true
						}
					}
				}
				_ = stream.Close()
			}
		}
	}
	return false
}

var Orm=cfg.Orm
func updateLockChannel() {
	var infos []tkrpc.LockHtlcPath
	now := time.Now().Add(-20 * time.Second)
	Orm.Find(&infos,"(curr_state=0 or curr_state=1 ) and created_at <?",now.Unix())
	for _, item := range infos {
		paths := strings.Split( item.Path,",")
		index := len(paths) - 1
		channelId := paths[index]
		channelInfo := &tkrpc.ChannelInfo{}
		if item.CurrState == 0 {
			// get anc check First channel, whether the path is invalid
			err:=Orm.First(channelInfo,tkrpc.ChannelInfo{ChannelId:channelId})
			if err == nil {
				if channelInfo.CurrState == tkrpc.ChannelState_LockByTracker {
					item.CurrState = 1
					_ = Orm.Save(&item)
				}
				if channelInfo.CurrState != tkrpc.ChannelState_LockByTracker {
					item.CurrState = 2
					_ = Orm.Save(&item)
					return
				}
			}
		}
		notifyObdFinish := true
		for index = len(paths) - 1; index >= 0; index-- {
			channelId = paths[index]
			channelInfo := &tkrpc.ChannelInfo{}
			err:=Orm.First(channelInfo,tkrpc.ChannelInfo{ChannelId:channelId})
			if err == nil {
				if SendChannelLockInfoToObd(channelInfo.ChannelId, channelInfo.PeerIda, channelInfo.ObdNodeIda,true) &&
					SendChannelLockInfoToObd(channelInfo.ChannelId, channelInfo.PeerIdb, channelInfo.ObdNodeIdb,true) {
					channelInfo.CurrState = tkrpc.ChannelState_CanUse
					_ = Orm.Save(channelInfo)
				} else {
					notifyObdFinish = false
				}
			}
		}
		if notifyObdFinish {
			item.CurrState = 2
			_ = Orm.Save(&item)
		}
	}
}

