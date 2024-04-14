package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type BannerModel struct {
	DB *sql.DB
}

type Content struct {
	Title string `json:"title" example:"some_title"`
	Text  string `json:"text" example:"some_text"`
	URL   string `json:"url" example:"some_url"`
}

type Banner struct {
	ID        int        `json:"banner_id" `
	TagID     []int      `json:"tag_ids" example:"0"`
	FeatureID int        `json:"feature_id" example:"0"`
	Content   Content    `json:"content"`
	IsActive  *bool      `json:"is_active" example:"true"`
	CreatedAt *time.Time `json:"created_at,omitempty" example:"2024-04-14"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" example:"2024-04-14"`
}

func (m *BannerModel) Get(tagId int, featureId int) (*Banner, error) {
	stmt := `SELECT title, text, url, visible from banners b 
                               join features f on b.feature_id = f.id 
                               join banner_feature_tag bt on b.id = bt.banner_id
                               join tags t on t.id = bt.tag_id
								where b.feature_id = $1 and tag_id = $2`
	b := Banner{}
	row := m.DB.QueryRow(stmt, featureId, tagId)

	err := row.Scan(&b.Content.Title, &b.Content.Text, &b.Content.URL, &b.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return &b, nil
}

func (m *BannerModel) GetAll(tagID, featureID, limit, offset int) (*[]Banner, error) {
	stmt := `SELECT DISTINCT b.id,  b.feature_id, title, text, url, visible, created_at, updated_at from banners b
                              join features f on b.feature_id = f.id
                              join banner_feature_tag bt on b.id = bt.banner_id
                              join tags t on t.id = bt.tag_id
								where case
								when $1 > 0 and $2 > 0 then tag_id = $1 and b.feature_id = $2
								when  $1 > 0 then tag_id = $1
								when $2 > 0 then b.feature_id = $2
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
	row := m.DB.QueryRow(stmt, b.Content.Title, b.Content.Text, b.Content.URL, *b.IsActive, b.FeatureID)
	var insertedID int
	err := row.Scan(&insertedID)
	if err != nil {
		return 0, err
	}
	return insertedID, nil
}

func (m *BannerModel) Update(b Banner) error {
	var stmt strings.Builder
	stmt.WriteString("UPDATE banners set ")
	args := make([]any, 0, 10)
	index := 2
	args = append(args, b.ID)
	if b.Content.Title != "" {
		stmt.WriteString(fmt.Sprintf("title = $%d , ", index))
		args = append(args, b.Content.Title)
		index++
	}
	if b.Content.Text != "" {
		stmt.WriteString(fmt.Sprintf("text = $%d , ", index))
		args = append(args, b.Content.Text)
		index++
	}
	if b.Content.URL != "" {
		stmt.WriteString(fmt.Sprintf("url = $%d , ", index))
		args = append(args, b.Content.URL)
		index++
	}
	if b.FeatureID != 0 {
		stmt.WriteString(fmt.Sprintf("feature_id = $%d , ", index))
		args = append(args, b.FeatureID)
		index++
	}
	if b.IsActive != nil {
		stmt.WriteString(fmt.Sprintf("visible = $%d , ", index))
		args = append(args, *b.IsActive)
		index++
	}
	stmt.WriteString(" updated_at = CURRENT_TIMESTAMP where id = $1")
	res, err := m.DB.Exec(stmt.String(), args...)
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return ErrNoRecord
	}

	return nil
}

func (m *BannerModel) getIDs(bannerID int) ([]int, error) {
	stmt := `SELECT tag_id from banner_feature_tag where banner_id = $1`
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
