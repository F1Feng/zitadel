package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/policy"
)

var (
	MailTextAddedEventType   = iamEventTypePrefix + policy.MailTextPolicyAddedEventType
	MailTextChangedEventType = iamEventTypePrefix + policy.MailTextPolicyChangedEventType
)

type MailTextAddedEvent struct {
	policy.MailTextAddedEvent
}

func NewMailTextAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	mailTextType,
	language,
	title,
	preHeader,
	subject,
	greeting,
	text,
	buttonText string,
) *MailTextAddedEvent {
	return &MailTextAddedEvent{
		MailTextAddedEvent: *policy.NewMailTextAddedEvent(
			eventstore.NewBaseEventForPush(ctx, aggregate, MailTextAddedEventType),
			mailTextType,
			language,
			title,
			preHeader,
			subject,
			greeting,
			text,
			buttonText),
	}
}

func MailTextAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.MailTextAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MailTextAddedEvent{MailTextAddedEvent: *e.(*policy.MailTextAddedEvent)}, nil
}

type MailTextChangedEvent struct {
	policy.MailTextChangedEvent
}

func NewMailTextChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	mailTextType,
	language string,
	changes []policy.MailTextChanges,
) (*MailTextChangedEvent, error) {
	changedEvent, err := policy.NewMailTextChangedEvent(
		eventstore.NewBaseEventForPush(ctx, aggregate, MailTextChangedEventType),
		mailTextType,
		language,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &MailTextChangedEvent{MailTextChangedEvent: *changedEvent}, nil
}

func MailTextChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.MailTextChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MailTextChangedEvent{MailTextChangedEvent: *e.(*policy.MailTextChangedEvent)}, nil
}
