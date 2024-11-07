package repository

import (
	"TelegrammBot/internal/domain/entity"
	"context"
	"database/sql"
	"log"
)

type BotRepository struct {
	db *sql.DB
}

func NewBotRepository(db *sql.DB) *BotRepository {
	return &BotRepository{db: db}
}

func (r *BotRepository) SearchAllByFilter(ctx context.Context) (string, error) {
	var suppliers []entity.ServiceSupplier
	var rows *sql.Rows
	var err error

	filters, ok := ctx.Value("searchFilters").(entity.SearchFilters)
	if !ok {
		return "Не удалось установить фильтры поиска", err
	}

	var cityID int
	var attID int
	err = r.db.QueryRow(`SELECT c.id, a.id FROM cities c join attendances a on a.name = $1 WHERE c.name = $2`, filters.RepairType, filters.City).Scan(&cityID, &attID)

	query := `SELECT name, business_type, address, email, phone_number FROM suppliers where att_id = $1 and city_id = $2`
	rows, err = r.db.Query(query, attID, cityID)

	if err != nil {
		log.Printf("Error while searching all suppliers: %v", err)
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var supplier entity.ServiceSupplier
		err = rows.Scan(&supplier.Name, &supplier.BusinessType, &supplier.Address, &supplier.Email, &supplier.PhoneNumber)
		if err != nil {
			log.Printf("Error while scanning supplier: %v", err)
		}
		suppliers = append(suppliers, supplier)
	}

	return formatServiceSuppliers(suppliers), nil
}

func formatServiceSuppliers(suppliers []entity.ServiceSupplier) string {
	var result string
	for _, supplier := range suppliers {
		result += "Наименование организации: " + supplier.Name + "\n" +
			"Тип организации: " + supplier.BusinessType + "\n" +
			"Адрес: " + supplier.Address + "\n" +
			"Email: " + supplier.Email + "\n" +
			"Телефон: " + supplier.PhoneNumber + "\n" +
			"--------------------------------------------------------" + "\n"
	}

	return result
}
