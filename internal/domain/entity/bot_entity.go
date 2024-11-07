package entity

type (
	ServiceSupplier struct {
		Name         string
		BusinessType string
		Address      string
		Email        string
		PhoneNumber  string
	}

	SearchFilters struct {
		RepairType string
		City       string
	}
)
