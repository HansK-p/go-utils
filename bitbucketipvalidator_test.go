package utils

import (
	"net"
	"testing"

	log "github.com/sirupsen/logrus"
)

type checkdatabiv struct {
	ipAddr         string
	expectedResult bool
}

func TestValidateIP(t *testing.T) {
	checkdatas := []checkdatabiv{
		{
			ipAddr:         "10.10.10.10",
			expectedResult: false,
		}, {
			ipAddr:         "20.20.20.20",
			expectedResult: false,
		}, { // "cidr": "18.234.32.128/25"
			ipAddr:         "18.234.32.3",
			expectedResult: false,
		}, {
			ipAddr:         "18.234.32.127",
			expectedResult: false,
		}, {
			ipAddr:         "18.234.32.128",
			expectedResult: true,
		}, {
			ipAddr:         "18.234.32.136",
			expectedResult: true,
		}, {
			ipAddr:         "18.234.32.136",
			expectedResult: true,
		}, { // "cidr": "34.218.168.212/32"
			ipAddr:         "34.218.168.211",
			expectedResult: false,
		}, {
			ipAddr:         "34.218.168.212",
			expectedResult: true,
		}, {
			ipAddr:         "34.218.168.213",
			expectedResult: false,
		},
	}
	logger := log.New()
	biv, err := NewBitbucketIPValidator(logger)
	if err != nil {
		t.Fatalf("Error creating a new BitbucketIpValidator: %v", err)
	}
	for _, checkdata := range checkdatas {
		ip := net.ParseIP(checkdata.ipAddr)
		valid := biv.ValidateIP(ip)
		if valid != checkdata.expectedResult {
			t.Fatalf("For IP %v the expected validation result is %v, but it was %v", ip, checkdata.expectedResult, valid)
		}
	}
}
