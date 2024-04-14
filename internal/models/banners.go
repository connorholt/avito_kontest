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
	ID        int        `json:"banner_id"`
	TagID     []int      `json:"tag_ids"`
	FeatureID int        `json:"feature_id"`
	Content   Content    `json:"content"`
	IsActive  *bool      `json:"is_active"`
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
	stmt := `SELECT DISTINCT b.id,  feature_id, title, text, url, visible, created_at, updated_at from banners b
                              join features f on b.feature_id = f.id
                              join banner_tag bt on b.id = bt.banner_id
                              join tags t on t.id = bt.tag_id
								where case
								when $1 > 0 and $2 > 0 then tag_id = $1 and feature_id = $2
								when  $1 > 0 then tag_id = $1
								when $2 > 0 then feature_id = $2
								else true
								end 
								order by b.id
								limit case 
								when $3 > 0 then $3
								else NULL
								END
								offset case 
								when $4 >= 0 then $4
								else 0
								END`

	capacity := 10
	if limit >= 0 {
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
		err := rows.Scan(&b.ID, &b.FeatureID, &b.Content.Title, &b.Content.Text, &b.Content.URL, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}

		b.TagID, err = m.getIDs(b.ID)
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

func (m *BannerModel) Delete(id int) error {
	stmt := `DELETE FROM banners where id = $1`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoRecord
		} else {
			return err
		}
	}
	return nil
}

func (m *BannerModel) Create(b Banner) (int, error) {
	stmt := `INSERT INTO banners (title, text, url, visible, feature_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := m.DB.QueryRow(stmt, b.Content.Title, b.Content.Text, b.Content.URL, b.IsActive, b.FeatureID)
	var insertedID int
	err := row.Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (m *BannerModel) InsertBannerTag(bannerID int, tagIDs []int) error {
	stmt := `INSERT INTO banner_tag VALUES ($1, $2)`

	for i := range tagIDs {
		_, err := m.DB.Exec(stmt, bannerID, tagIDs[i])
		if err != nil {
			return err
		}
	}

	return nil
}

//func (m *BannerModel) Update(b Banner) error {
//	stmt := `UPDATE banners set `
//	var bldr strings.Builder
//	if b.Content.Title != "" {
//		bldr.WriteString("title = $1, ")
//	}
//	if b.Content.Text != "" {
//		bldr.WriteString("text = $2, ")
//	}
//	if b.Content.URL != "" {
//		bldr.WriteString("url = $3, ")
//	}
//	if b.FeatureID != 0 {
//		bldr.WriteString("feature_id =$4")
//	}
//	//if b.IsActive != nil {
//	//	bldr.WriteString("feature_id =$5")
//	//}
//	_, err := m.DB.Exec(stmt, b.Content.Title, b.Content.Text, b.Content.URL, b.FeatureID, b.IsActive)
//}

func (m *BannerModel) getIDs(bannerID int) ([]int, error) {
	stmt := `SELECT tag_id from banner_tag where banner_id = $1`
	rows, err := m.DB.Query(stmt, bannerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var tagIDs []int
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tagIDs, nil
}
