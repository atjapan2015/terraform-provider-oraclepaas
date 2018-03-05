package oraclepaas

import (
	"fmt"
	"log"
	"strconv"
	"time"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPAASDatabaseServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPAASDatabaseServiceInstanceCreate,
		Read:   resourceOPAASDatabaseServiceInstanceRead,
		Delete: resourceOPAASDatabaseServiceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceStandardEdition),
					string(database.ServiceInstanceEnterpriseEdition),
					string(database.ServiceInstanceEnterpriseEditionHighPerformance),
					string(database.ServiceInstanceEnterpriseEditionExtremePerformance),
				}, true),
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(database.ServiceInstanceLevelBasic),
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceLevelPAAS),
					string(database.ServiceInstanceLevelBasic),
				}, true),
			},
			"shape": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceSubscriptionTypeHourly),
					string(database.ServiceInstanceSubscriptionTypeMonthly),
				}, true),
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_configuration": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_password": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"backup_destination": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(database.ServiceInstanceBackupDestinationNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceBackupDestinationBoth),
								string(database.ServiceInstanceBackupDestinationOSS),
								string(database.ServiceInstanceBackupDestinationNone),
							}, true),
						},
						"character_set": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "AL32UTF8",
						},
						"db_demo": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"disaster_recovery": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"failover_database": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Default:  false,
						},
						"golden_gate": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"is_rac": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"national_character_set": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  database.ServiceInstanceNCharSetUTF16,
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceNCharSetUTF16),
								string(database.ServiceInstanceNCharSetUTF8),
							}, true),
						},
						"pdb_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "pdb1",
						},
						"sid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "ORCL",
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "UTC",
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  database.ServiceInstanceTypeDB,
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceTypeDB),
							}, true),
						},
						"usable_storage": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(15, 2048),
						},
						"snapshot_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source_service_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"instantiate_from_backup": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_storage_container": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cloud_storage_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
							Computed:  true,
						},
						"cloud_storage_username": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"database_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"decryption_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"on_premise": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"service_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"wallet_file_content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"backups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_storage_container": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cloud_storage_username": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"cloud_storage_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"create_if_missing": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
					},
				},
			},
			"hybrid_disaster_recovery": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_storage_container": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"cloud_storage_username": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cloud_storage_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_reservations": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notification_email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"byol": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"high_performance_storage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cloud_storage_container": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"compute_site_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"dbaas_monitor_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"em_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"glassfish_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"identity_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceOPAASDatabaseServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating database service instance")

	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	input := database.CreateServiceInstanceInput{
		Name:                      d.Get("name").(string),
		Edition:                   database.ServiceInstanceEdition(d.Get("edition").(string)),
		IPReservations:            getStringList(d, "ip_reservations"),
		IsBYOL:                    d.Get("byol").(bool),
		Level:                     database.ServiceInstanceLevel(d.Get("level").(string)),
		Shape:                     database.ServiceInstanceShape(d.Get("shape").(string)),
		SubscriptionType:          database.ServiceInstanceSubscriptionType(d.Get("subscription_type").(string)),
		UseHighPerformanceStorage: d.Get("high_performance_storage").(bool),
		Version:                   database.ServiceInstanceVersion(d.Get("version").(string)),
		VMPublicKey:               d.Get("ssh_public_key").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	if v, ok := d.GetOk("notification_email"); ok {
		input.EnableNotification = true
		input.NotificationEmail = v.(string)
	}

	if v, ok := d.GetOk("ip_network"); ok {
		input.IPNetwork = v.(string)
	}

	if _, ok := d.GetOk("ip_reservations"); ok {
		input.IPReservations = getStringList(d, "ip_reservations")
	}

	if v, ok := d.GetOk("region"); ok {
		input.Region = v.(string)
	}

	if v, ok := d.GetOk("availability_domain"); ok {
		input.AvailabilityDomain = v.(string)
	}

	if v, ok := d.GetOk("subnet"); ok {
		input.Subnet = v.(string)
	}

	// Only the PaaS level can have a parameter.
	if input.Level == database.ServiceInstanceLevelPAAS {
		input.Parameter, err = expandParameter(d)
		if err != nil {
			return err
		}
	}

	info, err := client.CreateServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseServiceInstance: %+v", err)
	}

	d.SetId(info.Name)
	return resourceOPAASDatabaseServiceInstanceRead(d, meta)
}

func resourceOPAASDatabaseServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	getInput := database.GetServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetServiceInstance(&getInput)
	if err != nil {
		// DatabaseServiceInstance does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading database service instance %s: %+v", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of database service instance %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("region", result.Region)
	d.Set("availability_domain", result.AvailabilityDomain)
	d.Set("description", result.Description)
	d.Set("backup_destination", result.BackupDestination)
	d.Set("character_set", result.CharSet)
	d.Set("cloud_storage_container", result.CloudStorageContainer)
	d.Set("compute_site_name", result.ComputeSiteName)
	d.Set("connect_descriptor", result.ConnectDescriptor)
	d.Set("dbaas_monitor_url", result.DBAASMonitorURL)
	d.Set("edition", result.Edition)
	d.Set("em_url", result.EMURL)
	d.Set("failover_database", result.FailoverDatabase)
	d.Set("high_performance_storage", result.UseHighPerformanceStorage)
	d.Set("glassfish_url", result.GlassFishURL)
	d.Set("identity_domain", result.IdentityDomain)
	d.Set("ip_network", result.IPNetwork)
	d.Set("byol", result.IsBYOL)
	d.Set("level", result.Level)
	d.Set("national_character_set", result.NCharSet)
	d.Set("pdb_name", result.PDBName)
	d.Set("uri", result.URI)
	d.Set("shape", result.Shape)
	d.Set("sid", result.SID)
	d.Set("subnet", result.Subnet)
	d.Set("subscription_type", result.SubscriptionType)
	d.Set("timezone", result.Timezone)
	d.Set("version", result.Version)

	setAttributesFromConfig(d)

	return nil
}

// Certain values aren't received from the get call and need to be specified from the config
func setAttributesFromConfig(d *schema.ResourceData) {
	d.Set("disaster_recovery", d.Get("disaster_recovery"))

}

func resourceOPAASDatabaseServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	dbClient, err := getDatabaseClient(meta)
	if err != nil {
		return err
	}
	client := dbClient.ServiceInstanceClient()
	name := d.Id()

	log.Printf("[DEBUG] Deleting DatabaseServiceInstance: %v", name)

	input := database.DeleteServiceInstanceInput{
		Name: name,
	}
	if err := client.DeleteServiceInstance(&input); err != nil {
		return fmt.Errorf("Error deleting DatabaseServiceInstance: %+v", err)
	}
	return nil
}

func expandParameter(d *schema.ResourceData) (database.ParameterInput, error) {
	databaseConfigInfo := d.Get("database_configuration").([]interface{})
	attrs := databaseConfigInfo[0].(map[string]interface{})
	parameter := database.ParameterInput{
		AdminPassword:     attrs["admin_password"].(string),
		BackupDestination: database.ServiceInstanceBackupDestination(attrs["backup_destination"].(string)),
		CharSet:           attrs["character_set"].(string),
		DisasterRecovery:  attrs["disaster_recovery"].(bool),
		FailoverDatabase:  attrs["failover_database"].(bool),
		GoldenGate:        attrs["golden_gate"].(bool),
		IsRAC:             attrs["is_rac"].(bool),
		NCharSet:          database.ServiceInstanceNCharSet(attrs["national_character_set"].(string)),
		PDBName:           attrs["pdb_name"].(string),
		SID:               attrs["sid"].(string),
		Timezone:          attrs["timezone"].(string),
		Type:              database.ServiceInstanceType(attrs["type"].(string)),
		UsableStorage:     strconv.Itoa(attrs["usable_storage"].(int)),
	}

	if val, ok := attrs["snapshot_name"].(string); ok && val != "" {
		parameter.SnapshotName = val
	}
	if val, ok := attrs["source_service_name"].(string); ok && val != "" {
		parameter.SourceServiceName = val
	}
	if val, ok := attrs["db_demo"].(string); ok {
		addParam := database.AdditionalParameters{
			DBDemo: val,
		}
		parameter.AdditionalParameters = addParam
	}
	expandIbkup(d, &parameter)
	err := expandBackups(d, &parameter)
	if err != nil {
		return parameter, err
	}
	expandHDG(d, &parameter)

	return parameter, nil
}

