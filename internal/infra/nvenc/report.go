package nvenc

func BuildReport(requested, effective string, available bool) string {
	if requested == "auto" || requested == "" {
		if available {
			return "encoder:auto->nvenc"
		}
		return "encoder:auto->cpu"
	}
	return "encoder:" + requested + "->" + effective
}
