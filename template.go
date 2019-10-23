package http_render

type FuncMap map[string]interface{}

func (this FuncMap) Merge(m ...map[string]interface{}) FuncMap {
	if this == nil {
		this = m[0]
		m = m[1:]
	} else {
		n := FuncMap{}
		for k, v := range this {
			n[k] = v
		}
		this = n
	}
	for _, m := range m {
		for k, v := range m {
			this[k] = v
		}
	}
	return this
}

func (this *FuncMap) Update(m ...map[string]interface{}) FuncMap {
	if *this == nil {
		*this = m[0]
		m = m[1:]
	}
	for _, m := range m {
		for k, v := range m {
			(*this)[k] = v
		}
	}
	return *this
}
