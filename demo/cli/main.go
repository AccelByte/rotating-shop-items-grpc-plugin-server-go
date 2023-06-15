// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"accelbyte.net/rotating-shop-items-cli/cmd"
	"accelbyte.net/rotating-shop-items-cli/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cfg := config.Get()
	app := cmd.GetCLIApp(cfg)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
