package main

import (
	"bufio"
	"context"
	"flag"
	grpc2 "github.com/Dementir/decard/internal/grpc"
	"github.com/Dementir/decard/internal/logger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"os"
)

const serverAddr = "localhost:9001"

type client struct {
	cc grpc.ClientConnInterface
}

func main() {
	lg := logger.NewLogger("DEBUG")
	defer lg.Sync()

	inputPath := flag.String("i", "input.txt", "set config path")
	outputPath := flag.String("o", "output.txt", "set config path")
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		lg.Fatalw("fail to dial", "error", err)
	}
	defer conn.Close()

	lg.Info("run client")

	client := client{conn}
	points, err := client.ReadPoints(*inputPath)
	if err != nil {
		lg.Fatalw("cannot read points from file", "error", err)
	}

	result, err := client.SortPoints(context.Background(), points)
	if err != nil {
		lg.Fatalw("cannot sort points", "error", err)
	}

	err = client.SavePoints(*outputPath, result)
	if err != nil {
		lg.Fatalw("cannot save points to file", "error", err)
	}

	lg.Info(result)
}

func (c *client) ReadPoints(path string) (*grpc2.Points, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	points := make([]string, 0, 100)

	for scanner.Scan() {
		points = append(points, scanner.Text())
	}

	grpcPoint := grpc2.Points{
		Point: points,
	}

	return &grpcPoint, nil
}

func (c *client) SavePoints(path string, points *grpc2.Points) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	for _, point := range points.Point {

		_, err := w.WriteString(point + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *client) SortPoints(ctx context.Context, in *grpc2.Points, opts ...grpc.CallOption) (*grpc2.Points, error) {
	out := new(grpc2.Points)
	err := c.cc.Invoke(ctx, "/decard.Point/SortPoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
