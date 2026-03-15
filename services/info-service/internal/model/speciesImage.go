package model

type SpeciesImage struct {
	ID        int64  `db:"id"`
	SpeciesID int64  `db:"species_id"`
	ImagePath string `db:"image_path"`
	IsPrimary bool   `db:"is_primary"`
}
