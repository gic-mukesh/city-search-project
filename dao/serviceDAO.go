package dao

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"city-search-project/modelPojo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/xuri/excelize/v2"
)

var conService = ServiceDAO{}

func init() {
	conService.Server = "mongodb://localhost:27017/"
	conService.Database = "CityDB"
	conService.Collection = "Services"

	conService.Connect()

}

type ServiceDAO struct {
	Server     string
	Database   string
	Collection string
}

var CollectionService *mongo.Collection
var ctxService = context.TODO()

func (e *ServiceDAO) Connect() {
	clientOptions := options.Client().ApplyURI(e.Server)
	client, err := mongo.Connect(ctxService, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctxService, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = license.SetMeteredKey("7bf0687e6377d9e27406c5c6b8f26ad0fcb55cdd74e22275da7c7466cbc3d04f")
	CollectionService = client.Database(e.Database).Collection(e.Collection)
}

func (e *ServiceDAO) Insert(service modelPojo.Service) error {
	_, err := CollectionService.InsertOne(ctxService, service)

	if err != nil {
		return errors.New("unable to create new record")
	}

	return nil
}

func (e *ServiceDAO) FindByServiceName(name string) ([]*modelPojo.Service, error) {
	var service []*modelPojo.Service

	cur, err := CollectionService.Find(ctxService, bson.D{primitive.E{Key: "name", Value: name}})

	if err != nil {
		return service, errors.New("unable to query db")
	}

	for cur.Next(ctxService) {
		var e modelPojo.Service

		err := cur.Decode(&e)

		if err != nil {
			return service, err
		}

		service = append(service, &e)
	}

	if err := cur.Err(); err != nil {
		return service, err
	}

	cur.Close(ctxService)

	if len(service) == 0 {
		return service, mongo.ErrNoDocuments
	}

	return service, nil
}

func (e *ServiceDAO) DeleteService(name string) error {
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	res, err := CollectionService.DeleteOne(ctxService, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no category deleted")
	}

	return nil
}

func (epd *ServiceDAO) UpdateService(name string, service modelPojo.Service) error {
	filter := bson.D{primitive.E{Key: "name", Value: name}}

	update := bson.D{primitive.E{Key: "$set", Value: service}}

	e := &modelPojo.Service{}
	return CollectionService.FindOneAndUpdate(ctxService, filter, update).Decode(e)
}

func (e *ServiceDAO) FindByCategoryAndCity(search modelPojo.Search, option string) ([]*modelPojo.Service, string, error) {
	var service []*modelPojo.Service
	// var dataty []byte
	var cur *mongo.Cursor
	var err error
	str := "Given values are not valid "
	// option := "Pdf"

	os.MkdirAll("data/download", os.ModePerm)
	dir := "data/download/"
	file := "ServiceSearch" + fmt.Sprintf("%v", time.Now().Format("2006-01-02_3_4_5_pm"))
	if (search.CityName != "") && (search.ServiceType != "") {

		cur, err = CollectionService.Find(ctx, bson.D{primitive.E{Key: "city.city_name", Value: search.CityName}, primitive.E{Key: "classification.service_type", Value: search.ServiceType}})
	} else if search.CityName != "" {
		cur, err = CollectionService.Find(ctx, bson.D{primitive.E{Key: "city.city_name", Value: search.CityName}})
	} else if search.ServiceType != "" {
		cur, err = CollectionService.Find(ctx, bson.D{primitive.E{Key: "classification.service_type", Value: search.ServiceType}})

	}
	if err != nil {
		return service, file, errors.New("unable to query db")
	}

	for cur.Next(ctx) {
		var e modelPojo.Service
		err := cur.Decode(&e)
		if err != nil {
			return service, file, err
		}
		service = append(service, &e)
	}

	if service == nil {
		return service, file, errors.New(str)
	}

	if option == "Excel" {
		log.Println("Excel")
		errExcel := writeDataIntoExcel(dir, file, service)
		if errExcel != nil {
			return service, file, err
		}
		_, err = ioutil.ReadFile(dir + file + ".xlsx")
		if err != nil {
			return service, file, err
		}
	}

	if option == "Pdf" {
		log.Println("Pdf")
		_, errPdf := writeDataIntoPDFTable(dir, file, service)
		if errPdf != nil {
			fmt.Println(errPdf)
			return service, file, err
		}
		_, err2 := ioutil.ReadFile(dir + file + ".pdf")
		// fmt.Println(dataty)
		// fmt.Println("Data length", len(dataty))
		if err2 != nil {
			return service, file, err
		}
	}

	return service, file, nil
}

func writeDataIntoExcel(dir, file string, service []*modelPojo.Service) error {

	f := excelize.NewFile()
	f.SetSheetName("Sheet1", "SearchData")

	f.SetCellValue("SearchData", "A1", "ID")
	f.SetCellValue("SearchData", "C1", "Name")
	f.SetCellValue("SearchData", "D1", "Address")
	f.SetCellValue("SearchData", "E1", "Latitude")
	f.SetCellValue("SearchData", "F1", "Longitude")
	f.SetCellValue("SearchData", "G1", "Website")
	f.SetCellValue("SearchData", "H1", "ContactNumber")
	f.SetCellValue("SearchData", "J1", "City")
	f.SetCellValue("SearchData", "N1", "ServiceType")

	for i := range service {
		f.SetCellValue("SearchData", "A"+fmt.Sprintf("%v", i+2), service[i].ID)
		f.SetCellValue("SearchData", "C"+fmt.Sprintf("%v", i+2), service[i].Name)
		f.SetCellValue("SearchData", "D"+fmt.Sprintf("%v", i+2), service[i].Address)
		f.SetCellValue("SearchData", "E"+fmt.Sprintf("%v", i+2), service[i].Latitude)
		f.SetCellValue("SearchData", "F"+fmt.Sprintf("%v", i+2), service[i].Longitude)
		f.SetCellValue("SearchData", "G"+fmt.Sprintf("%v", i+2), service[i].Website)
		f.SetCellValue("SearchData", "H"+fmt.Sprintf("%v", i+2), service[i].ContactNumber)
		f.SetCellValue("SearchData", "J"+fmt.Sprintf("%v", i+2), service[i].City.CityName)
		f.SetCellValue("SearchData", "N"+fmt.Sprintf("%v", i+2), service[i].Classification.ServiceType)
	}

	if err := f.SaveAs(dir + file + ".xlsx"); err != nil {
		return err
	}
	return nil
}

func writeDataIntoPDFTable(dir, file string, service []*modelPojo.Service) (*creator.Creator, error) {

	c := creator.New()
	c.SetPageMargins(20, 20, 20, 20)

	// Create report fonts.
	// UniPDF supports a number of font-families, which can be accessed using model.
	// Here we are creating two fonts, a normal one and its bold version
	font, err := model.NewStandard14Font(model.HelveticaName)
	if err != nil {
		return c, err
	}

	// Bold font
	fontBold, err := model.NewStandard14Font(model.HelveticaBoldName)
	if err != nil {
		return c, err
	}

	// Generate basic usage chapter.
	if err := basicUsage(c, font, fontBold, service); err != nil {
		return c, err
	}

	err = c.WriteToFile(dir + file + ".pdf")
	if err != nil {
		return c, err
	}
	return c, nil
}

func basicUsage(c *creator.Creator, font, fontBold *model.PdfFont, service []*modelPojo.Service) error {
	// Create chapter.
	ch := c.NewChapter("Search Data")
	ch.SetMargins(0, 0, 10, 0)
	ch.GetHeading().SetFont(font)
	ch.GetHeading().SetFontSize(18)
	ch.GetHeading().SetColor(creator.ColorRGBFrom8bit(72, 86, 95))
	// You can also set inbuilt colors using creator
	// ch.GetHeading().SetColor(creator.ColorBlack)

	// Draw subchapters. Here we are only create horizontally aligned chapter.
	// You can also vertically align and perform other optimizations as well.
	// Check GitHub example for more.
	contentAlignH(c, ch, font, fontBold, service)

	// Draw chapter.
	if err := c.Draw(ch); err != nil {
		return err
	}

	return nil
}

func contentAlignH(c *creator.Creator, ch *creator.Chapter, font, fontBold *model.PdfFont, service []*modelPojo.Service) {
	// Create subchapter.
	// sc := ch.NewSubchapter("Content horizontal alignment")
	// sc.GetHeading().SetFontSize(10)
	// sc.GetHeading().SetColor(creator.ColorBlue)

	// Create table.
	table := c.NewTable(9)
	table.SetMargins(0, 0, 15, 0)

	drawCell := func(text string, font *model.PdfFont, align creator.CellHorizontalAlignment) {
		p := c.NewStyledParagraph()
		p.Append(text).Style.Font = font

		cell := table.NewCell()
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetHorizontalAlignment(align)
		cell.SetContent(p)
	}
	// Draw table header.
	drawCell("ID", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell("Name", fontBold, creator.CellHorizontalAlignmentRight)
	drawCell("Address", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell("Latitude", fontBold, creator.CellHorizontalAlignmentRight)
	drawCell("Longitude", fontBold, creator.CellHorizontalAlignmentLeft)
	drawCell("Website", fontBold, creator.CellHorizontalAlignmentCenter)
	drawCell("ContactNumber", fontBold, creator.CellHorizontalAlignmentRight)
	drawCell("City", fontBold, creator.CellHorizontalAlignmentCenter)
	drawCell("ServiceType", fontBold, creator.CellHorizontalAlignmentRight)

	// Draw table content.
	for i := range service {

		drawCell(fmt.Sprintf("%v", service[i].ID), font, creator.CellHorizontalAlignmentLeft)
		drawCell(service[i].Name, font, creator.CellHorizontalAlignmentCenter)
		drawCell(service[i].Address, font, creator.CellHorizontalAlignmentCenter)
		drawCell(fmt.Sprintf("%v", service[i].Latitude), font, creator.CellHorizontalAlignmentCenter)
		drawCell(fmt.Sprintf("%v", service[i].Longitude), font, creator.CellHorizontalAlignmentCenter)
		drawCell(service[i].Website, font, creator.CellHorizontalAlignmentCenter)
		drawCell(fmt.Sprintf("%v", service[i].ContactNumber), font, creator.CellHorizontalAlignmentCenter)
		drawCell(service[i].City.CityName, font, creator.CellHorizontalAlignmentCenter)
		drawCell(service[i].Classification.ServiceType, font, creator.CellHorizontalAlignmentCenter)
	}

	ch.Add(table)
}
