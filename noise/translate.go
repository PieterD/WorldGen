package noise

type Translate struct {
	tx, ty float64
	iw, ih float64
	ow, oh float64
}

func NewTranslate (w, h float64) (t *Translate) {
	t = new(Translate)
	t.tx = 0
	t.ty = 0
	t.iw = w
	t.ih = h
	t.ow = w
	t.oh = h
	return
}
func (t *Translate) Resize (w, h float64) {
	t.ow = w
	t.oh = h
}

func (t *Translate) Into (x, y float64) (nx, ny float64) {
	nx = x*(t.iw/t.ow)+t.tx
	ny = y*(t.ih/t.oh)+t.ty
	return
}

func (t *Translate) ZoomBox (x0, y0, x1, y1 float64) {
	if x1 < x0 {
		tmp := x1
		x1 = x0
		x0 = tmp
	}
	if y1 < y0 {
		tmp := y1
		y1 = y0
		y0 = tmp
	}
	if y0 < y1 && x0 < x1 {
		t.Zoom(x0, y0, x1-x0, y1-y0)
	}
}

func (t *Translate) Zoom (x, y, w, h float64) {
	minx, miny := t.Into(x, y)
	maxx, maxy := t.Into(x+w, y+h)
	t.tx = minx
	t.ty = miny
	t.iw = maxx-minx
	t.ih = maxy-miny
}


