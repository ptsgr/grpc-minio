package cmd

import (
	"log"
	"os"

	pb "github.com/ptsgr/grpc-minio/internal/server_grpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var GrpcClient pb.ServerClient

const (
	grpcPort = "grpc.port"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "GRPC client for put/get/delete file to minio",
	Long: `GRPC client for put/get/delete file to minio
	put: client put -f file.txt
	get: client get -n xxx/file.txt -o out_file.txt
	delete: client delete -n xxx/file.txt
	`,
}

func Execute() {
	if err := initConfig(); err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":"+viper.GetString(grpcPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	GrpcClient = pb.NewServerClient(conn)

	err = rootCmd.Execute()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {

}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("local")
	return viper.ReadInConfig()
}
