package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

var (
	p              *message.Printer = message.NewPrinter(language.English)
	nonDigitRegexp *regexp.Regexp   = regexp.MustCompile(`[^0-9\.]`)
)

func main() {
	hourlyRate := "60"
	var timeToCraft string

	var leatherUsageInFeet string
	var leatherPricePerFeet string
	var hardware string
	var transportation string
	var thread string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Leather Usage In Feet").Value(&leatherUsageInFeet),
			huh.NewInput().Title("Leather £E/feet^2").Value(&leatherPricePerFeet),
			huh.NewInput().Title("Hardware Cost").Value(&hardware),
			huh.NewInput().Title("Transportation Cost").Value(&transportation),
			huh.NewInput().Title("Thread Cost").Value(&thread),
		),
		huh.NewGroup(
			huh.NewInput().Title("Time to craft").Value(&timeToCraft),
			huh.NewInput().Title("Hourly Rate").Value(&hourlyRate),
		),
	)

	err := form.Run()
	if err != nil {
		l := logger.Get()
		l.Fatal("error running form: %v", zap.Error(err))
	}

	v := product.Variant{
		TimeToCraft:   floatToDuration(parseFloat(timeToCraft)),
		MaterialsCost: parseFloat(leatherUsageInFeet)*parseFloat(leatherPricePerFeet) + parseFloat(hardware) + parseFloat(thread) + parseFloat(transportation),
	}

	fmt.Printf("Estimated Price: £E %s\n", formatNum(v.EstPrice(parseFloat(hourlyRate))))
	fmt.Printf("Estimated Wholesale Price: £E %s\n", formatNum(v.EstWholesalePrice(parseFloat(hourlyRate))))
}

func floatToDuration(f float64) time.Duration {
	d := time.Duration(f * float64(time.Hour))
	if d < 0 {
		log.Fatalf("negative duration: %f", f)
	}
	return d
}

func formatNum[T ~int | ~int64 | ~int32 | ~int16 | ~int8 | ~float64 | ~float32](num T) string {
	return p.Sprintf("%v", number.Decimal(num))
}

func parseFloat(str string) float64 {
	if str == "" {
		return 0
	}
	str = nonDigitRegexp.ReplaceAllLiteralString(str, "")
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatalln("error parsing float", err)
	}
	return num
}
