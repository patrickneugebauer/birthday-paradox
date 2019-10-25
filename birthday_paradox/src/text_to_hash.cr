ITERATIONS = "iterations"
SAMPLE_SIZE = "sample-size"
PERCENT = "percent"
SECONDS = "seconds"

alias TextHash = Hash(String, String?)

class TextToHash
  property text

  def initialize(@text : String = "")
  end

  def convert : TextHash
    hash = TextHash.new

    hash[ITERATIONS] = match_and_capture(/#{ITERATIONS}[:]\s*(\d+)/)
    hash[SAMPLE_SIZE] = match_and_capture(/#{SAMPLE_SIZE.gsub("-", "[-]")}[:]\s*(\d+)/)
    hash[PERCENT] = match_and_capture(/#{PERCENT}[:]\s*([^\n]+)/)
    hash[SECONDS] = match_and_capture(/#{SECONDS}[:]\s*([^\n]+)/)

    hash
  end

  private def match_and_capture(r : Regex) : String?
    match = @text.match(r)
    match ? match.captures.first : nil
  end
end
