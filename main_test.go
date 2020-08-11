package main

import "testing"
import "os"
import "encoding/json"

func TestCheckRoot(t *testing.T) {
	if !checkRootId() {
		t.Error("This program must be run as root! (sudo)")
	}
}
func testCheckDNS(t *testing.T) {
	_, _, err := dnsQuery0("fire7.ru.", "1.1.1.1")
	if err != nil {
		t.Error("error: testCheckDNS")
	}
}
func TestReadConfigFpinger(t *testing.T) {
	file, err := os.Open("monvps.json")
	if err != nil {
		t.Error(err)
	}
	decoder := json.NewDecoder(file)
	ConfigFpinger := new(ConfigFpinger)
	err = decoder.Decode(&ConfigFpinger)
	if err != nil {
		t.Error("error parse config: ", err)
	}
}
