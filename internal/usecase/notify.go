package usecase

import (
	"context"
	"sync"

	"github.com/mvp-mogila/ozon-test-task/internal/delivery/graphql/dto"
	"github.com/mvp-mogila/ozon-test-task/internal/models"
	"github.com/mvp-mogila/ozon-test-task/pkg/logger"
)

type CommentSubscriber struct {
	id int
	ch chan *dto.Comment
}

type CommentNotificationService struct {
	mu          *sync.RWMutex
	lastID      int
	subscribers map[int][]CommentSubscriber
}

func NewCommentNotificationService() *CommentNotificationService {
	return &CommentNotificationService{
		mu:          &sync.RWMutex{},
		subscribers: make(map[int][]CommentSubscriber),
	}
}

func (n *CommentNotificationService) CreateSubscriber(ctx context.Context, postID int) (ch chan *dto.Comment, newSubID int, err error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.lastID++
	ch = make(chan *dto.Comment)
	newSubID = n.lastID

	newSub := CommentSubscriber{
		id: newSubID,
		ch: ch,
	}

	subs := n.subscribers[postID]
	newSubs := append(subs, newSub)
	n.subscribers[postID] = newSubs

	logger.Infof(ctx, "Cretated new subscriber %d on post %d", newSubID, postID)
	err = nil
	return
}

func (n *CommentNotificationService) NotifySubscribers(ctx context.Context, comment models.Comment) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	subs := n.subscribers[comment.PostID]
	for _, sub := range subs {
		commentDTO := dto.ConvertCommentModelToDTO(comment)
		sub.ch <- &commentDTO
	}
	return nil
}

func (n *CommentNotificationService) DeleteSubscriber(postID int, subID int) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	subs := n.subscribers[postID]
	newSubs := make([]CommentSubscriber, len(subs)-1)
	for _, sub := range subs {
		if sub.id == subID {
			close(sub.ch)
		}
		newSubs = append(newSubs, sub)
	}

	n.subscribers[postID] = newSubs
	return nil
}
