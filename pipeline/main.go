/*  патерн, который разделяет сложную задачу на более легкие, которые выполняются в отдельных горутинах*/
package main

import (
	"context"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Second))
	defer cancel()
}
