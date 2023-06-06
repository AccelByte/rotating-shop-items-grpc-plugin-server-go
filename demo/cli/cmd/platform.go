// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/catalog_changes"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/category"
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
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	servicePluginConfigService = &platform.ServicePluginConfigService{
		Client:          factory.NewPlatformClient(&configRepo),
		TokenRepository: &tokenRepo,
	}

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

	abViewName   = "Item Rotation Default View"
	displayOrder = int32(1)

	currencyCode      = "USD"
	currencyNamespace = "accelbyte"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type simpleItemInfo struct {
	id    string
	sku   string
	title string
}

type simpleSectionInfo struct {
	id    string
	items []*simpleItemInfo
}

func SetPlatformServiceGrpcTarget(c *cli.Context) error {
	ok, errOK := servicePluginConfigService.UpdateServicePluginConfigShort(&service_plugin_config.UpdateServicePluginConfigParams{
		Body:      &platformclientmodels.ServicePluginConfigUpdate{GrpcServerAddress: FlagGrpcUrl},
		Namespace: c.String(FlagNamespace),
	})
	if errOK != nil {
		logrus.Errorf("could not set the grpc url. %s", errOK)

		return errOK
	}

	logrus.Printf("set grpc url %s in namespace %s", ok.GrpcServerAddress, ok.Namespace)

	return nil
}

