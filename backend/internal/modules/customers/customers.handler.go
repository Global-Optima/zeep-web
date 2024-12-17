package customers

type CustomerHandler struct {
	service CustomerService
}

func NewEmployeeHandler(service CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}
