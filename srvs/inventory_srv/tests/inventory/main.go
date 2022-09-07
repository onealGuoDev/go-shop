package main

import (
	"context"
	"fmt"
	"go-shop/srvs/inventory_srv/proto"
	"google.golang.org/grpc"
	"sync"
)

var invClient proto.InventoryClient
var conn *grpc.ClientConn

func TestSetInv(goodsId, Num int32) {
	_, err := invClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     Num,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("設置庫存成功")
}

func TestInvDetail(goodsId int32) {
	rsp, err := invClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Num)
}

func TestSell(wg *sync.WaitGroup) {

	defer wg.Done()
	_, err := invClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
			//{GoodsId: 422, Num: 30},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("庫存扣減成功")
}

func TestReback() {
	_, err := invClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 30},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("歸還成功")
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	invClient = proto.NewInventoryClient(conn)
}

func main() {
	Init()
	//var i int32
	//for i = 421; i<=840; i++ {
	//	TestSetInv(i, 100)
	//}
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go TestSell(&wg)
	}

	wg.Wait()

	//TestInvDetail(421)
	//TestSell()
	//TestReback()
	conn.Close()
}
