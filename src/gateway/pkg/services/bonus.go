package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/pkg/models/privilege"
	"gateway/pkg/myjson"
	"io"
	"net/http"
	"time"
)

func GetPrivilegeShortInfo(bonusServiceAddress, username string) (*privilege.Privilege, error) {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/%s", bonusServiceAddress, username)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to privilege service: %s", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	privilege := &privilege.Privilege{}
	if res.StatusCode != http.StatusNotFound {
		if err = json.NewDecoder(res.Body).Decode(privilege); err != nil {
			return nil, fmt.Errorf("Failed to decode response: %s", err)
		}
	}

	return privilege, nil
}

func GetPrivilegeHistory(bonusServiceAddress string, privilegeID int) (*[]privilege.PrivilegeHistory, error) {
	requestURL := fmt.Sprintf("%s/api/v1/bonushistory/%d", bonusServiceAddress, privilegeID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("Failed to create an http request")
		return nil, err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to privilege service: %s", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			fmt.Println("Failed to close response body")
		}
	}(res.Body)

	privilegeHistory := &[]privilege.PrivilegeHistory{}
	if res.StatusCode != http.StatusNotFound {
		if err = json.NewDecoder(res.Body).Decode(privilegeHistory); err != nil {
			return nil, fmt.Errorf("Failed to decode response: %s", err)
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
		return fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Failed to create an http request")
		return err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed request to privilege service: %s", err)
	}

	return nil
}

func CreatePrivilege(bonusServiceAddress, username string, balance int) error {
	requestURL := fmt.Sprintf("%s/api/v1/bonus/privilege", bonusServiceAddress)

	record := &privilege.Privilege{
		Username: username,
		Balance:  balance,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Failed to create an http request")
		return err
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed request to privilege service: %s", err)
	}
	res.Body.Close()

	return nil
}
