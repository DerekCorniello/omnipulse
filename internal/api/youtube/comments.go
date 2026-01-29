// Package youtube provides comment fetching for YouTube videos.
package youtube

import (
	"context"

	"github.com/crossforge/crossforge/internal/data"
)

// Comments handles fetching YouTube comments.
type Comments struct {
	client *Client
}

// NewComments creates a new Comments instance.
func NewComments(client *Client) *Comments {
	return &Comments{client: client}
}

// GetVideoComments fetches comments for a specific video.
// Uses YouTube Data API: commentThreads.list
func (c *Comments) GetVideoComments(ctx context.Context, videoID string, maxResults int) ([]*data.Comment, error) {
	// TODO: Implement using YouTube Data API v3
	// Endpoint: GET https://www.googleapis.com/youtube/v3/commentThreads
	// Parameters:
	//   - part: snippet,replies
	//   - videoId: {video_id}
	//   - maxResults: 20-100
	//   - order: relevance or time
	//   - key: {api_key}
	//
	// Response includes:
	//   - snippet.topLevelComment.snippet.authorDisplayName
	//   - snippet.topLevelComment.snippet.textDisplay
	//   - snippet.topLevelComment.snippet.likeCount
	//   - snippet.topLevelComment.snippet.publishedAt
	//   - snippet.totalReplyCount
	return nil, nil
}

// GetChannelComments fetches recent comments across all channel videos.
// Uses YouTube Data API: commentThreads.list with allThreadsRelatedToChannelId
func (c *Comments) GetChannelComments(ctx context.Context, maxResults int) ([]*data.Comment, error) {
	// TODO: Implement using YouTube Data API v3
	// Endpoint: GET https://www.googleapis.com/youtube/v3/commentThreads
	// Parameters:
	//   - part: snippet
	//   - allThreadsRelatedToChannelId: {channel_id}
	//   - maxResults: 20-100
	//   - order: time
	//   - key: {api_key}
	//
	// Note: Requires the channel to have comments enabled
	// and may require OAuth for channels you own
	return nil, nil
}

// GetCommentReplies fetches replies to a specific comment thread.
// Uses YouTube Data API: comments.list
func (c *Comments) GetCommentReplies(ctx context.Context, parentID string, maxResults int) ([]*data.Comment, error) {
	// TODO: Implement using YouTube Data API v3
	// Endpoint: GET https://www.googleapis.com/youtube/v3/comments
	// Parameters:
	//   - part: snippet
	//   - parentId: {comment_id}
	//   - maxResults: 20-100
	//   - key: {api_key}
	return nil, nil
}
