// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package opcua

import (
	"fmt"

	"github.com/mainflux/mainflux/logger"
)

const protocol = "opcua"

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateThing creates thing  mfx:opc & opc:mfx route-map
	CreateThing(string, string) error

	// UpdateThing updates thing mfx:opc & opc:mfx route-map
	UpdateThing(string, string) error

	// RemoveThing removes thing mfx:opc & opc:mfx route-map
	RemoveThing(string) error

	// CreateChannel creates channel route-map
	CreateChannel(string, string) error

	// UpdateChannel updates chroute-map
	UpdateChannel(string, string) error

	// RemoveChannel removes channel route-map
	RemoveChannel(string) error

	// ConnectThing creates thing and channel connection route-map
	ConnectThing(string, string) error

	// DisconnectThing removes thing and channel connection route-map
	DisconnectThing(string, string) error
}

// Config OPC-UA Server
type Config struct {
	ServerURI string
	NodeID    string
	Policy    string
	Mode      string
	CertFile  string
	KeyFile   string
}

var _ Service = (*adapterService)(nil)

type adapterService struct {
	subscriber Subscriber
	thingsRM   RouteMapRepository
	channelsRM RouteMapRepository
	connectRM  RouteMapRepository
	cfg        Config
	logger     logger.Logger
}

// New instantiates the OPC-UA adapter implementation.
func New(sub Subscriber, thingsRM, channelsRM, connectRM RouteMapRepository, cfg Config, log logger.Logger) Service {
	return &adapterService{
		subscriber: sub,
		thingsRM:   thingsRM,
		channelsRM: channelsRM,
		connectRM:  connectRM,
		cfg:        cfg,
		logger:     log,
	}
}

func (as *adapterService) CreateThing(mfxDevID, opcuaNodeID string) error {
	return as.thingsRM.Save(mfxDevID, opcuaNodeID)
}

func (as *adapterService) UpdateThing(mfxDevID, opcuaNodeID string) error {
	return as.thingsRM.Save(mfxDevID, opcuaNodeID)
}

func (as *adapterService) RemoveThing(mfxDevID string) error {
	return as.thingsRM.Remove(mfxDevID)
}

func (as *adapterService) CreateChannel(mfxChanID, opcuaServerURI string) error {
	return as.channelsRM.Save(mfxChanID, opcuaServerURI)
}

func (as *adapterService) UpdateChannel(mfxChanID, opcuaServerURI string) error {
	return as.channelsRM.Save(mfxChanID, opcuaServerURI)
}

func (as *adapterService) RemoveChannel(mfxChanID string) error {
	return as.channelsRM.Remove(mfxChanID)
}

func (as *adapterService) ConnectThing(mfxChanID, mfxThingID string) error {
	serverURI, err := as.channelsRM.Get(mfxChanID)
	if err != nil {
		return err
	}

	nodeID, err := as.thingsRM.Get(mfxThingID)
	if err != nil {
		return err
	}

	as.cfg.NodeID = nodeID
	as.cfg.ServerURI = serverURI
	go as.subscribe(as.cfg)

	c := fmt.Sprintf("%s:%s", mfxChanID, mfxThingID)
	return as.connectRM.Save(c, c)
}

func (as *adapterService) DisconnectThing(mfxChanID, mfxThingID string) error {
	c := fmt.Sprintf("%s:%s", mfxChanID, mfxThingID)
	return as.connectRM.Remove(c)
}

func (as *adapterService) subscribe(cfg Config) {
	if err := as.subscriber.Subscribe(cfg); err != nil {
		as.logger.Warn(fmt.Sprintf("subscription failed: %s", err))
	}
}
