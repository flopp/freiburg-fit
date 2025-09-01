package app

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/flopp/freiburg-fit/internal/utils"
	"gopkg.in/yaml.v3"
)

type Venue struct {
	Name     string
	Location struct {
		Address     string
		Coordinates utils.LatLon
	}
	Description template.HTML
	Links       []struct {
		Title string
		URL   string
	}
	Category *Label
	Tags     []*Label
}

func (v Venue) Slug() string {
	return fmt.Sprintf("venue/%s", utils.SanitizeName(v.Name))
}

type Data struct {
	Venues        []*Venue
	Categories    []*Label
	CategoriesMap map[string]*Label
	Tags          []*Label
	TagsMap       map[string]*Label
}

func (d *Data) LoadCategory(yamlFileName string) error {
	type categoryData struct {
		ID          string
		Name        string
		Names       string
		Description template.HTML
	}

	file, err := os.Open(yamlFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var cd categoryData
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cd); err != nil {
		return err
	}

	category := &Label{
		ID:          utils.SanitizeName(cd.ID),
		Name:        cd.Name,
		Names:       cd.Names,
		Type:        "category",
		Description: cd.Description,
		Venues:      []*Venue{},
	}
	d.Categories = append(d.Categories, category)
	d.CategoriesMap[category.ID] = category

	return nil
}

func (d *Data) LoadVenue(yamlFileName string) error {
	type venueData struct {
		Name     string
		Location struct {
			Address     string
			Coordinates string
		}
		Description template.HTML
		Links       []struct {
			Title string
			URL   string
		} `yaml:"links"`
		Category string   `yaml:"category"`
		Tags     []string `yaml:"tags"`
	}

	file, err := os.Open(yamlFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var vd venueData
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&vd); err != nil {
		return err
	}

	venue := &Venue{
		Name:        vd.Name,
		Description: vd.Description,
		Links:       vd.Links,
	}
	d.Venues = append(d.Venues, venue)

	// location
	coords, err := utils.ParseLatLon(vd.Location.Coordinates)
	if err != nil {
		log.Printf("error parsing coordinates for venue %s: %v", vd.Name, err)
	}
	venue.Location.Address = vd.Location.Address
	venue.Location.Coordinates = coords

	// category
	categorySanitized := utils.SanitizeName(vd.Category)
	category, ok := d.CategoriesMap[categorySanitized]
	if !ok {
		category = &Label{
			ID:     categorySanitized,
			Name:   vd.Category,
			Names:  vd.Category,
			Type:   "category",
			Venues: []*Venue{},
		}
		d.Categories = append(d.Categories, category)
		d.CategoriesMap[category.ID] = category
	}
	venue.Category = category
	venue.Category.Venues = append(venue.Category.Venues, venue)

	// tags
	for _, tag := range vd.Tags {
		tagSanitized := utils.SanitizeName(tag)
		t, ok := d.TagsMap[tagSanitized]
		if !ok {
			t = &Label{
				ID:     tagSanitized,
				Name:   tag,
				Names:  tag,
				Type:   "tag",
				Venues: []*Venue{},
			}
			d.Tags = append(d.Tags, t)
			d.TagsMap[tagSanitized] = t
		}
		venue.Tags = append(venue.Tags, t)
		t.Venues = append(t.Venues, venue)
	}
	sort.Slice(venue.Tags, func(i, j int) bool {
		return venue.Tags[i].Slug() < venue.Tags[j].Slug()
	})

	return nil
}

func LoadData(dataFolder string) (*Data, error) {
	data := &Data{
		Categories:    []*Label{},
		CategoriesMap: make(map[string]*Label),
		Tags:          []*Label{},
		TagsMap:       make(map[string]*Label),
	}

	categoryFiles, err := ioutil.ReadDir(dataFolder + "/categories")
	if err != nil {
		return nil, err
	}
	for _, file := range categoryFiles {
		if filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		// Load categories from the YAML file
		fileName := filepath.Join(dataFolder, "categories", file.Name())
		if err := data.LoadCategory(fileName); err != nil {
			return nil, fmt.Errorf("error loading category from %s: %w", fileName, err)
		}
	}

	venueFiles, err := ioutil.ReadDir(dataFolder + "/venues")
	if err != nil {
		return nil, err
	}

	for _, file := range venueFiles {
		if filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		// Load venues from the YAML file
		fileName := filepath.Join(dataFolder, "venues", file.Name())
		if err := data.LoadVenue(fileName); err != nil {
			return nil, fmt.Errorf("error loading venue from %s: %w", fileName, err)
		}
	}

	// sort categories
	sort.Slice(data.Categories, func(i, j int) bool {
		return data.Categories[i].Slug() < data.Categories[j].Slug()
	})
	for _, category := range data.Categories {
		sort.Slice(category.Venues, func(i, j int) bool {
			return category.Venues[i].Slug() < category.Venues[j].Slug()
		})
	}

	// sort tags
	sort.Slice(data.Tags, func(i, j int) bool {
		return data.Tags[i].Slug() < data.Tags[j].Slug()
	})
	for _, tag := range data.Tags {
		sort.Slice(tag.Venues, func(i, j int) bool {
			return tag.Venues[i].Slug() < tag.Venues[j].Slug()
		})
	}

	return data, nil
}
