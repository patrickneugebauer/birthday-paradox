class Util
  def self.println(text : String) : Void
    puts text if ENV["DEV"]?
  end
end
