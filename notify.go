package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/NaySoftware/go-fcm"
)

func ReadConfigFpinger() (c *ConfigFpinger) {
	file, _ := os.Open("monvps.json")
	decoder := json.NewDecoder(file)
	ConfigFpinger := new(ConfigFpinger)
	err := decoder.Decode(&ConfigFpinger)
	if err != nil {
		log.Fatalln("error parse config: ", err)
	}
	return ConfigFpinger
}
func ReadConfigNotify() (c *ConfigFCM) {
	file, _ := os.Open("notify.json")
	decoder := json.NewDecoder(file)
	ConfigFCM := new(ConfigFCM)
	err := decoder.Decode(&ConfigFCM)
	if err != nil {
		log.Fatalln("error parse config: ", err)
	}
	return ConfigFCM
}
func NotifyRunFcm(title string, body string) {
	var NP fcm.NotificationPayload
	NP.Title = title
	NP.Body = body

	data := map[string]string{
		"msg": "Hello World1",
		"sum": "Happy Day",
	}

	ids := []string{
		CfgFCM.PushFCMNotification.IdsDeviceTokens[0],
	}

	xds := []string{}

	c := fcm.NewFcmClient(CfgFCM.PushFCMNotification.ApiKey)
	c.NewFcmRegIdsMsg(ids, data)
	c.AppendDevices(xds)
	c.SetNotificationPayload(&NP)
	_, err := c.Send()
	if err == nil {
		//		status.PrintResults()
	} else {
		fmt.Println(err)
	}
}
