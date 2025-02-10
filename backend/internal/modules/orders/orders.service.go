package orders

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sync"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/storeStocks"

	"github.com/Global-Optima/zeep-web/backend/pkg/utils/censor"

	storeAdditives "github.com/Global-Optima/zeep-web/backend/internal/modules/additives/storeAdditivies"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/product/storeProducts"
	"github.com/golang/freetype/truetype"
	"github.com/jung-kurt/gofpdf"

	"go.uber.org/zap"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"

	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications/details"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils/pdf"
)

type OrderService interface {
	GetOrders(filter types.OrdersFilterQuery) ([]types.OrderDTO, error)
	GetAllBaristaOrders(filter types.GetBaristaOrdersFilter) ([]types.OrderDTO, error)
	GetSubOrders(orderID uint) ([]types.SuborderDTO, error)
	GetStatusesCount(storeID uint) (types.OrderStatusesCountDTO, error)
	CreateOrder(storeId uint, createOrderDTO *types.CreateOrderDTO) (*data.Order, error)
	CompleteSubOrder(orderID, subOrderID uint) error
	CompleteSubOrderByBarcode(subOrderID uint) (*types.SuborderDTO, error)
	GeneratePDFReceipt(orderID uint) ([]byte, error)
	GetOrderBySubOrder(subOrderID uint) (*data.Order, error)
	GetOrderById(orderId uint) (types.OrderDTO, error)

	GetOrderDetails(orderID uint) (*types.OrderDetailsDTO, error)
	ExportOrders(filter *types.OrdersExportFilterQuery) ([]types.OrderExportDTO, error)
	GenerateSuborderBarcodePDF(suborderID uint) ([]byte, error)

	AcceptSubOrder(subOrderID uint) error
	AdvanceSubOrderStatus(subOrderID uint) (*types.SuborderDTO, error)
}

type orderValidationResults struct {
	ProductPrices  map[uint]float64
	ProductNames   map[uint]string
	AdditivePrices map[uint]float64
	AdditiveNames  map[uint]string
}

type orderService struct {
	orderRepo           OrderRepository
	storeProductRepo    storeProducts.StoreProductRepository
	storeAdditiveRepo   storeAdditives.StoreAdditiveRepository
	storeStockRepo      storeStocks.StoreStockRepository
	notificationService notifications.NotificationService
	logger              *zap.SugaredLogger
}

func NewOrderService(
	orderRepo OrderRepository,
	storeProductRepo storeProducts.StoreProductRepository,
	storeAdditiveRepo storeAdditives.StoreAdditiveRepository,
	storeStockRepo storeStocks.StoreStockRepository,
	notificationService notifications.NotificationService,
	logger *zap.SugaredLogger,
) OrderService {
	return &orderService{
		orderRepo:           orderRepo,
		storeProductRepo:    storeProductRepo,
		storeAdditiveRepo:   storeAdditiveRepo,
		storeStockRepo:      storeStockRepo,
		notificationService: notificationService,
		logger:              logger,
	}
}

var (
	fontFace      font.Face
	fontInitError error
	fontInitOnce  sync.Once
)

func (s *orderService) GetOrders(filter types.OrdersFilterQuery) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetOrders(filter)
	if err != nil {
		return nil, err
	}

	orderDTOs := make([]types.OrderDTO, 0)
	for _, order := range orders {
		orderDTOs = append(orderDTOs, types.ConvertOrderToDTO(&order))
	}

	return orderDTOs, nil
}

