# HKO Radar image Generator
Generate the HKO radar images into gif

This repository contains golang code for downloading JPEG images from the Hong Kong Observatory's website, converting them into a GIF file, and deleting temporary files.

## Functionality

The code performs the following tasks:

1. Downloads the latest 10 JPEG images from the Hong Kong Observatory's website based on the current time. (server time in UTC, there's `+ 8*time.Hour` in the code)
2. Converts the downloaded images into a GIF file with a delay between frames.
3. Deletes the temporary image files.

## Usage

1. Clone the repository:

   ```bash
   git clone git@github.com:macma/hko-radar-gif-generator.git
   cd hko-radar-gif-generator
   go mod init hko-radar-gif-generator
   go mod tidy
   go run main.go
   ```

## Example output
Below is the gif generated:

![generated gif](https://github.com/macma/hko-radar-gif-generator/blob/main/animation.gif?raw=true)

It's same as the animation for each picture play back as in [HKO playback](https://www.hko.gov.hk/tc/wxinfo/radars/radar.htm?pv_mode=playback)

## Image Copyright
The generated GIF file and the JPEG images used to create it are sourced from the Hong Kong Observatory (HKO) under their terms of use and copyright policy. For more information regarding the usage and restrictions, please refer to the [HKO Terms of Use](https://www.hko.gov.hk/en/appweb/applink.htm).
