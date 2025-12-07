// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "miniflux.app/v2/internal/ui"

import (
	"net/http"
	"time"

	"miniflux.app/v2/internal/http/request"
	"miniflux.app/v2/internal/http/response/html"
	jsonResponse "miniflux.app/v2/internal/http/response/json"
	"miniflux.app/v2/internal/http/route"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/ui/session"
	"miniflux.app/v2/internal/ui/view"
)

type unreadSnapshotResponse struct {
	UnreadCount         int           `json:"unread_count"`
	LastGlobalFeedCheck time.Time     `json:"last_global_feed_check"`
	Entries             model.Entries `json:"entries"`
	CountErrorFeeds     int           `json:"count_error_feeds"`
}

func (h *handler) showUnreadPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	lastGlobalFeedCheck, err := h.store.LastGlobalFeedCheck()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	offset := request.QueryIntParam(r, "offset", 0)
	builder := h.store.NewEntryQueryBuilder(user.ID)
	builder.WithStatus(model.EntryStatusUnread)
	builder.WithGloballyVisible()
	countUnread, err := builder.CountEntries()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	if offset >= countUnread {
		offset = 0
	}

	builder = h.store.NewEntryQueryBuilder(user.ID)
	builder.WithStatus(model.EntryStatusUnread)
	builder.WithSorting(user.EntryOrder, user.EntryDirection)
	builder.WithSorting("id", user.EntryDirection)
	builder.WithOffset(offset)
	builder.WithLimit(user.EntriesPerPage)
	builder.WithGloballyVisible()
	entries, err := builder.GetEntries()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	sess := session.New(h.store, request.SessionID(r))
	view := view.New(h.tpl, r, sess)
	view.Set("entries", entries)
	view.Set("pagination", getPagination(route.Path(h.router, "unread"), countUnread, offset, user.EntriesPerPage))
	view.Set("hasLastGlobalFeedCheck", !lastGlobalFeedCheck.IsZero())
	view.Set("lastGlobalFeedCheck", lastGlobalFeedCheck)
	view.Set("menu", "unread")
	view.Set("user", user)
	view.Set("countUnread", countUnread)
	view.Set("countErrorFeeds", h.store.CountUserFeedsWithErrors(user.ID))
	view.Set("hasSaveEntry", h.store.HasSaveEntry(user.ID))

	html.OK(w, r, view.Render("unread_entries"))
}

func (h *handler) unreadSnapshot(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		jsonResponse.ServerError(w, r, err)
		return
	}

	lastGlobalFeedCheck, err := h.store.LastGlobalFeedCheck()
	if err != nil {
		jsonResponse.ServerError(w, r, err)
		return
	}

	countBuilder := h.store.NewEntryQueryBuilder(user.ID)
	countBuilder.WithStatus(model.EntryStatusUnread)
	countBuilder.WithGloballyVisible()
	countUnread, err := countBuilder.CountEntries()
	if err != nil {
		jsonResponse.ServerError(w, r, err)
		return
	}

	builder := h.store.NewEntryQueryBuilder(user.ID)
	builder.WithStatus(model.EntryStatusUnread)
	builder.WithSorting(user.EntryOrder, user.EntryDirection)
	builder.WithSorting("id", user.EntryDirection)
	builder.WithOffset(0)
	builder.WithLimit(user.EntriesPerPage)
	builder.WithGloballyVisible()
	entries, err := builder.GetEntries()
	if err != nil {
		jsonResponse.ServerError(w, r, err)
		return
	}

	jsonResponse.OK(w, r, unreadSnapshotResponse{
		UnreadCount:         countUnread,
		LastGlobalFeedCheck: lastGlobalFeedCheck,
		Entries:             entries,
		CountErrorFeeds:     h.store.CountUserFeedsWithErrors(user.ID),
	})
}