func expandIbkup(d *schema.ResourceData, parameter *database.ParameterInput) {
	ibkupInfo := d.Get("instantiate_from_backup").([]interface{})
	if len(ibkupInfo) > 0 {
		attrs := ibkupInfo[0].(map[string]interface{})
		parameter.IBKUP = true
		parameter.IBKUPDatabaseID = attrs["cloud_storage_container"].(string)
		parameter.IBKUPOnPremise = attrs["on_premise"].(bool)
		if val, ok := attrs["cloud_storage_username"]; ok {
			parameter.IBKUPCloudStorageUser = val.(string)
		}
		if val, ok := attrs["cloud_storage_password"]; ok {
			parameter.IBKUPCloudStoragePassword = val.(string)
		}
		if val, ok := attrs["decryption_key"]; ok {
			parameter.IBKUPDecryptionKey = val.(string)
		}
		if val, ok := attrs["service_id"]; ok {
			parameter.IBKUPServiceID = val.(string)
		}
		if val, ok := attrs["wallet_file_content"]; ok {
			parameter.IBKUPWalletFileContent = val.(string)
		}
	}
}

func expandBackups(d *schema.ResourceData, parameter *database.ParameterInput) error {
	cloudStorageInfo := d.Get("backups").([]interface{})

	if parameter.BackupDestination == database.ServiceInstanceBackupDestinationBoth || parameter.BackupDestination == database.ServiceInstanceBackupDestinationOSS {
		if len(cloudStorageInfo) == 0 {
			return fmt.Errorf("`backups` must be set if `backup_destination` is set to `OSS` or `BOTH`")
		}
	}

	if len(cloudStorageInfo) > 0 {
		attrs := cloudStorageInfo[0].(map[string]interface{})
		parameter.CloudStorageContainer = attrs["cloud_storage_container"].(string)
		parameter.CreateStorageContainerIfMissing = attrs["create_if_missing"].(bool)
		if val, ok := attrs["cloud_storage_username"].(string); ok && val != "" {
			parameter.CloudStorageUsername = val
		}
		if val, ok := attrs["cloud_storage_password"].(string); ok && val != "" {
			parameter.CloudStoragePassword = val
		}
	}
	return nil
}

func expandHDG(d *schema.ResourceData, parameter *database.ParameterInput) error {
	hdgInfo := d.Get("hybrid_disaster_recovery").([]interface{})

	if len(hdgInfo) > 0 {
		if parameter.FailoverDatabase == true || parameter.IsRAC == true {
			return fmt.Errorf("`hybrid_disaster_recovery` cannot be set if `is_rac` or `failover_database` is set to true")
		}
		attrs := hdgInfo[0].(map[string]interface{})
		parameter.HDG = true
		parameter.HDGCloudStorageContainer = attrs["cloud_storage_container"].(string)
		// TODO read these values in the sdk like we do with cloud storage
		if val, ok := attrs["cloud_storage_username"].(string); ok && val != "" {
			parameter.HDGCloudStorageUser = val
		}
		if val, ok := attrs["cloud_storage_password"].(string); ok && val != "" {
			parameter.HDGCloudStoragePassword = val
		}
	}

	return nil
}