package main

import (
	"fmt"

	fcm "github.com/NaySoftware/go-fcm"
	"github.com/lastforeverzl/barkme-server/config"
	"github.com/lastforeverzl/barkme-server/message"
	"github.com/lastforeverzl/barkme-server/mydb"
)

type Fcm interface {
	SendNotification(message.Envelope, <-chan *mydb.TokensChan) <-chan *FcmResponse
}

type FcmClient struct {
	*fcm.FcmClient
}

type FcmResponse struct {
	Response *fcm.FcmResponseStatus
	Err      error
}

func initFcm(cfg *config.ServerConfig) *FcmClient {
	c := fcm.NewFcmClient(cfg.FcmServerKey)
	return &FcmClient{c}
}

func (client *FcmClient) SendNotification(envelope message.Envelope, in <-chan *mydb.TokensChan) <-chan *FcmResponse {
	out := make(chan *FcmResponse)
	var NP fcm.NotificationPayload
	NP.Title = fmt.Sprintf("%v %v", envelope.Username, envelope.Msg)
	NP.Body = "I'm barking"
	data := map[string]string{
		"sender": envelope.Username,
		"msg":    envelope.Msg,
		"token":  envelope.Token,
	}
	go func() {
		tokens := <-in
		client.NewFcmRegIdsMsg(tokens.Tokens, data)
		client.SetNotificationPayload(&NP)
		status, err := client.Send()
		out <- &FcmResponse{Response: status, Err: err}
		close(out)
	}()
	return out
}
