package seeds

import (
	"flight-data-api/models"
	"gorm.io/gorm"
	"math/rand"
	"syreclabs.com/go/faker"
	"time"
)

const (
	cnt = 10
)

type Seed struct {
	DB *gorm.DB
}

func Init(s Seed) error {
	cityNames := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose"}
	cities := make([]models.City, 0)
	for i := 0; i < cnt; i++ {
		cities = append(cities, models.City{
			Name: cityNames[i],
		})
	}
	if err := s.DB.Create(cities).Error; err != nil {
		return err
	}

	airplaneNames := []string{"Boeing 747", "Airbus A380", "Boeing 737", "Airbus A320", "Boeing 777", "Embraer E190", "Bombardier CRJ200", "Boeing 787", "Airbus A350", "Embraer E195"}
	airplanes := make([]models.Airplane, 0)
	for i := 0; i < cnt; i++ {
		airplanes = append(airplanes, models.Airplane{
			Name: airplaneNames[i],
		})
	}
	if err := s.DB.Create(airplanes).Error; err != nil {
		return err
	}

	cxlSits := make([]models.CancelingSituation, 0)
	cxlSits = append(cxlSits, models.CancelingSituation{
		Description: "des1",
		Data:        "data1",
	})
	cxlSits = append(cxlSits, models.CancelingSituation{
		Description: "des2",
		Data:        "data2",
	})
	if err := s.DB.Create(cxlSits).Error; err != nil {
		return err
	}

	airlineNames := []string{"American Airlines", "Delta Air Lines", "United Airlines", "Southwest Airlines", "JetBlue Airways", "Alaska Airlines", "Frontier Airlines", "Spirit Airlines", "Hawaiian Airlines", "Allegiant Air"}
	depTimes := make([]time.Time, 0)
	arrTimes := make([]time.Time, 0)
	for i := 0; i < cnt; i++ {
		depTimes = append(depTimes, faker.Date().Birthday(2, 3))
		arrTimes = append(arrTimes, depTimes[i].
			Add(time.Duration(rand.Intn(3)+1)*time.Hour).
			Add(time.Duration(rand.Intn(60))*time.Minute))
	}

	flights := make([]models.Flight, 0)
	for i := 0; i < cnt; i++ {
		flights = append(flights, models.Flight{
			DepCity:  cities[rand.Intn(10)],
			ArrCity:  cities[rand.Intn(10)],
			DepTime:  depTimes[i],
			ArrTime:  arrTimes[i],
			Airplane: airplanes[rand.Intn(10)],
			Airline:  airlineNames[i],
			Price:    int32(rand.Intn(1000) + 500),
			CxlSit:   cxlSits[rand.Intn(2)],
			LeftSeat: int32(rand.Intn(100) + 1),
		})
	}
	if err := s.DB.Create(flights).Error; err != nil {
		return err
	}

	return nil
}
