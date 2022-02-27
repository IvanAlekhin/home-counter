package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"home-counter/src/models"
)

var (
	configError = errors.New("can't find previous meters. Add it using post req /user/meters")
	dbReqError  = errors.New("something wrong with db request")
)

func GetUserConfig(user models.UserData) (models.UserConfig, error) {
	var q = `SELECT u.name, ut.cold_water::decimal, ut.hot_water::decimal, ut.out_water::decimal, ut.internet::decimal, 
       ut.electricity::decimal, umd.electricity as electricity_meter, umd.cold_water as cold_water_meter, 
       umd.hot_water as hot_water_meter
			 FROM "user" as u 
		     LEFT JOIN user_tariff as ut on u.id = ut.user_id
		     LEFT JOIN user_meter_data as umd on u.id = umd.user_id
		     WHERE u.id = $1`
	var userConfig = models.UserConfig{}

	err := models.DB.QueryRow(context.Background(), q, user.Id).Scan(&userConfig.Name, &userConfig.ColdWaterTariff,
		&userConfig.HotWaterTariff, &userConfig.OutWaterTariff, &userConfig.InternetTariff,
		&userConfig.ElectricityTariff, &userConfig.Electricity, &userConfig.ColdWater, &userConfig.HotWater)
	if err != nil {
		panic(err)
	}
	if !userConfig.HotWater.Valid || !userConfig.Electricity.Valid || !userConfig.ColdWater.Valid {
		return userConfig, configError
	}
	return userConfig, nil
}

func CreateOrUpdateUserMeters(user models.UserData, meters models.UserNewMetersReq) error {

	q := "UPDATE user_meter_data SET electricity=$1, cold_water=$2, hot_water=$3 WHERE user_id=$4"

	res, err := models.DB.Exec(context.Background(), q, meters.ElectricityMeter, meters.ColdWaterMeter, meters.HotWaterMeter, user.Id)
	if err != nil {
		return dbReqError
	}
	n := res.RowsAffected()
	if n == 0 {
		q = "INSERT INTO user_meter_data (user_id, electricity, cold_water, hot_water) VALUES ($1, $2, $3, $4)"
		_, err = models.DB.Exec(context.Background(), q, user.Id, meters.ElectricityMeter, meters.ColdWaterMeter, meters.HotWaterMeter)
		if err != nil {
			return dbReqError
		}
	}
	return nil
}

func CreateOrUpdateUserTariffs(user models.UserData, tariffs models.UserTariffsRequest) error {
	var getQ = `SELECT user_id FROM user_tariff WHERE user_id = $1`
	var updateQ = `UPDATE user_tariff 
                    SET electricity=$2, cold_water=$3, hot_water=$4, out_water=$5, internet=$6
                    WHERE user_id=$1`
	var insertQ = `INSERT INTO user_tariff (user_id, electricity, cold_water, hot_water, out_water, internet) 
             VALUES ($1, $2, $3, $4, $5, $6)`
	var q string
	var userId string

	dbErr := models.DB.QueryRow(context.Background(), getQ, user.Id).Scan(&userId)
	if dbErr != nil {
		switch dbErr {
		case sql.ErrNoRows:
			q = insertQ
		default:
			return dbReqError
		}
	} else if dbErr == nil {
		q = updateQ
	}

	_, err := models.DB.Exec(context.Background(), q, user.Id, tariffs.ElectricityTariff, tariffs.ColdWaterTariff, tariffs.HotWaterTariff,
		tariffs.OutWaterTariff, tariffs.InternetTariff)

	if err != nil {
		return dbReqError
	}

	return nil
}

func CountUserPayment(user models.UserData, meters models.UserNewMetersReq) (string, error) {
	// берем данные из запроса (помним про валидацию)
	userConfig, err := GetUserConfig(user)

	if err != nil {
		return "", err
	}

	// get other data from user config. counting, returning schema with integer. save previous data in request
	res := Counting(meters, userConfig)

	err = CreateOrUpdateUserMeters(user, meters)
	if err != nil {
		return res, err
	}

	fmt.Println(res)
	return res, nil
}

func Counting(newMeters models.UserNewMetersReq, userConfig models.UserConfig) string {
	electricity := int64(*newMeters.ElectricityMeter) - userConfig.Electricity.Int64
	hotWater := int64(*newMeters.HotWaterMeter) - userConfig.HotWater.Int64
	coldWater := int64(*newMeters.ColdWaterMeter) - userConfig.ColdWater.Int64
	outWater := hotWater + coldWater

	result := decimal.NewFromFloat(userConfig.ElectricityTariff.Float64).Mul(decimal.NewFromInt(electricity)).Add(
		decimal.NewFromFloat(userConfig.HotWaterTariff.Float64).Mul(decimal.NewFromInt(hotWater))).Add(
		decimal.NewFromFloat(userConfig.ColdWaterTariff.Float64).Mul(decimal.NewFromInt(coldWater))).Add(
		decimal.NewFromFloat(userConfig.OutWaterTariff.Float64).Mul(decimal.NewFromInt(outWater))).Add(
		decimal.NewFromFloat(userConfig.InternetTariff.Float64))

	return result.String()
}
