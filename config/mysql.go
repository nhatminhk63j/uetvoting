package config

import "fmt"

// MysqlConfig define mysql config.
type MysqlConfig struct {
	Schema   string `name:"mysql-schema" help:"Mysql schema" env:"MYSQL_SCHEMA" default:"uet_voting"`
	Host     string `name:"mysql-host" help:"Mysql host" env:"MYSQL_HOST" default:"localhost"`
	Port     int    `name:"mysql-port" help:"Mysql port" env:"MYSQL_PORT" default:"3306"`
	User     string `name:"mysql-auth" help:"Mysql auth" env:"MYSQL_USER" default:"uet"`
	Password string `name:"mysql-password" help:"Mysql password" env:"MYSQL_PASSWORD" default:"password"`

	MaxIdleConns    int `name:"max-idle-conns" env:"MAX_IDLE_CONNS" default:"3" help:"maximum number of connections in the idle connection pool"`
	MaxOpenConns    int `name:"max-open-conns" env:"MAX_OPEN_CONNS" default:"10" help:"maximum number of open connections to the database"`
	ConnMaxLifetime int `name:"conn-max-life-time" env:"CONN_MAX_LIFE_TIME" default:"60" help:"maximum amount of second a connection may be reused."`
}

// ToURI ...
func (cfg MysqlConfig) ToURI() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Schema)
}
