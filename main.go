package main

import (
	"city-search-project/dao"
	"city-search-project/modelPojo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const admin = "admin"
const posters = "posters"

var conCity = dao.CityDAO{}
var conService = dao.ServiceDAO{}
var conCategory = dao.CategoryDAO{}

func addCityDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token := r.Header.Get("tokenid")

	isAdmin := token == admin
	isPoster := token == posters

	if !(isAdmin || isPoster) {
		respondWithError(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var city modelPojo.City

	if err := json.NewDecoder(r.Body).Decode(&city); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if err := conCity.Insert(city); err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable To Insert Record")
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": " Record Inserted Successfully",
		})
	}
}

func getCityByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}

	cityName := strings.Split(r.URL.Path, "/")[2]

	city, err := conCity.FindByCityName(cityName)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, city)
}

func deleteCityByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
		return
	}
	var reqBody map[string]string

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	city := reqBody["city_name"]

	err := conCity.DeleteCity(city)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record deleted successfully",
	})
}

func updateCityByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}
	var city modelPojo.City
	err := json.NewDecoder(r.Body).Decode(&city)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	cityName := city.CityName

	_ = conCity.UpdateCity(cityName, city)

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record updated successfully",
	})
}

func addServiceDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token := r.Header.Get("tokenid")

	isAdmin := token == admin
	isPoster := token == posters

	if !(isAdmin || isPoster) {
		respondWithError(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var service modelPojo.Service

	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if err := conService.Insert(service); err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable To Insert Record")
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": " Record Inserted Successfully",
		})
	}
}

func getServiceByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}

	name := strings.Split(r.URL.Path, "/")[2]

	service, err := conService.FindByServiceName(name)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, service)
}

func getServiceByCityAndCategory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}

	name := strings.Split(r.URL.Path, "/")[2]

	service, err := conService.FindByServiceName(name)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, service)
}

func deleteServiceByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
		return
	}
	var reqBody map[string]string

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	name := reqBody["name"]

	err := conService.DeleteService(name)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record deleted successfully",
	})
}

func updateServiceByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}
	var service modelPojo.Service
	err := json.NewDecoder(r.Body).Decode(&service)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	name := service.Name

	_ = conService.UpdateService(name, service)

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record updated successfully",
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func addCategoryDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token := r.Header.Get("tokenid")

	isAdmin := token == admin
	isPoster := token == posters

	if !(isAdmin || isPoster) {
		respondWithError(w, http.StatusBadRequest, "Unauthorized")
		return
	}

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var category modelPojo.Classification

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	if err := conCategory.Insert(category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable To Insert Record")
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": " Record Inserted Successfully",
		})
	}
}

func getCategoryByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}

	serviceType := strings.Split(r.URL.Path, "/")[2]

	category, err := conCategory.FindByCategory(serviceType)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, category)
}

func deleteCategoryByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
		return
	}
	var reqBody map[string]string

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	category := reqBody["service_type"]

	err := conCategory.DeleteCategory(category)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record deleted successfully",
	})
}

func updateCategoryByName(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}
	var category modelPojo.Classification
	err := json.NewDecoder(r.Body).Decode(&category)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	serviceType := category.ServiceType

	_ = conCategory.UpdateCategory(serviceType, category)

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record updated successfully",
	})
}

func getServiceByCategoryAndCityInExcel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var search modelPojo.Search

	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	// searchData, err := conService.FindByCategoryAndCity(search)
	searchData, _, err := conService.FindByCategoryAndCity(search, "Excel")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, searchData)

}

func getServiceByCategoryAndCityInPdf(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
		return
	}

	var search modelPojo.Search

	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	// searchData, err := conService.FindByCategoryAndCity(search)
	searchData, _, err := conService.FindByCategoryAndCity(search, "Pdf")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, searchData)

}

func main() {
	http.HandleFunc("/service-add/", addServiceDetails)
	http.HandleFunc("/service-delete/", deleteServiceByName)
	http.HandleFunc("/service-edit/", updateServiceByName)
	http.HandleFunc("/service-search/", getServiceByName)

	http.HandleFunc("/city-add/", addCityDetails)
	http.HandleFunc("/city-delete/", deleteCityByName)
	http.HandleFunc("/city-edit/", updateCityByName)
	http.HandleFunc("/city-search/", getCityByName)

	http.HandleFunc("/category-add/", addCategoryDetails)
	http.HandleFunc("/category-delete/", deleteCategoryByName)
	http.HandleFunc("/category-edit/", updateCategoryByName)
	http.HandleFunc("/category-search/", getCategoryByName)

	http.HandleFunc("/category-city-search-excel/", getServiceByCategoryAndCityInExcel)
	http.HandleFunc("/category-city-search-pdf/", getServiceByCategoryAndCityInPdf)

	fmt.Println(" Main Method Excecuted ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
