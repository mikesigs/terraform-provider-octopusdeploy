package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateValueFunc([]string{
					"Untenanted",
					"TenantedOrUntenanted",
					"Tenanted",
				}),
			},

		},
	}
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountId := d.Id()
	account, err := client.Account.Get(accountId)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading account %s: %s", accountId, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("account_type", account.AccountType)
	d.Set("client_id", account.ClientId)
	d.Set("tenant_id", account.TenantId)
	d.Set("subscription_id", account.SubscriptionNumber)
	d.Set("client_secret", account.Password)
	d.Set("tenant_tags", account.TenantTags)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation)

	return nil
}

func buildAccountResource(d *schema.ResourceData) *octopusdeploy.Account {
	accountName := d.Get("name").(string)

	var accountType string
	var clientId string
	var tenantId string
	var subscriptionId string
	var clientSecret string
	var tenantTags []string
	var tenantedDeploymentParticipation string

	accountTypeInterface, ok := d.GetOk("account_type")
	if ok {
		accountType = accountTypeInterface.(string)
	}

	clientIdInterface, ok := d.GetOk("client_id")
	if ok {
		clientId = clientIdInterface.(string)
	}

	tenantIdInterface, ok := d.GetOk("tenant_id")
	if ok {
		tenantId = tenantIdInterface.(string)
	}

	subscriptionIdInterface, ok := d.GetOk("subscription_id")
	if ok {
		subscriptionId = subscriptionIdInterface.(string)
	}

	clientSecretInterface, ok := d.GetOk("client_secret")
	if ok {
		clientSecret = clientSecretInterface.(string)
	}

	tenantedDeploymentParticipationInterface, ok := d.GetOk("tenanted_deployment_participation")
	if ok {
		tenantedDeploymentParticipation = tenantedDeploymentParticipationInterface.(string)
	}

	tenantTagsInterface, ok := d.GetOk("tenant_tags")
	if ok {
		tenantTags = getSliceFromTerraformTypeList(tenantTagsInterface)
	}

	if tenantTags == nil {
		tenantTags = []string{}
	}

	var account = octopusdeploy.NewAccount(accountName, accountType)
	account.ClientId = clientId
	account.TenantId = tenantId
	account.Password = octopusdeploy.SensitiveValue{
		NewValue: clientSecret,
	}
	account.SubscriptionNumber = subscriptionId
	account.TenantTags = tenantTags
	account.TenantedDeploymentParticipation = tenantedDeploymentParticipation

	return account;
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newAccount := buildAccountResource(d)
	account, err := client.Account.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildAccountResource(d)
	account.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error updating account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountId := d.Id()

	err := client.Account.Delete(accountId)

	if err != nil {
		return fmt.Errorf("error deleting account id %s: %s", accountId, err.Error())
	}

	d.SetId("")
	return nil
}
