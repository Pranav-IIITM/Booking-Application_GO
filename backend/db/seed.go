package db

import (
	"context"

	"booking-backend/models"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func SeedSlots(ctx context.Context, firestoreClient *firestore.Client) error {
	slotsCollection := firestoreClient.Collection("slots")
	iter := slotsCollection.Limit(1).Documents(ctx)
	_, err := iter.Next()
	iter.Stop()
	if err == nil {
		return nil
	}
	if err != iterator.Done {
		return err
	}

	slots := map[string]models.Slot{
		"1": {Date: "2026-06-20", Time: "10:00 AM", Capacity: 10, BookedCount: 0},
		"2": {Date: "2026-06-20", Time: "12:00 PM", Capacity: 10, BookedCount: 0},
		"3": {Date: "2026-06-21", Time: "10:00 AM", Capacity: 5, BookedCount: 0},
		"4": {Date: "2026-06-21", Time: "02:00 PM", Capacity: 8, BookedCount: 0},
		"5": {Date: "2026-06-22", Time: "11:00 AM", Capacity: 6, BookedCount: 0},
	}

	for id, slot := range slots {
		if _, err := slotsCollection.Doc(id).Set(ctx, slot); err != nil {
			return err
		}
	}

	return nil
}
