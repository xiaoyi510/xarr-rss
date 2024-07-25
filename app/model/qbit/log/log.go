package log

type ApiLogMainReq struct {
	Normal        bool `json:"normal,omitempty"`        //Include Normal messages (default: true)
	Info          bool `json:"info,omitempty"`          //Include Info messages (default: true)
	Warning       bool `json:"warning,omitempty"`       //Include Warning messages (default: true)
	Critical      bool `json:"critical,omitempty"`      //Include Critical messages (default: true)
	Last_known_id int  `json:"last_Known_Id,omitempty"` //Exclude messages with "message id" <= Last_known_id (default: -1)
}
type ApiLogMainRes struct {
	Id        int    `json:"id,omitempty"`        //Id of the message
	Message   string `json:"message,omitempty"`   //Text of the message
	Timestamp int    `json:"timestamp,omitempty"` //Milliseconds since epoch
	Type      int    `json:"type"`                //Type of the message: Log::NORMAL: 1, Log::INFO: 2, Log::WARNING: 4, Log::CRITICAL: 8
}
