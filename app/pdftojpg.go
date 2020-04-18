package app

import "gopkg.in/gographics/imagick.v2/imagick"

const PDFResolution = 300

// ConvertPdfToJpg will take a filename of a pdf file and convert the file into an
// image which will be saved back to the same location. It will save the image as a
// high resolution jpg file with minimal compression.
func ConvertPdfToJpg(pdfName string, pages int) (result [][]byte, err error) {
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

	// Load the image file into imagick
	if err = mw.ReadImage(pdfName); err != nil {
		return
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err = mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
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
	return
}
