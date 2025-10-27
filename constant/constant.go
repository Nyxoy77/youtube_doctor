package constant

const (
	MissingChannelError = "Missing Channel Name"
	Gemini_Flash        = "gemini-2.5-flash"
	Channel             = "channel"
)

const (
	Prompt = `Analyze these YouTube video titles and provide improvements. 
    Return the response in this exact JSON format for each title:
    {
        "llmResponse": [
            {
                "prevTitle": "original title",
                "newTitle": "improved title",
                "reason": "explanation for the improvement"
            }
        ]
    }
    
    Only return valid JSON, no additional text.
    Here are the titles to analyze:%s`
)
