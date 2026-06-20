package models

type Slot struct {
	ID          string `firestore:"-" json:"id"`
	Date        string `firestore:"date" json:"date"`
	Time        string `firestore:"time" json:"time"`
	Capacity    int    `firestore:"capacity" json:"capacity"`
	BookedCount int    `firestore:"bookedCount" json:"bookedCount"`
}
