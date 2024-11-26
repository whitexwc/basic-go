//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:11309)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:10379",
	},
}
