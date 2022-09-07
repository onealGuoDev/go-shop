package shop_cart

import (
	"context"
	"go-shop/api/order_api/api"
	"go-shop/api/order_api/forms"
	"go-shop/api/order_api/proto"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-shop/api/order_api/global"
)

func List(ctx *gin.Context) {
	//获取购物车商品
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint)),
	})
	if err != nil {
		zap.S().Errorw("查詢購物車列表失敗")
		panic(err)
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Errorw("批量查詢購物車列表失敗")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if good.Id == item.GoodsId {
				tmpMap := map[string]interface{}{}
				tmpMap["id"] = item.Id
				tmpMap["goods_id"] = item.GoodsId
				tmpMap["good_name"] = good.Name
				tmpMap["good_image"] = good.GoodsFrontImage
				tmpMap["good_price"] = good.ShopPrice
				tmpMap["nums"] = item.Nums
				tmpMap["checked"] = item.Checked

				goodsList = append(goodsList, tmpMap)
			}
		}
	}
	reMap["data"] = goodsList
	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {

	itemForm := forms.ShopCartItemForm{}

	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("查詢商品失敗")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("查詢商品庫存失敗")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums": "庫存不足",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})

	if err != nil {
		zap.S().Errorw("新增到購物車失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Update(ctx *gin.Context) {
	// o/v1/421
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出錯",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	request := proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}

	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("更新購物車失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出錯",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除購物車記錄失敗")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
