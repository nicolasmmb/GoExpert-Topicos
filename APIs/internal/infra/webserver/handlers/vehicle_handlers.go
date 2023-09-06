package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.dev/nicolasmmb/GoExpert-Topicos/internal/dto"
	"github.dev/nicolasmmb/GoExpert-Topicos/internal/entity"
	"github.dev/nicolasmmb/GoExpert-Topicos/internal/infra/database"
	"github.dev/nicolasmmb/GoExpert-Topicos/pkg/filters"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type VehicleHandler struct {
	VehicleDB database.VehicleInterface
}

func NewVehicleHandler(db database.VehicleInterface) *VehicleHandler {
	return &VehicleHandler{VehicleDB: db}
}

func (v *VehicleHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if sort == "" {
		sort = "asc"
	}
	if sort != "asc" && sort != "desc" {
		http.Error(w, "sort must be asc or desc", http.StatusBadRequest)
		return
	}

	users, err := v.VehicleDB.FindAll(intPage, intLimit, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var vehicleOutput []dto.VehicleOutput

	for _, v := range users {
		vehicleOutput = append(vehicleOutput, dto.VehicleOutput{
			ID:             v.ID.String(),
			Value:          v.Value,
			Brand:          v.Brand,
			Model:          v.Model,
			ModelYear:      v.ModelYear,
			Fuel:           v.Fuel,
			FipeCode:       v.FipeCode,
			ReferenceMonth: v.ReferenceMonth,
			VehicleType:    v.VehicleType,
		})
	}

	json.NewEncoder(w).Encode(vehicleOutput)
	w.WriteHeader(http.StatusOK)
}

func (v *VehicleHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	db, _ := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	pag := filters.ParsePaginationRequest(r)

	vehicles := []entity.Vehicle{}

	ft := []string{"id", "brand", "model_year", "fuel", "fipe_code"}

	tx := db.Begin()
	tx = tx.Model(&entity.Vehicle{})

	for _, f := range ft {
		definition := filters.ParseParamRequest(r, f)
		if definition == nil {
			continue
		}
		query, value := definition.PreBuildQuery()

		log.Println(value)

		tx = tx.Where(query, value)

	}
	tx = tx.Offset(pag.Offset()).Limit(pag.Limit)

	tx.Find(&vehicles)

	var vehicleOutput []dto.VehicleOutput

	for _, v := range vehicles {
		vehicleOutput = append(vehicleOutput, dto.VehicleOutput{
			ID:             v.ID.String(),
			Value:          v.Value,
			Brand:          v.Brand,
			Model:          v.Model,
			ModelYear:      v.ModelYear,
			Fuel:           v.Fuel,
			FipeCode:       v.FipeCode,
			ReferenceMonth: v.ReferenceMonth,
			VehicleType:    v.VehicleType,
		})
	}

	json.NewEncoder(w).Encode(vehicleOutput)
	w.WriteHeader(http.StatusOK)

}
