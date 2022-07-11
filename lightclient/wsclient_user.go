package lightclient

import (
	"encoding/json"
	"errors"
	"github.com/omnilaboratory/obd/bean"
	"github.com/omnilaboratory/obd/bean/enum"
	"github.com/omnilaboratory/obd/config"
	"github.com/omnilaboratory/obd/service"
	"github.com/omnilaboratory/obd/tool"
	"github.com/omnilaboratory/obd/tracker/tkrpc"
	"github.com/tidwall/gjson"
)

func loginRetData(client Client) string {
	retData := make(map[string]interface{})
	retData["userPeerId"] = client.User.PeerId
	retData["nodePeerId"] = client.User.P2PLocalPeerId
	retData["nodeAddress"] = client.User.P2PLocalAddress
	retData["htlcFeeRate"] = config.HtlcFeeRate
	retData["htlcMaxFee"] = config.HtlcMaxFee
	retData["chainNodeType"] = config.ChainNodeType
	bytes, _ := json.Marshal(retData)
	return string(bytes)
}

func (client *Client) UserModule(msg bean.RequestMessage) (enum.SendTargetType, []byte, bool) {
	status := false
	var sendType = enum.SendTargetType_SendToNone
	var data string
	switch msg.Type {
	case enum.MsgType_UserLogin_2001:
		//for js sdk user
		walletId := gjson.Get(msg.Data, "wallet_id").String()

		//for admin user
		mnemonic := gjson.Get(msg.Data, "mnemonic").String()

		isAdmin := false

		//TODO 管理员验证
		if client.IsGRpcRequest {
			loginToken := gjson.Get(msg.Data, "login_token").Str
			isAdmin = service.CheckIsAdmin(loginToken)
		}
		peerId := ""
		if isAdmin {
			peerId = tool.GetUserPeerId(mnemonic)
		} else {
			peerId = tool.GetUserPeerId(walletId)
			//peerId = tool.GetUserPeerId(mnemonic)
		}
		if GlobalWsClientManager.OnlineClientMap[peerId] != nil {
			if GlobalWsClientManager.OnlineClientMap[peerId].User.IsAdmin {
				client.User = GlobalWsClientManager.OnlineClientMap[peerId].User
				GlobalWsClientManager.OnlineClientMap[peerId] = client
			} else {
				if isAdmin {
					client.User = GlobalWsClientManager.OnlineClientMap[peerId].User
					client.User.IsAdmin = true
					GlobalWsClientManager.OnlineClientMap[peerId] = client
				}
			}
		}
		if client.User != nil {
			if client.User.Mnemonic != mnemonic {
				_ = service.UserService.UserLogout(client.User)
				//sendInfoOnUserStateChange(client.User.PeerId)
				delete(GlobalWsClientManager.OnlineClientMap, client.User.PeerId)
				delete(service.OnlineUserMap, client.User.PeerId)
				client.User = nil
			}
		}

		if client.User != nil {
			data = loginRetData(*client)
			status = true
			client.SendToMyself(msg.Type, status, data)
			sendType = enum.SendTargetType_SendToSomeone
		} else {
			user := bean.User{
				Mnemonic:        mnemonic,
				PeerId:          peerId,
				P2PLocalAddress: localServerDest,
				P2PLocalPeerId:  P2PLocalNodeId,
				IsAdmin:         isAdmin,
			}
			var err error = nil
			if GlobalWsClientManager.OnlineClientMap[peerId] != nil {
				err = errors.New("user has login at other node")
			} else {
				err = service.UserService.UserLogin(&user)
				if err == nil {
					//sendInfoOnUserStateChange(user.PeerId)
				}
			}
			if err == nil {
				client.User = &user
				GlobalWsClientManager.OnlineClientMap[user.PeerId] = client
				service.OnlineUserMap[user.PeerId] = &user
				data = loginRetData(*client)
				status = true
				client.SendToMyself(msg.Type, status, data)
				sendType = enum.SendTargetType_SendToExceptMe
			} else {
				client.SendToMyself(msg.Type, status, err.Error())
				sendType = enum.SendTargetType_SendToSomeone
			}
		}
	case enum.MsgType_User_GetInfo_2009:
		if client.User != nil {
			data = loginRetData(*client)
			status = true
		} else {
			data = "please login"
		}
		client.SendToMyself(msg.Type, status, data)
		sendType = enum.SendTargetType_SendToSomeone
	case enum.MsgType_UserLogout_2002:
		if client.User != nil {
			exist := service.UserService.CheckExecutingTx(client.User)
			if exist == false {
				data = client.User.PeerId + " logout"
				status = true
			} else {
				data = "exist executing tx ,can not logout"
			}
			client.SendToMyself(msg.Type, status, data)

			if exist == false {
				user:=client.User
				ITclient.UpdateUserInfo(todo,&tkrpc.UpdateUserInfoReq{UserId:user.PeerId,NodeId:user.P2PLocalPeerId,P2PAddress:user.P2PLocalAddress,IsOnline: 2})

				_ = service.UserService.UserLogout(client.User)
				//sendInfoOnUserStateChange(client.User.PeerId)

				delete(GlobalWsClientManager.ClientsMap, client)
				delete(GlobalWsClientManager.OnlineClientMap, client.User.PeerId)
				delete(service.OnlineUserMap, client.User.PeerId)
				client.User = nil
			}
		} else {
			data = "please login"
			client.SendToMyself(msg.Type, status, data)
		}
		sendType = enum.SendTargetType_SendToSomeone
	case enum.MsgType_p2p_ConnectPeer_2003:
		remoteNodeAddress := gjson.Get(msg.Data, "remote_node_address")
		if remoteNodeAddress.Exists() == false {
			data = errors.New("remote_node_address not exist").Error()
		} else {
			localP2PAddress, err := connP2PNode(remoteNodeAddress.Str)
			if err != nil {
				data = err.Error()
			} else {
				status = true
				data = localP2PAddress
			}
		}
		client.SendToMyself(msg.Type, status, data)
		sendType = enum.SendTargetType_SendToSomeone
	case enum.MsgType_p2p_DisconnectPeer_2010:
		remoteNodeAddress := gjson.Get(msg.Data, "remote_node_address")
		if remoteNodeAddress.Exists() == false {
			data = errors.New("remote_node_address not exist").Error()
		} else {
			err := disConnP2PNode(remoteNodeAddress.Str)
			if err != nil {
				data = err.Error()
			} else {
				status = true
				data = "success"
			}
		}
		client.SendToMyself(msg.Type, status, data)
		sendType = enum.SendTargetType_SendToSomeone
	case enum.MsgType_GetObdNodeInfo_2005:
		bytes, err := json.Marshal(bean.CurrObdNodeInfo)
		if err != nil {
			data = err.Error()
		} else {
			status = true
			data = string(bytes)
		}
		client.SendToMyself(msg.Type, true, data)
		sendType = enum.SendTargetType_SendToSomeone
	case enum.MsgType_GetMiniBtcFundAmount_2006:
		fee := service.GetBtcMinerFundMiniAmount()
		data = tool.FloatToString(fee, 8)
		client.SendToMyself(msg.Type, true, data)
		sendType = enum.SendTargetType_SendToSomeone
	case enum.MsgType_HeartBeat_2007:
		client.SendToMyself(msg.Type, true, data)
		sendType = enum.SendTargetType_SendToSomeone
	case enum.MsgType_User_UpdateAdminToken_2008:
		oldLoginToken := gjson.Get(msg.Data, "current_password").Str
		newLoginToken := gjson.Get(msg.Data, "new_password").Str
		if client.User != nil && client.User.IsAdmin {
			err := service.UpdateAdminLoginToken(oldLoginToken, newLoginToken)
			if err != nil {
				data = err.Error()
			} else {
				data = newLoginToken
				status = true
			}
		} else {
			data = errors.New("you are not the admin or login as admin").Error()
		}
		client.SendToMyself(msg.Type, status, data)
		sendType = enum.SendTargetType_SendToSomeone
	// Process GetMnemonic
	case enum.MsgType_GetMnemonic_2004:
		if client.User != nil { // The user already login.
			client.SendToMyself(msg.Type, true, "already login")
			sendType = enum.SendTargetType_SendToSomeone
		} else {
			// get Mnemonic
			mnemonic, err := service.HDWalletService.Bip39GenMnemonic(256)
			if err == nil { //get  successful.
				data = mnemonic
				status = true
			} else {
				data = err.Error()
			}
			client.SendToMyself(msg.Type, status, data)
			sendType = enum.SendTargetType_SendToSomeone
		}
	}
	return sendType, []byte(data), status
}
