ITERATIONS = "iterations"
SAMPLE_SIZE = "sample-size"
PERCENT = "percent"
SECONDS = "seconds"

alias TextHash = Hash(String, String | Nil)

class TextToHash
  def self.unsafe(text : String) : TextHash
    iterations_regex = /#{ITERATIONS}[:]\s*(\d+)/
    sample_size_regex = /#{SAMPLE_SIZE.gsub("-", "[-]")}[:]\s*(\d+)/
    percent_regex = /#{PERCENT}[:]\s*([^\n]+)/
    seconds_regex = /#{SECONDS}[:]\s*([^\n]+)/
    hash = TextHash.new

    pairs = [
      {ITERATIONS, iterations_regex},
      {SAMPLE_SIZE, sample_size_regex},
      {PERCENT, percent_regex},
      {SECONDS, seconds_regex}
    ]

    text_hash = pairs.reduce(hash) do |acc, (name, regex)|
      match = text.match(regex)
      value = match ? match.captures.first : nil

      acc[name] = value
      acc
    end

    text_hash
  end
end

class BirthdayParadox
  property dir : String = "."
  property langs_to_run : Array(String) = [] of String
  property results : Array(Tuple(String, TextHash)) = [] of Tuple(String, TextHash)

  def initialize
    if ARGV.size > 0 && Dir.exists?(ARGV.first)
      @dir = ARGV.first
    end
  end

  def execute
    get_langs_to_run
    @langs_to_run.sort!
    execute_docker_for_langs

    pp @results
  end

  private def get_langs_to_run
    Dir.cd(@dir) do
      Dir.new(".").each_child do |dir|
        if File.file?("#{dir}/Dockerfile")
          @langs_to_run << dir
        end
      end
    end
  end

  private def execute_docker_for_langs
    langs = @langs_to_run.first(2)
    langs.each do |lang|
      Dir.cd("#{@dir}/#{lang}") do
        puts "#{lang} => Running..." if ENV["DEV"]?
        image_id = `docker build --no-cache --quiet .`
        output = `docker run --rm #{image_id}`
        result = TextToHash.unsafe(output)
        puts "#{lang} => #{result}" if ENV["DEV"]?
        @results << {lang, result}
      end
    end
  end
end

BirthdayParadox.new.execute
