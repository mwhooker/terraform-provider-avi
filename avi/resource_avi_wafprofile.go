/*
 * Copyright (c) 2017. Avi Networks.
 * Author: Gaurav Rastogi (grastogi@avinetworks.com)
 *
 */
package avi

import (
	"github.com/avinetworks/sdk/go/clients"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func ResourceWafProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"config": &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     ResourceWafConfigSchema(),
			Set: func(v interface{}) int {
				return 0
			},
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"files": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     ResourceWafDataFileSchema(),
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"tenant_ref": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"uuid": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func resourceAviWafProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAviWafProfileCreate,
		Read:   ResourceAviWafProfileRead,
		Update: resourceAviWafProfileUpdate,
		Delete: resourceAviWafProfileDelete,
		Schema: ResourceWafProfileSchema(),
	}
}

func ResourceAviWafProfileRead(d *schema.ResourceData, meta interface{}) error {
	s := ResourceWafProfileSchema()
	client := meta.(*clients.AviClient)
	var obj interface{}
	if uuid, ok := d.GetOk("uuid"); ok {
		path := "api/wafprofile/" + uuid.(string)
		err := client.AviSession.Get(path, &obj)
		if err != nil {
			d.SetId("")
			return nil
		}
	} else {
		d.SetId("")
		return nil
	}
	if _, err := ApiDataToSchema(obj, d, s); err == nil {
		if err != nil {
			log.Printf("[ERROR] in setting read object %v\n", err)
		}
	}
	return nil
}

func resourceAviWafProfileCreate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceWafProfileSchema()
	err := ApiCreateOrUpdate(d, meta, "wafprofile", s)
	if err == nil {
		err = ResourceAviWafProfileRead(d, meta)
	}
	return err
}

func resourceAviWafProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	s := ResourceWafProfileSchema()
	err := ApiCreateOrUpdate(d, meta, "wafprofile", s)
	if err == nil {
		err = ResourceAviWafProfileRead(d, meta)
	}
	return err
}

func resourceAviWafProfileDelete(d *schema.ResourceData, meta interface{}) error {
	objType := "wafprofile"
	client := meta.(*clients.AviClient)
	uuid := d.Get("uuid").(string)
	if uuid != "" {
		path := "api/" + objType + "/" + uuid
		err := client.AviSession.Delete(path)
		if err != nil && !(strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "204")) {
			log.Println("[INFO] resourceAviWafProfileDelete not found")
			return err
		}
		d.SetId("")
	}
	return nil
}