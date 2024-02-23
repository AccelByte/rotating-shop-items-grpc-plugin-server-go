// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/catalog_changes"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/category"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/currency"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/entitlement"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/item"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/section"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/service_plugin_config"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/store"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/view"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclientmodels"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/factory"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/platform"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/tests/integration"
	"github.com/go-openapi/strfmt"
	"github.com/urfave/cli/v2"
)

var (
	storeService = &platform.StoreService{
		Client:          factory.NewPlatformClient(&configRepo),
		TokenRepository: &tokenRepo,
	}
	categoryService = &platform.CategoryService{
		Client:          factory.NewPlatformClient(&configRepo),
		TokenRepository: &tokenRepo,
	}

	viewService = &platform.ViewService{
		Client:          factory.NewPlatformClient(&configRepo),
		TokenRepository: &tokenRepo,
	}

	catalogChangesService = &platform.CatalogChangesService{
		Client:          factory.NewPlatformClient(&configRepo),
		TokenRepository: &tokenRepo,
	}
	itemService = &platform.ItemService{
		Client:          factory.NewPlatformClient(&configRepo),
		TokenRepository: &tokenRepo,
	}

	sectionService = &platform.SectionService{
		Client:          factory.NewPlatformClient(&configRepo),
		TokenRepository: &tokenRepo,
	}

	abStoreName = "Item Rotation Plugin Demo Store"
	sLangs      = []string{"en"}
	sRegions    = []string{"US"}

	abViewName   = "Go Item Rotation Default View Demo/CLI 006"
	displayOrder = int32(1)

	currencyCode = "USDT1"
)

const AlphaChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const durationTwoDays = time.Hour * 24 * 2

type SimpleItemInfo struct {
	ID    string
	SKU   string
	Title string
}

type SimpleSectionInfo struct {
	ID    string
	Items []*SimpleItemInfo
}

func SetPlatformServiceGrpcTarget(c *cli.Context) error {
	servicePluginCfgWrapper := platform.ServicePluginConfigService{
		Client:           factory.NewPlatformClient(&configRepo),
		ConfigRepository: &configRepo,
		TokenRepository:  &tokenRepo,
	}

	if c.String(FlagGrpcUrl) != "" {
		fmt.Printf("(Custom Host: %s) ", c.String(FlagGrpcUrl))

		_, err := servicePluginCfgWrapper.UpdateSectionPluginConfigShort(&service_plugin_config.UpdateSectionPluginConfigParams{
			Namespace: c.String(FlagNamespace),
			Body: &platformclientmodels.SectionPluginConfigUpdate{
				ExtendType: Ptr(platformclientmodels.SectionPluginConfigUpdateExtendTypeCUSTOM),
				CustomConfig: &platformclientmodels.BaseCustomConfig{
					ConnectionType:    Ptr(platformclientmodels.BaseCustomConfigConnectionTypeINSECURE),
					GrpcServerAddress: Ptr(c.String(FlagGrpcUrl)),
				},
			},
		})

		return err
	}

	if c.String(FlagExtendAppName) != "" {
		fmt.Printf("(Extend App: %s) ", c.String(FlagExtendAppName))

		_, err := servicePluginCfgWrapper.UpdateSectionPluginConfigShort(&service_plugin_config.UpdateSectionPluginConfigParams{
			Namespace: c.String(FlagNamespace),
			Body: &platformclientmodels.SectionPluginConfigUpdate{
				ExtendType: Ptr(platformclientmodels.SectionPluginConfigUpdateExtendTypeAPP),
				AppConfig: &platformclientmodels.AppConfig{
					AppName: Ptr(c.String(FlagExtendAppName)),
				},
			},
		})

		return err
	}

	return nil
}

