:exclamation: **Work in Progress** :exclamation:

Pix
===

Pix is an image server able to apply image transformations (resize, crop, filters, nude detection and so on).

How it works
------------

Pix is a server designed to be an image repository. First you upload any picture, then you can request it using its unique identifier. You can apply transformations by providing a set of parameters in the URLs and voilà Pix give you exactly what you want.

Options
-------

* support for BMP, GIF, JPEG, PNG :heavy_check_mark:
* resize :heavy_check_mark:
* crop :heavy_check_mark:
* smart crop :heavy_check_mark:
* flip vertical / horizontal :heavy_check_mark:
* extension conversion :heavy_check_mark:
* command to delete cached images :heavy_check_mark:
* cache images :heavy_check_mark:
* nude detection
* JPEG interlace
* properly handle animated GIF
* custom URLs: create virtual hosts and associate a set of rules, like: `/user-pictures/file1.jpg`e

URLs
----

| Verb | Path | Description |
| ---- | ---- | ----------- |
| POST | /    | Upload images |
| GET  | /`{image_name}`.`{extension}`?`option1=value1...&optionN=valueN` | Retrieve images and optionally apply effects and transformations |

Libraries
---------

Here is the list of third-party libraries used in this project. I want to thank each author big time for their time and great work :+1:

- **[bild](https://github.com/anthonynsimon/bild)**: image manipulation
- **[mux](https://github.com/gorilla/mux)**: URL router
- **[godotenv](https://github.com/joho/godotenv)**: .env file reader for Go
- **[smartcrop](https://github.com/muesli/smartcrop)**: good image crops
- **[ksuid](https://github.com/segmentio/ksuid)**: K-Sortable Globally Unique IDs
