// Copyright 2017 EcoSystem Software LLP

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utilities

import (
	"image/jpeg"
	"os"
	"path"
	"strconv"

	"github.com/nfnt/resize"
)

//Deal with images without width specified and other formats TODO

//CreateImage takes an image name and width parameter and creates a resized image
//if one does not already exist on disk.  Either way, it returns the full path to the image ready to be served
func GetImage(image string, width string) (string, error) {

	//Cache the public image path
	p := "public/images_resized"

	//Create the composite target filename
	targetImageFileName := path.Join(p, width+"w", path.Base(image)) //Even if a path to the image is specified, this is dropped for caching

	//Check to see if there is a cached version of the image
	if _, err := os.Stat(targetImageFileName); os.IsNotExist(err) {

		//If there is no cached version of the image, then check to see if there is a source image
		if _, err := os.Stat(path.Join("public/images_source/", path.Base(image))); os.IsNotExist(err) {

			//If there is no source image with that name, use the placeholder image
			//Check to see if there is a cached resized placeholder image
			targetImageFileName = path.Join(p, width+"w", "image-not-found.jpg")
			if _, err := os.Stat(targetImageFileName); os.IsNotExist(err) {
				//If there isn't a cached placeholder at that size then make one
				err := makeImage(targetImageFileName, width)
				if err != nil {
					return "", err
				}
			}
		} else {
			//If there is a source image, then proceed to make a resized image and cache it
			err := makeImage(targetImageFileName, width)
			if err != nil {
				return "", err
			}
		}

	}

	return targetImageFileName, nil

}

func makeImage(targetImageFileName string, width string) error {

	//Cache the public image path
	p := "public/images_resized"

	reader, err := os.Open(path.Join("public/images_source/", path.Base(targetImageFileName)))
	defer reader.Close()
	if err != nil {
		//If reading the image fails,
		return err
	}
	// Try to decode jpeg into image.Image
	img, err := jpeg.Decode(reader)
	if err != nil {
		return err
	}
	// resize
	w64, err := strconv.ParseUint(width, 10, 64)
	w := uint(w64)
	m := resize.Resize(w, 0, img, resize.Lanczos3)

	// save
	os.Mkdir(path.Join(p, width+"w"), 0777)
	out, err := os.Create(targetImageFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return nil
}
