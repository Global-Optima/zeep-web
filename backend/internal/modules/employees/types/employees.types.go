package types

type EmployeeType string

const (
	StoreEmployee     EmployeeType = "STORE"
	WarehouseEmployee EmployeeType = "WAREHOUSE"
)

func ToEmployeeType(s string) EmployeeType {
	return EmployeeType(s)
}

func ToString(t EmployeeType) string {
	return string(t)
}
