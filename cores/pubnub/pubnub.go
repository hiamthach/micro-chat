package cores_pubnub

import (
	"fmt"

	"github.com/hiamthach/micro-chat/util"
	pubnub "github.com/pubnub/go/v7"
)

type IPubNubHelper interface {
	Init(config util.Config)
	InsWithUserId(uuid string) *pubnub.PubNub
}

// PubnubUtil is a wrapper for pubnub client
type PubNubHelper struct {
	IPubNubHelper
	config util.Config
}

// Init initializes pubnub client
func (p *PubNubHelper) Init(config util.Config) {
	p.config = config
}

// GetInstance returns pubnub client
func (p *PubNubHelper) InsWithUserId(uuid string) *pubnub.PubNub {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Errorf("InsWithUserId has error : %v", r))
		}
	}()

	config := pubnub.NewConfigWithUserId(pubnub.UserId(uuid))
	config.UseHTTP2 = false
	config.PublishKey = p.config.PubNubPublishKey
	config.SubscribeKey = p.config.PubNubSubscribeKey

	pn := pubnub.NewPubNub(config)

	return pn
}
