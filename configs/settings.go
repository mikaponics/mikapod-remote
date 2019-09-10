package configs

import "os"

const (
	MikapodStorageServiceAddress = "localhost:50051"             // Please do not change!
)

func GetIsRemoteUsingSSL() bool {
	return os.Getenv("MIKAPONICS_REMOTE_IS_USING_SSL") == "YES"
}
