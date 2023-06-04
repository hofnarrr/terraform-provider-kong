package kong

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kong/go-kong/kong"
)

func resourceKongRBACUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKongRBACUserCreate,
		ReadContext:   resourceKongRBACUserRead,
		// TODO: UpdateContext: resourceKongRBACUserUpdate,
		DeleteContext: resourceKongRBACUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			// TODO:
			// user_token
			// enabled
			// roles - https://docs.konghq.com/gateway/latest/admin-api/rbac/reference/#add-a-user-to-a-role
			// comment
		},
	}
}

func resourceKongRBACUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	rbacUserRequest := createKongRBACUserRequestFromResourceData(d)

	client := meta.(*config).adminClient.RBACUsers
	rbacUser, err := client.Create(ctx, rbacUserRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create kong rbacUser: %v error: %v", rbacUserRequest, err))
	}

	d.SetId(*rbacUser.ID)

	return resourceKongRBACUserRead(ctx, d, meta)
}

func resourceKongRBACUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.RBACUsers
	rbacUser, err := client.Get(ctx, kong.String(d.Id()))

	if !kong.IsNotFoundErr(err) && err != nil {
		return diag.FromErr(fmt.Errorf("could not find kong rbacUser: %v", err))
	}

	if rbacUser == nil {
		d.SetId("")
	} else {
		if rbacUser.Name != nil {
			err := d.Set("name", rbacUser.Name)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return diags
}

// TODO:
// func resourceKongRBACUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

// 	rbacUserRequest := createKongRBACUserRequestFromResourceData(d)

// 	client := meta.(*config).adminClient.RBACUsers
// 	rbacUser, err := client.Create(ctx, rbacUserRequest)
// 	if err != nil {
// 		return diag.FromErr(fmt.Errorf("failed to create kong rbacUser: %v error: %v", rbacUserRequest, err))
// 	}

// 	d.SetId(*rbacUser.ID)

// 	return resourceKongRBACUserRead(ctx, d, meta)
// }

func resourceKongRBACUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := meta.(*config).adminClient.RBACUsers
	err := client.Delete(ctx, kong.String(d.Id()))

	if err != nil {
		return diag.FromErr(fmt.Errorf("could not delete kong rbacUser: %v", err))
	}

	return diags
}

func createKongRBACUserRequestFromResourceData(d *schema.ResourceData) *kong.RBACUser {

	rbacUser := &kong.RBACUser{
		Name: readStringPtrFromResource(d, "name"),
	}
	if d.Id() != "" {
		rbacUser.ID = kong.String(d.Id())
	}
	return rbacUser
}
