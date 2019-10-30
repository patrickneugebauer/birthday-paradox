ITERATIONS = "iterations"
SAMPLE_SIZE = "sample-size"
PERCENT = "percent"
SECONDS = "seconds"

alias TextNamedTuple = NamedTuple(
  iterations: String?,
  sample_size: String?,
  percent: String?,
  seconds: String?
)

class TextToNamedTuple
  property text

  def initialize(@text : String = "")
  end

  def convert : TextNamedTuple
    iterations = match_and_capture(/#{ITERATIONS}[:]\s*(\d+)/)
    sample_size = match_and_capture(/#{SAMPLE_SIZE.gsub("-", "[-]")}[:]\s*(\d+)/)
    percent = match_and_capture(/#{PERCENT}[:]\s*([^\n]+)/)
    seconds = match_and_capture(/#{SECONDS}[:]\s*([^\n]+)/)

    TextNamedTuple.new(
      iterations: iterations,
      sample_size: sample_size,
      percent: percent,
      seconds: seconds
    )
  end

  private def match_and_capture(r : Regex) : String?
    match = @text.match(r)
    match ? match.captures.first : nil
  end
end
