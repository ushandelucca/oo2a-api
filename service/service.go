// Package service implements the management and the reporting use case.
package service

import (
	"MeasurementWeb/database"
	"MeasurementWeb/utils"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

type Unit int

const (
	Percent Unit = iota + 1
	Lph
)

func (u Unit) String() string {
	return [...]string{"", "Percent", "Lph"}[u]
}

// Measurement is the transport object for the measurement values.
type MeasurementDto struct {
	ID        uint
	Timestamp time.Time `validate:"required"`
	Sensor    string    `validate:"required"`
	Value     float64   `validate:"required"`
	Unit      Unit      `validate:"required"`
}

// MeasurementDto is the transport object for the measurement values.
type Oo2aModel struct {
	ValuesBegin time.Time
	ValuesEnd   time.Time

	LevelCurrent float64
	LevelValues  []float64
	LevelUnit    Unit

	PrecipitationCurrent float64
	PrecipitationValues  []float64
	PrecipitationUnit    Unit
}

type ManagementService interface {
	InitMeasurement() (err error)
	SaveMeasurement(m database.MeasurementDo) (entity database.MeasurementDo, err error)
}

// exported 'constructor'
func NewManagementService(config *utils.Conf) *managementService {
	db, err := database.NewMeasurementDB(config)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &managementService{db: db}
}

type managementService struct {
	db database.MeasurementDB
}

func (s *managementService) InitMeasurement() (err error) {
	err = s.db.SetupMeasurements()
	return err
}

func (s *managementService) SaveMeasurement(transferObject MeasurementDto) (err error) {
	validate := validator.New()
	err = validate.Struct(transferObject)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err)
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return
	}

	dataObject := database.MeasurementDo{}
	dataObject.ID = transferObject.ID
	dataObject.Timestamp = transferObject.Timestamp
	dataObject.Sensor = transferObject.Sensor
	dataObject.Value = transferObject.Value
	dataObject.Unit = transferObject.Unit.String()

	// maps only from do to dto --> use only for read
	// mapper := dto.Mapper{}
	// mapper.AddConvFunc(func(u Unit, mapper *dto.Mapper) string {
	// 	return fmt.Sprint(u)
	// })
	// err = mapper.Map(&dataObject, transferObject)

	if dataObject.ID == 0 {
		_, err = s.db.CreateMeasurement(dataObject)
	} else {
		_, err = s.db.UpdateMeasurement(dataObject)
	}

	if err != nil {
		return err
	}

	return nil
}
