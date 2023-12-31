// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package config

import (
	"github.com/caarlos0/env"
)

// note: this constant can be updated by Makefile before build
const defaultAccelByteBaseURL = "https://development.accelbyte.io"

const (
	EnvNamespace     = "AB_NAMESPACE"
	EnvBaseURL       = "AB_BASE_URL"
	EnvClientID      = "AB_CLIENT_ID"
	EnvClientSecret  = "AB_CLIENT_SECRET"
	EnvUsername      = "AB_USERNAME"
	EnvPassword      = "AB_PASSWORD"
	EnvGRPCServerURL = "GRPC_SERVER_URL"
	EnvExtendAppName = "EXTEND_APP_NAME"
	EnvRunMode       = "RUN_MODE"
	EnvCategoryPath  = "CATEGORY_PATH"
)

type Config struct {
	Namespace             string
	Username              string `env:"AB_USERNAME"`
	Password              string `env:"AB_PASSWORD"`
	AccelByteClientID     string `env:"AB_CLIENT_ID"`
	AccelByteClientSecret string `env:"AB_CLIENT_SECRET"`
	AccelByteBaseURL      string `env:"AB_BASE_URL"`
	AccelByteAccessToken  string `env:"AB_ACCESS_TOKEN" envDocs:"If not empty will be used as access token and no need to login"`
}

func Get() *Config {
	cfg := &Config{
		AccelByteBaseURL: defaultAccelByteBaseURL,
	}
	_ = env.Parse(cfg)

	return cfg
}
