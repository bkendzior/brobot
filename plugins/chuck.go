package brobot
	
type BrobotPlugin interface {
  getTriggers() []string
  trigger([]string)

  getDescription() string
}


func getDescription() string {
  return "A random Chuck Norris fact - specify an @ username for added hilarity"
}
