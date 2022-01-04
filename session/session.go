package session

import(
  "time"

  "github.com/BurntSushi/toml"
)

var (
  Config config
)

type config struct {
  Core struct {
    WebPort int
    IP string
    MaxThreads int
    MaxRequests int
    MaxAge time.Duration
    Graceful time.Duration
    RootPath string
  }
  TLS struct {
    Enable bool
    Port int
    CertFile string
    KeyFile string
  }
  ApiAuth struct {
    EnforceKey bool
    EnforceIP bool
    Key string
    IPAllowed map[string]int8
  }
  Verify struct {
    EnforceBan bool
    EnforceMap bool
    BanList map[string]int64
    MapList map[string]uint32
    SCHash uint32
  }
  Log struct {
    Dir string
  }
}

func LoadConfig(path string) error {
  _, err := toml.DecodeFile(path, &Config);
  if err != nil {
    return err
  }

  return nil
}