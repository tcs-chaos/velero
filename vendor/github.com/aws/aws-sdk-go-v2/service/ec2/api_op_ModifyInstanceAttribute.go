// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Modifies the specified attribute of the specified instance. You can specify
// only one attribute at a time. Note: Using this action to change the security
// groups associated with an elastic network interface (ENI) attached to an
// instance can result in an error if the instance has more than one ENI. To change
// the security groups associated with an ENI attached to an instance that has
// multiple ENIs, we recommend that you use the ModifyNetworkInterfaceAttribute
// action. To modify some attributes, the instance must be stopped. For more
// information, see Modify a stopped instance (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_ChangingAttributesWhileInstanceStopped.html)
// in the Amazon EC2 User Guide.
func (c *Client) ModifyInstanceAttribute(ctx context.Context, params *ModifyInstanceAttributeInput, optFns ...func(*Options)) (*ModifyInstanceAttributeOutput, error) {
	if params == nil {
		params = &ModifyInstanceAttributeInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "ModifyInstanceAttribute", params, optFns, c.addOperationModifyInstanceAttributeMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ModifyInstanceAttributeOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ModifyInstanceAttributeInput struct {

	// The ID of the instance.
	//
	// This member is required.
	InstanceId *string

	// The name of the attribute to modify. You can modify the following attributes
	// only: disableApiTermination | instanceType | kernel | ramdisk |
	// instanceInitiatedShutdownBehavior | blockDeviceMapping | userData |
	// sourceDestCheck | groupSet | ebsOptimized | sriovNetSupport | enaSupport |
	// nvmeSupport | disableApiStop | enclaveOptions
	Attribute types.InstanceAttributeName

	// Modifies the DeleteOnTermination attribute for volumes that are currently
	// attached. The volume must be owned by the caller. If no value is specified for
	// DeleteOnTermination , the default is true and the volume is deleted when the
	// instance is terminated. To add instance store volumes to an Amazon EBS-backed
	// instance, you must add them when you launch the instance. For more information,
	// see Update the block device mapping when launching an instance (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/block-device-mapping-concepts.html#Using_OverridingAMIBDM)
	// in the Amazon EC2 User Guide.
	BlockDeviceMappings []types.InstanceBlockDeviceMappingSpecification

	// Indicates whether an instance is enabled for stop protection. For more
	// information, see Stop Protection (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Stop_Start.html#Using_StopProtection)
	// .
	DisableApiStop *types.AttributeBooleanValue

	// If the value is true , you can't terminate the instance using the Amazon EC2
	// console, CLI, or API; otherwise, you can. You cannot use this parameter for Spot
	// Instances.
	DisableApiTermination *types.AttributeBooleanValue

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// Specifies whether the instance is optimized for Amazon EBS I/O. This
	// optimization provides dedicated throughput to Amazon EBS and an optimized
	// configuration stack to provide optimal EBS I/O performance. This optimization
	// isn't available with all instance types. Additional usage charges apply when
	// using an EBS Optimized instance.
	EbsOptimized *types.AttributeBooleanValue

	// Set to true to enable enhanced networking with ENA for the instance. This
	// option is supported only for HVM instances. Specifying this option with a PV
	// instance can make it unreachable.
	EnaSupport *types.AttributeBooleanValue

	// Replaces the security groups of the instance with the specified security
	// groups. You must specify the ID of at least one security group, even if it's
	// just the default security group for the VPC.
	Groups []string

	// Specifies whether an instance stops or terminates when you initiate shutdown
	// from the instance (using the operating system command for system shutdown).
	InstanceInitiatedShutdownBehavior *types.AttributeValue

	// Changes the instance type to the specified value. For more information, see
	// Instance types (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html)
	// in the Amazon EC2 User Guide. If the instance type is not valid, the error
	// returned is InvalidInstanceAttributeValue .
	InstanceType *types.AttributeValue

	// Changes the instance's kernel to the specified value. We recommend that you use
	// PV-GRUB instead of kernels and RAM disks. For more information, see PV-GRUB (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/UserProvidedKernels.html)
	// .
	Kernel *types.AttributeValue

	// Changes the instance's RAM disk to the specified value. We recommend that you
	// use PV-GRUB instead of kernels and RAM disks. For more information, see PV-GRUB (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/UserProvidedKernels.html)
	// .
	Ramdisk *types.AttributeValue

	// Enable or disable source/destination checks, which ensure that the instance is
	// either the source or the destination of any traffic that it receives. If the
	// value is true , source/destination checks are enabled; otherwise, they are
	// disabled. The default value is true . You must disable source/destination checks
	// if the instance runs services such as network address translation, routing, or
	// firewalls.
	SourceDestCheck *types.AttributeBooleanValue

	// Set to simple to enable enhanced networking with the Intel 82599 Virtual
	// Function interface for the instance. There is no way to disable enhanced
	// networking with the Intel 82599 Virtual Function interface at this time. This
	// option is supported only for HVM instances. Specifying this option with a PV
	// instance can make it unreachable.
	SriovNetSupport *types.AttributeValue

	// Changes the instance's user data to the specified value. If you are using an
	// Amazon Web Services SDK or command line tool, base64-encoding is performed for
	// you, and you can load the text from a file. Otherwise, you must provide
	// base64-encoded text.
	UserData *types.BlobAttributeValue

	// A new value for the attribute. Use only with the kernel , ramdisk , userData ,
	// disableApiTermination , or instanceInitiatedShutdownBehavior attribute.
	Value *string

	noSmithyDocumentSerde
}

type ModifyInstanceAttributeOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationModifyInstanceAttributeMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsEc2query_serializeOpModifyInstanceAttribute{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpModifyInstanceAttribute{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "ModifyInstanceAttribute"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addOpModifyInstanceAttributeValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opModifyInstanceAttribute(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opModifyInstanceAttribute(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "ModifyInstanceAttribute",
	}
}
