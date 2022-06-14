package system

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

type AuthConfig interface {
	IsEnforcingIP() bool
	IsKnownIP(ip string) bool
	IsEnforcingKey() bool
	GetKey() string
}

type MapConfig interface {
	VerifyMapName(name string, calculated uint32) bool
}

type BanConfig interface {
	IsSteamIdBanned(steamId string) bool
}

// Validate ApiConfig implements the following interfaces
var (
	_ AuthConfig = (*ApiConfig)(nil)
	_ MapConfig  = (*ApiConfig)(nil)
	_ BanConfig  = (*ApiConfig)(nil)
)

type ApiConfig struct {
	Core struct {
		DebugMode  bool
		Address    string
		Port       int
		MaxThreads int
		Graceful   time.Duration
		RootPath   string
		DBString   string
	}
	RateLimit struct {
		Enable      bool
		MaxRequests int
		MaxAge      time.Duration
	}
	Cert struct {
		Enable bool
		Domain string
	}
	ApiAuth struct {
		EnforceKey bool
		EnforceIP  bool
		Key        string
		IPListFile string
	}
	Verify struct {
		EnforceBan    bool
		EnforceMap    bool
		EnforceSC     bool
		MapListFile   string
		BanListFile   string
		AdminListFile string
		SCHash        uint32
	}
	Log struct {
		Level      string
		Dir        string
		ExpireTime string
	}
	iPList         map[string]bool
	iPListMutex    sync.RWMutex
	banList        map[string]bool
	banListMutex   sync.RWMutex
	mapList        map[string]uint32
	mapListMutex   sync.RWMutex
	adminList      map[string]bool
	adminListMutex sync.RWMutex
}

func LoadApiConfig(path string, dbg bool) (*ApiConfig, error) {
	var cfg ApiConfig

	switch filepath.Ext(path) {
	case ".toml", ".ini":
		if err := ini.MapTo(&cfg, path); err != nil {
			return nil, err
		}
	case ".yaml", ".json":
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported config type")
	}

	cfg.Core.DebugMode = dbg
	return &cfg, nil
}

func (c *ApiConfig) LoadIPList() error {
	c.iPListMutex.Lock()
	defer c.iPListMutex.Unlock()
	return loadJsonFile(c.ApiAuth.IPListFile, &c.iPList)
}

func (c *ApiConfig) LoadMapList() error {
	c.mapListMutex.Lock()
	defer c.mapListMutex.Unlock()
	return loadJsonFile(c.Verify.MapListFile, &c.mapList)
}

func (c *ApiConfig) LoadBanList() error {
	c.banListMutex.Lock()
	defer c.banListMutex.Unlock()
	return loadJsonFile(c.Verify.BanListFile, &c.banList)
}

func (c *ApiConfig) LoadAdminList() error {
	c.adminListMutex.Lock()
	defer c.adminListMutex.Unlock()
	return loadJsonFile(c.Verify.AdminListFile, &c.adminList)
}

func loadJsonFile(path string, container interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(file), container)
}

func (c *ApiConfig) IsEnforcingIP() bool {
	return c.ApiAuth.EnforceIP
}

func (c *ApiConfig) IsKnownIP(ip string) bool {
	c.iPListMutex.RLock()
	defer c.iPListMutex.RUnlock()

	_, ok := c.iPList[ip]
	return ok
}

func (c *ApiConfig) IsEnforcingKey() bool {
	return c.ApiAuth.EnforceKey
}

func (c *ApiConfig) GetKey() string {
	return c.ApiAuth.Key
}

func (c *ApiConfig) VerifyMapName(name string, calculated uint32) bool {
	c.mapListMutex.RLock()
	defer c.mapListMutex.RUnlock()

	return c.mapList[name] == calculated
}

func (c *ApiConfig) IsSteamIdAdmin(steamId string) bool {
	c.adminListMutex.RLock()
	defer c.adminListMutex.RUnlock()

	_, ok := c.adminList[steamId]
	return ok
}

func (c *ApiConfig) IsSteamIdBanned(steamId string) bool {
	c.banListMutex.RLock()
	defer c.banListMutex.RUnlock()

	_, ok := c.banList[steamId]
	return ok
}

func (c *ApiConfig) EnforceAndVerifyBanned(steamId string) bool {
	if c.Verify.EnforceBan {
		return c.IsSteamIdBanned(steamId)
	}
	return false
}
