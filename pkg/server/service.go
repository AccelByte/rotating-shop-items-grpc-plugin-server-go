// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
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
	logrus.Info("GetRotationItems Request")

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

	return &pb.GetRotationItemsResponse{
		Items:     responseItems,
		ExpiredAt: 0,
	}, nil
}

func (s SectionServiceServer) Backfill(ctx context.Context, request *pb.BackfillRequest) (*pb.BackfillResponse, error) {
	logrus.Info("Backfill Request")

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

	return &pb.BackfillResponse{BackfilledItems: newItems}, nil
}
