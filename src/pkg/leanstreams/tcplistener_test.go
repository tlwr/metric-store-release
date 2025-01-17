package leanstreams_test

import (
	"strconv"
	"testing"

	. "github.com/cloudfoundry/metric-store-release/src/pkg/leanstreams"
)

func TestListenTCPUsesDefaultMessageSize(t *testing.T) {
	cfg := TCPListenerConfig{
		Address:   FormatAddress("", strconv.Itoa(5031)),
		Callback:  func([]byte) error { return nil },
		TLSConfig: tlsConfig,
	}
	buffM, err := ListenTCP(cfg)
	if err != nil {
		t.Errorf("Could not Listen on port %d: %s", cfg.MaxMessageSize, err.Error())
	}
	if buffM.ConnConfig.MaxMessageSize != DefaultMaxMessageSize {
		t.Errorf("Expected Max Message Size to be %d, actually got %d", DefaultMaxMessageSize, buffM.ConnConfig.MaxMessageSize)
	}
}

func TestListenTCPUsesSpecifiedMaxMessageSize(t *testing.T) {
	cfg := TCPListenerConfig{
		MaxMessageSize: 8196,
		Address:        FormatAddress("", strconv.Itoa(5032)),
		Callback:       func([]byte) error { return nil },
		TLSConfig:      tlsConfig,
	}
	buffM, err := ListenTCP(cfg)
	if err != nil {
		t.Errorf("Could not Listen on port %d: %s", cfg.MaxMessageSize, err.Error())
	}
	if buffM.ConnConfig.MaxMessageSize != cfg.MaxMessageSize {
		t.Errorf("Expected Max Message Size to be %d, actually got %d", cfg.MaxMessageSize, buffM.ConnConfig.MaxMessageSize)
	}
}
