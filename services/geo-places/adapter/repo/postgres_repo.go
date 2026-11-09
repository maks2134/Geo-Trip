package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/makskozlov/atlasflow/services/geo-places/domain"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
)

type PostgresPlaceRepository struct {
	db *sql.DB
}

func NewPostgresPlaceRepository(db *sql.DB) *PostgresPlaceRepository {
	return &PostgresPlaceRepository{db: db}
}

func (r *PostgresPlaceRepository) UpsertPlace(ctx context.Context, p domain.Place) error {
	query := `
		INSERT INTO places (
			id, owner_id, name, category, status, country_iso2, city, address, geom, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			ST_SetSRID(ST_MakePoint($9, $10), 4326)::geography,
			$11, $12
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			category = EXCLUDED.category,
			status = EXCLUDED.status,
			country_iso2 = EXCLUDED.country_iso2,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			geom = EXCLUDED.geom,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := r.db.ExecContext(
		ctx, query,
		p.ID, p.OwnerID, p.Name, string(p.Category), string(p.Status),
		p.CountryISO2, p.City, p.Address,
		p.Location.Longitude, p.Location.Latitude,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert place sql: %w", err)
	}
	return nil
}

func (r *PostgresPlaceRepository) GetPlace(ctx context.Context, id, ownerID string) (*domain.Place, error) {
	query := `
		SELECT id, owner_id, name, category, status, country_iso2, city, address,
		       ST_Y(geom::geometry) AS lat, ST_X(geom::geometry) AS lng,
		       created_at, updated_at
		FROM places
		WHERE id = $1 AND ($2 = '' OR owner_id = $2);
	`
	var p domain.Place
	var category, status string
	var countryISO2, city, address sql.NullString

	err := r.db.QueryRowContext(ctx, query, id, ownerID).Scan(
		&p.ID, &p.OwnerID, &p.Name, &category, &status,
		&countryISO2, &city, &address,
		&p.Location.Latitude, &p.Location.Longitude,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get place sql: %w", err)
	}

	p.Category = domain.PlaceCategory(category)
	p.Status = domain.PlaceStatus(status)
	p.CountryISO2 = countryISO2.String
	p.City = city.String
	p.Address = address.String

	return &p, nil
}

func (r *PostgresPlaceRepository) BuildSearchSQL(filter usecase.SearchFilter) (string, []any) {
	var sb strings.Builder
	var args []any
	argIdx := 1

	sb.WriteString(`
		SELECT id, owner_id, name, category, status, country_iso2, city, address,
		       ST_Y(geom::geometry) AS lat, ST_X(geom::geometry) AS lng,
		       created_at, updated_at
		FROM places
		WHERE 1=1
	`)

	if filter.OwnerID != "" {
		sb.WriteString(fmt.Sprintf(" AND owner_id = $%d", argIdx))
		args = append(args, filter.OwnerID)
		argIdx++
	}

	if filter.Status != nil && *filter.Status != domain.StatusUnspecified {
		sb.WriteString(fmt.Sprintf(" AND status = $%d", argIdx))
		args = append(args, string(*filter.Status))
		argIdx++
	}

	if filter.Category != nil && *filter.Category != domain.CategoryUnspecified {
		sb.WriteString(fmt.Sprintf(" AND category = $%d", argIdx))
		args = append(args, string(*filter.Category))
		argIdx++
	}

	if filter.Query != "" {
		sb.WriteString(fmt.Sprintf(" AND (name ILIKE $%d OR city ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Query+"%")
		argIdx++
	}

	if filter.BBox != nil {
		sb.WriteString(fmt.Sprintf(`
			AND ST_Intersects(
				geom,
				ST_MakeEnvelope($%d, $%d, $%d, $%d, 4326)::geography
			)
		`, argIdx, argIdx+1, argIdx+2, argIdx+3))
		args = append(args, filter.BBox.MinLongitude, filter.BBox.MinLatitude, filter.BBox.MaxLongitude, filter.BBox.MaxLatitude)
		argIdx += 4
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = usecase.DefaultLimit
	} else if limit > usecase.MaxLimit {
		limit = usecase.MaxLimit
	}

	sb.WriteString(fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argIdx))
	args = append(args, limit)

	return sb.String(), args
}

func (r *PostgresPlaceRepository) SearchPlaces(ctx context.Context, filter usecase.SearchFilter) ([]domain.Place, error) {
	query, args := r.BuildSearchSQL(filter)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("search places sql: %w", err)
	}
	defer rows.Close()

	var places []domain.Place
	for rows.Next() {
		var p domain.Place
		var category, status string
		var countryISO2, city, address sql.NullString

		if err := rows.Scan(
			&p.ID, &p.OwnerID, &p.Name, &category, &status,
			&countryISO2, &city, &address,
			&p.Location.Latitude, &p.Location.Longitude,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan place sql: %w", err)
		}

		p.Category = domain.PlaceCategory(category)
		p.Status = domain.PlaceStatus(status)
		p.CountryISO2 = countryISO2.String
		p.City = city.String
		p.Address = address.String

		places = append(places, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return places, nil
}

func (r *PostgresPlaceRepository) MarkCountryVisited(ctx context.Context, ownerID, countryCode string, visited bool) error {
	if visited {
		query := `
			INSERT INTO visited_countries (user_id, country_iso2, visited_at)
			VALUES ($1, $2, NOW())
			ON CONFLICT (user_id, country_iso2) DO NOTHING;
		`
		_, err := r.db.ExecContext(ctx, query, ownerID, countryCode)
		if err != nil {
			return fmt.Errorf("mark country sql: %w", err)
		}
	} else {
		query := `DELETE FROM visited_countries WHERE user_id = $1 AND country_iso2 = $2;`
		_, err := r.db.ExecContext(ctx, query, ownerID, countryCode)
		if err != nil {
			return fmt.Errorf("unmark country sql: %w", err)
		}
	}
	return nil
}

func (r *PostgresPlaceRepository) GetVisitedCountries(ctx context.Context, ownerID string) ([]string, error) {
	query := `
		SELECT DISTINCT country_iso2
		FROM (
			SELECT country_iso2 FROM visited_countries WHERE user_id = $1
			UNION
			SELECT country_iso2 FROM places
			WHERE owner_id = $1 AND status = 'visited' AND country_iso2 IS NOT NULL AND country_iso2 != ''
		) combined
		ORDER BY country_iso2;
	`
	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("get visited countries sql: %w", err)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("scan country sql: %w", err)
		}
		codes = append(codes, strings.ToUpper(code))
	}

	return codes, rows.Err()
}

func (r *PostgresPlaceRepository) GetTravelerStats(ctx context.Context, ownerID string) (*domain.TravelerStats, error) {
	visitedCodes, err := r.GetVisitedCountries(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	var citiesVisited int32
	citiesQuery := `
		SELECT COUNT(DISTINCT city)
		FROM places
		WHERE owner_id = $1 AND status = 'visited' AND city IS NOT NULL AND city != '';
	`
	if err := r.db.QueryRowContext(ctx, citiesQuery, ownerID).Scan(&citiesVisited); err != nil {
		return nil, fmt.Errorf("stats cities sql: %w", err)
	}

	var pointsTotal int32
	pointsQuery := `SELECT COUNT(*) FROM places WHERE owner_id = $1;`
	if err := r.db.QueryRowContext(ctx, pointsQuery, ownerID).Scan(&pointsTotal); err != nil {
		return nil, fmt.Errorf("stats points sql: %w", err)
	}

	stats := domain.NewTravelerStats(visitedCodes, citiesVisited, pointsTotal, 0.0)
	return &stats, nil
}

func (r *PostgresPlaceRepository) GetVisitedCountryPolygons(ctx context.Context, countryCodes []string) ([]domain.CountryFeature, error) {
	if len(countryCodes) == 0 {
		return []domain.CountryFeature{}, nil
	}

	placeholders := make([]string, len(countryCodes))
	args := make([]any, len(countryCodes))
	for i, code := range countryCodes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = strings.ToUpper(code)
	}

	query := fmt.Sprintf(`
		SELECT iso2, ST_AsGeoJSON(geom) AS geojson
		FROM countries
		WHERE iso2 IN (%s)
		ORDER BY iso2;
	`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get country polygons sql: %w", err)
	}
	defer rows.Close()

	var features []domain.CountryFeature
	for rows.Next() {
		var feat domain.CountryFeature
		if err := rows.Scan(&feat.ISO2, &feat.GeoJSON); err != nil {
			return nil, fmt.Errorf("scan polygon sql: %w", err)
		}
		features = append(features, feat)
	}

	return features, rows.Err()
}