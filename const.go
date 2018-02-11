package srfax

const (
	// SRFax specific action verbs. Every POST request will use one of the following:
	actionQueueFax           = "Queue_Fax"
	actionGetFaxStatus       = "Get_FaxStatus"
	actionGetMulFaxStatus    = "Get_MultiFaxStatus"
	actionGetFaxInbox        = "Get_Fax_Inbox"
	actionGetFaxOutbox       = "Get_Fax_Outbox"
	actionForwardFax         = "Forward_Fax"
	actionRetrieveFax        = "Retrieve_Fax"
	actionUpdateViewedStatus = "Update_Viewed_Status"
	actionDeleteFax          = "Delete_Fax"
	actionStopFax            = "Stop_Fax"
	actionGetFaxUsage        = "Get_Fax_Usage"
)

const (
	// direction, "sDirection"
	inbound  = "IN"
	outbound = "OUT"
)

const (
	// fax type, "sFaxType"
	broadcast = "BROADCAST"
	single    = "SINGLE"
)

const (
	// "sMarkasViewed"
	yes = "Y" // READ
	no  = "N" // UNREAD
)

const (
	// "sFaxFormat"
	pdf = "PDF"
	tif = "TIF"
)
