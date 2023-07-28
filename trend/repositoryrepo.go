package trend

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/liweiyi88/gti/database"
)

type RankedTrendingRepository = map[int]TrendingRepository

type TrendingRepositoryRepo struct {
	db database.DB
}

func NewTrendingRepositoryRepo(db database.DB) *TrendingRepositoryRepo {
	return &TrendingRepositoryRepo{
		db: db,
	}
}

func (tr *TrendingRepositoryRepo) FindRankedTrendsByDate(ctx context.Context, date time.Time, language string) (RankedTrendingRepository, error) {
	lang := strings.TrimSpace(language)

	query := "SELECT * FROM trending_repositories WHERE trend_date = ? AND language is null"
	args := []any{date.Format("2006-01-02")}

	if lang != "" {
		query = "SELECT * FROM trending_repositories WHERE trend_date = ? AND language = ?"
		args = append(args, lang)
	}

	rows, err := tr.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rankedTrends := make(map[int]TrendingRepository, 0)

	for rows.Next() {
		var tr TrendingRepository

		if err := rows.Scan(&tr.Id, &tr.RepoFullName, &tr.Language, &tr.Rank, &tr.ScrapedAt, &tr.TrendDate, &tr.RepositoryId); err != nil {
			return rankedTrends, err
		}

		rankedTrends[tr.Rank] = tr
	}

	if err = rows.Err(); err != nil {
		return rankedTrends, err
	}

	return rankedTrends, nil
}

func (tr *TrendingRepositoryRepo) Save(ctx context.Context, trend TrendingRepository) error {
	query := "INSERT INTO `trending_repositories` (`full_name`, `language`, `rank`, `scraped_at`, `trend_date`) VALUES (?, ?, ?, ?, ?)"

	result, err := tr.db.ExecContext(ctx, query, trend.RepoFullName, trend.Language, trend.Rank, trend.ScrapedAt.Format("2006-01-02 15:04:05"), trend.TrendDate.Format("2006-01-02"))

	if err != nil {
		return fmt.Errorf("failed to exec insert query to db, error: %v", err)
	}

	_, err = result.RowsAffected()

	if err != nil {
		return fmt.Errorf("rows affected returns error: %v", err)
	}

	return nil
}

func (tr *TrendingRepositoryRepo) Update(ctx context.Context, trend TrendingRepository) error {
	query := "UPDATE `trending_repositories` SET full_name = ?, rank = ?, language = ?, scraped_at = ?, trend_date = ? WHERE id = ?"

	result, err := tr.db.ExecContext(ctx, query, trend.RepoFullName, trend.Rank, trend.Language, trend.ScrapedAt.Format("2006-01-02 15:04:05"), trend.TrendDate.Format("2006-01-02"), trend.Id)

	if err != nil {
		return fmt.Errorf("failed to run trending repo update query, trend id: %d, error: %v", trend.Id, err)
	}

	_, err = result.RowsAffected()

	if err != nil {
		return fmt.Errorf("rows affected returns error: %v", err)
	}

	return nil
}
