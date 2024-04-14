package models

func (b Banner) Validate() error {
	for _, id := range b.TagID {
		if id < 1 {
			return ErrInvalidData
		}
	}
	if b.Content.Title == "" || b.Content.Text == "" || b.Content.URL == "" || b.FeatureID < 1 {
		return ErrInvalidData
	}
	return nil
}
