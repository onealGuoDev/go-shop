package main

import (
	"context"
	"fmt"
	"go-shop/srvs/good_srv/proto"
	"google.golang.org/grpc"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetGoodsList() {
	rsp, err := brandClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin:    90,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestBatchGetGoods() {
	rsp, err := brandClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, good := range rsp.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestGetGoodsDetail() {
	rsp, err := brandClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Name)
	fmt.Println(rsp.DescImages)
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()
	//TestCreateUser()
	//TestGetGoodsList()
	//TestBatchGetGoods()
	TestGetGoodsDetail()

	conn.Close()
}
