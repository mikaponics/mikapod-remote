# Mikapod Remote
[![Go Report Card](https://goreportcard.com/badge/github.com/mikaponics/mikapod-remote)](https://goreportcard.com/report/github.com/mikaponics/mikapod-remote)

## Overview

The purpose of this application is to poll time-series data from our [Mikapod Soil Reader](https://github.com/mikaponics/mikapod-soil-reader) application and save it to the [Mikapod Storage](https://github.com/mikaponics/mikapod-storage) application. The interval of time is every one minute.

## Prerequisites

You must have the following installed before proceeding. If you are missing any one of these then you cannot begin.

* ``Go 1.12.7``

## Installation

1. Please visit the [Mikapod Storage](https://github.com/mikaponics/mikapod-storage) repository and setup that application on your device.

2. Get our latest code.

    ```
    go get -u github.com/mikaponics/mikapod-remote
    ```

5. Install the depencies for this project.

    ```
    go get -u google.golang.org/grpc
    ```

6. Setup our environment variable. Please change to the address of our remote server. Please adjust the sensor ID values based what was assugned frin the Mikaponics web service. If the data belongs to a different tenant please adjust the value, else leave as is!

    ```
    # Mikaponics Web Service
    export MIKAPONICS_REMOTE_APP_ADDRESS="localhost:50053"

    # Tenancy
    export MIKAPOD_TENANT_ID=1

    # Sensor IDs
    export MIKAPOD_HUMIDITY_SENSOR_ID=1
    export MIKAPOD_TEMPERATURE_SENSOR_ID=2
    export MIKAPOD_PRESSURE_SENSOR_ID=3
    export MIKAPOD_TEMPERATURE_BACKEND_SENSOR_ID=3
    export MIKAPOD_ALTITUDE_SENSOR_ID=4
    export MIKAPOD_ILLUMINANCE_SENSOR_ID=5
    export MIKAPOD_SOIL_MOISTURE_SENSOR_ID=6
    ```

7. Run our application.

    ```
    cd github.com/mikaponics/mikapod-remote
    go run main.go
    ```

## Production

The following instructions are specific to getting setup for [Raspberry Pi](https://www.raspberrypi.org/).

### Deployment

1. Please visit the [Mikapod Storage](https://github.com/mikaponics/mikapod-storage) repository and setup that application on your device.

2. (Optional) If already installed old golang with apt-get and you want to upgrade to the latest version. Run the following:

    ```
    sudo apt remove golang
    sudo apt-get autoremove
    source .profile
    ```

3. Install [Golang 1.11.8]():

    ```
    wget https://storage.googleapis.com/golang/go1.11.8.linux-armv6l.tar.gz
    sudo tar -C /usr/local -xzf go1.11.8.linux-armv6l.tar.gz
    export PATH=$PATH:/usr/local/go/bin # put into ~/.profile
    ```

4. Confirm we are using the correct version:

    ```
    go version
    ```

5. Install ``git``:

    ```
    sudo apt install git
    ```

6. Get our latest code.

    ```
    go get -u github.com/mikaponics/mikapod-remote
    ```

7. Install the depencies for this project.

    ```
    go get -u google.golang.org/grpc
    ```

8. Go to our application directory.

    ```
    cd ~/go/src/github.com/mikaponics/mikapod-remote
    ```

9. (Optional) Confirm our application builds on the raspberry pi device. You now should see a message saying ``gRPC server is running`` then the application is running.

    ```
    go run main.go
    ```

10. Build for the ARM device and install it in our ``~/go/bin`` folder:

    ```
    go install
    ```


### Operation

1. While being logged in as ``pi`` run the following:

    ```
    sudo vi /etc/systemd/system/mikapod-remote.service
    ```

2. Copy and paste the following contents.

    ```
    [Unit]
    Description=Mikapod Remote Daemon
    After=multi-user.target

    [Service]
    Type=idle
    ExecStart=/home/pi/go/bin/mikapod-remote
    Restart=on-failure
    KillSignal=SIGTERM

    [Install]
    WantedBy=multi-user.target
    ```

3. We can now start the Gunicorn service we created and enable it so that it starts at boot:

    ```
    sudo systemctl start mikapod-remote
    sudo systemctl enable mikapod-remote
    ```

4. Confirm our service is running.

    ```
    sudo systemctl status mikapod-remote.service
    ```

5. If the service is working correctly you should see something like this at the bottom:

    ```
    raspberrypi systemd[1]: Started Mikapod Poller Daemon.
    ```

6. Congradulations, you have setup instrumentation micro-service! All other micro-services can now poll the latest data from the soil-reader we have attached.

7. If you see any problems, run the following service to see what is wrong. More information can be found in [this article](https://unix.stackexchange.com/a/225407).

    ```
    sudo journalctl -u mikapod-remote
    ```

8. To reload the latest modifications to systemctl file.

    ```
    sudo systemctl daemon-reload
    ```

## Troubleshooting

If you get an errors on your Raspberry stating: ``go not found`` then please run the following:

```
export PATH=$PATH:/usr/local/go/bin # put into ~/.profile
```

## License

This application is licensed under the **BSD** license. See [LICENSE](LICENSE) for more information.