func CreateCurrency(c *cli.Context) error {
	currencyWrapper := platform.CurrencyService{
		Client:           factory.NewPlatformClient(&configRepo),
		ConfigRepository: &configRepo,
		TokenRepository:  &tokenRepo,
	}
	curr, err := currencyWrapper.GetCurrencySummaryShort(&currency.GetCurrencySummaryParams{
		CurrencyCode: currencyCode,
		Namespace:    c.String(FlagNamespace),
	})
	if err == nil && curr != nil {
		// already exists
		return nil
	}

	_, err = currencyWrapper.CreateCurrencyShort(&currency.CreateCurrencyParams{
		Namespace: c.String(FlagNamespace),
		Body: &platformclientmodels.CurrencyCreate{
			CurrencyCode:   Ptr(currencyCode),
			CurrencySymbol: "USDT1$",
			CurrencyType:   platformclientmodels.CurrencyCreateCurrencyTypeVIRTUAL,
			Decimals:       0,
		},
	})

	return err
}

func DeleteCurrency(c *cli.Context) error {
	currencyWrapper := platform.CurrencyService{
		Client:           factory.NewPlatformClient(&configRepo),
		ConfigRepository: &configRepo,
		TokenRepository:  &tokenRepo,
	}
	_, err := currencyWrapper.DeleteCurrencyShort(&currency.DeleteCurrencyParams{
		Namespace:    c.String(FlagNamespace),
		CurrencyCode: currencyCode,
	})

	return err
}

func CreateStore(c *cli.Context) (string, error) {
	// Clean up existing stores
	storeInfo, err := storeService.ListStoresShort(&store.ListStoresParams{
		Namespace: c.String(FlagNamespace),
	})
	if err != nil {
		return "", err
	}

	for _, s := range storeInfo {
		if !Val(s.Published) {
			_, _ = storeService.DeleteStoreShort(&store.DeleteStoreParams{
				Namespace: c.String(FlagNamespace),
				StoreID:   Val(s.StoreID),
			})
		}
	}

	ok, err := storeService.CreateStoreShort(&store.CreateStoreParams{
		Body: &platformclientmodels.StoreCreate{
			DefaultLanguage:    "en",
			DefaultRegion:      "US",
			Description:        "Description for item rotation grpc plugin demo store",
			SupportedLanguages: sLangs,
			SupportedRegions:   sRegions,
			Title:              &abStoreName,
		},
		Namespace: c.String(FlagNamespace),
	})
	if err != nil {
		return "", err
	}

	return Val(ok.StoreID), nil
}

func CreateCategory(c *cli.Context, storeId string) error {
	localization := make(map[string]string)
	localization["en"] = c.String(FlagCategoryPath)

	_, err := categoryService.CreateCategoryShort(&category.CreateCategoryParams{
		Body: &platformclientmodels.CategoryCreate{
			CategoryPath:             Ptr(c.String(FlagCategoryPath)),
			LocalizationDisplayNames: localization,
		},
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	})
	if err != nil {
		return err
	}

	return nil
}

func CreateStoreView(c *cli.Context, storeId string) (string, error) {
	localization := make(map[string]platformclientmodels.Localization)
	localization["en"] = platformclientmodels.Localization{
		Title: &abViewName,
	}
	ok, err := viewService.CreateViewShort(&view.CreateViewParams{
		Body: &platformclientmodels.ViewCreate{
			DisplayOrder:  &displayOrder,
			Localizations: localization,
			Name:          &abViewName,
		},
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	})
	if err != nil {
		return "", err
	}

	return Val(ok.ViewID), nil
}

func publishStoreChange(storeId string) string {
	inputCreate := &catalog_changes.PublishAllParams{
		StoreID:   storeId,
		Namespace: integration.NamespaceTest,
	}

	created, errCreate := catalogChangesService.PublishAllShort(inputCreate)
	if errCreate != nil {
		return ""
	}
	storeID := *created.StoreID

	return storeID
}

