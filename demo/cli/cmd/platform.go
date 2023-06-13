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
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/store"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclient/view"
	"github.com/AccelByte/accelbyte-go-sdk/platform-sdk/pkg/platformclientmodels"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/factory"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/platform"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/tests/integration"
	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"accelbyte.net/rotating-shop-items-cli/pkg/client/platformservice"
	"accelbyte.net/rotating-shop-items-cli/pkg/client/platformservice/openapi2/models"
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

	abViewName   = "Go Item Rotation Default View Demo/CLI"
	displayOrder = int32(1)

	currencyCode      = "USD"
	currencyNamespace = "accelbyte"
)

var platformClientSvc *platformservice.Client

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
	if c.String(FlagGrpcUrl) != "" {
		fmt.Printf("(Custom Host: %s) ", c.String(FlagGrpcUrl))

		return platformClientSvc.UpdateSectionPluginConfig(c.String(FlagNamespace), &models.SectionPluginConfigUpdate{
			ExtendType: Ptr(models.SectionPluginConfigUpdateExtendTypeCUSTOM),
			CustomConfig: &models.BaseCustomConfig{
				ConnectionType:    Ptr(models.BaseCustomConfigConnectionTypeINSECURE),
				GrpcServerAddress: Ptr(c.String(FlagGrpcUrl)),
			},
		})
	}

	if c.String(FlagExtendAppName) != "" {
		fmt.Printf("(Extend App: %s) ", c.String(FlagExtendAppName))

		return platformClientSvc.UpdateSectionPluginConfig(c.String(FlagNamespace), &models.SectionPluginConfigUpdate{
			ExtendType: Ptr(models.SectionPluginConfigUpdateExtendTypeAPP),
			AppConfig: &models.AppConfig{
				AppName: Ptr(c.String(FlagExtendAppName)),
			},
		})
	}

	return nil
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
		if Val(s.Published) == false {
			_, _ = storeService.DeleteStoreShort(&store.DeleteStoreParams{
				Namespace: c.String(FlagNamespace),
				StoreID:   Val(s.StoreID),
			})
		}
	}

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

