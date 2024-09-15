package pg

import (
	"context"
	"fmt"

	"dishdash.ru/internal/domain"
	"github.com/Vaniog/go-postgis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlaceRecommender struct {
	db *pgxpool.Pool
}

func NewPlaceRecommenderRepo(db *pgxpool.Pool) *PlaceRecommender {
	return &PlaceRecommender{db: db}
}

func (pr *PlaceRecommender) RecommendPlaces(
	ctx context.Context,
	opts domain.RecommendOpts,
	data domain.RecommendData,
) ([]*domain.Place, error) {
	query := `
	SELECT
		p.id,
		p.title,
		p.short_description,
		p.description,
		p.images,
		p.location,
		p.address,
		p.price_avg,
		p.review_rating,
		p.review_count,
		p.updated_at
	FROM place p
	JOIN place_tag pt ON p.id = pt.place_id
	JOIN tag t ON pt.tag_id = t.id

	WHERE t.name = ANY ($1)
	
	GROUP BY p.id

	ORDER BY
		$2 * ST_Distance(p.location, ST_GeogFromWkb($3)) +
		$4 * ABS(p.price_avg - $5)
`
	rows, err := pr.db.Query(ctx, query,
		data.Tags,
		opts.DistCoeff, postgis.PointS{SRID: 4326, X: data.Location.Lat, Y: data.Location.Lon},
		opts.PriceCoeff, data.PriceAvg,
	)
	if err != nil {
		return nil, fmt.Errorf("error while place recommending from db: %w", err)
	}
	defer rows.Close()

	var places []*domain.Place

	for rows.Next() {
		place, err := scanPlace(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return places, nil
}
