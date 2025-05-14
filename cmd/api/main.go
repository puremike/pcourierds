package main

import (
	_ "github.com/puremike/pcourierds/docs"
	"github.com/puremike/pcourierds/internal/env"
	"go.uber.org/zap"
)

var envGet *env.EnvConfig

//	@title			Courier Delivery System API
//	@version		1.0
//	@description	This is an API for a Courier Delivery System
// @contact.name   Puremike
// @contact.url    http://github.com/puremike
// @contact.email  digitalmarketfy@gmail.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//	@BasePath		/v1

func main() {

	envGet = env.GetEnv()

	cfg := &config{
		port: envGet.PORT,
		env:  envGet.ENV,
	}

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	app := &app{
		config: cfg,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.start(mux))
}
