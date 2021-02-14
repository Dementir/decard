package server

import (
	"context"
	"fmt"
	"github.com/Dementir/decard/internal/decard"
	grpc2 "github.com/Dementir/decard/internal/grpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	lg *zap.SugaredLogger
	grpc2.UnimplementedPointServer
}

func InitServer(addr string, lg *zap.SugaredLogger) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "cannot create tcp listner")
	}

	grpcServer := grpc.NewServer()
	server := Server{lg, grpc2.UnimplementedPointServer{}}

	grpc2.RegisterPointServer(grpcServer, &server)

	return grpcServer.Serve(listener)
}

func (s *Server) SortPoints(ctx context.Context, p *grpc2.Points) (*grpc2.Points, error) {
	s.lg.Info("start sort")
	inputPoint := p.Point

	points := make([]decard.Point, 0, len(inputPoint))

	for _, point := range inputPoint {
		splitCords := strings.Split(point, ";")

		if len(splitCords) != 2 {
			err := errors.New("string must have 2 point x and y")
			s.lg.Error(err)

			return nil, err
		}

		x, err := strconv.Atoi(splitCords[0])
		if err != nil {
			err = errors.Wrap(err, "cannot convert x cord from string to int")
			s.lg.Error(err)

			return nil, err
		}

		y, err := strconv.Atoi(splitCords[1])
		if err != nil {
			err = errors.Wrap(err, "cannot convert y cord from string to int")
			s.lg.Error(err)

			return nil, err
		}

		points = append(points, decard.Point{
			X: x,
			Y: y,
		})
	}

	points = decard.Decard(points)

	resultPoints := make([]string, 0, len(points))

	for _, point := range points {
		strPoint := fmt.Sprintf("%d;%d", point.X, point.Y)

		resultPoints = append(resultPoints, strPoint)
	}

	grpcPoint := grpc2.Points{Point: resultPoints}

	return &grpcPoint, nil
}
