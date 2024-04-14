package models

import (
	"database/sql"
)

type BannerTagModel struct {
	DB *sql.DB
}

func (m *BannerTagModel) Create(BannerID int, tagIDs []int) error {
	stmt := `INSERT INTO banner_tag (banner_id, tag_id) VALUES ($1, $2)`
	for i := range tagIDs {
		_, err := m.DB.Exec(stmt, BannerID, tagIDs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *BannerTagModel) Delete(BannerID int) error {
	stmt := `DELETE FROM banner_tag where banner_id = $1 `
	_, err := m.DB.Exec(stmt, BannerID)
	if err != nil {
		return err
	}

	return nil
}

func (m *BannerTagModel) Insert(bannerID int, tagIDs []int) error {
	stmt := `INSERT INTO banner_tag VALUES ($1, $2)`
	for i := range tagIDs {
		_, err := m.DB.Exec(stmt, bannerID, tagIDs[i])
		if err != nil {
			return err
		}
	}
	return nil
}
