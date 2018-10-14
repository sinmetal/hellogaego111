package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func SpannerSimpleQueryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := spannerService.SimpleQuery(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DONE"))
}

func NewSpannerClient(ctx context.Context, database string) (*spanner.Client, error) {
	return spanner.NewClient(ctx, database)
}

type SpannerService struct {
	Client *spanner.Client
}

func NewSpannerService(client *spanner.Client) *SpannerService {
	return &SpannerService{
		Client: client,
	}
}

func (s *SpannerService) SimpleQuery(ctx context.Context) error {
	iter := s.Client.Single().Query(ctx, spanner.NewStatement("SELECT 1"))
	defer iter.Stop()
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
