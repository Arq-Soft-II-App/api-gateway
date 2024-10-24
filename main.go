package main

import (
	envs "api-gateway/src/config/envs"
)

func main() {
	env := envs.LoadEnvs(".env")

	app := builder

}
