# gopushbullet #
[![GoDoc](https://godoc.org/github.com/kariudo/gopushbullet?status.svg)](https://godoc.org/github.com/kariudo/gopushbullet)
[![Build Status](https://travis-ci.org/kariudo/gopushbullet.svg?branch=master)](https://travis-ci.org/kariudo/gopushbullet)
[![Coverage Status](https://coveralls.io/repos/kariudo/gopushbullet/badge.svg)](https://coveralls.io/r/kariudo/gopushbullet)

A complete go package for interacting with the fantastic Pushbullet service.

## Status
Many major features complete, see notes below. Additional tests need to be written and a couple less common features. Also need some example code for the test file.

## Features

### Users
* Get User
* Set User preferences

### Pushes
* Send Pushes
 * Note
 * Link
 * Address
 * Checklist
 * File
   * File Uploads
* Delete a push
* Get push history
* Dismiss push

### Devices
* Get Devices

### Contacts
* Get Contacts
* Create Contacts
* Update Contact
* Delete Contact

### Channels
* Subscribe
* Unsubscribe
* Get channel info

## Todo
* Web Sockets
* OAuth account access
* Update a push (dismiss & update list items)
