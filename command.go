package goagi

// Answer executes AGI command "ANSWER"
// Answers channel if not already in answer state.
func (agi *AGI) Answer() (bool, error) {
	resp, err := agi.execute("ANSWER")
	if err != nil {
		return false, err
	}
	ok := resp.code == 200 && resp.result == "0"
	return ok, nil
}
