package internal // github.com/mikaponics/mikapod-remote/internal

import (
	// "context"
	"log"
	// "os"
	"time"
	// "fmt"

	"google.golang.org/grpc"
	// "github.com/golang/protobuf/ptypes/timestamp"

    // "github.com/mikaponics/mikapod-remote/configs"
	pb "github.com/mikaponics/mikapod-storage/api"
	// pb2 "github.com/mikaponics/mikapod-soil-reader/api"   //TODO: Uncommment
)

type MikapodRemote struct {
	timer *time.Timer
	ticker *time.Ticker
	done chan bool
	storageCon *grpc.ClientConn
	storage pb.MikapodStorageClient
	// remoteCon *grpc.ClientConn                             //TODO: Uncommment
	// remote pb2.MikapodSoilReaderClient                     //TODO: Uncommment
}

// Function will construct the Mikapod Remote application.
func InitMikapodRemote(mikapodStorageAddress string, mikaponicsRemoteServiceAddress string) (*MikapodRemote) {
	// Set up a direct connection to the `mikapod-storage` server.
	storageCon, err := grpc.Dial(mikapodStorageAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Set up our protocol buffer interface.
	storage := pb.NewMikapodStorageClient(storageCon)

    // TODO: Uncommment
    // // Set up a direct connection to the `mikapod-soil-remote` server.
	// remoteCon, remoteErr := grpc.Dial(mikaponicsRemoteServiceAddress, grpc.WithInsecure())
	// if remoteErr != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	//
	// // Set up our protocol buffer interface.
	// remote := pb2.NewMikapodSoilReaderClient(remoteCon)

	return &MikapodRemote{
		timer: nil,
		ticker: nil,
		done: make(chan bool, 1), // Create a execution blocking channel.
		storageCon: storageCon,
		storage: storage,
		// remoteCon: remoteCon,                              //TODO: Uncommment
		// remote: remote,                                    //TODO: Uncommment
	}
}

// Source: https://www.reddit.com/r/golang/comments/44tmti/scheduling_a_function_call_to_the_exact_start_of/
func minuteTicker() *time.Timer {
	// Current time
	now := time.Now()

	// Get the number of seconds until the next minute
	var d time.Duration
	d = time.Second * time.Duration(60-now.Second())

	// Time of the next tick
	nextTick := now.Add(d)

	// Subtract next tick from now
	diff := nextTick.Sub(time.Now())

	// Return new ticker
	return time.NewTimer(diff)
}


// Function will consume the main runtime loop and run the business logic
// of the Mikapod Remote application.
func (app *MikapodRemote) RunMainRuntimeLoop() {
	defer app.shutdown()

    // //TODO: UNCOMMENT WHEN READY
    // // DEVELOPERS NOTE:
	// // (1) The purpose of this block of code is to find the future date where
	// //     the minute just started, ex: 5:00 AM, 5:01, etc, and then start our
	// //     main runtime loop to run along for every minute afterwords.
	// // (2) If our application gets terminated by the user or system then we
	// //     terminate our timer.
    // log.Printf("Synching with local time...")
	// app.timer = minuteTicker()
	// select {
	// 	case <- app.timer.C:
	// 		log.Printf("Synchronized with local time.")
	// 		app.ticker = time.NewTicker(1 * time.Minute)
	// 	case <- app.done:
	// 		app.timer.Stop()
	// 		log.Printf("Interrupted timer.")
	// 		return
	// }

    // THIS CODE IS FOR TESTING, REMOVE WHEN READY TO USE, UNCOMMENT ABOVE.
	app.ticker = time.NewTicker(1 * time.Minute)

	//TODO: REMOVE WHEN READY.
	data := app.listTimeSeriesData()
	wasUploaded := app.uploadTimeSeriesData(data)
	if wasUploaded {
		app.deleteTimeSeriesData(data)
	}

    // DEVELOPERS NOTE:
	// (1) The purpose of this block of code is to run as a goroutine in the
	//     background as an anonymous function waiting to get either the
	//     ticker chan or app termination chan response.
	// (2) Main runtime loop's execution is blocked by the `done` chan which
	//     can only be triggered when this application gets a termination signal
	//     from the operating system.
	log.Printf("Remote is now running.")
	go func() {
        for {
            select {
	            case <- app.ticker.C:
					app.tick()
				case <- app.done:
					app.ticker.Stop()
					log.Printf("Interrupted ticker.")
					return
			}
		}
	}()
	<-app.done
}

// Function will tell the application to stop the main runtime loop when
// the process has been finished.
func (app *MikapodRemote) StopMainRuntimeLoop() {
	app.done <- true
}

func (app *MikapodRemote) shutdown()  {
    app.storageCon.Close()
	// app.remoteCon.Close()                                  //TODO: Uncommment
}

func (app *MikapodRemote) tick()  {
	log.Printf("Tick: Do something...")
	data := app.listTimeSeriesData()
	wasUploaded := app.uploadTimeSeriesData(data)
	if wasUploaded {
		app.deleteTimeSeriesData(data)
	}
}
