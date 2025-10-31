// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"fmt"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/api/dto"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"
	"gotribe-admin/pkg/util"
	"strings"

	"github.com/dengmengmian/ghelper/gconvert"
	"github.com/gin-gonic/gin"
)

type IProductController interface {
	GetProductInfo(c *gin.Context)          // 获取当前产品信息
	GetProducts(c *gin.Context)             // 获取产品列表
	CreateProduct(c *gin.Context)           // 创建产品
	UpdateProductByID(c *gin.Context)       // 更新产品
	BatchDeleteProductByIds(c *gin.Context) // 批量删除产品
}

type ProductController struct {
	ProductRepository     repository.IProductRepository
	ProductSpecRepository repository.IProductSpecRepository
	ProductSkuRepository  repository.IProductSkuRepository
}

// 构造函数
func NewProductController() IProductController {
	productRepository := repository.NewProductRepository(common.DB)
	productSpecRepository := repository.NewProductSpecRepository()
	productSku := repository.NewProductSkuRepository()
	productController := ProductController{
		ProductRepository:     productRepository,
		ProductSpecRepository: productSpecRepository,
		ProductSkuRepository:  productSku,
	}

	return productController
}

// 获取当前产品信息
// 获取当前产品信息
// GetProductInfo 获取当前产品信息
// @Summary      获取产品信息
// @Description  根据产品ID获取产品详细信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        productID path string true "产品ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /product/{productID} [get]
// @Security     BearerAuth
func (tc ProductController) GetProductInfo(c *gin.Context) {
	product, err := tc.ProductRepository.GetProductByProductID(c.Param("productID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}

	productInfoDto := dto.ToProductInfoDto(product)
	// 通过 productID 获取 sku 信息,并追加进去
	sku, err := tc.ProductSkuRepository.GetProductSkuByProductID(product.ProductID)
	if err != nil {
		response.InternalServerError(c, "获取SKU信息失败: "+err.Error())
		return
	}
	productInfoDto.SKU = dto.ToProductSkusDto(sku)
	response.Success(c, gin.H{
		"product": productInfoDto,
	}, common.Msg(c, common.MsgGetSuccess))
}

// GetProducts 获取产品列表
// @Summary      获取产品列表
// @Description  获取所有产品的列表
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        request query vo.ProductListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /product [get]
// @Security     BearerAuth
func (tc ProductController) GetProducts(c *gin.Context) {
	var req vo.ProductListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 获取
	product, total, err := tc.ProductRepository.GetProducts(&req)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgListFail)
		return
	}
	response.Success(c, gin.H{"products": dto.ToProductsDto(product), "total": total}, common.Msg(c, common.MsgListSuccess))
}

