package unit_test

import (
	"testing"
	"time"

	"shifty-backend/graph/mapper"
	"shifty-backend/graph/model"
	"shifty-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMapCreateCommentToEntity_WithParentID(t *testing.T) {
	postID := uuid.New().String()
	parentID := uuid.New().String()
	input := &model.CreateCommentInput{
		PostID:   postID,
		Content:  "Reply comment",
		ParentID: &parentID,
	}

	comment, err := mapper.MapCreateCommentToEntity(input)

	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse(postID), comment.PostID)
	assert.NotNil(t, comment.ParentID)
	assert.Equal(t, uuid.MustParse(parentID), *comment.ParentID)
}

func TestMapCommentEntityToModel_WithReplies(t *testing.T) {
	now := time.Now()
	postID := uuid.New()
	authorID := uuid.New()
	parentID := uuid.New()
	replyID := uuid.New()
	comment := &entity.Comment{
		ID:        parentID,
		Content:   "Parent comment",
		PostID:    postID,
		AuthorID:  authorID,
		CreatedAt: now,
		UpdatedAt: now,
		Replies: []entity.Comment{
			{ID: replyID, Content: "Reply comment", PostID: postID, AuthorID: authorID, ParentID: &parentID, CreatedAt: now, UpdatedAt: now},
		},
	}

	mappedComment, err := mapper.MapCommentEntityToModel(comment)

	assert.NoError(t, err)
	assert.Equal(t, parentID.String(), mappedComment.ID)
	assert.Len(t, mappedComment.Replies, 1)
	assert.Equal(t, replyID.String(), mappedComment.Replies[0].ID)
	assert.NotNil(t, mappedComment.Replies[0].ParentID)
	assert.Equal(t, parentID.String(), *mappedComment.Replies[0].ParentID)
}
