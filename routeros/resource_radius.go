package routeros

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const vrfKey = "vrf"

// https://help.mikrotik.com/docs/display/ROS/RADIUS#RADIUS-RADIUSClient
func ResourceRadius() *schema.Resource {
	resSchema := map[string]*schema.Schema{
		MetaResourcePath: PropResourcePath("/radius"),
		MetaId:           PropId(Id),

		"accounting_backup": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "An option whether the configuration is for the backup RADIUS server.",
		},
		"accounting_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1813,
			Description: "RADIUS server port used for accounting.",
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		"address": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "IPv4 or IPv6 address of RADIUS server.",
			ValidateFunc:  validation.IsIPAddress,
		},
		"authentication_port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1812,
			Description: "RADIUS server port used for authentication.",
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		"called_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "RADIUS calling station identifier.",
		},
		"certificate": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "none",
			Description: "Certificate to use for communication with RADIUS Server with RadSec enabled.",
		},
		KeyComment:  PropCommentRw,
		KeyDisabled: PropDisabledRw,
		"domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Microsoft Windows domain of client passed to RADIUS servers that require domain validation.",
		},
		"protocol": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "udp",
			Description: "An option specifies the protocol to use when communicating with the RADIUS Server.",
			ValidateFunc: validation.StringInSlice([]string{"radsec", "udp"}, false),
		},
		"realm": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Explicitly stated realm (user domain), so the users do not have to provide proper ISP domain name in the user name.",
		},
		"secret": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The shared secret to access the RADIUS server.",
		},
		"service": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A comma-separated list of router services that will use the RADIUS server. Possible values: `hotspot`, `login`, `ppp`, `wireless`, `dhcp`.",
		},
		"src_address": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Source IPv4/IPv6 address of the packets sent to the RADIUS server.",
			ValidateFunc: validation.IsIPAddress,
		},
		"timeout": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "300ms",
			Description: "A timeout, after which the request should be resent.",
			DiffSuppressFunc: TimeEquall,
		},
	}

	return &schema.Resource{
		CreateContext: DefaultCreate(resSchema),
		ReadContext: DefaultRead(resSchema),
		UpdateContext: DefaultUpdate(resSchema),
		DeleteContext: DefaultDelete(resSchema),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resSchema,
	}
}

// https://help.mikrotik.com/docs/display/ROS/RADIUS#RADIUS-ConnectionTerminatingfromRADIUS
func ResourceRadiusIncoming() *schema.Resource {
	resSchema := map[string]*schema.Schema{
		MetaResourcePath: PropResourcePath("/radius/incoming"),
		MetaId:           PropId(Name),

		"accept": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "An option whether to accept the unsolicited messages.",
		},
		"port": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3799,
			Description: "The port number to listen for the requests on.",
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		vrfKey: {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "VRF on which service is listening for incoming connections. This option is available in RouterOS starting from version 7.4.",
		},
	}

	return &schema.Resource{
		CreateContext: DefaultSystemCreate(resSchema),
		ReadContext:   compatibleSystemRead(resSchema),
		UpdateContext: DefaultSystemUpdate(resSchema),
		DeleteContext: DefaultSystemDelete(resSchema),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: resSchema,
	}
}

func compatibleSystemRead(s map[string]*schema.Schema) schema.ReadContextFunc {
	return func (ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		diags := SystemResourceRead(ctx, s, d, m)
		if _, exists := d.GetOk(vrfKey); exists {
			s[vrfKey].Default = "main"
		}

		return diags
	}
}