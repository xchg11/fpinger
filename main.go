package main

import (
	"log"
	"time"
)

var (
	CfgFCM     *ConfigFCM
	StatErr    = []int{}
	CfgFpinger *ConfigFpinger
)

func init() {
	CfgFCM = ReadConfigNotify()
	CfgFpinger = ReadConfigFpinger()

}
func main() {
	if !checkRootId() {
		log.Fatal("This program must be run as root! (sudo)")
	}
	for {
		// my commit
		//run pinger start ...
		if !CheckInet() { // если ошибка то ждем 1 минуту и опять проверяем
			time.Sleep(time.Second * 60) // 1 минута
			continue
		}
		for index0, listIP := range CfgFpinger.Fpinger.NsDNSChecks {
			_, _, err := dnsQuery0("fire7.ru.", listIP)
			StatErr = append(StatErr, 0)
			if err != nil {
				StatErr[index0]++

			} else {
				if StatErr[index0] >= 100 {
					notifyRun("VPS(днс) восстановлен. " + CfgFpinger.Fpinger.NsDNSChecks[index0])
					//					NotifyRunFcm("VPS RECOVERY", "VPS(днс) восстановлен. "+listVpsIp[index0])
					WriteLog("current time: " + time.Now().Format("2006-01-02 15:04:05") + " VPS(днс) восстановлен. " + CfgFpinger.Fpinger.NsDNSChecks[index0])
				}
				StatErr[index0] = 0
			}
			if StatErr[index0] == 2 { // если не работает 10 минут dns
				//notify
				notifyRun("VPS(днс) не отвечает " + CfgFpinger.Fpinger.NsDNSChecks[index0])
				//				NotifyRunFcm("VPS faild dns", "VPS(днс) не отвечает "+listVpsIp[index0])
				WriteLog("current time: " + time.Now().Format("2006-01-02 15:04:05") + " VPS(днс) не отвечает " + CfgFpinger.Fpinger.NsDNSChecks[index0])
				StatErr[index0] = 100
				//
			}

		}
		time.Sleep(time.Second * 60) // 1 минута
	}
}
