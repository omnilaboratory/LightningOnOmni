package rpc

import (
	"context"
	"errors"
	"github.com/omnilaboratory/obd/bean"
	"github.com/omnilaboratory/obd/bean/enum"
	"github.com/omnilaboratory/obd/proxy/pb"
	"github.com/omnilaboratory/obd/tool"
	"log"
)

func (s *RpcServer) OpenChannel(ctx context.Context, in *pb.OpenChannelRequest) (*pb.OpenChannelResponse, error) {
	log.Println("OpenChannel")
	if connObd == nil {
		return nil, errors.New("please login first")
	}

	if tool.CheckIsString(&in.RecipientInfo.RecipientNodePeerId) == false {
		return nil, errors.New("wrong recipient_node_peer_id")
	}

	if tool.CheckIsString(&in.RecipientInfo.RecipientUserPeerId) == false {
		return nil, errors.New("wrong recipient_user_peer_id")
	}

	if tool.CheckIsString(&in.NodePubkeyString) == false {
		return nil, errors.New("wrong node_pubkey_string")
	}

	if openChannelChan == nil {
		openChannelChan = make(chan bean.ReplyMessage)
	}

	channelOpen := bean.SendChannelOpen{
		FundingPubKey:      in.NodePubkeyString,
		FunderAddressIndex: int(in.NodePubkeyIndex),
		IsPrivate:          in.Private,
	}
	sendMsgToObd(channelOpen, in.RecipientInfo.RecipientNodePeerId, in.RecipientInfo.RecipientUserPeerId, enum.MsgType_SendChannelOpen_32)

	for {
		data := <-openChannelChan
		if data.Status == false {
			return nil, errors.New(data.Result.(string))
		}
		if data.Type == enum.MsgType_RecvChannelAccept_33 {
			log.Println(data.Result)
			resp := &pb.OpenChannelResponse{}
			resp.TemplateChannelId = data.Result.(map[string]interface{})["temporary_channel_id"].(string)
			return resp, nil
		}
	}
}

func (s *RpcServer) FundChannel(ctx context.Context, in *pb.FundChannelRequest) (*pb.FundChannelResponse, error) {
	log.Println("FundChannel")
	if connObd == nil {
		return nil, errors.New("please login first")
	}
	if tool.CheckIsString(&in.RecipientInfo.RecipientNodePeerId) == false {
		return nil, errors.New("wrong recipient_node_peer_id")
	}

	if tool.CheckIsString(&in.RecipientInfo.RecipientUserPeerId) == false {
		return nil, errors.New("wrong recipient_user_peer_id")
	}
	if tool.CheckIsString(&in.TemplateChannelId) == false {
		return nil, errors.New("wrong template_channel_id")
	}
	if in.BtcAmount < 0 {
		return nil, errors.New("wrong btc_amount")
	}
	if in.PropertyId < 0 {
		return nil, errors.New("wrong property_id")
	}
	if in.AssetAmount < 0 {
		return nil, errors.New("wrong asset_amount")
	}

	requestFunding := bean.SendRequestFunding{
		TemporaryChannelId: in.TemplateChannelId,
		BtcAmount:          in.BtcAmount,
		PropertyId:         in.PropertyId,
		AssetAmount:        in.AssetAmount,
	}

	if fundChannelChan == nil {
		fundChannelChan = make(chan bean.ReplyMessage)
	}

	sendMsgToObd(requestFunding, in.RecipientInfo.RecipientNodePeerId, in.RecipientInfo.RecipientUserPeerId, enum.MsgType_Funding_134)

	for {
		data := <-fundChannelChan
		if data.Status == false {
			return nil, errors.New(data.Result.(string))
		}
		if data.Type == enum.MsgType_ClientSign_AssetFunding_AliceSignRD_1134 {
			log.Println(data.Result)
			resp := &pb.FundChannelResponse{}
			resp.ChannelId = data.Result.(map[string]interface{})["channel_id"].(string)
			return resp, nil
		}
	}
}
func (s *RpcServer) RsmcPayment(ctx context.Context, in *pb.RsmcPaymentRequest) (*pb.RsmcPaymentResponse, error) {
	log.Println("RsmcPayment")
	if connObd == nil {
		return nil, errors.New("please login first")
	}
	if tool.CheckIsString(&in.RecipientInfo.RecipientNodePeerId) == false {
		return nil, errors.New("wrong recipient_node_peer_id")
	}

	if tool.CheckIsString(&in.RecipientInfo.RecipientUserPeerId) == false {
		return nil, errors.New("wrong recipient_user_peer_id")
	}
	if tool.CheckIsString(&in.ChannelId) == false {
		return nil, errors.New("wrong template_channel_id")
	}
	if in.Amount < 0 {
		return nil, errors.New("wrong amount")
	}

	request := bean.RequestCreateCommitmentTx{
		ChannelId: in.ChannelId,
		Amount:    in.Amount,
	}
	if rsmcChan == nil {
		rsmcChan = make(chan bean.ReplyMessage)
	}

	sendMsgToObd(request, in.RecipientInfo.RecipientNodePeerId, in.RecipientInfo.RecipientUserPeerId, enum.MsgType_CommitmentTx_SendCommitmentTransactionCreated_351)

	for {
		data := <-rsmcChan
		if data.Status == false {
			return nil, errors.New(data.Result.(string))
		}
		if data.Type == enum.MsgType_ClientSign_CommitmentTx_AliceSignC2a_360 {
			dataResult := data.Result.(map[string]interface{})
			resp := &pb.RsmcPaymentResponse{
				ChannelId: dataResult["channel_id"].(string),
				AmountA:   dataResult["amount_a"].(float64),
				AmountB:   dataResult["amount_b"].(float64),
			}
			return resp, nil
		}
	}

}
