package internal // github.com/mikaponics/mikapod-poller/internal

import (
	"context"
	"log"
	// // "os"
	"time"
	// "fmt"

	// "google.golang.org/grpc"
	// "github.com/golang/protobuf/ptypes/timestamp"

    // "github.com/mikaponics/mikapod-poller/configs"
	pb "github.com/mikaponics/mikapod-storage/api"
	// pb2 "github.com/mikaponics/mikapod-soil-reader/api"
)


func (app *MikapodRemote) listTimeSeriesData() ([]*TimeSeriesDatum){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := app.storage.ListTimeSeriesData(ctx, &pb.ListTimeSeriesDataRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

    // The following block of code will iterate through our `protocol buffer`
	// data and convert it to our golang `struct` data so our application can
	// use it.
	var list []*TimeSeriesDatum
	for _, v := range r.Data {
		if v.Id > 0 { // Only accept valid time series data!
			tsd := &TimeSeriesDatum{
	            Id:         v.Id,
	            Instrument: v.Instrument,
	            Value:      v.Value,
				Timestamp:  v.Timestamp,
	        }
	        list = append(list, tsd)
		}
    }
	return list
}

func (app *MikapodRemote) uploadTimeSeriesData(data []*TimeSeriesDatum) bool {
	for _, v := range data {
		log.Printf("TODO - UPLOAD DATA: %v", v)
	}
	return false
}

func (app *MikapodRemote) deleteTimeSeriesData(data []*TimeSeriesDatum) {
	var pks []int64
	for _, v := range data {
		pks = append(pks, v.Id)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := app.storage.DeleteTimeSeriesData(ctx, &pb.DeleteTimeSeriesDataByPKsRequest{
		Pks: pks,
	})
	if err != nil {
		log.Fatalf("could not add time-series data to storage: %v", err)
	}
}
