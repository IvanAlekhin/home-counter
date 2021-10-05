package models

type UserTariffsRequest struct {
	ElectricityTariff float64 `validate:"required" json:"electricity"`
	HotWaterTariff float64 `validate:"required" json:"hot_water"`
	ColdWaterTariff float64 `validate:"required" json:"cold_water"`
	OutWaterTariff float64 `validate:"required" json:"out_water"`
	InternetTariff float64 `json:"internet"`
}

type UserNewMetersReq struct {
	ElectricityMeter *int `schema:"electricity" json:"electricity" validate:"min=0,required"`
	HotWaterMeter *int `schema:"hot_water" json:"hot_water" validate:"min=0,required"`
	ColdWaterMeter *int `schema:"cold_water" json:"cold_water" validate:"min=0,required"`
}