func (s *orderService) GetAllBaristaOrders(filter types.GetBaristaOrdersFilter) ([]types.OrderDTO, error) {
	orders, err := s.orderRepo.GetAllBaristaOrders(filter)

	if err != nil {
		wrappedErr := fmt.Errorf("error getting barista orders: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	orderDTOs := make([]types.OrderDTO, 0)

	for _, order := range orders {
		orderDTOs = append(orderDTOs, types.ConvertOrderToDTO(&order))
	}

	return orderDTOs, nil
}

func (s *orderService) GetSubOrders(orderID uint) ([]types.SuborderDTO, error) {
	suborders, err := s.orderRepo.GetSubOrdersByOrderID(orderID)
	if err != nil {
		wrappedErr := fmt.Errorf("error getting suborders by orderID: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	var subOrderDTOs []types.SuborderDTO
	for _, suborder := range suborders {
		subOrderDTOs = append(subOrderDTOs, types.ConvertSuborderToDTO(&suborder))
	}
	return subOrderDTOs, nil
}

func (s *orderService) CreateOrder(storeID uint, createOrderDTO *types.CreateOrderDTO) (*data.Order, error) {
	censorValidator := censor.GetCensorValidator()

	if err := censorValidator.ValidateText(createOrderDTO.CustomerName); err != nil {
		s.logger.Error(err)
		return nil, err
	}

	storeProductSizeIDs, storeAdditiveIDs := RetrieveIDs(*createOrderDTO)

	validations, err := s.ValidationResults(storeID, storeProductSizeIDs, storeAdditiveIDs)
	if err != nil {
		wrappedErr := fmt.Errorf("validation failed: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	createOrderDTO.StoreID = storeID
	order, total := types.ConvertCreateOrderDTOToOrder(createOrderDTO, validations.ProductPrices, validations.AdditivePrices)

	order.Status = data.OrderStatusPending
	order.Total = total

	err = s.orderRepo.CreateOrder(&order)
	if err != nil {
		wrappedErr := fmt.Errorf("error creating order: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	notificationDetails := &details.NewOrderNotificationDetails{
		BaseNotificationDetails: details.BaseNotificationDetails{
			ID:           order.StoreID,
			FacilityName: order.Store.Name,
		},
		CustomerName: createOrderDTO.CustomerName,
		OrderID:      order.ID,
	}
	err = s.notificationService.NotifyNewOrder(notificationDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to notify new order: %w", err)
	}

	return &order, nil
}

func (s *orderService) CompleteSubOrder(orderID, subOrderID uint) error {
	err := s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusCompleted)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to complete suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return wrappedErr
	}

	allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(orderID, subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return wrappedErr
	}

	if allCompleted {
		order, err := s.orderRepo.GetOrderById(orderID)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to get order by id: %w", err)
			s.logger.Error(wrappedErr.Error())
			return wrappedErr
		}

		stockMap := make(map[uint]*data.StoreStock)

		err = s.deductProductSizeIngredientsFromStock(order, stockMap)
		if err != nil {
			return fmt.Errorf("failed to deduct ingredients for products: %w", err)
		}

		err = s.deductAdditiveIngredientsFromStock(order, stockMap)
		if err != nil {
			return fmt.Errorf("failed to deduct ingredients for additives: %w", err)
		}

		s.notifyLowStockIngredients(order, stockMap)

		var newStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		} else {
			newStatus = data.OrderStatusCompleted
		}

		err = s.orderRepo.UpdateOrderStatus(orderID, newStatus)

		if err != nil {
			wrappedErr := fmt.Errorf("failed to update order status: %w", err)
			s.logger.Error(wrappedErr.Error())
			return wrappedErr
		}
	}

	return nil
}

// initFont ensures we parse the font only once (thread-safe).
func initFont() {
	fontInitOnce.Do(func() {
		f, err := truetype.Parse(goregular.TTF)
		if err != nil {
			fontInitError = fmt.Errorf("failed to parse goregular font: %w", err)
			return
		}

		// Adjust Size as needed. This is a typical 12pt usage.
		fontFace = truetype.NewFace(f, &truetype.Options{
			Size: 12,
		})
	})
}

func (s *orderService) GenerateSuborderBarcodePDF(suborderID uint) ([]byte, error) {
	// Initialize font first
	initFont()
	if fontInitError != nil {
		return nil, fontInitError
	}

	// 1. Fetch suborder from repository
	suborder, err := s.orderRepo.GetSuborderByID(suborderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve suborder (id=%d): %w", suborderID, err)
	}
	if suborder == nil {
		return nil, fmt.Errorf("no suborder found for ID %d", suborderID)
	}
	// if suborder.Status != data.SubOrderStatusPending {
	// 	return nil, fmt.Errorf("barcode generation allowed only for pending suborders")
	// }

	// 2. Encode the barcode data using Code-128
	barcodeData := fmt.Sprintf("suborder-%d", suborder.ID)
	bcode, err := code128.Encode(barcodeData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode barcode data %q: %w", barcodeData, err)
	}

	// 3. Scale barcode to fit exactly within the image dimensions
	//    Suppose we want the final label to be 300px wide x 120px tall
	const labelWidth = 300
	const labelHeight = 120
	scaledBcode, err := barcode.Scale(bcode, labelWidth, labelHeight-20) // leave space for text
	if err != nil {
		return nil, fmt.Errorf("failed to scale barcode: %w", err)
	}

	// 4. Create a blank image with a white background
	finalImg := image.NewRGBA(image.Rect(0, 0, labelWidth, labelHeight))
	draw.Draw(finalImg, finalImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 5. Draw barcode on the image
	bcBounds := scaledBcode.Bounds()
	draw.Draw(finalImg, bcBounds, scaledBcode, image.Point{}, draw.Over)

	// 6. Draw the text under the barcode
	drawer := &font.Drawer{
		Dst:  finalImg,
		Src:  image.NewUniform(color.Black),
		Face: fontFace, // Make sure fontFace is defined somewhere
	}
	textWidth := drawer.MeasureString(barcodeData).Ceil()
	textX := (labelWidth - textWidth) / 2
	textY := labelHeight - 5

	drawer.Dot = fixed.Point26_6{
		X: fixed.I(textX),
		Y: fixed.I(textY),
	}
	drawer.DrawString(barcodeData)

	// 7. Encode image as PNG in memory
	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, finalImg); err != nil {
		return nil, fmt.Errorf("failed to encode final PNG: %w", err)
	}

	// 8. Convert PNG to PDF
	//    Instead of using "A4", we define a small custom page size:
	//
	//    1 pixel = ~0.264583 mm
	//    labelWidth (300 px)  ~ 79.375 mm
	//    labelHeight (120 px) ~ 31.75 mm
	//
	//    Add a small margin around your label if needed.
	const marginMM = 5.0
	imgWidthMM := float64(labelWidth) * 0.264583
	imgHeightMM := float64(labelHeight) * 0.264583

	// Create a new PDF with custom page size:
	// Use gofpdf.NewCustom so we can specify exact dimensions.
	pdfFile := gofpdf.NewCustom(&gofpdf.InitType{
		UnitStr: "mm",
		Size: gofpdf.SizeType{
			Wd: imgWidthMM + (2 * marginMM),
			Ht: imgHeightMM + (2 * marginMM),
		},
		FontDirStr: "",
	})
	pdfFile.SetAutoPageBreak(false, 0)
	pdfFile.AddPage()

	// Register and place the image
	imgOptions := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}
	pdfFile.RegisterImageOptionsReader("barcode.png", imgOptions, &imgBuf)

	// Place image with the margin as an offset
	pdfFile.ImageOptions(
		"barcode.png",
		marginMM,    // x-pos
		marginMM,    // y-pos
		imgWidthMM,  // image width in mm
		imgHeightMM, // image height in mm
		false,       // flow: false = free-floating
		imgOptions,
		0,
		"",
	)

	// 9. Output PDF as byte slice
	var pdfBuf bytes.Buffer
	if err := pdfFile.Output(&pdfBuf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfBuf.Bytes(), nil
}

func (s *orderService) CompleteSubOrderByBarcode(subOrderID uint) (*types.SuborderDTO, error) {
	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order for suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	err = s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusCompleted)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to complete suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	allCompleted, err := s.orderRepo.CheckAllSubordersCompleted(order.ID, subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to check suborder: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	if allCompleted {
		order, err := s.orderRepo.GetOrderById(order.ID)
		if err != nil {
			wrappedErr := fmt.Errorf("failed to get order by id: %w", err)
			s.logger.Error(wrappedErr.Error())
			return nil, wrappedErr
		}

		stockMap := make(map[uint]*data.StoreStock)

		err = s.deductProductSizeIngredientsFromStock(order, stockMap)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct ingredients for products: %w", err)
		}

		err = s.deductAdditiveIngredientsFromStock(order, stockMap)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct ingredients for additives: %w", err)
		}

		s.notifyLowStockIngredients(order, stockMap)

		var newStatus data.OrderStatus
		if order.DeliveryAddressID != nil {
			newStatus = data.OrderStatusInDelivery
		} else {
			newStatus = data.OrderStatusCompleted
		}

		err = s.orderRepo.UpdateOrderStatus(order.ID, newStatus)

		if err != nil {
			wrappedErr := fmt.Errorf("failed to update order status: %w", err)
			s.logger.Error(wrappedErr.Error())
			return nil, wrappedErr
		}
	}

	subOrder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return nil, err
	}
	response := types.ConvertSuborderToDTO(subOrder)

	return &response, nil
}

