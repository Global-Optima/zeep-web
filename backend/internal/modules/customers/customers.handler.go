package customers

type CustomerHandler struct {
	service CustomerService
}

func NewCustomerHandler(service CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}
