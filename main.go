package main

import (
	"api-gateway/src/config/builder"
	"api-gateway/src/config/envs"
)

func main() {
	env := envs.LoadEnvs(".env")

	app := builder.Build(env)

	app.Run(":" + env.Get("PORT"))
}