func (s *orderService) deductProductSizeIngredientsFromStock(order *data.Order, stockMap map[uint]*data.StoreStock) error {
	for _, suborder := range order.Suborders {
		updatedStocks, err := s.storeStockRepo.DeductStockByProductSizeTechCart(order.StoreID, suborder.StoreProductSizeID)
		if err != nil {
			return fmt.Errorf("failed to deduct from store stock: %w", err)
		}

		for _, stock := range updatedStocks {
			if existingStock, exists := stockMap[stock.IngredientID]; exists {
				existingStock.Quantity = stock.Quantity // Update latest quantity
			} else {
				stockMap[stock.IngredientID] = &stock
			}
		}
	}

	return nil
}

func (s *orderService) deductAdditiveIngredientsFromStock(order *data.Order, stockMap map[uint]*data.StoreStock) error {
	for _, suborder := range order.Suborders {
		for _, storeAdditive := range suborder.SuborderAdditives {
			updatedStocks, err := s.storeStockRepo.DeductStockByAdditiveTechCart(order.StoreID, storeAdditive.StoreAdditiveID)
			if err != nil {
				return fmt.Errorf("failed to deduct from store stock: %w", err)
			}

			for _, stock := range updatedStocks {
				if existingStock, exists := stockMap[stock.IngredientID]; exists {
					existingStock.Quantity = stock.Quantity // Update latest quantity
				} else {
					stockMap[stock.IngredientID] = &stock
				}
			}
		}
	}

	return nil
}

