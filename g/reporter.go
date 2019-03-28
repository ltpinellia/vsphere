package g

// ReportVCStatus report vc status
func ReportVCStatus(vc *VsphereConfig, hostname ...string) {
	if Config().Heartbeat.Enabled && Config().Heartbeat.Addr != "" {
		go reportVCStatus(vc, hostname...)
	}
}

func reportVCStatus(vc *VsphereConfig, hostname ...string) {
	var ip string
	var plugin string
	var sHostName string
	if vc.Extend {
		plugin = VERSION
	} else {
		plugin = "plugin not enabled"
	}
	if hostname != nil {
		sHostName = hostname[0]
		ip = hostname[0]
	} else {
		sHostName = vc.Hostname
		ip = InitVCIP(vc)
	}

	req := AgentReportRequest{
		Hostname:      sHostName,
		IP:            ip,
		AgentVersion:  VERSION,
		PluginVersion: plugin,
	}

	var resp SimpleRPCResponse
	err := HbsClient.Call("Agent.ReportStatus", req, &resp)
	if err != nil || resp.Code != 0 {
		Log.Warnln("[reporter.go] call Agent.ReportStatus fail:", err, "Request:", req, "Response:", resp)
	}
}
