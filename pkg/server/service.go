// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	pb "rotating-shop-items-grpc-plugin-server-go/pkg/pb"
)

const upperLimit = 24

type SectionServiceServer struct {
	pb.UnimplementedSectionServer
}

func NewSectionServiceServer() (*SectionServiceServer, error) {
	return &SectionServiceServer{}, nil
}

func (s SectionServiceServer) GetRotationItems(ctx context.Context, request *pb.GetRotationItemsRequest) (*pb.GetRotationItemsResponse, error) {
	logJson("GetRotationItems Request: ", request)

	inputCount := len(request.GetSectionObject().GetItems())
	currentPoint := time.Now().Hour()
	selectedIndex := int(math.Floor((float64(inputCount) / float64(upperLimit)) * float64(currentPoint)))
	selectedItem := request.GetSectionObject().GetItems()[selectedIndex]

	responseItems := []*pb.SectionItemObject{
		{
			ItemId:  selectedItem.ItemId,
			ItemSku: selectedItem.ItemSku,
		},
	}

	resp := pb.GetRotationItemsResponse{
		Items:     responseItems,
		ExpiredAt: 0,
	}
	logJson("GetRotationItems Response: ", &resp)

	return &resp, nil
}

func (s SectionServiceServer) Backfill(ctx context.Context, request *pb.BackfillRequest) (*pb.BackfillResponse, error) {
	logJson("Backfill Request: ", request)

	var newItems []*pb.BackfilledItemObject

	for _, item := range request.GetItems() {
		if item.Owned {
			// if an item is owned by user, then replace it with new item id.
			// item id will be generated randomly for example purpose.
			newItem := &pb.BackfilledItemObject{
				ItemId: strings.ReplaceAll(uuid.NewString(), "-", ""),
				Index:  item.Index,
			}

			newItems = append(newItems, newItem)
		}
	}

	resp := &pb.BackfillResponse{BackfilledItems: newItems}
	logJson("Backfill Response: ", resp)

	return resp, nil
}

func logJson(msg string, data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	logrus.Infof("%s%s", msg, string(jsonData))
}
