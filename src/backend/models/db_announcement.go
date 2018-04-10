package models

import (
// "github.com/rakutentech/dotconf-assist/src/backend/settings"
)

func SaveAnnouncement(announcement Announcement) error {
	res := mysqldb.Save(&announcement)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func GetAnnouncements() ([]Announcement, error) {
	var announcements []Announcement
	res := mysqldb.Order("created_at DESC").Find(&announcements)
	if res.Error != nil { // res.Error is nil even if no record found
		return nil, res.Error
	}

	return announcements, nil
}

func GetAnnouncement(id string) (Announcement, error) {
	var announcement Announcement
	res := mysqldb.Where("id = ?", id).Find(&announcement)
	if res.Error != nil { //record not found
		return Announcement{}, res.Error
	}
	return announcement, nil
}

func UpdateAnnouncement(id string, newAnnouncement Announcement) error {
	var announcement Announcement
	res := mysqldb.Where("id = ? ", id).Find(&announcement)
	if res.Error != nil { //record not found
		return res.Error
	}
	announcement.Content = newAnnouncement.Content
	return SaveAnnouncement(announcement)
}

func DeleteAnnouncement(id string) error {
	announcement, err := GetAnnouncement(id)
	if err != nil {
		return err
	}
	return mysqldb.Delete(&announcement).Error
}
