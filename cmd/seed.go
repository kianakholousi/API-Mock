package cmd

import (
	"flight-data-api/config"
	"flight-data-api/database"
	"flight-data-api/models"
	"github.com/spf13/cobra"
	"math/rand"
	"syreclabs.com/go/faker"
	"time"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database",
	Long: `This command seeds the database with initial data.

You must specify a custom configuration file in YAML format using the --config flag. By default, this command will not run without a configuration file.

It is recommended to run this command after running the migrations to ensure that database is available.
	
Usage:
seed --config [path]
	
Flags:
	-c, --config string   Path to custom configuration file in YAML format (required)
	-h, --help            help for seed`,
	Run: func(cmd *cobra.Command, args []string) {
		seedDB()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to custom configuration file in YAML format")
	err := seedCmd.MarkFlagRequired("config")
	if err != nil {
		panic(err)
	}
}

func seedDB() {
	cfg, err := config.Init(config.Params{FilePath: configPath, FileType: "yaml"})
	if err != nil {
		panic(err)
	}

	db, err := database.InitDB(cfg.Database)
	if err != nil {
		panic(err)
	}

	cityNames := []string{"London", "Paris", "Tokyo", "Seoul", "Sydney", "Oslo", "Berlin", "Rome", "Madrid", "Athens"}
	cities := make([]models.City, 0, 10)
	for i := 0; i < 10; i++ {
		cities = append(cities, models.City{
			Name: cityNames[i],
		})
	}
	if err := db.Create(cities).Error; err != nil {
		panic(err)
	}

	airplaneNames := []string{"AirbusA320", "Boeing737", "EmbraerE190", "BombardierCRJ900", "SukhoiSSJ100", "Fokker100", "ATR72", "DeHavillandDash8", "MitsubishiMRJ", "ComacARJ21"}
	airplanes := make([]models.Airplane, 0, 10)
	for i := 0; i < 10; i++ {
		airplanes = append(airplanes, models.Airplane{
			Name: airplaneNames[i],
		})
	}
	if err := db.Create(airplanes).Error; err != nil {
		panic(err)
	}

	cxlSits := make([]models.CancelingSituation, 0, 2)
	cxlSits = append(cxlSits, models.CancelingSituation{
		Description: "des1",
		Data:        "data1",
	})
	cxlSits = append(cxlSits, models.CancelingSituation{
		Description: "des2",
		Data:        "data2",
	})
	if err := db.Create(cxlSits).Error; err != nil {
		panic(err)
	}

	depTimes := make([]time.Time, 0, 20)
	arrTimes := make([]time.Time, 0, 20)
	depCityInd := make([]int, 0, 20)
	arrCityInd := make([]int, 0, 20)
	for i := 0; i < 20; i++ {
		depTimes = append(depTimes, faker.Date().Birthday(2, 3))
		arrTimes = append(arrTimes, depTimes[i].
			Add(time.Duration(rand.Intn(3)+1)*time.Hour).
			Add(time.Duration(rand.Intn(60))*time.Minute))
		depCityInd = append(depCityInd, rand.Intn(10))
		arrCityInd = append(arrCityInd, (depCityInd[i]+rand.Intn(9)+1)%10)
	}

	airlineNames := []string{"Emirates", "QatarAirways", "EtihadAirways", "TurkishAirlines", "Lufthansa"}
	flights := make([]models.Flight, 0, 20)
	for i := 0; i < 20; i++ {
		flights = append(flights, models.Flight{
			DepCity:        cities[depCityInd[i]],
			ArrCity:        cities[arrCityInd[i]],
			DepTime:        depTimes[i],
			ArrTime:        arrTimes[i],
			Airplane:       airplanes[rand.Intn(10)],
			Airline:        airlineNames[rand.Intn(5)],
			Price:          int32(rand.Intn(1000) + 500),
			CxlSit:         cxlSits[rand.Intn(2)],
			RemainingSeats: int32(rand.Intn(100) + 1),
		})
	}
	if err := db.Create(flights).Error; err != nil {
		panic(err)
	}
}
