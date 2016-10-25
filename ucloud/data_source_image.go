package ucloud

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceImageRead,

		Schema: map[string]*schema.Schema{
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"os_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"image_name_regexp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"most_recent": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"image_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"image_size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"features": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

type imageSort []*client.UHostImage

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	return a[i].CreateTime < a[j].CreateTime
}

// Returns the most recent AMI out of a slice of images.
func mostRecentImage(images []*client.UHostImage) *client.UHostImage {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

func setImageMeta(d *schema.ResourceData, image *client.UHostImage) {
	d.SetId(image.ImageId)
	d.Set("image_id", image.ImageId)
	d.Set("zone", image.Zone)
	d.Set("image_name", image.ImageName)
	d.Set("image_type", image.ImageType)
	d.Set("image_size", image.ImageSize)
	d.Set("os_type", image.OsType)
	d.Set("os_name", image.OsName)
	d.Set("state", image.State)
	d.Set("image_description", image.ImageDescription)
	d.Set("create_time", image.CreateTime)
	d.Set("features", image.Features)
}

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.Client)

	var nameRegexp *regexp.Regexp
	if v, ok := d.GetOk("image_name_regexp"); ok {
		nameRegexp = regexp.MustCompile(v.(string))
	}

	params := client.DescribeImageRequest{Limit: 200}
	if v, ok := d.GetOk("zone"); ok {
		params.Zone = v.(string)
	}
	if v, ok := d.GetOk("image_type"); ok {
		params.ImageType = v.(string)
	}
	if v, ok := d.GetOk("os_type"); ok {
		params.OsType = v.(string)
	}
	if v, ok := d.GetOk("image_id"); ok {
		params.ImageId = v.(string)
	}

	var filteredImages []*client.UHostImage
	var resp client.DescribeImageResponse

	for done := true; done; done = len(resp.ImageSet) == 0 {
		err := apiClient.Call(&params, &resp)
		if err != nil {
			return err
		}

		if nameRegexp != nil {
			for _, image := range resp.ImageSet {
				if nameRegexp.MatchString(image.ImageName) {
					filteredImages = append(filteredImages, image)
				}
			}
		} else {
			filteredImages = append(filteredImages, resp.ImageSet...)
		}

		params.Offset += params.Limit
	}

	var image *client.UHostImage
	if len(filteredImages) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(filteredImages) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] ucloud_image - multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			image = mostRecentImage(filteredImages)
		} else {
			return fmt.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		image = filteredImages[0]
	}

	log.Printf("[DEBUG] ucloud_image - Single Image found: %s", image.ImageId)
	setImageMeta(d, image)

	return nil
}
