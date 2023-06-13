// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"accelbyte.net/rotating-shop-items-cli/config"
	"accelbyte.net/rotating-shop-items-cli/pkg/client/platformservice"
)

const DefaultCategoryPath = "/goRotatingShopItemsDemoCli"

const (
	FlagNamespace     = "namespace"
	FlagBaseUrl       = "baseURL"
	FlagClientId      = "clientId"
	FlagClientSecret  = "clientSecret"
	FlagUsername      = "username"
	FlagPassword      = "password"
	FlagCategoryPath  = "categoryPath"
	FlagGrpcUrl       = "grpcTarget"
	FlagRunMode       = "runMode"
	FlagExtendAppName = "extendAppName"
)

func getRotatingShopItemHandler(c *cli.Context, cfg *config.Config) error {
	clientSvc, err := platformservice.NewClient(c.String(FlagBaseUrl), &tokenRepo)
	if err != nil {
		return err
	}
	platformClientSvc = clientSvc

	userInfo, err := TokenGrantV3(c)
	if err != nil {
		return err
	}

	log.Infof("getting grpc server response... ")

	log.Infof("Configuring platform service grpc target... ")
	err = SetPlatformServiceGrpcTarget(c)
	if err != nil {
		return err
	}

	log.Infof("Creating store... ")
	storeId, err := CreateStore(c)
	if err != nil {
		return err
	}

	defer func() {
		DeleteStore(c, storeId)
	}()

	log.Infof("Creating category... ")
	err = CreateCategory(c, storeId)
	if err != nil {
		return err
	}

	log.Infof("Creating store view... ")
	viewId, err := CreateStoreView(c, storeId)
	if err != nil {
		return err
	}

	log.Infof("Creating section with Items...")
	sInfo, err := CreateSectionsWithItems(c, storeId, viewId, 10, true)
	if err != nil {
		return err
	}

	if c.String(FlagRunMode) == "backfill" {
		log.Infof("Enabling custom backfill for section... ")
		err = enableFixedRotationWithCustomBackfillForSection(c, storeId, sInfo.ID, true)
		if err != nil {
			return err
		}
	} else {
		log.Infof("Enabling custom rotation for section... ")
		err = enableCustomRotationForSection(c, storeId, sInfo.ID, true)
		if err != nil {
			return err
		}
	}

	log.Infof("Getting active sections's rotation Items... ")
	resp, err := GetSectionRotationItems(c, Val(userInfo.UserID), viewId)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

func GetCLIApp(cfg *config.Config) *cli.App {
	dockerCmd := &cli.Command{
		Name:  "rotatingShop",
		Usage: "command to invoke the rotating shop grpc server through cli",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    FlagNamespace,
				Usage:   "Namespace",
				Aliases: []string{"n"},
				EnvVars: []string{config.EnvNamespace},
			},
			&cli.StringFlag{
				Name:    FlagBaseUrl,
				Usage:   "Base URL",
				Aliases: []string{"b"},
				EnvVars: []string{config.EnvBaseURL},
			},
			&cli.StringFlag{
				Name:    FlagClientId,
				Usage:   "Client Id",
				Aliases: []string{"i"},
				EnvVars: []string{config.EnvClientID},
			},
			&cli.StringFlag{
				Name:    FlagClientSecret,
				Usage:   "Client Secret",
				Aliases: []string{"s"},
				EnvVars: []string{config.EnvClientSecret},
			},
			&cli.StringFlag{
				Name:    FlagUsername,
				Usage:   "Username",
				Aliases: []string{"u"},
				EnvVars: []string{config.EnvUsername},
			},
			&cli.StringFlag{
				Name:    FlagPassword,
				Usage:   "Password",
				Aliases: []string{"p"},
				EnvVars: []string{config.EnvPassword},
			},
			&cli.StringFlag{
				Name:    FlagCategoryPath,
				Usage:   "Category Path. example `/` or `/customitemrotationtest` ",
				Aliases: []string{"a"},
				EnvVars: []string{config.EnvCategoryPath},
				Value:   DefaultCategoryPath,
			},
			&cli.StringFlag{
				Name:    FlagGrpcUrl,
				Usage:   "GRPC Target",
				Aliases: []string{"g"},
				EnvVars: []string{config.EnvGRPCServerURL},
			},
			&cli.StringFlag{
				Name:    FlagExtendAppName,
				Usage:   "Extend App Name",
				Aliases: []string{"e"},
				EnvVars: []string{config.EnvExtendAppName},
			},
			&cli.StringFlag{
				Name:     FlagRunMode,
				Usage:    "either `backfill` or `rotation`",
				Aliases:  []string{"m"},
				EnvVars:  []string{config.EnvRunMode},
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			// TODO test
			//return getRotatingShopItemHandler(c, cfg)

			fmt.Println("Variables")
			fmt.Printf("%s: %s\n", FlagNamespace, c.String(FlagNamespace))
			fmt.Printf("%s: %s\n", FlagBaseUrl, c.String(FlagBaseUrl))
			fmt.Printf("%s: %s\n", FlagClientId, c.String(FlagClientId))
			fmt.Printf("%s: %s\n", FlagClientSecret, c.String(FlagClientSecret))
			fmt.Printf("%s: %s\n", FlagUsername, c.String(FlagUsername))
			fmt.Printf("%s: %s\n", FlagPassword, c.String(FlagPassword))
			fmt.Printf("%s: %s\n", FlagCategoryPath, c.String(FlagCategoryPath))
			fmt.Printf("%s: %s\n", FlagGrpcUrl, c.String(FlagGrpcUrl))
			fmt.Printf("%s: %s\n", FlagExtendAppName, c.String(FlagExtendAppName))
			fmt.Printf("%s: %s\n", FlagRunMode, c.String(FlagRunMode))
			return getRotatingShopItemHandler(c, cfg)
		},
	}

	return &cli.App{
		Name:  "rotating-shop-Items-cli",
		Usage: fmt.Sprintf("AccelByte rotating shop item CLI tools (default base url: %s)", cfg.AccelByteBaseURL),
		Commands: []*cli.Command{
			dockerCmd,
		},
	}
}
