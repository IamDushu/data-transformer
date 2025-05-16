package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bearbin/go-age"
)

type DonorProfile struct {
	Answers map[string]AnswerBlock `json:"answers"`
	User    User                   `json:"user"`
	Photos  []Photo                `json:"photos"`
	Program string                 `json:"program"`
}

type Photo struct {
	CroppedSource string `json:"cropped_source"`
}

type User struct {
	ID           string       `json:"id"`
	DonorCode    string       `json:"donorCode"`
	DateOfBirth  string       `json:"dateOfBirth"`
	FreezeMember FreezeMember `json:"freeze_member"`
}

type FreezeMember struct {
	ProfileBio string `json:"profile_bio"`
}

type AnswerBlock struct {
	Answer   AnswerValue `json:"answer"`
	Question Question    `json:"question"`
}

type AnswerValue struct {
	Value interface{} `json:"value"`
}

type Question struct {
	Label string `json:"label"`
}

type OutputRecord struct {
	UserID              string `json:"user"`
	UserImage           string `json:"user_image"`
	Text                string `json:"text"`
	JobTitle            string `json:"job_title"`
	ArtisticAbility     int    `json:"artistic_ability"`
	AthleticAbility     int    `json:"athletic_ability"`
	MathematicalAbility int    `json:"mathematical_ability"`
	ScientificAbility   int    `json:"scientific_ability"`
	SingingAbility      int    `json:"singing_ability"`
	HairType            string `json:"hair_type"`
	HairColor           string `json:"hair_color"`
	EducationLevel      string `json:"education_level"`
	JewishAncestry      string `json:"jewish_ancestry"`
	HeightFt            int    `json:"height_ft"`
	HeightIn            int    `json:"height_in"`
	LogicalCreative     string `json:"logical_creative"`
	SeriousSilly        string `json:"serious_silly"`
	IntrovertExtrovert  string `json:"introvert_extrovert"`
	RelationshipPrefs   string `json:"relationship_preferences"`
	Passions            string `json:"passions"`
	GoalsInLife         string `json:"goals_in_life"`
	GreatestStrengths   string `json:"greatest_strengths"`
	PerfectDay          string `json:"perfect_day"`
	DinnerParty         string `json:"dinner_party"`
	Motivation          string `json:"motivation"`
	MessageToIPs        string `json:"message_to_ips"`
	Book                string `json:"book"`
	Movie               string `json:"movie"`
	Food                string `json:"food"`
	Allergies           string `json:"allergies"`
	DentalWork          string `json:"dental_work"`
	Dimples             string `json:"dimples"`
	EggRetrieval        string `json:"egg_retrieval"`
	Freckles            string `json:"freckles"`
	Siblings            string `json:"siblings"`
	Complexion          string `json:"complexion"`
	Diet                string `json:"diet"`
	DominantHand        string `json:"dominant_hand"`
	EyeColor            string `json:"eye_color"`
	HairTexture         string `json:"hair_texture"`
	MaritalStatus       string `json:"marital_status"`
	VisionQuality       string `json:"vision_quality"`
	Weight              int    `json:"weight"`
	DonorCode           string `json:"donorCode"`
	ProfileBio          string `json:"profile_bio"`
	Age                 int    `json:"age"`
}

