package models

import (
	"database/sql"
	"errors"
	"time"
)

type BannerModel struct {
	DB *sql.DB
}

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type Banner struct {
	Id        int        `json:"banner_id"`
	TagID     int        `json:"tag_ids"`
	FeatureID int        `json:"feature_id"`
	Content   Content    `json:"content"`
	IsActive  bool       `json:"is_active"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m *BannerModel) Get(tagId int, featureId int) (*Content, error) {
	stmt := `SELECT title, text, url from banners b 
                               join features f on b.feature_id = f.id 
                               join banner_tag bt on b.id = bt.banner_id
                               join tags t on t.id = bt.tag_id
								where feature_id = $1 and tag_id = $2`
	c := Content{}
	row := m.DB.QueryRow(stmt, featureId, tagId)

	err := row.Scan(&c.Title, &c.Text, &c.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return &c, nil
}

func (m *BannerModel) GetAll(tagID, featureID, limit, offset int) (*[]Banner, error) {
	stmt := `SELECT b.id, tag_id, feature_id, title, text, url, visible, created_at, updated_at from banners b
                              join features f on b.feature_id = f.id
                              join banner_tag bt on b.id = bt.banner_id
                              join tags t on t.id = bt.tag_id
								where case
								when $1 != -1 and $2 != -1 then tag_id = $1 and feature_id = $2
								when  $1 != -1 then tag_id = $1
								when $2 != -1 then feature_id = $2
								else true
								end 
								order by b.id
								limit case 
								when $3 != -1 then $3
								else NULL
								END
								offset case 
								when $4 != -1 then $4
								else 0
								END`

	capacity := 10
	if limit != -1 {
		capacity = limit
	}
	banners := make([]Banner, 0, capacity)
	rows, err := m.DB.Query(stmt, tagID, featureID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		b := Banner{}
		err := rows.Scan(&b.Id, &b.TagID, &b.FeatureID, &b.Content.Title, &b.Content.Text, &b.Content.URL, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		banners = append(banners, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &banners, nil
}

//stmt := `SELECT b.id, tag_id, feature_id, title, text, url, visible, created_at, updated_at from banners b
//                              join features f on b.feature_id = f.id
//                              join banner_tag bt on b.id = bt.banner_id
//                              join tags t on t.id = bt.tag_id
//								where ($1 != 0 and $2 != 0 and tag_id = $1 and feature_id = $2)
//								or ( $1 != 0 and tag_id = $1)
//								or ($2 != 0 and feature_id = $2)
//								order by b.id`
