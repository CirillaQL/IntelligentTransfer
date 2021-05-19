/*
	获取航班信息包
*/
package service

//爬虫获取航班信息
import (
	"IntelligentTransfer/pkg/logger"
	"context"
	"log"
	"time"

	pb "IntelligentTransfer/shift"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

// GetShiftInfo gRPC接口，用于与python的航班信息爬虫进行交互
func GetShiftInfo(shift, date string) (string, string) {
	//建立链接
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//创建gRPC客户段
	Client := pb.NewGetShiftServiceClient(conn)
	// 设定请求超时时间 10s
	context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	//获取请求
	response, err := Client.GetShift(ctx, &pb.GetShiftReq{
		ShiftNumber: shift,
		Date:        date,
	})
	if err != nil || response == nil {
		logger.ZapLogger.Sugar().Errorf("gRPC service failed. err: %+v ", err)
		return "", ""
	}
	return response.TakeoffTime, response.LandingTime
}
