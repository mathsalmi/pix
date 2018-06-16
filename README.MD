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
* nude detection
* interlace
* extension conversion
* custom URLs: create virtual hosts and associate a set of rules, like: `/user-pictures/file1.jpg`
* work like a pipe (effects/transformations are applied from left to right order)

URLs
----

| Verb | Path | Description |
| ---- | ---- | ----------- |
| POST | /    | Upload images |
| GET  | /`{image_name}`.`{extension}`?`option1=value1...&optionN=valueN` | Retrieve images and optionally apply effects and transformations |
