package config

type RaftConfig struct {
	Config struct {
		TermInfoPath string `json:"term_info_path" yaml:"term_info_path"`
		LogPath      string `json:"log_path" yaml:"log_path"`
		Servers      []struct {
			Port  uint32 `json:"port" yaml:"port"`
			Local bool   `json:"local" yaml:"local"`
		} `json:"servers" yaml:"servers"`
	} `json:"config" yaml:"config"`
}
