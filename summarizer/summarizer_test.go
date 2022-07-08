package summarizer

import (
	"errors"
	"reflect"
	"testing"
)

type fakeProvider struct {
	entries []HistoryEntry
	doError bool
}

func (f fakeProvider) GetEntries() ([]HistoryEntry, error) {
	if f.doError {
		return nil, errors.New("client error")
	}
	return f.entries, nil
}

func Test_summarizer_GetSummary(t *testing.T) {
	type fields struct {
		client listEntryProvider
	}
	tests := []struct {
		name    string
		fields  fields
		want    Summary
		wantErr bool
	}{
		{name: "should return summary", fields: fields{client: fakeProvider{entries: []HistoryEntry{
			{Name: "consumerName", Candy: "apple", Eaten: 1},
		}}}, want: Summary{
			SummaryEntry{Name: "consumerName", FavouriteSnack: "apple", TotalSnacks: 1},
		}, wantErr: false},
		{name: "should return total of all candy consumed per consumer", fields: fields{client: fakeProvider{entries: []HistoryEntry{
			{Name: "consumerName", Candy: "apple", Eaten: 100},
			{Name: "consumerName", Candy: "apple", Eaten: 200},
			{Name: "consumerName", Candy: "orange", Eaten: 30},
		}}}, want: Summary{
			SummaryEntry{Name: "consumerName", FavouriteSnack: "apple", TotalSnacks: 330},
		}, wantErr: false},
		{name: "should properly identify most consumed candy per consumer", fields: fields{client: fakeProvider{entries: []HistoryEntry{
			{Name: "consumerName", Candy: "apple", Eaten: 100},
			{Name: "consumerName", Candy: "apple", Eaten: 200},
			{Name: "consumerName", Candy: "orange", Eaten: 310},
		}}}, want: Summary{
			SummaryEntry{Name: "consumerName", FavouriteSnack: "orange", TotalSnacks: 610},
		}, wantErr: false},
		{name: "should reduce the consumers", fields: fields{client: fakeProvider{entries: []HistoryEntry{
			{Name: "consumerName", Candy: "apple", Eaten: 100},
			{Name: "consumerName", Candy: "apple", Eaten: 200},
			{Name: "consumer2", Candy: "apple", Eaten: 20},
		}}}, want: Summary{
			SummaryEntry{Name: "consumerName", FavouriteSnack: "apple", TotalSnacks: 300}, SummaryEntry{Name: "consumer2", FavouriteSnack: "apple", TotalSnacks: 20},
		}, wantErr: false},
		{name: "should return summary sorted by largest consumer", fields: fields{client: fakeProvider{entries: []HistoryEntry{
			{Name: "consumer2", Candy: "apple", Eaten: 20},
			{Name: "consumerName", Candy: "apple", Eaten: 100},
			{Name: "consumerName", Candy: "apple", Eaten: 200},
		}}}, want: Summary{
			SummaryEntry{Name: "consumerName", FavouriteSnack: "apple", TotalSnacks: 300}, SummaryEntry{Name: "consumer2", FavouriteSnack: "apple", TotalSnacks: 20},
		}, wantErr: false},
		{name: "should return error when client fails", fields: fields{client: fakeProvider{doError: true}}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				client: tt.fields.client,
			}
			got, err := s.GetSummary()
			if (err != nil) != tt.wantErr {
				t.Errorf("summarizer.GetSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("summarizer.GetSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}
