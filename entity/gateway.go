package entity

type GatewayConfig struct {
	ListenAddr string  `mapstructure:"listenAddr"`
	MySql      string  `mapstructure:"mysql"`
	Routes     []Route `mapstructure:"routes"`
}
