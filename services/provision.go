package services

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"
	"sync"
)

type ProvisionService struct {
	devicesPrefix string
	wgDir         string
	mu            sync.Mutex
}

func NewProvisionService(devicesPrefix, wgDir string) *ProvisionService {
	return &ProvisionService{devicesPrefix: devicesPrefix, wgDir: wgDir}
}

func (s *ProvisionService) Provision() (configBytes []byte, deviceName string, token string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, err := s.loadDevices()
	if err != nil {
		return nil, "", "", fmt.Errorf("load devices: %w", err)
	}

	deviceName, err = s.generateUniqueName(existing)
	if err != nil {
		return nil, "", "", fmt.Errorf("generate name: %w", err)
	}

	cmd := exec.Command("./easy-wg-quick", deviceName)
	cmd.Dir = s.wgDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, "", "", fmt.Errorf("easy-wg-quick: %s: %w", string(out), err)
	}

	confPath := path.Join(s.wgDir, "wgclient_"+deviceName+".conf")
	configBytes, err = os.ReadFile(confPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("read config: %w", err)
	}

	if err := s.appendDevice(deviceName); err != nil {
		return nil, "", "", fmt.Errorf("append device: %w", err)
	}

	tokenBytes := make([]byte, 32)
	if _, err := crand.Read(tokenBytes); err != nil {
		return nil, "", "", fmt.Errorf("generate token: %w", err)
	}
	token = hex.EncodeToString(tokenBytes)

	return configBytes, deviceName, token, nil
}

func (s *ProvisionService) loadDevices() ([]string, error) {
	devicesFile := path.Join(s.wgDir, "devices.txt")
	data, err := os.ReadFile(devicesFile)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), nil
}

func (s *ProvisionService) generateUniqueName(existing []string) (string, error) {
	for {
		randomNum := rand.IntN(900) + 100
		name := fmt.Sprintf("%s%d", s.devicesPrefix, randomNum)
		if !slices.Contains(existing, name) {
			return name, nil
		}
	}
}

func (s *ProvisionService) appendDevice(deviceName string) error {
	f, err := os.OpenFile(path.Join(s.wgDir, "devices.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, deviceName)
	return err
}
