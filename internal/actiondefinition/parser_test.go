package actiondefinition

import (
	"reflect"
	"slices"
	"testing"
)

func TestParser(t *testing.T) {
	testConfig := `{
		"uid": "test_action",
		"type": "workflow-action",
		"config": {
			"actionUrl": "https://webhook.site/c5db800b-e3c3-4a16-af81-1e19f6d8b6a0",
			"isPublished": false,
			"executionRules": [
				{
					"labelName": "success",
					"conditions": {
						"status": ["success"]
					}
				},
				{
					"labelName": "failure",
					"conditions": {
						"status": ["failure"]
					}
				}
			],
			"functions": [],
			"inputFieldDependencies": [],
			"inputFields": [
				{
					"isRequired": false,
					"supportedValueTypes": ["STATIC_VALUE"],
					"typeDefinition": {
						"fieldType": "text",
						"name": "label",
						"options": [],
						"type": "string"
					}
				},
				{
					"isRequired": true,
					"supportedValueTypes": ["STATIC_VALUE"],
					"typeDefinition": {
						"fieldType": "text",
						"name": "value",
						"options": [],
						"type": "string"
					}
				}
			],
			"labels": {
				"en": {
					"actionName": "Test Action",
					"executionRules": {
						"success": "Action Succeeded!",
						"failure": "Action Failed."
					},
					"inputFieldLabels": {
						"label": "Label",
						"value": "Value"
					},
					"inputFieldDescriptions": {
						"label": "",
						"value": ""
					},
					"inputFieldOptionLabels": {
						"label": {},
						"value": {}
					},
					"outputFieldLabels": {
						"status": "Status"
					}
				}
			},
			"functions": [
				{
					"functionType": "PRE_ACTION_EXECUTION",
					"functionSource": "exports.main = function(event, callback) { return callback(transformRequest(event)); }\nfunction transformRequest(request) { return { webhookUrl: request.webhookUrl, body: JSON.stringify(request.fields), contentType: 'application/x-www-form-urlencoded', accept: 'application/json', httpMethod: 'POST' }; }"
				},
				{
					"functionType": "PRE_FETCH_OPTIONS",
					"id": "widgetSize",
					"functionSource": "exports.main = function(event, callback) { return callback(transformRequest(event)); }\nfunction transformRequest(request) { return { webhookUrl: request.webhookUrl + '?color=' + request.fields.widgetColor.value, body: JSON.stringify(request.fields), httpMethod: 'GET' }; }"
				}
			],
			"objectTypes": ["CONTACT"],
			"outputFields": [
				{
					"typeDefinition": {
						"externalOptions": false,
						"options": [
							{
								"displayOrder": 1,
								"label": "Success",
								"value": "success"
							},
							{
								"displayOrder": 2,
								"label": "Failure",
								"value": "failure"
							}
						],
						"type": "enumeration",
						"name": "status",
						"label": "Status"
					}
				}
			],
			"supportedClients": [{ "client": "WORKFLOWS" }]
		}
	}`

	expectedInputFields := []InputField{
		{
			TypeDefinition: StringTypeDefinition{
				Name:      "label",
				Type:      "string",
				FieldType: "text",
			},
			SupportedValueTypes: []string{"STATIC_VALUE"},
			IsRequired:          false,
		},
		{
			TypeDefinition: StringTypeDefinition{
				Name:      "value",
				Type:      "string",
				FieldType: "text",
			},
			SupportedValueTypes: []string{"STATIC_VALUE"},
			IsRequired:          true,
		},
	}

	expectedFunctions := []Function{
		ActionFunction{
			FunctionType:   PreActionExecution,
			FunctionSource: "exports.main = function(event, callback) { return callback(transformRequest(event)); }\nfunction transformRequest(request) { return { webhookUrl: request.webhookUrl, body: JSON.stringify(request.fields), contentType: 'application/x-www-form-urlencoded', accept: 'application/json', httpMethod: 'POST' }; }",
		},
		OptionFunction{
			FunctionType:   PreFetchOptions,
			FunctionSource: "exports.main = function(event, callback) { return callback(transformRequest(event)); }\nfunction transformRequest(request) { return { webhookUrl: request.webhookUrl + '?color=' + request.fields.widgetColor.value, body: JSON.stringify(request.fields), httpMethod: 'GET' }; }",
			Id:             "widgetSize",
		},
	}

	expectedExecutionRules := []ExecutionRule{
		{
			LabelName: "success",
			Conditions: map[string]any{
				"status": []any{"success"},
			},
		},
		{
			LabelName: "failure",
			Conditions: map[string]any{
				"status": []any{"failure"},
			},
		},
	}

	definition, err := Parse([]byte(testConfig))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !slices.EqualFunc(definition.Config.InputFields, expectedInputFields, func(e1, e2 InputField) bool {
		if e1.IsRequired != e2.IsRequired {
			return false
		}

		if e1.TypeDefinition != e2.TypeDefinition {
			return false
		}

		return slices.Equal(e1.SupportedValueTypes, e2.SupportedValueTypes)
	}) {
		t.Errorf("[ActionDef]: Inputfields: Expected %v, got %v", expectedInputFields, definition.Config.InputFields)
	}
	if !slices.Equal(definition.Config.Functions, expectedFunctions) {
		t.Errorf("[ActionDef]: Functions: Expected %v, got %v", expectedFunctions, definition.Config.Functions)
	}
	if !slices.EqualFunc(definition.Config.ExecutionRules, expectedExecutionRules, func(r1, r2 ExecutionRule) bool {
		eqConditions := reflect.DeepEqual(r1.Conditions, r2.Conditions)
		return r1.LabelName == r2.LabelName && eqConditions
	}) {
		t.Errorf("[ActionDef]: ExecutionRules: Expected %v, got %v", expectedExecutionRules, definition.Config.ExecutionRules)
	}
}
