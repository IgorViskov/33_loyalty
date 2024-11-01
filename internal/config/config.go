package config

import "net/url"

type AppConfig struct {
	RunHost             string   `env:"RUN_ADDRESS"`
	DBURI               string   `env:"DATABASE_URI"`
	AccrualHost         *url.URL `env:"ACCRUAL_SYSTEM_ADDRESS"`
	MaxParallelRequests int      `env:"MAX_PARALLEL_REQUESTS"`
	PeriodRequests      int      `env:"PERIOD_REQUESTS"`
}

func HostNameParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		conf.RunHost = flagValue
		return nil
	}
}

func AccrualHostParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		u, err := url.Parse(flagValue)
		if err != nil {
			return err
		}
		conf.AccrualHost = u
		return nil
	}
}

func DBURIParser(conf *AppConfig) func(flagValue string) error {
	return func(flagValue string) error {
		conf.DBURI = flagValue
		return nil
	}
}
