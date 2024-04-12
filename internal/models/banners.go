package models

import (
	"database/sql"
	"errors"
	"time"
)

type BannerModel struct {
	DB *sql.DB
}

type Banner struct {
	Id          int        `json:"banner_id"`
	Title       string     `json:"title"`
	Text        string     `json:"text"`
	URL         string     `json:"url"`
	Createad_at *time.Time `json:",omitempty"`
	Updated_at  *time.Time `json:",omitempty"`
}

func (m *BannerModel) Get(tagId int, featureId int) (*Banner, error) {
	stmt := `SELECT title, text, url from banners b 
                               join features f on b.feature_id = f.id 
                               join banner_tag bt on b.id = bt.banner_id
                               join tags t on t.id = bt.tag_id
								where feature_id = $1 and tag_id = $2`
	b := Banner{}
	row := m.DB.QueryRow(stmt, featureId, tagId)
	var titleNull, descNull sql.NullString
	err := row.Scan(&titleNull, &descNull, &b.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	if titleNull.Valid {
		b.Text = titleNull.String
	}
	if descNull.Valid {
		b.Text = descNull.String
	}
	return &b, nil
}

//func (m *BannerModel) GetAllBanners(tagID, featureID, limit int) ([]*Banner, error) {
//
//	stmt := `SELECT id, tag_id, feature_id, title, text, url from banners b
//                               join features f on b.feature_id = f.id
//                               join banner_tag bt on b.id = bt.banner_id
//                               join tags t on t.id = bt.tag_id
//								where feature_id = $1 and tag_id = $2 limit $3`
//}
