package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/aws"
)

const (
	outputFolder   = "web_service"
	outputFilename = "web_service"
)

func main() {
	// remove output folder if exists
	_ = os.RemoveAll(outputFolder)

	d, err := diagram.New(func(options *diagram.Options) {
		options.Name = outputFolder
	}, diagram.Filename(outputFilename), diagram.Label("Web Service"), diagram.Direction(string(diagram.BottomToTop)))
	if err != nil {
		log.Fatal(err)
	}

	lb := aws.Network.ElasticLoadBalancing(diagram.NodeLabel("lb"))
	web := aws.Compute.Ec2(diagram.NodeLabel("web"))
	db := aws.Database.Rds(diagram.NodeLabel("db"))

	batchCluster := diagram.NewGroup("batch").Label("Batches")
	for i := 1; i <= 2; i++ {
		batchCluster.Add(aws.Compute.Ec2(diagram.NodeLabel(fmt.Sprintf("batch_%d", i))))
	}

	d.Connect(lb, web)
	d.Connect(web, db)
	d.Group(batchCluster)
	batchCluster.ConnectAllTo(db.ID(), func(options *diagram.EdgeOptions) {
		options.Label = "r/w"
	})
	batchCluster.ConnectAllFrom(db.ID(), func(options *diagram.EdgeOptions) {
		options.Label = "trigger"
		options.Style = "dotted"
		options.Color = "darkgreen"
	})

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("bash", "-c", fmt.Sprintf(
		"cd %[1]s && cat %[2]s.dot | dot -Tpng > %[2]s.png && open %[2]s.png", outputFolder, outputFilename,
	)).Run(); err != nil {
		log.Fatal(err)
	}
}
