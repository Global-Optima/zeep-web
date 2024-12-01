package types

type EmployeeType string

const (
	StoreEmployee     EmployeeType = "Store"
	WarehouseEmployee EmployeeType = "Warehouse"
)

func ToEmployeeType(s string) EmployeeType {
	return EmployeeType(s)
}

func ToString(t EmployeeType) string {
	return string(t)
}
