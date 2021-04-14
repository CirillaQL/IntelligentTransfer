/*
	获取航班信息包
*/
package service

//爬虫获取航班信息
import (
	"context"
	"fmt"
	"log"
	"time"

	pb "IntelligentTransfer/shift"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func Grpct() {
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	Client := pb.NewGetShiftServiceClient(conn)

	// 设定请求超时时间 3s
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	// UserIndex 请求
	Reponse, err := Client.GetShift(ctx, &pb.GetShiftReq{
		ShiftNumber: "HU7601",
		Date:        "2021-04-14",
	})
	fmt.Println(Reponse)
}
