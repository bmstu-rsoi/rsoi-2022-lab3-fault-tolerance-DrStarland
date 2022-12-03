package services

import (
	"bytes"
	"fmt"
	"gateway/pkg/models/privilege"
	"gateway/pkg/myjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetPrivilegeShortInfo(bonusServiceAddress, username string) (*privilege.Privilege, error) {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/%s", bonusServiceAddress, username)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Println("Failed to create an http request")
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to privilege service: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	res.Body.Close()

	privilege := &privilege.Privilege{}
	if res.StatusCode != http.StatusNotFound {
		if err = myjson.From(body, privilege); err != nil {
			return nil, fmt.Errorf("failed to decode response: %s", err)
		}
	}

	return privilege, nil
}

func GetPrivilegeHistory(bonusServiceAddress string, privilegeID int) (*[]privilege.PrivilegeHistory, error) {
	requestURL := fmt.Sprintf("%s/api/v1/bonushistory/%d", bonusServiceAddress, privilegeID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to privilege service: %s", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	res.Body.Close()

	privilegeHistory := &[]privilege.PrivilegeHistory{}
	if res.StatusCode != http.StatusNotFound {
		if err = myjson.From(body, privilegeHistory); err != nil {
			return nil, fmt.Errorf("failed to decode response: %s", err)
		}
	}

	return privilegeHistory, nil
}

func CreatePrivilegeHistoryRecord(bonusServiceAddress, uid, date, optype string, ID, diff int) error {
	requestURL := fmt.Sprintf("%s/api/v1/bonus", bonusServiceAddress)

	record := &privilege.PrivilegeHistory{
		PrivilegeID:   ID,
		TicketUID:     uid,
		Date:          date,
		BalanceDiff:   diff,
		OperationType: optype,
	}
	data, err := myjson.To(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		log.Println("Failed to create an http request")
		return err
	}

	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to privilege service: %s", err)
	}
	res.Body.Close()

	return nil
}

func CreatePrivilege(bonusServiceAddress, username string, balance int) error {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/privilege", bonusServiceAddress)

	record := &privilege.Privilege{
		Username: username,
		Balance:  balance,
	}
	data, err := myjson.To(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		log.Println("Failed to create an http request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to privilege service: %s", err)
	}
	res.Body.Close()

	return nil
}

func UpdatePrivilege(bonusServiceAddress, username string, balance int) error {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/privilege", bonusServiceAddress)

	priv := &privilege.Privilege{
		Username: username,
		Balance:  balance,
	}
	data, err := myjson.To(priv)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewReader(data))
	if err != nil {
		log.Println("Failed to create an http request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Name", username)

	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to privilege service: %s", err)
	}
	res.Body.Close()

	return nil
}
