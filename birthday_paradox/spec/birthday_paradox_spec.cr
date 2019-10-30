require "./spec_helper"

describe TextToNamedTuple do
  it "should not find any text" do
    instance = TextToNamedTuple.new("hello world")
    hash = instance.convert

    hash[ITERATIONS].should eq(nil)
    hash[SAMPLE_SIZE].should eq(nil)
    hash[PERCENT].should eq(nil)
    hash[SECONDS].should eq(nil)
  end

  it "should find some text" do
    instance = TextToNamedTuple.new(<<-TEXT)
      iterations: 123
      percent: 78.00
    TEXT
    hash = instance.convert

    hash[ITERATIONS].should eq("123")
    hash[SAMPLE_SIZE].should eq(nil)
    hash[PERCENT].should eq("78.00")
    hash[SECONDS].should eq(nil)
  end

  it "should find correct text" do
    instance = TextToNamedTuple.new(<<-TEXT)
      iterations: 123
      sample-size: 456
      percent: 78.00
      seconds: 0.89
    TEXT
    hash = instance.convert

    hash[ITERATIONS].should eq("123")
    hash[SAMPLE_SIZE].should eq("456")
    hash[PERCENT].should eq("78.00")
    hash[SECONDS].should eq("0.89")
  end

  it "should find correct text excluding incorrect text" do
    instance = TextToNamedTuple.new(<<-TEXT)
      kjldsjflkasdj;flkajsd;flkjsad;lfkjas;ldkfjasld
      iterations: 123
      sample-size: 456
      kjsdofkjsdkfjlskdjflkasjdflkjsa;d
      percent: 78.00
      seconds: 0.89
      oskdfjlkjweofkjwekfj
    TEXT
    hash = instance.convert

    hash[ITERATIONS].should eq("123")
    hash[SAMPLE_SIZE].should eq("456")
    hash[PERCENT].should eq("78.00")
    hash[SECONDS].should eq("0.89")
  end
end
