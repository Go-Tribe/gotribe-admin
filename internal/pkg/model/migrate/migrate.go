// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"fmt"
	"strconv"
	"strings"

	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// DBAutoMigrate 自动迁移表结构
func DBAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Admin{},
		&model.UserEvent{},
		&model.Role{},
		&model.Menu{},
		&model.Api{},
		&model.Post{},
		&model.PostTag{},
		&model.Tag{},
		&model.Category{},
		&model.Column{},
		&model.Comment{},
		&model.Project{},
		&model.Config{},
		&model.Resource{},
		&model.AdScene{},
		&model.Ad{},
		&model.PointLog{},
		&model.PointDeduction{},
		&model.PointAvailable{},
		&model.SystemConfig{},
		&model.OperationLog{},
		&model.ThirdPartyAccounts{},
		&model.Feedback{},
	)
	if err != nil {
		panic(fmt.Sprintf("database migration failed: %v", err))
	}

	normalizeLegacyUserContacts(db)
	ensureNoDuplicateAPIs(db)
	ensureNoDuplicateUsers(db)
	normalizeLegacyAPIDescriptions(db)
	recreateUniqueCompositeIndex(db, &model.Api{}, "idx_api_path_method", "path", "method")
	recreateUniqueCompositeIndex(db, &model.User{}, "idx_user_project_username", "project_id", "username")
	recreateUniqueCompositeIndex(db, &model.User{}, "idx_user_project_email", "project_id", "email")
	recreateUniqueCompositeIndex(db, &model.User{}, "idx_user_project_phone", "project_id", "phone")
	ensureCompositeIndex(db, &model.PointLog{}, "idx_point_log_project_user_created_at", "project_id", "user_id", "created_at")
	ensureCompositeIndex(db, &model.PointLog{}, "idx_point_log_project_created_at", "project_id", "created_at")
	ensureCompositeIndex(db, &model.PointLog{}, "idx_point_log_user_created_at", "user_id", "created_at")
	ensureCompositeIndex(db, &model.PointDeduction{}, "idx_point_deduction_project_user_created_at", "project_id", "user_id", "created_at")
	migrateLegacyPostTags(db)
}

func ensureCompositeIndex(db *gorm.DB, value interface{}, name string, columns ...string) {
	if db.Migrator().HasIndex(value, name) {
		return
	}
	if err := db.Migrator().CreateIndex(value, name); err == nil {
		return
	}

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(value); err != nil {
		panic(fmt.Sprintf("database schema parse failed for %s: %v", name, err))
	}
	tableName := stmt.Schema.Table
	if err := db.Exec(buildCreateIndexSQL(tableName, name, columns...)).Error; err != nil {
		panic(fmt.Sprintf("database index migration failed for %s: %v", name, err))
	}
}

func recreateUniqueCompositeIndex(db *gorm.DB, value interface{}, name string, columns ...string) {
	if db.Migrator().HasIndex(value, name) {
		if err := db.Migrator().DropIndex(value, name); err != nil {
			panic(fmt.Sprintf("drop index %s failed: %v", name, err))
		}
	}

	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(value); err != nil {
		panic(fmt.Sprintf("database schema parse failed for %s: %v", name, err))
	}
	if err := db.Exec(buildCreateUniqueIndexSQL(stmt.Quote(stmt.Schema.Table), quoteColumns(stmt, columns), name)).Error; err != nil {
		panic(fmt.Sprintf("database unique index migration failed for %s: %v", name, err))
	}
}

func buildCreateIndexSQL(tableName, indexName string, columns ...string) string {
	return fmt.Sprintf(
		"CREATE INDEX %s ON %s (%s)",
		indexName,
		tableName,
		strings.Join(columns, ", "),
	)
}

func buildCreateUniqueIndexSQL(tableName string, columns []string, indexName string) string {
	return fmt.Sprintf(
		"CREATE UNIQUE INDEX %s ON %s (%s)",
		indexName,
		tableName,
		strings.Join(columns, ", "),
	)
}

func migrateLegacyPostTags(db *gorm.DB) {
	postTable := tableName(db, &model.Post{})
	quotedPostTable := quotedTableName(db, &model.Post{})
	if !db.Migrator().HasColumn(postTable, "tag") {
		return
	}

	type legacyPostTag struct {
		ID  uint
		Tag string
	}

	var legacyPosts []legacyPostTag
	if err := db.Raw(fmt.Sprintf("SELECT id, tag FROM %s WHERE tag IS NOT NULL AND tag <> ''", quotedPostTable)).Scan(&legacyPosts).Error; err != nil {
		panic(fmt.Sprintf("legacy post tag migration failed: %v", err))
	}

	for _, legacyPost := range legacyPosts {
		tagIDs := parseLegacyTagIDs(legacyPost.Tag)
		for _, tagID := range tagIDs {
			postTag := model.PostTag{
				PostID: legacyPost.ID,
				TagID:  tagID,
			}
			if err := db.Where("post_id = ? AND tag_id = ?", postTag.PostID, postTag.TagID).FirstOrCreate(&postTag).Error; err != nil {
				panic(fmt.Sprintf("legacy post tag backfill failed for post %d: %v", legacyPost.ID, err))
			}
		}
	}
}

