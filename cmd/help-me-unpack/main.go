package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/aryuuu/hackattic/client"
	"github.com/aryuuu/hackattic/constants"
	"github.com/joho/godotenv"
)

type Problem struct {
	Bytes string `json:"bytes"`
}

type Solution struct {
	Int             int32   `json:"int"`
	Uint            uint32  `json:"uint"`
	Short           int16   `json:"short"`
	Float           float32 `json:"float"`
	Double          float64 `json:"double"`
	BigEndianDouble float64 `json:"big_endian_double"`
}

type SolutionResp struct {
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file %w", err)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")
	problemPath := fmt.Sprintf("/challenges/help_me_unpack/problem?access_token=%s", accessToken)
	client := client.NewGeoIPClient(constants.BaseURL)
	ctx := context.Background()

	problem := Problem{}
	err = client.GetProblem(ctx, problemPath, &problem)
	if err != nil {
		log.Println("failed to get problem,", err)
		panic(1)
	}

	log.Println(problem.Bytes)
	plain, err := base64.StdEncoding.DecodeString(problem.Bytes)
	if err != nil {
		log.Println("failed to decode base64 string")
		panic(1)
	}

	log.Println("plain as bytes", plain)
	log.Println("plain as string", string(plain))

	solution := Solution{
		Int:             unpackInt32(plain[0:4]),
		Uint:            unpackUint32(plain[4:8]),
		Short:           unpackInt16(plain[8:10]),
		Float:           math.Float32frombits(unpackUint32(plain[12:16])),
		Double:          math.Float64frombits(unpackUint64(plain[16:24])),
		BigEndianDouble: math.Float64frombits(unpackUint64BigEndian(plain[24:])),
	}

	solvePath := fmt.Sprintf("/challenges/help_me_unpack/solve?access_token=%s", accessToken)
	solres := SolutionResp{}
	client.PostSolution(ctx, solvePath, solution, &solres)

}

func unpackInt32(bytes []byte) int32 {
	var result int32

	for i := len(bytes) - 1; i >= 0; i-- {
		result <<= 8
		result |= int32(bytes[i])
	}

	return result
}

func unpackUint32(bytes []byte) uint32 {
	var result uint32

	for i := len(bytes) - 1; i >= 0; i-- {
		result <<= 8
		result |= uint32(bytes[i])
	}

	return result
}

func unpackInt16(bytes []byte) int16 {
	var result int16

	for i := len(bytes) - 1; i >= 0; i-- {
		result <<= 8
		result |= int16(bytes[i])
	}

	return result
}

func unpackUint64(bytes []byte) uint64 {
	var result uint64

	for i := len(bytes) - 1; i >= 0; i-- {
		result <<= 8
		result |= uint64(bytes[i])
	}

	return result
}

func unpackUint64BigEndian(bytes []byte) uint64 {
	var result uint64

	for _, b := range bytes {
		result <<= 8
		result |= uint64(b)
	}

	return result
}
