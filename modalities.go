package dlframework

type Modality string

const (
	UnknownModality              Modality = "unknown_modality"
	ImageObjectDetectionModality Modality = "image_object_detection"
	ImageClassificationModality  Modality = "image_classification"
	ImageSegmentationModality    Modality = "image_segmentation"
	ImageEnhancementModality     Modality = "image_enhancement"
)
