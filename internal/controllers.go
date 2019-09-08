package internal // github.com/mikaponics/mikapod-poller/internal

import (
	"context"
	"log"
	"time"

	pb "github.com/mikaponics/mikapod-storage/api"
	pb2 "github.com/mikaponics/mikaponics-thing/api"
)


func (app *MikapodRemote) listTimeSeriesData() ([]*TimeSeriesDatum){
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := app.storage.ListTimeSeriesData(ctx, &pb.ListTimeSeriesDataRequest{})
	if err != nil {
		log.Fatalf("could not list from storage: %v", err)
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
    // Convert our `struct` formatted list to be of `protocol buffer`
	// formatted list which we can use in our `grpc` output.
	var list []*pb2.SetTimeSeriesDatumRequest
	for _, v := range data {
        // Create our `protocol buffer` single time-series datum object.
        ri := &pb2.SetTimeSeriesDatumRequest{
			TenantId:   GetTenantId(),
            SensorId:   GetSensorIdFromInstrumentId(v.Instrument),
            Value:      v.Value,
			Timestamp:  v.Timestamp,
        }

		log.Printf("UPLOADING %v", ri)

		// Attach our single time-series datum object to our `protocol buffer`
		// list of time-series data.
        list = append(list, ri)
    }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := app.remote.SetTimeSeriesData(ctx, &pb2.SetTimeSeriesDataRequest{
		Data: list,
	})
	if err != nil {
		log.Fatalf("could not add time-series data to remote: %v", err)
	}

    //TODO: PLEASE SWAP THE FOLLOWING COMMENTS WHEN THE REMOTE SERVER IS WORKING.
	// return err == nil
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
		log.Fatalf("could not delete time-series data to storage: %v", err)
	}
}
