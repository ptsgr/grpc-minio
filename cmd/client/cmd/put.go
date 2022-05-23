package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	pb "github.com/ptsgr/grpc-minio/internal/server_grpc"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "put command send file to MinIO server",
	Long: `put command send file to MinIO server
		use -f/--file for set file for sending
		`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Println("Please fix file param: ", err.Error())
			os.Exit(1)
		}

		fileData, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Println("File not exist: ", err.Error())
			os.Exit(1)
		}

		resp, err := GrpcClient.Put(context.Background(), &pb.PutRequest{
			Filename: fileName,
			FileData: fileData,
		})
		if err != nil {
			log.Println("Put file error: ", err.Error())
			os.Exit(1)
		}

		log.Println(resp.Message)

	},
}

func init() {
	rootCmd.AddCommand(putCmd)
	putCmd.Flags().StringP("file", "f", "", "File for sending")
}