func CreateItems(c *cli.Context, storeId, itemDiff string, itemCount int, doPublish bool) ([]*SimpleItemInfo, error) {
	iType := platformclientmodels.ItemCreateItemTypeSEASON
	eType := platformclientmodels.ItemCreateEntitlementTypeDURABLE
	sType := platformclientmodels.ItemCreateSeasonTypeTIER
	status := platformclientmodels.ItemCreateStatusACTIVE
	cType := platformclientmodels.RegionDataItemCurrencyTypeREAL

	categoryPath := getCategoryPath(c)

	iLocalization := make(map[string]platformclientmodels.Localization)
	iRegionData := make(map[string][]platformclientmodels.RegionDataItemDTO)

	var nItems []*SimpleItemInfo

	for i := 0; i < itemCount; i++ {
		price := int32((i + 1) * 2)
		var nItemInfo SimpleItemInfo
		nItemInfo.Title = fmt.Sprint("Item " + itemDiff + " Titled " + strconv.FormatInt(int64(i+1), 10))
		nItemInfo.SKU = fmt.Sprint("SKU_" + itemDiff + "_" + strconv.FormatInt(int64(i+1), 10))

		iLocalization["en"] = platformclientmodels.Localization{
			Title: &nItemInfo.Title,
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
				Name:            &nItemInfo.Title,
				ItemType:        &iType,
				CategoryPath:    &categoryPath,
				EntitlementType: &eType,
				SeasonType:      sType,
				Status:          &status,
				Listable:        true,
				Purchasable:     true,
				Sku:             nItemInfo.SKU,
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

		// add the Items
		nItemInfo.ID = *ok.ItemID
		nItems = append(nItems, &nItemInfo)
	}

	if doPublish {
		publishStoreChange(storeId)
		logrus.Infof("publish storeId %s.", storeId)

		return nItems, nil
	}

	return nItems, nil
}

func CreateSectionsWithItems(c *cli.Context, storeId, viewId string, itemCount int, doPublish bool) (*SimpleSectionInfo, error) {
	itemDiff := RandStringBytes(6)
	items, errOK := CreateItems(c, storeId, itemDiff, itemCount, doPublish)
	if errOK != nil {
		logrus.Errorf("could not create Items section. %s", errOK.Error())

		return nil, errOK
	}

	var sectionItems []*platformclientmodels.SectionItem
	for _, itemField := range items {
		sectionItems = append(sectionItems, &platformclientmodels.SectionItem{
			ID:  &itemField.ID,
			Sku: itemField.SKU,
		})
	}
	sectionTitle := itemDiff + " Section"
	localz := make(map[string]platformclientmodels.Localization)
	localz["en"] = platformclientmodels.Localization{
		Title: &sectionTitle,
	}

	startDate := time.Now()
	ok, err := sectionService.CreateSectionShort(&section.CreateSectionParams{
		Body: &platformclientmodels.SectionCreate{
			Active:       true,
			DisplayOrder: 1,
			StartDate:    strfmt.DateTime(startDate),
			EndDate:      strfmt.DateTime(startDate.Add(time.Hour * 24)),
			FixedPeriodRotationConfig: &platformclientmodels.FixedPeriodRotationConfig{
				BackfillType: platformclientmodels.FixedPeriodRotationConfigBackfillTypeNONE,
				Rule:         platformclientmodels.FixedPeriodRotationConfigRuleSEQUENCE,
			},
			Items:         sectionItems,
			Localizations: localz,
			Name:          &sectionTitle,
			RotationType:  platformclientmodels.SectionCreateRotationTypeFIXEDPERIOD,
			ViewID:        viewId,
		},
		Namespace: c.String(FlagNamespace),
		StoreID:   storeId,
	})
	if err != nil {
		logrus.Errorf("could not create section. %s", errOK.Error())

		return nil, errOK
	}

	result := &SimpleSectionInfo{
		ID:    *ok.SectionID,
		Items: items,
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
		return fmt.Errorf("no store ID stored")
	}

	startDate := time.Now()
	_, err := sectionService.UpdateSectionShort(&section.UpdateSectionParams{
		Body: &platformclientmodels.SectionUpdate{
			Active: true,
			FixedPeriodRotationConfig: &platformclientmodels.FixedPeriodRotationConfig{
				BackfillType: platformclientmodels.FixedPeriodRotationConfigBackfillTypeCUSTOM,
				Rule:         platformclientmodels.FixedPeriodRotationConfigRuleSEQUENCE,
			},
			RotationType: platformclientmodels.SectionUpdateRotationTypeFIXEDPERIOD,
			StartDate:    strfmt.DateTime(startDate),
			EndDate:      strfmt.DateTime(startDate.Add(time.Hour * 24)),
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

func GetSectionRotationItems(c *cli.Context, userId, viewId string) ([]*SimpleSectionInfo, error) {
	activeSections, errOK := sectionService.PublicListActiveSectionsShort(&section.PublicListActiveSectionsParams{
		Namespace: c.String(FlagNamespace),
		UserID:    userId,
		ViewID:    &viewId,
	})
	if errOK != nil {
		logrus.Errorf("could not get active sessions. %s", errOK)

		return nil, errOK
	}

	var iSection []*SimpleSectionInfo
	var rItems []*platformclientmodels.ItemInfo
	for _, activeSection := range activeSections {
		for _, currentRotationItem := range activeSection.CurrentRotationItems {
			rItems = append(rItems, currentRotationItem)
			var sectionInfo SimpleSectionInfo
			sectionInfo.ID = Val(activeSection.SectionID)

			if rItems != nil && len(rItems) != 0 {
				var items []*SimpleItemInfo
				for _, it := range rItems {
					var itemInfo SimpleItemInfo
					itemInfo.ID = *it.ItemID
					itemInfo.SKU = it.Sku
					itemInfo.Title = *it.Title

					items = append(items, &itemInfo)
				}
				sectionInfo.Items = items
			} else {
				sectionInfo.Items = []*SimpleItemInfo{}
			}

			iSection = append(iSection, &sectionInfo)
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
