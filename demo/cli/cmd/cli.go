// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package cmd

import (
	"encoding/json"
	"fmt"

	"accelbyte.net/rotating-shop-items-cli/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	FlagNamespace    = "namespace"
	FlagBaseUrl      = "baseURL"
	FlagClientId     = "clientId"
	FlagClientSecret = "clientSecret"
	FlagUsername     = "username"
	FlagPassword     = "password"
	FlagCategoryPath = "categoryPath"
	FlagGrpcUrl      = "grpcTarget"
	FlagRunMode      = "runMode"
)

var (
	userId string
)

func getRotatingShopItemHandler(c *cli.Context, cfg *config.Config) error {
	if cfg.AccelByteAccessToken == "" {
		token, err := TokenGrantV3(c)
		if err != nil {
			return err
		}

		if token.UserID != "" {
			userId = token.UserID
		} else {
			userId = "123test"
		}

	}

	log.Infof("getting grpc server response... ")

	log.Infof("Configuring platform service grpc target... ")
	err := SetPlatformServiceGrpcTarget(c)
	if err != nil {
		return err
	}

	log.Infof("Creating store... ")
	storeId, err := CreateStore(c)
	if err != nil {
		return err
	}

	defer DeleteStore(c, storeId)

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

	log.Infof("Creating section with items...")
	sInfo, err := CreateSectionsWithItems(c, storeId, viewId, 10, true)
	if err != nil {
		return err
	}

	if c.String(FlagRunMode) == "backfill" {
		log.Infof("Enabling custom backfill for section... ")
		err = enableFixedRotationWithCustomBackfillForSection(c, storeId, sInfo.id, true)
		if err != nil {
			return err
		}
	} else {
		log.Infof("Enabling custom rotation for section... ")
		err = enableCustomRotationForSection(c, storeId, sInfo.id, true)
		if err != nil {
			return err
		}
	}

	log.Infof("Getting active sections's rotation items... ")
	resp, err := GetSectionRotationItems(c, userId, viewId)
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
			},
			&cli.StringFlag{
				Name:    FlagBaseUrl,
				Usage:   "Base URL",
				Aliases: []string{"b"},
			},
			&cli.StringFlag{
				Name:    FlagClientId,
				Usage:   "Client Id",
				Aliases: []string{"i"},
			},
			&cli.StringFlag{
				Name:    FlagClientSecret,
				Usage:   "Client Secret",
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:    FlagUsername,
				Usage:   "Username",
				Aliases: []string{"u"},
			},
			&cli.StringFlag{
				Name:    FlagPassword,
				Usage:   "Password",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    FlagCategoryPath,
				Usage:   "Category Path. example `/` or `/customitemrotationtest` ",
				Aliases: []string{"a"},
			},
			&cli.StringFlag{
				Name:    FlagGrpcUrl,
				Usage:   "GRPC Target",
				Aliases: []string{"g"},
			},
			&cli.StringFlag{
				Name:    FlagRunMode,
				Usage:   "either `backfill` or `rotation`",
				Aliases: []string{"m"},
			},
		},
		Action: func(c *cli.Context) error {
			return getRotatingShopItemHandler(c, cfg)
		},
	}

	return &cli.App{
		Name:  "rotating-shop-items-cli",
		Usage: fmt.Sprintf("AccelByte rotating shop item CLI tools (default base url: %s)", cfg.AccelByteBaseURL),
		Commands: []*cli.Command{
			dockerCmd,
		},
	}
}
