package internal // github.com/mikaponics/mikapod-remote/internal

import (
	// "context"
	"log"
	// "os"
	"time"
	// "fmt"
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// "github.com/golang/protobuf/ptypes/timestamp"

    "github.com/mikaponics/mikapod-remote/configs"
	pb "github.com/mikaponics/mikapod-storage/api"
	pb2 "github.com/mikaponics/mikaponics-thing/api"
)

type MikapodRemote struct {
	timer *time.Timer
	ticker *time.Ticker
	done chan bool
	storageCon *grpc.ClientConn
	storage pb.MikapodStorageClient
	remoteCon *grpc.ClientConn
	remote pb2.MikaponicsThingClient
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

    // DEVELOPERS NOTE:
	// (1) If the REMOTE web service is using SSL then we need to force it without
	//     checking the source of origin. This option should be set if you are
	//     running this code in production environment. Special thanks to the
	//     following articles:
	//     - https://www.nginx.com/blog/nginx-1-13-10-grpc/#grpc_pass
	//     - https://itnext.io/effectively-communicate-between-microservices-de7252ba2f3c
	//     - https://www.youtube.com/watch?v=bhiJfNDWRsY
	//     - https://itnext.io/practical-guide-to-securing-grpc-connections-with-go-and-tls-part-1-f63058e9d6d1
	//     - https://itnext.io/practical-guide-to-securing-grpc-connections-with-go-and-tls-part-2-994ef93b8ea9
	// (2) If the REMOTE web service is not using SSL certificate, which is
	//     typically found on your developer environment then run that code.
    var remoteCon *grpc.ClientConn
	var remoteErr error
	if configs.GetIsRemoteUsingSSL() {
		creds := credentials.NewTLS( &tls.Config{ InsecureSkipVerify: true } )
		remoteCon, remoteErr = grpc.Dial(
			mikaponicsRemoteServiceAddress,
			// grpc.WithInsecure(),
			grpc.WithTransportCredentials( creds ),
			grpc.WithTimeout(10*time.Second),
			grpc.WithUnaryInterceptor(unaryInterceptor), // Ex. Added `UnaryInterceptor`.
		)
	} else {
		remoteCon, remoteErr = grpc.Dial(
			mikaponicsRemoteServiceAddress,
			grpc.WithInsecure(),
			grpc.WithTimeout(10*time.Second),
			grpc.WithUnaryInterceptor(unaryInterceptor), // Ex. Added `UnaryInterceptor`.
		)
	}
	if remoteErr != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Set up our protocol buffer interface.
	remote := pb2.NewMikaponicsThingClient(remoteCon)

	return &MikapodRemote{
		timer: nil,
		ticker: nil,
		done: make(chan bool, 1), // Create a execution blocking channel.
		storageCon: storageCon,
		storage: storage,
		remoteCon: remoteCon,
		remote: remote,
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

	//TODO: REMOVE WHEN READY.
	data := app.listTimeSeriesData()
	wasUploaded := app.uploadTimeSeriesData(data)
	if wasUploaded {
		app.deleteTimeSeriesData(data)
	}

    // Setup a background timer which will upload the time-series data to the
	// remote `Mikaponics web service`.
	app.ticker = time.NewTicker(1 * time.Minute)

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
	// app.timer.Stop()
    app.storageCon.Close()
	app.remoteCon.Close()
}

func (app *MikapodRemote) tick()  {
	log.Printf("Uploading local data to remote web service.")
	data := app.listTimeSeriesData()
	wasUploaded := app.uploadTimeSeriesData(data)
	if wasUploaded {
		app.deleteTimeSeriesData(data)
	}
}
