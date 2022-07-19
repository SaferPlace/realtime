package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"

	"safer.place/internal/address"
	"safer.place/internal/address/eircodesaferplace"
	"safer.place/internal/address/roughprefix"
	"safer.place/internal/language"
	"safer.place/internal/score"
	scorerv1 "safer.place/internal/score/v1"
	"safer.place/internal/stations"
	"safer.place/internal/web"
	"safer.place/internal/webserver"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
}

var templateFuncs = template.FuncMap{
	"html": func(value string) template.HTML {
		return template.HTML(value)
	},
	"times": func(n int) []int {
		return make([]int, n)
	},
	"subtract": func(a, b int) int {
		return a - b
	},
}

func run() error {
	var cfg Config
	if err := envconfig.Process("saferplace", &cfg); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	log.Printf("config: %+v", cfg)

	var opts []webserver.Option

	langs, err := language.Languages()
	if err != nil {
		return fmt.Errorf("unable to load languages: %w", err)
	}
	opts = append(opts, webserver.Languages(langs))

	// For now we just want something, we don't care what
	var addrResolvers []address.Resolver
	for _, r := range cfg.AddressResolvers.Order {
		switch r {
		case "roughprefix":
			addrResolvers = append(addrResolvers, roughprefix.New())
		case "eircode":
			addrResolvers = append(addrResolvers, eircodesaferplace.New(
				cfg.AddressResolvers.EircodeAddr,
				cfg.AddressResolvers.EircodeToken,
			))
		}
	}

	opts = append(opts, webserver.AddressResolvers(addrResolvers...))

	// Parse the templates
	tmpl := template.Must(template.New("").
		Funcs(templateFuncs).
		ParseFS(web.Templates, "**.html"),
	)
	opts = append(opts, webserver.Templates(tmpl))

	var scorer score.Scorer
	switch cfg.Scorer {
	case "v1":
		scorer = scorerv1.New(stations.New())
	}
	opts = append(opts, webserver.Scorer(scorer))

	ws := webserver.New(opts...)
	if err := ws.Run(cfg.Port); err != nil {
		return fmt.Errorf("webserver failure: %w", err)
	}

	return nil
}

type Config struct {
	Port             int              `envconfig:"PORT" default:"8080"`
	AddressResolvers AddressResolvers `split_words:"true"`
	Scorer           string           `default:"v1"`

	Font      string
	FancyFont string `split_words:"true"`
}

type AddressResolvers struct {
	EircodeAddr  string   `split_words:"true" default:"https://eircode.safer.place"`
	EircodeToken string   `split_words:"true"`
	Order        []string `default:"roughprefix"`
}