func (s *orderService) notifyLowStockIngredients(order *data.Order, stockMap map[uint]*data.StoreStock) {
	for _, stock := range stockMap {
		if stock.Quantity <= stock.LowStockThreshold {
			notificationDetails := &details.StoreWarehouseRunOutDetails{
				BaseNotificationDetails: details.BaseNotificationDetails{
					ID:           order.StoreID,
					FacilityName: order.Store.Name,
				},
				StockItem:   stock.Ingredient.Name,
				StockItemID: stock.IngredientID,
			}
			err := s.notificationService.NotifyStoreWarehouseRunOut(notificationDetails)
			if err != nil {
				s.logger.Errorf("failed to send notification: %v", err)
			}
		}
	}
}

func (s *orderService) GeneratePDFReceipt(orderID uint) ([]byte, error) {
	order, err := s.orderRepo.GetOrderById(orderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order details: %w", err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	detailsPDF := pdf.PDFReceiptDetails{
		OrderID:   order.ID,
		StoreID:   order.StoreID,
		OrderDate: order.CreatedAt.Format("2006-01-02 15:04:05"),
		Total:     order.Total,
	}

	for _, suborder := range order.Suborders {
		pdfSuborder := pdf.PDFSubOrder{
			ProductName: fmt.Sprintf("Product #%d", suborder.StoreProductSizeID),
			Price:       suborder.Price,
			Status:      string(suborder.Status),
		}

		for _, storeAdditive := range suborder.SuborderAdditives {
			pdfSuborder.Additives = append(pdfSuborder.Additives, pdf.PDFAdditive{
				Name:  fmt.Sprintf("Additive #%d", storeAdditive.StoreAdditiveID),
				Price: storeAdditive.Price,
			})
		}

		detailsPDF.SubOrders = append(detailsPDF.SubOrders, pdfSuborder)
	}

	return pdf.GeneratePDFReceipt(detailsPDF)
}

func (s *orderService) GetStatusesCount(storeID uint) (types.OrderStatusesCountDTO, error) {
	countsMap, err := s.orderRepo.GetStatusesCount(storeID)
	if err != nil {
		wrappedErr := fmt.Errorf("error couting statuses: %w", err)
		s.logger.Error(wrappedErr.Error())
		return types.OrderStatusesCountDTO{}, wrappedErr
	}

	dto := types.OrderStatusesCountDTO{
		PENDING:     countsMap[data.OrderStatusPending],
		PREPARING:   countsMap[data.OrderStatusPreparing],
		COMPLETED:   countsMap[data.OrderStatusCompleted],
		IN_DELIVERY: countsMap[data.OrderStatusInDelivery],
		DELIVERED:   countsMap[data.OrderStatusDelivered],
		CANCELLED:   countsMap[data.OrderStatusCancelled],
	}

	dto.ALL = dto.PENDING + dto.PREPARING + dto.COMPLETED + dto.IN_DELIVERY + dto.DELIVERED + dto.CANCELLED

	return dto, nil
}

