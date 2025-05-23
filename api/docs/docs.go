// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/tasks": {
            "get": {
                "description": "Get all tasks with optional sorting",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Get all tasks",
                "parameters": [
                    {
                        "enum": [
                            "CreateAsc",
                            "CreateDesc",
                            "PriorityAsc",
                            "PriorityDesc",
                            "DeadlineAsc",
                            "DeadlineDesc"
                        ],
                        "type": "string",
                        "description": "Sorting",
                        "name": "sorting",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Task"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "post": {
                "description": "Create a new task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Create a task",
                "parameters": [
                    {
                        "description": "Task",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/DTOs.CreateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/errors.ApplicationError"
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/tasks/{id}": {
            "put": {
                "description": "Update task",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Update task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Task",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/DTOs.UpdateTaskRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/DTOs.TaskResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/errors.ApplicationError"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/errors.ApplicationError"
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            },
            "delete": {
                "description": "Delete task by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Delete task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/errors.ApplicationError"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/errors.ApplicationError"
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/tasks/{id}/toggle": {
            "patch": {
                "description": "Change task's status",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Toggle task's status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Task",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/DTOs.ToggleTaskStatusRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/DTOs.TaskResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/errors.ApplicationError"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/errors.ApplicationError"
                        }
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "DTOs.CreateTaskRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "deadline": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "enum": [
                        "Low",
                        "Medium",
                        "High",
                        "Critical"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/enums.Priority"
                        }
                    ]
                }
            }
        },
        "DTOs.TaskResponse": {
            "type": "object",
            "required": [
                "createdAt",
                "id",
                "name",
                "priority",
                "status"
            ],
            "properties": {
                "changedAt": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "deadline": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "$ref": "#/definitions/enums.Priority"
                },
                "status": {
                    "$ref": "#/definitions/enums.Status"
                }
            }
        },
        "DTOs.ToggleTaskStatusRequest": {
            "type": "object",
            "required": [
                "isDone"
            ],
            "properties": {
                "isDone": {
                    "type": "boolean"
                }
            }
        },
        "DTOs.UpdateTaskRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "deadline": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "enum": [
                        "Low",
                        "Medium",
                        "High",
                        "Critical"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/enums.Priority"
                        }
                    ]
                }
            }
        },
        "enums.Priority": {
            "type": "string",
            "enum": [
                "Low",
                "Medium",
                "High",
                "Critical"
            ],
            "x-enum-varnames": [
                "Low",
                "Medium",
                "High",
                "Critical"
            ]
        },
        "enums.Status": {
            "type": "string",
            "enum": [
                "Active",
                "Completed",
                "Overdue",
                "Late"
            ],
            "x-enum-varnames": [
                "Active",
                "Completed",
                "Overdue",
                "Late"
            ]
        },
        "errors.ApplicationError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "errors": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "statusCode": {
                    "type": "integer"
                }
            }
        },
        "models.Task": {
            "type": "object",
            "properties": {
                "changedAt": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "deadline": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "$ref": "#/definitions/enums.Priority"
                },
                "status": {
                    "$ref": "#/definitions/enums.Status"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
