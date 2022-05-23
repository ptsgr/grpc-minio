package cmd

import (
	"context"
	"log"
	"os"

	pb "github.com/ptsgr/grpc-minio/internal/server_grpc"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command remove file from MinIO server",
	Long: `delete command remove file from MinIO server
		use -n/--filename for set MinIO filename
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName, err := cmd.Flags().GetString("filename")
		if err != nil {
			log.Println("Please fix filename param: ", err.Error())
			os.Exit(1)
		}

		resp, err := GrpcClient.Delete(context.Background(), &pb.DeleteRequest{
			Filename: fileName,
		})
		if err != nil {
			log.Println("Delete file error: ", err.Error())
			os.Exit(1)
		}

		log.Println(resp.Message)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("filename", "n", "", "Filename for delete")
}
