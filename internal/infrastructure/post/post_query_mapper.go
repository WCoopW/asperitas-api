package post

import (
	"fmt"
	"slices"
	"strings"

	domain "reddit/internal/domain/post"
)

type (
	PostQueryBuilder struct{}
	joinSpec         struct {
		selectCols string
		joins      string
	}
)

const baseQuery = `SELECT posts.id, posts.author_id, posts.content_type, posts.title, posts.url_content, posts.text_content, posts.score, posts.views, posts.upvote_percentage, posts.created_at`

func (qb *PostQueryBuilder) BuildQuery(request *domain.PostRequest) (string, []any) {
	query := strings.Builder{}
	query.WriteString(baseQuery)
	if request == nil {
		query.WriteString(" FROM posts")
		return query.String(), nil
	}
	spec := qb.populatedFields(&request.Populated)
	if spec.selectCols != "" {
		query.WriteString(", " + spec.selectCols)
	}
	query.WriteString(" FROM posts")
	if spec.joins != "" {
		query.WriteString(" " + spec.joins)
	}
	whereClause, args := qb.filterPosts(request.PostFilter)
	query.WriteString(whereClause)
	return query.String(), args
}

func (qb *PostQueryBuilder) filterPosts(filter domain.PostFilter) (string, []any) {
	var conditions []string
	var args []any

	if filter.ID != "" {
		args = append(args, filter.ID)
		conditions = append(conditions, fmt.Sprintf("posts.id = $%d", len(args)))
	}
	if filter.AuthorID != "" {
		args = append(args, filter.AuthorID)
		conditions = append(conditions, fmt.Sprintf("posts.author_id = $%d", len(args)))
	}
	if filter.Category != "" {
		args = append(args, "%"+filter.Category+"%")
		conditions = append(conditions, fmt.Sprintf("category ILIKE $%d", len(args)))
	}
	if len(conditions) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (qb *PostQueryBuilder) populatedFields(populate *[]domain.PostPopulatedFields) joinSpec {
	if populate == nil || len(*populate) == 0 {
		return joinSpec{}
	}
	var selectCols []string
	var joins []string
	if slices.Contains(*populate, domain.AuthorPopulate) {
		joins = append(joins, "LEFT JOIN users ON posts.author_id = users.id")
		selectCols = append(selectCols, "users.id AS author_id", "users.username")

	}
	return joinSpec{
		selectCols: strings.Join(selectCols, ", "),
		joins:      strings.Join(joins, " "),
	}
}
