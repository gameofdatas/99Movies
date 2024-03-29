package io

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"99Movies/models"
)

func TestJSONReader_GetMovies(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *json.Decoder
		output  func() map[string]int
		wantErr bool
	}{
		{
			name: "successfully parse",
			setup: func() *json.Decoder {
				input := []byte(`[
									{"title":"Star Wars","year":1977},
									{"title":"Star Wars The Force Awakens","year":2015}
								 ]`)
				decoder := json.NewDecoder(bytes.NewReader(input))
				return decoder
			},
			output: func() map[string]int {
				m := make(map[string]int)
				m["Star Wars"] = 1977
				m["Star Wars The Force Awakens"] = 2015
				return m
			},
			wantErr: false,
		},
		{
			name: "error while parsing",
			setup: func() *json.Decoder {
				input := []byte(`[
									{"abc"},
									{"def"}
								 ]`)
				decoder := json.NewDecoder(bytes.NewReader(input))
				return decoder
			},
			output:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		decoder := tt.setup()
		reader := JSONReader{
			Parser: decoder,
			File:   nil,
		}
		items, err := reader.GetMovies()
		if err != nil && tt.wantErr == false {
			t.Errorf("unexpected error %s", err)
			return
		}
		if tt.wantErr == true && err == nil {
			t.Errorf("wanted an error during %s", tt.name)
			return
		}
		if tt.wantErr == true && err != nil {
			// :)
			return
		}
		if !reflect.DeepEqual(tt.output(), items) {
			t.Errorf("GetMovies(%v)=%v, wanted %v", tt.output(), items, tt.output())
		}
	}
}

func TestJSONReader_GetReviews(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *json.Decoder
		checkReview func(chan models.Reviews)
		checkError  func(chan error)
		wantErr     bool
	}{
		{
			name: "successfully parse",
			setup: func() *json.Decoder {
				input := []byte(`[
									{"title":"Star Wars","review":"Great, this film was","score":77}
								 ]`)
				decoder := json.NewDecoder(bytes.NewReader(input))
				return decoder
			},
			checkReview: func(out chan models.Reviews) {
				for review := range out {
					expectedReview := models.Reviews{
						Title:  "Star Wars",
						Review: "Great, this film was",
						Score:  77,
					}
					if !reflect.DeepEqual(expectedReview, review) {
						t.Errorf("GetReview()=%v, wanted %v", review, expectedReview)

					}
				}
			},
			wantErr: false,
		},
		{
			name: "error while parsing",
			setup: func() *json.Decoder {
				input := []byte(`[
									{"abc"},
									{"def"}
								 ]`)
				decoder := json.NewDecoder(bytes.NewReader(input))
				return decoder
			},
			checkError: func(errors chan error) {
				for err := range errors {
					expectedErrorString := "could not decode reviews invalid character '}' after object key"
					if !reflect.DeepEqual(err.Error(), expectedErrorString) {
						t.Errorf("GetReview()=%v, wanted %v", err, expectedErrorString)
					}
					return
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		decoder := tt.setup()
		reader := JSONReader{
			Parser: decoder,
			File:   nil,
		}
		out := make(chan models.Reviews)
		errors := make(chan error)
		go reader.GetReviews(out, errors)
		if tt.wantErr {
			tt.checkError(errors)
		} else {
			tt.checkReview(out)
		}
	}
}
