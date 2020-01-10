package replication

import (
	"fmt"
	"net"
	"regexp"

	"github.com/rs/zerolog/log"
)

// Address is an IP address or a hostname of a FReD node.
type Address struct {
	Addr string
	IsIP bool
}

// ParseAddress parses a string and returns it as an Address if it is a valid address.
// It returns true if the address is an IP address (opposed to a hostname).
// If it cannot parse the address, it returns an error != nil.
// https://stackoverflow.com/questions/42479410/how-to-determine-if-given-string-is-a-hostname-or-an-ip-address
func ParseAddress(a string) (Address, error) {
	ip := net.ParseIP(a)
	if ip != nil {
		return Address{
			Addr: string(ip),
			IsIP: true,
		}, nil
	}

	// i hate regex
	// https://stackoverflow.com/questions/3026957/how-to-validate-a-domain-name-using-regex-php
	matched, err := regexp.MatchString(`^(?:[-A-Za-z0-9]+\.)+[A-Za-z]{2,}$`, a)

	if err != nil {
		log.Err(err).Msg("")
		return Address{}, err
	}

	if !matched {
		return Address{}, fmt.Errorf("replication.address: could not validate %s as IP address or hostname", a)
	}

	return Address{
		Addr: a,
		IsIP: false,
	}, nil
}