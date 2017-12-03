package dlframework

const (
	dlframework_swagger = `{
  "swagger": "2.0",
  "info": {
    "title": "CarML DLFramework",
    "version": "0.2.18",
    "description": "CarML (Cognitive ARtifacts for Machine Learning) is a framework allowing people to develop and deploy machine learning models. It allows machine learning (ML) developers to publish and evaluate their models, users to experiment with different models and frameworks through a web user interface or a REST api, and system architects to capture system resource usage to inform future system and hardware configuration.",
    "contact": {
      "name": "Abdul Dakkak, Cheng Li",
      "url": "https://github.com/rai-project/carml"
    },
    "license": {
      "name": "NCSA/UIUC",
      "url": "https://raw.githubusercontent.com/rai-project/dlframework/master/LICENSE.TXT"
    }
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
    "/predict/close": {
      "post": {
        "summary": "Close a predictor clear it's memory.",
        "operationId": "Close",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkPredictorCloseResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkPredictor"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/dataset": {
      "post": {
        "summary": "Dataset method receives a single dataset and runs\nthe predictor on all elements of the dataset.",
        "description": "The result is a prediction feature list.",
        "operationId": "Dataset",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkFeaturesResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkDatasetRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/images": {
      "post": {
        "summary": "Image method receives a list base64 encoded images and runs\nthe predictor on all the images.",
        "description": "The result is a prediction feature list for each image.",
        "operationId": "Images",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkFeaturesResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkImagesRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/open": {
      "post": {
        "summary": "Opens a predictor and returns an id where the predictor\nis accessible. The id can be used to perform inference\nrequests.",
        "operationId": "Open",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkPredictor"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkPredictorOpenRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/reset": {
      "post": {
        "summary": "Clear method clears the internal cache of the predictors",
        "operationId": "Reset",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkResetResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkResetRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/stream/dataset": {
      "post": {
        "summary": "Dataset method receives a single dataset and runs\nthe predictor on all elements of the dataset.",
        "description": "The result is a prediction feature stream.",
        "operationId": "DatasetStream",
        "responses": {
          "200": {
            "description": "(streaming responses)",
            "schema": {
              "$ref": "#/definitions/dlframeworkFeatureResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkDatasetRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/stream/images": {
      "post": {
        "summary": "Image method receives a list base64 encoded images and runs\nthe predictor on all the images.",
        "description": "The result is a prediction feature stream for each image.",
        "operationId": "ImagesStream",
        "responses": {
          "200": {
            "description": "(streaming responses)",
            "schema": {
              "$ref": "#/definitions/dlframeworkFeatureResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkImagesRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/stream/urls": {
      "post": {
        "summary": "Image method receives a stream of urls and runs\nthe predictor on all the urls. The",
        "description": "The result is a prediction feature stream for each url.",
        "operationId": "URLsStream",
        "responses": {
          "200": {
            "description": "(streaming responses)",
            "schema": {
              "$ref": "#/definitions/dlframeworkFeatureResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkURLsRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/predict/urls": {
      "post": {
        "summary": "Image method receives a stream of urls and runs\nthe predictor on all the urls. The",
        "description": "The result is a prediction feature stream for each url.",
        "operationId": "URLs",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkFeaturesResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dlframeworkURLsRequest"
            }
          }
        ],
        "tags": [
          "Predict"
        ]
      }
    },
    "/registry/frameworks/agent": {
      "get": {
        "operationId": "FrameworkAgents",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkAgents"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "query",
            "required": false,
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
    "/registry/frameworks/manifest": {
      "get": {
        "operationId": "FrameworkManifests",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkFrameworkManifestsResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "query",
            "required": false,
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
    "/registry/models/agent": {
      "get": {
        "operationId": "ModelAgents",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkAgents"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "framework_version",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "model_name",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "model_version",
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
    "/registry/models/manifest": {
      "get": {
        "operationId": "ModelManifests",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/dlframeworkModelManifestsResponse"
            }
          }
        },
        "parameters": [
          {
            "name": "framework_name",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "framework_version",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "model_name",
            "in": "query",
            "required": false,
            "type": "string"
          },
          {
            "name": "model_version",
            "in": "query",
            "required": false,
            "type": "string"
          }
        ],
        "tags": [
          "Registry"
        ]
      }
    }
  },
  "definitions": {
    "DatasetRequestDataset": {
      "type": "object",
      "properties": {
        "category": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "ExecutionOptionsTraceLevel": {
      "type": "string",
      "enum": [
        "NO_TRACE",
        "STEP_TRACE",
        "FRAMEWORK_TRACE",
        "CPU_ONLY_TRACE",
        "HARDWARE_TRACE",
        "FULL_TRACE"
      ],
      "default": "NO_TRACE"
    },
    "ImagesRequestImage": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "title": "An id used to identify the output feature: maps to input_id for output"
        },
        "data": {
          "type": "string",
          "format": "byte",
          "title": "The image is base64 encoded"
        }
      }
    },
    "URLsRequestURL": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "title": "An id used to identify the output feature: maps to input_id for output"
        },
        "data": {
          "type": "string"
        }
      }
    },
    "dlframeworkCPUOptions": {
      "type": "object"
    },
    "dlframeworkDatasetRequest": {
      "type": "object",
      "properties": {
        "predictor": {
          "$ref": "#/definitions/dlframeworkPredictor"
        },
        "dataset": {
          "$ref": "#/definitions/DatasetRequestDataset"
        },
        "options": {
          "$ref": "#/definitions/dlframeworkPredictionOptions"
        }
      }
    },
    "dlframeworkExecutionOptions": {
      "type": "object",
      "properties": {
        "trace_level": {
          "$ref": "#/definitions/ExecutionOptionsTraceLevel"
        },
        "timeout_in_ms": {
          "type": "string",
          "format": "int64",
          "description": "Time to wait for operation to complete in milliseconds."
        },
        "device_count": {
          "type": "object",
          "additionalProperties": {
            "type": "integer",
            "format": "int32"
          },
          "description": "Map from device type name (e.g., \"CPU\" or \"GPU\" ) to maximum\nnumber of devices of that type to use.  If a particular device\ntype is not found in the map, the system picks an appropriate\nnumber."
        },
        "cpu_options": {
          "$ref": "#/definitions/dlframeworkCPUOptions",
          "description": "Options that apply to all CPUs."
        },
        "gpu_options": {
          "$ref": "#/definitions/dlframeworkGPUOptions",
          "description": "Options that apply to all GPUs."
        }
      }
    },
    "dlframeworkFeature": {
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
        },
        "metadata": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "dlframeworkFeatureResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "request_id": {
          "type": "string"
        },
        "input_id": {
          "type": "string"
        },
        "features": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkFeature"
          }
        },
        "metadata": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "dlframeworkFeaturesResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "responses": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkFeatureResponse"
          }
        }
      }
    },
    "dlframeworkGPUOptions": {
      "type": "object",
      "properties": {
        "per_process_gpu_memory_fraction": {
          "type": "number",
          "format": "double",
          "description": "A value between 0 and 1 that indicates what fraction of the\navailable GPU memory to pre-allocate for each process.  1 means\nto pre-allocate all of the GPU memory, 0.5 means the process\nallocates ~50% of the available GPU memory."
        },
        "allocator_type": {
          "type": "string",
          "description": "The type of GPU allocation strategy to use.\n\nAllowed values:\n\"\": The empty string (default) uses a system-chosen default\n    which may change over time.\n\n\"BFC\": A \"Best-fit with coalescing\" algorithm, simplified from a\n       version of dlmalloc."
        },
        "visible_device_list": {
          "type": "string",
          "description": "A comma-separated list of GPU ids that determines the 'visible'\nto 'virtual' mapping of GPU devices.  For example, if TensorFlow\ncan see 8 GPU devices in the process, and one wanted to map\nvisible GPU devices 5 and 3 as \"/device:GPU:0\", and \"/device:GPU:1\", then\none would specify this field as \"5,3\".  This field is similar in spirit to\nthe CUDA_VISIBLE_DEVICES environment variable, except it applies to the\nvisible GPU devices in the process.\n\nNOTE: The GPU driver provides the process with the visible GPUs\nin an order which is not guaranteed to have any correlation to\nthe *physical* GPU id in the machine.  This field is used for\nremapping \"visible\" to \"virtual\", which means this operates only\nafter the process starts.  Users are required to use vendor\nspecific mechanisms (e.g., CUDA_VISIBLE_DEVICES) to control the\nphysical to visible device mapping prior to invoking TensorFlow."
        },
        "force_gpu_compatible": {
          "type": "boolean",
          "format": "boolean",
          "description": "Force all tensors to be gpu_compatible. On a GPU-enabled TensorFlow,\nenabling this option forces all CPU tensors to be allocated with Cuda\npinned memory. Normally, TensorFlow will infer which tensors should be\nallocated as the pinned memory. But in case where the inference is\nincomplete, this option can significantly speed up the cross-device memory\ncopy performance as long as it fits the memory.\nNote that this option is not something that should be\nenabled by default for unknown or very large models, since all Cuda pinned\nmemory is unpageable, having too much pinned memory might negatively impact\nthe overall host system performance."
        }
      }
    },
    "dlframeworkImagesRequest": {
      "type": "object",
      "properties": {
        "predictor": {
          "$ref": "#/definitions/dlframeworkPredictor"
        },
        "images": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ImagesRequestImage"
          },
          "title": "A list of Base64 encoded images"
        },
        "options": {
          "$ref": "#/definitions/dlframeworkPredictionOptions"
        }
      }
    },
    "dlframeworkPredictionOptions": {
      "type": "object",
      "properties": {
        "request_id": {
          "type": "string"
        },
        "feature_limit": {
          "type": "integer",
          "format": "int32"
        },
        "batch_size": {
          "type": "integer",
          "format": "int64"
        },
        "execution_options": {
          "$ref": "#/definitions/dlframeworkExecutionOptions"
        },
        "agent": {
          "type": "string"
        }
      }
    },
    "dlframeworkPredictor": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        }
      }
    },
    "dlframeworkPredictorCloseResponse": {
      "type": "object"
    },
    "dlframeworkPredictorOpenRequest": {
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
        "options": {
          "$ref": "#/definitions/dlframeworkPredictionOptions"
        }
      }
    },
    "dlframeworkResetRequest": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        },
        "predictor": {
          "$ref": "#/definitions/dlframeworkPredictor"
        }
      }
    },
    "dlframeworkResetResponse": {
      "type": "object",
      "properties": {
        "predictor": {
          "$ref": "#/definitions/dlframeworkPredictor"
        }
      }
    },
    "dlframeworkURLsRequest": {
      "type": "object",
      "properties": {
        "predictor": {
          "$ref": "#/definitions/dlframeworkPredictor"
        },
        "urls": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/URLsRequestURL"
          }
        },
        "options": {
          "$ref": "#/definitions/dlframeworkPredictionOptions"
        }
      }
    },
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
        },
        "weights_checksum": {
          "type": "string"
        },
        "graph_checksum": {
          "type": "string"
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
    "dlframeworkAgent": {
      "type": "object",
      "properties": {
        "host": {
          "type": "string"
        },
        "port": {
          "type": "string"
        },
        "hostname": {
          "type": "string"
        },
        "architecture": {
          "type": "string"
        },
        "hasgpu": {
          "type": "boolean",
          "format": "boolean"
        },
        "cpuinfo": {
          "type": "string"
        },
        "gpuinfo": {
          "type": "string"
        },
        "frameworks": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkFrameworkManifest"
          }
        },
        "metadata": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          }
        }
      }
    },
    "dlframeworkAgents": {
      "type": "object",
      "properties": {
        "agents": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkAgent"
          }
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
    "dlframeworkFrameworkManifestsResponse": {
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
        },
        "hidden": {
          "type": "boolean",
          "format": "boolean"
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
    "dlframeworkModelManifestsResponse": {
      "type": "object",
      "properties": {
        "manifests": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dlframeworkModelManifest"
          }
        }
      }
    }
  },
  "host": "carml.org",
  "basePath": "/api",
  "externalDocs": {
    "url": "https://rai-project.github.io/carml"
  }
}
`
	swagger_info = `{
	"info": {
		"title": "CarML DLFramework",
		"description":
			"CarML (Cognitive ARtifacts for Machine Learning) is a framework allowing people to develop and deploy machine learning models. It allows machine learning (ML) developers to publish and evaluate their models, users to experiment with different models and frameworks through a web user interface or a REST api, and system architects to capture system resource usage to inform future system and hardware configuration.",
		"version": "0.2.18",
		"contact": {
			"name": "Abdul Dakkak, Cheng Li",
			"url": "https://github.com/rai-project/carml"
		},
		"license": {
			"name": "NCSA/UIUC",
			"url": "https://raw.githubusercontent.com/rai-project/dlframework/master/LICENSE.TXT"
		}
	},
	"host": "carml.org",
	"basePath": "/api",
	"externalDocs": {
		"url": "https://rai-project.github.io/carml"
	}
}
`
)
