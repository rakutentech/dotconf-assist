package models

import (
// "github.com/rakutentech/dotconf-assist/src/backend/settings"
)

func SaveUnitPrice(unitPrice UnitPrice) error {
	res := mysqldb.Save(&unitPrice)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func GetUnitPrices() ([]UnitPrice, error) {
	var unitPrices []UnitPrice
	res := mysqldb.Order("id").Find(&unitPrices)
	if res.Error != nil { // res.Error is nil even if no record found
		return nil, res.Error
	}
	return unitPrices, nil
}

func GetUnitPrice(id int) (UnitPrice, error) {
	var unitPrice UnitPrice
	res := mysqldb.Where("id = ?", id).Find(&unitPrice)
	if res.Error != nil { //record not found
		return UnitPrice{}, res.Error
	}
	return unitPrice, nil
}

func UpdateUnitPrice(id int, newUnitPrice UnitPrice) error {
	var unitPrice UnitPrice
	res := mysqldb.Where("id = ?", id).Find(&unitPrice)
	if res.Error != nil { //record not found
		return res.Error
	}
	unitPrice.ServicePrice = newUnitPrice.ServicePrice
	unitPrice.StoragePrice = newUnitPrice.StoragePrice
	res = mysqldb.Save(&unitPrice)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteUnitPrice(id int) error {
	unitPrice, err := GetUnitPrice(id)
	if err != nil {
		return err
	}
	return mysqldb.Delete(&unitPrice).Error
}
