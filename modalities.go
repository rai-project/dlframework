package dlframework

type Modality string

const (
	UnknownModality                   Modality = "unknown_modality"
	ImageObjectDetectionModality      Modality = "image_object_detection"
	ImageClassificationModality       Modality = "image_classification"
	ImageInstanceSegmentationModality Modality = "image_instance_segmentation"
	ImageSemanticSegmentationModality Modality = "image_semantic_segmentation"
	ImageEnhancementModality          Modality = "image_enhancement"
	ImageCaptioningModality           Modality = "image_captioning"
)
