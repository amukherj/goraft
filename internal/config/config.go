package config

type RaftConfig struct {
	Config struct {
		ServerID     *int32 `json:"-" yaml:"-"`
		TermInfoPath string `json:"term_info_path" yaml:"term_info_path"`
		LogPath      string `json:"log_path" yaml:"log_path"`
		Servers      []struct {
			Port  uint32 `json:"port" yaml:"port"`
			Local bool   `json:"local" yaml:"local"`
		} `json:"servers" yaml:"servers"`
	} `json:"config" yaml:"config"`
}

func (rc *RaftConfig) GetServerID() int32 {
	if rc.Config.ServerID == nil {
		for n, svr := range rc.Config.Servers {
			if svr.Local {
				id := int32(n + 1)
				rc.Config.ServerID = &id
				break
			}
		}
	}
	if rc.Config.ServerID != nil {
		return *rc.Config.ServerID
	}
	return 0
}
