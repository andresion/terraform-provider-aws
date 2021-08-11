package cloudwatch_test

import (
	"strings"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestValidEventCustomEventBusEventSourceName(t *testing.T) {
	cases := []struct {
		Value   string
		IsValid bool
	}{
		{
			Value:   "",
			IsValid: false,
		},
		{
			Value:   "default",
			IsValid: false,
		},
		{
			Value:   "aws.partner/example.com/test/" + sdkacctest.RandStringFromCharSet(227, sdkacctest.CharSetAlpha),
			IsValid: true,
		},
		{
			Value:   "aws.partner/example.com/test/" + sdkacctest.RandStringFromCharSet(228, sdkacctest.CharSetAlpha),
			IsValid: false,
		},
		{
			Value:   "aws.partner/example.com/test/12345ab-cdef-1235",
			IsValid: true,
		},
		{
			Value:   "/test0._1-",
			IsValid: false,
		},
		{
			Value:   "test0._1-",
			IsValid: false,
		},
	}
	for _, tc := range cases {
		_, errors := validEventCustomEventBusEventSourceName(tc.Value, "aws_cloudwatch_event_bus_event_source_name")
		isValid := len(errors) == 0
		if tc.IsValid && !isValid {
			t.Errorf("expected %q to return valid, but did not", tc.Value)
		} else if !tc.IsValid && isValid {
			t.Errorf("expected %q to not return valid, but did", tc.Value)
		}
	}
}

func TestValidEventCustomEventBusName(t *testing.T) {
	cases := []struct {
		Value   string
		IsValid bool
	}{
		{
			Value:   "",
			IsValid: false,
		},
		{
			Value:   "default",
			IsValid: false,
		},
		{
			Value:   sdkacctest.RandStringFromCharSet(256, sdkacctest.CharSetAlpha),
			IsValid: true,
		},
		{
			Value:   sdkacctest.RandStringFromCharSet(257, sdkacctest.CharSetAlpha),
			IsValid: false,
		},
		{
			Value:   "aws.partner/example.com/test/12345ab-cdef-1235",
			IsValid: true,
		},
		{
			Value:   "/test0._1-",
			IsValid: true,
		},
		{
			Value:   "test0._1-",
			IsValid: true,
		},
	}
	for _, tc := range cases {
		_, errors := validEventCustomEventBusName(tc.Value, "aws_cloudwatch_event_bus")
		isValid := len(errors) == 0
		if tc.IsValid && !isValid {
			t.Errorf("expected %q to return valid, but did not", tc.Value)
		} else if !tc.IsValid && isValid {
			t.Errorf("expected %q to not return valid, but did", tc.Value)
		}
	}
}

func TestValidDashboardName(t *testing.T) {
	validNames := []string{
		"HelloWorl_d",
		"hello-world",
		"hello-world-012345",
	}
	for _, v := range validNames {
		_, errors := validDashboardName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid CloudWatch dashboard name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"special@character",
		"slash/in-the-middle",
		"dot.in-the-middle",
		strings.Repeat("W", 256), // > 255
	}
	for _, v := range invalidNames {
		_, errors := validDashboardName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid CloudWatch dashboard name", v)
		}
	}
}

func TestValidEventBusNameOrARN(t *testing.T) {
	validNames := []string{
		"HelloWorl_d",
		"hello-world",
		"hello.World0125",
		"aws.partner/mongodb.com/stitch.trigger/something",        // nosemgrep: domain-names
		"arn:aws:events:us-east-1:123456789012:event-bus/default", // lintignore:AWSAT003,AWSAT005
	}
	for _, v := range validNames {
		_, errors := validEventBusNameOrARN(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid CW event rule name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"special@character",
		"arn:aw:events:us-east-1:123456789012:event-bus/default", // lintignore:AWSAT003,AWSAT005
	}
	for _, v := range invalidNames {
		_, errors := validEventBusNameOrARN(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid CW event rule name", v)
		}
	}
}

func TestValidEventRuleName(t *testing.T) {
	validNames := []string{
		"HelloWorl_d",
		"hello-world",
		"hello.World0125",
	}
	for _, v := range validNames {
		_, errors := validEventRuleName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid CW event rule name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"special@character",
		"slash/in-the-middle",
		// Length > 64
		"TooLooooooooooooooooooooooooooooooooooooooooooooooooooooooongName",
	}
	for _, v := range invalidNames {
		_, errors := validEventRuleName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid CW event rule name", v)
		}
	}
}

