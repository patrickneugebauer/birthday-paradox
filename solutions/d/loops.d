import std.stdio;
import std.random;
import std.conv;
import std.datetime.systime;
import core.time;

void main(string[] args) {
  // constants and variables
  immutable auto start = Clock.currTime();
  immutable auto iterations = to!int(args[1]);
  immutable auto sampleSize = 23;
  auto rnd = Random();
  auto count = 0;

  // loop and count
  for (int i = 0; i < iterations; i++) {
    int[365] data;
    for (int s =0; s < sampleSize; s++) {
      immutable auto num = uniform(0, 365, rnd);
      if (data[num] == 1) {
        count++;
        break;
      } else {
        data[num] = 1;
      }
    }
  }

  // calcs
  immutable auto percent = to!float(count) / iterations * 100;
  immutable auto finish = Clock.currTime();
  immutable auto milliseconds = (finish - start).total!"msecs";
  immutable auto seconds = to!float(milliseconds) / 1000;

  // output
  writefln("iterations: %d", iterations);
  writefln("sample-size: %d", sampleSize);
  writefln("percent: %.2f", percent);
  writefln("seconds: %.3f", seconds);
}