// CreateProduct 创建产品
// @Summary      创建产品
// @Description  创建一个新的产品
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        request body vo.CreateProductRequest true "创建产品请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /product [post]
// @Security     BearerAuth
func (tc ProductController) CreateProduct(c *gin.Context) {
	var req vo.CreateProductRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	// 校验参数req.sku
	if len(req.SKU) > 0 {
		for _, sku := range req.SKU {
			if sku.CostPrice <= 0 {
				response.ValidationFail(c, "成本价必填")
				return
			}
			if sku.UnitPrice <= 0 {
				response.ValidationFail(c, "市场价必填")
				return
			}
			if sku.MarketPrice <= 0 {
				response.ValidationFail(c, "市场价必填")
				return
			}
			if sku.Quantity <= 0 {
				response.ValidationFail(c, "库存必填")
				return
			}
			if sku.UnitPoint <= 0 {
				response.ValidationFail(c, "积分数必填")
				return
			}
			if gconvert.IsEmpty(sku.Title) {
				response.ValidationFail(c, "商品名必填")
				return
			}
		}
	}

	tx, err := tc.ProductRepository.BeginTx()
	if err != nil {
		response.InternalServerError(c, "开始事务失败: "+err.Error())
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	imageStr := strings.Join(req.Image, ",")
	product := model.Product{
		Title:         req.Title,
		Content:       req.Content,
		Description:   req.Description,
		Image:         imageStr,
		Video:         req.Video,
		ProductNumber: req.ProductNumber,
		CategoryID:    req.CategoryID,
		ProjectID:     req.ProjectID,
		BuyLimit:      req.BuyLimit,
		Enable:        req.Enable,
		ProductSpec:   req.ProductSpec,
		HtmlContent:   req.HtmlContent,
		Tag:           req.Tag,
	}

	productInfo, err := tc.ProductRepository.CreateProduct(tx, &product)
	if err != nil {
		tx.Rollback()
		response.HandleDatabaseError(c, err, common.MsgCreateFail)
		return
	}
	// 创建产品成功后创建 SKU
	if len(req.SKU) > 0 {
		for _, sku := range req.SKU {
			productSku := model.ProductSku{
				//SKUID:         gid.GenShortID(),
				ProductID:     productInfo.ProductID,
				CostPrice:     util.MoneyUtil.YuanToCents(sku.CostPrice),
				EnableDefault: uint(sku.EnableDefault),
				Image:         sku.Image,
				MarketPrice:   util.MoneyUtil.YuanToCents(sku.MarketPrice),
				Quantity:      uint(sku.Quantity),
				Title:         sku.Title,
				UnitPrice:     util.MoneyUtil.YuanToCents(sku.UnitPrice),
				UnitPoint:     util.MoneyUtil.YuanToCents(sku.UnitPoint),
				ProjectID:     productInfo.ProjectID,
			}
			if _, err := tc.ProductRepository.CreateProductSku(tx, &productSku); err != nil {
				tx.Rollback()
				response.HandleDatabaseError(c, err, common.MsgCreateFail)
				return
			}
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.InternalServerError(c, "提交事务失败: "+err.Error())
		return
	}
	response.Success(c, gin.H{"product": dto.ToProductInfoDto(*productInfo)}, common.Msg(c, common.MsgCreateSuccess))
}

// UpdateProductByID 更新产品
// @Summary      更新产品
// @Description  根据产品ID更新产品信息
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        productID path string true "产品ID"
// @Param        request body vo.CreateProductRequest true "更新产品请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /product/{productID} [patch]
// @Security     BearerAuth
func (tc ProductController) UpdateProductByID(c *gin.Context) {
	var req vo.CreateProductRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}
	// 校验参数req.sku
	if len(req.SKU) > 0 {
		for _, sku := range req.SKU {
			if sku.CostPrice <= 0 {
				response.ValidationFail(c, "成本价必填")
				return
			}
			if sku.UnitPrice <= 0 {
				response.ValidationFail(c, "市场价必填")
				return
			}
			if sku.MarketPrice <= 0 {
				response.ValidationFail(c, "市场价必填")
				return
			}
			if sku.Quantity <= 0 {
				response.ValidationFail(c, "库存必填")
				return
			}
			if sku.UnitPoint <= 0 {
				response.ValidationFail(c, "积分数必填")
				return
			}
			if gconvert.IsEmpty(sku.Title) {
				response.ValidationFail(c, "商品名必填")
				return
			}
		}
	}
	imageStr := strings.Join(req.Image, ",")
	// 根据path中的ProductID获取产品信息
	oldProduct, err := tc.ProductRepository.GetProductByProductID(c.Param("productID"))
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgGetFail)
		return
	}

	tx, err := tc.ProductRepository.BeginTx()
	if err != nil {
		response.InternalServerError(c, "开始事务失败: "+err.Error())
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			response.InternalServerError(c, "更新产品过程中发生错误: "+fmt.Sprintf("%v", r))
		}
	}()

	oldProduct.Title = req.Title
	oldProduct.Content = req.Content
	oldProduct.HtmlContent = req.HtmlContent
	oldProduct.Description = req.Description
	oldProduct.Image = imageStr
	oldProduct.Video = req.Video
	oldProduct.ProductNumber = req.ProductNumber
	oldProduct.CategoryID = req.CategoryID
	oldProduct.ProjectID = req.ProjectID
	oldProduct.BuyLimit = req.BuyLimit
	oldProduct.Enable = req.Enable
	oldProduct.ProductSpec = req.ProductSpec
	oldProduct.Tag = req.Tag
	// 更新产品
	err = tc.ProductRepository.UpdateProduct(tx, &oldProduct)
	if err != nil {
		tx.Rollback()
		response.HandleDatabaseError(c, err, common.MsgUpdateFail)
		return
	}

	// 更新商品SKU
	if len(req.SKU) > 0 {
		for _, sku := range req.SKU {
			productSku, err := tc.ProductRepository.GetProductSkuByProductSkuID(tx, sku.SKUID)
			if err != nil {
				tx.Rollback()
				response.InternalServerError(c, "获取需要更新的产品SKU信息失败: "+err.Error())
				return
			}
			productSku.CostPrice = util.MoneyUtil.YuanToCents(sku.CostPrice)
			productSku.EnableDefault = sku.EnableDefault
			productSku.Image = sku.Image
			productSku.MarketPrice = util.MoneyUtil.YuanToCents(sku.MarketPrice)
			productSku.Quantity = sku.Quantity
			productSku.Title = sku.Title
			productSku.UnitPrice = util.MoneyUtil.YuanToCents(sku.UnitPrice)
			productSku.UnitPoint = util.MoneyUtil.YuanToCents(sku.UnitPoint)
			productSku.ProductID = oldProduct.ProductID

			err = tc.ProductRepository.UpdateProductSku(tx, productSku)
			if err != nil {
				tx.Rollback()
				response.HandleDatabaseError(c, err, common.MsgUpdateFail)
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response.InternalServerError(c, "提交事务失败: "+err.Error())
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgUpdateSuccess))
}

// BatchDeleteProductByIds 批量删除产品
// @Summary      批量删除产品
// @Description  根据产品ID列表批量删除产品
// @Tags         产品管理
// @Accept       json
// @Produce      json
// @Param        request body vo.DeleteProductsRequest true "删除产品请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /product [delete]
// @Security     BearerAuth
func (tc ProductController) BatchDeleteProductByIds(c *gin.Context) {
	var req vo.DeleteProductsRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		response.HandleBindError(c, err)
		return
	}
	// 参数校验
	if err := common.Validate.Struct(&req); err != nil {
		response.HandleValidationError(c, err)
		return
	}

	// 前端传来的产品ID
	reqProductIds := strings.Split(req.ProductIds, ",")
	err := tc.ProductRepository.BatchDeleteProductByIds(reqProductIds)
	if err != nil {
		response.HandleDatabaseError(c, err, common.MsgDeleteFail)
		return
	}

	response.Success(c, nil, common.Msg(c, common.MsgDeleteSuccess))

}
