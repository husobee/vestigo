package vestigo

import "net/http"

// handler - internal structure for specifying which handlers belong to a particular route
type handler struct {
	Connect        http.HandlerFunc
	Delete         http.HandlerFunc
	Get            http.HandlerFunc
	Head           http.HandlerFunc
	Options        http.HandlerFunc
	Patch          http.HandlerFunc
	Post           http.HandlerFunc
	Put            http.HandlerFunc
	Trace          http.HandlerFunc
	allowedMethods string
}

// CopyTo - Copy the handler to another handler passed in by reference
func (h *handler) CopyTo(v *handler) {
	v.Get = h.Get
	v.Connect = h.Connect
	v.Delete = h.Delete
	v.Get = h.Get
	v.Head = h.Head
	v.Options = h.Options
	v.Patch = h.Patch
	v.Post = h.Post
	v.Put = h.Put
	v.Trace = h.Trace
	v.allowedMethods = h.allowedMethods
}

// addToAllowedMethods - Add a method to the allowed methods for this route
func (h *handler) addToAllowedMethods(method string) {
	if h.allowedMethods == "" {
		h.allowedMethods = method
	} else {
		h.allowedMethods = h.allowedMethods + ", " + method
	}
}

// AddMethodHandler - Add a method/handler pair to the handler structure
func (h *handler) AddMethodHandler(method string, handler http.HandlerFunc) {
	l := len(method)
	firstChar := method[0]
	secondChar := method[1]
	if h != nil {
		if l == 3 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4745 {
				h.addToAllowedMethods(method)
				h.Get = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5055 {
				h.addToAllowedMethods(method)
				h.Put = handler
			}
		} else if l == 4 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x504f {
				h.addToAllowedMethods(method)
				h.Post = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4845 {
				h.addToAllowedMethods(method)
				h.Head = handler
			}
		} else if l == 5 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5452 {
				h.addToAllowedMethods(method)
				h.Trace = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5041 {
				h.addToAllowedMethods(method)
				h.Patch = handler
			}
		} else if l == 6 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4445 {
				h.addToAllowedMethods(method)
				h.Delete = handler
			}
		} else if l == 7 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4f50 {
				h.addToAllowedMethods(method)
				h.Options = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x434f {
				h.addToAllowedMethods(method)
				h.Connect = handler
			}
		}
	}
}

// GetMethodHandler - Get a method/handler pair from the handler structure
func (h *handler) GetMethodHandler(method string) (http.HandlerFunc, string) {
	l := len(method)
	firstChar := method[0]
	secondChar := method[1]
	if l == 3 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4745 {
			return h.Get, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5055 {
			return h.Put, h.allowedMethods
		}
	} else if l == 4 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x504f {
			return h.Post, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4845 {
			return h.Head, h.allowedMethods
		}
	} else if l == 5 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5452 {
			return h.Trace, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5041 {
			return h.Patch, h.allowedMethods
		}
	} else if l == 6 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4445 {
			return h.Delete, h.allowedMethods
		}
	} else if l == 7 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4f50 {
			return h.Options, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x434f {
			return h.Connect, h.allowedMethods
		}
	}
	return nil, h.allowedMethods
}
