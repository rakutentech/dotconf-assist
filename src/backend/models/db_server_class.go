package models

import (
// "github.com/rakutentech/dotconf-assist/src/backend/settings"
// "time"
)

type ID struct {
	ID int
}

func saveServerClassForwarders(serverClassID int, forwarderIDs []int) error {
	for _, id := range forwarderIDs {
		var scFwdr ServerClassForwarder
		scFwdr.ServerClassID = serverClassID
		scFwdr.ForwarderID = id
		res := mysqldb.Save(&scFwdr)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func getServerClassFwdrIDs(scID int) ([]int, error) {
	var fwdrIDs []int
	var ID int
	// var IDs []ID
	// rows, err := mysqldb.Table("server_class_forwarders").Select("forwarder_id").Where("server_class_id = ?", scID).Scan(&IDs)
	rows, err := mysqldb.Table("server_class_forwarders").Select("forwarder_id").Where("server_class_id = ?", scID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&ID)
		fwdrIDs = append(fwdrIDs, ID)
	}
	return fwdrIDs, nil
}

func getServerClassFwdrs(scID int) ([]string, error) {
	var fwdrs []string
	var forwarders []Forwarder
	res := mysqldb.Where("id IN (?)", mysqldb.Table("server_class_forwarders").
		Select("forwarder_id").Where("server_class_id = ?", scID).
		QueryExpr()).Find(&forwarders)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, fwdr := range forwarders {
		fwdrs = append(fwdrs, fwdr.Name)
	}
	return fwdrs, nil
}

func updateServerClassForwarders(serverClassID int, forwarderIDs []int) error {
	var IDs2Delete []int
	var IDs2Add []int
	fwdrIDs, _ := getServerClassFwdrIDs(serverClassID)
	for _, oldID := range fwdrIDs {
		var found = false
		for _, newID := range forwarderIDs {
			if newID == oldID {
				found = true
				break
			}
		}
		if !found { //not found, to delete
			IDs2Delete = append(IDs2Delete, oldID)
		}
	}

	for _, newID := range forwarderIDs {
		var found = false
		for _, oldID := range fwdrIDs {
			if newID == oldID {
				found = true
				break
			}
		}
		if !found { //not found, to add
			IDs2Add = append(IDs2Add, newID)
		}
	}

	err := saveServerClassForwarders(serverClassID, IDs2Add)
	err = deleteServerClassForwarders(serverClassID, IDs2Delete)
	return err
}

func deleteServerClassForwarders(serverClassID int, forwarderIDs []int) error {
	for _, ID := range forwarderIDs {
		mysqldb.Where("server_class_id = ? AND forwarder_id = ?", serverClassID, ID).Delete(ServerClassForwarder{})
	}
	return nil
}

func SaveServerClass(sc ServerClass) error {
	res := mysqldb.Save(&sc)
	if res.Error != nil {
		return res.Error
	}
	var scID ID
	mysqldb.Table("server_classes").Select("id").Where("name = ? AND user_name = ?", sc.Name, sc.UserName).Scan(&scID)
	saveServerClassForwarders(scID.ID, sc.ForwarderIDs)
	return nil
}

func GetServerClasses(envUser []string, isAdmin bool) ([]ServerClass, error) {
	var serverClasses []ServerClass
	if isAdmin {
		res := mysqldb.Where("env = ?", envUser[0]).Order("id").Find(&serverClasses)
		if res.Error != nil { // res.Error is nil even if no record found
			return nil, res.Error
		}
	} else {
		res := mysqldb.Where("env = ? AND user_name = ?", envUser[0], envUser[1]).Order("id").Find(&serverClasses)
		if res.Error != nil { // res.Error is nil even if no record found
			return nil, res.Error
		}
	}

	for i, _ := range serverClasses {
		// serverClasses[i].ForwarderIDs, _ = getServerClassFwdrIDs(serverClasses[i].ID)
		serverClasses[i].Forwarders, _ = getServerClassFwdrs(serverClasses[i].ID)
	}

	return serverClasses, nil
}

func GetServerClassesByIDs(scIDs []int, getForwarders bool) ([]ServerClass, error) {
	var serverClasses []ServerClass
	res := mysqldb.Table("server_classes").Where("id IN (?)", scIDs).Find(&serverClasses)
	if res.Error != nil { // res.Error is nil even if no record found
		return nil, res.Error
	}
	if getForwarders {
		for i, sc := range serverClasses {
			var forwarders []Forwarder
			res := mysqldb.Select("name").Where("id IN (?)", mysqldb.Table("server_class_forwarders").
				Select("forwarder_id").
				Where("server_class_id = ?", sc.ID).
				QueryExpr()).
				Find(&forwarders)
			if res.Error != nil { // res.Error is nil even if no record found
				return nil, res.Error
			}
			for _, fwdr := range forwarders {
				serverClasses[i].Forwarders = append(serverClasses[i].Forwarders, fwdr.Name)
			}
		}
	}

	// settings.WriteDebugLog(serverClasses)
	return serverClasses, nil
}

func GetServerClass(name, user string) (ServerClass, error) {
	var serverClass ServerClass
	res := mysqldb.Where("name = ? AND user_name = ?", name, user).Find(&serverClass)
	if res.Error != nil { //record not found
		return ServerClass{}, res.Error
	}
	return serverClass, nil
}

func UpdateServerClass(name, user string, newServerClass ServerClass) error {
	var serverClass ServerClass
	res := mysqldb.Where("name = ? AND user_name = ?", name, user).Find(&serverClass)
	if res.Error != nil { //record not found
		return res.Error
	}
	serverClass.Name = newServerClass.Name
	res = mysqldb.Save(&serverClass)
	if res.Error != nil {
		return res.Error
	}
	return updateServerClassForwarders(serverClass.ID, newServerClass.ForwarderIDs)
}

func DeleteServerClass(name, user string) error {
	serverClass, err := GetServerClass(name, user)
	if err != nil {
		return err
	}
	res := mysqldb.Delete(&serverClass)
	if res.Error != nil {
		return res.Error
	}
	return mysqldb.Where("server_class_id = ?", serverClass.ID).Delete(ServerClassForwarder{}).Error
}
