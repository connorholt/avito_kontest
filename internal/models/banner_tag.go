package models

import (
	"database/sql"
)

type BannerTagModel struct {
	DB *sql.DB
}

func (m *BannerTagModel) Create(bannerID, featureID int, tagIDs []int) error {
	stmt := `INSERT INTO banner_feature_tag (banner_id, feature_id, tag_id) VALUES ($1, $2, $3)`
	for i := range tagIDs {
		_, err := m.DB.Exec(stmt, bannerID, featureID, tagIDs[i])
		if err != nil {
			if IsErrorCode(err, UniqueViolationErr) {
				return ErrDuplicateFeatureTag
			}
			return err
		}
	}
	return nil
}

func (m *BannerTagModel) Delete(BannerID int) error {
	stmt := `DELETE FROM banner_feature_tag where banner_id = $1 `
	_, err := m.DB.Exec(stmt, BannerID)
	if err != nil {
		return err
	}

	return nil
}
