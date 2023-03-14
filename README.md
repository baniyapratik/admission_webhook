
```agsl
// AdmissionReview describes an admission review request/response
type AdmissionReview struct {
metav1.TypeMeta `json:"inline"`
// Request describes the attributes for the admission request
Request *AdmissionRequest `json:"request,omitempty"`
// Response describes the attributes for the admission response
Response *AdmissionResponse `json:"response,omitempty"`
}
```