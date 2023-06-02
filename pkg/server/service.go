// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	pb "rotating-shop-items-grpc-plugin-server-go/pkg/pb"
)

type SectionServiceServer struct {
	pb.UnimplementedSectionServer
}

func NewSectionServiceServer() (*SectionServiceServer, error) {
	return &SectionServiceServer{}, nil
}

func (s SectionServiceServer) GetRotationItems(ctx context.Context, request *pb.GetRotationItemsRequest) (*pb.GetRotationItemsResponse, error) {
	logrus.Info("GetRotationItems")

	var items []*pb.SectionItemObject
	inputCount := float64(len(items))

	upperLimit := float64(24)
	currentPoint := float64(time.Now().Hour())
	selectedIndex := math.Floor((inputCount / upperLimit) * currentPoint)

	var selectedItem *pb.SectionItemObject
	for i, item := range request.SectionObject.Items {
		if i == int(selectedIndex) {
			selectedItem = item

			items = append(items, &pb.SectionItemObject{
				ItemId:  selectedItem.ItemId,
				ItemSku: selectedItem.ItemSku,
			})
		}
	}

	return &pb.GetRotationItemsResponse{
		Items:     items,
		ExpiredAt: 0,
	}, nil
}

func (s SectionServiceServer) Backfill(ctx context.Context, request *pb.BackfillRequest) (*pb.BackfillResponse, error) {
	logrus.Info("Backfill")

	var newItems []*pb.BackfilledItemObject

	for i, item := range request.GetItems() {
		if item.Owned {
			//if and item is owned by user, then replace it with new item id.
			//item id will be generated randomly for example purpose.

			newItem := &pb.BackfilledItemObject{
				ItemId:  strings.ReplaceAll(uuid.NewString(), "-", ""),
				ItemSku: strconv.FormatInt(int64(i), rand.Int()),
				Index:   int32(i),
			}

			newItems = append(newItems, newItem)
		}
	}

	return &pb.BackfillResponse{BackfilledItems: newItems}, nil
}
