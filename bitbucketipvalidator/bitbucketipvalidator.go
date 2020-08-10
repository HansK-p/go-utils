package bitbucketipvalidator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	layout      = "2006-01-02T15:04:05.000000"
	ipRangesURL = "https://ip-ranges.atlassian.com"
)

type timestamp struct {
	time.Time
}

func (ts *timestamp) UnmarshalJSON(b []byte) error {
	fmt.Printf("Unmarshal timestamp %v\n", string(b))
	// Convert to string and remove quotes
	s := strings.Trim(string(b), "\"")

	// Parse the time using the layout
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}

	// Assign the parsed time to our variable
	ts.Time = t
	return nil
}

type cidr struct {
	net.IPNet
}

func (cidr *cidr) UnmarshalJSON(b []byte) error {
	// Convert to string and remove quotes
	s := strings.Trim(string(b), "\"")

	_, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}

	cidr.IP = ipNet.IP
	cidr.Mask = ipNet.Mask
	return nil
}

type ipRange struct {
	Network string
	MaskLen int  `json:"mask_len,omitempty"`
	Cidr    cidr // net.IPNet
	Mask    string
}

type iPRanges struct {
	CreationDate timestamp
	SyncToken    int64
	IPRanges     []ipRange `json:"items"`
}

// BitbucketIPValidator will check if a given IP range is within the Bitbucket public IP range
type BitbucketIPValidator struct {
	logger   *log.Entry
	ipRanges *iPRanges
}

// NewBitbucketIPValidator will initialize a BitbucketIPValidator and return it
func NewBitbucketIPValidator(logger *log.Logger) (*BitbucketIPValidator, error) {
	entry := logger.WithFields(log.Fields{"module": "BitbucketIPValidator"})
	entry.Infof("Downloading IP rules from %v", ipRangesURL)
	res, err := http.Get(ipRangesURL)
	if err != nil {
		entry.Warnf("Error downloading from URL %v: %v", ipRangesURL, err)
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		entry.Warnf("Error downloading body from URL %v: %v", ipRangesURL, err)
		return nil, err
	}
	entry.Tracef("IP rules downloaded: %s\n", body)
	var ipRanges iPRanges
	err = json.Unmarshal([]byte(body), &ipRanges)
	if err != nil {
		entry.Warnf("Error unmarshaling ip ranged downloaded from %v: %v", ipRangesURL, err)
		return nil, err
	}

	biv := BitbucketIPValidator{
		logger:   entry,
		ipRanges: &ipRanges,
	}
	return &biv, nil
}

// ValidateIP Will validate if the given IP address is a Bitbucket public IP address
func (biv *BitbucketIPValidator) ValidateIP(ip net.IP) bool {
	isABitbucketIP := false
	for _, ipRange := range biv.ipRanges.IPRanges {
		if ipRange.Cidr.Contains(ip) {
			biv.logger.Infof("IP %v matches Cidr %v", ip, ipRange.Cidr)
			isABitbucketIP = true
		} else {
			biv.logger.Debugf("IP %v does not match cidr %v", ip, ipRange.Cidr)
		}
	}
	return isABitbucketIP
}