func TestValidEC2AutomateARN(t *testing.T) {
	validNames := []string{
		"arn:aws:automate:us-east-1:ec2:reboot",    //lintignore:AWSAT003,AWSAT005
		"arn:aws:automate:us-east-1:ec2:recover",   //lintignore:AWSAT003,AWSAT005
		"arn:aws:automate:us-east-1:ec2:stop",      //lintignore:AWSAT003,AWSAT005
		"arn:aws:automate:us-east-1:ec2:terminate", //lintignore:AWSAT003,AWSAT005
	}
	for _, v := range validNames {
		_, errors := validEC2AutomateARN(v, "test_property")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid ARN: %q", v, errors)
		}
	}

	invalidNames := []string{
		"",
		"arn:aws:elasticbeanstalk:us-east-1:123456789012:environment/My App/MyEnvironment", // lintignore:AWSAT003,AWSAT005 // Beanstalk
		"arn:aws:iam::123456789012:user/David",                                             // lintignore:AWSAT005          // IAM User
		"arn:aws:rds:eu-west-1:123456789012:db:mysql-db",                                   // lintignore:AWSAT003,AWSAT005 // RDS
		"arn:aws:s3:::my_corporate_bucket/exampleobject.png",                               // lintignore:AWSAT005          // S3 object
		"arn:aws:events:us-east-1:319201112229:rule/rule_name",                             // lintignore:AWSAT003,AWSAT005 // CloudWatch Rule
		"arn:aws:lambda:eu-west-1:319201112229:function:myCustomFunction",                  // lintignore:AWSAT003,AWSAT005 // Lambda function
		"arn:aws:lambda:eu-west-1:319201112229:function:myCustomFunction:Qualifier",        // lintignore:AWSAT003,AWSAT005 // Lambda func qualifier
		"arn:aws-us-gov:s3:::corp_bucket/object.png",                                       // lintignore:AWSAT005          // GovCloud ARN
		"arn:aws-us-gov:kms:us-gov-west-1:123456789012:key/some-uuid-abc123",               // lintignore:AWSAT003,AWSAT005 // GovCloud KMS ARN
	}
	for _, v := range invalidNames {
		_, errors := validEC2AutomateARN(v, "test_property")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid ARN", v)
		}
	}
}

func TestValidLogGroupName(t *testing.T) {
	validNames := []string{
		"ValidLogGroupName",
		"ValidLogGroup.Name",
		"valid/Log-group",
		"1234",
		"YadaValid#0123",
		"Also_valid-name",
		strings.Repeat("W", 512),
	}
	for _, v := range validNames {
		_, errors := validLogGroupName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Log Group name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"also $ invalid",
		"This . is also %% invalid@!)+(",
		"*",
		"",
		// length > 512
		strings.Repeat("W", 513),
	}
	for _, v := range invalidNames {
		_, errors := validLogGroupName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Log Group name", v)
		}
	}
}

func TestValidLogGroupNamePrefix(t *testing.T) {
	validNames := []string{
		"ValidLogGroupName",
		"ValidLogGroup.Name",
		"valid/Log-group",
		"1234",
		"YadaValid#0123",
		"Also_valid-name",
		strings.Repeat("W", 483),
	}
	for _, v := range validNames {
		_, errors := validLogGroupNamePrefix(v, "name_prefix")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Log Group name prefix: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"also $ invalid",
		"This . is also %% invalid@!)+(",
		"*",
		"",
		// length > 483
		strings.Repeat("W", 484),
	}
	for _, v := range invalidNames {
		_, errors := validLogGroupNamePrefix(v, "name_prefix")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Log Group name prefix", v)
		}
	}
}

func TestValidLogMetricFilterName(t *testing.T) {
	validNames := []string{
		"YadaHereAndThere",
		"Valid-5Metric_Name",
		"This . is also %% valid@!)+(",
		"1234",
		strings.Repeat("W", 512),
	}
	for _, v := range validNames {
		_, errors := validLogMetricFilterName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Log Metric Filter Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"*",
		// length > 512
		strings.Repeat("W", 513),
	}
	for _, v := range invalidNames {
		_, errors := validLogMetricFilterName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Log Metric Filter Name", v)
		}
	}
}

func TestValidLogMetricTransformationName(t *testing.T) {
	validNames := []string{
		"YadaHereAndThere",
		"Valid-5Metric_Name",
		"This . is also %% valid@!)+(",
		"1234",
		"",
		strings.Repeat("W", 255),
	}
	for _, v := range validNames {
		_, errors := validLogMetricFilterTransformationName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Log Metric Filter Transformation Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Here is a name with: colon",
		"and here is another * invalid name",
		"also $ invalid",
		"*",
		// length > 255
		strings.Repeat("W", 256),
	}
	for _, v := range invalidNames {
		_, errors := validLogMetricFilterTransformationName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Log Metric Filter Transformation Name", v)
		}
	}
}
