package mutate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	v1beta1 "k8s.io/api/admission/v1beta1"
)

func TestMutateJSONWithoutTTL(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"apiVersion": "argoproj.io/v1alpha1",
				"kind": "Workflow",
				"metadata": {
					"creationTimestamp": "2019-11-30T09:00:01Z",
					"generation": 5,
					"labels": {
						"pipeline/persistedFinalState": "true",
						"pipeline/runid": "b5522a9e-b181-497c-9d49-473a4f6cdcd7",
						"scheduledworkflows.kubeflow.org/isOwnedByScheduledWorkflow": "true",
						"scheduledworkflows.kubeflow.org/scheduledWorkflowName": "sandtradeedcdailyrunhwn89",
						"scheduledworkflows.kubeflow.org/workflowEpoch": "1575104400",
						"scheduledworkflows.kubeflow.org/workflowIndex": "3",
						"workflows.argoproj.io/completed": "true",
						"workflows.argoproj.io/phase": "Failed"
					},
					"name": "sandtradeedcdailyrunhwn89-3-3509392613",
					"namespace": "kubeflow",
					"ownerReferences": [
						{
							"apiVersion": "kubeflow.org/v1beta1",
							"blockOwnerDeletion": true,
							"controller": true,
							"kind": "ScheduledWorkflow",
							"name": "sandtradeedcdailyrunhwn89",
							"uid": "bf28c1fb-10f7-11ea-a12d-42010a4c1005"
						}
					],
					"resourceVersion": "7360201",
					"selfLink": "/apis/argoproj.io/v1alpha1/namespaces/kubeflow/workflows/sandtradeedcdailyrunhwn89-3-3509392613",
					"uid": "cb09d71e-134f-11ea-a12d-42010a4c1005"
				},
				"spec": {
					"arguments": {
						"parameters": [
							{
								"name": "preprocessing-blob-name",
								"value": "data/preprocessing/input.csv"
							},
							{
								"name": "train-features-blob-name",
								"value": "train.csv"
							},
							{
								"name": "forecast-features-blob-name",
								"value": "forecast.csv"
							},
							{
								"name": "bucket",
								"value": "sandtrade_production"
							},
							{
								"name": "r-forecast-blob-name",
								"value": "forecast_r.csv"
							},
							{
								"name": "prophet-forecast-blob-name",
								"value": "forecast_prophet.csv"
							},
							{
								"name": "ensamble-output",
								"value": "output.csv"
							}
						]
					},
					"entrypoint": "sandtrade-demand-forecasting",
					"serviceAccountName": "pipeline-runner",
					"templates": [
						{
							"container": {
								"args": [
									"--vanilla",
									"main.R",
									"-i",
									"{{inputs.parameters.preprocessing-blob-name}}",
									"-l",
									"tmp.csv",
									"-t",
									"{{inputs.parameters.train-features-blob-name}}",
									"-f",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"-b",
									"{{inputs.parameters.bucket}}"
								],
								"command": [
									"Rscript"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/feature-generation-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "preprocessing-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "feature-generation",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"post_processing.py",
									"--linearmodelrpredictions",
									"{{inputs.parameters.r-forecast-blob-name}}",
									"--linearmodelprophetpredictions",
									"{{inputs.parameters.prophet-forecast-blob-name}}",
									"--localfilenamelinear",
									"tmp_linear.csv",
									"--localfilenameprophet",
									"tmp_prophet.csv",
									"--localoutputfilename",
									"{{inputs.parameters.ensamble-output}}",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--secret_name",
									"database",
									"--secret_path",
									"/secret/test-credentials"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/postprocessing-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/test-credentials",
										"name": "mysecretvolume"
									},
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "ensamble-output"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "post-processing-image",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								},
								{
									"name": "mysecretvolume",
									"secret": {
										"secretName": "database"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"main.py",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--destination_blob_name",
									"{{inputs.parameters.preprocessing-blob-name}}",
									"--localfilename",
									"tmp.csv",
									"--secret_name",
									"database",
									"--secret_path",
									"/secret/test-credentials"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/preprocessing-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/test-credentials",
										"name": "mysecretvolume"
									},
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "preprocessing-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "preprocess",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								},
								{
									"name": "mysecretvolume",
									"secret": {
										"secretName": "database"
									}
								}
							]
						},
						{
							"dag": {
								"tasks": [
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "preprocessing-blob-name",
													"value": "{{inputs.parameters.preprocessing-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"preprocess"
										],
										"name": "feature-generation",
										"template": "feature-generation"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "ensamble-output",
													"value": "{{inputs.parameters.ensamble-output}}"
												},
												{
													"name": "prophet-forecast-blob-name",
													"value": "{{inputs.parameters.prophet-forecast-blob-name}}"
												},
												{
													"name": "r-forecast-blob-name",
													"value": "{{inputs.parameters.r-forecast-blob-name}}"
												}
											]
										},
										"dependencies": [
											"train-predict-using-propeht-model",
											"train-r-linear-model"
										],
										"name": "post-processing-image",
										"template": "post-processing-image"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "preprocessing-blob-name",
													"value": "{{inputs.parameters.preprocessing-blob-name}}"
												}
											]
										},
										"name": "preprocess",
										"template": "preprocess"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "prophet-forecast-blob-name",
													"value": "{{inputs.parameters.prophet-forecast-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"feature-generation"
										],
										"name": "train-predict-using-propeht-model",
										"template": "train-predict-using-propeht-model"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "r-forecast-blob-name",
													"value": "{{inputs.parameters.r-forecast-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"feature-generation"
										],
										"name": "train-r-linear-model",
										"template": "train-r-linear-model"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "ensamble-output"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "preprocessing-blob-name"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "sandtrade-demand-forecasting",
							"outputs": {}
						},
						{
							"container": {
								"args": [
									"main.py",
									"--modelpath",
									".",
									"--inputtrain",
									"{{inputs.parameters.train-features-blob-name}}",
									"--inputforecast",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"--localfilename",
									"tmp.csv",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--outputforecast",
									"{{inputs.parameters.prophet-forecast-blob-name}}"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/train-predict-python-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "train-predict-using-propeht-model",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"--vanilla",
									"main.R",
									"-t",
									"{{inputs.parameters.train-features-blob-name}}",
									"-f",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"-l",
									"tmp_train.csv",
									"-k",
									"tmp_forecast.csv",
									"-j",
									"tmp_output.csv",
									"-o",
									"{{inputs.parameters.r-forecast-blob-name}}",
									"-b",
									"{{inputs.parameters.bucket}}"
								],
								"command": [
									"Rscript"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/train-predict-r-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "train-r-linear-model",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						}
					]
				},
				"status": {
					"finishedAt": "2019-11-30T09:00:17Z",
					"nodes": {
						"sandtradeedcdailyrunhwn89-3-3509392613": {
							"children": [
								"sandtradeedcdailyrunhwn89-3-3509392613-884332321"
							],
							"displayName": "sandtradeedcdailyrunhwn89-3-3509392613",
							"finishedAt": "2019-11-30T09:00:17Z",
							"id": "sandtradeedcdailyrunhwn89-3-3509392613",
							"inputs": {
								"parameters": [
									{
										"name": "bucket",
										"value": "sandtrade_production"
									},
									{
										"name": "ensamble-output",
										"value": "output.csv"
									},
									{
										"name": "forecast-features-blob-name",
										"value": "forecast.csv"
									},
									{
										"name": "preprocessing-blob-name",
										"value": "data/preprocessing/input.csv"
									},
									{
										"name": "prophet-forecast-blob-name",
										"value": "forecast_prophet.csv"
									},
									{
										"name": "r-forecast-blob-name",
										"value": "forecast_r.csv"
									},
									{
										"name": "train-features-blob-name",
										"value": "train.csv"
									}
								]
							},
							"name": "sandtradeedcdailyrunhwn89-3-3509392613",
							"phase": "Failed",
							"startedAt": "2019-11-30T09:00:01Z",
							"templateName": "sandtrade-demand-forecasting",
							"type": "DAG"
						},
						"sandtradeedcdailyrunhwn89-3-3509392613-884332321": {
							"boundaryID": "sandtradeedcdailyrunhwn89-3-3509392613",
							"displayName": "preprocess",
							"finishedAt": "2019-11-30T09:00:16Z",
							"id": "sandtradeedcdailyrunhwn89-3-3509392613-884332321",
							"inputs": {
								"parameters": [
									{
										"name": "bucket",
										"value": "sandtrade_production"
									},
									{
										"name": "preprocessing-blob-name",
										"value": "data/preprocessing/input.csv"
									}
								]
							},
							"message": "failed with exit code 1",
							"name": "sandtradeedcdailyrunhwn89-3-3509392613.preprocess",
							"phase": "Failed",
							"startedAt": "2019-11-30T09:00:02Z",
							"templateName": "preprocess",
							"type": "Pod"
						}
					},
					"phase": "Failed",
					"startedAt": "2019-11-30T09:00:01Z"
				}

			}
		}
	}` // Example of workflow

	response, err := Mutate([]byte(rawJSON))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"add","path":"/spec/ttlSecondsAfterFinished","value":36000}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")
}

func TestMutateJSONWithTTL(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"apiVersion": "argoproj.io/v1alpha1",
				"kind": "Workflow",
				"metadata": {
					"creationTimestamp": "2019-11-30T09:00:01Z",
					"generation": 5,
					"labels": {
						"pipeline/persistedFinalState": "true",
						"pipeline/runid": "b5522a9e-b181-497c-9d49-473a4f6cdcd7",
						"scheduledworkflows.kubeflow.org/isOwnedByScheduledWorkflow": "true",
						"scheduledworkflows.kubeflow.org/scheduledWorkflowName": "sandtradeedcdailyrunhwn89",
						"scheduledworkflows.kubeflow.org/workflowEpoch": "1575104400",
						"scheduledworkflows.kubeflow.org/workflowIndex": "3",
						"workflows.argoproj.io/completed": "true",
						"workflows.argoproj.io/phase": "Failed"
					},
					"name": "sandtradeedcdailyrunhwn89-3-3509392613",
					"namespace": "kubeflow",
					"ownerReferences": [
						{
							"apiVersion": "kubeflow.org/v1beta1",
							"blockOwnerDeletion": true,
							"controller": true,
							"kind": "ScheduledWorkflow",
							"name": "sandtradeedcdailyrunhwn89",
							"uid": "bf28c1fb-10f7-11ea-a12d-42010a4c1005"
						}
					],
					"resourceVersion": "7360201",
					"selfLink": "/apis/argoproj.io/v1alpha1/namespaces/kubeflow/workflows/sandtradeedcdailyrunhwn89-3-3509392613",
					"uid": "cb09d71e-134f-11ea-a12d-42010a4c1005"
				},
				"spec": {
					"arguments": {
						"parameters": [
							{
								"name": "preprocessing-blob-name",
								"value": "data/preprocessing/input.csv"
							},
							{
								"name": "train-features-blob-name",
								"value": "train.csv"
							},
							{
								"name": "forecast-features-blob-name",
								"value": "forecast.csv"
							},
							{
								"name": "bucket",
								"value": "sandtrade_production"
							},
							{
								"name": "r-forecast-blob-name",
								"value": "forecast_r.csv"
							},
							{
								"name": "prophet-forecast-blob-name",
								"value": "forecast_prophet.csv"
							},
							{
								"name": "ensamble-output",
								"value": "output.csv"
							}
						]
					},
					"entrypoint": "sandtrade-demand-forecasting",
					"serviceAccountName": "pipeline-runner",
					"templates": [
						{
							"container": {
								"args": [
									"--vanilla",
									"main.R",
									"-i",
									"{{inputs.parameters.preprocessing-blob-name}}",
									"-l",
									"tmp.csv",
									"-t",
									"{{inputs.parameters.train-features-blob-name}}",
									"-f",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"-b",
									"{{inputs.parameters.bucket}}"
								],
								"command": [
									"Rscript"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/feature-generation-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "preprocessing-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "feature-generation",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"post_processing.py",
									"--linearmodelrpredictions",
									"{{inputs.parameters.r-forecast-blob-name}}",
									"--linearmodelprophetpredictions",
									"{{inputs.parameters.prophet-forecast-blob-name}}",
									"--localfilenamelinear",
									"tmp_linear.csv",
									"--localfilenameprophet",
									"tmp_prophet.csv",
									"--localoutputfilename",
									"{{inputs.parameters.ensamble-output}}",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--secret_name",
									"database",
									"--secret_path",
									"/secret/test-credentials"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/postprocessing-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/test-credentials",
										"name": "mysecretvolume"
									},
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "ensamble-output"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "post-processing-image",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								},
								{
									"name": "mysecretvolume",
									"secret": {
										"secretName": "database"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"main.py",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--destination_blob_name",
									"{{inputs.parameters.preprocessing-blob-name}}",
									"--localfilename",
									"tmp.csv",
									"--secret_name",
									"database",
									"--secret_path",
									"/secret/test-credentials"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/preprocessing-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/test-credentials",
										"name": "mysecretvolume"
									},
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "preprocessing-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "preprocess",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								},
								{
									"name": "mysecretvolume",
									"secret": {
										"secretName": "database"
									}
								}
							]
						},
						{
							"dag": {
								"tasks": [
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "preprocessing-blob-name",
													"value": "{{inputs.parameters.preprocessing-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"preprocess"
										],
										"name": "feature-generation",
										"template": "feature-generation"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "ensamble-output",
													"value": "{{inputs.parameters.ensamble-output}}"
												},
												{
													"name": "prophet-forecast-blob-name",
													"value": "{{inputs.parameters.prophet-forecast-blob-name}}"
												},
												{
													"name": "r-forecast-blob-name",
													"value": "{{inputs.parameters.r-forecast-blob-name}}"
												}
											]
										},
										"dependencies": [
											"train-predict-using-propeht-model",
											"train-r-linear-model"
										],
										"name": "post-processing-image",
										"template": "post-processing-image"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "preprocessing-blob-name",
													"value": "{{inputs.parameters.preprocessing-blob-name}}"
												}
											]
										},
										"name": "preprocess",
										"template": "preprocess"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "prophet-forecast-blob-name",
													"value": "{{inputs.parameters.prophet-forecast-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"feature-generation"
										],
										"name": "train-predict-using-propeht-model",
										"template": "train-predict-using-propeht-model"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "r-forecast-blob-name",
													"value": "{{inputs.parameters.r-forecast-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"feature-generation"
										],
										"name": "train-r-linear-model",
										"template": "train-r-linear-model"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "ensamble-output"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "preprocessing-blob-name"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "sandtrade-demand-forecasting",
							"outputs": {}
						},
						{
							"container": {
								"args": [
									"main.py",
									"--modelpath",
									".",
									"--inputtrain",
									"{{inputs.parameters.train-features-blob-name}}",
									"--inputforecast",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"--localfilename",
									"tmp.csv",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--outputforecast",
									"{{inputs.parameters.prophet-forecast-blob-name}}"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/train-predict-python-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "train-predict-using-propeht-model",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"--vanilla",
									"main.R",
									"-t",
									"{{inputs.parameters.train-features-blob-name}}",
									"-f",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"-l",
									"tmp_train.csv",
									"-k",
									"tmp_forecast.csv",
									"-j",
									"tmp_output.csv",
									"-o",
									"{{inputs.parameters.r-forecast-blob-name}}",
									"-b",
									"{{inputs.parameters.bucket}}"
								],
								"command": [
									"Rscript"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/train-predict-r-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "train-r-linear-model",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						}
					],
					"ttlSecondsAfterFinished": 1800
				},
				"status": {
					"finishedAt": "2019-11-30T09:00:17Z",
					"nodes": {
						"sandtradeedcdailyrunhwn89-3-3509392613": {
							"children": [
								"sandtradeedcdailyrunhwn89-3-3509392613-884332321"
							],
							"displayName": "sandtradeedcdailyrunhwn89-3-3509392613",
							"finishedAt": "2019-11-30T09:00:17Z",
							"id": "sandtradeedcdailyrunhwn89-3-3509392613",
							"inputs": {
								"parameters": [
									{
										"name": "bucket",
										"value": "sandtrade_production"
									},
									{
										"name": "ensamble-output",
										"value": "output.csv"
									},
									{
										"name": "forecast-features-blob-name",
										"value": "forecast.csv"
									},
									{
										"name": "preprocessing-blob-name",
										"value": "data/preprocessing/input.csv"
									},
									{
										"name": "prophet-forecast-blob-name",
										"value": "forecast_prophet.csv"
									},
									{
										"name": "r-forecast-blob-name",
										"value": "forecast_r.csv"
									},
									{
										"name": "train-features-blob-name",
										"value": "train.csv"
									}
								]
							},
							"name": "sandtradeedcdailyrunhwn89-3-3509392613",
							"phase": "Failed",
							"startedAt": "2019-11-30T09:00:01Z",
							"templateName": "sandtrade-demand-forecasting",
							"type": "DAG"
						},
						"sandtradeedcdailyrunhwn89-3-3509392613-884332321": {
							"boundaryID": "sandtradeedcdailyrunhwn89-3-3509392613",
							"displayName": "preprocess",
							"finishedAt": "2019-11-30T09:00:16Z",
							"id": "sandtradeedcdailyrunhwn89-3-3509392613-884332321",
							"inputs": {
								"parameters": [
									{
										"name": "bucket",
										"value": "sandtrade_production"
									},
									{
										"name": "preprocessing-blob-name",
										"value": "data/preprocessing/input.csv"
									}
								]
							},
							"message": "failed with exit code 1",
							"name": "sandtradeedcdailyrunhwn89-3-3509392613.preprocess",
							"phase": "Failed",
							"startedAt": "2019-11-30T09:00:02Z",
							"templateName": "preprocess",
							"type": "Pod"
						}
					},
					"phase": "Failed",
					"startedAt": "2019-11-30T09:00:01Z"
				}
			}
		}
	}` // Example of workflow// Example of workflow

	response, err := Mutate([]byte(rawJSON))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")
}
func TestMutateJSONWithTTLToLong(t *testing.T) {
	rawJSON := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1beta1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"apiVersion": "argoproj.io/v1alpha1",
				"kind": "Workflow",
				"metadata": {
					"creationTimestamp": "2019-11-30T09:00:01Z",
					"generation": 5,
					"labels": {
						"pipeline/persistedFinalState": "true",
						"pipeline/runid": "b5522a9e-b181-497c-9d49-473a4f6cdcd7",
						"scheduledworkflows.kubeflow.org/isOwnedByScheduledWorkflow": "true",
						"scheduledworkflows.kubeflow.org/scheduledWorkflowName": "sandtradeedcdailyrunhwn89",
						"scheduledworkflows.kubeflow.org/workflowEpoch": "1575104400",
						"scheduledworkflows.kubeflow.org/workflowIndex": "3",
						"workflows.argoproj.io/completed": "true",
						"workflows.argoproj.io/phase": "Failed"
					},
					"name": "sandtradeedcdailyrunhwn89-3-3509392613",
					"namespace": "kubeflow",
					"ownerReferences": [
						{
							"apiVersion": "kubeflow.org/v1beta1",
							"blockOwnerDeletion": true,
							"controller": true,
							"kind": "ScheduledWorkflow",
							"name": "sandtradeedcdailyrunhwn89",
							"uid": "bf28c1fb-10f7-11ea-a12d-42010a4c1005"
						}
					],
					"resourceVersion": "7360201",
					"selfLink": "/apis/argoproj.io/v1alpha1/namespaces/kubeflow/workflows/sandtradeedcdailyrunhwn89-3-3509392613",
					"uid": "cb09d71e-134f-11ea-a12d-42010a4c1005"
				},
				"spec": {
					"arguments": {
						"parameters": [
							{
								"name": "preprocessing-blob-name",
								"value": "data/preprocessing/input.csv"
							},
							{
								"name": "train-features-blob-name",
								"value": "train.csv"
							},
							{
								"name": "forecast-features-blob-name",
								"value": "forecast.csv"
							},
							{
								"name": "bucket",
								"value": "sandtrade_production"
							},
							{
								"name": "r-forecast-blob-name",
								"value": "forecast_r.csv"
							},
							{
								"name": "prophet-forecast-blob-name",
								"value": "forecast_prophet.csv"
							},
							{
								"name": "ensamble-output",
								"value": "output.csv"
							}
						]
					},
					"entrypoint": "sandtrade-demand-forecasting",
					"serviceAccountName": "pipeline-runner",
					"templates": [
						{
							"container": {
								"args": [
									"--vanilla",
									"main.R",
									"-i",
									"{{inputs.parameters.preprocessing-blob-name}}",
									"-l",
									"tmp.csv",
									"-t",
									"{{inputs.parameters.train-features-blob-name}}",
									"-f",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"-b",
									"{{inputs.parameters.bucket}}"
								],
								"command": [
									"Rscript"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/feature-generation-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "preprocessing-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "feature-generation",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"post_processing.py",
									"--linearmodelrpredictions",
									"{{inputs.parameters.r-forecast-blob-name}}",
									"--linearmodelprophetpredictions",
									"{{inputs.parameters.prophet-forecast-blob-name}}",
									"--localfilenamelinear",
									"tmp_linear.csv",
									"--localfilenameprophet",
									"tmp_prophet.csv",
									"--localoutputfilename",
									"{{inputs.parameters.ensamble-output}}",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--secret_name",
									"database",
									"--secret_path",
									"/secret/test-credentials"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/postprocessing-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/test-credentials",
										"name": "mysecretvolume"
									},
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "ensamble-output"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "post-processing-image",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								},
								{
									"name": "mysecretvolume",
									"secret": {
										"secretName": "database"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"main.py",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--destination_blob_name",
									"{{inputs.parameters.preprocessing-blob-name}}",
									"--localfilename",
									"tmp.csv",
									"--secret_name",
									"database",
									"--secret_path",
									"/secret/test-credentials"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/preprocessing-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/test-credentials",
										"name": "mysecretvolume"
									},
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "preprocessing-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "preprocess",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								},
								{
									"name": "mysecretvolume",
									"secret": {
										"secretName": "database"
									}
								}
							]
						},
						{
							"dag": {
								"tasks": [
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "preprocessing-blob-name",
													"value": "{{inputs.parameters.preprocessing-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"preprocess"
										],
										"name": "feature-generation",
										"template": "feature-generation"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "ensamble-output",
													"value": "{{inputs.parameters.ensamble-output}}"
												},
												{
													"name": "prophet-forecast-blob-name",
													"value": "{{inputs.parameters.prophet-forecast-blob-name}}"
												},
												{
													"name": "r-forecast-blob-name",
													"value": "{{inputs.parameters.r-forecast-blob-name}}"
												}
											]
										},
										"dependencies": [
											"train-predict-using-propeht-model",
											"train-r-linear-model"
										],
										"name": "post-processing-image",
										"template": "post-processing-image"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "preprocessing-blob-name",
													"value": "{{inputs.parameters.preprocessing-blob-name}}"
												}
											]
										},
										"name": "preprocess",
										"template": "preprocess"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "prophet-forecast-blob-name",
													"value": "{{inputs.parameters.prophet-forecast-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"feature-generation"
										],
										"name": "train-predict-using-propeht-model",
										"template": "train-predict-using-propeht-model"
									},
									{
										"arguments": {
											"parameters": [
												{
													"name": "bucket",
													"value": "{{inputs.parameters.bucket}}"
												},
												{
													"name": "forecast-features-blob-name",
													"value": "{{inputs.parameters.forecast-features-blob-name}}"
												},
												{
													"name": "r-forecast-blob-name",
													"value": "{{inputs.parameters.r-forecast-blob-name}}"
												},
												{
													"name": "train-features-blob-name",
													"value": "{{inputs.parameters.train-features-blob-name}}"
												}
											]
										},
										"dependencies": [
											"feature-generation"
										],
										"name": "train-r-linear-model",
										"template": "train-r-linear-model"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "ensamble-output"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "preprocessing-blob-name"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "sandtrade-demand-forecasting",
							"outputs": {}
						},
						{
							"container": {
								"args": [
									"main.py",
									"--modelpath",
									".",
									"--inputtrain",
									"{{inputs.parameters.train-features-blob-name}}",
									"--inputforecast",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"--localfilename",
									"tmp.csv",
									"--bucket",
									"{{inputs.parameters.bucket}}",
									"--outputforecast",
									"{{inputs.parameters.prophet-forecast-blob-name}}"
								],
								"command": [
									"python3"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/train-predict-python-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "prophet-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "train-predict-using-propeht-model",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						},
						{
							"container": {
								"args": [
									"--vanilla",
									"main.R",
									"-t",
									"{{inputs.parameters.train-features-blob-name}}",
									"-f",
									"{{inputs.parameters.forecast-features-blob-name}}",
									"-l",
									"tmp_train.csv",
									"-k",
									"tmp_forecast.csv",
									"-j",
									"tmp_output.csv",
									"-o",
									"{{inputs.parameters.r-forecast-blob-name}}",
									"-b",
									"{{inputs.parameters.bucket}}"
								],
								"command": [
									"Rscript"
								],
								"env": [
									{
										"name": "GOOGLE_APPLICATION_CREDENTIALS",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									},
									{
										"name": "CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE",
										"value": "/secret/gcp-credentials/user-gcp-sa.json"
									}
								],
								"image": "gcr.io/ds-production-259110/train-predict-r-forecasting",
								"imagePullPolicy": "Always",
								"name": "",
								"resources": {},
								"volumeMounts": [
									{
										"mountPath": "/secret/gcp-credentials",
										"name": "gcp-credentials-user-gcp-sa"
									}
								]
							},
							"inputs": {
								"parameters": [
									{
										"name": "bucket"
									},
									{
										"name": "forecast-features-blob-name"
									},
									{
										"name": "r-forecast-blob-name"
									},
									{
										"name": "train-features-blob-name"
									}
								]
							},
							"metadata": {},
							"name": "train-r-linear-model",
							"outputs": {},
							"volumes": [
								{
									"name": "gcp-credentials-user-gcp-sa",
									"secret": {
										"secretName": "user-gcp-sa"
									}
								}
							]
						}
					],
					"ttlSecondsAfterFinished": 86401
				},
				"status": {
					"finishedAt": "2019-11-30T09:00:17Z",
					"nodes": {
						"sandtradeedcdailyrunhwn89-3-3509392613": {
							"children": [
								"sandtradeedcdailyrunhwn89-3-3509392613-884332321"
							],
							"displayName": "sandtradeedcdailyrunhwn89-3-3509392613",
							"finishedAt": "2019-11-30T09:00:17Z",
							"id": "sandtradeedcdailyrunhwn89-3-3509392613",
							"inputs": {
								"parameters": [
									{
										"name": "bucket",
										"value": "sandtrade_production"
									},
									{
										"name": "ensamble-output",
										"value": "output.csv"
									},
									{
										"name": "forecast-features-blob-name",
										"value": "forecast.csv"
									},
									{
										"name": "preprocessing-blob-name",
										"value": "data/preprocessing/input.csv"
									},
									{
										"name": "prophet-forecast-blob-name",
										"value": "forecast_prophet.csv"
									},
									{
										"name": "r-forecast-blob-name",
										"value": "forecast_r.csv"
									},
									{
										"name": "train-features-blob-name",
										"value": "train.csv"
									}
								]
							},
							"name": "sandtradeedcdailyrunhwn89-3-3509392613",
							"phase": "Failed",
							"startedAt": "2019-11-30T09:00:01Z",
							"templateName": "sandtrade-demand-forecasting",
							"type": "DAG"
						},
						"sandtradeedcdailyrunhwn89-3-3509392613-884332321": {
							"boundaryID": "sandtradeedcdailyrunhwn89-3-3509392613",
							"displayName": "preprocess",
							"finishedAt": "2019-11-30T09:00:16Z",
							"id": "sandtradeedcdailyrunhwn89-3-3509392613-884332321",
							"inputs": {
								"parameters": [
									{
										"name": "bucket",
										"value": "sandtrade_production"
									},
									{
										"name": "preprocessing-blob-name",
										"value": "data/preprocessing/input.csv"
									}
								]
							},
							"message": "failed with exit code 1",
							"name": "sandtradeedcdailyrunhwn89-3-3509392613.preprocess",
							"phase": "Failed",
							"startedAt": "2019-11-30T09:00:02Z",
							"templateName": "preprocess",
							"type": "Pod"
						}
					},
					"phase": "Failed",
					"startedAt": "2019-11-30T09:00:01Z"
				}
			}
		}
	}` // Example of workflow// Example of workflow

	response, err := Mutate([]byte(rawJSON))
	if err != nil {
		t.Errorf("failed to mutate AdmissionRequest %s with error %s", string(response), err)
	}

	r := v1beta1.AdmissionReview{}
	err = json.Unmarshal(response, &r)
	assert.NoError(t, err, "failed to unmarshal with error %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"add","path":"/spec/ttlSecondsAfterFinished","value":36000}]`, string(rr.Patch))
	assert.Contains(t, rr.AuditAnnotations, "mutateme")
}
