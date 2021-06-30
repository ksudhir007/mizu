package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mizuserver/pkg/database"
	"mizuserver/pkg/models"
	"mizuserver/pkg/utils"
	"mizuserver/pkg/validation"
	"reflect"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/martian/har"
	jsonpath "github.com/yalp/jsonpath"
	yaml "gopkg.in/yaml.v2"
)

const (
	OrderDesc = "desc"
	OrderAsc  = "asc"
	LT        = "lt"
	GT        = "gt"
)

var (
	operatorToSymbolMapping = map[string]string{
		LT: "<",
		GT: ">",
	}
	operatorToOrderMapping = map[string]string{
		LT: OrderDesc,
		GT: OrderAsc,
	}
)

func GetEntries(c *fiber.Ctx) error {
	entriesFilter := &models.EntriesFilter{}
	if err := c.QueryParser(entriesFilter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	err := validation.Validate(entriesFilter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	order := operatorToOrderMapping[entriesFilter.Operator]
	operatorSymbol := operatorToSymbolMapping[entriesFilter.Operator]
	var entries []models.MizuEntry
	database.GetEntriesTable().
		Order(fmt.Sprintf("timestamp %s", order)).
		Where(fmt.Sprintf("timestamp %s %v", operatorSymbol, entriesFilter.Timestamp)).
		Omit("entry"). // remove the "big" entry field
		Limit(entriesFilter.Limit).
		Find(&entries)

	if len(entries) > 0 && order == OrderDesc {
		// the entries always order from oldest to newest so we should revers
		utils.ReverseSlice(entries)
	}

	// Convert to base entries
	baseEntries := make([]models.BaseEntryDetails, 0, entriesFilter.Limit)
	fmt.Println(baseEntries)
	for _, entry := range entries {
		baseEntries = append(baseEntries, utils.GetResolvedBaseEntry(entry))
	}
	return c.Status(fiber.StatusOK).JSON(baseEntries)
}

func GetHARs(c *fiber.Ctx) error {
	entriesFilter := &models.HarFetchRequestBody{}
	order := OrderDesc
	if err := c.QueryParser(entriesFilter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	err := validation.Validate(entriesFilter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	var timestampFrom, timestampTo int64

	if entriesFilter.From < 0 {
		timestampFrom = 0
	} else {
		timestampFrom = entriesFilter.From
	}
	if entriesFilter.To <= 0 {
		timestampTo = time.Now().UnixNano() / int64(time.Millisecond)
	} else {
		timestampTo = entriesFilter.To
	}

	var entries []models.MizuEntry
	database.GetEntriesTable().
		Where(fmt.Sprintf("timestamp BETWEEN %v AND %v", timestampFrom, timestampTo)).
		Order(fmt.Sprintf("timestamp %s", order)).
		Find(&entries)

	if len(entries) > 0 {
		// the entries always order from oldest to newest so we should revers
		utils.ReverseSlice(entries)
	}

	harsObject := map[string]*models.ExtendedHAR{}

	for _, entryData := range entries {
		var harEntry har.Entry
		_ = json.Unmarshal([]byte(entryData.Entry), &harEntry)

		sourceOfEntry := entryData.ResolvedSource
		fileName := fmt.Sprintf("%s.har", sourceOfEntry)
		if harOfSource, ok := harsObject[fileName]; ok {
			harOfSource.Log.Entries = append(harOfSource.Log.Entries, &harEntry)
		} else {
			var entriesHar []*har.Entry
			entriesHar = append(entriesHar, &harEntry)
			harsObject[fileName] = &models.ExtendedHAR{
				Log: &models.ExtendedLog{
					Version: "1.2",
					Creator: &models.ExtendedCreator{
						Creator: &har.Creator{
							Name:    "mizu",
							Version: "0.0.2",
						},
						Source: sourceOfEntry,
					},
					Entries: entriesHar,
				},
			}
		}
	}

	retObj := map[string][]byte{}
	for k, v := range harsObject {
		bytesData, _ := json.Marshal(v)
		retObj[k] = bytesData
	}
	buffer := utils.ZipData(retObj)
	return c.Status(fiber.StatusOK).SendStream(buffer)
}

func GetFullEntries(c *fiber.Ctx) error {
	entriesFilter := &models.HarFetchRequestBody{}
	order := OrderDesc
	if err := c.QueryParser(entriesFilter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}
	err := validation.Validate(entriesFilter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	var timestampFrom, timestampTo int64

	if entriesFilter.From < 0 {
		timestampFrom = 0
	} else {
		timestampFrom = entriesFilter.From
	}
	if entriesFilter.To <= 0 {
		timestampTo = time.Now().UnixNano() / int64(time.Millisecond)
	} else {
		timestampTo = entriesFilter.To
	}

	var entries []models.MizuEntry
	database.GetEntriesTable().
		Where(fmt.Sprintf("timestamp BETWEEN %v AND %v", timestampFrom, timestampTo)).
		Order(fmt.Sprintf("timestamp %s", order)).
		Find(&entries)

	if len(entries) > 0 {
		// the entries always order from oldest to newest so we should revers
		utils.ReverseSlice(entries)
	}

	entriesArray := make([]har.Entry, 0)
	for _, entryData := range entries {
		var harEntry har.Entry
		_ = json.Unmarshal([]byte(entryData.Entry), &harEntry)
		entriesArray = append(entriesArray, harEntry)
	}
	return c.Status(fiber.StatusOK).JSON(entriesArray)
}

func GetEntry(c *fiber.Ctx) error {
	var entryData models.EntryData
	database.GetEntriesTable().
		Select("entry", "resolvedDestination").
		Where(map[string]string{"entryId": c.Params("entryId")}).
		First(&entryData)

	var fullEntry har.Entry
	unmarshallErr := json.Unmarshal([]byte(entryData.Entry), &fullEntry)
	utils.CheckErr(unmarshallErr)
	resultPolicyToSend := matchRequestPolicy(fullEntry)
	if entryData.ResolvedDestination != "" {
		fullEntry.Request.URL = utils.SetHostname(fullEntry.Request.URL, entryData.ResolvedDestination)
	}
	var fullEntryWithPolicy models.FullEntryWithPolicy
	fullEntryWithPolicy.RulesMatched = resultPolicyToSend
	fullEntryWithPolicy.Entry = fullEntry
	return c.Status(fiber.StatusOK).JSON(fullEntryWithPolicy)
}

func matchRequestPolicy(fullEntry har.Entry) []models.RulesMatched {
	enforcePolicy, _ := decodeEnforcePolicy()
	var resultPolicyToSend []models.RulesMatched
	for _, value := range enforcePolicy.Rules {
		if value.Type == "json" {
			var bodyJsonMap interface{}
			_ = json.Unmarshal(fullEntry.Response.Content.Text, &bodyJsonMap)
			fmt.Println("217 ", bodyJsonMap)
			var result models.RulesMatched
			out, err := jsonpath.Read(bodyJsonMap, value.Key)
			if err != nil {
				fmt.Println("220 ", err)
				continue
			}
			fmt.Println("222 ", out, " ", reflect.TypeOf(out))
			var matchValue bool
			if reflect.TypeOf(out).Kind() == reflect.String {
				matchValue, err = regexp.MatchString(value.Value, out.(string))
			} else {
				val := fmt.Sprint(out)
				matchValue, err = regexp.MatchString(value.Value, val)
			}
			fmt.Println("225 ", matchValue)
			if matchValue {
				fmt.Printf("%s matched with value %v", value.Name, value.Value)
				result.Matched = true
				result.Rule = value
				resultPolicyToSend = append(resultPolicyToSend, result)
			} else {
				result.Matched = false
				result.Rule = value
				resultPolicyToSend = append(resultPolicyToSend, result)
			}
		} else if value.Type == "header" {
			for j := range fullEntry.Response.Headers {
				matchKey, _ := regexp.MatchString(value.Key, fullEntry.Response.Headers[j].Name)
				if matchKey {
					matchValue, _ := regexp.MatchString(value.Value, fullEntry.Response.Headers[j].Value)
					var result models.RulesMatched
					if matchValue {
						result.Matched = true
						result.Rule = value
						resultPolicyToSend = append(resultPolicyToSend, result)
					} else {
						result.Matched = true
						result.Rule = value
						resultPolicyToSend = append(resultPolicyToSend, result)
					}
				}
			}
		} else {

		}
	}
	return resultPolicyToSend
}

func decodeEnforcePolicy() (models.RulesPolicy, error) {
	content, err := ioutil.ReadFile("/app/enforce-policy/enforce-policy.yaml")
	enforcePolicy := models.RulesPolicy{}
	if err != nil {
		return enforcePolicy, err
	}
	err = yaml.Unmarshal([]byte(content), &enforcePolicy)
	if err != nil {
		return enforcePolicy, err
	}
	invalidIndex := enforcePolicy.ValidateRulesPolicy()
	if len(invalidIndex) != 0 {
		for i := range invalidIndex {
			fmt.Println("only json and header types are supported on rule")
			enforcePolicy.RemoveNotValidPolicy(invalidIndex[i])
		}
	}
	return enforcePolicy, nil
}

func DeleteAllEntries(c *fiber.Ctx) error {
	database.GetEntriesTable().
		Where("1 = 1").
		Delete(&models.MizuEntry{})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "Success",
	})

}

func GetGeneralStats(c *fiber.Ctx) error {
	sqlQuery := "SELECT count(*) as count, min(timestamp) as min, max(timestamp) as max from mizu_entries"
	var result struct {
		Count int
		Min   int
		Max   int
	}
	database.GetEntriesTable().Raw(sqlQuery).Scan(&result)
	return c.Status(fiber.StatusOK).JSON(&result)
}
