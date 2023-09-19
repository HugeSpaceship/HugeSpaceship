package db

import (
	"HugeSpaceship/pkg/common/model/lbp_xml/photos"
	"context"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

const photoInsertSQL = `INSERT INTO photos (domain, author, small, medium, large, plan, slotType, slotField) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id;`

func InsertPhoto(ctx context.Context, photo photos.UploadPhoto, author uuid.UUID, domain uint) (id uint64, err error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	tx, err := conn.Begin(ctx)
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

	row := tx.QueryRow(ctx, photoInsertSQL, domain, author, photo.Small, photo.Medium, photo.Large, photo.Plan, photo.Slot.Type, slotField)
	err = row.Scan(&id)
	if err != nil {
		er2 := tx.Rollback(ctx)
		err = errors.Join(err, er2)
		return
	}

	for _, subject := range photo.Subjects.Subjects {
		userID := uuid.NullUUID{Valid: false}

		if subject.NpHandle.Username != "" { // Get the userid if it is a valid user, i.e. not a local one
			user, err := GetUserByName(ctx, subject.NpHandle.Username)
			if err == nil {
				userID.Valid = true
				userID.UUID = user.ID
			}
		}

		x1, y1, x2, y2, err := photos.PhotoSubjectBounds(subject.Bounds).GetBounds()
		if err != nil {
			continue
		}

		_, err = tx.Exec(ctx, "INSERT INTO photo_subjects VALUES ($1,$2,$3,$4,$5,$6,$7);",
			id, userID, subject.DisplayName, x1, y1, x2, y2,
		)
		if err != nil {
			return 0, err
		}
	}

	err = tx.Commit(ctx)

	return
}

func GetPhotos(ctx context.Context, by uuid.UUID, pageSize, pageStart uint64, domain uint) (photos photos.Photos, err error) {
	conn := ctx.Value("conn").(*pgxpool.Conn)

	err = pgxscan.Select(ctx, conn, &photos.Photos, "SELECT * FROM photos WHERE author = $1 AND domain = $2 LIMIT $3 OFFSET $4",
		by, domain, pageSize, pageStart,
	)
	return
}
