package main

// CodeOK is for a 200 response
const CodeOK int = 200

// CodeNoContent is for a 204 response
const CodeNoContent int = 204

// CodeInvalid is for a 400 response
const CodeInvalid int = 400

// CodeUnauth is for a 401 response
const CodeUnauth int = 401

const errInvalidSecret = "Invalid request secret"
const errMissingEvent = "Missing event value"
const errInvalidBody = "Invalid POST body"
const errInvalidHeader = "Invalid request headers"
