package main

import (
	"bufio"
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"time"

	"github.com/0xAX/notificator"
	"github.com/miekg/dns"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var (
	notify           *notificator.Notificator
	CheckInetErrorSt = 0
)

func CheckInet() bool {
	errCnt := 0
	CntPacketSendPing := 4
	st := SetPing{CfgFpinger.Fpinger.ArrayIpAddrsPing, time.Second * 2, CntPacketSendPing}
	Stat := Pinger(st)
	for _, err := range Stat.RespErrPacket {
		if err != nil {
			//			WriteLog("current time: " + time.Now().Format("2006-01-02 15:04:05") + " error ping: " + Stat.ArrayIpAddr[i])
			errCnt++
		}
	}
	if errCnt != CntPacketSendPing*2 {
		if CheckInetErrorSt == 1 {
			WriteLog("current time: " + time.Now().Format("2006-01-02 15:04:05") + " появился интернет...")
			notifyRun("появился локальный интернет...")
		}
		CheckInetErrorSt = 0
		return true
	} else {
		if CheckInetErrorSt == 1 {
			return false
		}
		WriteLog("current time: " + time.Now().Format("2006-01-02 15:04:05") + " пропал интернет...")
		notifyRun("пропал локальный интернет...")
		CheckInetErrorSt = 1
		return false
	}
}
func Pinger(param0 SetPing) RespPacketPing {
	var elapsed time.Duration
	var ArrIpList []string
	var ArrRespPacket []time.Duration
	var ArrRespError []error
	w1 := param0.TimeWaitPacket
	for CntPacketSendPing := 0; CntPacketSendPing < param0.CntPacketSend; CntPacketSendPing++ {
		for _, IpAddr := range param0.ArrayIpAddr {
			time.Sleep(time.Second * 1)
			ArrIpList = append(ArrIpList, IpAddr)
			start := time.Now()
			c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
			if err != nil {
				ArrRespError = append(ArrRespError, err)
				ArrRespPacket = append(ArrRespPacket, time.Second*0)
				continue
			}
			wm := icmp.Message{
				Type: ipv4.ICMPTypeEcho, Code: 0,
				Body: &icmp.Echo{
					ID: os.Getpid() & 0xDaAa, Seq: 1,
					//					Data: []byte("SuPeRCodEr@1144-1024"),
					Data: []byte("my coder..."),
				},
			}
			wb, err := wm.Marshal(nil)
			if err != nil {
				ArrRespError = append(ArrRespError, err)
				ArrRespPacket = append(ArrRespPacket, time.Second*0)
				continue
			}
			if _, err := c.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(IpAddr)}); err != nil {
				ArrRespError = append(ArrRespError, err)
				ArrRespPacket = append(ArrRespPacket, time.Second*0)
				c.Close()
				continue
			}
			rb := make([]byte, 1500)
			timeoutDuration := w1
			c.SetReadDeadline(time.Now().Add(timeoutDuration))
			n, _, err := c.ReadFrom(rb)
			if err != nil {
				elapsed = time.Since(start)
				ArrRespError = append(ArrRespError, err)
				ArrRespPacket = append(ArrRespPacket, time.Second*0)
				c.Close()
				continue

			}
			rm, err := icmp.ParseMessage(1, rb[:n])
			if rm.Code != 0 {
				ArrRespError = append(ArrRespError, errors.New("destination unreachable"))
				ArrRespPacket = append(ArrRespPacket, time.Second*0)
				c.Close()
				continue
			}
			if err != nil {
				ArrRespError = append(ArrRespError, err)
				ArrRespPacket = append(ArrRespPacket, time.Second*0)
				c.Close()
				continue
			}
			switch rm.Type {
			case ipv4.ICMPTypeEchoReply:
				elapsed = time.Since(start)
				ArrRespPacket = append(ArrRespPacket, elapsed)
				ArrRespError = append(ArrRespError, err)

			default:
				elapsed = time.Since(start)
			}
			defer func() {
				c.Close()

			}()
		}
	}
	return RespPacketPing{ArrayIpAddr: ArrIpList,
		RespArrayTimePacket: ArrRespPacket, RespErrPacket: ArrRespError}
}

// get interface and IPv4
func getIfaceIpv4() ([]string, []string, error) {
	var ip net.IP
	var ArrIFace []string
	var ArrIPv4 []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ArrIFace, ArrIPv4, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil && !v.IP.IsLoopback() {
					ip = v.IP
					ArrIPv4 = append(ArrIPv4, ip.String())
					ArrIFace = append(ArrIFace, i.Name)
				}
			}
		}
		if err != nil {
			return ArrIFace, ArrIPv4, err
		}
	}
	return ArrIFace, ArrIPv4, err
}
func strTm(st string) (time.Duration, bool) {
	status_1 := true
	var TimeWaitPacket time.Duration
	switch strings.ToLower(st[len(st)-1 : len(st)]) {
	case "s":
		st_sec0, err := strconv.Atoi(st[0 : len(st)-1])
		if err != nil {
			TimeWaitPacket = time.Second * 10 //default 1 secound time wait
			status_1 = false
		}
		TimeWaitPacket = time.Second * time.Duration(st_sec0)
	case "m":
		st_sec1, err := strconv.Atoi(st[0 : len(st)-1])
		if err != nil {
			TimeWaitPacket = time.Minute * 1 //default 1 minute time wait
			status_1 = false
		}
		TimeWaitPacket = time.Minute * time.Duration(st_sec1)
	default:
		TimeWaitPacket = time.Second * 10 //default 1 secound time wait
	}
	return TimeWaitPacket, status_1
}

func dnsQuery0(dnsname1 string, targetIPDNS string) (string, time.Duration, error) {
	var result_dnsname string
	var result_rtt time.Duration
	var err error
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{dnsname1, dns.TypeA, dns.ClassINET}
	c := new(dns.Client)
	in, rtt, err := c.Exchange(m1, targetIPDNS+":53")
	if err == nil {
		if len(in.Answer) == 0 {
			err = errors.New("error answer size 0")
			return result_dnsname, result_rtt, err
		}
		if t, ok := in.Answer[0].(*dns.A); ok {
			result_dnsname = t.A.String()
			result_rtt = rtt
		}
	} else {
		return result_dnsname, result_rtt, err
	}
	return result_dnsname, result_rtt, err
}

//
func WriteLog(mytxt interface{}) {
	file, err := os.OpenFile("monitor.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("File does not exists or cannot be created")
		os.Exit(1)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, mytxt)
	w.Flush()

}
func notifyRun(s string) {
	notify = notificator.New(notificator.Options{
		DefaultIcon: "./logo/monitoring.png",
		AppName:     "Monitoring",
	})
	err := notify.Push("vps testing...", s, "", notificator.UR_NORMAL)
	if err != nil {
		WriteLog("current time: " + time.Now().Format("2006-01-02 15:04:05") + "Error Notify Push: " + err.Error())
		log.Println("error: ")
	}
}

// check run root
func checkRootId() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		log.Fatal(err)
	}

	if i == 0 {
		return true
	} else {
		return false
	}
}
