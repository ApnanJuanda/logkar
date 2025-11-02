package product

import (
	generalRepository "bsnack/domain/api/general/repository"
	"bsnack/domain/api/product/model"
	"bsnack/domain/api/product/repository"
	"bsnack/lib/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type ProductServiceInterface interface {
	InsertSize(req model.InsertSizeRequest) (responseCode int, err error)
	GetAllSize() (resp []model.Size, responseCode int, err error)
	InsertFlavor(req model.InsertFlavorRequest) (responseCode int, err error)
	GetAllFlavor() (resp []model.Flavor, responseCode int, err error)
	InsertProductType(req model.InsertProductTypeRequest) (responseCode int, err error)
	GetProductType() (resp []model.ProductType, responseCode int, err error)
	InsertProduct(req model.InsertProductRequest) (responseCode int, err error)
	InsertProductDetail(req model.InsertProductDetailRequest) (responseCode int, err error)
	GetListProduct(req model.GetProductRequest) (resp []model.GetProductResponse, count int64, responseCode int, err error)
}

type productService struct {
	Repository  repository.ProductRepositoryInterface
	GeneralRepo generalRepository.GeneralRepositoryInterface
}

func NewProductService(repository repository.ProductRepositoryInterface, generalRepo generalRepository.GeneralRepositoryInterface) ProductServiceInterface {
	return &productService{
		Repository:  repository,
		GeneralRepo: generalRepo,
	}
}

