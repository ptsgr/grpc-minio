package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pb "github.com/ptsgr/grpc-minio/internal/server_grpc"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get command return file from MinIO server",
	Long: `get command return file from MinIO server
		use -n/--filename for set MinIO filename
		use -o/--out for set output filename
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName, err := cmd.Flags().GetString("filename")
		if err != nil {
			log.Println("Please fix filename param: ", err.Error())
			os.Exit(1)
		}

		outFileName, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Println("Please fix output filename param: ", err.Error())
			os.Exit(1)
		}

		fileExtension := filepath.Ext(fileName)
		outFileExtension := filepath.Ext(outFileName)

		if outFileExtension == "" && fileExtension != "" {
			outFileName += fileExtension
		}

		resp, err := GrpcClient.Get(context.Background(), &pb.GetRequest{
			Filename: fileName,
		})
		if err != nil || resp.Code != pb.StatusCode_Ok {
			log.Println("Get file error: ", err.Error())
			os.Exit(1)
		}

		err = ioutil.WriteFile(outFileName, resp.FileData, 0644)
		if err != nil {
			log.Println("Cannot create output file: ", err.Error())
			os.Exit(1)
		}

		log.Println(resp.Message)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("filename", "n", "", "Filename for getting")
	getCmd.Flags().StringP("out", "o", "output_file", "Output file name")

}
