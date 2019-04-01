package dlframework

func (b *BoundingBox) Width() float32 {
	if b == nil {
		return 0
	}
	return b.Xmax - b.Xmin
}

func (b *BoundingBox) Height() float32 {
	if b == nil {
		return 0
	}
	return b.Ymax - b.Ymin
}

func (b *BoundingBox) Area() float32 {
	return b.Width() * b.Height()
}

// returns a list of the form [xmin, ymin, XMAX, YMAX]
func (b *BoundingBox) ToxyXY() []float32 {
	if b == nil {
		return nil
	}
	return []float32{b.Xmin, b.Ymin, b.Xmax, b.Ymax}
}

// returns a list of the form [xmin, ymin, width, height]
func (b *BoundingBox) ToXYWH() []float32 {
	if b == nil {
		return nil
	}
	return []float32{b.Xmin, b.Ymin, b.Width(), b.Height()}
}