func CreateItems(c *cli.Context, storeId, itemDiff string, itemCount int, doPublish bool) ([]*SimpleItemInfo, error) {
	var nItems []*SimpleItemInfo

	for i := 0; i < itemCount; i++ {
		price := int32((i + 1) * 2)
		var nItemInfo SimpleItemInfo
		nItemInfo.Title = fmt.Sprint("Item " + itemDiff + " Titled " + strconv.FormatInt(int64(i+1), 10))
		nItemInfo.SKU = fmt.Sprint("SKU_" + itemDiff + "_" + strconv.FormatInt(int64(i+1), 10))

		localization := map[string]platformclientmodels.Localization{
			"en": {
				Title: &nItemInfo.Title,
			},
		}
		regionData := map[string][]platformclientmodels.RegionDataItemDTO{
			"US": {
				{
					CurrencyCode:      &currencyCode,
					CurrencyNamespace: Ptr(c.String(FlagNamespace)),
					CurrencyType:      Ptr(platformclientmodels.RegionDataItemCurrencyTypeREAL),
					Price:             &price,
				},
			},
		}

		ok, err := itemService.CreateItemShort(&item.CreateItemParams{
			Body: &platformclientmodels.ItemCreate{
				Features:        []string{"go-demo-cli"},
				Tags:            []string{"tags"},
				Name:            &nItemInfo.Title,
				ItemType:        Ptr(platformclientmodels.ItemCreateItemTypeINGAMEITEM),
				CategoryPath:    Ptr(c.String(FlagCategoryPath)),
				EntitlementType: Ptr(platformclientmodels.ItemCreateEntitlementTypeDURABLE),
				SeasonType:      platformclientmodels.ItemCreateSeasonTypeTIER,
				Status:          Ptr(platformclientmodels.ItemCreateStatusACTIVE),
				Listable:        true,
				Purchasable:     true,
				Sku:             nItemInfo.SKU,
				Localizations:   localization,
				RegionData:      regionData,
			},
			Namespace: c.String(FlagNamespace),
			StoreID:   storeId,
		})
		if err != nil {
			return nil, err
		}

		// add the Items
		nItemInfo.ID = Val(ok.ItemID)
		nItems = append(nItems, &nItemInfo)
	}

	if doPublish {
		publishStoreChange(storeId)
	}

	return nItems, nil
}

func CreateSectionsWithItems(c *cli.Context, storeId, viewId string, itemCount int, doPublish bool) (*SimpleSectionInfo, []*SimpleItemInfo, error) {
	itemDiff := RandomString(AlphaChars, 6)
	items, err := CreateItems(c, storeId, itemDiff, itemCount, doPublish)
	if err != nil {
		return nil, nil, err
	}

	var sectionItems = make([]*platformclientmodels.SectionItem, len(items))
	for idx, itemField := range items {
		sectionItems[idx] = &platformclientmodels.SectionItem{
			ID:  &itemField.ID,
			Sku: itemField.SKU,
		}
	}

	sectionTitle := itemDiff + " Section"
	localization := make(map[string]platformclientmodels.Localization)
	localization["en"] = platformclientmodels.Localization{
		Title: &sectionTitle,
	}

	startDate := time.Now()
	ok, err := sectionService.CreateSectionShort(&section.CreateSectionParams{
		Body: &platformclientmodels.SectionCreate{
			Active:       true,
			DisplayOrder: 1,
			StartDate:    strfmt.DateTime(startDate),
			EndDate:      strfmt.DateTime(startDate.Add(durationTwoDays)),
			FixedPeriodRotationConfig: &platformclientmodels.FixedPeriodRotationConfig{
				BackfillType: platformclientmodels.FixedPeriodRotationConfigBackfillTypeNONE,
				Rule:         platformclientmodels.FixedPeriodRotationConfigRuleSEQUENCE,
			},
			Items:         sectionItems,
			Localizations: localization,
			Name:          &sectionTitle,
			RotationType:  platformclientmodels.SectionCreateRotationTypeFIXEDPERIOD,
			ViewID:        viewId,
		},
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	})
	if err != nil {
		return nil, nil, err
	}

	result := &SimpleSectionInfo{
		ID:    *ok.SectionID,
		Items: items,
	}

	if doPublish {
		publishStoreChange(storeId)
	}

	return result, items, nil
}

