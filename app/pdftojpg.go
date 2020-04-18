package app

import "gopkg.in/gographics/imagick.v2/imagick"

const PDFResolution = 300

type image []byte

// convertPdfToJpg will take a PDF as input and convert this into a slice of
// images (1 per input page)
func convertPdfToImage(pdfContent image, pages int) (result []image, err error) {
	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err = mw.SetResolution(PDFResolution, PDFResolution); err != nil {
		return
	}

	// Load the content into imagick
	if err = mw.ReadImageBlob(pdfContent); err != nil {
		return
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err = mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
		return
	}

	// Set any compression (100 = max quality)
	if err = mw.SetCompressionQuality(100); err != nil {
		return
	}

	// Convert into JPG
	if err = mw.SetFormat("png"); err != nil {
		return
	}

	for p := 0; p < pages; p++ {
		// Select each page of pdf
		mw.SetIteratorIndex(p)

		// Get blob
		blob := mw.GetImageBlob()
		result = append(result, blob)
	}

	return result, err
}
