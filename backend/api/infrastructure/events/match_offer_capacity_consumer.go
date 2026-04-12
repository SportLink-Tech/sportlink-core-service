package events

import (
	"context"
	matchofferevent "sportlink/api/application/matchoffer/events"
	"sportlink/api/application/matchoffer/usecases"
	"sportlink/pkg/log"
)

// MatchOfferCapacityConsumer listens for MatchOfferCapacityReachedEvent and
// triggers the confirm match offer use case to auto-create the match.
type MatchOfferCapacityConsumer struct {
	ch        <-chan matchofferevent.MatchOfferCapacityReachedEvent
	confirmUC *usecases.ConfirmMatchOfferUC
}

func NewMatchOfferCapacityConsumer(
	ch <-chan matchofferevent.MatchOfferCapacityReachedEvent,
	confirmUC *usecases.ConfirmMatchOfferUC,
) *MatchOfferCapacityConsumer {
	return &MatchOfferCapacityConsumer{ch: ch, confirmUC: confirmUC}
}

// Start launches the consumer goroutine. It stops when ctx is cancelled or
// the channel is closed.
func (c *MatchOfferCapacityConsumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case event, ok := <-c.ch:
				if !ok {
					return
				}
				c.handleEvent(ctx, event)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *MatchOfferCapacityConsumer) handleEvent(ctx context.Context, event matchofferevent.MatchOfferCapacityReachedEvent) {
	_, err := c.confirmUC.Invoke(ctx, usecases.ConfirmMatchOfferInput{
		MatchOfferID:   event.MatchOfferID,
		OwnerAccountID: event.OwnerAccountID,
	})
	if err != nil {
		log.GetLogger(ctx).Error("auto-confirm failed for match offer "+event.MatchOfferID, err)
	}
}
