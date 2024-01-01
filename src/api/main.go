package main

import "prospector/cmd"

//	@title						Prospector API
//	@description				API backend for Prospector
//	@version					1.0
//	@host						prospector.ie
//	@schemes					https
//	@BasePath					/api
//	@securityDefinitions.basic	BasicAuth

func main() {
	cmd.Execute()
}
