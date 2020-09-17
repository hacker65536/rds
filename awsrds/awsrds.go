package awsrds

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	color "github.com/logrusorgru/aurora"
)

type Rds struct {
	svc *rds.Client
}

func New() *Rds {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	c := &Rds{
		svc: rds.New(cfg),
	}
	return c
}
func (r *Rds) ListRds(filter string) {
	reg := regexp.MustCompile(filter)

	params := &rds.DescribeDBInstancesInput{}
	req := r.svc.DescribeDBInstancesRequest(params)

	p := rds.NewDescribeDBInstancesPaginator(req)

	for p.Next(context.TODO()) {
		page := p.CurrentPage()

		for _, v := range page.DBInstances {
			if reg.MatchString(aws.StringValue(v.DBInstanceIdentifier)) {
				//tags := make(map[string]string)

				var dbstatus string
				dbstatus = aws.StringValue(v.DBInstanceStatus)

				var dbst interface{}
				switch dbstatus {
				case "available":
					dbst = color.Green(dbstatus)
				default:
					dbst = color.Yellow(dbstatus)

				}
				fmt.Printf("%s [%s %s]",
					color.Bold(aws.StringValue(v.DBInstanceIdentifier)),
					aws.StringValue(v.DBInstanceClass),
					dbst,
				)
				fmt.Println()

				t := r.ListTags(v.DBInstanceArn)
				if len(t.ListTagsForResourceOutput.TagList) != 0 {
					r.Printmap(t)
				}
			}
		}
	}

	if err := p.Err(); err != nil {
		fmt.Println(err)
	}
}

func (r *Rds) ListTags(rdsarn *string) *rds.ListTagsForResourceResponse {
	params := &rds.ListTagsForResourceInput{
		ResourceName: rdsarn,
	}
	req := r.svc.ListTagsForResourceRequest(params)
	resp, err := req.Send(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
func (r *Rds) Printmap(t *rds.ListTagsForResourceResponse) {

	for _, v := range t.ListTagsForResourceOutput.TagList {
		//tags[aws.StringValue(v.Key)] = aws.StringValue(v.Value)
		//		fmt.Printf(`"` + aws.StringValue(v.Key) + `"`)
		//		fmt.Printf(`:`)
		//		fmt.Printf(`"` + aws.StringValue(v.Value) + `"`)
		//		fmt.Println()
		//
		fmt.Printf("%s", color.Cyan(aws.StringValue(v.Key)))
		fmt.Printf("\t= ")
		fmt.Printf("%s", color.Cyan(aws.StringValue(v.Value)))
		fmt.Println()
	}
}
