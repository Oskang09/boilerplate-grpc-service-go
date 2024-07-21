package bootstrap

import (
	"context"
	"net"
	"os"

	"service/app/config"

	"github.com/RevenueMonster/sqlike/plugin/opentracing"
	"github.com/RevenueMonster/sqlike/sql/instrumented"
	"github.com/RevenueMonster/sqlike/sqlike"
	"github.com/RevenueMonster/sqlike/sqlike/options"
	"github.com/go-sql-driver/mysql"
)

func (bs *Bootstrap) initMySQL() *Bootstrap {

	if config.IsProduction() {
		client := sqlike.MustConnect(context.Background(), "mysql",
			options.Connect().
				SetHost(config.DBHost).
				SetPort(config.DBPort).
				SetUsername(config.DBUser).
				SetPassword(config.DBPassword),
		)
		client.SetPrimaryKey("Key")

		bs.Database = client.Database(config.DBName)
	} else {
		dbConfig := mysql.NewConfig()
		dbConfig.User = config.DBUser
		dbConfig.Passwd = config.DBPassword
		dbConfig.ParseTime = true
		dbConfig.Net = "tcp"
		dbConfig.Addr = net.JoinHostPort(config.DBHost, config.DBPort)
		connector, err := mysql.NewConnector(dbConfig)
		if err != nil {
			panic(err)
		}

		itpr := opentracing.NewInterceptor(
			opentracing.WithDBInstance("sqlike"),
			opentracing.WithDBUser(dbConfig.User),
			opentracing.WithQuery(true),
		)
		client := sqlike.MustConnectDB(context.Background(), "mysql", instrumented.WrapConnector(connector, itpr))
		client.SetPrimaryKey("Key")

		bs.Database = client.Database(config.DBName)
	}

	os.Stdout.WriteString("Connected to MySQL!\n")
	return bs
}
