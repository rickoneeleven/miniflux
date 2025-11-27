// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "miniflux.app/v2/internal/storage"

import (
	"database/sql"
	"time"
)

// LastGlobalFeedCheck returns the most recent feed check timestamp across all enabled feeds.
func (s *Storage) LastGlobalFeedCheck() (time.Time, error) {
	var result sql.NullTime

	err := s.db.QueryRow(`SELECT max(checked_at) FROM feeds WHERE disabled IS FALSE`).Scan(&result)
	if err != nil {
		return time.Time{}, err
	}

	if !result.Valid {
		return time.Time{}, nil
	}

	return result.Time, nil
}

