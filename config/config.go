// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type EtcdbeatConfig struct {
	Period *int64
	Port   *string
	Host   *string

	// Authentication for BasicAuth
	Authentication struct {
		Enable   *bool
		Username *string
		Password *string
	}

	Statistics struct {
		Leader *bool
		Self   *bool
		Store  *bool
	}
}

type ConfigSettings struct {
	Input EtcdbeatConfig
}
