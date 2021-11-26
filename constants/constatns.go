package constants

import (
	"errors"
	"os"

	"github.com/magiconair/properties"
)

var Config = map[string]string{}

func init() {

	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		f, err := os.Create("./config.properties")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		p := properties.NewProperties()
		p.SetValue("serverAddr", "localhost:3000")
		p.SetValue("cname", "controller-A")
		p.Write(f, properties.UTF8)
		Config["serverAddr"] = "localhost:3000"
		Config["cname"] = "controller-A"
		return
	}

	p := properties.MustLoadFile("./config.properties", properties.UTF8)
	Config["serverAddr"] = p.GetString("serverAddr", "localhost:3000")
	Config["cname"] = p.GetString("cname", "controller-A")
	Config["cid"] = p.GetString("cid", "")
}

func Set(key, value string) {

	var p *properties.Properties
	if _, err := os.Stat("./config.properties"); errors.Is(err, os.ErrNotExist) {
		p = properties.NewProperties()
	} else {
		p = properties.MustLoadFile("./config.properties", properties.UTF8)
		os.Remove("./config.properties")
	}
	f, err := os.Create("./config.properties")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p.SetValue(key, value)
	p.Write(f, properties.UTF8)
	Config[key] = value
}