func main() {
	file, err := os.Open("donorprofiles.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var donors []DonorProfile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&donors)
	if err != nil {
		panic(err)
	}

	// Explicitly ordered fields
	fieldsToExtract := []string{
		"passions",
		"goals_in_life",
		"greatest_strengths",
		"perfect_day",
		"dinner_party",
		"motivation",
		"message_to_ips",
		"book",
		"movie",
		// -----------
		"food",
	}

	var outputRecords []OutputRecord

	for _, donor := range donors {
		textParts := []string{}

		// Add context intro
		textParts = append(textParts, "the following is a question and answer profile of an egg donor:")

		// Add selected fields in order
		for _, key := range fieldsToExtract {
			answer, exists := donor.Answers[key]
			if !exists {
				continue
			}

			value := answer.Answer.Value
			if value == nil || fmt.Sprintf("%v", value) == "" {
				continue
			}

			// Default label
			label := strings.TrimSpace(answer.Question.Label)

			// Custom label overrides
			switch key {
			case "book":
				label = "What are your favorite books"
			case "movie":
				label = "What is your favorite movie"
			case "food":
				label = "What is your favorite food"
			}

			textLine := fmt.Sprintf("%s: %v.", label, value)
			textParts = append(textParts, textLine)
		}

		// Add profile bio if present
		if strings.TrimSpace(donor.User.FreezeMember.ProfileBio) != "" {
			textParts = append(textParts, fmt.Sprintf("Donor Information: %s.", donor.User.FreezeMember.ProfileBio))
		}

		// Combine all text parts
		text := strings.Join(textParts, " ")

		// Clean text (lowercase, replace underscores)
		cleanText := strings.ToLower(strings.ReplaceAll(text, "_", " "))

		birthDate, _ := time.Parse(time.RFC3339, donor.User.DateOfBirth)
		age := age.Age(birthDate)

		record := OutputRecord{
			UserID:              donor.User.ID,
			UserImage:           donor.Photos[0].CroppedSource,
			Text:                cleanText,
			DonorCode:           donor.User.DonorCode,
			ProfileBio:          donor.User.FreezeMember.ProfileBio,
			Age:                 age,
			JobTitle:            getStringAnswer(donor, "job_title"),
			ArtisticAbility:     getIntAnswer(donor, "artistic_ability"),
			AthleticAbility:     getIntAnswer(donor, "athletic_ability"),
			MathematicalAbility: getIntAnswer(donor, "mathematical_ability"),
			ScientificAbility:   getIntAnswer(donor, "scientific_ability"),
			SingingAbility:      getIntAnswer(donor, "singing_ability"),
			HairType:            getStringAnswer(donor, "hair_type"),
			HairColor:           getStringAnswer(donor, "hair_color"),
			EducationLevel:      getStringAnswer(donor, "education_level"),
			JewishAncestry:      getStringAnswer(donor, "jewish_ancestry"),
			HeightFt:            getIntAnswer(donor, "height_ft"),
			HeightIn:            getIntAnswer(donor, "height_in"),
			LogicalCreative:     getStringAnswer(donor, "logical_creative"),
			SeriousSilly:        getStringAnswer(donor, "serious_silly"),
			IntrovertExtrovert:  getStringAnswer(donor, "introvert_extrovert"),
			RelationshipPrefs:   getStringAnswer(donor, "relationship_preferences"),
			GoalsInLife:         getStringAnswer(donor, "goals_in_life"),
			Passions:            getStringAnswer(donor, "passions"),
			GreatestStrengths:   getStringAnswer(donor, "greatest_strengths"),
			PerfectDay:          getStringAnswer(donor, "perfect_day"),
			DinnerParty:         getStringAnswer(donor, "dinner_party"),
			Motivation:          getStringAnswer(donor, "motivation"),
			MessageToIPs:        getStringAnswer(donor, "message_to_ips"),
			Book:                getStringAnswer(donor, "book"),
			Movie:               getStringAnswer(donor, "movie"),
			Food:                getStringAnswer(donor, "food"),
			Allergies:           getStringAnswer(donor, "allergies"),
			DentalWork:          getStringAnswer(donor, "dental_work"),
			Dimples:             getStringAnswer(donor, "dimples"),
			EggRetrieval:        getStringAnswer(donor, "egg_retrieval"),
			Freckles:            getStringAnswer(donor, "freckles"),
			Siblings:            getStringAnswer(donor, "siblings"),
			Complexion:          getStringAnswer(donor, "complexion"),
			Diet:                getStringAnswer(donor, "diet"),
			DominantHand:        getStringAnswer(donor, "dominant_hand"),
			EyeColor:            getStringAnswer(donor, "eye_color"),
			HairTexture:         getStringAnswer(donor, "hair_texture"),
			MaritalStatus:       getStringAnswer(donor, "marital_status"),
			VisionQuality:       getStringAnswer(donor, "vision_quality"),
			Weight:              getIntAnswer(donor, "weight"),
		}

		outputRecords = append(outputRecords, record)
	}

	// Write to JSON file
	outfile, err := os.Create("a2_profiles.json")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	encoder := json.NewEncoder(outfile)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(outputRecords)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Cleaned donor profiles written to a2_profiles.json")
}

func getIntAnswer(donor DonorProfile, key string) int {
	answer, exists := donor.Answers[key]
	if !exists || answer.Answer.Value == nil {
		return 0
	}
	if val, ok := answer.Answer.Value.(float64); ok {
		return int(val)
	}
	return 0
}

func getStringAnswer(donor DonorProfile, key string) string {
	answer, exists := donor.Answers[key]
	if !exists || answer.Answer.Value == nil {
		return ""
	}
	if val, ok := answer.Answer.Value.(string); ok {
		return val
	}
	return ""
}