func parseLegacyTagIDs(raw string) []uint {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	values := strings.Split(raw, ",")
	seen := make(map[uint]struct{}, len(values))
	tagIDs := make([]uint, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		tagID, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			continue
		}
		id := uint(tagID)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		tagIDs = append(tagIDs, id)
	}
	return tagIDs
}

func ensureNoDuplicateAPIs(db *gorm.DB) {
	apiTable := quotedTableName(db, &model.Api{})
	type duplicateAPI struct {
		Path   string
		Method string
		Count  int64
	}

	var duplicate duplicateAPI
	err := db.Raw(fmt.Sprintf(`
		SELECT path, method, COUNT(*) AS count
		FROM %s
		GROUP BY path, method
		HAVING COUNT(*) > 1
		LIMIT 1
	`, apiTable)).Scan(&duplicate).Error
	if err != nil {
		panic(fmt.Sprintf("duplicate api check failed: %v", err))
	}
	if duplicate.Count > 1 {
		panic(fmt.Sprintf("duplicate api records found for path=%s method=%s", duplicate.Path, duplicate.Method))
	}
}

func ensureNoDuplicateUsers(db *gorm.DB) {
	userTable := quotedTableName(db, &model.User{})
	checks := []struct {
		name   string
		column string
		sql    string
	}{
		{
			name:   "username",
			column: "username",
			sql: fmt.Sprintf(`
				SELECT project_id, username AS value, COUNT(*) AS count
				FROM %s
				GROUP BY project_id, username
				HAVING COUNT(*) > 1
				LIMIT 1
			`, userTable),
		},
		{
			name:   "email",
			column: "email",
			sql: fmt.Sprintf(`
				SELECT project_id, email AS value, COUNT(*) AS count
				FROM %s
				WHERE email IS NOT NULL AND email <> ''
				GROUP BY project_id, email
				HAVING COUNT(*) > 1
				LIMIT 1
			`, userTable),
		},
		{
			name:   "phone",
			column: "phone",
			sql: fmt.Sprintf(`
				SELECT project_id, phone AS value, COUNT(*) AS count
				FROM %s
				WHERE phone IS NOT NULL AND phone <> ''
				GROUP BY project_id, phone
				HAVING COUNT(*) > 1
				LIMIT 1
			`, userTable),
		},
	}

	type duplicateUser struct {
		ProjectId uint
		Value     string
		Count     int64
	}

	for _, check := range checks {
		var duplicate duplicateUser
		if err := db.Raw(check.sql).Scan(&duplicate).Error; err != nil {
			panic(fmt.Sprintf("duplicate user %s check failed: %v", check.name, err))
		}
		if duplicate.Count > 1 {
			panic(fmt.Sprintf("duplicate user %s found for project_id=%d value=%s", check.name, duplicate.ProjectId, duplicate.Value))
		}
	}
}

func normalizeLegacyAPIDescriptions(db *gorm.DB) {
	apiTable := quotedTableName(db, &model.Api{})
	if err := db.Exec(
		fmt.Sprintf(`UPDATE %s SET "desc" = ? WHERE path = ? AND method = ? AND "desc" <> ?`, apiTable),
		"获取积分列表", "/point", "GET", "获取积分列表",
	).Error; err != nil {
		panic(fmt.Sprintf("normalize api descriptions failed: %v", err))
	}
}

func normalizeLegacyUserContacts(db *gorm.DB) {
	userTable := quotedTableName(db, &model.User{})
	statements := []string{
		fmt.Sprintf("UPDATE %s SET email = NULL WHERE email = ''", userTable),
		fmt.Sprintf("UPDATE %s SET phone = NULL WHERE phone = ''", userTable),
	}
	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			panic(fmt.Sprintf("legacy user contact normalization failed: %v", err))
		}
	}
}

func tableName(db *gorm.DB, value interface{}) string {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(value); err != nil {
		panic(fmt.Sprintf("database schema parse failed: %v", err))
	}
	return stmt.Schema.Table
}

func quotedTableName(db *gorm.DB, value interface{}) string {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(value); err != nil {
		panic(fmt.Sprintf("database schema parse failed: %v", err))
	}
	return stmt.Quote(stmt.Schema.Table)
}

func quoteColumns(stmt *gorm.Statement, columns []string) []string {
	quoted := make([]string, 0, len(columns))
	for _, column := range columns {
		quoted = append(quoted, stmt.Quote(column))
	}
	return quoted
}
