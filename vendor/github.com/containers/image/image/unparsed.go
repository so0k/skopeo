package image

import "github.com/containers/image/types"

// UnparsedImage implements types.UnparsedImage .
type UnparsedImage struct {
	src            types.ImageSource
	cachedManifest []byte // A private cache for Manifest(); nil if not yet known.
	// A private cache for Manifest(), may be the empty string if guessing failed.
	// Valid iff cachedManifest is not nil.
	cachedManifestMIMEType string
	cachedSignatures       [][]byte // A private cache for Signatures(); nil if not yet known.
}

// UnparsedFromSource returns a types.UnparsedImage implementation for source.
// The caller must call .Close() on the returned UnparsedImage.
//
// UnparsedFromSource “takes ownership” of the input ImageSource and will call src.Close()
// when the image is closed.  (This does not prevent callers from using both the
// UnparsedImage and ImageSource objects simultaneously, but it means that they only need to
// keep a reference to the UnparsedImage.)
func UnparsedFromSource(src types.ImageSource) *UnparsedImage {
	return &UnparsedImage{src: src}
}

// Reference returns the reference used to set up this source, _as specified by the user_
// (not as the image itself, or its underlying storage, claims).  This can be used e.g. to determine which public keys are trusted for this image.
func (i *UnparsedImage) Reference() types.ImageReference {
	return i.src.Reference()
}

// Close removes resources associated with an initialized UnparsedImage, if any.
func (i *UnparsedImage) Close() {
	i.src.Close()
}

// Manifest is like ImageSource.GetManifest, but the result is cached; it is OK to call this however often you need.
func (i *UnparsedImage) Manifest() ([]byte, string, error) {
	if i.cachedManifest == nil {
		m, mt, err := i.src.GetManifest()
		if err != nil {
			return nil, "", err
		}
		i.cachedManifest = m
		i.cachedManifestMIMEType = mt
	}
	return i.cachedManifest, i.cachedManifestMIMEType, nil
}

// Signatures is like ImageSource.GetSignatures, but the result is cached; it is OK to call this however often you need.
func (i *UnparsedImage) Signatures() ([][]byte, error) {
	if i.cachedSignatures == nil {
		sigs, err := i.src.GetSignatures()
		if err != nil {
			return nil, err
		}
		i.cachedSignatures = sigs
	}
	return i.cachedSignatures, nil
}