func (s *orderService) ValidationResults(storeID uint, storeProductSizeIDs, storeAdditiveIDs []uint) (*orderValidationResults, error) {
	productPrices, productNames, err := ValidateStoreProductSizes(storeID, storeProductSizeIDs, s.storeProductRepo)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	additivePrices, additiveNames, err := ValidateStoreAdditives(storeID, storeAdditiveIDs, s.storeAdditiveRepo)
	if err != nil {
		return nil, fmt.Errorf("additive validation failed: %w", err)
	}

	return &orderValidationResults{
		ProductPrices:  productPrices,
		ProductNames:   productNames,
		AdditivePrices: additivePrices,
		AdditiveNames:  additiveNames,
	}, nil
}

func ValidateStoreAdditives(storeID uint, storeAdditiveIDs []uint, repo storeAdditives.StoreAdditiveRepository) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	additiveNames := make(map[uint]string)
	for _, id := range storeAdditiveIDs {
		storeAdditive, err := repo.GetStoreAdditiveByID(storeID, id)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid store additive ID: %d", id)
		}
		if storeAdditive == nil {
			return nil, nil, fmt.Errorf("store additive with ID %d is nil", id)
		}
		if storeAdditive.Additive.Name == "" {
			return nil, nil, fmt.Errorf("store additive with ID %d has an empty name", id)
		}

		price := storeAdditive.Additive.BasePrice
		if storeAdditive.StorePrice == nil {
			price = *storeAdditive.StorePrice
		}
		prices[id] = price
		additiveNames[id] = storeAdditive.Additive.Name
	}
	return prices, additiveNames, nil
}

func ValidateStoreProductSizes(storeID uint, storeProductSizeIDs []uint, repo storeProducts.StoreProductRepository) (map[uint]float64, map[uint]string, error) {
	prices := make(map[uint]float64)
	productNames := make(map[uint]string)
	for _, id := range storeProductSizeIDs {
		storeProductSize, err := repo.GetStoreProductSizeById(storeID, id)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid store product size ID: %d", id)
		}
		if storeProductSize == nil {
			return nil, nil, fmt.Errorf("store product size with ID %d is nil", id)
		}
		if storeProductSize.ProductSize.Product.Name == "" {
			return nil, nil, fmt.Errorf("product size with ID %d has an associated product with an empty name", id)
		}

		price := storeProductSize.ProductSize.BasePrice
		if storeProductSize.StorePrice == nil {
			price = *storeProductSize.StorePrice
		}

		prices[id] = price
		productNames[id] = storeProductSize.StoreProduct.Product.Name

	}
	return prices, productNames, nil
}

func RetrieveIDs(createOrderDTO types.CreateOrderDTO) ([]uint, []uint) {
	var storeProductSizeIDs []uint
	var storeAdditiveIDs []uint
	for _, product := range createOrderDTO.Suborders {
		storeProductSizeIDs = append(storeProductSizeIDs, product.StoreProductSizeID)
		storeAdditiveIDs = append(storeAdditiveIDs, product.StoreAdditivesIDs...)
	}
	return storeProductSizeIDs, storeAdditiveIDs
}

func (s *orderService) GetOrderBySubOrder(subOrderID uint) (*data.Order, error) {
	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order for suborder %d: %w", subOrderID, err)
		s.logger.Error(wrappedErr.Error())
		return nil, wrappedErr
	}

	return order, nil
}

func (s *orderService) GetOrderById(orderId uint) (types.OrderDTO, error) {
	order, err := s.orderRepo.GetOrderById(orderId)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to fetch order with ID %d: %w", orderId, err)
		s.logger.Error(wrappedErr.Error())
		return types.OrderDTO{}, wrappedErr
	}

	orderDTO := types.ConvertOrderToDTO(order)
	return orderDTO, nil
}

func (s *orderService) GetOrderDetails(orderID uint) (*types.OrderDetailsDTO, error) {
	order, err := s.orderRepo.GetOrderDetails(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, nil
	}

	return types.MapToOrderDetailsDTO(order), nil
}

