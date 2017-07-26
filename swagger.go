package dlframework

const (
	dlframework_swagger = `{
  "swagger": "2.0",
  "info": {
    "title": "CarML DLFramework",
    "version": "1.0.0",
    "description": "CarML (Cognitive ARtifacts for Machine Learning) is a framework allowing people to develop and deploy machine learning models. It allows machine learning (ML) developers to publish and evaluate their models, users to experiment with different models and frameworks through a web user interface or a REST api, and system architects to capture system resource usage to inform future system and hardware configuration."
  },
  "schemes": [
    "http",
    "https"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/v1/framework/{framework_name}/info": {
      "get": {
        "operationId": "GetFrameworkManifest",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkFrameworkManifest"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "framework_version",
            "in": "query",
            "required": false,
            "type": "string"
          }
        ],
        "tags": [
          "Registry"
        ]
      }
    },
    "/v1/framework/{framework_name}/model/{model_name}/info": {
      "post": {
        "operationId": "GetFrameworkModelManifest",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkModelManifest"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "model_name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkGetFrameworkModelManifestRequest"
            }
          }
        ],
        "tags": [
          "Registry"
        ]
      }
    },
    "/v1/framework/{framework_name}/model/{model_name}/predict": {
      "post": {
        "operationId": "Predict",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkPredictResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "model_name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkPredictRequest"
            }
          }
        ],
        "tags": [
          "Predictor"
        ]
      }
    },
    "/v1/framework/{framework_name}/models": {
      "get": {
        "operationId": "GetFrameworkModels",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkGetModelManifestsResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "framework_version",
            "in": "query",
            "required": false,
            "type": "string"
          }
        ],
        "tags": [
          "Registry"
        ]
      }
    },
    "/v1/frameworks": {
      "get": {
        "operationId": "GetFrameworkManifests",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkGetFrameworkManifestsResponse"
            }
          }
        },
        "tags": [
          "Registry"
        ]
      }
    },
    "/v1/model/{model_name}/info": {
      "post": {
        "operationId": "GetModelManifest",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkModelManifest"
            }
          }
        },
        "parameters": [
          {
            "name": "model_name",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkGetModelManifestRequest"
            }
          }
        ],
        "tags": [
          "Registry"
        ]
      }
    },
    "/v1/models": {
      "get": {
        "operationId": "GetModelManifests",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkGetModelManifestsResponse"
            }
          }
        },
        "tags": [
          "Registry"
        ]
      }
    }
  },
  "definitions": {
    "ModelManifestModel": {
      "type": "object",
      "properties": {
        "base_url": {
          "type": "string"
        },
        "weights_path": {
          "type": "string"
        },
        "graph_path": {
          "type": "string"
        },
        "is_archive": {
          "type": "boolean",
          "format": "boolean"
        }
      }
    },
    "TypeParameter": {
      "type": "object",
      "properties": {
        "value": {
          "type": "string"
        }
      }
    },
    "dlframeworkContainerHardware": {
      "type": "object",
      "properties": {
        "gpu": {
          "type": "string"
        },
        "cpu": {
          "type": "string"
        }
      }
    },
    "dlframeworkErrorStatus": {
      "type": "object",
      "properties": {
        "ok": {
          "type": "boolean",
          "format": "boolean"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "dlframeworkFrameworkManifest": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "version": {
          "type": "string"
        },
        "container": {
          "type": "object",
          "additionalProperties": {
            "$ref": "#/definitions/dlframeworkContainerHardware"
          }
        }
      }
    },
    "dlframeworkGetFrameworkManifestRequest": {
      "type": "object",
      "properties": {
        "framework_name": {
          "type": "string"
        },
        "framework_version": {
          "type": "string"
        }
      }
    },
    "dlframeworkGetFrameworkManifestsResponse": {
      "type": "object",
      "properties": {
        "manifests": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkFrameworkManifest"
          }
        }
      }
    },
    "dlframeworkGetFrameworkModelManifestRequest": {
      "type": "object",
      "properties": {
        "framework_name": {
          "type": "string"
        },
        "framework_version": {
          "type": "string"
        },
        "model_name": {
          "type": "string"
        },
        "model_version": {
          "type": "string"
        }
      }
    },
    "dlframeworkGetModelManifestRequest": {
      "type": "object",
      "properties": {
        "model_name": {
          "type": "string"
        },
        "model_version": {
          "type": "string"
        }
      }
    },
    "dlframeworkGetModelManifestsResponse": {
      "type": "object",
      "properties": {
        "manifests": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkModelManifest"
          }
        }
      }
    },
    "dlframeworkModelManifest": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "version": {
          "type": "string"
        },
        "framework": {
          "$ref": "#/definitions/dlframeworkFrameworkManifest"
        },
        "container": {
          "type": "object",
          "additionalProperties": {
            "$ref": "#/definitions/dlframeworkContainerHardware"
          }
        },
        "description": {
          "type": "string"
        },
        "reference": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "license": {
          "type": "string"
        },
        "inputs": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkModelManifestType"
          }
        },
        "output": {
          "$ref": "#/definitions/dlframeworkModelManifestType"
        },
        "before_preprocess": {
          "type": "string"
        },
        "preprocess": {
          "type": "string"
        },
        "after_preprocess": {
          "type": "string"
        },
        "before_postprocess": {
          "type": "string"
        },
        "postprocess": {
          "type": "string"
        },
        "after_postprocess": {
          "type": "string"
        },
        "model": {
          "$ref": "#/definitions/ModelManifestModel"
        },
        "attributes": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "dlframeworkModelManifestType": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string"
        },
        "description": {
          "type": "string"
        },
        "parameters": {
          "type": "object",
          "additionalProperties": {
            "$ref": "#/definitions/TypeParameter"
          }
        }
      }
    },
    "dlframeworkNull": {
      "type": "object"
    },
    "dlframeworkPredictRequest": {
      "type": "object",
      "properties": {
        "model_name": {
          "type": "string"
        },
        "model_version": {
          "type": "string"
        },
        "framework_name": {
          "type": "string"
        },
        "framework_version": {
          "type": "string"
        },
        "limit": {
          "type": "integer",
          "format": "int32"
        },
        "data": {
          "type": "string",
          "format": "byte"
        },
        "url": {
          "type": "string"
        }
      }
    },
    "dlframeworkPredictResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "features": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkPredictionFeature"
          }
        },
        "error": {
          "$ref": "#/definitions/dlframeworkErrorStatus"
        }
      }
    },
    "dlframeworkPredictionFeature": {
      "type": "object",
      "properties": {
        "index": {
          "type": "string",
          "format": "int64"
        },
        "name": {
          "type": "string"
        },
        "probability": {
          "type": "number",
          "format": "float"
        }
      }
    }
  }
}
`
	swagger_info = `{
  "info": {
    "title": "CarML DLFramework",
    "version": "1.0.0",
    "description": "CarML (Cognitive ARtifacts for Machine Learning) is a framework allowing people to develop and deploy machine learning models. It allows machine learning (ML) developers to publish and evaluate their models, users to experiment with different models and frameworks through a web user interface or a REST api, and system architects to capture system resource usage to inform future system and hardware configuration."
  }
}
`
)
