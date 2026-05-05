package prompts

const (
	// SleepCycleEvaluateMemory is the prompt used to evaluate if a memory is useful.
	SleepCycleEvaluateMemory = `Evaluate this memory: "%s". Is it a critical long-term fact/preference, or is it obsolete/useless data? Respond with a JSON object containing "action": "direct_response" and "response": "KEEP" or "response": "DISCARD".`

	// SleepCycleConsolidateLogs is the prompt used to extract facts and preferences from chat logs.
	SleepCycleConsolidateLogs = `Extract all new factual information, user preferences, and system states from this transcript. Output as a JSON array of concise statements. Do not wrap the JSON array in any markdown, just output the JSON string array.

Transcript:
%s`
)
