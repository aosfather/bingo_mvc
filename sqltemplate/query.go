package sqltemplate

/**

 */
//分页
type Page struct {
	Size  int
	Index int
	Count int
}

func (this *Page) getStart() int {
	if this.Index > 0 {
		return this.Size * (this.Index - 1)
	}
	return 0
}

func (this *Page) getEnd() int {
	if this.Index > 0 {
		return this.Size * this.Index
	}

	return this.Size - 1
}