func (s *productService) InsertSize(req model.InsertSizeRequest) (responseCode int, err error) {
	var list_data = []model.Size{}
	for _, value := range req.ListName {
		seqNumber, err := s.GeneralRepo.GetDataSequence(nil, "S")
		if err != nil {
			log.Printf("ERROR GetDataSeq size: %v", err)
			return http.StatusInternalServerError, err
		}
		sizeId := fmt.Sprintf("%s%04d", "S", seqNumber)
		data := model.Size{
			ID:        sizeId,
			Name:      value,
			CreatedAt: utils.GetLocaltime(),
			CreatedBy: req.SellerEmail,
		}
		list_data = append(list_data, data)
	}
	err = s.Repository.InsertSize(list_data)
	if err != nil {
		log.Printf("ERROR InsertSize: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) GetAllSize() (resp []model.Size, responseCode int, err error) {
	resp, err = s.Repository.GetAllSize()
	if err != nil {
		log.Printf("ERROR GetAllSize: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) InsertFlavor(req model.InsertFlavorRequest) (responseCode int, err error) {
	var list_data = []model.Flavor{}
	for _, value := range req.ListName {
		seqNumber, err := s.GeneralRepo.GetDataSequence(nil, "F")
		if err != nil {
			log.Printf("ERROR GetDataSeq flavor: %v", err)
			return http.StatusInternalServerError, err
		}
		sizeId := fmt.Sprintf("%s%04d", "F", seqNumber)
		data := model.Flavor{
			ID:        sizeId,
			Name:      value,
			CreatedAt: utils.GetLocaltime(),
			CreatedBy: req.SellerEmail,
		}
		list_data = append(list_data, data)
	}
	err = s.Repository.InsertFlavor(list_data)
	if err != nil {
		log.Printf("ERROR InsertFlavor: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) InsertProductType(req model.InsertProductTypeRequest) (responseCode int, err error) {
	var list_data = []model.ProductType{}
	for _, value := range req.ListName {
		seqNumber, err := s.GeneralRepo.GetDataSequence(nil, "T")
		if err != nil {
			log.Printf("ERROR GetDataSeq product type: %v", err)
			return http.StatusInternalServerError, err
		}
		productTypeId := fmt.Sprintf("%s%04d", "T", seqNumber)
		data := model.ProductType{
			ID:        productTypeId,
			Name:      value,
			CreatedAt: utils.GetLocaltime(),
			CreatedBy: req.SellerEmail,
		}
		list_data = append(list_data, data)
	}
	err = s.Repository.InsertProductType(list_data)
	if err != nil {
		log.Printf("ERROR InsertProductType: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) GetProductType() (resp []model.ProductType, responseCode int, err error) {
	resp, err = s.Repository.GetAllProductType()
	if err != nil {
		log.Printf("ERROR GetAllProductType: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) GetAllFlavor() (resp []model.Flavor, responseCode int, err error) {
	resp, err = s.Repository.GetAllFlavor()
	if err != nil {
		log.Printf("ERROR GetAllFlavor: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) InsertProduct(req model.InsertProductRequest) (responseCode int, err error) {
	tx := s.GeneralRepo.BeginTrans()
	defer func() {
		if err != nil {
			s.GeneralRepo.RollbackTrans(tx)
			return
		}
		s.GeneralRepo.CommitTrans(tx)
	}()
	localTime := utils.GetLocaltime()
	var listNewProduct = []model.Product{}
	for _, value := range req.ListProduct {
		productType, err := s.Repository.GetProductTypeByID(tx, value.TypeId)
		if err != nil {
			log.Printf("ERROR GetProductTypeByID %v : %v", value.TypeId, err)
			return http.StatusBadRequest, err
		}
		if productType.ID == "" {
			log.Printf("ERROR GetProductTypeByID %v is not found", value.TypeId)
			return http.StatusBadRequest, errors.New("product type is not found")
		}

		seqNumber, err := s.GeneralRepo.GetDataSequence(tx, "P")
		if err != nil {
			log.Printf("ERROR GetDataSeq flavor: %v", err)
			return http.StatusInternalServerError, err
		}
		productID := fmt.Sprintf("%s%04d", "P", seqNumber)
		newProduct := model.Product{
			ID:        productID,
			Name:      value.Name,
			SellerID:  req.SellerID,
			TypeID:    value.TypeId,
			CreatedAt: localTime,
			CreatedBy: req.SellerEmail,
		}
		listNewProduct = append(listNewProduct, newProduct)
	}

	err = s.Repository.InsertProduct(tx, listNewProduct)
	if err != nil {
		log.Printf("ERROR InsertProduct: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) InsertProductDetail(req model.InsertProductDetailRequest) (responseCode int, err error) {
	tx := s.GeneralRepo.BeginTrans()
	defer func() {
		if err != nil {
			s.GeneralRepo.RollbackTrans(tx)
			return
		}
		s.GeneralRepo.CommitTrans(tx)
	}()
	localTime := utils.GetLocaltime()
	var listNewProductDetail = []model.ProductDetail{}
	for _, productInfo := range req.ListProductInfo {
		size, err := s.Repository.GetSizeByID(tx, productInfo.SizeID)
		if err != nil {
			log.Printf("ERROR GetSizeByID %v : %v", productInfo.SizeID, err)
			return http.StatusBadRequest, err
		}
		if size.ID == "" {
			log.Printf("ERROR GetSizeByID %v is not found", productInfo.SizeID)
			return http.StatusBadRequest, errors.New("size is not found")
		}

		flavor, err := s.Repository.GetFlavorByID(tx, productInfo.FlavorID)
		if err != nil {
			log.Printf("ERROR GetFlavorByID %v : %v", productInfo.FlavorID, err)
			return http.StatusBadRequest, err
		}
		if flavor.ID == "" {
			log.Printf("ERROR GetFlavorByID %v is not found", productInfo.FlavorID)
			return http.StatusBadRequest, errors.New("flavor is not found")
		}

		product, err := s.Repository.GetProductByID(tx, productInfo.ProductID)
		if err != nil {
			log.Printf("ERROR GetProductByID %v : %v", productInfo.ProductID, err)
			return http.StatusBadRequest, err
		}
		if product.ID == "" {
			log.Printf("ERROR GetProductByID %v is not found", productInfo.ProductID)
			return http.StatusBadRequest, errors.New("product is not found")
		}

		newProductDetail := model.ProductDetail{
			ProductID: product.ID,
			SizeID:    productInfo.SizeID,
			FlavorID:  productInfo.FlavorID,
			Price:     productInfo.Price,
			Stock:     productInfo.Stock,
			CreatedAt: localTime,
			CreatedBy: req.SellerEmail,
		}
		listNewProductDetail = append(listNewProductDetail, newProductDetail)
	}
	// save product detail
	err = s.Repository.InsertProductDetail(tx, listNewProductDetail)
	if err != nil {
		log.Printf("ERROR InsertProductDetail: %v", err)
		return http.StatusInternalServerError, err
	}
	responseCode = http.StatusOK
	return
}

func (s *productService) GetListProduct(req model.GetProductRequest) (resp []model.GetProductResponse, count int64, responseCode int, err error) {
	resp, count, err = s.Repository.GetListProduct(req)
	if err != nil {
		log.Printf("ERROR GetListProduct: %v", err)
		responseCode = http.StatusInternalServerError
		return
	}
	responseCode = http.StatusOK
	return
}