func CreateStore(c *cli.Context) (string, error) {
	ok, errOK := storeService.CreateStoreShort(&store.CreateStoreParams{
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
	if errOK != nil {
		logrus.Errorf("could not create store. %s", errOK)

		return "", errOK
	}

	return *ok.StoreID, nil
}

func CreateCategory(c *cli.Context, storeId string) error {
	categoryPath := getCategoryPath(c)
	localz := make(map[string]string)
	localz["en"] = categoryPath

	ok, errOK := categoryService.CreateCategoryShort(&category.CreateCategoryParams{
		Body: &platformclientmodels.CategoryCreate{
			CategoryPath:             &categoryPath,
			LocalizationDisplayNames: localz,
		},
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	})
	if errOK != nil {
		logrus.Errorf("could not create category. %s", errOK)

		return errOK
	}

	logrus.Printf("category created with path. %s", *ok.CategoryPath)

	return nil
}

func CreateStoreView(c *cli.Context, storeId string) (string, error) {
	localization := make(map[string]platformclientmodels.Localization)
	localization["en"] = platformclientmodels.Localization{
		Description:     "",
		LongDescription: "",
		Title:           &abViewName,
	}
	ok, errOK := viewService.CreateViewShort(&view.CreateViewParams{
		Body: &platformclientmodels.ViewCreate{
			DisplayOrder:  &displayOrder,
			Localizations: localization,
			Name:          &abViewName,
		},
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	})
	if errOK != nil {
		logrus.Errorf("could not create view. %s", errOK.Error())

		return "", errOK
	}

	return *ok.ViewID, nil
}

func publishStoreChange(storeId string) string {
	inputCreate := &catalog_changes.PublishAllParams{
		StoreID:   storeId,
		Namespace: integration.NamespaceTest,
	}

	created, errCreate := catalogChangesService.PublishAllShort(inputCreate)
	if errCreate != nil {
		logrus.Error(errCreate.Error())

		return ""
	}
	storeID := *created.StoreID

	return storeID
}

func CreateItems(c *cli.Context, storeId, itemDiff string, itemCount int, doPublish bool) ([]*simpleItemInfo, error) {
	iType := platformclientmodels.ItemCreateItemTypeSEASON
	eType := platformclientmodels.ItemCreateEntitlementTypeDURABLE
	sType := platformclientmodels.ItemCreateSeasonTypeTIER
	status := platformclientmodels.ItemCreateStatusACTIVE
	cType := platformclientmodels.RegionDataItemCurrencyTypeREAL

	categoryPath := getCategoryPath(c)

	iLocalization := make(map[string]platformclientmodels.Localization)
	iRegionData := make(map[string][]platformclientmodels.RegionDataItemDTO)

	var nItems []*simpleItemInfo

	for i := 0; i < itemCount; i++ {
		price := int32((i + 1) * 2)
		var nItemInfo simpleItemInfo
		nItemInfo.title = fmt.Sprint("Item " + itemDiff + " Titled " + strconv.FormatInt(int64(i+1), 10))
		nItemInfo.sku = fmt.Sprint("SKU_" + itemDiff + "_" + strconv.FormatInt(int64(i+1), 10))

		iLocalization["en"] = platformclientmodels.Localization{
			Title: &nItemInfo.title,
		}

		regionData := platformclientmodels.RegionDataItemDTO{
			CurrencyCode:      &currencyCode,
			CurrencyNamespace: &currencyNamespace,
			CurrencyType:      &cType,
			Price:             &price,
		}
		iRegionData["US"] = append(iRegionData["key"], regionData)

		ok, errOK := itemService.CreateItemShort(&item.CreateItemParams{
			Body: &platformclientmodels.ItemCreate{
				Features:        []string{"test"},
				Tags:            []string{"tags"},
				Name:            &nItemInfo.title,
				ItemType:        &iType,
				CategoryPath:    &categoryPath,
				EntitlementType: &eType,
				SeasonType:      sType,
				Status:          &status,
				Listable:        true,
				Purchasable:     true,
				Sku:             nItemInfo.sku,
				Localizations:   iLocalization,
				RegionData:      iRegionData,
			},
			Namespace: c.String(FlagNamespace),
			StoreID:   storeId,
		})
		if errOK != nil {
			logrus.Errorf("could not create item. %s", errOK.Error())

			return nil, errOK
		}

		// add the items
		nItemInfo.id = *ok.ItemID
		nItems = append(nItems, &nItemInfo)
	}

	if doPublish {
		publishStoreChange(storeId)
		logrus.Infof("publish storeId %s.", storeId)

		return nItems, nil
	}

	return nItems, nil
}

func CreateSectionsWithItems(c *cli.Context, storeId, viewId string, itemCount int, doPublish bool) (*simpleSectionInfo, error) {
	itemDiff := RandStringBytes(6)
	items, errOK := CreateItems(c, storeId, itemDiff, itemCount, doPublish)
	if errOK != nil {
		logrus.Errorf("could not create items section. %s", errOK.Error())

		return nil, errOK
	}

	var sectionItems []*platformclientmodels.SectionItem
	for _, itemField := range items {
		sectionItems = append(sectionItems, &platformclientmodels.SectionItem{
			ID:  &itemField.id,
			Sku: itemField.sku,
		})
	}
	sectionTitle := itemDiff + " Section"
	localz := make(map[string]platformclientmodels.Localization)
	localz["en"] = platformclientmodels.Localization{
		Title: &sectionTitle,
	}

	ok, err := sectionService.CreateSectionShort(&section.CreateSectionParams{
		Body: &platformclientmodels.SectionCreate{
			Active:       true,
			DisplayOrder: 1,
			EndDate:      strfmt.DateTime(time.Date(2025, 8, 9, 0, 0, 0, 0, time.UTC)),
			FixedPeriodRotationConfig: &platformclientmodels.FixedPeriodRotationConfig{
				BackfillType: platformclientmodels.FixedPeriodRotationConfigBackfillTypeNONE,
				Rule:         platformclientmodels.FixedPeriodRotationConfigRuleSEQUENCE,
			},
			Items:         sectionItems,
			Localizations: localz,
			Name:          &sectionTitle,
			RotationType:  platformclientmodels.SectionCreateRotationTypeFIXEDPERIOD,
			StartDate:     strfmt.DateTime(time.Now()),
			ViewID:        viewId,
		},
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	})
	if err != nil {
		logrus.Errorf("could not create section. %s", errOK.Error())

		return nil, errOK
	}

	result := &simpleSectionInfo{
		id:    *ok.SectionID,
		items: items,
	}

	if doPublish {
		publishStoreChange(storeId)
		logrus.Infof("publish storeId %s.", storeId)

		return result, nil
	}

	return result, nil
}

func enableFixedRotationWithCustomBackfillForSection(c *cli.Context, storeId, sectionId string, doPublish bool) error {
	if storeId == "" {
		return fmt.Errorf("no store id stored")
	}

	sectionTitle := "Update Section"
	localz := make(map[string]platformclientmodels.Localization)
	localz["en"] = platformclientmodels.Localization{
		Title: &sectionTitle,
	}
	_, err := sectionService.UpdateSectionShort(&section.UpdateSectionParams{
		Body: &platformclientmodels.SectionUpdate{
			FixedPeriodRotationConfig: &platformclientmodels.FixedPeriodRotationConfig{
				BackfillType: platformclientmodels.FixedPeriodRotationConfigBackfillTypeCUSTOM,
				Rule:         platformclientmodels.FixedPeriodRotationConfigRuleSEQUENCE,
			},
			RotationType:  platformclientmodels.SectionUpdateRotationTypeFIXEDPERIOD,
			Localizations: localz,
			Name:          &sectionTitle,
		},
		Namespace: c.String(FlagNamespace),
		SectionID: sectionId,
		StoreID:   storeId,
	})
	if err != nil {
		logrus.Errorf("could not update section for custom backfill. %s", err.Error())

		return err
	}

	if doPublish {
		publishStoreChange(storeId)
		logrus.Infof("publish storeId %s.", storeId)

		return nil
	}

	return nil
}

func enableCustomRotationForSection(c *cli.Context, storeId, sectionId string, doPublish bool) error {
	if storeId == "" {
		return fmt.Errorf("no store id stored")
	}

	_, err := sectionService.UpdateSectionShort(&section.UpdateSectionParams{
		Body: &platformclientmodels.SectionUpdate{
			RotationType: platformclientmodels.SectionUpdateRotationTypeCUSTOM,
		},
		Namespace: c.String(FlagNamespace),
		SectionID: sectionId,
		StoreID:   storeId,
	})
	if err != nil {
		logrus.Errorf("could not update section. %s", err.Error())

		return err
	}

	if doPublish {
		publishStoreChange(storeId)
		logrus.Infof("publish storeId %s.", storeId)

		return nil
	}

	return nil
}

func GetSectionRotationItems(c *cli.Context, userId, viewId string) ([]*simpleSectionInfo, error) {
	activeSections, errOK := sectionService.PublicListActiveSectionsShort(&section.PublicListActiveSectionsParams{
		Namespace: c.String(FlagNamespace),
		UserID:    userId,
		ViewID:    &viewId,
	})
	if errOK != nil {
		logrus.Errorf("could not get active sessions. %s", errOK)

		return nil, errOK
	}

	var iSection []*simpleSectionInfo
	var rItems []*platformclientmodels.ItemInfo
	for _, activeSection := range activeSections {
		for _, currentRotationItem := range activeSection.CurrentRotationItems {
			rItems = append(rItems, currentRotationItem)
			var sectionInfo *simpleSectionInfo
			sectionInfo.id = *activeSection.SectionID

			if rItems != nil && len(rItems) != 0 {
				var items []*simpleItemInfo
				for _, it := range rItems {
					var itemInfo *simpleItemInfo
					itemInfo.id = *it.ItemID
					itemInfo.sku = it.Sku
					itemInfo.title = *it.Title

					items = append(items, itemInfo)
				}
				sectionInfo.items = items
			} else {
				sectionInfo.items = []*simpleItemInfo{}
			}

			iSection = append(iSection, sectionInfo)
		}
	}

	return iSection, nil
}

func DeleteStore(c *cli.Context, storeId string) (*platformclientmodels.StoreInfo, error) {
	inputDelete := &store.DeleteStoreParams{
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	}

	ok, errOK := storeService.DeleteStoreShort(inputDelete)
	if errOK != nil {
		logrus.Errorf("could not delete store %s", errOK)

		return nil, errOK
	}

	return ok, nil
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		rand.Seed(time.Now().UnixNano())
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func getCategoryPath(c *cli.Context) string {
	if c.String(FlagCategoryPath) == "" {
		return "/customitemrotationtest"
	}

	return c.String(FlagCategoryPath)
}
