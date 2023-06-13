// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AccelByte/accelbyte-go-sdk/iam-sdk/pkg/iamclient/users"
	"github.com/AccelByte/accelbyte-go-sdk/iam-sdk/pkg/iamclientmodels"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/factory"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/utils/auth"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	configRepo = *auth.DefaultConfigRepositoryImpl()
	tokenRepo  = *auth.DefaultTokenRepositoryImpl()
)

func TokenGrantV3(c *cli.Context) (*iamclientmodels.ModelUserResponseV3, error) {
	err := os.Setenv("AB_CLIENT_ID", c.String(FlagClientId))
	err = os.Setenv("AB_CLIENT_SECRET", c.String(FlagClientSecret))
	err = os.Setenv("AB_BASE_URL", c.String(FlagBaseUrl))

	oauth := &iam.OAuth20Service{
		Client:           factory.NewIamClient(&configRepo),
		ConfigRepository: &configRepo,
		TokenRepository:  &tokenRepo,
	}

	err = oauth.LoginUser(c.String(FlagUsername), c.String(FlagPassword))
	if err != nil {
		logrus.Print("failed login user.")
	} else {
		logrus.Print("successful login.")
	}

	usersService := &iam.UsersService{
		Client:           factory.NewIamClient(&configRepo),
		ConfigRepository: &configRepo,
		TokenRepository:  &tokenRepo,
	}
	userInfo, err := usersService.PublicGetMyUserV3Short(&users.PublicGetMyUserV3Params{})
	if err != nil {
		log.Fatalf("Get user info failed: %s\n", err)
	}
	fmt.Printf("\tUser: %s %s\n", userInfo.UserName, Val(userInfo.UserID))

	return userInfo, nil
}