func enableFixedRotationWithCustomBackfillForSection(c *cli.Context, storeId, sectionId string, doPublish bool) error {
	if storeId == "" {
		return fmt.Errorf("no store ID stored")
	}

	startDate := time.Now()
	_, err := sectionService.UpdateSectionShort(&section.UpdateSectionParams{
		Body: &platformclientmodels.SectionUpdate{
			Active: true,
			FixedPeriodRotationConfig: &platformclientmodels.FixedPeriodRotationConfig{
				BackfillType: platformclientmodels.FixedPeriodRotationConfigBackfillTypeCUSTOM,
				Rule:         platformclientmodels.FixedPeriodRotationConfigRuleSEQUENCE,
				Duration:     24 * 60,
				ItemCount:    3,
			},
			RotationType: platformclientmodels.SectionUpdateRotationTypeFIXEDPERIOD,
			StartDate:    strfmt.DateTime(startDate),
			EndDate:      strfmt.DateTime(startDate.Add(durationTwoDays)),
		},
		Namespace: c.String(FlagNamespace),
		SectionID: sectionId,
		StoreID:   storeId,
	})
	if err != nil {
		return err
	}

	if doPublish {
		publishStoreChange(storeId)
	}

	return nil
}

func enableCustomRotationForSection(c *cli.Context, storeId, sectionId string, doPublish bool) error {
	if storeId == "" {
		return fmt.Errorf("no store ID stored")
	}

	startDate := time.Now()
	_, err := sectionService.UpdateSectionShort(&section.UpdateSectionParams{
		Body: &platformclientmodels.SectionUpdate{
			Active:       true,
			RotationType: platformclientmodels.SectionUpdateRotationTypeCUSTOM,
			StartDate:    strfmt.DateTime(startDate),
			EndDate:      strfmt.DateTime(startDate.Add(time.Hour * 24)),
		},
		Namespace: c.String(FlagNamespace),
		SectionID: sectionId,
		StoreID:   storeId,
	})
	if err != nil {
		return err
	}

	if doPublish {
		publishStoreChange(storeId)
	}

	return nil
}

func GetSectionRotationItems(c *cli.Context, userId, viewId string) ([]*SimpleSectionInfo, error) {
	activeSections, err := sectionService.PublicListActiveSectionsShort(&section.PublicListActiveSectionsParams{
		Namespace: c.String(FlagNamespace),
		UserID:    userId,
		ViewID:    &viewId,
	})
	if err != nil {
		return nil, err
	}

	var sectionList = make([]*SimpleSectionInfo, len(activeSections))
	for idx, activeSection := range activeSections {
		var sectionInfo = SimpleSectionInfo{
			ID:    Val(activeSection.SectionID),
			Items: []*SimpleItemInfo{},
		}

		for _, currentRotationItem := range activeSection.CurrentRotationItems {
			sectionInfo.Items = append(sectionInfo.Items, &SimpleItemInfo{
				ID:    Val(currentRotationItem.ItemID),
				SKU:   currentRotationItem.Sku,
				Title: Val(currentRotationItem.Title),
			})
		}

		sectionList[idx] = &sectionInfo
	}

	return sectionList, nil
}

func DeleteStore(c *cli.Context, storeId string) error {
	inputDelete := &store.DeleteStoreParams{
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	}

	_, err := storeService.DeleteStoreShort(inputDelete)
	if err != nil {
		return err
	}

	_, err = storeService.DeletePublishedStoreShort(&store.DeletePublishedStoreParams{
		Namespace: c.String(FlagNamespace),
	})

	return fmt.Errorf("delete published store error: %w", err)
}

func GrantEntitlement(c *cli.Context, storeID string, userID string, itemID string, count int32) (string, error) {
	entitlementWrapper := platform.EntitlementService{
		Client:           factory.NewPlatformClient(&configRepo),
		ConfigRepository: &configRepo,
		TokenRepository:  &tokenRepo,
	}
	entitlementInfo, err := entitlementWrapper.GrantUserEntitlementShort(&entitlement.GrantUserEntitlementParams{
		Namespace: c.String(FlagNamespace),
		UserID:    userID,
		Body: []*platformclientmodels.EntitlementGrant{
			{
				ItemID:        Ptr(itemID),
				Quantity:      Ptr(count),
				Source:        platformclientmodels.EntitlementGrantSourcePURCHASE,
				StoreID:       storeID,
				ItemNamespace: Ptr(c.String(FlagNamespace)),
			},
		},
	})
	if err != nil {
		return "", err
	}
	if len(entitlementInfo) <= 0 {
		return "", fmt.Errorf("could not grant item to user")
	}

	return Val(entitlementInfo[0].ID), nil
}
