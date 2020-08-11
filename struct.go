package main

import (
	"time"
)

type RespPacketPing struct {
	ArrayIpAddr         []string
	RespArrayTimePacket []time.Duration
	RespErrPacket       []error
}
type SetPing struct {
	ArrayIpAddr    []string
	TimeWaitPacket time.Duration
	CntPacketSend  int
}
type ConfigFCM struct {
	PushFCMNotification FcmSettings
}
type FcmSettings struct {
	ApiKey          string
	IdsDeviceTokens []string
}
type ConfigFpinger struct {
	Fpinger FpingerSet
}
type FpingerSet struct {
	ArrayIpAddrsPing []string
	NsDNSChecks      []string
	Status           int32
	IntervalCheck    string
	ServiceSetStatus int16
}
