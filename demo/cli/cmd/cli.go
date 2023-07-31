// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/section"
	"github.com/urfave/cli/v2"

	"accelbyte.net/rotating-shop-items-cli/config"
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

func runRotatingShopDemo(c *cli.Context, cfg *config.Config) error {
	fmt.Printf("Login to AccelByte... ")
	userInfo, err := TokenGrantV3(c)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")
	fmt.Printf("\tUser: %s %s\n", userInfo.UserName, Val(userInfo.UserID))

	fmt.Printf("Configuring platform service grpc target... ")
	err = SetPlatformServiceGrpcTarget(c)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	fmt.Printf("Creating store... ")
	storeId, err := CreateStore(c)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")
	fmt.Printf("\tStoreID: %s\n", storeId)

	fmt.Printf("Setting up currency (%s)... ", currencyCode)
	err = CreateCurrency(c)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	defer func() {
		fmt.Print("Deleting currency... ")
		if err := DeleteCurrency(c); err != nil {
			return
		}
		fmt.Println("[OK]")

		fmt.Print("Deleting store...")
		if _, err := DeleteStore(c, storeId); err != nil {
			return
		}
		fmt.Println("[OK]")

	}()

	fmt.Printf("Creating category %s... ", c.String(FlagCategoryPath))
	err = CreateCategory(c, storeId)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	fmt.Printf("Creating store view... ")
	viewId, err := CreateStoreView(c, storeId)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")
	fmt.Printf("\tViewID: %s\n", viewId)

	itemCount := 5
	fmt.Printf("Creating section with %d items... ", itemCount)
	sectionInfo, items, err := CreateSectionsWithItems(c, storeId, viewId, itemCount, true)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")
	fmt.Printf("\tSectionID: %s\n", sectionInfo.ID)

	// --------
	fmt.Println("[Testing BackFill]")
	fmt.Printf("Granting entitlement for item: %s to user: %s... ", items[0].ID, Val(userInfo.UserID))
	_, err = GrantEntitlement(c, storeId, Val(userInfo.UserID), items[0].ID, 1)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	fmt.Printf("Enabling custom backfill for section %s... ", sectionInfo.ID)
	err = enableFixedRotationWithCustomBackfillForSection(c, storeId, sectionInfo.ID, true)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	fmt.Printf("Getting active sections's rotation items... ")
	_, err = GetSectionRotationItems(c, Val(userInfo.UserID), viewId)
	// assuming the customer doesn't modify sample app, it should be response with random
	// item id which basically not exists in accelbyte and result in not found error
	_, isNotFoundErr := err.(*section.PublicListActiveSectionsNotFound)
	if err != nil && !isNotFoundErr {
		return err
	}
	fmt.Printf(" expect not found: %t ", isNotFoundErr)
	fmt.Println("[OK]")

	// --------
	fmt.Println("[Testing Custom Rotation Items]")
	fmt.Printf("Enabling custom rotation for section %s... ", sectionInfo.ID)
	err = enableCustomRotationForSection(c, storeId, sectionInfo.ID, true)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	fmt.Printf("Getting active sections's rotation items... ")
	resp, err := GetSectionRotationItems(c, Val(userInfo.UserID), viewId)
	if err != nil {
		return err
	}
	fmt.Println("[OK]")

	jsonData, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("Current rotation response: %s\n", string(jsonData))

	fmt.Println("[SUCCESS]")

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
			return runRotatingShopDemo(c, cfg)
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
