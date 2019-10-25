require "./util"
require "./text_to_hash"

class BirthdayParadox
  property dir = "."
  property langs_to_run = [] of String
  property results = [] of Tuple(String, TextHash)

  def initialize
    if ARGV.size > 0 && Dir.exists?(ARGV.first)
      @dir = ARGV.first
    end
    @text_to_hash = TextToHash.new
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
        Util.println("#{lang} => Running...")

        image_id = `docker build --no-cache --quiet .`
        output = `docker run --rm #{image_id}`
        @text_to_hash.text = output
        result = @text_to_hash.convert

        Util.println("#{lang} => #{result}")

        @results << {lang, result}
      end
    end
  end
end

BirthdayParadox.new.execute
