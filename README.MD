Pix
===

Pix is an image server able to apply image transformations (resize, crop, filters, nude detection and so on).

__Pix is still a Work in Progress__

How it works
------------

Pix is a server designed to be an image repository. First you upload any picture, then you can request it using its unique identifier. You can apply transformations by providing a set of parameters in the URLs and voilà Pix give you exactly what you want.

Options
-------

* crop
* resize
* nude detection
* interlace
* extension conversion
* custom URLs: create virtual hosts and associate a set of rules, like: `/user-pictures/file1.jpg`

URLs
----

* Upload images

```
POST /
```

* Retrieving images and optionally apply transformations

```
GET /image_name.extension?option1=value1&option2=value2...&optionN=valueN
```
