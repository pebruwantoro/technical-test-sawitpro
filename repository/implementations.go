package repository

import (
	"context"
)

func (r *Repository) CreateEstate(ctx context.Context, input Estate) (result Estate, err error) {
	var id string
	err = r.Db.QueryRowContext(ctx, `
		INSERT INTO estates (id, width, length)
		VALUES ($1, $2, $3)
		returning id;
	`,
		input.Id,
		input.Width,
		input.Length,
	).Scan(&id)
	if err != nil {
		return
	}

	result.Id = id

	return
}

func (r *Repository) CreateEstateTree(ctx context.Context, input EstateTree) (result EstateTree, err error) {
	err = r.Db.QueryRowContext(ctx, `
		INSERT INTO trees (id, estate_id, x, y, height)
		VALUES ($1, $2, $3, $4, $5)
		returning id;
	`,
		input.Id,
		input.EstateId,
		input.X,
		input.Y,
		input.Height,
	).Scan(&result.Id)
	if err != nil {
		return
	}

	result = input

	return
}

func (r *Repository) GetStatsByEstateId(ctx context.Context, id string) (result StatsEstate, err error) {
	err = r.Db.QueryRowContext(ctx, `
	    SELECT 
			COALESCE(COUNT(*), 0) AS count, 
			COALESCE(MAX(height), 0) AS max_height, 
			COALESCE(MIN(height), 0) AS min_height, 
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height), 0) AS median_height
		FROM trees
		WHERE estate_id = $1;
	`, id).Scan(
		&result.Count,
		&result.Max,
		&result.Min,
		&result.Median,
	)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetEstateById(ctx context.Context, id string) (result Estate, err error) {
	err = r.Db.QueryRowContext(ctx, `
		SELECT id, width, length FROM estates WHERE id = $1;
	`, id).Scan(
		&result.Id,
		&result.Width,
		&result.Length,
	)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetTreesByEstateId(ctx context.Context, id string) (result []EstateTree, err error) {
	rows, err := r.Db.QueryContext(ctx, `
        SELECT id, estate_id, x, y, height FROM trees WHERE estate_id = $1;
    `, id)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tree EstateTree
		err = rows.Scan(
			&tree.Id,
			&tree.EstateId,
			&tree.X,
			&tree.Y,
			&tree.Height,
		)
		if err != nil {
			return
		}
		result = append(result, tree)
	}

	return
}
