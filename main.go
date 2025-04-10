package main

import (
	"fmt"

	"github.com/acme/sample/command"
	"github.com/acme/sample/event"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewChildcareClassActor() actor.Actor {
	return &ChildcareClassActor{
		state: new(event.ActivityAdded),
	}
}

type ChildcareClassActor struct {
	persistence.Mixin
	state *event.ActivityAdded
}

// Receive is the entry point for messages
func (a *ChildcareClassActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *persistence.RequestSnapshot:
		a.PersistSnapshot(a.state)
	case *persistence.ReplayComplete:
		fmt.Printf("Replay complete for Child %s\n", ctx.Self().Id)
	case *command.AddObservation:
		evt := &event.ObservationAdded{
			ChildID:  msg.ChildID,
			Note:     msg.Note,
			DateTime: timestamppb.New(msg.DateTime),
		}
		sender := ctx.Sender()
		a.applyEvent(evt)
		ctx.Send(msg.ReplyTo, evt)
		ctx.Send(sender, "Observation recorded.")
	case *command.AddNap:
		evt := &event.NapAdded{
			ChildID:         msg.ChildID,
			DurationMinutes: msg.DurationMinutes,
			StartTime:       timestamppb.New(msg.StartTime),
		}
		sender := ctx.Sender()
		a.applyEvent(evt)
		ctx.Send(msg.ReplyTo, evt)
		ctx.Send(sender, "Nap recorded.")
	case *command.AddMeal:
		evt := &event.MealAdded{
			ChildID:  msg.ChildID,
			MealType: msg.MealType,
			DateTime: timestamppb.New(msg.DateTime),
			Menu:     msg.Menu,
		}
		sender := ctx.Sender()
		a.applyEvent(evt)
		ctx.Send(msg.ReplyTo, evt)
		ctx.Send(sender, "Meal recorded.")
	case *command.PrintActivities:
		fmt.Println(a.state.Details)
	case *event.ObservationAdded, *event.NapAdded, *event.MealAdded:
		a.applyEvent(msg.(proto.Message))
	}
}

func (a *ChildcareClassActor) applyEvent(msg proto.Message) {
	if !a.Recovering() {
		a.PersistReceive(msg)
	}
	switch e := msg.(type) {
	case *event.ObservationAdded:
		// 実際にはドメインロジックなりいろんな処理が必要です
		r, _ := anypb.New(e)
		a.state.Details = append(a.state.Details, r)
	case *event.NapAdded:
		r, _ := anypb.New(e)
		a.state.Details = append(a.state.Details, r)
	case *event.MealAdded:
		r, _ := anypb.New(e)
		a.state.Details = append(a.state.Details, r)
	}
}
