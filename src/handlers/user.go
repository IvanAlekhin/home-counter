package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"home-counter/src/models"
	"home-counter/src/services"
	"net/http"
	"strconv"
)


func Counter (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there...")
}


func UserDataHandler (w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(models.UserData)
	userConfig, err := services.GetUserConfig(user)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	// if even one tariff not set, user should create it
	if !userConfig.ElectricityTariff.Valid || !userConfig.HotWaterTariff.Valid || !userConfig.ColdWaterTariff.Valid {
		http.Error(w, "User tariffs not found. Create it.", 404)
		return
	}
	res, err := json.Marshal(userConfig)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(res))
}

func UserTariffsHandler (w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("User").(models.UserData)
	tariffs := models.UserTariffsRequest{}
	err := json.NewDecoder(r.Body).Decode(&tariffs)

	validate := validator.New()
	err = validate.Struct(tariffs)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), 403)
		return
		}

	err = services.CreateOrUpdateUserTariffs(user, tariffs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "success")
}


func UserMeterCountHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(models.UserData)
	elMet, err := strconv.Atoi(r.URL.Query().Get("electricity"))

	if err != nil {
		fmt.Fprintf(w, "can't parse electricity param %s", r.URL.Query().Get("electricity"))
		http.Error(w, "", 403)
		return
	}
	HWMet, err := strconv.Atoi(r.URL.Query().Get("hot_water"))
	if err != nil {
		fmt.Fprintf(w, "can't parse hot_water param %s", r.URL.Query().Get("hot_water"))
		http.Error(w, "", 403)
		return
	}
	CWMet, err := strconv.Atoi(r.URL.Query().Get("cold_water"))
	if err != nil {
		fmt.Fprintf(w, "can't parse cold_water param %s", r.URL.Query().Get("cold_water"))
		http.Error(w, "", 403)
		return
	}

	metersData := models.UserNewMetersReq{
		ElectricityMeter: &elMet,
		HotWaterMeter: &HWMet,
		ColdWaterMeter: &CWMet,
	}
	validate := validator.New()
	err = validate.Struct(metersData)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), 403)
		return
	}

	res, err := services.CountUserPayment(user, metersData)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return
	}

	fmt.Fprint(w, "New sum to pay is ", res)

}

func CreateOrUpdateUserMetersHandler (w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("User").(models.UserData)
	meters := models.UserNewMetersReq{}
	err := json.NewDecoder(r.Body).Decode(&meters)

	validate := validator.New()
	err = validate.Struct(meters)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		http.Error(w, validationErrors.Error(), 403)
		return
	}

	err = services.CreateOrUpdateUserMeters(user, meters)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "Succeess")

}