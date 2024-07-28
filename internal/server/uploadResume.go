package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"recruit-sys/internal/models"
	"strings"
)

func (s *Server) HandleUploadResume(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value(UserID).(float64)
	if !ok {
		return NewAPIError(http.StatusBadGateway, fmt.Errorf("unable to proceed"))
	}

	user, err := s.db.SelectUserWhereID(userID)
	if err != nil {
		return err
	}

	if err := r.ParseMultipartForm(2 << 20); err != nil {
		return err
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return NewAPIError(400, err)
	}
	
	if filepath.Ext(handler.Filename) != ".pdf" && filepath.Ext(handler.Filename) != ".docx" {
		return NewAPIError(http.StatusBadRequest,fmt.Errorf("only pdf and docx files accepted"))
	}
	defer file.Close()

	destDir := "/home/shaurya/code/git/recruit-sys/uploads/"
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}

	userIdString := fmt.Sprintf("%f", userID)
	uniqueFileName := userIdString + filepath.Ext(handler.Filename)
	filePath := filepath.Join(destDir, uniqueFileName)

	dest, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		return err
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	profile, err := postResumeFile(fileContent)
	if err != nil {
		return err
	}

	profile.User = user
	profile.UserID = int(userID)
	profile.ResumeFileAddress = filePath

	if err = s.db.CreateProfile(&profile); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
	return nil
}

// function to post resume to 3rd party API
func postResumeFile(file []byte) (models.Profile, error) {
	apiURL := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(file))
	if err != nil {
		return models.Profile{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apikey", apiKey)

	client := &http.Client{}
	fmt.Println("profile")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return models.Profile{}, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return models.Profile{}, err
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", responseBody)

	profile, err := BytesToProfile(responseBody)
	if err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}

// To convert API response to valid Profile format
func BytesToProfile(response []byte) (models.Profile, error) {
	var respProfile ResponseProfile

	if err := json.Unmarshal(response, &respProfile); err != nil {
		fmt.Println(err)
		return models.Profile{}, err
	}

	education := ""
	for _, edu := range respProfile.Education {
		if education == "" {
			education = edu.Name
		} else {
			education = education + ", " + edu.Name
		}
	}

	experience := ""
	for _, exp := range respProfile.Experience {
		expStr := fmt.Sprintf("%s at %s, %s (%s - %s)", exp.Title, exp.Organization, exp.Location, exp.DateStart, exp.DateEnd)
		if experience == "" {
			experience = expStr
		} else {
			experience += "; " + expStr
		}
	}

	profile := models.Profile{
		Skills:     strings.Join(respProfile.Skills, ", "),
		Experience: experience,
		Education:  education,
		Name:       respProfile.Name,
		Email:      respProfile.Email,
		Phone:      respProfile.Phone,
	}

	return profile, nil
}

type ResponseProfile struct {
	Skills    []string `json:"skills"`
	Education []struct {
		Name string `json:"name"`
	} `json:"education"`
	Experience []struct {
		Title        string   `json:"title"`
		Dates        []string `json:"dates"`
		DateStart    string   `json:"date_start"`
		DateEnd      string   `json:"date_end"`
		Location     string   `json:"location"`
		Organization string   `json:"organization"`
	} `json:"experience"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
