class Notes
{
  // should sampleSize be get-only, readonly, or const
  // const - never change - compile time constant
  // readonly - change at compile-time
  // get-only - change at runtime

  // property vs field vs local
  // property is an abstraction over fields
  // fields should be used if data is private
  // locals are local variables
  // properties should be used if data is public
  // private static readonly int sampleSize = 23;
  // faster to read from the class and pass in to function
  // than to read from the class over and over in the function...

  // js const can be simulated with readonly ref
  // need to look up immutability in c#

  // private static Random rnd = new Random();
  // sharing random makes the app fast - but it appears to be skipping
  // time is used to seed random, so to ensure each thread is distinctly different, we should pass in some var

  static void Main(string[] args)
  { }
}
