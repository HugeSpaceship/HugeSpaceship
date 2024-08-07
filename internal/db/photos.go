package db

import (
	"context"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/common"
	"github.com/HugeSpaceship/HugeSpaceship/internal/model/lbp_xml/photos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

const photoInsertSQL = `INSERT INTO photos (domain, author, small, medium, large, plan, slotType, slotField) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id;`

func InsertPhoto(conn *pgxpool.Conn, photo *photos.UploadPhoto, author uuid.UUID, domain uint) (id uint64, err error) {

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return
	}
	var slotField string

	switch photo.Slot.Type {
	case "user", "developer":
		slotField = strconv.FormatInt(photo.Slot.ID, 10)
	case "pod", "moon", "local":
		slotField = photo.Slot.RootLevel
	}

	row := tx.QueryRow(context.Background(), photoInsertSQL, domain, author, photo.Small, photo.Medium, photo.Large, photo.Plan, photo.Slot.Type, slotField)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback(context.Background())
		return
	}

	for _, subject := range photo.Subjects.Subjects {
		userID := uuid.NullUUID{Valid: false}

		if subject.NpHandle.Username != "" { // Get the userid if it is a valid user, i.e. not a local one
			// the game is not currently required for photos
			user, err := GetUserByName(conn, subject.NpHandle.Username, common.LBP2)
			if err == nil {
				userID.Valid = true
				userID.UUID = user.ID
			}
		}

		x1, y1, x2, y2, err := photos.PhotoSubjectBounds(subject.Bounds).GetBounds()
		if err != nil {
			continue
		}

		_, err = tx.Exec(context.Background(), "INSERT INTO photo_subjects VALUES ($1,$2,$3,$4,$5,$6,$7);",
			id, userID, subject.DisplayName, x1, y1, x2, y2,
		)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit(context.Background())

	return
}

func GetPhotos(conn *pgxpool.Conn, by uuid.UUID, pageSize, pageStart uint64, domain uint) (outPhotos photos.Photos, err error) {

	const photosSQL = "SELECT * FROM photos WHERE author = $1 AND domain = $2 LIMIT $3 OFFSET $4"

	rows, err := conn.Query(context.Background(), photosSQL, by, domain, pageSize, pageStart)
	if err != nil {
		return
	}
	outPhotos.Photos, err = pgx.CollectRows(rows, pgx.RowToStructByName[photos.Photo])

	return
}