func (s *orderService) ExportOrders(filter *types.OrdersExportFilterQuery) ([]types.OrderExportDTO, error) {
	orders, err := s.orderRepo.GetOrdersForExport(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders for export: %w", err)
	}

	exports := make([]types.OrderExportDTO, len(orders))
	for i, order := range orders {
		exports[i] = types.ToOrderExportDTO(&order, order.Store.Name)
	}

	return exports, nil
}

func (s *orderService) AcceptSubOrder(subOrderID uint) error {
	suborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve suborder: %w", err)
	}
	if suborder == nil {
		return fmt.Errorf("suborder %d not found", subOrderID)
	}

	if suborder.Status != data.SubOrderStatusPending {
		return fmt.Errorf("suborder %d is not pending", subOrderID)
	}

	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, data.SubOrderStatusPreparing); err != nil {
		return fmt.Errorf("failed to update suborder status: %w", err)
	}

	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	suborders, err := s.orderRepo.GetSubOrdersByOrderID(order.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch suborders for order %d: %w", order.ID, err)
	}

	allAccepted := true
	for _, so := range suborders {
		if so.Status == data.SubOrderStatusPending {
			allAccepted = false
			break
		}
	}

	// If all suborders have been accepted (i.e. transitioned to PREPARING or later)
	// and the order is still in the PENDING state, update the order's status.
	if allAccepted && order.Status == data.OrderStatusPending {
		if err := s.orderRepo.UpdateOrderStatus(order.ID, data.OrderStatusPreparing); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}
	return nil
}

var allowedTransitions = map[data.SubOrderStatus]data.SubOrderStatus{
	data.SubOrderStatusPending:   data.SubOrderStatusPreparing,
	data.SubOrderStatusPreparing: data.SubOrderStatusCompleted,
}

func (s *orderService) AdvanceSubOrderStatus(subOrderID uint) (*types.SuborderDTO, error) {
	suborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve suborder: %w", err)
	}
	if suborder == nil {
		return nil, fmt.Errorf("suborder %d not found", subOrderID)
	}

	currentStatus := suborder.Status
	nextStatus, ok := allowedTransitions[currentStatus]
	if !ok {
		return nil, fmt.Errorf("no allowed transition from status %s", currentStatus)
	}

	if err := s.orderRepo.UpdateSubOrderStatus(subOrderID, nextStatus); err != nil {
		return nil, fmt.Errorf("failed to update suborder status: %w", err)
	}

	order, err := s.orderRepo.GetOrderBySubOrderID(subOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve order for suborder %d: %w", subOrderID, err)
	}

	suborders, err := s.orderRepo.GetSubOrdersByOrderID(order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch suborders for order %d: %w", order.ID, err)
	}

	if nextStatus == data.SubOrderStatusCompleted {
		allCompleted := true
		for _, so := range suborders {
			if so.Status != data.SubOrderStatusCompleted {
				allCompleted = false
				break
			}
		}
		if allCompleted {
			var newOrderStatus data.OrderStatus
			// Example logic: if the order has a delivery address, mark it as IN_DELIVERY;
			// otherwise, mark it as COMPLETED.
			if order.DeliveryAddressID != nil {
				newOrderStatus = data.OrderStatusInDelivery
			} else {
				newOrderStatus = data.OrderStatusCompleted
			}
			if err := s.orderRepo.UpdateOrderStatus(order.ID, newOrderStatus); err != nil {
				return nil, fmt.Errorf("failed to update order status: %w", err)
			}
		}
	} else if nextStatus == data.SubOrderStatusPreparing {
		allNotPending := true
		for _, so := range suborders {
			if so.Status == data.SubOrderStatusPending {
				allNotPending = false
				break
			}
		}

		if allNotPending && order.Status == data.OrderStatusPending {
			if err := s.orderRepo.UpdateOrderStatus(order.ID, data.OrderStatusPreparing); err != nil {
				return nil, fmt.Errorf("failed to update order status: %w", err)
			}
		}
	}

	updatedSuborder, err := s.orderRepo.GetSuborderByID(subOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated suborder: %w", err)
	}
	dto := types.ConvertSuborderToDTO(updatedSuborder)
	return &dto, nil
}
