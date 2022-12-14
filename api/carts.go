package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"
	"strings"
)

func (api *API) AddCart(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(model.ErrorResponse{Error: "http: named cookie not present"})
            return
        }
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal server error"})
        return
	}

	product := r.FormValue("product")
	if product == "" {
		w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Request Product Not Found"})
        return
	}

	var list []model.Product
	var totalPrice float64
	for _, formList := range r.Form {
		for _, v := range formList {
			item := strings.Split(v, ",")
			p, _ := strconv.ParseFloat(item[2], 64)
			q, _ := strconv.ParseFloat(item[3], 64)
			total := p * q
			list = append(list, model.Product{
				Id:       item[0],
				Name:     item[1],
				Price:    item[2],
				Quantity: item[3],
				Total:    total,
			})
			totalPrice += total
		}
	}

	username := fmt.Sprintf("%s", r.Context().Value("username"))

	cart := model.Cart{
		Name:       username,
		Cart:       list,
		TotalPrice: totalPrice,
	}


	_, err = api.cartsRepo.CartUserExist(cart.Name)
	if err != nil {
		api.cartsRepo.AddCart(cart)
	} else {
		api.cartsRepo.UpdateCart(cart)
	}
	api.dashboardView(w, r)
}
